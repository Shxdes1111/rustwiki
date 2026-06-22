# Rust Wiki — база знаний по игре Rust

## Локальный запуск (всё на одной машине)

```bash
docker compose up
```

- **Фронт:** http://localhost:5173
- **Бэкенд:** http://localhost:8080
- **PG Admin:** http://localhost:5050 (admin@admin.com / admin)

## Тест на телефоне (локальная сеть)

1. Добавить `host: '0.0.0.0'` в `frontend/vite.config.ts` (уже добавлено)
2. `.env`: `VITE_API_BASE=` (пусто, через Vite proxy)
3. Телефон в той же WiFi — открыть `http://<IP-ПК>:5173`
4. Найти IP: `ip a` или `ifconfig`

## Показ через интернет (Cloudflare Tunnel)

```bash
# Установить cloudflared
curl -sL https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64 -o ~/.local/bin/cloudflared
chmod +x ~/.local/bin/cloudflared

# Запустить туннель (проксирует фронт, Vite сам проксит /api и /uploads на бэк)
cloudflared tunnel --url http://localhost:5173
```

Туннель выдаст URL вида `https://что-то.trycloudflare.com` — его можно открыть с телефона/показать кому угодно. Один туннель для всего (фронт + API через Vite proxy).

## Актуальный туннель

<!-- обновлять при перезапуске -->
- **Текущий URL:** https://grace-conduct-beast-fluid.trycloudflare.com