package metrix

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/sk000f/metrix/pkg/collector"
	"github.com/sk000f/metrix/pkg/collector/gitlab"
	"github.com/sk000f/metrix/pkg/storage/mongo"
)

// Start initialises and configures the application
func Start() error {

	cfg := SetupConfig()

	gl := &gitlab.GitLab{
		Token: cfg.GitLabToken,
		URL:   cfg.GitLabURL,
	}

	r := new(mongo.DB)
	r.ConnStr = cfg.DBConnString

	c := collector.NewService(gl, r)

	err := c.RefreshData()
	if err != nil {
		fmt.Printf("Error: %v", err.Error())
		return err
	}

	return nil
}

// SetupConfig configures application based on environment variables
func SetupConfig() *Config {

	cfg := new(Config)

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	cfg.GitLabURL = os.Getenv("METRIX_GITLAB_URL")
	cfg.GitLabToken = os.Getenv("METRIX_GITLAB_TOKEN")
	cfg.DBConnString = os.Getenv("METRIX_DB_CONN_STRING")

	return cfg
}

// Config stores configuration values
type Config struct {
	GitLabURL    string
	GitLabToken  string
	DBConnString string
}
