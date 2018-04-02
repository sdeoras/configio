package configio

import "context"

// NewConfigFileManager returns instance of ConfigManager interface.
// ConfigManager should be used by clients requiring full config management
// features
func NewConfigFileManager(ctx context.Context) ConfigManager {
	return newConfigFileManager(ctx)
}

// NewConfigFileReader returns instance of ConfigReader interface.
// ConfigReader should be used by clients who just wish to read config and/or
// for clients who should be given read permissions only
func NewConfigFileReader(ctx context.Context) ConfigReader {
	return newConfigFileManager(ctx)
}

// NewConfigFileWriter returns instance of ConfigWriter interface.
// ConfigWriter should be used by clients who just wish to write config and/or
// for clients who should be given write permissions only
func NewConfigFileWriter(ctx context.Context) ConfigWriter {
	return newConfigFileManager(ctx)
}

// NewConfigFileReadWriter returns instance of ConfigReadWriter interface.
// ConfigReadWriter should be used by clients who wish to read and write config
// and are not going to perform a watch
func NewConfigFileReadWriter(ctx context.Context) ConfigReadWriter {
	return newConfigFileManager(ctx)
}

// newConfigFileManager returns instance of configFileManager struct
func newConfigFileManager(ctx context.Context) *configFileManager {
	return new(configFileManager).Init(ctx)
}
