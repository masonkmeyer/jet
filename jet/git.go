package jet

import (
	"fmt"
	"sort"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type (
	// Branch represents a git branch
	Branch struct {
		Name string
		When time.Time
	}
	// Commit represents a git commit
	Commit struct {
		Hash    string
		Author  string
		Message string
		When    time.Time
	}
	// Git is a wrapper around the git command
	Git struct {
		repo *git.Repository
	}
)

// NewGit creates a new Git instance
func NewGit() (*Git, error) {
	r, err := git.PlainOpen(".")

	if err != nil {
		return nil, err
	}

	return &Git{
		repo: r,
	}, nil
}

// IsClean returns true if the current branch is clean (no uncommited changes)
func (g Git) IsClean() bool {
	w, err := g.repo.Worktree()

	if err != nil {
		return false
	}

	status, err := w.Status()

	if err != nil {
		return false
	}

	return status.IsClean()
}

// Checkout executes a git checkout command with the branch name
// and returns the output and the command that was executed
func (g Git) Checkout(branchName string) error {
	w, err := g.repo.Worktree()

	if err != nil {
		return err
	}

	branchRefName := plumbing.NewBranchReferenceName(branchName)
	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(branchRefName),
	})

	return err
}

// ListBranches executes a git branch command with --list and any other provided args
func (g Git) ListBranches() []Branch {
	branches, _ := g.repo.Branches()

	results := []Branch{}
	branches.ForEach(func(ref *plumbing.Reference) error {
		commit, err := g.repo.CommitObject(ref.Hash())

		if err != nil {
			return err
		}

		results = append(results, Branch{
			Name: ref.Name().Short(),
			When: commit.Author.When,
		})

		return nil
	})

	// sort branches by most recent commit
	sort.Slice(results, func(i, j int) bool {
		return results[i].When.After(results[j].When)
	})

	return results
}

// Log executes a git log command with any provided args and --no-pager to prevent interactive paging
func (g Git) Logs(branchName string, n int) []Commit {
	name := plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branchName))
	ref, _ := g.repo.Reference(name, true)

	yearAgo := time.Now().AddDate(-1, 0, 0)
	logIter, _ := g.repo.Log(&git.LogOptions{From: ref.Hash(), Since: &yearAgo})

	results := []Commit{}
	logIter.ForEach(func(c *object.Commit) error {
		if len(results) >= n {
			return nil
		}

		results = append(results, Commit{
			Hash:    c.Hash.String()[:7],
			Author:  c.Author.Name,
			Message: c.Message,
			When:    c.Author.When,
		})

		return nil
	})

	return results
}
