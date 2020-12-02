package gitlab

import (
	"fmt"

	"github.com/sk000f/metrix/pkg/collector"
	gl "github.com/xanzy/go-gitlab"
)

// GitLab represents a GitLab server
type GitLab struct {
	Token string
	URL   string
}

// RefreshData gets latest deployment data from CI server and saves to repository
func (g *GitLab) RefreshData(r collector.Repository) error {

	c, err := g.SetupClient(g.Token, g.URL)
	if err != nil {
		fmt.Printf("Error: %v", err.Error())
		return err
	}

	p := g.UpdateProjects(c, r)

	g.UpdateDeployments(p, c, r)

	return nil
}

// UpdateProjects gets all projects from GitLab and stores them in the repository
func (g *GitLab) UpdateProjects(c *gl.Client, r collector.Repository) []*collector.Project {

	// get all projects
	p, err := g.GetProjects(c, getProjectListOptions())
	if err != nil {
		fmt.Printf("Error: %v", err.Error())
	}

	// save projects to repository
	r.SaveProjects(p)

	return p
}

// UpdateDeployments gets deployments for all projects from GitLab and stores them in the repository
func (g *GitLab) UpdateDeployments(p []*collector.Project, c *gl.Client, r collector.Repository) {

	// get and update deployments for projects
	for _, proj := range p {
		d, _ := g.GetDeployments(proj, c, getDeploymentListOptions())
		for _, dep := range d {
			r.SaveDeployment(dep)
		}
	}
}

// GetProjects lists all projects from specified GitLab server
func (g *GitLab) GetProjects(client *gl.Client, opt *gl.ListProjectsOptions) ([]*collector.Project, error) {

	p := []*collector.Project{}

	for {

		projects, resp, err := client.Projects.ListProjects(opt)
		if err != nil {
			fmt.Printf("Error: %v", err.Error())
			return nil, err
		}

		// iterate over projects and convert to metrix representation
		for _, pr := range projects {
			p = append(p, &collector.Project{
				ID:                pr.ID,
				Name:              pr.Name,
				Path:              pr.Path,
				PathWithNamespace: pr.PathWithNamespace,
				Namespace:         pr.Namespace.FullPath,
				WebURL:            pr.WebURL,
			})
		}

		// Exit the loop when we've seen all pages.
		if resp.CurrentPage >= resp.TotalPages {
			break
		}

		opt.Page = resp.NextPage
	}

	return p, nil
}

// GetDeployments lists all Deployments for the specified Project
func (g *GitLab) GetDeployments(p *collector.Project, client *gl.Client, opt *gl.ListProjectDeploymentsOptions) ([]*collector.Deployment, error) {

	d := []*collector.Deployment{}

	for {

		deployments, resp, err := client.Deployments.ListProjectDeployments(p.ID, opt)
		if err != nil {
			fmt.Printf("Error: %v", err.Error())
			return nil, err
		}

		// iterate over deployments and convert to metrix representation
		for _, dep := range deployments {

			if dep.Environment.Name == "production" &&
				(dep.Status == "success" || dep.Status == "failed") {

				d = append(d, &collector.Deployment{
					ID:               dep.ID,
					Status:           dep.Status,
					EnvironmentName:  dep.Environment.Name,
					ProjectID:        p.ID,
					ProjectName:      p.Name,
					ProjectPath:      p.Path,
					ProjectNamespace: p.Namespace,
					PipelineID:       dep.Deployable.Pipeline.ID,
					FinishedAt:       dep.Deployable.FinishedAt,
					Duration:         dep.Deployable.Duration,
				})
			}

		}

		// Exit the loop when we've seen all pages.
		if resp.CurrentPage >= resp.TotalPages {
			break
		}

		opt.Page = resp.NextPage
	}

	return d, nil
}

// SetupClient returns a GitLab client with the specified base URL
func (g *GitLab) SetupClient(token, baseURL string) (*gl.Client, error) {
	client, err := gl.NewClient(token, gl.WithBaseURL(baseURL))
	if err != nil {
		fmt.Printf("Error: %v", err.Error())
		return nil, err
	}

	return client, nil
}

func getProjectListOptions() *gl.ListProjectsOptions {
	return &gl.ListProjectsOptions{
		ListOptions: gl.ListOptions{Page: 1, PerPage: 20},
		Simple:      gl.Bool(false),
	}
}

func getDeploymentListOptions() *gl.ListProjectDeploymentsOptions {
	return &gl.ListProjectDeploymentsOptions{
		ListOptions: gl.ListOptions{Page: 1, PerPage: 20},
		Environment: gl.String("production"),
	}
}
