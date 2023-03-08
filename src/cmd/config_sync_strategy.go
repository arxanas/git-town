package cmd

import (
	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/run"
	"github.com/spf13/cobra"
)

func syncStrategyCommand(repo *git.ProdRepo) *cobra.Command {
	var globalFlag bool
	syncStrategyCmd := cobra.Command{
		Use:   "sync-strategy [(merge | rebase)]",
		Short: "Displays or sets your sync strategy",
		Long: `Displays or sets your sync strategy

The sync strategy specifies what strategy to use
when merging remote tracking branches into local feature branches.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return printSyncStrategy(globalFlag, repo)
			}
			return setSyncStrategy(globalFlag, repo, args[0])
		},
		Args: cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return ValidateIsRepository(repo)
		},
	}
	syncStrategyCmd.Flags().BoolVar(&globalFlag, "global", false, "Displays or sets the global sync strategy")
	return &syncStrategyCmd
}

func printSyncStrategy(globalFlag bool, repo *git.ProdRepo) error {
	var strategy config.SyncStrategy
	var err error
	if globalFlag {
		strategy, err = repo.Config.SyncStrategyGlobal()
	} else {
		strategy, err = repo.Config.SyncStrategy()
	}
	if err != nil {
		return err
	}
	cli.Println(strategy)
	return nil
}

func setSyncStrategy(globalFlag bool, repo *git.ProdRepo, value string) error {
	syncStrategy, err := config.ToSyncStrategy(value)
	if err != nil {
		return err
	}
	if globalFlag {
		return repo.Config.SetSyncStrategyGlobal(syncStrategy, run.Silent)
	}
	return repo.Config.SetSyncStrategy(syncStrategy, run.Silent)
}
