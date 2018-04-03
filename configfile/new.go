package configfile

import (
	"context"

	"github.com/sdeoras/configio"
)

// NewManager returns instance of ConfigManager interface.
// ConfigManager should be used by clients requiring full config management
// features
func NewManager(ctx context.Context) (configio.ConfigManager, error) {
	return newManager(ctx)
}

// NewReader returns instance of ConfigReader interface.
// ConfigReader should be used by clients who just wish to read config and/or
// for clients who should be given read permissions only
func NewReader(ctx context.Context) (configio.ConfigReader, error) {
	return newManager(ctx)
}

// NewWriter returns instance of ConfigWriter interface.
// ConfigWriter should be used by clients who just wish to write config and/or
// for clients who should be given write permissions only
func NewWriter(ctx context.Context) (configio.ConfigWriter, error) {
	return newManager(ctx)
}

// NewReadWriter returns instance of ConfigReadWriter interface.
// ConfigReadWriter should be used by clients who wish to read and write config
// and are not going to perform a watch
func NewReadWriter(ctx context.Context) (configio.ConfigReadWriter, error) {
	return newManager(ctx)
}

// NewWatcher returns an instance of ConfigWatcher interface.
// ConfigWatcher should be used by clients who only wish to watch config changes
func NewWatcher(ctx context.Context) (configio.ConfigWatcher, error) {
	return newManager(ctx)
}

// newManager returns instance of manager struct
func newManager(ctx context.Context) (*manager, error) {
	return new(manager).Init(ctx)
}
