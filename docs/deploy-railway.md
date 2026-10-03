# 部署到 Railway

一個 Railway project，五個服務：PostgreSQL、核心 API、AI 組長、官網、管理後台。都從 GitHub 的 `geodown-ops/bodhi-alliance` 自動部署，`main` 有新 commit 就會重新部署。

## 一次性設定

> 2026-10-03 已照這份步驟在 Geodown 的 Railway 建好 project `bodhi-alliance`。Railway 已停用設定檔（Config as Code），所以 `server/railway.toml`、`apps/railway.toml` 只當參考，下面的設定都直接填在各服務的 Settings。

1. 在 Railway 建立新 project，加入 **PostgreSQL**（服務名保持 `Postgres`）。
2. 新增四個空服務，名稱一定要是 `api`、`guide`、`web`、`admin`（變數裡的 `${{api.RAILWAY_PUBLIC_DOMAIN}}` 依名稱對應），各自設定：

| 服務 | Root Directory | Dockerfile Path | Healthcheck | Watch Paths | Variables |
| --- | --- | --- | --- | --- | --- |
| `api` | `/server` | （自動找 `server/Dockerfile`） | `/healthz` | `/server/**` | `SERVICE=api`、`DATABASE_URL=${{Postgres.DATABASE_URL}}`、`BOOTSTRAP_ADMIN_EMAIL`、`BOOTSTRAP_ADMIN_PASSWORD`、`ALLOWED_ORIGINS`、`TRUSTED_PROXIES` |
| `guide` | `/server` | （同上） | `/healthz` | `/server/**` | `SERVICE=guide`、`DATABASE_URL=${{Postgres.DATABASE_URL}}`、`ANTHROPIC_API_KEY`、`ALLOWED_ORIGINS`、`TRUSTED_PROXIES` |
| `web` | （留空，用 repo 根目錄） | `apps/Dockerfile` | — | `/apps/**`、`/package.json`、`/package-lock.json`、`/index.html` | `APP=web`、`VITE_API_BASE=https://${{api.RAILWAY_PUBLIC_DOMAIN}}`、`VITE_GUIDE_BASE=https://${{guide.RAILWAY_PUBLIC_DOMAIN}}` |
| `admin` | （留空） | `apps/Dockerfile` | — | （同 web） | `APP=admin`、`VITE_API_BASE`、`VITE_GUIDE_BASE`（同上） |

   四個服務的 Restart Policy 都設 On Failure。`SERVICE`、`APP`、`VITE_*` 會被 Railway 當成同名的 build arg 傳進 Dockerfile。
3. 四個服務都到 Settings → Networking 按 **Generate Domain**（之後可以換成自己的網域）。
4. `ALLOWED_ORIGINS` 填官網與後台的網址，逗號分隔，例如 `https://web-production-xxxx.up.railway.app,https://admin-production-xxxx.up.railway.app`。
5. `TRUSTED_PROXIES`：Railway 前面有一層代理，限流要看 `X-Forwarded-For` 才分得出不同訪客。設 `0.0.0.0/0,::/0` 表示信任所有來源轉送的位址；Railway 的服務只能從它的代理連進來，所以這樣設是安全的。
6. 最後把四個服務的來源接到 GitHub repo `geodown-ops/bodhi-alliance`、分支 `main`，就會開始第一次建置。
7. 第一次部署完，用 `BOOTSTRAP_ADMIN_EMAIL`／`BOOTSTRAP_ADMIN_PASSWORD` 登入後台，到「帳號」新增其他人，再把 `BOOTSTRAP_ADMIN_PASSWORD` 從變數刪掉（帳號已建好，之後不再需要）。

## 注意

- `VITE_*` 變數是在建置時寫進網頁的，改了之後要重新部署 `web`／`admin` 才會生效。
- AI 組長的費用上限在後台「AI 組長設定」調整，預設每月 50 美元。
- 資料庫遷移由 `api` 與 `guide` 啟動時自動執行，不需要手動操作。
- 官網 `web` 的首頁就是根目錄原本整頁的招募頁 `index.html`（建置時複製成 `recruit.html`，由 `apps/Caddyfile` 在 `/` 提供），其他頁面照常走 Vue。GitHub Pages 也仍提供同一份招募頁。
