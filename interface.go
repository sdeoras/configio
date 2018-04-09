// configio defines a generic interface for config management
package configio

import (
	"context"
	"io"
)

// Config defines the behavior of a type that can act as a config parameter
// Any such type should provide serialization methods and ability to get a key
type Config interface {
	Key() string
	Marshal() ([]byte, error)
	Unmarshal(b []byte) error
}

// ConfigManager defines an interface to access config params
type ConfigManager interface {
	// Unmarshal unmarshals into marshaler
	Unmarshal(config Config) error
	// Marshal marshals data in marshaler
	Marshal(config Config) error
	// Watch registers a callback function
	Watch(name string, data interface{},
		f func(ctx context.Context, data interface{}, err error) <-chan error) <-chan struct{}
	io.Closer
}

// ConfigReadWriter defines an interface to perform read/write on config params
type ConfigReadWriter interface {
	// Unmarshal unmarshals into marshaler
	Unmarshal(config Config) error
	// Marshal marshals data in marshaler
	Marshal(config Config) error
	io.Closer
}

// ConfigReader defines an interface to perform read operation on config params
type ConfigReader interface {
	// Unmarshal unmarshals into marshaler
	Unmarshal(config Config) error
	io.Closer
}

// ConfigWriter defines an interface to perform write operation on config params
type ConfigWriter interface {
	// Marshal marshals data in marshaler
	Marshal(config Config) error
	io.Closer
}

// ConfigWatcher defines an interface to perform a watch on config changes.
// Watch registers a function that gets executed on config changes.
// If an error occurs on function execution, the function is called again with that
// error passed in as an input argument and is removed from the registry
type ConfigWatcher interface {
	// Watch registers a callback function
	Watch(name string, data interface{},
		f func(ctx context.Context, data interface{}, err error) <-chan error) <-chan struct{}
	io.Closer
}
