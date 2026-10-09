# 8. 雲端部署設計與維運

## 8.1 拓撲

```mermaid
flowchart LR
  gh[GitHub<br/>geodown-ops/bodhi-alliance<br/>main] -- push 觸發 --> rw
  subgraph rw[Railway project bodhi-alliance · production]
    web[web<br/>Caddy]
    admin[admin<br/>Caddy]
    api[api<br/>Go]
    guide[guide<br/>Go]
    pg[(Postgres)]
    api --> pg
    guide --> pg
  end
  wix[Wix DNS<br/>sunnylife.world] -- CNAME www --> web
  wix -- CNAME admin --> admin
  wix -- 轉址 apex → www --> web
```

| 服務 | 建置 | Root Directory | Watch Paths | 健康檢查 | 公開網域 |
| --- | --- | --- | --- | --- | --- |
| `Postgres` | Railway PostgreSQL 外掛（本機與 CI 用 16 版） | — | — | — | 只有內網 |
| `api` | `server/Dockerfile`，build arg `SERVICE=api` | `/server` | `/server/**` | `/healthz` | `api-production-9e87.up.railway.app` |
| `guide` | `server/Dockerfile`，`SERVICE=guide` | `/server` | `/server/**` | `/healthz` | `guide-production-42bb.up.railway.app` |
| `web` | `apps/Dockerfile`，`APP=web` | repo 根目錄 | `/apps/**`、`/package.json`、`/package-lock.json` | — | `www.sunnylife.world`（另有 `web-production-47350.up.railway.app`） |
| `admin` | `apps/Dockerfile`，`APP=admin` | repo 根目錄 | 同 web | — | `admin.sunnylife.world`（另有 `admin-production-d74c.up.railway.app`） |

- 四個應用服務 Restart Policy 皆為 On Failure。Railway 已停用設定檔（Config as Code），`server/railway.toml`、`apps/railway.toml` 只當參考，設定都在各服務 Settings。
- 服務變數會被當成同名 build arg 傳入 Dockerfile（`SERVICE`、`APP`、`VITE_*`）。
- Go 服務映像：`golang:1.26` 建置 → `gcr.io/distroless/static-debian12`，以 `nonroot` 執行，`GIN_MODE=release`，監聽 Railway 給的 `PORT`。
- 前端映像：`node:22` 建置 → `caddy:2-alpine`，SPA 路由一律回 `index.html`，`/assets/*` 一年 immutable 快取，zstd／gzip 壓縮。

## 8.2 環境變數（只列名稱與用途；祕密值見第 9 章的存放位置）

| 服務 | 變數 | 祕密 | 說明 |
| --- | --- | --- | --- |
| api | `SERVICE=api` | | 建置參數 |
| api | `DATABASE_URL=${{Postgres.DATABASE_URL}}` | 引用 | 資料庫連線 |
| api | `ALLOWED_ORIGINS` | | `https://www.sunnylife.world,https://sunnylife.world,https://admin.sunnylife.world` 與 web、admin 的 Railway 網址 |
| api | `TRUSTED_PROXIES=0.0.0.0/0,::/0` | | Railway 代理後才看得到真實 IP |
| api | `BOOTSTRAP_ADMIN_EMAIL` | | 第一位超級管理員的 email |
| api | `BOOTSTRAP_ADMIN_PASSWORD` | ✔ | 只在建立第一位管理員時使用，建好後應清空 |
| api | `BODHI_CHAIN_NETWORK=polygon` | | 第 6 章 |
| api | `BODHI_CHAIN_OPERATOR_KEY` | ✔ | 營運地址（金庫）私鑰 |
| api | `BODHI_CHAIN_MEMBER_SEED` | ✔ | 會員地址推導種子 |
| api | `BODHI_CHAIN_RPC`、`BODHI_CHAIN_EXPLORER`、`BODHI_CHAIN_CONTRACT` | | 選填 |
| guide | `SERVICE=guide`、`DATABASE_URL`、`ALLOWED_ORIGINS`、`TRUSTED_PROXIES` | | 同 api |
| guide | `ANTHROPIC_API_KEY` | ✔ | Claude API |
| guide | `GUIDE_MODEL`、`GUIDE_NAME` | | 選填，預設 `claude-opus-5-5`、`Sunny` |
| guide | `AZURE_SPEECH_KEY` | ✔ | 選填，雲端語音 |
| guide | `AZURE_SPEECH_REGION`、`TTS_VOICE`、`TTS_RATE` | | 選填 |
| web、admin | `APP`、`VITE_API_BASE=https://${{api.RAILWAY_PUBLIC_DOMAIN}}`、`VITE_GUIDE_BASE=https://${{guide.RAILWAY_PUBLIC_DOMAIN}}` | | 建置時寫入網頁；改了要重新部署 |

