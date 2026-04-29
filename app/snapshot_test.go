package app

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPreviewTexSnapshots(t *testing.T) {
	root := repoRoot(t)
	work, err := os.MkdirTemp("", "planner-snapshot-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(work) })
	if err := os.Mkdir(filepath.Join(work, "out"), 0755); err != nil {
		t.Fatal(err)
	}

	overridePath := filepath.Join(work, "snapshot.yaml")
	if err := os.WriteFile(overridePath, []byte(`publicholidays:
  countrycodes: []
  countrycode: ""
  custom: []
  icsfiles: []
events:
  files: []
  custom:
    - date: "2026-01-08"
      name: "Snapshot event"
      shortName: "Snap"
      types: [Event]
`), 0600); err != nil {
		t.Fatal(err)
	}

	configFiles := []string{
		filepath.Join(root, "cfg", "base.yaml"),
		filepath.Join(root, "cfg", "rm2.base.yaml"),
		filepath.Join(root, "cfg", "rm2.desk.weekly.yaml"),
		filepath.Join(root, "cfg", "rm2.desk.monthly.yaml"),
		filepath.Join(root, "cfg", "template_desk_monthly_rm2.yaml"),
		overridePath,
	}

	previous, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(work); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(previous); err != nil {
			t.Fatal(err)
		}
	})

	t.Setenv("PLANNER_YEAR", "2026")
	args := []string{"plannergen", "--preview", "--config", strings.Join(configFiles, ",")}
	if err := New().RunContext(context.Background(), args); err != nil {
		t.Fatal(err)
	}

	for _, name := range []string{"snapshot.tex", "title.tex", "annual.tex", "monthly.tex", "weekly.tex"} {
		gotPath := filepath.Join(work, "out", name)
		wantPath := filepath.Join(root, "app", "testdata", "snapshots", name)
		assertSnapshot(t, wantPath, gotPath)
	}
}

func assertSnapshot(t *testing.T, wantPath string, gotPath string) {
	t.Helper()
	got, err := os.ReadFile(gotPath)
	if err != nil {
		t.Fatal(err)
	}
	if os.Getenv("UPDATE_SNAPSHOTS") == "1" {
		if err := os.MkdirAll(filepath.Dir(wantPath), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(wantPath, got, 0600); err != nil {
			t.Fatal(err)
		}
		return
	}

	want, err := os.ReadFile(wantPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Fatalf("snapshot mismatch for %s; run UPDATE_SNAPSHOTS=1 go test ./app to update", filepath.Base(wantPath))
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Clean(filepath.Join(wd, ".."))
}
