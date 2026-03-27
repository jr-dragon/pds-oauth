package data

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

type AppConfig struct {
	Name    string `default:"pds-auth"`
	Version string `default:"0.0.0"`
}

type ServerConfig struct {
	Addr string `mapstructure:"addr" env:"SERVER_ADDR" default:"8000"`
}

type Config struct {
	App    AppConfig    `default:""`
	Server ServerConfig `default:""`
}

func NewConfig(name, version string, paths ...string) (*Config, error) {
	config.WithOptions(config.ParseEnv, config.ParseTime, config.ParseDefault)
	config.AddDriver(yaml.Driver)

	cfg := &Config{}

	var cfgPaths []string
	for _, p := range paths {
		stat, err := os.Stat(p)
		if err != nil {
			return nil, err
		}

		if !stat.IsDir() {
			cfgPaths = append(cfgPaths, p)
			continue
		}

		filepath.WalkDir(p, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() {
				cfgPaths = append(cfgPaths, path)
			}

			return nil
		})
	}

	if err := config.LoadFiles(cfgPaths...); err != nil {
		return nil, err
	}

	if err := config.BindStruct("", cfg); err != nil {
		return nil, err
	}

	if name != "" {
		cfg.App.Name = name
	}
	if version != "" {
		cfg.App.Version = version
	}

	return cfg, nil
}
