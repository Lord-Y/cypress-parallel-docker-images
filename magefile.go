//go:build mage
// +build mage

package main

import (
	"bytes"
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
			cypress: "7.4.0",
			cli:     "v0.1.1",
		},
		{
			cypress: "7.4.0",
			cli:     "v0.1.0",
		},
		{
			cypress: "7.4.0",
			cli:     "v0.0.5",
		},
		{
			cypress: "7.3.0",
			cli:     "v0.0.5",
		},
		{
			cypress: "7.2.0",
			cli:     "v0.0.5",
		},
		{
			cypress: "7.2.0",
			cli:     "v0.0.4",
		},
		{
			cypress: "7.2.0",
			cli:     "v0.0.3",
		},
		{
			cypress: "7.2.0",
			cli:     "v0.0.2",
		},
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

// Build Build and publish docker images if doesn't exist
func Build() (err error) {
	mg.Deps(InstallDeps)
	var tags []string
	z, err := docker.GetDockerImages()
	if err != nil {
		log.Error().Err(err).Msg("Fail to retrieve docker image list")
		return
	}
	for _, tag := range z {
		tags = append(tags, tag.Metadata.Container.Tags[0])
	}
	log.Info().Msgf("Actual tags %+v", tags)

	if len(tags) > 0 {
		for _, image := range images {
			var m buildImage
			m.image = image
			if !tools.StringInSlice(
				fmt.Sprintf("%s-%s", m.image.cypress, strings.TrimPrefix(m.image.cli, "v")),
				tags,
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
		var cmdOutput bytes.Buffer
		cmd1 := exec.Command(
			"echo",
			strings.TrimSpace(os.Getenv("GITHUB_TOKEN")),
		)
		cmd2 := exec.Command(
			"sudo",
			"docker",
			"login",
			"https://docker.pkg.github.com",
			"-u",
			docker.Owner,
			"--password-stdin",
		)
		cmd2.Stdin, _ = cmd1.StdoutPipe()
		cmd2.Stdout = &cmdOutput
		err = cmd2.Start()
		if err != nil {
			log.Error().Err(err).Msg("Fail to run cmd2")
			return
		}
		err = cmd1.Run()
		if err != nil {
			log.Error().Err(err).Msg("Fail to run cmd1")
			return
		}
		err = cmd2.Wait()
		if err != nil {
			log.Error().Err(err).Msg("Fail to login to github docker registry")
			return
		}
		log.Info().Msgf("%s", string(cmdOutput.Bytes()))

		output, err = exec.Command(
			"sudo",
			"docker",
			"push",
			fmt.Sprintf("%s:%s-%s", ghr, m.image.cypress, strings.TrimPrefix(m.image.cli, "v")),
		).Output()
		log.Info().Msgf("%s", output)
		if err != nil {
			return
		}
	}
	return
}

// Install dependencies
func InstallDeps() error {
	os.Setenv("GO111MODULE", "on")
	fmt.Println("Installing Deps...")
	cmd := exec.Command("go", "mod", "download")
	return cmd.Run()
}
