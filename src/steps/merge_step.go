package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/git-town/git-town/v7/src/run"
)

// MergeStep merges the branch with the given name into the current branch.
type MergeStep struct {
	EmptyStep
	Branch      string
	previousSha string
}

func (step *MergeStep) CreateAbortStep() Step {
	return &AbortMergeStep{}
}

func (step *MergeStep) CreateContinueStep() Step {
	return &ContinueMergeStep{}
}

func (step *MergeStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) {
	return &ResetToShaStep{Hard: true, Sha: step.previousSha}, nil
}

func (step *MergeStep) Run(repo *git.ProdRepo, connector hosting.Connector) error {
	var err error
	step.previousSha, err = repo.Runner.CurrentSha(run.Silent)
	if err != nil {
		return err
	}
	return repo.Runner.MergeBranchNoEdit(step.Branch, run.Logging)
}
