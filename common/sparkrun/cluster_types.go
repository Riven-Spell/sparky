package sparkrun

// ClusterSummary is the saved definition of a cluster, as returned
// by `sparkrun cluster show`, and reused as the element type of
// `sparkrun cluster list` and the body of `sparkrun cluster
// default` (where the latter may also be the JSON literal null,
// surfaced as a nil *ClusterSummary).
type ClusterSummary struct {
	Name        string   `json:"name"`
	Hosts       []string `json:"hosts"`
	Description string   `json:"description"`
	User        string   `json:"user"`
	Default     bool     `json:"default"`
}

// ClusterStatusResult is the body of `sparkrun cluster status
// --json`. It is keyed by cluster ID for multi-host groups and
// exposes solo entries as a flat list. Per-host SSH/connect errors
// land in Errors rather than aborting the call -- the wrapper
// therefore does not treat them as fatal.
//
// Meta carries a recipe_state blob that is large and frequently
// version-dependent; we type it as any so callers can decode the
// fields they care about without this package baking in fields it
// doesn't use. See 0-sparky-docs/Sparkrun output shapes/Cluster
// status.md for the full shape.
type ClusterStatusResult struct {
	Groups          map[string]ClusterGroup `json:"groups"`
	SoloEntries     []ClusterSoloEntry      `json:"solo_entries"`
	IdleHosts       []string                `json:"idle_hosts"`
	PendingOps      []any                   `json:"pending_ops"`
	Errors          map[string]string       `json:"errors"`
	TotalContainers int                     `json:"total_containers"`
	HostCount       int                     `json:"host_count"`
}

// ClusterGroup is a multi-host workload. The agent reads .Meta and
// .Containers to decide which models it can serve on localhost.
type ClusterGroup struct {
	Meta       any          `json:"meta"`
	Containers []ClusterBox `json:"containers"`
	Label      string       `json:"label"`
}

// ClusterSoloEntry is a single-host workload.
type ClusterSoloEntry struct {
	ClusterID string `json:"cluster_id"`
	Meta      any    `json:"meta"`
	Host      string `json:"host"`
	Status    string `json:"status"`
	Image     string `json:"image"`
	Label     string `json:"label"`
}

// ClusterBox describes one container in a [ClusterGroup].
type ClusterBox struct {
	Host   string `json:"host"`
	Role   string `json:"role"`
	Status string `json:"status"`
	Image  string `json:"image"`
}

// ClusterCheckJobResult is the body of `sparkrun cluster check-job
// --json`. On a non-zero exit the wrapper still attempts to decode
// the body; if it succeeds, the parsed struct is returned on the
// [ExitError]; if not, the raw bytes are exposed via [ExitError.Raw].
type ClusterCheckJobResult struct {
	Running           bool            `json:"running"`
	ClusterID         string          `json:"cluster_id"`
	Healthy           *bool           `json:"healthy"`
	Metadata          map[string]any  `json:"metadata"`
	ContainerStatuses map[string]bool `json:"container_statuses"`
	Hosts             []string        `json:"hosts"`
}
