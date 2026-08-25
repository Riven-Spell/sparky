# Sparkrun wrapper — plan

## Goals
- Replace shell-outs to `sparkrun` with a typed Go API in `common/sparkrun`.
- Provide a `Client` interface so `agent` and `manager` can be unit-tested without invoking the real binary.
- Obey `0-sparky-docs/Code Style.md` (one file per purpose, options bags with varargs, fully labeled struct fields, unexported internals, explicit error discarding).

## Non-goals (for now)
- A `MockClient` — added only when a concrete test in `agent` or `manager` needs one.
- A "remote" implementation that proxies through the manager HTTP API — out of scope until asked.
- Streaming logs / handle-coupling from `Run` — caller resolves `cluster_id` themselves.
- Handling non-`--json` output for commands that have `--json` — always pass `--json` when available.

## Package

Path: `common/sparkrun`. Import: `github.com/Riven-Spell/sparky/common/sparkrun`. Lives alongside `common/models`, `common/enum`, etc.

## File layout

```
common/sparkrun/
  plan.md              this document
  sparkrun.go          package doc, DefaultBinary constant
  errors.go            SparkrunError + ErrorKind + AsSparkrunError
  target.go            RecipeNameOrFile, RecipeNameOrJobID, constructors
  host_diff.go         HostsList + HostsDiff sealed interfaces;
                       Hosts/HostsFromFile/AddHosts/RemoveHosts/CompoundHostDiff
  key_value_opts.go    KeyValueOpts map[string]string -> sorted --<name> key=value (ArgProvider)
  transfer_mode.go     TransferMode enum (auto/local/push/delegated)
  transfer_interface.go TransferInterface enum (auto/cx7/mgmt)
  topology.go          Topology enum (none/direct/switch/ring)
  cli_client.go        cliClient + NewCliClient + CliOption + runCmd/jsonCmd/streamCmd/plainCmd
  interface.go         Client interface (grows as verbs land)
  cluster_types.go     ClusterSummary + ClusterStatusResult + ClusterGroup +
                       ClusterSoloEntry + ClusterBox + ClusterCheckJobResult

  # `cluster` verb — DONE
  cluster_show.go          ClusterShow
  cluster_status.go        ClusterStatus + ClusterStatusOptions
  cluster_list.go          ClusterList
  cluster_default.go       ClusterDefault
  cluster_create.go        ClusterCreate + ClusterCreateOptions  cluster_delete.go        ClusterDelete (always passes --force; no options)
  cluster_update.go        ClusterUpdate + ClusterUpdateOptions
  cluster_set_default.go   ClusterSetDefault
  cluster_unset_default.go ClusterUnsetDefault
  cluster_import.go        ClusterImport + ClusterImportOptions
  cluster_check_job.go     ClusterCheckJob + ClusterCheckJobOptions
  cluster_monitor.go       ClusterMonitor + ClusterMonitorOptions

  # lifecycle — NOT YET
  run.go               Run + RunOptions
  stop.go              Stop + StopOptions, StopAll + StopAllOptions
  logs.go              Logs + LogsOptions (io.ReadCloser)
  status.go            Status + StatusOptions + StatusResult   # note: top-level `status`, not cluster

  # recipe / registry / proxy / export / arena / benchmark / setup / tune / update — NOT YET
  recipe.go            RecipeList/Search/Show/Vram/Validate
  registry.go          RegistryList/Add/Remove/Update/...
  proxy.go             ProxyStart/Stop/Status/Load/Unload/...
  export.go            ExportRecipe/RunningRecipe/Systemd
  arena.go             ArenaLogin/Logout/Status/Benchmark
  benchmark.go         BenchmarkRun
  setup.go             Setup*
  tune.go              TuneVllm, TuneSglang (io.ReadCloser)
  update.go            Update
```

One top-level command per file. **No** shared options bags — each command declares its own.

## Implementation status

