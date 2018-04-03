package simpleconfig

import "encoding/json"

// Config defines a typical config params
type Config struct {
	Name     string
	Value    int
	ReadOnly bool
}

// Rand populates config struct with random data
func (config *Config) Rand() *Config {
	config.Name = "mypd"
	config.Value = 500
	config.ReadOnly = true
	return config
}

// Marshal defines serialization for the receiver type
func (config *Config) Marshal() ([]byte, error) {
	return json.MarshalIndent(config, "", "  ")
}

// Unmarshal defines deserialization for the receiver type
func (config *Config) Unmarshal(b []byte) error {
	return json.Unmarshal(b, config)
}
