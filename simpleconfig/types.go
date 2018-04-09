package simpleconfig

import (
	"encoding/json"

	"github.com/google/uuid"
)

// Config defines a typical config params
type Config struct {
	key      string
	Name     string
	Value    int
	ReadOnly bool
}

// Rand populates config struct with random data
func (config *Config) Rand() *Config {
	config.Name = "mypd"
	config.Value = 500
	config.ReadOnly = true
	config.key = uuid.New().String()
	return config
}

func (config *Config) Init(key string) *Config {
	config.key = key
	return config
}

func (config *Config) Key() string {
	return config.key
}

// Marshal defines serialization for the receiver type
func (config *Config) Marshal() ([]byte, error) {
	return json.MarshalIndent(config, "", "  ")
}

// Unmarshal defines deserialization for the receiver type
func (config *Config) Unmarshal(b []byte) error {
	return json.Unmarshal(b, config)
}
