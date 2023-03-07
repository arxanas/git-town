package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/git-town/git-town/v7/src/run"
)

// DeleteLocalBranchStep deletes the branch with the given name,
// optionally in a safe or unsafe way.
type DeleteLocalBranchStep struct {
	EmptyStep
	Branch    string
	Force     bool
	branchSha string
}

func (step *DeleteLocalBranchStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) {
	return &CreateBranchStep{Branch: step.Branch, StartingPoint: step.branchSha}, nil
}

func (step *DeleteLocalBranchStep) Run(repo *git.ProdRepo, connector hosting.Connector) error {
	var err error
	step.branchSha, err = repo.Runner.ShaForBranch(step.Branch, run.Silent)
	if err != nil {
		return err
	}
	hasUnmergedCommits, err := repo.Runner.BranchHasUnmergedCommits(step.Branch, run.Silent)
	if err != nil {
		return err
	}
	return repo.Runner.DeleteLocalBranch(step.Branch, step.Force || hasUnmergedCommits, run.Logging)
}
