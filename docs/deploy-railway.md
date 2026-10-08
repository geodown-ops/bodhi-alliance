# 部署到 Railway

一個 Railway project，五個服務：PostgreSQL、核心 API、AI 組長、官網、管理後台。都從 GitHub 的 `geodown-ops/bodhi-alliance` 自動部署，`main` 有新 commit 就會重新部署。

## 一次性設定

> 2026-10-03 已照這份步驟在 Geodown 的 Railway 建好 project `bodhi-alliance`。Railway 已停用設定檔（Config as Code），所以 `server/railway.toml`、`apps/railway.toml` 只當參考，下面的設定都直接填在各服務的 Settings。

1. 在 Railway 建立新 project，加入 **PostgreSQL**（服務名保持 `Postgres`）。
2. 新增四個空服務，名稱一定要是 `api`、`guide`、`web`、`admin`（變數裡的 `${{api.RAILWAY_PUBLIC_DOMAIN}}` 依名稱對應），各自設定：

| 服務 | Root Directory | Dockerfile Path | Healthcheck | Watch Paths | Variables |
| --- | --- | --- | --- | --- | --- |
| `api` | `/server` | （自動找 `server/Dockerfile`） | `/healthz` | `/server/**` | `SERVICE=api`、`DATABASE_URL=${{Postgres.DATABASE_URL}}`、`BOOTSTRAP_ADMIN_EMAIL`、`BOOTSTRAP_ADMIN_PASSWORD`、`ALLOWED_ORIGINS`、`TRUSTED_PROXIES` |
| `guide` | `/server` | （同上） | `/healthz` | `/server/**` | `SERVICE=guide`、`DATABASE_URL=${{Postgres.DATABASE_URL}}`、`ANTHROPIC_API_KEY`、`ALLOWED_ORIGINS`、`TRUSTED_PROXIES`；選填 `AZURE_SPEECH_KEY`、`AZURE_SPEECH_REGION`（Sunny 的雲端語音，見下方） |
| `web` | （留空，用 repo 根目錄） | `apps/Dockerfile` | — | `/apps/**`、`/package.json`、`/package-lock.json` | `APP=web`、`VITE_API_BASE=https://${{api.RAILWAY_PUBLIC_DOMAIN}}`、`VITE_GUIDE_BASE=https://${{guide.RAILWAY_PUBLIC_DOMAIN}}` |
| `admin` | （留空） | `apps/Dockerfile` | — | （同 web） | `APP=admin`、`VITE_API_BASE`、`VITE_GUIDE_BASE`（同上） |

   四個服務的 Restart Policy 都設 On Failure。`SERVICE`、`APP`、`VITE_*` 會被 Railway 當成同名的 build arg 傳進 Dockerfile。
3. 四個服務都到 Settings → Networking 按 **Generate Domain**（之後可以換成自己的網域）。
4. `ALLOWED_ORIGINS` 填官網與後台的網址，逗號分隔，例如 `https://web-production-xxxx.up.railway.app,https://admin-production-xxxx.up.railway.app`。
5. `TRUSTED_PROXIES`：Railway 前面有一層代理，限流要看 `X-Forwarded-For` 才分得出不同訪客。設 `0.0.0.0/0,::/0` 表示信任所有來源轉送的位址；Railway 的服務只能從它的代理連進來，所以這樣設是安全的。
6. 最後把四個服務的來源接到 GitHub repo `geodown-ops/bodhi-alliance`、分支 `main`，就會開始第一次建置。
7. 第一次部署完，用 `BOOTSTRAP_ADMIN_EMAIL`／`BOOTSTRAP_ADMIN_PASSWORD` 登入後台，到「帳號」新增其他人，再把 `BOOTSTRAP_ADMIN_PASSWORD` 從變數刪掉（帳號已建好，之後不再需要）。

## 自訂網域 sunnylife.world

2026-10-05 起官網改用 `sunnylife.world`（原本是協會在 Wix 上的網站）。

