package config

import "time"

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Postgres PostgresConfig `yaml:"postgres"`
	Log      LogConfig      `yaml:"log"`
}

type LogConfig struct {
	Level string `yaml:"level"`
}

type ServerConfig struct {
	Address      string        `yaml:"address"`
	IdleTimeout  time.Duration `yaml:"idleTimeout"`
	ReadTimeout  time.Duration `yaml:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`
}

type PostgresConfig struct {
	Name         string        `yaml:"name"`
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	User         string        `yaml:"user"`
	Password     string        `yaml:"password"`
	SSLMode      string        `yaml:"sslmode"`
	MaxOpenConns int           `yaml:"maxOpenConns"`
	MaxIdleConns int           `yaml:"maxIdleConns"`
	MaxIdleTime  time.Duration `yaml:"maxIdleTime"`
}
