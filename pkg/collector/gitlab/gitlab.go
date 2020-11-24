package gitlab

import (
	"fmt"

	gl "github.com/xanzy/go-gitlab"
)

// GetProjects lists all projects from specified GitLab server
func GetProjects(client *gl.Client, opt *gl.ListProjectsOptions) ([]*Project, error) {

	p := []*Project{}

	for {
		projects, resp, err := client.Projects.ListProjects(opt)
		if err != nil {
			return nil, err
		}

		// iterate over projects and convert to metrix representation
		for _, pr := range projects {
			p = append(p, &Project{
				ID:                pr.ID,
				Name:              pr.Name,
				NameWithNamespace: pr.NameWithNamespace,
				WebURL:            pr.WebURL,
			})
		}

		fmt.Printf("Page: %v, TotalPages: %v, NextPage: %v", opt.Page, resp.TotalPages, resp.NextPage)

		// Exit the loop when we've seen all pages.
		if resp.CurrentPage >= resp.TotalPages {
			break
		}

		opt.Page = resp.NextPage
	}

	return p, nil
}

// GetDeployments lists all Deployments for the specified Project
func GetDeployments(pid int, client *gl.Client, opt *gl.ListProjectDeploymentsOptions) ([]*Deployment, error) {
	d := []*Deployment{}

	for {
		deployments, resp, err := client.Deployments.ListProjectDeployments(pid, opt)
		if err != nil {
			return nil, err
		}

		// iterate over deployments and convert to metrix representation
		for _, dep := range deployments {
			d = append(d, &Deployment{
				ID:              dep.ID,
				Status:          dep.Deployable.Status,
				EnvironmentName: dep.Environment.Name,
				PipelineID:      dep.Deployable.Pipeline.ID,
			})

			// Exit the loop when we've seen all pages.
			if resp.CurrentPage >= resp.TotalPages {
				break
			}

			opt.Page = resp.NextPage
		}

		return d, nil
	}

}

// SetupClient returns a GitLab client with the specified base URL
func SetupClient(token, baseURL string) (*gl.Client, error) {
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

// Deployment represents metrix view of a GitLab deployment object
type Deployment struct {
	ID              int
	Status          string
	EnvironmentName string
	PipelineID      int
}