| Verb | Status | Smoke-tested? |
|---|---|---|
| `cluster show`       | ✅ done | yes |
| `cluster status`     | ✅ done | yes |
| `cluster list`       | ✅ done | yes |
| `cluster default`    | ✅ done | yes |
| `cluster create`     | ✅ done | yes (incl. `--default` flag; HostsList; executor/scheduler/max-gpu-mem-util) |
| `cluster delete`     | ✅ done | yes |
| `cluster update`     | ✅ done | yes (all HostsDiff modes incl. CompoundHostDiff; infer-hardware/executor/scheduler/max-gpu-mem-util) |
| `cluster set-default`| ✅ done | yes |
| `cluster unset-default` | ✅ done | yes |
| `cluster import`     | ✅ done | yes (svd provider; dry-run+real) |
| `cluster check-job`  | ✅ done | yes (both error paths) |
| `cluster monitor`    | ✅ done | yes (NDJSON stream) |
| everything else      | not yet | n/a |

## Design

### Client interface (`interface.go`)

Every CLI verb is a method. Methods take `context.Context`. Required positional args are direct parameters; options are a varargs options bag. `Stop` and `StopAll` are separate methods because `--all` is a distinct mode of `stop`.

The interface grows as verbs are implemented. **Current state (cluster pass):**

```go
type Client interface {
    // cluster discovery / management
    ClusterShow(ctx context.Context, name string, opts ...ClusterShowOptions) (*ClusterSummary, error)
    ClusterList(ctx context.Context, opts ...ClusterListOptions) ([]ClusterSummary, error)
    ClusterDefault(ctx context.Context, opts ...ClusterDefaultOptions) (*ClusterSummary, error)
    ClusterStatus(ctx context.Context, opts ClusterStatusOptions) (*ClusterStatusResult, error)
    ClusterCreate(ctx context.Context, name string, opts ...ClusterCreateOptions) error
    ClusterDelete(ctx context.Context, name string, opts ...ClusterDeleteOptions) error
    ClusterUpdate(ctx context.Context, name string, opts ClusterUpdateOptions) error
    ClusterSetDefault(ctx context.Context, name string, opts ...ClusterSetDefaultOptions) error
    ClusterUnsetDefault(ctx context.Context, opts ...ClusterUnsetDefaultOptions) error
    ClusterImport(ctx context.Context, opts ...ClusterImportOptions) error
    ClusterCheckJob(ctx context.Context, target RecipeNameOrJobID, opts ...ClusterCheckJobOptions) (*ClusterCheckJobResult, error)
    ClusterMonitor(ctx context.Context, opts ClusterMonitorOptions) (stdout, stderr io.Reader, kill func(), err error)
}
```

Methods that take no positional target or that produce no structured body return just `error` (e.g. `ClusterCreate`, `ClusterDelete`, `ClusterUpdate`, `ClusterSetDefault`, `ClusterUnsetDefault`). Methods that have a typed body return `(*Result, error)`.

`ClusterCheckJob` is the exception: it can return a structured body even on non-zero exit (sparkrun emits JSON regardless). The wrapper therefore returns `(nil, *SparkrunError)` with the parsed body in `Details["result"]` on exit 1 — callers type-assert to read it. See "Non-zero-exit-with-body" below.

### Target overload (`target.go`)

Two specific interfaces for call-site type safety. `RecipeName` is intentionally usable as both.

```go
type RecipeNameOrFile interface { recipeRef() string }
type RecipeNameOrJobID interface { workloadRef() string }

type recipeName string
func (r recipeName) recipeRef() string   { return string(r) }
func (r recipeName) workloadRef() string { return string(r) }
func RecipeName(name string) recipeName { return recipeName(name) }

type recipeFile string
func (r recipeFile) recipeRef() string { return string(r) }
func RecipeFile(path string) recipeFile { return recipeFile(path) }

type jobID string
func (j jobID) workloadRef() string { return string(j) }
func JobID(id string) jobID { return jobID(id) }
```

| Method | Interface | Accepts |
|---|---|---|
| `Run`, `RecipeShow`, `RecipeVram`, `RecipeValidate` | `RecipeNameOrFile` | `RecipeName`, `RecipeFile` |
| `Stop`, `Logs`, `ClusterCheckJob`, `ExportRunningRecipe`, `ExportSystemd` | `RecipeNameOrJobID` | `RecipeName`, `JobID` |

`StopAll` doesn't take a target — listed above only for completeness of the per-method rules.

### Host list / HostsDiff overload (`host_diff.go`)

Verbs that accept a host list or mutate a cluster's host list resolve that choice through two sealed interfaces:

```go
type HostsList interface { cmdline.ArgProvider; hostsList() }   // replacement only
type HostsDiff interface { cmdline.ArgProvider; hostDiff() }    // any mutation
```

