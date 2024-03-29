package gitlab_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/sk000f/metrix/pkg/collector"
	"github.com/sk000f/metrix/pkg/collector/gitlab"
	gl "github.com/xanzy/go-gitlab"
)

func TestGitLabProjects(t *testing.T) {
	t.Run("get single Project from GitLab", func(t *testing.T) {

		mux, server, client, g := setupMockGitLabClient(t)
		defer teardown(server)

		mux.HandleFunc("/api/v4/projects", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `[{
					"id": 1, 
					"name": "test",
					"path": "test",
					"path_with_namespace": "test/test", 
					"web_url": "http://test.com/test/test",
					"namespace" :{
						"full_path": "test/test"
					}
					}]`)
		})

		want := []*collector.Project{
			{
				ID:                1,
				Name:              "test",
				Path:              "test",
				PathWithNamespace: "test/test",
				Namespace:         "test/test",
				WebURL:            "http://test.com/test/test",
			}}

		got, err := g.GetProjects(client, getProjectListOptions())
		if err != nil {
			t.Errorf("Error getting Projects: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v; wanted %+v", got, want)
		}
	})

	t.Run("get multiple pages of Projects from GitLab", func(t *testing.T) {
		mux, server, client, g := setupMockGitLabClient(t)
		defer teardown(server)

		mux.HandleFunc("/api/v4/projects", func(w http.ResponseWriter, r *http.Request) {

			if r.URL.Query()["page"][0] == "1" {
				w.Header().Set("X-Page", "1")
				w.Header().Set("X-Total-Pages", "2")
				w.Header().Set("X-Next-Page", "2")
				fmt.Fprint(w, `[{
					"id": 1, 
					"name": "test", 
					"path": "test",
					"path_with_namespace": "test/test", 
					"web_url": "http://test.com/test/test",
					"namespace" :{
						"full_path": "test/test"
					}
				}
				]`)
			}

			if r.URL.Query()["page"][0] == "2" {
				w.Header().Set("X-Page", "2")
				w.Header().Set("X-Total-Pages", "2")
				w.Header().Set("X-Next-Page", "2")
				fmt.Fprint(w, `[
				{
					"id": 2, 
					"name": "test", 
					"path": "test",
					"path_with_namespace": "test/test", 
					"web_url": "http://test.com/test/test",
					"namespace" :{
						"full_path": "test/test"
					}
				},
				{
					"id": 3, 
					"name": "test", 
					"path": "test",
					"path_with_namespace": "test/test", 
					"web_url": "http://test.com/test/test",
					"namespace" :{
						"full_path": "test/test"
					}
				}
				]`)
			}
		})

		want := []*collector.Project{
			{
				ID:                1,
				Name:              "test",
				Path:              "test",
				PathWithNamespace: "test/test",
				Namespace:         "test/test",
				WebURL:            "http://test.com/test/test",
			},
			{
				ID:                2,
				Name:              "test",
				Path:              "test",
				PathWithNamespace: "test/test",
				Namespace:         "test/test",
				WebURL:            "http://test.com/test/test",
			},
			{
				ID:                3,
				Name:              "test",
				Path:              "test",
				PathWithNamespace: "test/test",
				Namespace:         "test/test",
				WebURL:            "http://test.com/test/test",
			},
		}

		opt := getProjectListOptions()

		got, err := g.GetProjects(client, opt)
		if err != nil {
			t.Errorf("Error getting Projects: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v; wanted %+v", got, want)
		}
	})

	t.Run("update projects in repository", func(t *testing.T) {

		mux, server, client, g := setupMockGitLabClient(t)
		defer teardown(server)

		mux.HandleFunc("/api/v4/projects", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `[{
					"id": 1, 
					"name": "test", 
					"path": "test",
					"path_with_namespace": "test/test", 
					"web_url": "http://test.com/test/test",
					"namespace" :{
						"full_path": "test/test"
					}
					}]`)
		})

		mockRepository := new(mockRepo)

		want := []*collector.Project{
			{
				ID:                1,
				Name:              "test",
				Path:              "test",
				PathWithNamespace: "test/test",
				Namespace:         "test/test",
				WebURL:            "http://test.com/test/test",
			}}

		g.UpdateProjects(client, mockRepository)

		if !reflect.DeepEqual(mockRepository.ProjectData, want) {
			t.Errorf("got %+v; wanted %+v", mockRepository.ProjectData, want)
		}
	})
}

