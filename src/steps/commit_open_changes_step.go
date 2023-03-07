package steps

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/git-town/git-town/v7/src/run"
)

// CommitOpenChangesStep commits all open changes as a new commit.
// It does not ask the user for a commit message, but chooses one automatically.
type CommitOpenChangesStep struct {
	EmptyStep
	previousSha string
}

func (step *CommitOpenChangesStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) {
	return &ResetToShaStep{Sha: step.previousSha}, nil
}

func (step *CommitOpenChangesStep) Run(repo *git.ProdRepo, connector hosting.Connector) error {
	var err error
	step.previousSha, err = repo.Runner.CurrentSha(run.Silent)
	if err != nil {
		return err
	}
	err = repo.Runner.StageFiles(run.Logging, "-A")
	if err != nil {
		return err
	}
	currentBranch, err := repo.Runner.CurrentBranch(run.Silent)
	if err != nil {
		return err
	}
	return repo.Runner.CommitStagedChanges(fmt.Sprintf("WIP on %s", currentBranch), run.Logging)
}
