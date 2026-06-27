# sparky config and agent plan

grug write plan for sparky. sparky is model router and dashboard for dgx spark cluster.

## goals

- keep complexity demon away.
- no ongoing ssh connections. ssh is slow dinosaur.
- passive agents on nodes. manager does polling via http.

## how things fit together

### host discovery

manager run local sparkrun command on startup:
`sparkrun cluster show default --json`

this give hosts list and ssh user. no manual ip list needed.

### monitoring and health

each node run sparky agent. agent serve http on port 8081.

agent endpoint:
`GET /status`
agent execute local command:
`sparkrun cluster status -H localhost --json`
this return local containers and status. no ssh needed.

manager polls agent:
`GET http://<node>:8081/status`
if http fail, node is offline. if succeed, parse json to update dashboard.

### ssh verification

on startup, manager run quick one-off check:
`ssh -o ConnectTimeout=2 <user>@<host> "echo OK"`
this warn user if passwordless ssh not working.

## systemd configuration

### agent service

```ini
[Unit]
Description=Sparky Agent
After=network.target

[Service]
Type=simple
User=sparky
Environment="SPARKY_AGENT_PORT=8081"
ExecStart=/usr/local/bin/sparky agent
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

### manager service

```ini
[Unit]
Description=Sparky Manager
After=network.target

[Service]
Type=simple
User=clusteruser
Environment="SPARKY_CLUSTER_NAME=default"
Environment="SPARKY_MANAGER_PORT=8080"
ExecStart=/usr/local/bin/sparky manager
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```
