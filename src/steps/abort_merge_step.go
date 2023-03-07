package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/git-town/git-town/v7/src/run"
)

// AbortMergeStep aborts the current merge conflict.
type AbortMergeStep struct {
	EmptyStep
}

func (step *AbortMergeStep) Run(repo *git.ProdRepo, connector hosting.Connector) error {
	return repo.Runner.AbortMerge(run.Logging)
}
