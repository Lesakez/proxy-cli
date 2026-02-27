# proxy-cli

Локальный HTTP прокси-сервер с маршрутизацией трафика по доменам.

Запускаешь `proxy-cli`, указываешь его как HTTP прокси в браузере или приложении (`127.0.0.1:5555`) — и трафик к нужным доменам уходит через настроенный прокси, остальное идёт напрямую.

---

## Установка

```bash
git clone https://github.com/Lesakez/proxy-cli
cd proxy-cli
go build -o proxy-cli.exe .
```

---

## Использование

```
proxy-cli [OPTIONS]

  --port    int     Локальный порт (по умолчанию: 5555)
  --config  string  Путь к конфигу (по умолчанию: ~/.proxy-cli/config.json)
  --version         Вывести версию
```

```bash
proxy-cli --port 5555 --config ./config.json
```

---

## Конфигурация

```json
[
  {
    "name": "Мой прокси",
    "enabled": true,
    "scheme": "SOCKS5",
    "host": "proxy.example.com",
    "port": 1080,
    "auth": {
      "credentials": { "username": "user", "password": "pass" },
      "token": ""
    },
    "rules": [
      {
        "name": "YouTube",
        "hosts": ["youtu.be", "*.youtube.com", "*.googlevideo.com"]
      },
      { "name": "Discord", "hosts": ["*discord*.*"] }
    ]
  }
]
```

Можно добавить несколько прокси — у каждого свои правила.  
`"enabled": false` — отключает прокси без удаления из конфига.

### Поддерживаемые схемы

| `scheme` | Протокол |
|----------|----------|
| `HTTP` | HTTP CONNECT туннель |
| `HTTPS` | HTTP CONNECT туннель |
| `SOCKS5` | SOCKS5 с авторизацией по логину/паролю |

### Wildcard-паттерны

| Паттерн | Совпадает с |
|---------|-------------|
| `*.youtube.com` | `www.youtube.com`, `m.youtube.com` |
| `*discord*.*` | `discord.com`, `cdn.discordapp.com` |
| `api.ipify.org` | `api.ipify.org` (точное совпадение) |

### Авторизация

| Метод | Поле конфига |
|-------|-------------|
| Логин/пароль (HTTP) | `auth.credentials.username` + `password` |
| Токен (HTTP) | `auth.token` → `Proxy-Authorization: Bearer <token>` |
| Логин/пароль (SOCKS5) | `auth.credentials.username` + `password` |

---

## Структура проекта

```
proxy-cli/
├── main.go              # Точка входа, CLI флаги
├── config/
│   ├── config.go        # Структуры конфига и загрузчик JSON
│   └── config_test.go
├── filter/
│   ├── filter.go        # Wildcard-матчинг доменов
│   └── filter_test.go
├── proxy/
│   ├── proxy.go         # Локальный TCP сервер, горутина на соединение
│   └── tunnel.go        # HTTP CONNECT / SOCKS5 туннель
└── config.example.json
```

---

## Тесты

```bash
go test ./...
```

---

## Лицензия

MIT
