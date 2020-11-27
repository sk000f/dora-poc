package metrix

import (
	"os"

	"github.com/sk000f/metrix/pkg/collector"

	"github.com/sk000f/metrix/pkg/collector/gitlab"
)

// Start initialises and configures the application
func Start() {

	cfg := SetupConfig()

	gl := &gitlab.GitLab{
		Token: cfg.GitLabToken,
		URL:   cfg.GitLabURL,
	}

	r := new(mockRepo)

	c := collector.NewService(gl, r)

	c.RefreshData()
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

type mockRepo struct {
	ProjectData    []*collector.Project
	DeploymentData []*collector.Deployment
}

func (m *mockRepo) SaveProjects(p []*collector.Project) {
	for _, proj := range p {
		m.ProjectData = append(m.ProjectData, proj)
	}
}

func (m *mockRepo) SaveDeployment(d *collector.Deployment) {
	m.DeploymentData = append(m.DeploymentData, d)
}
