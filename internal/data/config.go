package data

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

type AppConfig struct {
	Name    string `default:"pds-auth"`
	Version string `default:"0.0.0"`

	Key      []byte
	KeyValue string `mapstructure:"key" env:"APP_ENV"`
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
	if cfg.App.KeyValue == "" {
		return nil, errors.New("missing application key")
	} else {
		var err error
		switch {
		case strings.HasPrefix(cfg.App.KeyValue, "hex:"):
			trimmed := strings.TrimPrefix(cfg.App.KeyValue, "hex:")
			if cfg.App.Key, err = hex.DecodeString(trimmed); err != nil {
				return nil, fmt.Errorf("failed to decode key from hex: %w", err)
			}
		case strings.HasPrefix(cfg.App.KeyValue, "base64:"):
			trimmed := strings.TrimPrefix(cfg.App.KeyValue, "base64:")
			if cfg.App.Key, err = base64.StdEncoding.DecodeString(trimmed); err != nil {
				return nil, fmt.Errorf("failed to decode key from base64: %w", err)
			}
		}
	}

	return cfg, nil
}