func TestGitLabDeployments(t *testing.T) {
	t.Run("get all deployments from GitLab", func(t *testing.T) {
		mux, server, client, g := setupMockGitLabClient(t)
		defer teardown(server)

		mux.HandleFunc("/api/v4/projects/1/deployments", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `[
			{
				"id": 1, 
				"status": "success",
				"environment": {
					"name": "production"
				}, 
				"deployable": { 
					"finished_at": "2020-10-06T15:30:53.355Z",
					"duration": 123.45,
					"pipeline": {
						"id": 1
					}
				}
			}
			]`)
		})

		p := &collector.Project{
			ID:                1,
			Name:              "test",
			Path:              "test",
			PathWithNamespace: "test/test",
			Namespace:         "test/test",
			WebURL:            "http://test.com/test/test",
		}

		timestamp, e := time.Parse(time.RFC3339, "2020-10-06T15:30:53.355Z")
		if e != nil {
			t.Errorf(e.Error())
		}

		want := []*collector.Deployment{{
			ID:               1,
			Status:           "success",
			EnvironmentName:  "production",
			ProjectID:        1,
			ProjectName:      "test",
			ProjectPath:      "test",
			ProjectNamespace: "test/test",
			PipelineID:       1,
			FinishedAt:       &timestamp,
			Duration:         123.45,
		}}

		got, err := g.GetDeployments(p, client, getDeploymentListOptions())
		if err != nil {
			t.Errorf("Error getting Deployments: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v; wanted %+v", got, want)
		}
	})

	t.Run("get multiple pages of deployments from GitLab", func(t *testing.T) {

		mux, server, client, g := setupMockGitLabClient(t)
		defer teardown(server)

		mux.HandleFunc("/api/v4/projects/1/deployments", func(w http.ResponseWriter, r *http.Request) {

			if r.URL.Query()["page"][0] == "1" {
				w.Header().Set("X-Page", "1")
				w.Header().Set("X-Total-Pages", "2")
				w.Header().Set("X-Next-Page", "2")
				fmt.Fprint(w, `[{
						"id": 1,
						"status": "success",
						"environment": {
							"name": "production"
							},
							"deployable": {
								"finished_at": "2020-10-06T15:30:53.355Z",
								"duration": 123.45,
								"pipeline": {
									"id": 1
								}
							}
						}
						]`)
			}

			if r.URL.Query()["page"][0] == "2" {
				w.Header().Set("X-Page", "2")
				w.Header().Set("X-Total-Pages", "2")
				w.Header().Set("X-Next-Page", "2")
				fmt.Fprint(w, `[{
						"id": 2,
						"status": "success",
						"environment": {
							"name": "production"
							},
							"deployable": {
								"finished_at": "2020-10-06T15:30:53.355Z",
								"duration": 123.45,
								"pipeline": {
									"id": 2
								}
							}
						}
						]`)
			}
		})

		p := &collector.Project{
			ID:                1,
			Name:              "test",
			Path:              "test",
			PathWithNamespace: "test/test",
			Namespace:         "test/test",
			WebURL:            "http://test.com/test/test",
		}

		timestamp, e := time.Parse(time.RFC3339, "2020-10-06T15:30:53.355Z")
		if e != nil {
			t.Errorf(e.Error())
		}

		want := []*collector.Deployment{
			{
				ID:               1,
				Status:           "success",
				EnvironmentName:  "production",
				ProjectID:        1,
				ProjectName:      "test",
				ProjectPath:      "test",
				ProjectNamespace: "test/test",
				PipelineID:       1,
				FinishedAt:       &timestamp,
				Duration:         123.45,
			},
			{
				ID:               2,
				Status:           "success",
				EnvironmentName:  "production",
				ProjectID:        1,
				ProjectName:      "test",
				ProjectPath:      "test",
				ProjectNamespace: "test/test",
				PipelineID:       2,
				FinishedAt:       &timestamp,
				Duration:         123.45,
			},
		}

		opt := getDeploymentListOptions()

		got, err := g.GetDeployments(p, client, opt)
		if err != nil {
			t.Errorf("Error getting Deployments: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v; wanted %+v", got, want)
		}
	})

	t.Run("filter out non production deployments from GitLab", func(t *testing.T) {
		mux, server, client, g := setupMockGitLabClient(t)
		defer teardown(server)

		mux.HandleFunc("/api/v4/projects/1/deployments", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `[
					{
						"id": 1,
						"status": "success",
						"environment": {
							"name": "production"
						},
						"deployable": {
							"finished_at": "2020-10-06T15:30:53.355Z",
							"duration": 123.45,
							"pipeline": {
								"id": 1
							}
						}
					},
					{
						"id": 2,
						"status": "success",
						"environment": {
							"name": "staging"
						},
						"deployable": {
							"finished_at": "2020-10-06T15:30:53.355Z",
							"duration": 123.45,
							"pipeline": {
								"id": 2
							}
						}
					},
					{
						"id": 3,
						"status": "success",
						"environment": {
							"name": "production"
						},
						"deployable": {
							"finished_at": "2020-10-06T15:30:53.355Z",
							"duration": 123.45,
							"pipeline": {
								"id": 3
							}
						}
					}
					]`)
		})

		p := &collector.Project{
			ID:                1,
			Name:              "test",
			Path:              "test",
			PathWithNamespace: "test/test",
			Namespace:         "test/test",
			WebURL:            "http://test.com/test/test",
		}

		timestamp, e := time.Parse(time.RFC3339, "2020-10-06T15:30:53.355Z")
		if e != nil {
			t.Errorf(e.Error())
		}

		want := []*collector.Deployment{
			{
				ID:               1,
				Status:           "success",
				EnvironmentName:  "production",
				ProjectID:        1,
				ProjectName:      "test",
				ProjectPath:      "test",
				ProjectNamespace: "test/test",
				PipelineID:       1,
				FinishedAt:       &timestamp,
				Duration:         123.45,
			},
			{
				ID:               3,
				Status:           "success",
				EnvironmentName:  "production",
				ProjectID:        1,
				ProjectName:      "test",
				ProjectPath:      "test",
				ProjectNamespace: "test/test",
				PipelineID:       3,
				FinishedAt:       &timestamp,
				Duration:         123.45,
			}}

		got, err := g.GetDeployments(p, client, getDeploymentListOptions())
		if err != nil {
			t.Errorf("Error getting Deployments: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v; wanted %+v", got, want)
		}
	})

	t.Run("filter out deployments not successful or failed from GitLab", func(t *testing.T) {
		mux, server, client, g := setupMockGitLabClient(t)
		defer teardown(server)

		mux.HandleFunc("/api/v4/projects/1/deployments", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `[
					{
						"id": 1,
						"status": "success",
						"environment": {
							"name": "production"
						},
						"deployable": {
							"finished_at": "2020-10-06T15:30:53.355Z",
							"duration": 123.45,
							"pipeline": {
								"id": 1
							}
						}
					},
					{
						"id": 2,
						"status": "pending",
						"environment": {
							"name": "production"
						},
						"deployable": {
							"finished_at": "2020-10-06T15:30:53.355Z",
							"duration": 123.45,
							"pipeline": {
								"id": 2
							}
						}
					},
					{
						"id": 3,
						"status": "failed",
						"environment": {
							"name": "production"
						},
						"deployable": {
							"finished_at": "2020-10-06T15:30:53.355Z",
							"duration": 123.45,
							"pipeline": {
								"id": 3
							}
						}
					}
					]`)
		})

		p := &collector.Project{
			ID:                1,
			Name:              "test",
			Path:              "test",
			PathWithNamespace: "test/test",
			Namespace:         "test/test",
			WebURL:            "http://test.com/test/test",
		}

		timestamp, e := time.Parse(time.RFC3339, "2020-10-06T15:30:53.355Z")
		if e != nil {
			t.Errorf(e.Error())
		}

		want := []*collector.Deployment{
			{
				ID:               1,
				Status:           "success",
				EnvironmentName:  "production",
				ProjectID:        1,
				ProjectName:      "test",
				ProjectPath:      "test",
				ProjectNamespace: "test/test",
				PipelineID:       1,
				FinishedAt:       &timestamp,
				Duration:         123.45,
			},
			{
				ID:               3,
				Status:           "failed",
				EnvironmentName:  "production",
				ProjectID:        1,
				ProjectName:      "test",
				ProjectPath:      "test",
				ProjectNamespace: "test/test",
				PipelineID:       3,
				FinishedAt:       &timestamp,
				Duration:         123.45,
			}}

		got, err := g.GetDeployments(p, client, getDeploymentListOptions())
		if err != nil {
			t.Errorf("Error getting Deployments: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v; wanted %+v", got, want)
		}
	})

	t.Run("update deployment in repository", func(t *testing.T) {

		mux, server, client, g := setupMockGitLabClient(t)
		defer teardown(server)

		mux.HandleFunc("/api/v4/projects/1/deployments", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `[{
					"id": 1,
					"status": "success",
					"environment": {
						"name": "production"
						},
						"deployable": {
							"finished_at": "2020-10-06T15:30:53.355Z",
							"duration": 123.45,
							"pipeline": {
								"id": 1
							}
						}
					}]`)
		})

		mockRepository := new(mockRepo)

		p := []*collector.Project{
			{
				ID:                1,
				Name:              "test",
				Path:              "test",
				PathWithNamespace: "test/test",
				Namespace:         "test/test",
				WebURL:            "http://test.com/test/test",
			}}

		timestamp, e := time.Parse(time.RFC3339, "2020-10-06T15:30:53.355Z")
		if e != nil {
			t.Errorf(e.Error())
		}

		want := []*collector.Deployment{{
			ID:               1,
			Status:           "success",
			EnvironmentName:  "production",
			ProjectID:        1,
			ProjectName:      "test",
			ProjectPath:      "test",
			ProjectNamespace: "test/test",
			PipelineID:       1,
			FinishedAt:       &timestamp,
			Duration:         123.45,
		}}

		g.UpdateDeployments(p, client, mockRepository)

		if !reflect.DeepEqual(mockRepository.DeploymentData, want) {
			t.Errorf("got %+v; wanted %+v", mockRepository.DeploymentData, want)
		}
	})
}

