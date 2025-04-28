# listmonk proxy

To set

```
LISTMONK_URL

API_USER

API_TOKEN

HINDI_LIST

ENGLISH_LIST
```

Endpoints

```
ENG to HIN: /proxy/switch_to_hin
HIN to ENG: /proxy/switch_to_english
BOTH: /proxy/use_both
```

## Instructions

1. Setup listmonk binary
2. Add config.toml for local listmonk instance
```toml
[app]
# Interface and port where the app will run its webserver.  The default value
# of localhost will only listen to connections from the current machine. To
# listen on all interfaces use '0.0.0.0'. To listen on the default web address
# port, use port 80 (this will require running with elevated permissions).
address = "0.0.0.0:80"

# Database.
[db]
host = "localhost"
port = 5432
user = "listmonk"
password = "listmonk"

# Ensure that this database has been created in Postgres.
database = "listmonk"

ssl_mode = "disable"
max_open = 25
max_idle = 25
max_lifetime = "300s"

# Optional space separated Postgres DSN params. eg: "application_name=listmonk gssencmode=disable"
params = ""
```
3. Create Systemd service for listmonk
```
[Unit]
Description=Secure Proxy for Listmonk Newsletter Manager
After=network.target listmonk.service

[Service]
Environment="API_USER=apiuser2"
Environment="LISTMONK_URL=0.0.0.0:80"
#Environment="API_TOKEN=gNMAnvuMN7gcoaRYlPqADZGEqVxtBDqG"
Environment="API_TOKEN=C2iqFzpXp39R5k9bT7YjLCb8dR9lrVWn"
Environment="HINDI_LIST_1M=5"
Environment="HINDI_LIST_3M=7"
Environment="ENGLISH_LIST_3M=8"
Environment="ENGLISH_LIST_1M=6"
Type=simple
User=root
Group=root
WorkingDirectory=/home/ubuntu
ExecStart=/home/ubuntu/listmonk_proxy
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

4. Create Systemd service for listmonk_proxy
```
[Unit]
Description=Secure Proxy for Listmonk Newsletter Manager
After=network.target listmonk.service

[Service]
Environment="API_USER=apiuser2"
Environment="LISTMONK_URL=0.0.0.0:80"
#Environment="API_TOKEN=gNMAnvuMN7gcoaRYlPqADZGEqVxtBDqG"
Environment="API_TOKEN=C2iqFzpXp39R5k9bT7YjLCb8dR9lrVWn"
Environment="HINDI_LIST_1M=5"
Environment="HINDI_LIST_3M=7"
Environment="ENGLISH_LIST_3M=8"
Environment="ENGLISH_LIST_1M=6"
Type=simple
User=root
Group=root
WorkingDirectory=/home/ubuntu
ExecStart=/home/ubuntu/listmonk_proxy
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```
