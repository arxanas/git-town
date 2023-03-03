package cmd

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/dialog"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/git-town/git-town/v7/src/steps"
	"github.com/git-town/git-town/v7/src/stringslice"
	"github.com/spf13/cobra"
)

func syncCmd(repo *git.ProdRepo) *cobra.Command {
	var allFlag bool
	var dryRunFlag bool
	syncCmd := cobra.Command{
		Use:   "sync",
		Short: "Updates the current branch with all relevant changes",
		Long: fmt.Sprintf(`Updates the current branch with all relevant changes

Synchronizes the current branch with the rest of the world.

When run on a feature branch
- syncs all ancestor branches
- pulls updates for the current branch
- merges the parent branch into the current branch
- pushes the current branch

When run on the main branch or a perennial branch
- pulls and pushes updates for the current branch
- pushes tags

If the repository contains an "upstream" remote,
syncs the main branch with its upstream counterpart.
You can disable this by running "git config %s false".`, config.SyncUpstream),
		Run: func(cmd *cobra.Command, args []string) {
			config, err := determineSyncConfig(allFlag, repo)
			if err != nil {
				cli.Exit(err)
			}
			stepList, err := syncBranchesSteps(config, repo)
			if err != nil {
				cli.Exit(err)
			}
			runState := runstate.New("sync", stepList)
			err = runstate.Execute(runState, repo, nil)
			if err != nil {
				cli.Exit(err)
			}
		},
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := ValidateIsRepository(repo); err != nil {
				return err
			}
			if dryRunFlag {
				currentBranch, err := repo.Silent.CurrentBranch()
				if err != nil {
					return err
				}
				repo.DryRun.Activate(currentBranch)
			}
			if err := validateIsConfigured(repo); err != nil {
				return err
			}
			exit, err := handleUnfinishedState(repo, nil)
			if err != nil {
				return err
			}
			if exit {
				os.Exit(0)
			}
			return nil
		},
		GroupID: "basic",
	}
	syncCmd.Flags().BoolVar(&allFlag, "all", false, "Sync all local branches")
	syncCmd.Flags().BoolVar(&dryRunFlag, "dry-run", false, "Print the commands but don't run them")
	return &syncCmd
}

type syncConfig struct {
	branchesToSync            []string
	branchesWithDeletedRemote []string // local branches whose tracking branches have been deleted
	hasOrigin                 bool
	initialBranch             string
	isOffline                 bool
	mainBranch                string
	shouldPushTags            bool
}

func determineSyncConfig(allFlag bool, repo *git.ProdRepo) (*syncConfig, error) {
	hasOrigin, err := repo.Silent.HasOrigin()
	if err != nil {
		return nil, err
	}
	isOffline, err := repo.Config.IsOffline()
	if err != nil {
		return nil, err
	}
	mainbranch := repo.Config.MainBranch()
	if hasOrigin && !isOffline {
		err := repo.Logging.Fetch()
		if err != nil {
			return nil, err
		}
	}
	branchesWithDeletedRemote, err := repo.Silent.LocalBranchesWithDeletedTrackingBranches()
	if err != nil {
		return nil, err
	}
	initialBranch, err := repo.Silent.CurrentBranch()
	if err != nil {
		return nil, err
	}
	parentDialog := dialog.ParentBranches{}
	var branchesToSync []string
	var shouldPushTags bool
	if allFlag {
		branches, err := repo.Silent.LocalBranchesMainFirst()
		if err != nil {
			return nil, err
		}
		err = parentDialog.EnsureKnowsParentBranches(branches, repo)
		if err != nil {
			return nil, err
		}
		branchesToSync = branches
		shouldPushTags = true
	} else {
		err = parentDialog.EnsureKnowsParentBranches([]string{initialBranch}, repo)
		if err != nil {
			return nil, err
		}
		branchesToSync = append(repo.Config.AncestorBranches(initialBranch), initialBranch)
		shouldPushTags = !repo.Config.IsFeatureBranch(initialBranch)
	}
	return &syncConfig{
		branchesToSync:            branchesToSync,
		branchesWithDeletedRemote: branchesWithDeletedRemote,
		hasOrigin:                 hasOrigin,
		initialBranch:             initialBranch,
		isOffline:                 isOffline,
		mainBranch:                mainbranch,
		shouldPushTags:            shouldPushTags,
	}, nil
}

// syncBranchesSteps provides the step list for the "git sync" command.
func syncBranchesSteps(config *syncConfig, repo *git.ProdRepo) (runstate.StepList, error) {
	list := runstate.StepListBuilder{}
	for _, branch := range config.branchesToSync {
		syncStepsForBranch(&list, branch, config, repo)
	}
	list.Add(&steps.CheckoutStep{Branch: config.initialBranch})
	if config.hasOrigin && config.shouldPushTags && !config.isOffline {
		list.Add(&steps.PushTagsStep{})
	}
	list.Wrap(runstate.WrapOptions{RunInGitRoot: true, StashOpenChanges: true}, repo)
	return list.Result()
}

// hasDeletedTrackingBranch indicates whether the tracking branch of the branch with the given name was deleted on the remote.
func (sc *syncConfig) hasDeletedTrackingBranch(branch string) bool {
	return stringslice.Contains(sc.branchesWithDeletedRemote, branch)
}

// syncStepsForBranch adds the steps to sync the branch with the given name to the given step list.
func syncStepsForBranch(list *runstate.StepListBuilder, branch string, config *syncConfig, repo *git.ProdRepo) {
	if config.hasDeletedTrackingBranch(branch) {
		deleteBranchSteps(list, branch, config, repo)
	} else {
		updateBranchSteps(list, branch, true, config.branchesWithDeletedRemote, repo)
	}
}

