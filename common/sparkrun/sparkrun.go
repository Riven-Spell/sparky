// Package sparkrun wraps the `sparkrun` CLI used to manage LLM workloads
// across a DGX Spark cluster. The package exposes a typed [Client] that
// the agent and manager binaries depend on, so they can be unit-tested
// against a fake without invoking the real binary.
//
// One file per sparkrun verb; methods follow the options-bag style
// described in 0-sparky-docs/Code Style.md. Streaming commands
// (e.g. ClusterMonitor) return stdout and stderr io.Reader streams
// of the raw subprocess output -- the caller decodes per line.
//
// This pass implements only the `cluster` verb; other verbs (Run, Stop,
// Logs, recipe*, registry*, etc.) are added as consumers need them.
package sparkrun

// DefaultBinary is the executable name looked up via $PATH when no
// explicit binary path is configured via [CliClientOptions.WithBinaryPath].
const DefaultBinary = "sparkrun"