Constructors:

```go
func Hosts(values []string) HostsList            // --hosts a,b,c
func HostsFromFile(path string) HostsList        // --hosts-file <path>
func AddHosts(values []string) HostsDiff         // --add-host a,b,c
func RemoveHosts(values []string) HostsDiff      // --remove-host a,b,c
func CompoundHostDiff(added, removed []string) HostsDiff // --add-host ... --remove-host ...
```

`HostsList` (a replacement) is accepted by verbs that only ever replace the whole host set (`cluster create`, `cluster check-job`). `HostsDiff` is the superset used by `cluster update`; every `HostsList` is also a `HostsDiff`. The unexported `hostsList()`/`hostDiff()` methods seal both interfaces so only types defined in this package qualify. Incremental edits are emitted as a single comma-separated repeat per flag (sparkrun accepts `--add-host a,b,c`). Each concrete type carries a `var _ HostsDiff = ...` compile-time assertion.

The wrapper passes whatever `HostsDiff`/`HostsList` the caller chose. sparkrun itself rejects combinations of `--hosts`/`--hosts-file` with `--add-host`/`--remove-host`; we deliberately don't re-validate.

### Error type (`errors.go`)

Categorization is the value; `Details` is for callers who want to introspect. Message always begins with literal `"sparkrun"`.

```go
type ErrorKind string

const (
    ErrorKindParse  ErrorKind = "parse"   // --json parse failure (exit 0, body undecodable)
    ErrorKindExec   ErrorKind = "exec"    // could not start subprocess
    ErrorKindExit   ErrorKind = "exit"    // non-zero exit code (incl. usage errors)
    ErrorKindTarget ErrorKind = "target"  // wrong target kind for command
    ErrorKindUsage  ErrorKind = "usage"   // invalid arguments (click exit-2)
)
```

Note: **exit-2 from click is classified as `ErrorKindUsage`**, all other non-zero exits as `ErrorKindExit`. A non-zero exit with a non-JSON body is **not** a parse error — it's an `ErrorKindExit` with the body in `Details["stdout"]`/`["stderr"]`. `ErrorKindParse` is reserved for the rare "exit 0, body undecodable" case.

```go
type SparkrunError struct {
    Kind       ErrorKind
    Subcommand string
    Details    map[string]any
    Err        error
}

func (e *SparkrunError) Error() string { ... }
func (e *SparkrunError) Unwrap() error { return e.Err }
func AsSparkrunError(err error) (*SparkrunError, bool) { ... }
```

`Subcommand` is the full nested path: `"cluster"`, `"cluster status"`, `"export running-recipe"`. The helper `subcommandPath(args)` in `cli_client.go` derives it from the leading non-flag tokens.

Standard `Details` keys set by `cliClient`:
- `ErrorKindExit` → `{"exit_code": int, "stderr": string, "stdout": string}` — stderr is always captured.
- `ErrorKindParse` → `{"raw": []byte, "stderr": string}`.
- `ErrorKindExec` → `{"binary": string, "args": []string, "stderr": string, "stdout": string}`.
- `ErrorKindUsage` → same as `ErrorKindExit` (just a different `Kind`).

Callers can type-assert `*SparkrunError` (or use `AsSparkrunError`) and read these; they're never in the formatted message.

### Non-zero-exit-with-body

Some verbs — `cluster check-job` — emit a JSON body even on a non-zero exit. `ClusterCheckJob` therefore:

- returns `(*ClusterCheckJobResult, nil)` on exit 0;
- returns `(nil, *SparkrunError{..., Details: {"result": *ClusterCheckJobResult}})` on exit 1 with a parseable body;
- returns `(nil, *SparkrunError{..., Details: {"raw": []byte}})` on exit 1 with a non-JSON body;
- returns `(nil, *SparkrunError{..., Details: {"stderr": string}})` on exit 1 with an empty body.

Callers that don't care about the body can simply treat any non-nil error as "not running". Callers that want the body should `AsSparkrunError` and read `Details["result"]`.

### `cliClient` (`cli_client.go`)

Unexported type. Fields unexported. Options pattern.

