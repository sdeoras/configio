package simpleconfig

import "encoding/json"

// Config defines config params
type Config struct {
	PDName   string
	PDSize   string
	ReadOnly bool
}

// Rand populates config struct with random data
func (config *Config) Rand() *Config {
	config.PDName = "mypd"
	config.PDSize = "500Gi"
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
