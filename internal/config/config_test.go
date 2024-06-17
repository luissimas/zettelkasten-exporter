package config

import (
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_DefaultValues(t *testing.T) {
	t.Setenv("INFLUXDB_URL", "http://localhost:8086")
	t.Setenv("INFLUXDB_TOKEN", "any-token")
	t.Setenv("INFLUXDB_ORG", "any-org")
	t.Setenv("INFLUXDB_BUCKET", "any-bucket")
	t.Setenv("ZETTELKASTEN_DIRECTORY", "/any/dir")
	c, err := LoadConfig()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := Config{
		InfluxDBURL:              "http://localhost:8086",
		InfluxDBToken:            "any-token",
		InfluxDBOrg:              "any-org",
		InfluxDBBucket:           "any-bucket",
		CollectionInterval:       time.Minute * 5,
		CollectHistoricalMetrics: true,
		LogLevel:                 slog.LevelInfo,
		ZettelkastenDirectory:    "/any/dir",
		ZettelkastenGitBranch:    "main",
		IgnoreFiles:              []string{".git", ".obsidian", ".trash", "README.md"},
	}
	assert.Equal(t, expected, c)
}

func TestLoadConfig_PartialEnv(t *testing.T) {
	t.Setenv("INFLUXDB_URL", "http://localhost:8086")
	t.Setenv("INFLUXDB_TOKEN", "any-token")
	t.Setenv("INFLUXDB_ORG", "any-org")
	t.Setenv("INFLUXDB_BUCKET", "any-bucket")
	t.Setenv("LOG_LEVEL", "DEBUG")
	t.Setenv("ZETTELKASTEN_DIRECTORY", "/any/dir")
	c, err := LoadConfig()
	if assert.NoError(t, err) {
		expected := Config{
			InfluxDBURL:              "http://localhost:8086",
			InfluxDBToken:            "any-token",
			InfluxDBOrg:              "any-org",
			InfluxDBBucket:           "any-bucket",
			CollectionInterval:       time.Minute * 5,
			CollectHistoricalMetrics: true,
			LogLevel:                 slog.LevelDebug,
			ZettelkastenDirectory:    "/any/dir",
			ZettelkastenGitBranch:    "main",
			IgnoreFiles:              []string{".git", ".obsidian", ".trash", "README.md"},
		}
		assert.Equal(t, expected, c)
	}
}

func TestLoadConfig_FullEnvDirectory(t *testing.T) {
	t.Setenv("INFLUXDB_URL", "http://localhost:8086")
	t.Setenv("INFLUXDB_TOKEN", "any-token")
	t.Setenv("INFLUXDB_ORG", "any-org")
	t.Setenv("INFLUXDB_BUCKET", "any-bucket")
	t.Setenv("COLLECTION_INTERVAL", "2h")
	t.Setenv("COLLECT_HISTORICAL_METRICS", "false")
	t.Setenv("LOG_LEVEL", "WARN")
	t.Setenv("ZETTELKASTEN_DIRECTORY", "/any/dir")
	t.Setenv("IGNORE_FILES", ".obsidian,test,/something/another,dir/file.md")
	c, err := LoadConfig()
	if assert.NoError(t, err) {
		expected := Config{
			InfluxDBURL:              "http://localhost:8086",
			InfluxDBToken:            "any-token",
			InfluxDBOrg:              "any-org",
			InfluxDBBucket:           "any-bucket",
			CollectionInterval:       time.Hour * 2,
			CollectHistoricalMetrics: false,
			LogLevel:                 slog.LevelWarn,
			ZettelkastenDirectory:    "/any/dir",
			ZettelkastenGitBranch:    "main",
			IgnoreFiles:              []string{".obsidian", "test", "/something/another", "dir/file.md"},
		}
		assert.Equal(t, expected, c)
	}
}

func TestLoadConfig_FullEnvGit(t *testing.T) {
	t.Setenv("INFLUXDB_URL", "http://localhost:8086")
	t.Setenv("INFLUXDB_TOKEN", "any-token")
	t.Setenv("INFLUXDB_ORG", "any-org")
	t.Setenv("INFLUXDB_BUCKET", "any-bucket")
	t.Setenv("COLLECTION_INTERVAL", "15m")
	t.Setenv("COLLECT_HISTORICAL_METRICS", "false")
	t.Setenv("LOG_LEVEL", "ERROR")
	t.Setenv("ZETTELKASTEN_GIT_URL", "https://github.com/user/zettel")
	t.Setenv("IGNORE_FILES", ".obsidian,test,/something/another,dir/file.md")
	c, err := LoadConfig()
	if assert.NoError(t, err) {
		expected := Config{
			InfluxDBURL:              "http://localhost:8086",
			InfluxDBToken:            "any-token",
			InfluxDBOrg:              "any-org",
			InfluxDBBucket:           "any-bucket",
			CollectionInterval:       time.Minute * 15,
			CollectHistoricalMetrics: false,
			LogLevel:                 slog.LevelError,
			ZettelkastenGitURL:       "https://github.com/user/zettel",
			ZettelkastenGitBranch:    "main",
			IgnoreFiles:              []string{".obsidian", "test", "/something/another", "dir/file.md"},
		}
		assert.Equal(t, expected, c)
	}
}

func TestLoadConfig_Validate(t *testing.T) {
	data := []struct {
		name        string
		shouldError bool
		env         map[string]string
	}{
		{
			name:        "missing source",
			shouldError: true,
			env: map[string]string{
				"LOG_LEVEL": "INFO",
			},
		},
		{
			name:        "both sources",
			shouldError: true,
			env: map[string]string{
				"LOG_LEVEL":              "INFO",
				"ZETTELKASTEN_DIRECTORY": "/any/dir",
				"ZETTELKASTEN_GIT_URL":   "any-string",
			},
		},
		{
			name:        "valid config",
			shouldError: false,
			env: map[string]string{
				"LOG_LEVEL":                  "INFO",
				"ZETTELKASTEN_GIT_URL":       "any-url",
				"ZETTELKASTEN_GIT_BRANCH":    "any-branch",
				"COLLECTION_INTERVAL":        "15m",
				"COLLECT_HISTORICAL_METRICS": "false",
				"INFLUXDB_URL":               "http://localhost:8086",
				"INFLUXDB_TOKEN":             "any-token",
				"INFLUXDB_ORG":               "any-org",
				"INFLUXDB_BUCKET":            "any-bucket",
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
