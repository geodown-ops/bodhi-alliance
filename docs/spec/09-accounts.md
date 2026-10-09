# 9. 關鍵帳號與存取

> **本章不含、也永遠不應寫入任何金鑰、密碼、私鑰或助記詞的實際值。** 只記錄帳號／服務、用途、擁有者、存放位置與環境變數名稱。祕密值只存在表中「存放位置」，並由擁有者另外備份到密碼管理器。

## 9.1 帳號總表

| # | 帳號／服務 | 用途 | 擁有者 | 存放位置／登入方式 | 相關變數 |
| --- | --- | --- | --- | --- | --- |
| 1 | GitHub 帳號 `geodown-ops` | 原始碼 repo `geodown-ops/bodhi-alliance`、PR 審核與合併、GitHub Pages（舊招募頁） | Geodown | GitHub 登入（Geodown 自管） | — |
| 2 | Claude GitHub App | Claude 開分支、PR、看 CI | Geodown 安裝授權 | GitHub → Settings → Applications | — |
| 3 | Railway 帳號 | 託管 project `bodhi-alliance`（5 個服務）；所有正式環境祕密都在這裡 | Geodown | Railway 登入；Claude 經 Railway connector 存取 | 各服務 Variables |
| 4 | Wix 帳號 | `sunnylife.world` 網域註冊與 DNS；Wix 轉址根網域 | Geodown／世界佛教教育協會 | Wix 登入 → Domains → Manage DNS Records | — |
| 5 | Anthropic API 金鑰 | Sunny 呼叫 Claude API | Geodown（帳單） | Railway `guide` 變數 | `ANTHROPIC_API_KEY` |
| 6 | Azure AI Speech（尚未申請） | Sunny 雲端語音 | 待定 | 將放在 Railway `guide` 變數 | `AZURE_SPEECH_KEY`、`AZURE_SPEECH_REGION` |
| 7 | 菩提幣營運地址（兼金庫）`0xA32c80dB53945A42F8a34005aD84681c436b875a` | 部署合約、持有 5 億枚、發入會贈幣、付 gas、合約 owner | 平台（Geodown 保管） | 私鑰只在 Railway `api` 變數；Geodown 須另外備份到密碼管理器 | `BODHI_CHAIN_OPERATOR_KEY` |
| 8 | 會員地址推導種子 | 推導每位會員的鏈上地址與私鑰 | 平台（Geodown 保管） | Railway `api` 變數；須另外備份 | `BODHI_CHAIN_MEMBER_SEED` |
| 9 | Polygon（POL）入金來源 | 替營運地址補 gas | Geodown | Geodown 自己的交易所或錢包 | — |
| 10 | 後台第一位超級管理員 | admin.sunnylife.world 登入、建立其他帳號 | Geodown | 資料庫 `app_user`（bcrypt）；email 在 Railway `api` 變數；密碼只有 Geodown 知道 | `BOOTSTRAP_ADMIN_EMAIL`、`BOOTSTRAP_ADMIN_PASSWORD`（建好後應清空） |
| 11 | 其他後台帳號（知識管理員、中心管理員） | 後台分工 | 由超級管理員在後台「帳號」建立 | 資料庫 | — |
| 12 | Railway Postgres 連線 | api、guide 連資料庫 | Railway 自動產生 | Railway `Postgres` 服務，其他服務以引用取得 | `DATABASE_URL` |
| 13 | Skybox AI（Blockade Labs）付費帳號 | 生成 Sunny 場景全景圖 | Geodown | Claude 雲端環境的環境變數（只供開發用，不在正式環境） | `BLOCKADE_API_KEY` |
| 14 | Meshy 帳號 | 生成梅花鹿、蓮花 3D 模型 | Geodown | Claude 雲端環境的環境變數（只供開發用） | `meshy_API_KEY` |
| 15 | claude.ai project「Bodhi coin」 | Claude 協作、專案檔案、connector（GitHub、Railway） | Geodown | claude.ai | — |
| 16 | 本機原型 bodhi-guide | Sunny 原型（localhost:8080） | Geodown 的電腦 | 本機 `.env` | `ANTHROPIC_API_KEY`（本機） |

不需要帳號的外部資源：PolygonScan（公開查詢）、PublicNode／dRPC 公開 RPC 節點、Sketchfab 菩提樹模型（CC BY 4.0，免登入下載）。

## 9.2 祕密清單與備份責任

| 祕密 | 位置 | 遺失後果 | 輪替方式 |
| --- | --- | --- | --- |
| `BODHI_CHAIN_OPERATOR_KEY` | Railway `api` | 金庫 5 億枚與合約管理權無法動用 | 無法輪替金庫本身；只能部署新合約（新總量）或先把幣轉到新地址並 `setInstitutional` |
| `BODHI_CHAIN_MEMBER_SEED` | Railway `api` | 會員地址仍在、幣仍在，但平台無法再替會員簽名 | 換種子＝所有會員換新地址，舊地址的幣無法搬移 |
| `ANTHROPIC_API_KEY` | Railway `guide` | Sunny 對話停止 | 在 Anthropic Console 重新產生後更新變數並重新部署 guide |
| `AZURE_SPEECH_KEY` | Railway `guide`（尚未設定） | 退回瀏覽器語音 | Azure 入口網站重新產生 |
| 超級管理員密碼 | 資料庫雜湊 | 無法登入後台 | 另一位超級管理員重建帳號；或暫時設定新的 `BOOTSTRAP_ADMIN_EMAIL`＋密碼重新部署 api 建立新帳號 |
| `DATABASE_URL` | Railway 自動管理 | — | Railway 重設 Postgres 憑證 |

**最優先待辦**：確認營運私鑰與會員種子已由 Geodown 從 Railway 複製到密碼管理器（2026-10-08 主網切換時已提醒）。測試鏈（Amoy）時期的金鑰已被覆蓋，不再存在。

## 9.3 存取原則

1. 祕密只放 Railway 服務變數，不寫進 repo、對話、文件或 Claude 記憶；repo 只有 `.env.example` 假值。
2. Claude 可經 Railway connector 讀寫變數與看日誌，但不能登入後台（不知道管理員密碼），雲端環境也連不到 `*.up.railway.app`。
3. 真實金錢的動作（替營運地址入 POL）只由 Geodown 執行。
4. 新人加入：GitHub 加協作者、Railway 加成員、後台建帳號並給最小角色。
5. 有人離開：撤 GitHub／Railway 權限、後台刪帳號（會一併登出）；若他看過鏈上祕密，評估是否搬移金庫資產。
