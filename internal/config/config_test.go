package config

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_DefaultValues(t *testing.T) {
	t.Setenv("ZETTELKASTEN_DIRECTORY", "/any/dir")
	c, err := LoadConfig()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := Config{
		IP:                    "0.0.0.0",
		Port:                  10018,
		LogLevel:              slog.LevelInfo,
		ZettelkastenDirectory: "/any/dir",
		ZettelkastenGitBranch: "main",
		IgnoreFiles:           []string{".git", ".obsidian", ".trash", "README.md"},
	}
	assert.Equal(t, expected, c)
}

func TestLoadConfig_PartialEnv(t *testing.T) {
	t.Setenv("PORT", "4444")
	t.Setenv("LOG_LEVEL", "DEBUG")
	t.Setenv("ZETTELKASTEN_DIRECTORY", "/any/dir")
	c, err := LoadConfig()
	if assert.NoError(t, err) {
		expected := Config{
			IP:                    "0.0.0.0",
			Port:                  4444,
			LogLevel:              slog.LevelDebug,
			ZettelkastenDirectory: "/any/dir",
			ZettelkastenGitBranch: "main",
			IgnoreFiles:           []string{".git", ".obsidian", ".trash", "README.md"},
		}
		assert.Equal(t, expected, c)
	}
}

func TestLoadConfig_FullEnvDirectory(t *testing.T) {
	t.Setenv("IP", "127.0.0.1")
	t.Setenv("PORT", "4444")
	t.Setenv("LOG_LEVEL", "DEBUG")
	t.Setenv("ZETTELKASTEN_DIRECTORY", "/any/dir")
	t.Setenv("IGNORE_FILES", ".obsidian,test,/something/another,dir/file.md")
	c, err := LoadConfig()
	if assert.NoError(t, err) {
		expected := Config{
			IP:                    "127.0.0.1",
			Port:                  4444,
			LogLevel:              slog.LevelDebug,
			ZettelkastenDirectory: "/any/dir",
			ZettelkastenGitBranch: "main",
			IgnoreFiles:           []string{".obsidian", "test", "/something/another", "dir/file.md"},
		}
		assert.Equal(t, expected, c)
	}
}

func TestLoadConfig_FullEnvGit(t *testing.T) {
	t.Setenv("IP", "127.0.0.1")
	t.Setenv("PORT", "4444")
	t.Setenv("LOG_LEVEL", "DEBUG")
	t.Setenv("ZETTELKASTEN_GIT_URL", "https://github.com/user/zettel")
	t.Setenv("IGNORE_FILES", ".obsidian,test,/something/another,dir/file.md")
	c, err := LoadConfig()
	if assert.NoError(t, err) {
		expected := Config{
			IP:                    "127.0.0.1",
			Port:                  4444,
			LogLevel:              slog.LevelDebug,
			ZettelkastenGitURL:    "https://github.com/user/zettel",
			ZettelkastenGitBranch: "main",
			IgnoreFiles:           []string{".obsidian", "test", "/something/another", "dir/file.md"},
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
				"IP":        "0.0.0.0",
				"PORT":      "4444",
				"LOG_LEVEL": "INFO",
			},
		},
		{
			name:        "both sources",
			shouldError: true,
			env: map[string]string{
				"IP":                     "0.0.0.0",
				"PORT":                   "4444",
				"LOG_LEVEL":              "INFO",
				"ZETTELKASTEN_DIRECTORY": "/any/dir",
				"ZETTELKASTEN_GIT_URL":   "any-string",
			},
		},
		{
			name:        "invalid ip",
			shouldError: true,
			env: map[string]string{
				"IP":                     "any-string",
				"PORT":                   "4444",
				"LOG_LEVEL":              "INFO",
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
				"ZETTELKASTEN_DIRECTORY": "/any/dir",
			},
		},
		{
			name:        "valid config",
			shouldError: false,
			env: map[string]string{
				"IP":                      "0.0.0.0",
				"PORT":                    "4444",
				"LOG_LEVEL":               "INFO",
				"ZETTELKASTEN_GIT_URL":    "any-url",
				"ZETTELKASTEN_GIT_BRANCH": "any-branch",
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