func TestRefreshData(t *testing.T) {
	t.Run("refresh data successfully", func(t *testing.T) {

		_, _, _, g := setupMockGitLabClient(t)

		r := new(mockRepo)

		err := g.RefreshData(r)

		if err != nil {
			t.Errorf("Unexpected error: %v", err.Error())
		}

	})
}

func getProjectListOptions() *gl.ListProjectsOptions {
	return &gl.ListProjectsOptions{
		ListOptions: gl.ListOptions{Page: 1, PerPage: 1},
		Simple:      gl.Bool(false),
	}
}

func getDeploymentListOptions() *gl.ListProjectDeploymentsOptions {
	return &gl.ListProjectDeploymentsOptions{
		ListOptions: gl.ListOptions{Page: 1, PerPage: 1},
		Environment: gl.String("production"),
		Status:      gl.String("success"),
	}
}

func setupMockGitLabClient(t *testing.T) (*http.ServeMux, *httptest.Server, *gl.Client, *gitlab.GitLab) {

	mux := http.NewServeMux()

	server := httptest.NewServer(mux)

	g := &gitlab.GitLab{
		Token: "",
		URL:   server.URL,
	}

	client, err := g.SetupClient(g.Token, g.URL)
	if err != nil {
		server.Close()
		t.Fatalf("Error creating mock GitLab client: %v", err)
	}

	return mux, server, client, g
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

func teardown(server *httptest.Server) {
	server.Close()
}
