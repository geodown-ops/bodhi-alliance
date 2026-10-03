# 菩提幣聯盟

世界佛教教育協會菩提幣系統的官網、管理後台與後端服務。架構見 [docs/architecture.md](docs/architecture.md)，資料模型見 [docs/data-model.md](docs/data-model.md)。

目前是「階段 A 官網先行」：七個選單的公開頁、覺行小組報名、共好企業登記，以及線上覺行小組 AI 組長與它的知識庫後台。錢包、核發與券要等法務結論（見架構文件〈分期〉）。

根目錄的 `index.html` 是原本的招募頁，正式站上線前繼續放在 GitHub Pages。

## 目錄

| 路徑 | 內容 |
| --- | --- |
| `apps/web` | 官網（Vue 3 + Quasar），開發時 http://localhost:5173 |
| `apps/admin` | 管理後台，開發時 http://localhost:5174 |
| `server/cmd/api` | 核心 API：登入、帳號、覺行小組、報名與登記（:8081） |
| `server/cmd/guide` | AI 組長服務：對話、帶領共修、知識庫、角色設定（:8082） |
| `server/internal` | 兩個服務共用的登入、資料庫與遷移 |

## 本機開發

需要 Go 1.26、Node 22、PostgreSQL 16（或 Docker）。

```bash
cp .env.example .env            # 填入管理員帳密；要試 AI 組長就填 ANTHROPIC_API_KEY
docker compose up -d            # 資料庫 + api + guide
npm ci
npm run dev:web                 # 官網
npm run dev:admin               # 後台，用 .env 裡的管理員帳密登入
```

不用 Docker 的話，自己起 PostgreSQL，再分別執行：

```bash
cd server
export DATABASE_URL=postgres://postgres:postgres@localhost:5432/bodhi?sslmode=disable
BOOTSTRAP_ADMIN_EMAIL=admin@example.org BOOTSTRAP_ADMIN_PASSWORD=change-me-please go run ./cmd/api
ANTHROPIC_API_KEY=... go run ./cmd/guide
```

兩個服務啟動時都會自動套用資料庫遷移（`server/internal/db/migrations`）。第一次啟動 api 時，若 `BOOTSTRAP_ADMIN_EMAIL` 還沒有帳號，會建立一個聯盟管理員。

## 測試

```bash
cd server && TEST_DATABASE_URL=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable go test ./...
npm test
```

需要資料庫的測試會為每個測試套件建一個暫時的資料庫，跑完刪掉；沒設 `TEST_DATABASE_URL` 時會跳過。

## 環境變數

| 變數 | 服務 | 說明 |
| --- | --- | --- |
| `DATABASE_URL` | 兩者 | PostgreSQL 連線字串 |
| `ADDR` | 兩者 | 監聽位址，預設 api `:8081`、guide `:8082` |
| `ALLOWED_ORIGINS` | 兩者 | 允許呼叫的前端網址，逗號分隔 |
| `TRUSTED_PROXIES` | 兩者 | 前方反向代理的 IP 或網段，逗號分隔；設了才採信 X-Forwarded-For 來限流。沒設時用連線來源 IP |
| `BOOTSTRAP_ADMIN_EMAIL`、`BOOTSTRAP_ADMIN_PASSWORD` | api | 第一位聯盟管理員 |
| `ANTHROPIC_API_KEY` | guide | 沒設時對話關閉，知識庫後台照常可用 |
| `GUIDE_MODEL` | guide | 預設 `claude-opus-5-5` |
| `GUIDE_NAME` | guide | 第一次啟動時的組長名字，預設 Sunny；之後在後台改 |
| `VITE_API_BASE`、`VITE_GUIDE_BASE` | 前端建置 | 正式環境 api 與 guide 的網址（例如 `https://api.example.org`）。官網有 `/guide` 頁面，所以 guide 服務要用自己的網域，不能和官網同網域掛在 `/guide` 底下 |

## 角色

| 角色 | 可以做什麼 |
| --- | --- |
| 聯盟管理員 | 全部：報名與登記、覺行小組、知識庫、AI 組長設定、帳號 |
| 知識管理員 | 只有知識庫與 AI 組長設定 |
