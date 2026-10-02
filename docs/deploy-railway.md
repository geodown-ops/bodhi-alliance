# 部署到 Railway

一個 Railway project，五個服務：PostgreSQL、核心 API、AI 組長、官網、管理後台。都從 GitHub 的 `geodown-ops/bodhi-alliance` 自動部署，`main` 有新 commit 就會重新部署。

## 一次性設定

1. 在 Railway 建立新 project，加入 **PostgreSQL**。
2. 新增四個服務，來源都選 GitHub repo `geodown-ops/bodhi-alliance`，分支 `main`，各自設定：

| 服務 | Settings → Root Directory | Settings → Config file | Variables |
| --- | --- | --- | --- |
| `api` | `server` | `server/railway.toml` | `SERVICE=api`、`DATABASE_URL=${{Postgres.DATABASE_URL}}`、`BOOTSTRAP_ADMIN_EMAIL`、`BOOTSTRAP_ADMIN_PASSWORD`、`ALLOWED_ORIGINS`、`TRUSTED_PROXIES` |
| `guide` | `server` | `server/railway.toml` | `SERVICE=guide`、`DATABASE_URL=${{Postgres.DATABASE_URL}}`、`ANTHROPIC_API_KEY`、`ALLOWED_ORIGINS`、`TRUSTED_PROXIES` |
| `web` | （留空，用 repo 根目錄） | `apps/railway.toml` | `APP=web`、`VITE_API_BASE=https://${{api.RAILWAY_PUBLIC_DOMAIN}}`、`VITE_GUIDE_BASE=https://${{guide.RAILWAY_PUBLIC_DOMAIN}}` |
| `admin` | （留空） | `apps/railway.toml` | `APP=admin`、`VITE_API_BASE`、`VITE_GUIDE_BASE`（同上） |

3. 四個服務都到 Settings → Networking 按 **Generate Domain**（之後可以換成自己的網域）。
4. `ALLOWED_ORIGINS` 填官網與後台的網址，逗號分隔，例如 `https://web-production-xxxx.up.railway.app,https://admin-production-xxxx.up.railway.app`。
5. `TRUSTED_PROXIES`：Railway 前面有一層代理，限流要看 `X-Forwarded-For` 才分得出不同訪客。設 `0.0.0.0/0,::/0` 表示信任所有來源轉送的位址；Railway 的服務只能從它的代理連進來，所以這樣設是安全的。
6. 第一次部署完，用 `BOOTSTRAP_ADMIN_EMAIL`／`BOOTSTRAP_ADMIN_PASSWORD` 登入後台，到「帳號」新增其他人，再把 `BOOTSTRAP_ADMIN_PASSWORD` 從變數刪掉（帳號已建好，之後不再需要）。

## 注意

- `VITE_*` 變數是在建置時寫進網頁的，改了之後要重新部署 `web`／`admin` 才會生效。
- AI 組長的費用上限在後台「AI 組長設定」調整，預設每月 50 美元。
- 資料庫遷移由 `api` 與 `guide` 啟動時自動執行，不需要手動操作。
- 原本的招募頁（根目錄 `index.html`）仍由 GitHub Pages 提供；正式站上線後可以把網域指向 `web`，再決定招募頁要不要下架。
