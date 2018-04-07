package configfile

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sdeoras/configio"
)

// NewManager returns instance of ConfigManager interface.
// ConfigManager should be used by clients requiring full config management
// features
func NewManager(ctx context.Context, options ...interface{}) (configio.ConfigManager, error) {
	return newManager(ctx, options...)
}

// NewReader returns instance of ConfigReader interface.
// ConfigReader should be used by clients who just wish to read config and/or
// for clients who should be given read permissions only
func NewReader(ctx context.Context, options ...interface{}) (configio.ConfigReader, error) {
	return newManager(ctx, options...)
}

// NewWriter returns instance of ConfigWriter interface.
// ConfigWriter should be used by clients who just wish to write config and/or
// for clients who should be given write permissions only
func NewWriter(ctx context.Context, options ...interface{}) (configio.ConfigWriter, error) {
	return newManager(ctx, options...)
}

// NewReadWriter returns instance of ConfigReadWriter interface.
// ConfigReadWriter should be used by clients who wish to read and write config
// and are not going to perform a watch
func NewReadWriter(ctx context.Context, options ...interface{}) (configio.ConfigReadWriter, error) {
	return newManager(ctx, options...)
}

// NewWatcher returns an instance of ConfigWatcher interface.
// ConfigWatcher should be used by clients who only wish to watch config changes
func NewWatcher(ctx context.Context, options ...interface{}) (configio.ConfigWatcher, error) {
	return newManager(ctx, options...)
}

// newManager returns instance of manager struct
func newManager(ctx context.Context, options ...interface{}) (*manager, error) {
	m := new(manager).Init(ctx)

	opt := make(map[string]interface{})
	for i, option := range options {
		if i%2 != 0 {
			continue
		}
		key, ok := option.(string)
		if !ok {
			return nil, fmt.Errorf("options need to be in the format key, value. key is a string")
		}

		if i < len(options)-1 {
			opt[strings.ToLower(key)] = options[i+1]
		}
	}

	if option, present := opt[fileOption]; present {
		if file, ok := option.(string); !ok {
			return nil, fmt.Errorf("option value for file should be a string")
		} else {
			if err := m.setConfigFile(file); err != nil {
				return nil, err
			}
		}
	} else {
		home := os.Getenv("HOME")
		file := filepath.Join(home, ".config", DefaultConfigDir, DefaultConfigFile)
		if err := m.setConfigFile(file); err != nil {
			return nil, err
		}
	}

	return m, nil
}
