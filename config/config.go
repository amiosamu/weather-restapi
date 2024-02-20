package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"path"
)


type (
	Config struct {
		App   `yaml:"app"`
		HTTP  `yaml:"http"`
		Mongo `yaml:"mongo"`
	}
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}
	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}
	Mongo struct {
		Addr     string `yaml:"addr"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	}
)

func NewConfig(confPath string) (*Config, error) {
	cfg := &Config{}
	err := cleanenv.ReadConfig(path.Join("./", confPath), cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}
	err = cleanenv.UpdateEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("error updating config file: %w", err)
	}
	return cfg, nil
}