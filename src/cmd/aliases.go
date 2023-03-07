package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/run"
	"github.com/spf13/cobra"
)

func aliasCommand(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:   "aliases (add | remove)",
		Short: "Adds or removes default global aliases",
		Long: `Adds or removes default global aliases

Global aliases make Git Town commands feel like native Git commands.
When enabled, you can run "git hack" instead of "git town hack".

Does not overwrite existing aliases.

This can conflict with other tools that also define Git aliases.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "add":
				err := addAliases(repo)
				if err != nil {
					return err
				}
			case "remove":
				err := removeAliases(repo)
				if err != nil {
					return err
				}
			}
			return fmt.Errorf(`invalid argument %q. Please provide either "add" or "remove"`, args[0])
		},
		Args:    cobra.ExactArgs(1),
		GroupID: "setup",
	}
}

func addAliases(repo *git.ProdRepo) error {
	for _, aliasType := range config.AliasTypes() {
		_, err1 := repo.Config.AddGitAlias(aliasType, run.Logging)
		if err1 != nil {
			return err1
		}
	}
	return nil
}

func removeAliases(repo *git.ProdRepo) error {
	for _, aliasType := range config.AliasTypes() {
		existingAlias := repo.Config.GitAlias(aliasType)
		if existingAlias == "town "+string(aliasType) {
			_, err1 := repo.Config.RemoveGitAlias(string(aliasType), run.Logging)
			if err1 != nil {
				return err1
			}
		}
	}
	return nil
}
