package steps

import (
	"errors"

	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/git-town/git-town/v7/src/run"
)

// RestoreOpenChangesStep restores stashed away changes into the workspace.
type RestoreOpenChangesStep struct {
	EmptyStep
}

func (step *RestoreOpenChangesStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) {
	return &StashOpenChangesStep{}, nil
}

func (step *RestoreOpenChangesStep) Run(repo *git.ProdRepo, connector hosting.Connector) error {
	err := repo.Runner.PopStash(run.Logging)
	if err != nil {
		return errors.New("conflicts between your uncommmitted changes and the main branch")
	}
	return nil
}
