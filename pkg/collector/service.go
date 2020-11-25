package collector

// Service provides functionality for updating CI data
type Service struct {
	ci CIServer
	r  Repository
}

// CIServer provides functionality for getting data from CI server
type CIServer interface {
	RefreshData(r Repository)
}

// Repository provides access to data storage
type Repository interface {
	SaveData()
}

// NewService creates a collector with required dependencies
func NewService(ci CIServer, r Repository) *Service {
	return &Service{ci, r}
}

// RefreshData collects data from CI server and saves in data repository
func (s *Service) RefreshData() {
	s.ci.RefreshData(s.r)
}
