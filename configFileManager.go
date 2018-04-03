package configio

import (
	"context"
	"io/ioutil"
	"sync"

	"github.com/sirupsen/logrus"
)

const (
	configFile = "/tmp/config.json"
)

// configFileManager implements several config management interfaces for a
// backend that is a file on the disk
type configFileManager struct {
	mu  sync.Mutex
	ctx context.Context
	cb  map[string]*callbackData
	log *logrus.Entry
}

// Init initializes newly instantiated configFileManager
func (m *configFileManager) Init(ctx context.Context) *configFileManager {
	m.cb = make(map[string]*callbackData)
	m.ctx = ctx
	m.log = logrus.WithField("manager", "configFileManager")
	return m
}

// Get reads config file, unmarshals it and returns as Marshaler interface
func (m *configFileManager) Get() (Marshaler, error) {
	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		m.log.Error(err)
		return nil, err
	}

	config := new(Config)
	if err := config.Unmarshal(b); err != nil {
		m.log.Error(err)
		return nil, err
	}

	return config, nil
}

// Set serializes input config and writes to config file.
// Furthermore, it runs through registered callbacks
func (m *configFileManager) Set(config Marshaler) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	b, err := config.Marshal()
	if err != nil {
		m.log.Error(err)
		return err
	}

	if err := ioutil.WriteFile(configFile, b, 0666); err != nil {
		m.log.Error(err)
		return err
	}

	for name := range m.cb {
		go m.execCallback(name, config)
	}

	return nil
}

// Watch registers a function to watch on config changes and returns a channel on which clients can watch
func (m *configFileManager) Watch(name string, data interface{}, f func(ctx context.Context, data interface{}, err error) <-chan error) <-chan Marshaler {
	cbd := new(callbackData)
	cbd.f = f
	cbd.c = make(chan Marshaler)
	cbd.data = data
	m.cb[name] = cbd
	return cbd.c
}

// execCallback executes callback
func (m *configFileManager) execCallback(name string, config Marshaler) {
	cbd := m.cb[name]
	err := cbd.f(m.ctx, cbd.data, cbd.err)
	readConfig, sentConfirmation := false, false
	for {
		select {
		case cbd.c <- config:
			m.log.Info(name, " read config")
			readConfig = true
		case cbd.err = <-err:
			if cbd.err != nil {
				m.log.Info(name, " executed unsuccessfully")
				go func() {
					select {
					case <-cbd.f(m.ctx, cbd.data, cbd.err):
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
