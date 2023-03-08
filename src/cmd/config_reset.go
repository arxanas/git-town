package cmd

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/run"
	"github.com/spf13/cobra"
)

func resetConfigCommand(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:   "reset",
		Short: "Resets your Git Town configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return repo.Config.RemoveLocalGitConfiguration(run.Silent)
		},
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return ValidateIsRepository(repo)
		},
	}
}
