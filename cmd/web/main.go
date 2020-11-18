package main

import (
	"os"

	"github.com/sk000f/metrix/pkg/metrix"
)

func main() {

	cfg := new(metrix.Config)

	cfg.GitlabURL = os.Getenv("METRIX_GITLAB_URL")
	cfg.GitlabToken = os.Getenv("METRIX_GITLAB_TOKEN")

	metrix.Start(cfg)
}
