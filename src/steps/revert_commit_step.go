package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/git-town/git-town/v7/src/run"
)

// RevertCommitStep reverts the commit with the given sha.
type RevertCommitStep struct {
	EmptyStep
	Sha string
}

func (step *RevertCommitStep) Run(repo *git.ProdRepo, connector hosting.Connector) error {
	return repo.Runner.RevertCommit(step.Sha, run.Logging)
}
