package config

import "fmt"

type PullBranchStrategy string

const (
	PullBranchStrategyMerge  = "merge"
	PullBranchStrategyRebase = "rebase"
)

func ToPullBranchStrategy(text string) (PullBranchStrategy, error) {
	switch text {
	case "merge":
		return PullBranchStrategyMerge, nil
	case "rebase", "":
		return PullBranchStrategyRebase, nil
	default:
		return PullBranchStrategyMerge, fmt.Errorf("unknown pull branch strategy: %q", text)
	}
}

func (pbs PullBranchStrategy) String() string {
	return string(pbs)
}
