package config

import "sync"

type Config struct {
	Prompt string
}

var conf *Config
var once sync.Once

func loadConfig() *Config {
	return &Config{
		Prompt: "traSH",
	}
}

func GetConfig() *Config {
	once.Do(func() {
		conf = loadConfig()
	})
	return conf
}
