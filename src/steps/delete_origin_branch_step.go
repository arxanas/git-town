package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/git-town/git-town/v7/src/run"
)

// DeleteOriginBranchStep deletes the current branch from the origin remote.
type DeleteOriginBranchStep struct {
	EmptyStep
	Branch     string
	IsTracking bool
	NoPushHook bool
	branchSha  string
}

func (step *DeleteOriginBranchStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) {
	if step.IsTracking {
		return &CreateTrackingBranchStep{Branch: step.Branch, NoPushHook: step.NoPushHook}, nil
	}
	return &CreateRemoteBranchStep{Branch: step.Branch, Sha: step.branchSha, NoPushHook: step.NoPushHook}, nil
}

func (step *DeleteOriginBranchStep) Run(repo *git.ProdRepo, connector hosting.Connector) error {
	if !step.IsTracking {
		trackingBranch := repo.Runner.TrackingBranch(step.Branch, run.Silent)
		var err error
		step.branchSha, err = repo.Runner.ShaForBranch(trackingBranch, run.Silent)
		if err != nil {
			return err
		}
	}
	return repo.Runner.DeleteRemoteBranch(step.Branch, run.Logging)
}
