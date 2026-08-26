package sparkrun

import (
	"context"
	"io"

	"github.com/Riven-Spell/sparky/common/sparkrun/sparkrun_models"
)

// Client is the surface that the agent and manager binaries consume.
// Methods take a context.Context; required positional arguments are
// direct parameters; everything else is a varargs options bag.
//
// This pass implements the `cluster` and `recipe` verbs. Run/Stop/
// Logs and registry/proxy/etc. methods are added to this interface as
// they are implemented -- see common/sparkrun/plan.md for the order.
type Client interface {
	// ClusterShow returns the saved definition of a named cluster.
	ClusterShow(ctx context.Context, name string, opts ...ClusterShowOptions) (*sparkrun_models.ClusterSummary, error)

	// ClusterList returns every saved cluster as a slice.
	// The slice may be empty; the wrapper does not return an error
	// in that case.
	ClusterList(ctx context.Context, opts ...ClusterListOptions) ([]sparkrun_models.ClusterSummary, error)

	// ClusterDefault returns the currently-default cluster, or
	// (nil, nil) if no default is set.
	ClusterDefault(ctx context.Context, opts ...ClusterDefaultOptions) (*sparkrun_models.ClusterSummary, error)

	// ClusterStatus lists sparkrun containers running on cluster
	// hosts. opts.Cluster / opts.Hosts select the host set; opts.DryRun
	// is mutually exclusive with --json and
	// causes the call to return whatever dry-run text sparkrun
	// emits (the wrapper does not attempt to parse it).
	ClusterStatus(ctx context.Context, opts ClusterStatusOptions) (*sparkrun_models.ClusterStatusResult, error)

	// ClusterCreate creates a new named cluster. Hosts (or a hosts
	// file) is mandatory -- sparkrun itself rejects calls without
	// one.
	ClusterCreate(ctx context.Context, name string, opts ...ClusterCreateOptions) error

	// ClusterDelete removes a saved cluster. The wrapper always
	// passes --force; non-interactive callers cannot answer the
	// [y/N] prompt otherwise.
	ClusterDelete(ctx context.Context, name string, opts ...ClusterDeleteOptions) error

	// ClusterUpdate mutates an existing cluster. opts.Hosts is
	// optional; if set it must be one of [Hosts], [HostsFromFile],
	// [AddHosts], or [RemoveHosts] (the latter two are repeatable;
	// sparkrun rejects mixing them with --hosts / --hosts-file).
	ClusterUpdate(ctx context.Context, name string, opts ClusterUpdateOptions) error

	// ClusterSetDefault marks a saved cluster as the default.
	ClusterSetDefault(ctx context.Context, name string, opts ...ClusterSetDefaultOptions) error

	// ClusterUnsetDefault clears the default cluster setting.
	ClusterUnsetDefault(ctx context.Context, opts ...ClusterUnsetDefaultOptions) error

	// ClusterImport imports an external cluster config (e.g. a
	// spark-vllm-docker .env) into a sparkrun cluster. Currently only
	// the svd provider is used; opts.Path is the config file.
	ClusterImport(ctx context.Context, opts ...ClusterImportOptions) error

	// ClusterCheckJob reports whether a recipe (by name) or a
	// specific running workload (by cluster ID) is up. The returned
	// *ClusterCheckJobResult is non-nil on exit 0; on a non-zero
	// exit the result is wrapped in *ExitError.Result() so callers
	// that care about the "not running" body can still read it. A
	// nil result with no error indicates "no such recipe or cluster
	// ID".
	ClusterCheckJob(ctx context.Context, target RecipeNameOrJobID, opts ...ClusterCheckJobOptions) (*sparkrun_models.ClusterCheckJobResult, error)

	// ClusterMonitor streams per-host CPU/RAM/GPU metrics. The
	// returned stdout [io.Reader] produces NDJSON -- one JSON object
	// per line -- and stderr produces diagnostic output until the
	// subprocess exits. The underlying process is terminated by
	// cancelling the context or by calling the returned kill function;
	// callers read until io.EOF rather than closing the streams.
	// An [ExitError] surfaces non-zero exits.
	ClusterMonitor(ctx context.Context, opts ClusterMonitorOptions) (stdout, stderr io.Reader, kill func(), err error)

	// Logs streams the log output for a running recipe (by name) or a
	// specific running workload (by cluster ID). The returned stdout
	// [io.Reader] produces the log lines as sparkrun emits them; stderr
	// carries sparkrun diagnostics. The subprocess is terminated by
	// cancelling the context or by calling the returned kill function;
	// callers read until io.EOF. An [ExitError] surfaces non-zero
	// exits.
	Logs(ctx context.Context, target RecipeNameOrJobID, opts ...LogsOptions) (stdout, stderr io.Reader, kill func(), err error)

	// RecipeList returns every available recipe from all registries.
	// query filters by name/model/description; an empty string lists
	// all. Returns a (possibly empty) slice.
	RecipeList(ctx context.Context, query string, opts ...RecipeListOptions) ([]sparkrun_models.RecipeSummary, error)

	// RecipeSearch finds recipes by name, model, or description.
	// Returns a (possibly empty) slice.
	RecipeSearch(ctx context.Context, query string, opts ...RecipeSearchOptions) ([]sparkrun_models.RecipeSummary, error)

	// RecipeShow returns the fully-normalized recipe for a recipe name
	// or a recipe-YAML file.
	RecipeShow(ctx context.Context, target RecipeNameOrFile, opts ...RecipeShowOptions) (*sparkrun_models.RecipeDetail, error)

	// RecipeVram estimates VRAM usage for a recipe on a DGX Spark.
	RecipeVram(ctx context.Context, target RecipeNameOrFile, opts ...RecipeVramOptions) (*sparkrun_models.RecipeVramEstimate, error)

	// RecipeValidate checks that a recipe (by name or file) is valid.
	RecipeValidate(ctx context.Context, target RecipeNameOrFile, opts ...RecipeValidateOptions) (*sparkrun_models.RecipeValidateResult, error)
}
