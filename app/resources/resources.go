package resources

import (
	"os"
	"path/filepath"
	"strings"
)

const EnvDir = "PLANNER_RESOURCE_DIR"

func Resolve(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	if _, err := os.Stat(path); err == nil {
		return path
	}
	if dir := os.Getenv(EnvDir); dir != "" {
		return filepath.Join(dir, path)
	}

	return path
}

func Bundled(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	if dir := os.Getenv(EnvDir); dir != "" {
		return filepath.Join(dir, path)
	}

	return Resolve(path)
}

func BundledList(paths string) string {
	if paths == "" {
		return ""
	}

	parts := strings.Split(paths, ",")
	for i, path := range parts {
		parts[i] = Bundled(path)
	}

	return strings.Join(parts, ",")
}
