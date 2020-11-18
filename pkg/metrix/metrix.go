package metrix

import (
	"fmt"
	"log"

	"github.com/sk000f/metrix/pkg/collector/gitlab"
)

// Start initialises and configures the application
func Start(cfg *Config) {
	fmt.Println(cfg.GitlabURL)
	fmt.Println(cfg.GitlabToken)

	client, err := gitlab.SetupClient(cfg.GitlabToken, cfg.GitlabURL)
	if err != nil {
		log.Fatal(err)
	}

	projects, _, err := client.Projects.ListProjects(nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(projects)
}

// Config stores configuration values
type Config struct {
	GitlabURL   string
	GitlabToken string
}
