package config

import (
	"errors"
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
	ZettelkastenDirectory string     `koanf:"zettelkasten_directory" validate:"requiredWithout:ZettelkastenGitURL"`
	ZettelkastenGitURL    string     `koanf:"zettelkasten_git_url" validate:"requiredWithout:ZettelkastenDirectory" validate:"url/isURL"`
	ZettelkastenGitBranch string     `koanf:"zettelkasten_git_branch"`
	LogLevel              slog.Level `koanf:"log_level"`
	IgnoreFiles           []string   `koanf:"ignore_files"`
}

func LoadConfig() (Config, error) {
	k := koanf.New(".")

	// Set default values
	k.Load(structs.Provider(Config{
		IP:                    "0.0.0.0",
		Port:                  10018,
		LogLevel:              slog.LevelInfo,
		IgnoreFiles:           []string{".git", ".obsidian", ".trash", "README.md"},
		ZettelkastenGitBranch: "main",
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
	if cfg.ZettelkastenGitURL != "" && cfg.ZettelkastenDirectory != "" {
		return Config{}, errors.New("ZettelkastenGitURL and ZettelkastenDirectory cannot be provided together")
	}

	return cfg, nil
}
