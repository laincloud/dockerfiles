package cmd

import (
	"archive/tar"
	"bytes"
	"context"
	"errors"
	"github.com/docker/docker/api"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/laincloud/dockerfiles/src/core"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
)

var (
	oldRegistryHost  string
	oldOrganization  string
	newRegistryHost  string
	newOrganization  string
	aptMirrorHost    string
	dockerHost       string
	dockerApiVersion string
)

var retagCmd = &cobra.Command{
	Use:   "retag",
	Short: "retag images from ${oldRegistryHost}/${oldOrganization}/${repository}:${tag} to ${newRegistryHost}/${newOrganization}/${repository}:${tag}",
	RunE:  retag,
}

func init() {
	retagCmd.Flags().StringVar(&commit1, "commit1", "origin/master", "previous commit")
	retagCmd.Flags().StringVar(&commit2, "commit2", "HEAD", "current commit")
	retagCmd.Flags().StringVar(&oldRegistryHost, "old-registry-host", "docker.io", "the old registry host who serves this image")
	retagCmd.Flags().StringVar(&oldOrganization, "old-organization", "laincloud", "the old organization build this image")
	retagCmd.Flags().StringVar(&newRegistryHost, "new-registry-host", "", "the new registry host who serves this image")
	retagCmd.Flags().StringVar(&newOrganization, "new-organization", "", "the new organization build this image")
	retagCmd.Flags().StringVar(&aptMirrorHost, "apt-mirror-host", "", "apt mirror host")
	retagCmd.Flags().StringVar(&dockerHost, "docker-host", "/var/run/docker.sock", "docker host")
	retagCmd.Flags().StringVar(&dockerApiVersion, "docker-api-version", "", "docker api version")
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

	if dockerApiVersion == "" {
		dockerApiVersion = api.DefaultVersion
	}

	if oldRegistryHost == newRegistryHost && oldOrganization == newOrganization {
		return errors.New("old-registry-host == new-registry-host && old-organization == new-organization")
	}

	dockerClient, err := client.NewClient(client.DefaultDockerHost, dockerApiVersion, nil, nil)
	if err != nil {
		return err
	}

	diffFiles, err := util.Diff(commit1, commit2)
	if err != nil {
		return err
	}

	diffImages, err := util.GetContext2Images(diffFiles)
	if err != nil {
		return err
	}

	for _, image := range diffImages {
		for _, tag := range image.Tags {
			log.Println("retag " + image.Repository + ":" + tag)
			if aptMirrorHost != "" {
				dockerfile := "FROM " + oldRegistryHost + "/" + oldOrganization + "/" + image.Repository + ":" + tag +
					"\nRUN sed -i 's|deb.debian.org|" + aptMirrorHost + "|g' /etc/apt/sources.list && sed -i 's|security.debian.org|" + aptMirrorHost + "/debian-security|g' /etc/apt/sources.list"

				buf := new(bytes.Buffer)

				tw := tar.NewWriter(buf)

				hdr := &tar.Header{
					Name: "Dockerfile",
					Mode: 0600,
					Size: int64(len(dockerfile)),
				}
				if err := tw.WriteHeader(hdr); err != nil {
					return err
				}
				if _, err := tw.Write([]byte(dockerfile)); err != nil {
					return err
				}
				if err := tw.Close(); err != nil {
					return err
				}
				buildOptions := types.ImageBuildOptions{Tags: []string{newRegistryHost + "/" + newOrganization + "/" + image.Repository + ":" + tag}}
				buildResponse, err := dockerClient.ImageBuild(context.Background(), buf, buildOptions)
				if err != nil {
					return err
				}
				response, err := ioutil.ReadAll(buildResponse.Body)
				if err != nil {
					return err
				}
				if err := buildResponse.Body.Close(); err != nil {
					return err
				}
				log.Println(string(response))
			} else {
				err := dockerClient.ImageTag(context.Background(), oldRegistryHost+"/"+oldOrganization+"/"+image.Repository+":"+tag, newRegistryHost+"/"+newOrganization+"/"+image.Repository+":"+tag)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
