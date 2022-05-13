package nodejs

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRuntimeFeatures(t *testing.T) {
	t.Parallel()
	type testCase struct {
		name              string
		sourceMapsEnabled bool
	}

	wd, err := os.Getwd()

	if err != nil {
		t.Error(err)
		return
	}

	for _, v := range [...]testCase{
		{"callsites", false},
		{"imports", false},
		{"sourcemaps_external", false},
		{"sourcemaps_external", true},
		{"sourcemaps_inline", false},
		{"sourcemaps_inline", true},
	} {
		func(tc testCase) {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				cmd := exec.Command("./prybar-nodejs", filepath.Join(wd, "testdata/test_"+tc.name))

				cmd.Dir = filepath.Join(wd, "../..")

				if tc.sourceMapsEnabled {
					cmd.Env = append(os.Environ(), "ENABLE_SOURCE_MAPS=1")
				}

				var out strings.Builder
				cmd.Stderr = &out

				err := cmd.Run()
				if err != nil {
					t.Error(err, "\n"+out.String())
				}
			})
		}(v)
	}
}
