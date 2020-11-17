package gitlab_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/sk000f/metrix/pkg/collector/gitlab"
	gl "github.com/xanzy/go-gitlab"
)

func TestGitLabProjects(t *testing.T) {
	t.Run("get Projects from GitLab", func(t *testing.T) {

		mux, server, client := setupMockGitLabClient(t)
		defer teardown(server)

		mux.HandleFunc("/api/v4/projects", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `[{"id":1}]`)
		})

		want := []*gitlab.Project{{ID: 1}}

		got, err := gitlab.GetProjects(client, getProjectListOptions())

		if err != nil {
			t.Errorf("Error getting Projects: %v", err)
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
		OrderBy:     gl.String("name"),
		Sort:        gl.String("asc"),
		Search:      gl.String("query"),
		Simple:      gl.Bool(true),
		Visibility:  gl.Visibility(gl.PublicVisibility),
	}
}

func setupMockGitLabClient(t *testing.T) (*http.ServeMux, *httptest.Server, *gl.Client) {

	mux := http.NewServeMux()

	server := httptest.NewServer(mux)

	client, err := gitlab.SetupClient(server.URL)
	if err != nil {
		server.Close()
		t.Fatalf("Error creating mock GitLab client: %v", err)
	}

	return mux, server, client
}

func teardown(server *httptest.Server) {
	server.Close()
}
