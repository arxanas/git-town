package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/git-town/git-town/v7/src/run"
)

// PreserveCheckoutHistoryStep does stuff.
type PreserveCheckoutHistoryStep struct {
	EmptyStep
	InitialBranch                     string
	InitialPreviouslyCheckedOutBranch string
}

func (step *PreserveCheckoutHistoryStep) Run(repo *git.ProdRepo, connector hosting.Connector) error {
	expectedPreviouslyCheckedOutBranch, err := repo.Runner.ExpectedPreviouslyCheckedOutBranch(step.InitialPreviouslyCheckedOutBranch, step.InitialBranch, run.Silent)
	if err != nil {
		return err
	}
	// NOTE: errors are not a failure condition here --> ignoring them
	previouslyCheckedOutBranch, _ := repo.Runner.PreviouslyCheckedOutBranch(run.Silent)
	if expectedPreviouslyCheckedOutBranch == previouslyCheckedOutBranch {
		return nil
	}
	currentBranch, err := repo.Runner.CurrentBranch(run.Silent)
	if err != nil {
		return err
	}
	err = repo.Runner.CheckoutBranch(expectedPreviouslyCheckedOutBranch, run.Silent)
	if err != nil {
		return err
	}
	return repo.Runner.CheckoutBranch(currentBranch, run.Silent)
}
