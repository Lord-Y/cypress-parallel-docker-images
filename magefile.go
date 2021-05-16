// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Lord-Y/cypress-parallel-docker-images/docker"
	"github.com/Lord-Y/golang-tools/logger"
	"github.com/Lord-Y/golang-tools/tools"
	"github.com/magefile/mage/mg"
	"github.com/rs/zerolog/log"
)

var (
	images = []image{
		{
			cypress: "7.2.0",
			cli:     "v0.0.1",
		},
	}
	ghr = fmt.Sprintf("docker.pkg.github.com/%s/%s/%s", strings.ToLower(docker.Owner), docker.Repository, docker.Repository)
)

type image struct {
	cypress string
	cli     string
}

type buildImage struct {
	image   image
	publish bool
}

func init() {
	os.Setenv("LOGGER_TYPE", "shell")
	logger.SetLoggerLogLevel()
}

// Build and publish docker images if doesn't exist
func Build() (err error) {
	mg.Deps(InstallDeps)
	z, err := docker.GetDockerImages()
	if err != nil {
		log.Error().Err(err).Msg("Fail to retrieve docker image list")
		return
	}
	if len(z) > 0 {
		for _, image := range images {
			var m buildImage
			m.image = image
			if !tools.StringInSlice(
				fmt.Sprintf("%s:%s-%s", ghr, m.image.cypress, strings.TrimPrefix(m.image.cli, "v")),
				z[0].Metadata.Docker.Tags,
			) {
				if strings.TrimSpace(os.Getenv("PUBLISH_DOCKER_IMAGES")) != "" {
					m.publish = true
				}
				err = m.build()
			}
		}
		return
	}
	for _, image := range images {
		var m buildImage
		m.image = image
		if strings.TrimSpace(os.Getenv("PUBLISH_DOCKER_IMAGES")) != "" {
			m.publish = true
		}
		err = m.build()
	}
	return
}

func (m *buildImage) build() (err error) {
	output, err := exec.Command(
		"sudo",
		"docker",
		"build",
		"--build-arg",
		fmt.Sprintf(`CYPRESS_DOCKER_IMAGE_VERSION=%s`, m.image.cypress),
		"--build-arg",
		fmt.Sprintf(`CYPRESS_PARALLEL_CLI=%s`, m.image.cli),
		"-t",
		fmt.Sprintf("%s:%s-%s", ghr, m.image.cypress, strings.TrimPrefix(m.image.cli, "v")),
		".",
	).Output()
	log.Info().Msgf("%s", output)
	if err != nil {
		return
	}
	if m.publish {
		output, err = exec.Command(
			"echo",
			strings.TrimSpace(os.Getenv("GITHUB_TOKEN")),
			"|",
			"sudo",
			"docker",
			"login",
			"https://docker.pkg.github.com",
			"-u",
			docker.Owner,
			"--password-stdin",
		).Output()
		log.Info().Msgf("%s", output)
		if err != nil {
			return
		}

		output, err = exec.Command(
			"sudo",
			"docker",
			"push",
			fmt.Sprintf("%s:%s-%s", ghr, m.image.cypress, strings.TrimPrefix(m.image.cli, "v")),
		).Output()
		log.Info().Msgf("%s", output)
	}
	return
}

// Manage your deps, or running package managers.
func InstallDeps() error {
	os.Setenv("GO111MODULE", "on")
	fmt.Println("Installing Deps...")
	cmd := exec.Command("go", "mod", "download")
	return cmd.Run()
}
