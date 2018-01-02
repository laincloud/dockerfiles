package cmd

import (
	"errors"

	"github.com/laincloud/dockerfiles/util"
	"github.com/spf13/cobra"
)

var (
	oldRegistryHost string
	oldOrganization string
	newRegistryHost string
	newOrganization string
)

var retagCmd = &cobra.Command{
	Use:   "retag",
	Short: "retag images from ${oldRegistryHost}/${oldOrganization}/${repository}:${tag} to ${newRegistryHost}/${newOrganization}/${repository}:${tag}",
	RunE:  retag,
}

func init() {
	retagCmd.Flags().StringVar(&commit1, "commit1", "origin/master", "previous commit")
	retagCmd.Flags().StringVar(&commit2, "commit2", "HEAD", "current commit")
	retagCmd.Flags().StringVar(&oldRegistryHost, "old-registry-host", "", "the old registry host who serves this image")
	retagCmd.Flags().StringVar(&oldOrganization, "old-organization", "laincloud", "the old organization build this image")
	retagCmd.Flags().StringVar(&newRegistryHost, "new-registry-host", "", "the new registry host who serves this image")
	retagCmd.Flags().StringVar(&newOrganization, "new-organization", "", "the new organization build this image")
	rootCmd.AddCommand(retagCmd)
}

func retag(cmd *cobra.Command, args []string) error {
	if commit1 == "" {
		return errors.New("--commit1 is required")
	}

	if commit2 == "" {
		return errors.New("--commit2 is required")
	}

	if oldOrganization == "" {
		return errors.New("--old-organization is required")
	}

	if newOrganization == "" {
		return errors.New("--new-organization is required")
	}

	if oldRegistryHost == newRegistryHost && oldOrganization == newOrganization {
		return errors.New("old-registry-host == new-registry-host && old-organization == new-organization")
	}

	return util.Make(util.Args{
		Command:         util.Retag,
		Commit1:         commit1,
		Commit2:         commit2,
		NewOrganization: newOrganization,
		NewRegistryHost: newRegistryHost,
		OldOrganization: oldOrganization,
		OldRegistryHost: oldRegistryHost,
	})
}