// deleteBranchSteps adds the steps to delete the branch with the given name locally.
func deleteBranchSteps(list *runstate.StepListBuilder, branch string, config *syncConfig, repo *git.ProdRepo) {
	if config.initialBranch == branch {
		list.Add(&steps.CheckoutStep{Branch: config.mainBranch})
	}
	parent := repo.Config.ParentBranch(branch)
	if parent != "" {
		for _, child := range repo.Config.ChildBranches(branch) {
			list.Add(&steps.SetParentStep{Branch: child, ParentBranch: parent})
		}
		list.Add(&steps.DeleteParentBranchStep{Branch: branch})
	}
	if repo.Config.IsPerennialBranch(branch) {
		list.Add(&steps.RemoveFromPerennialBranchesStep{Branch: branch})
	}
	list.Add(&steps.DeleteLocalBranchStep{Branch: branch})
}

//nolint:nestif
func updateBranchSteps(list *runstate.StepListBuilder, branch string, pushBranch bool, branchesWithDeletedRemote []string, repo *git.ProdRepo) {
	isFeatureBranch := repo.Config.IsFeatureBranch(branch)
	syncStrategy := repo.Config.SyncStrategy()
	hasOrigin := list.Bool(repo.Silent.HasOrigin())
	pushHook := list.Bool(repo.Config.PushHook())
	isOffline := list.Bool(repo.Config.IsOffline())
	if !hasOrigin && !isFeatureBranch {
		return
	}
	list.Add(&steps.CheckoutStep{Branch: branch})
	if isFeatureBranch {
		updateFeatureBranchSteps(list, branch, branchesWithDeletedRemote, repo)
	} else {
		updatePerennialBranchSteps(list, branch, repo)
	}
	if pushBranch && hasOrigin && !isOffline {
		hasTrackingBranch := list.Bool(repo.Silent.HasTrackingBranch(branch))
		if !hasTrackingBranch {
			list.Add(&steps.CreateTrackingBranchStep{Branch: branch})
			return
		}
		if !isFeatureBranch {
			list.Add(&steps.PushBranchStep{Branch: branch})
			return
		}
		pushFeatureBranchSteps(list, branch, syncStrategy, pushHook)
	}
}

func updateFeatureBranchSteps(list *runstate.StepListBuilder, branch string, branchesWithDeletedRemote []string, repo *git.ProdRepo) {
	syncStrategy := repo.Config.SyncStrategy()
	hasTrackingBranch := list.Bool(repo.Silent.HasTrackingBranch(branch))
	if hasTrackingBranch {
		syncTrackingBranchSteps(list, repo.Silent.TrackingBranch(branch), syncStrategy)
	}
	// TODO: the last non-deleted parent branch here
	ancestorBranches := repo.Config.AncestorBranches(branch)
	ancestorBranches = stringslice.RemoveMany(ancestorBranches, branchesWithDeletedRemote)
	newParentBranch := stringslice.Last(ancestorBranches)
	if newParentBranch == nil {
		return
	}
	syncParentSteps(list, repo.Config.ParentBranch(branch), syncStrategy)
}

func updatePerennialBranchSteps(list *runstate.StepListBuilder, branch string, repo *git.ProdRepo) {
	hasTrackingBranch := list.Bool(repo.Silent.HasTrackingBranch(branch))
	if hasTrackingBranch {
		syncTrackingBranchSteps(list, repo.Silent.TrackingBranch(branch), repo.Config.PullBranchStrategy())
	}
	mainBranch := repo.Config.MainBranch()
	hasUpstream := list.Bool(repo.Silent.HasRemote("upstream"))
	shouldSyncUpstream := list.Bool(repo.Config.ShouldSyncUpstream())
	if mainBranch == branch && hasUpstream && shouldSyncUpstream {
		list.Add(&steps.FetchUpstreamStep{Branch: mainBranch})
		list.Add(&steps.RebaseBranchStep{Branch: fmt.Sprintf("upstream/%s", mainBranch)})
	}
}

// syncTrackingBranchStep provides the steps to sync the given tracking branch into the current branch.
func syncTrackingBranchSteps(list *runstate.StepListBuilder, trackingBranch, syncStrategy string) {
	switch syncStrategy {
	case "merge":
		list.Add(&steps.MergeStep{Branch: trackingBranch})
	case "rebase":
		list.Add(&steps.RebaseBranchStep{Branch: trackingBranch})
	default:
		list.Fail("unknown syncStrategy value: %q", syncStrategy)
	}
}

// syncParentSteps provides the steps to sync the given parent branch into the current branch.
func syncParentSteps(list *runstate.StepListBuilder, parentBranch, syncStrategy string) {
	switch syncStrategy {
	case "merge":
		list.Add(&steps.MergeStep{Branch: parentBranch})
	case "rebase":
		list.Add(&steps.RebaseBranchStep{Branch: parentBranch})
	default:
		list.Fail("unknown syncStrategy value: %q", syncStrategy)
	}
}

func pushFeatureBranchSteps(list *runstate.StepListBuilder, branch, syncStrategy string, pushHook bool) {
	switch syncStrategy {
	case "merge":
		list.Add(&steps.PushBranchStep{Branch: branch, NoPushHook: !pushHook})
	case "rebase":
		list.Add(&steps.PushBranchStep{Branch: branch, ForceWithLease: true})
	default:
		list.Fail("unknown syncStrategy value: %q", syncStrategy)
	}
}
