package config

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var DefaultConfig = Config{
	Log: Log{
		Level:            "info",
		AddSource:        false,
		ToFile:           false,
		BaseName:         "main",
		Directory:        "log",
		RotateSize:       "1GB",
		RotateInterval:   "24h",
		RotateAtMidnight: true,
	},
}

type Config struct {
	Log Log `koanf:"log"`
}

type Log struct {
	Level            string `koanf:"level"`
	AddSource        bool   `koanf:"add_source"`
	ToFile           bool   `koanf:"to_file"`
	BaseName         string `koanf:"base_name"`
	Directory        string `koanf:"directory"`
	RotateSize       string `koanf:"rotate_size"`
	RotateInterval   string `koanf:"rotate_interval"`
	RotateAtMidnight bool   `koanf:"rotate_at_midnight"`
}

func Load(cfgFile string) (Config, error) {
	k := koanf.New(".")
	f := file.Provider(cfgFile)
	c := DefaultConfig

	err := k.Load(f, toml.Parser())
	if err != nil {
		return c, fmt.Errorf("error loading config: %w", err)
	}

	envPrefix := k.String("env.prefix")
	err = k.Load(env.Provider(envPrefix, ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, envPrefix)), "__", ".", -1)
	}), nil)
	if err != nil {
		return c, fmt.Errorf("error loading env: %w", err)
	}

	err = k.Unmarshal("", &c)
	if err != nil {
		return c, fmt.Errorf("error unmarshalling config: %w", err)
	}

	return c, nil
}

func Watch(cfgFile string, levelVar *slog.LevelVar, logger *slog.Logger) error {
	f := file.Provider(cfgFile)
	return f.Watch(func(event interface{}, err error) {
		if err != nil {
			logger.Error("watch error", "err", err)
			return
		}

		logger.Info("config changed. Reloading ...")
		k := koanf.New(".")
		if err := k.Load(f, toml.Parser()); err != nil {
			logger.Error("error loading config", "err", err)
			return
		}
		logger.Info("config", "config", k.Raw())

		err = levelVar.UnmarshalText(k.Bytes("log.level"))
		if err != nil {
			levelVar.Set(slog.LevelInfo)
			logger.Warn("invalid log level, use info instead")
		}
	})
}