1. `web` 服務已加上自訂網域 `www.sunnylife.world` 與 `sunnylife.world`，Railway 會給每個網域一筆 CNAME 和一筆 `_railway-verify` TXT（值在 web → Settings → Networking）。
2. `api`、`guide` 的 `ALLOWED_ORIGINS` 已加上 `https://sunnylife.world,https://www.sunnylife.world`。
3. 網域是在 Wix 買的，DNS 只能在 Wix 的「Domains → Manage DNS Records」改：`www` 的 CNAME 改成 Railway 給的值，再加上 `_railway-verify.www` 的 TXT（Wix 的主機名稱欄只填 `_railway-verify.www`）。少了 TXT，Railway 不會核發 HTTPS 憑證。
4. Wix 不能改名稱伺服器，根網域也不能設 CNAME，而 Railway 的根網域需要 CNAME（或 ALIAS），所以不帶 www 的 `sunnylife.world` 仍指向 Wix，由 Wix 轉址到 `https://www.sunnylife.world/`。Railway 上的 `sunnylife.world` 自訂網域因此用不到；之後若把網域轉到支援 CNAME 攤平的服務（例如 Cloudflare），才會用到它。

## 注意

- `VITE_*` 變數是在建置時寫進網頁的，改了之後要重新部署 `web`／`admin` 才會生效。
- AI 組長的費用上限在後台「AI 組長設定」調整，預設每月 50 美元。
- 資料庫遷移由 `api` 與 `guide` 啟動時自動執行，不需要手動操作。
- 原本的招募頁（根目錄 `index.html`）仍由 GitHub Pages 提供；正式站上線後可以把網域指向 `web`，再決定招募頁要不要下架。

## Sunny 的雲端語音（選填）

`guide` 設定 `AZURE_SPEECH_KEY` 後，Sunny 改用 Azure AI Speech 的臺灣華語神經語音唸回答（預設 `zh-TW-HsiaoChenNeural`），比瀏覽器內建語音自然；沒有設定、或某一句取不到時，網站自動改用瀏覽器語音。

| 變數 | 預設 | 說明 |
|---|---|---|
| `AZURE_SPEECH_KEY` | （空，表示不用雲端語音） | Azure 入口網站 → 建立「語音服務」資源 → 金鑰與端點 → 金鑰 1 |
| `AZURE_SPEECH_REGION` | `eastasia` | 該資源的區域（位置） |
| `TTS_VOICE` | `zh-TW-HsiaoChenNeural` | 也可試 `zh-TW-HsiaoYuNeural`（女）、`zh-TW-YunJheNeural`（男） |
| `TTS_RATE` | （空，正常語速） | SSML 語速，例如 `-5%` 慢一點 |

免費方案（F0）每月 50 萬字；網站一次送一句，`/guide/tts` 每個 IP 有頻率限制，組長關閉或超過每月預算時也一併停用。


## 菩提幣上鏈（api 服務）

`api` 設定下面兩個金鑰後就會啟用入會贈幣上鏈；沒設定時個人頁不顯示鏈上區塊。金鑰只放在 Railway 變數，不要寫進程式碼庫。

| 變數 | 預設 | 說明 |
| --- | --- | --- |
| `BODHI_CHAIN_OPERATOR_KEY` | （空，表示關閉） | 營運地址的私鑰（hex）。測試鏈上它也是金庫；地址要有一點 POL 付手續費 |
| `BODHI_CHAIN_MEMBER_SEED` | （空，表示關閉） | 32 位元組以上的 hex 密鑰，會員地址由它推導；要另外備份，換掉會員地址就會變 |
| `BODHI_CHAIN_RPC` | `https://rpc-amoy.polygon.technology` | 鏈的 JSON-RPC 節點 |
| `BODHI_CHAIN_EXPLORER` | `https://amoy.polygonscan.com` | 區塊瀏覽器，個人頁的連結用 |
| `BODHI_CHAIN_CONTRACT` | （空，表示自動部署） | 已部署的 BodhiCoin 地址；空的時候 api 自己部署一次並記在資料庫 |

測試鏈的手續費用 Polygon 官方水龍頭（faucet.polygon.technology）領 Amoy POL 到營運地址。上 mainnet 前要換新的金鑰、改用多簽金庫，並等法務結論。
