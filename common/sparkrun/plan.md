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
  sparkrun.go          package doc, default-binary constant
  interface.go         Client interface (one method per verb)
  cli_client.go        cliClient — real impl (unexported type, options pattern)
  errors.go            SparkrunError + ErrorKind enum
  target.go            RecipeNameOrFile, RecipeNameOrJobID, constructors
  run.go               Run + RunOptions
  stop.go              Stop + StopOptions
                       StopAll + StopAllOptions
  logs.go              Logs + LogsOptions + LogsResult (io.ReadCloser)
  status.go            Status + StatusOptions + StatusResult
  cluster.go           ClusterShow/Status/Create/Delete/Default/Update/...
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

## Design

### Client interface (`interface.go`)
Every CLI verb is a method. Methods take `context.Context`. Required positional args are direct parameters; options are a varargs options bag. `Stop` and `StopAll` are separate methods because `--all` is a distinct mode of `stop`.

```go
type Client interface {
    // lifecycle
    Run(ctx context.Context, recipe RecipeNameOrFile, opts ...RunOptions) error
    Stop(ctx context.Context, target RecipeNameOrJobID, opts ...StopOptions) error
    StopAll(ctx context.Context, opts ...StopAllOptions) error
    Logs(ctx context.Context, target RecipeNameOrJobID, opts ...LogsOptions) (io.ReadCloser, error)

    // cluster discovery / management
    Status(ctx context.Context, opts StatusOptions) (*StatusResult, error)
    ClusterShow(ctx context.Context, name string) (*ClusterShowResult, error)
    ClusterStatus(ctx context.Context, opts ClusterStatusOptions) (*ClusterStatusResult, error)
    ClusterDefault(ctx context.Context) (*ClusterDefaultResult, error)
    ClusterList(ctx context.Context) (*ClusterListResult, error)
    ClusterCreate(ctx context.Context, name string, opts ...ClusterCreateOptions) (*ClusterCreateResult, error)
    ClusterDelete(ctx context.Context, name string, opts ...ClusterDeleteOptions) error
    ClusterUpdate(ctx context.Context, name string, opts ClusterUpdateOptions) error
    ClusterSetDefault(ctx context.Context, name string) error
    ClusterUnsetDefault(ctx context.Context) error
    ClusterCheckJob(ctx context.Context, target RecipeNameOrJobID, opts ...ClusterCheckJobOptions) (*ClusterCheckJobResult, error)
    ClusterMonitor(ctx context.Context, opts ClusterMonitorOptions) (io.ReadCloser, error)

    // recipe discovery / introspection
    RecipeList(ctx context.Context, opts RecipeListOptions) (*RecipeListResult, error)
    RecipeSearch(ctx context.Context, query string, opts ...RecipeSearchOptions) (*RecipeSearchResult, error)
    RecipeShow(ctx context.Context, recipe RecipeNameOrFile, opts ...RecipeShowOptions) (*RecipeShowResult, error)
    RecipeVram(ctx context.Context, recipe RecipeNameOrFile, opts ...RecipeVramOptions) (*RecipeVramResult, error)
    RecipeValidate(ctx context.Context, recipe RecipeNameOrFile) (*RecipeValidateResult, error)

    // registry, proxy, export, arena, benchmark, setup, tune, update — same pattern
}
```

`Run`, `Stop`, `StopAll`, `ClusterCreate`, `ClusterDelete`, etc. return just `error`. Caller resolves `cluster_id` via `ClusterStatus`.

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

### Error type (`errors.go`)
Categorization is the value; `Details` is for callers who want to introspect. Message always begins with literal `"sparkrun"`.

```go
type ErrorKind string

const (
    ErrorKindParse  ErrorKind = "parse"   // --json parse failure
    ErrorKindExec   ErrorKind = "exec"    // could not start subprocess
    ErrorKindExit   ErrorKind = "exit"    // non-zero exit code
    ErrorKindTarget ErrorKind = "target"  // wrong target kind for command
    ErrorKindUsage  ErrorKind = "usage"   // invalid arguments
)

type SparkrunError struct {
    Kind       ErrorKind
    Subcommand string
    Details    map[string]any
    Err        error
}

func (e *SparkrunError) Error() string {
    if e.Subcommand == "" {
        return fmt.Sprintf("sparkrun: %s: %v", e.Kind, e.Err)
    }
    return fmt.Sprintf("sparkrun %s: %s: %v", e.Subcommand, e.Kind, e.Err)
}

func (e *SparkrunError) Unwrap() error { return e.Err }
```

