package configio

import "context"

// Marshaler defines the behavior of a type that can act as a config parameter
// Any such type should provide serialization methods
type Marshaler interface {
	Marshal() ([]byte, error)
	Unmarshal(b []byte) error
}

// ConfigManager defines an interface to access config params
type ConfigManager interface {
	ConfigReadWriter
	ConfigWatcher
}

// ConfigReadWriter defines an interface to perform read/write on config params
type ConfigReadWriter interface {
	ConfigReader
	ConfigWriter
}

// ConfigReader defines an interface to perform read operation on config params
type ConfigReader interface {
	Get() (Marshaler, error)
}

// ConfigWriter defines an interface to perform write operation on config params
type ConfigWriter interface {
	Set(config Marshaler) error
}

// ConfigWatcher defines an interface to perform a watch on config changes.
// Watch registers a function that gets executed on config changes.
// If an error occurs on function execution, the function is called again with that
// error passed in as an input argument and is removed from the registry
type ConfigWatcher interface {
	Watch(name string, data interface{}, f func(ctx context.Context, data interface{}, err error) <-chan error) <-chan Marshaler
}
