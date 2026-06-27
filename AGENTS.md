# Sparky

Go CLI tool (`module sparky`, Go 1.26.3, uses `github.com/spf13/cobra`) for managing LLMs across a DGX Spark cluster. Two binaries: **manager** (web UI + router, port 8080) and **agent** (per-node health reporter, port 8081).

## Architecture

- **Manager** — runs on one host. Discovers cluster nodes by executing `sparkrun cluster show default --json` on startup. Polls each agent via `GET http://<node>:8081/v1/health`. Exposes OpenAI-compatible API (`/v1/chat/completions`, `/v1/models`, `/v1/completions`) and manager APIs (`/v1/config/reload`, `/v1/agents`, `/v1/agents/health`). Auto-loads models to available nodes; evicts after configurable idle time (default 15 min) to prevent thrashing.
- **Agent** — runs on each DGX Spark node. Self-discovers by running `sparkrun cluster status -H localhost --json` every 30-60s. Serves `/v1/health` exposing loaded models and VRAM usage. Reports anomalies (crashed models) to manager via `POST /v1/agents/health`.

## Critical external dependency

`sparkrun` (DGX Spark CLI) is required at runtime — both manager and agent invoke it. Model recipes are referenced by address (e.g. `@official/llama-3-70b`). All `sparkrun` interaction is process-execution; the codebase must not assume its API contract without checking.

## Commands

- `go build -o sparky ./...` — build
- `sparky agent` — run agent (needs `SPARKY_AGENT_PORT`, default 8081)
- `sparky manager` — run manager (needs `SPARKY_CLUSTER_NAME=default`, `SPARKY_MANAGER_PORT`, default 8080)

## Agent discovery flow

1. Manager shells out `sparkrun cluster show default --json` → gets host list + SSH user
2. Manager SSH-checks `echo OK` against each host (connect timeout 2s)
3. Manager polls `GET http://<node>:8081/v1/health` every 30-60s
4. If HTTP fails, node marked offline

## Deployment

systemd services under `/etc/systemd/system/`:
- `sparky-agent.service` — runs as user `sparky`, ExecStart `/usr/local/bin/sparky agent`
- `sparky-manager.service` — runs as user `clusteruser`, ExecStart `/usr/local/bin/sparky manager`

## Gotchas

- **No SSH from manager to agents after startup** — all monitoring is HTTP polling. Don't add ongoing SSH connections.
- **Model load is slow** — time-gate eviction is mandatory. Without it, models thrash in/out under memory pressure.
- **Agent uses local IP as agent_id** — this is how the manager identifies the node.
- **`go.mod` is bare** — no lint/test/typecheck config exists yet. Add if needed.

## Key files

- `sparky-docs/openapi.yaml` — canonical API surface (agents should reference this, not the code)
- `sparky-docs/Plan.md` — architecture rationale
- `implementation_plan.md` — systemd configs and discovery flow summary
- `.opencode/opencode.json` — plugins (websearch, dcp)