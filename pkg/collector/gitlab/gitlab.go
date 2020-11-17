package gitlab

import (
	"github.com/xanzy/go-gitlab"
)

// SetupClient returns a GitLab client with the specified base URL
func SetupClient(baseURL string) (*gitlab.Client, error) {
	client, err := gitlab.NewClient("", gitlab.WithBaseURL(baseURL))
	if err != nil {
		return nil, err
	}

	return client, nil
}
