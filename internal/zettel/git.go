package zettel

import (
	"errors"
	"io/fs"
	"log/slog"
	"os"

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
	slog.Info("Git repository open", slog.Any("repo", repo))

	w, err := repo.Worktree()
	if err != nil {
		slog.Error("Unexpected error when getting git repository worktree", slog.Any("error", err), slog.Any("repo", repo))
		return err
	}

	slog.Info("Pulling from repository", slog.Any("repo", repo))
	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil {
		slog.Error("Unexpected error when pulling from git repository", slog.Any("error", err), slog.Any("repo", repo))
		return err
	}

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
