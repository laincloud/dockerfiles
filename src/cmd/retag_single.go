package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"text/template"

	"github.com/spf13/cobra"
)

const (
	dockerfileTemplate = `
FROM {{.OldRegistryHost}}/{{.OldOrganization}}/{{.RepositoryAndTag}}
RUN sed -i 's|deb.debian.org|{{.AptMirrorHost}}|g' /etc/apt/sources.list && \
    sed -i 's|security.debian.org|{{.AptMirrorHost}}/debian-security|g' /etc/apt/sources.list
`
)

type DockerfileData struct {
	OldRegistryHost  string
	OldOrganization  string
	RepositoryAndTag string
	AptMirrorHost    string
}

var retagSingleCmd = &cobra.Command{
	Use:   "retag-single",
	Short: "retag an image from ${oldRegistryHost}/${oldOrganization}/${repository}:${tag} to ${newRegistryHost}/${newOrganization}/${repository}:${tag}",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("exact one argument is required(format: `repository:tag`)")
		}

		return nil
	},
	RunE: retagSingle,
}

func init() {
	retagSingleCmd.Flags().StringVar(&oldRegistryHost, "old-registry-host", "docker.io", "the old registry host who serves this image")
	retagSingleCmd.Flags().StringVar(&oldOrganization, "old-organization", "laincloud", "the old organization build this image")
	retagSingleCmd.Flags().StringVar(&newRegistryHost, "new-registry-host", "", "the new registry host who serves this image")
	retagSingleCmd.Flags().StringVar(&newOrganization, "new-organization", "", "the new organization build this image")
	retagSingleCmd.Flags().StringVar(&aptMirrorHost, "apt-mirror-host", "", "apt mirror host")
	rootCmd.AddCommand(retagSingleCmd)
}

func retagSingle(cmd *cobra.Command, args []string) error {
	if oldRegistryHost == "" {
		return errors.New("--old-registry-host is required")
	}

	if oldOrganization == "" {
		return errors.New("--old-organization is required")
	}

	if newRegistryHost == "" {
		return errors.New("--new-registry-host is required")
	}

	if newOrganization == "" {
		return errors.New("--new-organization is required")
	}

	if newRegistryHost == oldRegistryHost && newOrganization == oldOrganization {
		return errors.New("new-registry-host == old-registry-host && new-organization == old-organization")
	}

	if aptMirrorHost != "" {
		t := template.Must(template.New("dockerfile").Parse(dockerfileTemplate))
		var buf bytes.Buffer
		if err := t.Execute(&buf, DockerfileData{
			OldRegistryHost:  oldRegistryHost,
			OldOrganization:  oldOrganization,
			RepositoryAndTag: args[0],
			AptMirrorHost:    aptMirrorHost,
		}); err != nil {
			return err
		}

		dockerBuild := exec.Command("docker", "build", "-t", fmt.Sprintf("%s/%s/%s", newRegistryHost, newOrganization, args[0]), "-")
		dockerBuild.Stdin = &buf
		dockerBuild.Stdout = os.Stdout
		dockerBuild.Stderr = os.Stderr
		return dockerBuild.Run()
	}

	dockerRetag := exec.Command("docker", "tag", fmt.Sprintf("%s/%s/%s", oldRegistryHost, oldOrganization, args[0]), fmt.Sprintf("%s/%s/%s", newRegistryHost, newOrganization, args[0]))
	dockerRetag.Stdout = os.Stdout
	dockerRetag.Stderr = os.Stderr
	return dockerRetag.Run()
}
