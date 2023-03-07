package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/git-town/git-town/v7/src/run"
)

// ContinueRebaseStep finishes an ongoing rebase operation
// assuming all conflicts have been resolved by the user.
type ContinueRebaseStep struct {
	EmptyStep
}

func (step *ContinueRebaseStep) CreateAbortStep() Step {
	return &AbortRebaseStep{}
}

func (step *ContinueRebaseStep) CreateContinueStep() Step {
	return step
}

func (step *ContinueRebaseStep) Run(repo *git.ProdRepo, connector hosting.Connector) error {
	hasRebaseInProgress, err := repo.Runner.HasRebaseInProgress(run.Silent)
	if err != nil {
		return err
	}
	if hasRebaseInProgress {
		return repo.Runner.ContinueRebase(run.Logging)
	}
	return nil
}
