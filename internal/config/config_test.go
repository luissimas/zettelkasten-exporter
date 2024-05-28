package config

import (
	"log/slog"
	"testing"
	"time"
)

func TestLoadConfig_DefaultValues(t *testing.T) {
	t.Setenv("ZETTELKASTEN_DIRECTORY", "/any/dir")
	c, err := LoadConfig()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := Config{
		IP:                    "0.0.0.0",
		Port:                  6969,
		LogLevel:              slog.LevelInfo,
		ScrapeInterval:        time.Minute * 5,
		ZettelkastenDirectory: "/any/dir",
	}
	if c != expected {
		t.Errorf("Expected %v, got: %v", expected, c)
	}
}

func TestLoadConfig_PartialEnv(t *testing.T) {
	t.Setenv("PORT", "4444")
	t.Setenv("LOG_LEVEL", "DEBUG")
	t.Setenv("ZETTELKASTEN_DIRECTORY", "/any/dir")
	c, err := LoadConfig()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := Config{
		IP:                    "0.0.0.0",
		Port:                  4444,
		LogLevel:              slog.LevelDebug,
		ScrapeInterval:        time.Minute * 5,
		ZettelkastenDirectory: "/any/dir",
	}
	if c != expected {
		t.Errorf("Expected %v, got: %v", expected, c)
	}
}

func TestLoadConfig_FullEnv(t *testing.T) {
	t.Setenv("IP", "127.0.0.1")
	t.Setenv("PORT", "4444")
	t.Setenv("LOG_LEVEL", "DEBUG")
	t.Setenv("SCRAPE_INTERVAL", "5m")
	t.Setenv("ZETTELKASTEN_DIRECTORY", "/any/dir")
	c, err := LoadConfig()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := Config{
		IP:                    "127.0.0.1",
		Port:                  4444,
		LogLevel:              slog.LevelDebug,
		ScrapeInterval:        time.Minute * 5,
		ZettelkastenDirectory: "/any/dir",
	}
	if c != expected {
		t.Errorf("Expected %v, got: %v", expected, c)
	}
}

func TestLoadConfig_Validate(t *testing.T) {
	data := []struct {
		name        string
		shouldError bool
		env         map[string]string
	}{
		{
			name:        "missing directory",
			shouldError: true,
			env: map[string]string{
				"IP":        "0.0.0.0",
				"PORT":      "4444",
				"LOG_LEVEL": "INFO",
			},
		},
		{
			name:        "invalid ip",
			shouldError: true,
			env: map[string]string{
				"IP":                     "any-string",
				"PORT":                   "4444",
				"LOG_LEVEL":              "INFO",
				"SCRAPE_INTERVAL":        "5m",
				"ZETTELKASTEN_DIRECTORY": "/any/dir",
			},
		},
		{
			name:        "invalid port",
			shouldError: true,
			env: map[string]string{
				"IP":                     "0.0.0.0",
				"PORT":                   "-1",
				"LOG_LEVEL":              "INFO",
				"SCRAPE_INTERVAL":        "5m",
				"ZETTELKASTEN_DIRECTORY": "/any/dir",
			},
		},
		{
			name:        "invalid interval",
			shouldError: true,
			env: map[string]string{
				"IP":                     "0.0.0.0",
				"PORT":                   "4444",
				"LOG_LEVEL":              "INFO",
				"SCRAPE_INTERVAL":        "5",
				"ZETTELKASTEN_DIRECTORY": "/any/dir",
			},
		},
		{
			name:        "valid config",
			shouldError: false,
			env: map[string]string{
				"IP":                     "0.0.0.0",
				"PORT":                   "4444",
				"LOG_LEVEL":              "INFO",
				"SCRAPE_INTERVAL":        "5m",
				"ZETTELKASTEN_DIRECTORY": "/any/dir",
			},
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			for k, v := range d.env {
				t.Setenv(k, v)
			}
			_, err := LoadConfig()
			if err == nil && d.shouldError {
				t.Errorf("Expected error, got: %v", err)
			} else if err != nil && !d.shouldError {
				t.Errorf("Expected no error, got: %v", err)
			}
		})
	}
}
