package metrix

import (
	"fmt"
)

// Start initialises and configures the application
func Start(cfg *Config) {
	fmt.Println(cfg.GitlabURL)
	fmt.Println(cfg.GitlabToken)
}

// Config stores configuration values
type Config struct {
	GitlabURL   string
	GitlabToken string
}
