package gitlab

import (
	gl "github.com/xanzy/go-gitlab"
)

// GetProjects lists all projects from specified GitLab server
func GetProjects(client *gl.Client, opt *gl.ListProjectsOptions) ([]*Project, error) {

	projects, _, err := client.Projects.ListProjects(opt)
	if err != nil {
		return nil, err
	}

	// iterate over projects and convert to metrix representation
	p := []*Project{}
	for _, pr := range projects {
		p = append(p, &Project{
			ID:                pr.ID,
			Name:              pr.Name,
			NameWithNamespace: pr.NameWithNamespace,
			WebURL:            pr.WebURL,
		})
	}

	return p, nil
}

// SetupClient returns a GitLab client with the specified base URL
func SetupClient(baseURL, token string) (*gl.Client, error) {
	client, err := gl.NewClient(token, gl.WithBaseURL(baseURL))
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Project represents metrix view of a GitLab project object
type Project struct {
	ID                int
	Name              string
	NameWithNamespace string
	WebURL            string
}
