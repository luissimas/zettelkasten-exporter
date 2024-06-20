package zettelkasten

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// GitZettelkasten represents a Zettelkasten versioned using git.
type GitZettelkasten struct {
	RootPath string
	URL      string
	Branch   string
	Token    string
}

// NewGitZettelkasten creates a new GitZettelkasten.
func NewGitZettelkasten(url, branch, token string) GitZettelkasten {
	return GitZettelkasten{RootPath: "/tmp/zettelkasten-exporter", URL: url, Branch: branch, Token: token}
}

// GetRoot retrieves the root of the zettelkasten git repository
func (g GitZettelkasten) GetRoot() fs.FS {
	return os.DirFS(g.RootPath)
}

// Ensure makes sure that the git repository is valid and updated with the
// latest changes from the remote.
func (g GitZettelkasten) Ensure() error {
	repo, err := git.PlainOpen(g.RootPath)
	if errors.Is(err, git.ErrRepositoryNotExists) {
		repo, err = cloneRepository(g.URL, g.Branch, g.RootPath, g.Token)
		if err != nil {
			slog.Error("Unexpected error when cloning git repository", slog.Any("error", err), slog.String("path", g.RootPath))
			return err
		}
	} else if err != nil {
		slog.Error("Unexpected error when opening git repository", slog.Any("error", err), slog.String("path", g.RootPath))
		return err
	}
	slog.Debug("Git repository open", slog.String("url", g.URL), slog.String("branch", g.Branch))
	w, err := repo.Worktree()
	if err != nil {
		slog.Error("Unexpected error when getting git repository worktree", slog.Any("error", err))
		return err
	}
	branch, err := repo.Branch(g.Branch)
	if err != nil {
		slog.Error("Unexpected error when getting git repository branch", slog.Any("error", err))
		return err
	}

	rev, err := repo.ResolveRevision(plumbing.Revision(branch.Name))
	if err != nil {
		slog.Error("Unexpected error when getting git repository remote revision", slog.Any("error", err))
		return err
	}
	err = w.Reset(&git.ResetOptions{
		Commit: *rev,
		Mode:   git.HardReset,
	})
	if err != nil {
		slog.Error("Unexpected error when reseting git repository", slog.Any("error", err))
		return err
	}

	slog.Info("Pulling from repository", slog.String("url", g.URL), slog.String("branch", g.Branch))
	start := time.Now()
	err = forcePullRepository(repo, g.Token)
	if err != nil {
		slog.Error("Unexpected error when pulling from git repository", slog.Any("error", err), slog.String("url", g.URL), slog.String("branch", g.Branch))
		return err
	}
	slog.Info("Pulled changes from repository", slog.Duration("duration", time.Since(start)))

	return nil
}

// WalkHistory calls `walkFunc` for each point in the zettelkasten history.
func (g GitZettelkasten) WalkHistory(walkFunc WalkFunc) error {
	repo, err := git.PlainOpen(g.RootPath)
	if err != nil {
		slog.Error("Unexpected error when opening git repository", slog.Any("error", err), slog.String("path", g.RootPath))
		return err
	}
	originalHead, err := repo.Head()
	if err != nil {
		slog.Error("Unexpected error when getting git repository head", slog.Any("error", err))
		return err
	}
	originalHash := originalHead.Hash()
	w, err := repo.Worktree()
	if err != nil {
		slog.Error("Unexpected error when getting git repository worktree", slog.Any("error", err))
		return err
	}
	log, err := repo.Log(&git.LogOptions{Order: git.LogOrderCommitterTime})
	err = log.ForEach(func(c *object.Commit) error {
		slog.Debug("Walking commit", slog.String("sha", c.Hash.String()), slog.String("message", c.Message), slog.Time("date", c.Committer.When))
		err = w.Reset(&git.ResetOptions{
			Commit: c.Hash,
			Mode:   git.HardReset,
		})
		if err != nil {
			slog.Error("Unexpected error when reseting git repository", slog.Any("error", err))
			return err
		}
		err := walkFunc(g.GetRoot(), c.Committer.When)
		if err != nil {
			slog.Error("Error when walking commit", slog.String("hash", c.Hash.String()), slog.Any("error", err))
			return err
		}
		return nil
	})
	err = w.Reset(&git.ResetOptions{
		Commit: originalHash,
		Mode:   git.HardReset,
	})
	if err != nil {
		slog.Error("Unexpected error when reseting git repository", slog.Any("error", err))
		return err
	}
	return nil
}

func cloneRepository(url, branch, target, token string) (*git.Repository, error) {
	slog.Info("Cloning git repository", slog.String("url", url), slog.String("branch", branch), slog.String("target", target))
	cloneOptions := git.CloneOptions{
		URL:           url,
		Depth:         1,
		SingleBranch:  true,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
	}
	if token != "" {
		cloneOptions.Auth = &http.BasicAuth{
			Username: "git",
			Password: token,
		}
	}
	repo, err := git.PlainClone(target, false, &cloneOptions)
	if err != nil {
		slog.Error("Could not clone git repository", slog.String("url", url), slog.String("branch", branch), slog.String("target", target))
		return nil, err
	}
	slog.Info("Git repository cloned")
	return repo, err
}

func forcePullRepository(repo *git.Repository, token string) error {
	w, err := repo.Worktree()
	if err != nil {
		slog.Error("Unexpected error when getting git repository worktree", slog.Any("error", err))
		return err
	}

	// NOTE: instead of just pulling, we fetch and then hard reset to
	// account for the case of force pushes to the remote branch
	fetchOptions := git.FetchOptions{Depth: 2147483647}
	if token != "" {
		fetchOptions.Auth = &http.BasicAuth{
			Username: "git",
			Password: token,
		}
	}
	err = repo.Fetch(&fetchOptions)

	if errors.Is(err, git.NoErrAlreadyUpToDate) {
		slog.Info("Already up to date with remote repository, no changes pulled")
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
		slog.Error("Unexpected error when reseting git repository", slog.Any("error", err))
		return err
	}
	return nil
}
