package zettelkasten

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

// GitZettelkasten represents a Zettelkasten versioned using git.
type GitZettelkasten struct {
	rootPath string
	url      string
	branch   string
	token    string
}

// NewGitZettelkasten creates a new GitZettelkasten.
func NewGitZettelkasten(url, branch, token string) GitZettelkasten {
	return GitZettelkasten{rootPath: "/tmp/zettelkasten-exporter", url: url, branch: branch, token: token}
}

// GetRoot retrieves the root of the zettelkasten git repository
func (g GitZettelkasten) GetRoot() fs.FS {
	return os.DirFS(g.rootPath)
}

// Ensure makes sure that the git repository is valid and updated with the latest changes from the remote.
func (g GitZettelkasten) Ensure() error {
	f, err := os.Stat(g.rootPath)
	if errors.Is(err, fs.ErrNotExist) {
		slog.Info("Zettelkasten root path does not exist, will create it", slog.String("path", g.rootPath))
		err = os.Mkdir(g.rootPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("error creating directory: %w", err)
		}
		slog.Info("Cloning repository", slog.String("path", g.rootPath), slog.String("url", g.url), slog.String("branch", g.branch))
		start := time.Now()
		err = g.cloneRepository()
		slog.Info("Repository cloned", slog.String("path", g.rootPath), slog.Duration("duration", time.Since(start)))
		if err != nil {
			return fmt.Errorf("error ensuring repository: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("error stating directory: %w", err)
	} else if !f.IsDir() {
		return errors.New("root path is not a directory")
	}

	slog.Info("Pulling repository")
	start := time.Now()
	err = g.pullRepository()
	if err != nil {
		return fmt.Errorf("error pulling repository: %w", err)
	}
	slog.Info("Pulled repository", slog.Duration("duration", time.Since(start)))

	return nil
}

// WalkHistory calls `walkFunc` for each point in the zettelkasten history.
func (g GitZettelkasten) WalkHistory(walkFunc WalkFunc) error {
	log, err := g.execInRoot("git", "log", "--reverse", "--pretty=format:%h %ad", "--date=iso")
	if err != nil {
		return fmt.Errorf("error walking zettelkasten history: %w", err)
	}
	lines := strings.Split(log, "\n")
	for _, line := range lines {
		splited := strings.SplitN(line, " ", 2)
		commit := splited[0]
		date, err := time.Parse("2006-01-02 15:04:05 -0700", splited[1])
		if err != nil {
			return fmt.Errorf("error parsing commit date: %s", err)
		}
		slog.Info("Walking commit", slog.String("commit", commit), slog.Time("date", date))
		start := time.Now()
		_, err = g.execInRoot("git", "reset", "--hard", commit)
		if err != nil {
			return fmt.Errorf("error reseting repository %w", err)
		}
		err = walkFunc(g.GetRoot(), date)
		if err != nil {
			return fmt.Errorf("error walking history: %w", err)
		}
		slog.Info("Walked commit", slog.String("commit", commit), slog.Duration("duration", time.Since(start)))
	}
	_, err = g.execInRoot("git", "reset", "--hard", fmt.Sprintf("origin/%s", g.branch))
	if err != nil {
		return fmt.Errorf("error reseting repository: %w", err)
	}
	return nil
}

func (g GitZettelkasten) cloneRepository() error {
	authenticatedUrl, err := makeAuthenticatedUrl(g.url, g.token)
	if err != nil {
		return fmt.Errorf("error clonning repository: %w", err)
	}
	err = exec.Command("git", "clone", authenticatedUrl, g.rootPath).Run()
	if err != nil {
		return fmt.Errorf("error clonning repository: %w", err)
	}
	return nil
}

func (g GitZettelkasten) pullRepository() error {
	_, err := g.execInRoot("git", "fetch", "origin")
	if err != nil {
		return fmt.Errorf("error fetching from repository: %w", err)
	}
	_, err = g.execInRoot("git", "reset", "--hard", fmt.Sprintf("origin/%s", g.branch))
	if err != nil {
		return fmt.Errorf("error reseting repository: %w", err)
	}
	return nil
}

func (g GitZettelkasten) execInRoot(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	cmd.Dir = g.rootPath
	result, err := cmd.Output()
	return string(result), err
}

func makeAuthenticatedUrl(rawURL, token string) (string, error) {
	if token == "" {
		return rawURL, nil
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("error creating authenticated url: %w", err)
	}
	parsed.User = url.UserPassword("git", token)
	return parsed.String(), nil
}
