package config

import (
	"log/slog"
	"strings"

	"github.com/gookit/validate"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/structs"
)

type Config struct {
	IP                    string     `koanf:"ip" validate:"required|ip"`
	Port                  int        `koanf:"port" validate:"required|uint"`
	LogLevel              slog.Level `koanf:"log_level"`
	ZettelkastenDirectory string     `koanf:"zettelkasten_directory" validate:"required"`
}

func LoadConfig() (Config, error) {
	k := koanf.New(".")

	// Set default values
	k.Load(structs.Provider(Config{
		IP:       "0.0.0.0",
		Port:     6969,
		LogLevel: slog.LevelInfo,
	}, "koanf"), nil)

	// Load env variables
	k.Load(env.Provider("", ".", strings.ToLower), nil)

	// Unmarshal into config struct
	var cfg Config
	k.Unmarshal("", &cfg)

	// Validate config
	v := validate.Struct(cfg)
	if !v.Validate() {
		return Config{}, v.Errors
	}

	return cfg, nil
}
