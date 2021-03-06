package configfile

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"encoding/json"

	"github.com/fsnotify/fsnotify"
	"github.com/google/uuid"
	"github.com/sdeoras/configio"
	"github.com/sirupsen/logrus"
)

// manager implements several config management interfaces for a
// backend that is a file on the disk
type manager struct {
	mu          sync.Mutex
	ctx         context.Context
	cb          map[string]*configio.Callback
	log         *logrus.Entry
	file        string
	watcher     *fsnotify.Watcher
	watchCtx    context.Context
	watchCancel context.CancelFunc
	watchClosed bool
}

// Init initializes newly instantiated manager
func (m *manager) Init(ctx context.Context) *manager {
	m.log = logrus.WithField("manager", "configio")
	m.cb = make(map[string]*configio.Callback)
	m.ctx = ctx

	return m
}

// Close closes and performs cleanup if any
func (m *manager) Close() error {
	return m.closeWatch()
}

// close watch
func (m *manager) closeWatch() error {
	if !m.watchClosed {
		// cancel watch context
		m.watchCancel()
		m.watchClosed = true

		return m.watcher.Close()
	} else {
		return nil
	}
}

// setConfigFile sets location of config file other than the default
func (m *manager) setConfigFile(fileName string) error {
	log := m.log.WithField("func", "setConfigFile")

	// mkdir and fstat check
	if err := initIfNotExists(fileName); err != nil {
		log.Error(err)
		return err
	}

	// update file
	m.file = fileName

	// recreate new watcher and context
	if err := m.initWatch(); err != nil {
		return err
	}

	return nil
}

func (m *manager) initWatch() error {
	log := m.log.WithField("func", "initWatch")
	m.watchClosed = false

	// recreate new watcher and context
	if watcher, err := fsnotify.NewWatcher(); err != nil {
		log.Error(err)
		return err
	} else {
		m.watcher = watcher
		m.watchCtx, m.watchCancel = context.WithCancel(m.ctx)
	}

	// start watching
	go m.watch()

	if err := m.watcher.Add(m.file); err != nil {
		log.Error(err)
		// cancel watch context
		m.watchCancel()
		return err
	}

	return nil
}

// Unmarshal reads config file, unmarshals it into configio.Config
func (m *manager) Unmarshal(config configio.Config) error {
	log := m.log.WithField("func", "Unmarshal")

	if len(config.Key()) == 0 {
		return fmt.Errorf("config key is empty")
	}

	data := make(map[string][]byte)

	// read from file and unmarshal into map
	if b, err := ioutil.ReadFile(m.file); err != nil {
		log.Error(err)
		return err
	} else {
		if err := json.Unmarshal(b, &data); err != nil {
			return err
		}
	}

	// find value for key and unmarshal into object
	if b, present := data[config.Key()]; !present {
		log.WithField("key", config.Key()).WithField("file", m.file).Error("no data available")
		return fmt.Errorf("no data available")
	} else {
		if err := config.Unmarshal(b); err != nil {
			log.Error(err)
			return err
		}
	}

	return nil
}

// Marshal serializes input config and writes to config file.
// Furthermore, it runs through registered callbacks
func (m *manager) Marshal(config configio.Config) error {
	log := m.log.WithField("func", "Marshal")
	m.mu.Lock()
	defer m.mu.Unlock()

	data := make(map[string][]byte)

	// read from file and unmarshal into map
	if b, err := ioutil.ReadFile(m.file); err != nil {
		log.Error(err)
		return err
	} else {
		if err := json.Unmarshal(b, &data); err != nil {
			return err
		}
	}

	if b, err := config.Marshal(); err != nil {
		log.Error(err)
		return err
	} else {
		// set config in map
		data[config.Key()] = b
	}

	dir, _ := filepath.Split(m.file)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Error(err)
		return err
	}

	if b, err := json.MarshalIndent(data, "", "  "); err != nil {
		log.Error(err)
		return err
	} else {
		if err := ioutil.WriteFile(m.file, b, 0666); err != nil {
			log.Error(err)
			return err
		}
	}

	return nil
}

// Watch registers a function to watch on config changes and returns a channel on which clients can watch
func (m *manager) Watch(name string, data interface{}, f func(ctx context.Context, data interface{}, err error) <-chan error) <-chan struct{} {
	m.mu.Lock()
	defer m.mu.Unlock()

	cbd := new(configio.Callback)
	cbd.Func = f
	cbd.Chan = make(chan struct{})
	cbd.Data = data
	m.cb[name] = cbd
	return cbd.Chan
}

// watch watches for file changes
func (m *manager) watch() {
	log := m.log.WithField("func", "watch")
	if m.watchClosed {
		return
	}

	log.Info("starting watch")
	for {
		if m.watchClosed {
			break
		}
		select {
		case event := <-m.watcher.Events:
			log.WithField("file", event.Name).Info(event.Op)
			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Rename == fsnotify.Rename {
				m.mu.Lock()
				for name, cbd := range m.cb {
					name, cbd := name, cbd
					go m.execCallback(name, cbd)
				}
				m.mu.Unlock()
			}
			if event.Op&fsnotify.Remove == fsnotify.Remove {
				log.Info("file removed, stopping watch")
				return
			}
		case err := <-m.watcher.Errors:
			log.Error(err)
		case <-m.watchCtx.Done():
			log.Info("context done, stopping watch")
			return
		}
	}

	log.Info("stopping watch")
	return
}

// execCallback executes callback
func (m *manager) execCallback(name string, cbd *configio.Callback) {
	log := m.log.WithField("func", "execCallback")
	log.WithField("callback", name).Info("executing")

	err := cbd.Func(m.ctx, cbd.Data, cbd.Err)
	readConfig, sentConfirmation := false, false
	for {
		select {
		case cbd.Chan <- struct{}{}:
			log.WithField("callback", name).Info("received notification")
			readConfig = true
		case cbd.Err = <-err:
			if cbd.Err != nil {
				log.WithField("callback", name).Info("executed unsuccessfully")
				go func() {
					select {
					case <-cbd.Func(m.ctx, cbd.Data, cbd.Err):
					case <-m.ctx.Done():
					}
				}()
				delete(m.cb, name)
			} else {
				log.WithField("callback", name).Info("executed successfully")
			}
			sentConfirmation = true
		case <-m.ctx.Done():
			log.Error("context done, returning")
			return
		}
		if readConfig && sentConfirmation {
			return
		}
	}
}

func initIfNotExists(fileName string) error {
	// mkdir and touch
	dir, _ := filepath.Split(fileName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	if fileInfo, err := os.Stat(fileName); err != nil {
		if !os.IsNotExist(err) {
			return err
		} else {
			data := make(map[string][]byte)
			data[uuid.New().String()] = []byte{0, 1, 2}
			if b, err := json.MarshalIndent(data, "", "  "); err != nil {
				return err
			} else {
				if err := ioutil.WriteFile(fileName, b, 0666); err != nil {
					return err
				}
			}
		}
	} else {
		if fileInfo.IsDir() {
			return fmt.Errorf("input option filename is a directory")
		}
	}

	return nil
}
