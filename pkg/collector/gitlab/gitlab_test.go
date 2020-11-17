package gitlab_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/sk000f/metrix/pkg/collector/gitlab"
	gogitlab "github.com/xanzy/go-gitlab"
)

func TestGitLabProjects(t *testing.T) {
	t.Run("get Projects from GitLab", func(t *testing.T) {

		mux, server, client := setupMockGitLabClient(t)
		defer teardown(server)

		mux.HandleFunc("/api/v4/projects", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `[{"id":1}]`)
		})

		opt := &gogitlab.ListProjectsOptions{
			ListOptions: gogitlab.ListOptions{2, 3},
			Archived:    gogitlab.Bool(true),
			OrderBy:     gogitlab.String("name"),
			Sort:        gogitlab.String("asc"),
			Search:      gogitlab.String("query"),
			Simple:      gogitlab.Bool(true),
			Visibility:  gogitlab.Visibility(gogitlab.PublicVisibility),
		}

		want := []*gogitlab.Project{{ID: 1}}
		got, _, err := client.Projects.ListProjects(opt)
		if err != nil {
			t.Errorf("Error getting Projects: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v; wanted %+v", got, want)
		}
	})
}

func setupMockGitLabClient(t *testing.T) (*http.ServeMux, *httptest.Server, *gogitlab.Client) {

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
