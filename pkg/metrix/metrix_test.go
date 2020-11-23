package metrix_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/sk000f/metrix/pkg/metrix"
)

func TestMetrix(t *testing.T) {
	t.Run("configuration values set correctly", func(t *testing.T) {
		os.Setenv("METRIX_GITLAB_URL", "https://example.com")
		os.Setenv("METRIX_GITLAB_TOKEN", "1234567890")

		want := &metrix.Config{GitLabURL: "https://example.com", GitLabToken: "1234567890"}
		got := metrix.SetupConfig()

		if !reflect.DeepEqual(got, want) {
			t.Errorf("want %v; got %v", want, got)
		}

		os.Unsetenv("METRIX_GITLAB_URL")
		os.Unsetenv("METRIX_GITLAB_TOKEN")
	})
}
