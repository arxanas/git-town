package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/git-town/git-town/v7/src/run"
)

// PushTagsStep pushes newly created Git tags to origin.
type PushTagsStep struct {
	EmptyStep
}

func (step *PushTagsStep) Run(repo *git.ProdRepo, connector hosting.Connector) error {
	return repo.Runner.PushTags(run.Logging)
}
