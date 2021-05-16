// Package docker permit to retrieve docker images list
package docker

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Lord-Y/golang-tools/httprequests"
	"github.com/rs/zerolog/log"
)

// https://docs.github.com/en/rest/reference/packages

var (
	Owner       = "Lord-Y"
	Repository  = "cypress-parallel-docker-images"
	githubToken = strings.TrimSpace(os.Getenv("GITHUB_TOKEN"))
)

type Response struct {
	ID             int    `form:"id" json:"workload"`
	Name           string `form:"name" json:"name"`
	Url            string `form:"url" json:"url"`
	PackageHtmlURL string `form:"package_html_url" json:"package_html_url"`
	License        string `form:"license" json:"license"`
	CreatedAT      string `form:"created_at" json:"created_at"`
	UpdatedAT      string `form:"updated_at" json:"updated_at"`
	HtmlURL        string `form:"html_url" json:"html_url"`
	Metadata       metadata
}

type metadata struct {
	PackageType string `form:"package_type" json:"package_type"`
	Docker      docker
}

type docker struct {
	Tags []string `form:"tags" json:"tags"`
}

func GetDockerImages() (gr []Response, err error) {
	headers := make(map[string]string)
	headers["Accept"] = "application/vnd.github.v3+json"
	headers["Authorization"] = fmt.Sprintf("bearer %s", githubToken)

	body, resp, err := httprequests.PerformRequests(headers, "GET", fmt.Sprintf("https://api.github.com/user/packages/docker/%s/versions", Repository), "", "")
	if err != nil {
		log.Error().Err(err).Msgf("Fail to retrieve docker images informations, status code %d", resp.StatusCode)
		return
	}
	if resp.StatusCode == 200 {
		err = json.Unmarshal(body, &gr)
	}
	return
}