`Subcommand` is the full nested path: `"run"`, `"cluster status"`, `"export running-recipe"`.

Standard `Details` keys set by `cliClient`:
- `ErrorKindExit` → `{"exit_code": int, "stderr": string}` — stderr is always captured unless silenced via `WithStderr(io.Discard)`.
- `ErrorKindParse` → `{"raw": []byte}`
- `ErrorKindTarget` → `{"expected": []string, "got": string}`
- `ErrorKindExec` → `{"binary": string, "args": []string}`

Callers can type-assert `*SparkrunError` and read these; they're never in the formatted message.

### `cliClient` (`cli_client.go`)
Unexported type. Fields unexported. Options pattern.

```go
type cliClient struct {
    binaryPath string
    env        []string
    stderr     io.Writer
}

type CliOption func(*cliClient)

func WithBinaryPath(path string) CliOption { ... }
func WithEnv(env []string) CliOption       { ... }
func WithStderr(w io.Writer) CliOption     { ... } // default os.Stderr; pass io.Discard to silence

func NewCliClient(opts ...CliOption) Client { return &cliClient{ /* defaults */ } }
```

`NewCliClient` returns the `Client` interface so `cliClient` stays unexported.

### Streaming commands — `Logs`, `ClusterMonitor`, `TuneVllm`, `TuneSglang`
Raw streaming via `io.Pipe`, no buffering. A non-zero exit returns an error on `ReadCloser.Read()` after any consumed output.

```go
// shape (sketched — internal to cli_client.go)
func (c *cliClient) streamCmd(ctx context.Context, args ...string) (io.ReadCloser, error) {
    cmd := exec.CommandContext(ctx, c.binaryPath, args...)
    cmd.Env = append(os.Environ(), c.env...)

    pr, pw := io.Pipe()
    cmd.Stdout = pw

    var stderrBuf bytes.Buffer
    cmd.Stderr = &stderrBuf // or io.MultiWriter(c.stderr, &stderrBuf)

    go func() {
        waitErr := cmd.Wait()
        if waitErr != nil {
            pw.CloseWithError(&SparkrunError{
                Kind:       ErrorKindExit,
                Subcommand: <sub>,
                Err:        waitErr,
                Details:    map[string]any{"exit_code": cmd.ProcessState.ExitCode(), "stderr": stderrBuf.String()},
            })
            return
        }
        pw.Close() // signals io.EOF
    }()

    if err := cmd.Start(); err != nil { ... }
    return pr, nil
}
```

Caller reads until they get either `io.EOF` (clean exit) or a `*SparkrunError` (non-zero exit). On `Close()`, the underlying `cmd` is killed.

### Options and results
- Options bags per command, not merged.
- Mandatory → non-pointer; optional → pointer; lists/maps → bare (per Code Style).
- `StopOptions` and `StopAllOptions` are *separate* types (no `--all` in `StopOptions`; no target in `StopAllOptions`).
- `LogsOptions.Tail *int` → `--tail N`.
- Result structs defined alongside their method, in the same `.go` file.
- Commands with `--json` always pass it; `cliClient` parses into the typed result.

### JSON shapes
Each new `--json` command gets a doc file under `0-sparky-docs/Sparkrun output shapes/`. Already captured: `Cluster show.md`, `Cluster status.md`. We add the rest as we implement.

## Per-command checklist (stage 2)

1. `cluster show` — JSON shape known. **Priority 1** (manager startup).
2. `cluster status` — JSON shape known. **Priority 1** (agent health).
3. `cluster list`, `cluster default` — **Priority 2** (manager).
4. `run`, `stop`, `stop --all`, `logs` — **Priority 3** (manager lifecycle).
5. `recipe list/search/show/vram/validate` — **Priority 4** (WebUI).
6. `cluster create/delete/update/set-default/unset-default/check-job/monitor` — **Priority 5** (WebUI).
7. `registry*`, `proxy*`, `export*`, `arena*`, `benchmark*`, `setup*`, `tune*`, `update` — **Priority 6+** as consumers need them.

For each: add method to `interface.go`, implement on `cliClient` with `--json` parsing, capture JSON shape doc.
