package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/git-town/git-town/v7/src/run"
)

// RebaseBranchStep rebases the current branch
// against the branch with the given name.
type RebaseBranchStep struct {
	EmptyStep
	Branch      string
	previousSha string
}

func (step *RebaseBranchStep) CreateAbortStep() Step {
	return &AbortRebaseStep{}
}

func (step *RebaseBranchStep) CreateContinueStep() Step {
	return &ContinueRebaseStep{}
}

func (step *RebaseBranchStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) {
	return &ResetToShaStep{Hard: true, Sha: step.previousSha}, nil
}

func (step *RebaseBranchStep) Run(repo *git.ProdRepo, connector hosting.Connector) error {
	var err error
	step.previousSha, err = repo.Runner.CurrentSha(run.Silent)
	if err != nil {
		return err
	}
	err = repo.Runner.Rebase(step.Branch, run.Logging)
	if err != nil {
		repo.Runner.CurrentBranchCache.Invalidate()
	}
	return err
}