```go
type cliClient struct {
    binaryPath string
    env        []string
    timeout    time.Duration
    maxBuf     *int
}

type CliOption func(*cliClient)

func WithBinaryPath(path string) CliOption          { ... }
func WithEnv(env []string) CliOption                { ... }
func WithStreamTimeout(d time.Duration) CliOption   { ... }
func WithStreamBuffer(size *int) CliOption          { ... }

func NewCliClient(opts ...CliOption) Client { return &cliClient{ /* defaults */ } }
```

`NewCliClient` returns the `Client` interface so `cliClient` stays unexported.

Three internal helpers split the verb surface:

- `runCmd(ctx, args...)` → `(stdout, stderr, error)`. Always classifies non-zero exits as `ErrorKindExit` (or `ErrorKindUsage` on click exit-2) and `ErrorKindExec` on start failures. Captures both stdout and stderr into `Details`.
- `jsonCmd(ctx, &out, args...)` → appends `--json`, calls `runCmd`, decodes stdout into `out`. Empty body is a no-op. Surfaces `ErrorKindParse` only when the body is non-empty but undecodable.
- `plainCmd(ctx, args...)` → calls `runCmd`, discards stdout (kept only in `Details` on error). For verbs that don't have `--json`.

### Streaming commands — `ClusterMonitor`, `Logs`, `TuneVllm`, `TuneSglang`

Raw streaming via `util.BufferedPipe`, which caps the write-side buffer so a slow
reader can't block the subprocess writer indefinitely. Both stdout and stderr use
separate pipes sharing the same buffer-size configuration. Returns
`(stdout, stderr io.Reader, err error)`. The subprocess is terminated when the
command exits or when the provided context is cancelled (`exec.CommandContext`
kills on `ctx.Done`), so callers read until `io.EOF` rather than closing the
streams. A non-zero exit returns an error on `Read()` after any consumed output.

```go
func (c *cliClient) streamCmd(ctx context.Context, subcommand string, args ...string) (stdout, stderr io.Reader, err error) {
    cmd := exec.CommandContext(ctx, c.binaryPath, args...)
    if len(c.env) > 0 { cmd.Env = append(os.Environ(), c.env...) }

    stdoutPipe := util.NewBufferedPipe(util.BufferedPipeOptions{Ctx: ctx, MaxBuf: c.maxBuf})
    stderrPipe := util.NewBufferedPipe(util.BufferedPipeOptions{Ctx: ctx, MaxBuf: c.maxBuf})
    cmd.Stdout = stdoutPipe
    cmd.Stderr = stderrPipe

    if err := cmd.Start(); err != nil { ... }
    go func() {
        waitErr := cmd.Wait()
        if waitErr != nil {
            stdoutPipe.Close()
            stderrPipe.Close()
            return
        }
        stdoutPipe.Close()
        stderrPipe.Close()
    }()
    return stdout, stderrPipe, nil
}
```

The caller reads until they get either `io.EOF` (clean exit) or a `*SparkrunError`
(non-zero exit). There is no explicit close: cancelling the context kills the
subprocess, which `exec.CommandContext` turns into a non-zero exit on `Wait`.

Buffer size is configured via `CliClientOptions.WithStreamBuffer(*int)`. When nil,
the pipe grows without bound (legacy behavior). A 1 MB cap is recommended for
long-running streams like `ClusterMonitor` to keep memory bounded.

### Host fields: replacement vs. plain list

Host fields come in two shapes:

- Replacement verbs (`cluster create`, `cluster check-job`) take an `HostsList` (`Hosts(...)` or `HostsFromFile(...)`), never a raw `[]string` — see "Host list / HostsDiff overload" above.
- Verbs that merely select an existing host set (`cluster status`, `cluster monitor`) take a plain `Hosts []string` (comma-joined into a single `--hosts` flag). Callers that need a single host pass a one-element slice. There is no single-string `HostsFile` field — file-based selection is only expressed via `HostsFromFile(path)` on the replacement verbs.

