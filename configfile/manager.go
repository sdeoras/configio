package configfile

import (
	"context"
	"io/ioutil"
	"sync"

	"github.com/sdeoras/configio"
	"github.com/sirupsen/logrus"
)

// manager implements several config management interfaces for a
// backend that is a file on the disk
type manager struct {
	mu   sync.Mutex
	ctx  context.Context
	cb   map[string]*configio.Callback
	log  *logrus.Entry
	file string
}

// Init initializes newly instantiated manager
func (m *manager) Init(ctx context.Context) *manager {
	m.cb = make(map[string]*configio.Callback)
	m.ctx = ctx
	m.file = DefaultConfigFile
	m.log = logrus.WithField("manager", "manager")
	return m
}

// SetConfigFile sets location of config file other than the default
func (m *manager) SetConfigFile(fileName string) {
	m.file = fileName
}

// Unmarshal reads config file, unmarshals it into configio.Marshaler
func (m *manager) Unmarshal(config configio.Marshaler) error {
	b, err := ioutil.ReadFile(DefaultConfigFile)
	if err != nil {
		m.log.Error(err)
		return err
	}

	if err := config.Unmarshal(b); err != nil {
		m.log.Error(err)
		return err
	}

	return nil
}

// Marshal serializes input config and writes to config file.
// Furthermore, it runs through registered callbacks
func (m *manager) Marshal(config configio.Marshaler) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	b, err := config.Marshal()
	if err != nil {
		m.log.Error(err)
		return err
	}

	if err := ioutil.WriteFile(DefaultConfigFile, b, 0666); err != nil {
		m.log.Error(err)
		return err
	}

	for name := range m.cb {
		go m.execCallback(name, config)
	}

	return nil
}

// Watch registers a function to watch on config changes and returns a channel on which clients can watch
func (m *manager) Watch(name string, data interface{}, f func(ctx context.Context, data interface{}, err error) <-chan error) <-chan configio.Marshaler {
	cbd := new(configio.Callback)
	cbd.Func = f
	cbd.Chan = make(chan configio.Marshaler)
	cbd.Data = data
	m.cb[name] = cbd
	return cbd.Chan
}

// execCallback executes callback
func (m *manager) execCallback(name string, config configio.Marshaler) {
	cbd := m.cb[name]
	err := cbd.Func(m.ctx, cbd.Data, cbd.Err)
	readConfig, sentConfirmation := false, false
	for {
		select {
		case cbd.Chan <- config:
			m.log.Info(name, " read config")
			readConfig = true
		case cbd.Err = <-err:
			if cbd.Err != nil {
				m.log.Info(name, " executed unsuccessfully")
				go func() {
					select {
					case <-cbd.Func(m.ctx, cbd.Data, cbd.Err):
					case <-m.ctx.Done():
					}
				}()
				delete(m.cb, name)
			} else {
				m.log.Info(name, " executed successfully")
			}
			sentConfirmation = true
		case <-m.ctx.Done():
			m.log.Error(name, "context done")
			return
		}
		if readConfig && sentConfirmation {
			return
		}
	}
}
