# proxy-cli

A lightweight CLI tool written in Go that runs a local HTTP proxy with **domain-based traffic routing**.

Route traffic for specific domains through an upstream proxy, while everything else connects directly. Configured via a simple JSON file.

Русская версия: [README_RU.md](./README_RU.md)

---

## How It Works

```
Your App → proxy-cli (127.0.0.1:5555)
                │
                ├─ domain matches rules → upstream proxy
                └─ no match            → direct connection
```

Point your browser or application at `127.0.0.1:5555` as an HTTP proxy. proxy-cli intercepts each `CONNECT` request, checks the destination host against your configured rules, and routes accordingly.

---

## Install

```bash
git clone https://github.com/Lesakez/proxy-cli
cd proxy-cli
go build -o proxy-cli.exe .
```

---

## Usage

```
proxy-cli [OPTIONS]

  --port    int     Local port to listen on (default: 5555)
  --config  string  Path to config file    (default: ~/.proxy-cli/config.json)
  --version         Print version and exit
```

```bash
proxy-cli --port 5555 --config ./config.json
```

---

## Configuration

```json
[
  {
    "name": "HTTP Proxy",
    "enabled": true,
    "scheme": "HTTP",
    "host": "proxy.example.com",
    "port": 8080,
    "auth": {
      "credentials": { "username": "user", "password": "pass" },
      "token": ""
    },
    "rules": [
      {
        "name": "YouTube",
        "hosts": ["youtu.be", "*.youtube.com", "*.googlevideo.com"]
      }
    ]
  },
  {
    "name": "SOCKS5 Proxy",
    "enabled": true,
    "scheme": "SOCKS5",
    "host": "socks.example.com",
    "port": 1080,
    "auth": {
      "credentials": { "username": "user", "password": "pass" },
      "token": ""
    },
    "rules": [
      { "name": "Discord", "hosts": ["*discord*.*"] },
      { "name": "Check IP", "hosts": ["api.ipify.org"] }
    ]
  }
]
```

Multiple proxy entries are supported — each with its own set of rules.  
Set `"enabled": false` to disable a proxy without removing it.

### Supported Schemes

| `scheme` | Protocol |
|----------|----------|
| `HTTP` | HTTP CONNECT tunnel |
| `HTTPS` | HTTP CONNECT tunnel over TLS |
| `SOCKS5` | SOCKS5 with optional username/password auth |

### Wildcard Patterns

| Pattern | Matches |
|---------|---------|
| `*.youtube.com` | `www.youtube.com`, `m.youtube.com` |
| `*discord*.*` | `discord.com`, `cdn.discordapp.com` |
| `api.ipify.org` | `api.ipify.org` (exact) |

### Authentication

| Method | Config field |
|--------|-------------|
| Basic auth (HTTP) | `auth.credentials.username` + `password` |
| Token (HTTP) | `auth.token` → `Proxy-Authorization: Bearer <token>` |
| SOCKS5 auth | `auth.credentials.username` + `password` |

---

## Project Structure

```
proxy-cli/
├── main.go              # Entry point, CLI flags
├── config/
│   ├── config.go        # Config structs and JSON loader
│   └── config_test.go
├── filter/
│   ├── filter.go        # Wildcard host matching
│   └── filter_test.go
├── proxy/
│   ├── proxy.go         # Local TCP server, goroutine per connection
│   └── tunnel.go        # HTTP CONNECT tunnel + bidirectional pipe
└── config.example.json
```

---

## Tests

```bash
go test ./...
```

---

## License

MIT
