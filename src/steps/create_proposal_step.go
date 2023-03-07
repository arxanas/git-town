package steps

import (
	"github.com/git-town/git-town/v7/src/browser"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/git-town/git-town/v7/src/run"
)

// CreateProposalStep creates a new pull request for the current branch.
type CreateProposalStep struct {
	EmptyStep
	Branch string
}

func (step *CreateProposalStep) Run(repo *git.ProdRepo, connector hosting.Connector) error {
	parentBranch := repo.Config.ParentBranch(step.Branch)
	prURL, err := connector.NewProposalURL(step.Branch, parentBranch)
	if err != nil {
		return err
	}
	browser.Open(prURL, run.Logging)
	return nil
}
