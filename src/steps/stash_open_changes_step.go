package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/git-town/git-town/v7/src/run"
)

type StashOpenChangesStep struct {
	EmptyStep
}

func (step *StashOpenChangesStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) {
	return &RestoreOpenChangesStep{}, nil
}

func (step *StashOpenChangesStep) Run(repo *git.ProdRepo, connector hosting.Connector) error {
	return repo.Runner.Stash(run.Logging)
}