Applied to:
- `ClusterCheckJobOptions.Hosts` — `HostsList`
- `ClusterCreateOptions.Hosts` — `HostsList`
- `ClusterUpdateOptions.Hosts` — `HostsDiff`
- `ClusterMonitorOptions.Hosts` — `[]string`
- `ClusterStatusOptions.Hosts` — `[]string`
- Options bags per command, not merged.
- Mandatory → non-pointer; optional → pointer; lists/maps → bare (per Code Style).
- Repeatable `--executor-opt key=value` (and other `-o`-style) options are expressed via `KeyValueOpts` (`key_value_opts.go`), an [ArgProvider] over a `map[string]string` that emits one `--<name> key=value` per entry in sorted-key order (deterministic). The flag name comes from the field's `cmd:"..."` tag. Bare `--executor`/`--scheduler` selectors and `--max-gpu-mem-util` scalar are plain pointer fields.
- `StopOptions` and `StopAllOptions` are *separate* types (no `--all` in `StopOptions`; no target in `StopAllOptions`).
- `ClusterDelete` always passes `--force`: a non-interactive caller has no way to answer sparkrun's [y/N] prompt. If a future caller needs the prompt, an options bag is reintroduced then.
- `LogsOptions.Tail *int` → `--tail N`.
- `ClusterStatusOptions.DryRun` and `--json` are mutually exclusive (per sparkrun); the wrapper omits `--json` when `DryRun` is set and returns `(nil, nil)` for the result.
- Result structs live in `cluster_types.go` (cluster scope) or alongside their method (future verbs).
- Commands with `--json` always pass it; `cliClient` parses into the typed result.

### Importing clusters (`cluster_import.go`)

`sparkrun cluster import` imports an external cluster config (e.g. a
spark-vllm-docker `.env`) into a sparkrun cluster. It takes a provider
subcommand (`svd` | `eugr`, plus plugin-provided ones) and a positional
path. `ClusterImport` currently hardcodes the `svd` provider and its
`Path` positional; `--name`/`--default`/`--dry-run` are the options
bag. The deprecated `--from-spark-vllm-docker-env` form is folded into
the `svd` path.

### JSON shapes

Each new `--json` command gets a doc file under `0-sparky-docs/Sparkrun output shapes/`.

Captured so far:
- `Cluster show.md` — includes the `user` field (recently added by sparkrun; older docs may omit it).
- `Cluster status.md` — flat top-level with `groups` / `solo_entries` / `idle_hosts` / `pending_ops` / `errors` / `total_containers` / `host_count`.
- `Cluster list.md` — bare JSON array of `Cluster show` shape. Empty list is `[]`.
- `Cluster default.md` — `Cluster show` shape, or literal `null` if no default is set.
- `Cluster check-job.md` — JSON emitted even on exit 1; field semantics documented.

Deferred:
- `Cluster monitor.md` — NDJSON; complex; will be captured when a real consumer (WebUI?) needs it.

## Per-command checklist (stage 2)

1. `cluster show` — **DONE** ✅
2. `cluster status` — **DONE** ✅
3. `cluster list`, `cluster default` — **DONE** ✅
4. `run`, `stop`, `stop --all`, `logs` — **Priority 3** (manager lifecycle).
5. `recipe list/search/show/vram/validate` — **Priority 4** (WebUI).
6. `cluster create/delete/update/set-default/unset-default/import/check-job/monitor` — **DONE** ✅
7. `registry*`, `proxy*`, `export*`, `arena*`, `benchmark*`, `setup*`, `tune*`, `update` — **Priority 6+** as consumers need them.

For each not-yet-implemented verb: add method to `interface.go`, implement on `cliClient` (with `--json` parsing where applicable), capture JSON shape doc.

## Verification

Each pass is smoke-tested against the live `default` cluster by a small program in `/tmp/opencode/smoke*/main.go`:

- `ClusterShow`/`ClusterList`/`ClusterDefault` — read-only probes.
- `ClusterCreate` + `ClusterDelete` — throwaway cluster names; round-trips the description/user fields.
- `ClusterSetDefault`/`ClusterUnsetDefault` — snapshots and restores the current default.
- `ClusterUpdate` — snapshots host list, exercises every `HostsDiff` mode, asserts round-trip.
- `ClusterCheckJob` — both error paths: nonexistent-recipe (stderr-only) and real-cluster-id (JSON body even on exit 1).
- `ClusterMonitor` — opens the stream, reads 2 NDJSON lines, decodes each.

No unit tests this pass — `MockClient` will be added when a concrete `agent` or `manager` test needs it. The options-bag argument emission and stream handling are covered by unit tests in `cluster_create_test.go`, `cluster_update_test.go`, `cluster_check_job_test.go`, and `cluster_monitor_test.go`.
