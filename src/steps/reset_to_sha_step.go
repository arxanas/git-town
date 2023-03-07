package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/git-town/git-town/v7/src/run"
)

// ResetToShaStep undoes all commits on the current branch
// all the way until the given SHA.
type ResetToShaStep struct {
	EmptyStep
	Hard bool
	Sha  string
}

func (step *ResetToShaStep) Run(repo *git.ProdRepo, connector hosting.Connector) error {
	currentSha, err := repo.Runner.CurrentSha(run.Silent)
	if err != nil {
		return err
	}
	if step.Sha == currentSha {
		return nil
	}
	return repo.Runner.ResetToSha(step.Sha, step.Hard, run.Logging)
}
