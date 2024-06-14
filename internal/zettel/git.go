package zettel

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"time"

	"github.com/luissimas/zettelkasten-exporter/internal/config"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
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
		repo, err = cloneRepository(g.Config.ZettelkastenGitURL, g.Config.ZettelkastenGitBranch, g.RepositoryPath)
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
		slog.Error("Unexpected error when getting git repository worktree", slog.Any("error", err))
		return err
	}
	branch, err := repo.Branch(g.Config.ZettelkastenGitBranch)
	if err != nil {
		slog.Error("Unexpected error when getting git repository branch", slog.Any("error", err))
		return err
	}

	rev, err := repo.ResolveRevision(plumbing.Revision(branch.Name))
	if err != nil {
		slog.Error("Unexpected error when getting git repository remote revision", slog.Any("error", err))
		return err
	}

	w.Reset(&git.ResetOptions{Commit: *rev})

	slog.Info("Pulling from repository", slog.String("url", g.Config.ZettelkastenGitURL), slog.String("branch", g.Config.ZettelkastenGitBranch))
	start := time.Now()
	err = forcePullRepository(repo)
	if err != nil {
		slog.Error("Unexpected error when pulling from git repository", slog.Any("error", err), slog.String("url", g.Config.ZettelkastenGitURL), slog.String("branch", g.Config.ZettelkastenGitBranch))
		return err
	}
	slog.Info("Pulled changes from repository", slog.Duration("duration", time.Since(start)))

	return nil
}

func (g *GitZettel) WalkHistory(walkFunc func(time.Time) error) error {
	repo, err := git.PlainOpen(g.RepositoryPath)
	if err != nil {
		slog.Error("Unexpected error when opening git repository", slog.Any("error", err), slog.String("path", g.RepositoryPath))
		return err
	}
	w, err := repo.Worktree()
	if err != nil {
		slog.Error("Unexpected error when getting git repository worktree", slog.Any("error", err))
		return err
	}
	log, err := repo.Log(&git.LogOptions{Order: git.LogOrderCommitterTime})
	log.ForEach(func(c *object.Commit) error {
		slog.Debug("Walking commit", slog.String("sha", c.Hash.String()), slog.String("message", c.Message), slog.Time("date", c.Committer.When))
		w.Reset(&git.ResetOptions{Commit: c.Hash, Mode: git.HardReset})
		err := walkFunc(c.Committer.When)
		if err != nil {
			slog.Error("Error when walking commit", slog.String("hash", c.Hash.String()), slog.Any("error", err))
			return err
		}
		return nil
	})
	return nil
}

func cloneRepository(url, branch, target string) (*git.Repository, error) {
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

func forcePullRepository(repo *git.Repository) error {
	w, err := repo.Worktree()
	if err != nil {
		slog.Error("Unexpected error when getting git repository worktree", slog.Any("error", err))
		return err
	}

	// NOTE: instead of just pulling, we fetch and then hard reset to
	// account for the case of force pushes to the remote branch
	err = repo.Fetch(&git.FetchOptions{Depth: 2147483647})

	if errors.Is(err, git.NoErrAlreadyUpToDate) {
		slog.Info("Already up to date with remote repository, no changes pulled")
		return nil
	} else if err != nil {
		slog.Error("Unexpected error when fetching from git repository", slog.Any("error", err))
		return err
	}

	head, err := repo.Head()
	if err != nil {
		slog.Error("Unexpected error when getting git repository head", slog.Any("error", err))
		return err
	}

	branch, err := repo.Branch(head.Name().Short())
	if err != nil {
		slog.Error("Unexpected error when getting git repository branch", slog.Any("error", err))
		return err
	}

	rev, err := repo.ResolveRevision(plumbing.Revision(fmt.Sprintf("remotes/%s/%s", branch.Remote, head.Name().Short())))
	if err != nil {
		slog.Error("Unexpected error when getting git repository remote revision", slog.Any("error", err))
		return err
	}

	err = w.Reset(&git.ResetOptions{
		Commit: *rev,
		Mode:   git.HardReset,
	})
	if err != nil {
		slog.Error("Unexpected error when reseting from git repository", slog.Any("error", err))
		return err
	}
	return nil
}
