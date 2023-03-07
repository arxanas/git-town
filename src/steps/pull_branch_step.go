package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/git-town/git-town/v7/src/run"
)

// PullBranchStep updates the branch with the given name with commits from its remote.
type PullBranchStep struct {
	EmptyStep
	Branch string
}

func (step *PullBranchStep) Run(repo *git.ProdRepo, connector hosting.Connector) error {
	return repo.Runner.Pull(run.Logging)
}
