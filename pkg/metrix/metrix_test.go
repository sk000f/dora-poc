package metrix_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/sk000f/metrix/pkg/collector"
	"github.com/sk000f/metrix/pkg/metrix"
)

func TestMetrix(t *testing.T) {
	t.Run("configuration values set correctly", func(t *testing.T) {

		os.Setenv("METRIX_ENV", "dev")
		os.Setenv("METRIX_GITLAB_URL", "https://example.com")
		os.Setenv("METRIX_GITLAB_TOKEN", "1234567890")

		want := &metrix.Config{GitLabURL: "https://example.com", GitLabToken: "1234567890"}
		got := metrix.SetupConfig()

		if !reflect.DeepEqual(got, want) {
			t.Errorf("want %v; got %v", want, got)
		}

		os.Unsetenv("METRIX_ENV")
		os.Unsetenv("METRIX_GITLAB_URL")
		os.Unsetenv("METRIX_GITLAB_TOKEN")
	})

	t.Run("application starts and executes correctly", func(t *testing.T) {

		os.Setenv("METRIX_ENV", "dev")
		os.Setenv("METRIX_GITLAB_URL", "https://example.com")
		os.Setenv("METRIX_GITLAB_TOKEN", "1234567890")

		err := metrix.Start()

		if err != nil {
			t.Errorf("Unexpected error thrown: %v", err.Error())
		}

		os.Unsetenv("METRIX_ENV")
		os.Unsetenv("METRIX_GITLAB_URL")
		os.Unsetenv("METRIX_GITLAB_TOKEN")
	})
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
