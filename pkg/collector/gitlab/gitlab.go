package gitlab

import (
	"github.com/sk000f/metrix/pkg/collector"
	"github.com/xanzy/go-gitlab"
	gl "github.com/xanzy/go-gitlab"
)

// RefreshData gets latest deployment data from CI server and saves to repository
func RefreshData(c *gl.Client, r collector.Repository) {

	p := UpdateProjects(c, r)

	UpdateDeployments(p, c, r)

}

// UpdateProjects gets all projects from GitLab and stores them in the repository
func UpdateProjects(c *gitlab.Client, r collector.Repository) []*collector.Project {

	// get all projects
	p, err := GetProjects(c, nil)
	if err != nil {

	}

	// save projects to repository
	r.SaveProjects(p)

	return p
}

// UpdateDeployments gets deployments for all projects from GitLab and stores them in the repository
func UpdateDeployments(p []*collector.Project, c *gitlab.Client, r collector.Repository) {

	// get and update deployments for projects
	for _, proj := range p {
		d, _ := GetDeployments(proj.ID, c, nil)
		for _, dep := range d {
			r.SaveDeployment(dep)
		}
	}
}

// GetProjects lists all projects from specified GitLab server
func GetProjects(client *gl.Client, opt *gl.ListProjectsOptions) ([]*collector.Project, error) {

	p := []*collector.Project{}

	for {
		projects, resp, err := client.Projects.ListProjects(opt)
		if err != nil {
			return nil, err
		}

		// iterate over projects and convert to metrix representation
		for _, pr := range projects {
			p = append(p, &collector.Project{
				ID:                pr.ID,
				Name:              pr.Name,
				NameWithNamespace: pr.NameWithNamespace,
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
func GetDeployments(pid int, client *gl.Client, opt *gl.ListProjectDeploymentsOptions) ([]*collector.Deployment, error) {
	d := []*collector.Deployment{}

	for {
		deployments, resp, err := client.Deployments.ListProjectDeployments(pid, opt)
		if err != nil {
			return nil, err
		}

		// iterate over deployments and convert to metrix representation
		for _, dep := range deployments {
			d = append(d, &collector.Deployment{
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
