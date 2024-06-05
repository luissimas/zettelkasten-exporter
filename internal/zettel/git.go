package zettel

import (
	"errors"
	"io/fs"
	"log/slog"
	"os"
	"time"

	"github.com/luissimas/zettelkasten-exporter/internal/config"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type GitZettel struct {
	Config         config.Config
	RepositoryPath string
}

func NewGitZettel(cfg config.Config) *GitZettel {
	return &GitZettel{RepositoryPath: "/tmp/zettelkasten-exporter", Config: cfg}
}

// GetRoot retrieves the root of the zettelkasten git repository
func (g *GitZettel) GetRoot() fs.FS {
	return os.DirFS(g.RepositoryPath)
}

// Ensure makes sure that the git repository is valid and updated with the
// latest changes from the remote.
func (g *GitZettel) Ensure() error {
	repo, err := git.PlainOpen(g.RepositoryPath)
	if errors.Is(err, git.ErrRepositoryNotExists) {
		repo, err = cloneGitRepository(g.Config.ZettelkastenGitURL, g.Config.ZettelkastenGitBranch, g.RepositoryPath)
		if err != nil {
			slog.Error("Unexpected error when cloning git repository", slog.Any("error", err), slog.String("path", g.RepositoryPath))
			return err
		}
	} else if err != nil {
		slog.Error("Unexpected error when opening git repository", slog.Any("error", err), slog.String("path", g.RepositoryPath))
		return err
	}
	slog.Debug("Git repository open", slog.String("url", g.Config.ZettelkastenGitURL), slog.String("branch", g.Config.ZettelkastenGitBranch))

	w, err := repo.Worktree()
	if err != nil {
		slog.Error("Unexpected error when getting git repository worktree", slog.Any("error", err), slog.String("url", g.Config.ZettelkastenGitURL), slog.String("branch", g.Config.ZettelkastenGitBranch))
		return err
	}

	slog.Info("Pulling from repository", slog.String("url", g.Config.ZettelkastenGitURL), slog.String("branch", g.Config.ZettelkastenGitBranch))
	start := time.Now()
	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	if errors.Is(err, git.NoErrAlreadyUpToDate) {
		slog.Info("Already up to date with remote repository, no changes pulled", slog.Duration("duration", time.Since(start)))
		return nil
	} else if err != nil {
		slog.Error("Unexpected error when pulling from git repository", slog.Any("error", err), slog.String("url", g.Config.ZettelkastenGitURL), slog.String("branch", g.Config.ZettelkastenGitBranch))
		return err
	}
	slog.Info("Pulled changes from repository", slog.Duration("duration", time.Since(start)))
	return nil
}

func cloneGitRepository(url, branch, target string) (*git.Repository, error) {
	slog.Info("Cloning git repository", slog.String("url", url), slog.String("branch", branch), slog.String("target", target))
	repo, err := git.PlainClone(target, false, &git.CloneOptions{
		URL:           url,
		Depth:         1,
		SingleBranch:  true,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
	})
	if err != nil {
		slog.Error("Could not clone git repository", slog.String("url", url), slog.String("branch", branch), slog.String("target", target))
		return nil, err
	}
	slog.Info("Git repository cloned")
	return repo, err
}
