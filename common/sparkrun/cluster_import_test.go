package sparkrun

import "testing"

func TestClusterImportOptionsArgs(t *testing.T) {
	name := "lab"
	got := mustArgs(t, []string{"cluster", "import", "svd", "/tmp/foo.env"}, ClusterImportOptions{
		Name:       &name,
		SetDefault: true,
		DryRun:     true,
	})
	want := []string{
		"cluster", "import", "svd", "/tmp/foo.env",
		"--name", "lab",
		"--default",
		"--dry-run",
	}
	assertEqualArgs(t, got, want)
}

func TestClusterImportOptionsPathOnly(t *testing.T) {
	got := mustArgs(t, []string{"cluster", "import", "svd", "/f"}, ClusterImportOptions{})
	assertEqualArgs(t, got, []string{"cluster", "import", "svd", "/f"})
}
