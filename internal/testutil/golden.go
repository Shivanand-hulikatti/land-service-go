package testutil

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// GoldenPath returns the path to a file under docs/golden/.
func GoldenPath(name string) string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("runtime.Caller failed")
	}
	// testutil -> internal -> land-services-go
	root := filepath.Join(filepath.Dir(file), "..", "..")
	return filepath.Join(root, "docs", "golden", name)
}

// ReadGolden loads a golden JSON fixture.
func ReadGolden(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(GoldenPath(name))
	if err != nil {
		t.Fatalf("read golden %s: %v", name, err)
	}
	return data
}
