package config

import "time"

type Config struct {
	Server ServerConfig `yaml:"server"`
}

type ServerConfig struct {
	Address      string        `yaml:"address"`
	IdleTimeout  time.Duration `yaml:"idleTimeout"`
	ReadTimeout  time.Duration `yaml:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`
}
