package sparkrun

import (
	"context"

	"github.com/Riven-Spell/generic/list_tools"
	"github.com/Riven-Spell/sparky/common/cmdline"
)

// ClusterImportOptions configures [Client.ClusterImport].
//
// sparkrun's `cluster import` imports an external cluster config
// (e.g. a spark-vllm-docker .env) into a named sparkrun cluster.
type ClusterImportOptions struct {
	// Path is the external config file to import (a spark-vllm-docker
	// .env). It is the positional argument to `cluster import svd PATH`.
	Path string
	// Name is the cluster name to import into (--name). When unset it
	// defaults to a name derived from the source.
	Name *string `cmd:"name"`
	// SetDefault, when true, sets the imported cluster as the default
	// (--default).
	SetDefault bool `cmd:"default"`
	// DryRun shows what would be done without importing (--dry-run).
	DryRun bool `cmd:"dry-run"`
}

// ClusterImport imports an external cluster config into a sparkrun
// cluster.
//
// `sparkrun cluster import svd <path> [--name] [--default] [--dry-run]`
func (c *cliClient) ClusterImport(ctx context.Context, opts ...ClusterImportOptions) error {
	var o = list_tools.FirstOrZero(opts)
	args, err := cmdline.BuildArgs([]string{"cluster", "import", "svd", o.Path}, o)
	if err != nil {
		return err
	}
	return c.plainCmd(ctx, args...)
}
