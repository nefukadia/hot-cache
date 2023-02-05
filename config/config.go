package config

type Config struct {
	Listen     string `yaml:"listen"`
	Auth       string `yaml:"auth"`
	Heartbeat  int    `yaml:"heartbeat"`
	ReadBuffer int    `yaml:"read_buffer"`
}
