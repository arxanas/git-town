package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/git-town/git-town/v7/src/run"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/spf13/cobra"
)

func continueCmd(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:   "continue",
		Short: "Restarts the last run git-town command after having resolved conflicts",
		RunE: func(cmd *cobra.Command, args []string) error {
			runState, err := runstate.Load(repo)
			if err != nil {
				return fmt.Errorf("cannot load previous run state: %w", err)
			}
			if runState == nil || !runState.IsUnfinished() {
				return fmt.Errorf("nothing to continue")
			}
			hasConflicts, err := repo.Runner.HasConflicts(run.Silent)
			if err != nil {
				return err
			}
			if hasConflicts {
				return fmt.Errorf("you must resolve the conflicts before continuing")
			}
			connector, err := hosting.NewConnector(&repo.Config, &repo.Runner, cli.PrintConnectorAction)
			if err != nil {
				return err
			}
			return runstate.Execute(runState, repo, connector)
		},
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := ValidateIsRepository(repo); err != nil {
				return err
			}
			return validateIsConfigured(repo)
		},
		GroupID: "errors",
	}
}
