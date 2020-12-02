package collector

import "fmt"

// Service provides functionality for updating CI data
type Service struct {
	ci CIServer
	r  Repository
}

// CIServer provides functionality for getting data from CI server
type CIServer interface {
	RefreshData(r Repository) error
}

// Repository provides access to data storage
type Repository interface {
	SaveProjects(p []*Project)
	SaveDeployment(d *Deployment)
}

// Project represents metrix view of a GitLab project object
type Project struct {
	ID                int
	Name              string
	PathWithNamespace string
	Namespace         string
	WebURL            string
}

// Deployment represents metrix view of a GitLab deployment object
type Deployment struct {
	ID              int
	Status          string
	EnvironmentName string
	PipelineID      int
}

// NewService creates a collector with required dependencies
func NewService(ci CIServer, r Repository) *Service {
	return &Service{ci, r}
}

// RefreshData collects data from CI server and saves in data repository
func (s *Service) RefreshData() error {
	err := s.ci.RefreshData(s.r)
	if err != nil {
		fmt.Printf("Error: %v", err.Error())
		return err
	}
	return nil
}
