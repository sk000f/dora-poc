package gitlab_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/sk000f/metrix/pkg/collector"
	"github.com/sk000f/metrix/pkg/collector/gitlab"
	gl "github.com/xanzy/go-gitlab"
)

func TestGitLabProjects(t *testing.T) {
	t.Run("get single Project from GitLab", func(t *testing.T) {

		mux, server, client := setupMockGitLabClient(t)
		defer teardown(server)

		mux.HandleFunc("/api/v4/projects", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `[{
					"id": 1, 
					"name": "test", 
					"name_with_namespace": "test/test", 
					"web_url": "http://test.com/test/test"
					}]`)
		})

		want := []*collector.Project{
			{
				ID:                1,
				Name:              "test",
				NameWithNamespace: "test/test",
				WebURL:            "http://test.com/test/test",
			}}

		got, err := gitlab.GetProjects(client, getProjectListOptions())
		if err != nil {
			t.Errorf("Error getting Projects: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v; wanted %+v", got, want)
		}
	})

	t.Run("get multiple pages of Projects from GitLab", func(t *testing.T) {
		mux, server, client := setupMockGitLabClient(t)
		defer teardown(server)

		mux.HandleFunc("/api/v4/projects", func(w http.ResponseWriter, r *http.Request) {

			if r.URL.Query()["page"][0] == "1" {
				w.Header().Set("X-Page", "1")
				w.Header().Set("X-Total-Pages", "2")
				w.Header().Set("X-Next-Page", "2")
				fmt.Fprint(w, `[{
					"id": 1, 
					"name": "test", 
					"name_with_namespace": "test/test", 
					"web_url": "http://test.com/test/test"}
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
					"name_with_namespace": "test/test", 
					"web_url": "http://test.com/test/test"
				},
				{
					"id": 3, 
					"name": "test", 
					"name_with_namespace": "test/test", 
					"web_url": "http://test.com/test/test"
				}
				]`)
			}
		})

		want := []*collector.Project{
			{
				ID:                1,
				Name:              "test",
				NameWithNamespace: "test/test",
				WebURL:            "http://test.com/test/test",
			},
			{
				ID:                2,
				Name:              "test",
				NameWithNamespace: "test/test",
				WebURL:            "http://test.com/test/test",
			},
			{
				ID:                3,
				Name:              "test",
				NameWithNamespace: "test/test",
				WebURL:            "http://test.com/test/test",
			},
		}

		opt := getProjectListOptions()

		got, err := gitlab.GetProjects(client, opt)
		if err != nil {
			t.Errorf("Error getting Projects: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v; wanted %+v", got, want)
		}
	})
}

func TestGitLabDeploymemts(t *testing.T) {
	t.Run("get all deployments from GitLab", func(t *testing.T) {
		mux, server, client := setupMockGitLabClient(t)
		defer teardown(server)

		mux.HandleFunc("/api/v4/projects/1/deployments", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `[{
				"id": 1, 
				"environment": {
					"name": "production"
					}, 
					"deployable": {
						"status": "success", 
						"pipeline": {
							"id": 1
						}
					}
				}]`)
		})

		want := []*collector.Deployment{{
			ID:              1,
			Status:          "success",
			EnvironmentName: "production",
			PipelineID:      1,
		}}

		got, err := gitlab.GetDeployments(1, client, getDeploymentListOptions())
		if err != nil {
			t.Errorf("Error getting Deployments: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v; wanted %+v", got, want)
		}
	})
}

func getProjectListOptions() *gl.ListProjectsOptions {
	return &gl.ListProjectsOptions{
		ListOptions: gl.ListOptions{Page: 1, PerPage: 1},
		Archived:    gl.Bool(true),
		OrderBy:     gl.String("id"),
		Sort:        gl.String("asc"),
		Search:      gl.String("query"),
		Simple:      gl.Bool(true),
		Visibility:  gl.Visibility(gl.PublicVisibility),
	}
}

func getDeploymentListOptions() *gl.ListProjectDeploymentsOptions {
	return &gl.ListProjectDeploymentsOptions{
		ListOptions: gl.ListOptions{Page: 1, PerPage: 1},
		OrderBy:     gl.String("id"),
		Sort:        gl.String("asc"),
		Environment: gl.String("production"),
		Status:      gl.String("success"),
	}
}

func setupMockGitLabClient(t *testing.T) (*http.ServeMux, *httptest.Server, *gl.Client) {

	mux := http.NewServeMux()

	server := httptest.NewServer(mux)

	client, err := gitlab.SetupClient("", server.URL)
	if err != nil {
		server.Close()
		t.Fatalf("Error creating mock GitLab client: %v", err)
	}

	return mux, server, client
}

func teardown(server *httptest.Server) {
	server.Close()
}
