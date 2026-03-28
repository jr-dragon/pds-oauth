package data

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

type AppConfig struct {
	Name    string `default:"pds-auth"`
	Version string `default:"0.0.0"`

	URL      *url.URL `mapstructure:"-"`
	URLValue string   `mapstructure:"url" env:"APP_URL"`

	Key      []byte `mapstructure:"-"`
	KeyValue string `mapstructure:"key" env:"APP_KEY"`
}

func (c AppConfig) IsLocal() bool {
	u, err := url.ParseRequestURI(c.URLValue)
	if err != nil {
		// handle invalid URL format
		return false
	}

	host := u.Host
	// Remove port if present (e.g., "localhost:8080" becomes "localhost")
	if strings.Contains(host, ":") {
		h, _, err := net.SplitHostPort(host)
		if err != nil {
			// handle error in splitting host and port
			return false
		}
		host = h
	}

	// Check if the hostname is a common loopback identifier
	switch host {
	case "localhost", "127.0.0.1", "::1", "":
		// An empty host might indicate a file:/// style URL or similar, which is local
		return true
	}

	// You can also check if the resolved IP is a loopback address
	// This is more complex and usually not necessary if checking the hostname string
	return false
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
	if u, err := url.Parse(cfg.App.URLValue); err != nil {
		return nil, err
	} else {
		cfg.App.URL = u
	}

	return cfg, nil
}
