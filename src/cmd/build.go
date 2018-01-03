package cmd

import (
	"errors"

	"github.com/laincloud/dockerfiles/src/core"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "build docker images",
	RunE:  build,
}

func init() {
	buildCmd.Flags().StringVar(&commit1, "commit1", "origin/master", "previous commit")
	buildCmd.Flags().StringVar(&commit2, "commit2", "HEAD", "current commit")
	rootCmd.AddCommand(buildCmd)
}

func build(cmd *cobra.Command, args []string) error {
	if commit1 == "" {
		return errors.New("--commit1 is required")
	}

	if commit2 == "" {
		return errors.New("--commit2 is required")
	}

	return util.Make(util.Args{
		Command: util.Build,
		Commit1: commit1,
		Commit2: commit2,
	})
}
