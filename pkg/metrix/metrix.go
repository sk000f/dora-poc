package metrix

import (
	"fmt"
	"log"
	"os"

	"github.com/sk000f/metrix/pkg/collector/gitlab"
)

// Start initialises and configures the application
func Start() {

	cfg := SetupConfig()

	client, err := gitlab.SetupClient(cfg.GitLabToken, cfg.GitLabURL)
	if err != nil {
		log.Fatal(err)
	}

	projects, _, err := client.Projects.ListProjects(nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(projects)
}

// SetupConfig configures application based on environment variables
func SetupConfig() *Config {
	cfg := new(Config)

	cfg.GitLabURL = os.Getenv("METRIX_GITLAB_URL")
	cfg.GitLabToken = os.Getenv("METRIX_GITLAB_TOKEN")

	return cfg
}

// Config stores configuration values
type Config struct {
	GitLabURL   string
	GitLabToken string
}
