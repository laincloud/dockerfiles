package cmd

import (
	"errors"

	"github.com/laincloud/dockerfiles/src/core"
	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "pull docker images",
	RunE:  pull,
}

func init() {
	pullCmd.Flags().StringVar(&commit1, "commit1", "origin/master", "previous commit")
	pullCmd.Flags().StringVar(&commit2, "commit2", "HEAD", "current commit")
	pullCmd.Flags().StringVarP(&registryHost, "registry-host", "r", "", "the registry host who serves this image")
	pullCmd.Flags().StringVarP(&organization, "organization", "o", "laincloud", "the organization build this image")
	rootCmd.AddCommand(pullCmd)
}

func pull(cmd *cobra.Command, args []string) error {
	if commit1 == "" {
		return errors.New("--commit1 is required")
	}

	if commit2 == "" {
		return errors.New("--commit2 is required")
	}

	if organization == "" {
		return errors.New("--organization is required")
	}

	return util.Make(util.Args{
		Command:      util.Pull,
		Commit1:      commit1,
		Commit2:      commit2,
		Organization: organization,
		RegistryHost: registryHost,
	})
}
