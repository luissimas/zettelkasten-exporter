package config

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/gookit/validate"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/structs"
)

type Config struct {
	ZettelkastenDirectory    string        `koanf:"zettelkasten_directory" validate:"requiredWithout:ZettelkastenGitURL"`
	ZettelkastenGitURL       string        `koanf:"zettelkasten_git_url" validate:"requiredWithout:ZettelkastenDirectory|url"`
	ZettelkastenGitBranch    string        `koanf:"zettelkasten_git_branch"`
	ZettelkastenGitToken     string        `koanf:"zettelkasten_git_token"`
	LogLevel                 slog.Level    `koanf:"log_level"`
	IgnoreFiles              []string      `koanf:"ignore_files"`
	CollectionInterval       time.Duration `koanf:"collection_interval"`
	CollectHistoricalMetrics bool          `koanf:"collect_historical_metrics"`
	VictoriaMetricsURL       string        `koanf:"victoriametrics_url" validate:"fullUrl"`
	InfluxDBURL              string        `koanf:"influxdb_url" validate:"fullUrl"`
	InfluxDBToken            string        `koanf:"influxdb_token" validate:"requiredWith:InfluxDBURL"`
	InfluxDBOrg              string        `koanf:"influxdb_org" validate:"requiredWith:InfluxDBURL"`
	InfluxDBBucket           string        `koanf:"influxdb_bucket" validate:"requiredWith:InfluxDBURL"`
}

func LoadConfig() (Config, error) {
	k := koanf.New(".")

	// Set default values
	err := k.Load(structs.Provider(Config{
		LogLevel:                 slog.LevelInfo,
		IgnoreFiles:              []string{".git", ".obsidian", ".trash", "README.md"},
		ZettelkastenGitBranch:    "main",
		CollectionInterval:       time.Minute * 5,
		CollectHistoricalMetrics: true,
	}, "koanf"), nil)
	if err != nil {
		return Config{}, fmt.Errorf("error loading default config values: %w", err)
	}

	// Load env variables
	err = k.Load(env.ProviderWithValue("", ".", func(key, value string) (string, interface{}) {
		key = strings.ToLower(key)
		if key == "collection_interval" {
			parsedValue, err := parseCollectionInterval(value)
			if err != nil {
				slog.Warn("Error parsing collection_interval", slog.Any("error", err))
				return key, ""
			}
			return key, parsedValue
		}
		return key, value
	}), nil)
	if err != nil {
		return Config{}, fmt.Errorf("error loading env variables: %w", err)
	}

	// Unmarshal into config struct
	var cfg Config
	err = k.Unmarshal("", &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("error unmarshalling config: %w", err)
	}

	// Validate config
	v := validate.Struct(cfg)
	if !v.Validate() {
		return Config{}, v.Errors
	}
	if cfg.ZettelkastenGitURL != "" && cfg.ZettelkastenDirectory != "" {
		return Config{}, errors.New("ZettelkastenGitURL and ZettelkastenDirectory cannot be provided together")
	}
	if cfg.VictoriaMetricsURL != "" && cfg.InfluxDBURL != "" {
		return Config{}, errors.New("InfluxDBURL and VictoriaMetricsURL cannot be provided together")
	}
	if cfg.VictoriaMetricsURL == "" && cfg.InfluxDBURL == "" {
		return Config{}, errors.New("Either InfluxDBURL or VictoriaMetricsURL must be provided")
	}

	return cfg, nil
}

func (c Config) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("ZettelkastenDirectory", c.ZettelkastenDirectory),
		slog.String("ZettelkastenGitURL", c.ZettelkastenGitURL),
		slog.String("ZettelkastenGitToken", "[REDACTED]"),
		slog.String("ZettelkastenGitBranch", c.ZettelkastenGitBranch),
		slog.String("LogLevel", c.LogLevel.String()),
		slog.Any("IgnoreFiles", c.IgnoreFiles),
		slog.Duration("CollectionInterval", c.CollectionInterval),
		slog.Bool("CollectHistoricalMetrics", c.CollectHistoricalMetrics),
		slog.String("VictoriaMetricsURL", c.VictoriaMetricsURL),
		slog.String("InfluxDBURL", c.InfluxDBURL),
		slog.String("InfluxDBToken", "[REDACTED]"),
		slog.String("InfluxDBOrg", c.InfluxDBOrg),
		slog.String("InfluxDBBucket", c.InfluxDBBucket),
	)
}

func parseCollectionInterval(value string) (time.Duration, error) {
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("invalid config argument: %w", err)
	}
	return parsed, nil
}