guide 必須有自己的網域，不能掛在官網的 `/guide` 底下（官網本身有 `/guide` 頁面）。

## 8.3 網域與 DNS

- 網域 `sunnylife.world` 在 **Wix** 註冊，名稱伺服器 `ns6/ns7.wixdns.net`，DNS 只能在 Wix「Domains → Manage DNS Records」修改。沒有 MX 紀錄。
- `www`：CNAME 指向 Railway web 給的目標，另有 TXT `_railway-verify.www`（Wix 主機欄只填 `_railway-verify.www`）。少了 TXT，Railway 不發 HTTPS 憑證。
- `admin`：CNAME 指向 Railway admin 給的目標，TXT `_railway-verify.admin`。
- 根網域：Wix 不能 CNAME 根網域，所以維持 Wix 的 A 紀錄，由 Wix 轉址到 `https://www.sunnylife.world/`。Railway 上的 `sunnylife.world` 自訂網域目前用不到；之後若轉到支援 CNAME 攤平的 DNS（例如 Cloudflare）才會用到。
- 憑證由 Railway 自動簽發（Let's Encrypt）；卡在驗證時可在 Railway 重試憑證。
- 實際 CNAME 目標值在 Railway 各服務 Settings → Networking。

## 8.4 發布流程

```mermaid
flowchart LR
  br[功能分支] --> pr[Pull Request] --> ci[GitHub Actions CI] --> ok{Geodown 在 thread 打 merge}
  ok --> main[合併到 main] --> build[Railway 依 Watch Paths 重建受影響的服務] --> mig[api／guide 啟動時自動遷移] --> live[上線]
```

**CI**（`.github/workflows/ci.yml`，PR 與 push main 觸發）：

| Job | 內容 |
| --- | --- |
| server | Postgres 16 服務容器；`gofmt` 檢查、`go vet ./...`、`go test -race ./...`（需要資料庫的測試各建暫時資料庫） |
| apps | Node 22；`npm ci`、`npm run build`（含 `vue-tsc`）、`npm test`（vitest） |

**資料庫遷移**：只新增、不修改已上線的遷移檔；破壞性變更（刪欄位、改型別）要分兩次發布，先讓新舊程式都能用。

## 8.5 本機開發

需要 Go 1.26、Node 22、PostgreSQL 16（或 Docker）。

```bash
cp .env.example .env            # 填入本機管理員帳密；要試 Sunny 就填自己的 ANTHROPIC_API_KEY
docker compose up -d            # db + api(:8081) + guide(:8082)
npm ci
npm run dev:web                 # http://localhost:5173
npm run dev:admin               # http://localhost:5174
```

測試：`cd server && TEST_DATABASE_URL=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable go test ./...`，前端 `npm test`。上鏈功能在本機不設金鑰時自動關閉；測試用 go-ethereum 模擬鏈。

## 8.6 維運

| 項目 | 做法 |
| --- | --- |
| 監看 | Railway 各服務 Deploy Logs；健康檢查 `/healthz`；上鏈 worker 日誌前綴 `chain:` |
| Sunny 費用 | 後台「Sunny訓練設定」看本月用量；預算預設 50 美元／月 |
| 上鏈 gas | 營運地址 POL 餘額（PolygonScan）；不足時日誌出現 `insufficient funds` |
| 備份 | Railway Postgres 的備份功能；重大遷移前先手動匯出（`pg_dump`） |
| 回滾 | Railway 對該服務選前一個成功部署 Redeploy；資料庫遷移不會自動回滾 |
| 停用 Sunny | 後台關閉對話，或清空 `ANTHROPIC_API_KEY` |
| 停用上鏈 | 清空 `BODHI_CHAIN_OPERATOR_KEY` 或 `BODHI_CHAIN_MEMBER_SEED`（**先備份**，清掉就找不回） |
| 舊招募頁 | 根目錄 `index.html` 仍在 GitHub Pages，正式站已取代 |
| Railway 殘留 | 自動匯入產生的 project `divine-reprieve` 從未部署，可由 Geodown 刪除 |
