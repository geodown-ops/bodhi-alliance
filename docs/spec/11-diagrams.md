# 11. 服務架構圖與流程圖

本章集中放全系統的架構圖與主要流程圖，以 Mermaid 撰寫（GitHub 直接顯示）。同一組圖另有 PNG／SVG 版本放在專案檔案 `spec/diagrams/`。各章另有局部的圖。

| 圖 | 內容 |
| --- | --- |
| 11.1 | 服務架構圖 |
| 11.2 | 雲端部署架構圖 |
| 11.3 | 會員註冊與入會贈幣上鏈 |
| 11.4 | 共修發起到菩提幣審核發放 |
| 11.5 | 活動與送審狀態 |
| 11.6 | Sunny 問答與語音 |
| 11.7 | 知識庫上架 |
| 11.8 | 發布部署流程 |
| 11.9 | 登入與權限檢查 |

## 11.1 服務架構圖

```mermaid
flowchart TB
  subgraph users["使用者端"]
    visitor["訪客／會員<br/>瀏覽器・手機"]
    staff["管理者<br/>超級管理員・中心管理員・知識管理員"]
  end

  subgraph front["前端（靜態站，Caddy）"]
    web["web 官網<br/>Vue 3 + Quasar<br/>Sunny 3D 場景 three.js / VRM<br/>語音播放 voice.ts"]
    admin["admin 管理後台<br/>Vue 3 + Quasar"]
  end

  subgraph backend["後端（Go + Gin）"]
    api["<b>api 核心服務</b><br/>auth 登入・帳號・角色<br/>members 會員・名冊<br/>events 共修活動・菩提幣審核・錢包<br/>apply 覺行小組・報名・共好企業登記<br/>org 中心・場域・共好企業<br/>association 協會通知・行事曆・會刊<br/>chain 鏈上查詢 API"]
    worker["<b>上鏈 worker</b><br/>（在 api 內，每 20 秒）"]
    guide["<b>guide Sunny 服務</b><br/>chat 對話・帶領共修（SSE 串流）<br/>知識庫・版本・試問<br/>Sunny 訓練設定・月預算<br/>tts 語音代理<br/>seed 內建知識自動上架"]
    api --- worker
  end

  subgraph data["資料"]
    pg[("PostgreSQL<br/>public：會員・活動・帳本・協會・鏈上紀錄<br/>guide：知識・用量・設定")]
  end

  subgraph ext["外部服務"]
    claude["Anthropic Claude API"]
    azure["Azure AI Speech（選用）"]
    rpc["Polygon 公開 RPC 節點"]
    polygon["Polygon 主網<br/>BodhiCoin ERC-20"]
    scan["PolygonScan"]
  end

  visitor --> web
  staff --> admin
  web -- "REST・Bearer token" --> api
  web -- "REST・SSE" --> guide
  admin -- "REST・Bearer token" --> api
  admin -- "REST" --> guide
  api --> pg
  guide --> pg
  worker --> pg
  guide --> claude
  guide --> azure
  worker --> rpc --> polygon
  web -. "交易・合約連結" .-> scan
```

## 11.2 雲端部署架構圖

```mermaid
flowchart LR
  dev["開發者／Claude<br/>功能分支"] --> gh["GitHub<br/>geodown-ops/bodhi-alliance"]
  gh -- "PR" --> ci["GitHub Actions CI<br/>Go: gofmt・vet・test -race<br/>Apps: build・vitest"]
  gh -- "main 有新 commit" --> rw

  subgraph rw["Railway project bodhi-alliance（production）"]
    direction TB
    webs["web<br/>apps/Dockerfile APP=web<br/>node:22 建置 → caddy:2"]
    admins["admin<br/>apps/Dockerfile APP=admin"]
    apis["api<br/>server/Dockerfile SERVICE=api<br/>golang:1.26 → distroless"]
    guides["guide<br/>server/Dockerfile SERVICE=guide"]
    pgs[("Postgres")]
    vars["服務變數（祕密只在這裡）<br/>ANTHROPIC_API_KEY・BODHI_CHAIN_*<br/>BOOTSTRAP_ADMIN_*・ALLOWED_ORIGINS"]
    apis --> pgs
    guides --> pgs
    vars -.-> apis
    vars -.-> guides
  end

  subgraph dns["Wix DNS：sunnylife.world"]
    www["www → CNAME"]
    adm["admin → CNAME"]
    apex["根網域 → Wix 轉址到 www"]
  end

  www --> webs
  adm --> admins
  apex --> www
  webs -- "VITE_API_BASE" --> apis
  webs -- "VITE_GUIDE_BASE" --> guides
  admins --> apis
  admins --> guides
```

## 11.3 會員註冊與入會贈幣上鏈

```mermaid
sequenceDiagram
  autonumber
  actor M as 會員
  participant W as 官網 /join
  participant A as api
  participant DB as PostgreSQL
  participant K as 上鏈 worker（api 內）
  participant R as Polygon RPC
  participant C as BodhiCoin 合約

  M->>W: 填實名、email、密碼，選覺行小組／協會
  W->>A: POST /api/volunteers
  A->>A: 限流 5/分、honeypot、密碼 ≥ 10 字
  A->>DB: 交易：app_user（bcrypt）＋ volunteer
  A-->>W: 建立成功
  W->>A: POST /api/auth/login
  A->>DB: user_session（只存 token 雜湊）
  A-->>W: token（存 localStorage）
  W->>M: 進入個人頁 /me

  loop 每 20 秒
    K->>DB: 找還沒有本鏈地址或贈幣的會員
    K->>K: HMAC-SHA256(會員種子, 會員 id) 推導地址
    K->>DB: member_chain_account、chain_grant = pending
    K->>R: 金庫 transfer(會員地址, 1.00 BODHI)
    R->>C: 執行轉帳
    K->>DB: chain_grant = sent（tx_hash）
    R-->>K: 交易收據
    K->>DB: chain_grant = confirmed（block_number）
  end

  M->>W: 個人頁「鏈上菩提幣」
  W->>A: GET /api/me/chain
  A-->>W: 地址、餘額、交易連結（PolygonScan）
```

## 11.4 共修發起到菩提幣審核發放

```mermaid
flowchart TB
  start(["會員登入"]) --> create["發起共修<br/>POST /api/me/events<br/>日期・開始・結束・人數 ≥ 3・場域"]
  create --> open["活動 open<br/>公開在覺行小組頁"]
  open --> join["其他會員報名<br/>參加者 participant／協辦 helper<br/>POST /me/events/:id/join"]
  join --> held{"活動結束了？"}
  open --> cancel["發起人取消（未送審前）"] --> cancelled(["cancelled"])
  held -- "否" --> join
  held -- "是" --> claim["發起人送審<br/>出席人數 ≥ 3、活動回報<br/>POST /me/events/:id/claim"]
  claim --> submitted["coin_claim = submitted"]
  submitted --> review["超級管理員：後台「菩提幣審核」<br/>系統依時數給建議金額<br/>未滿 4h 300/h、4h 1,500、8h 3,000"]
  review --> decide{"核准？"}
  decide -- "退回（附說明）" --> rejected(["rejected<br/>會員在個人頁看到原因"])
  decide -- "核准（可調整每人金額）" --> tx["同一筆資料庫交易<br/>claim_recipient：發起人＋協辦<br/>coin_ledger：kind = issue"]
  tx --> wallet(["個人頁錢包餘額增加<br/>GET /api/me/wallet"])
  wallet -.-> future["之後（第 10 章）：<br/>決策小組梯級、分級審核、上鏈"]
```

## 11.5 活動與送審狀態

```mermaid
stateDiagram-v2
  [*] --> open: 發起
  open --> cancelled: 發起人取消
  open --> ended: 結束時間已過
  ended --> submitted: 發起人送審
  submitted --> approved: 超級管理員核准，寫入帳本
  submitted --> rejected: 超級管理員退回
  cancelled --> [*]
  approved --> [*]
  rejected --> [*]
  note right of ended
    ended 不是資料庫欄位，
    由 ends_at 與現在時間判斷
  end note
```

## 11.6 Sunny 問答與語音

```mermaid
sequenceDiagram
  autonumber
  actor U as 使用者
  participant S as 官網 GuideStage／3D 場景
  participant G as guide 服務
  participant DB as PostgreSQL（guide schema）
  participant L as Claude API
  participant T as Azure Speech
  participant B as 瀏覽器語音

  S->>G: GET /guide/info
  G-->>S: 名字、是否可用、是否有雲端語音
  U->>S: 輸入問題
  S->>S: 角色睜眼、鏡頭推近
  S->>G: POST /guide/chat（Accept: text/event-stream）
  G->>G: 限流 20/30 秒；檢查開關與本月預算
  G->>DB: 讀角色設定與已上架知識段落
  G->>G: 知識 ≤ 15 萬字整份放入，否則挑 12 段<br/>（角色＋回答方式＋知識，知識加快取）
  G->>L: Messages API 串流
  loop 逐段回傳
    L-->>G: 文字片段
    G-->>S: SSE event: delta
    S->>S: 湊滿一句
    alt 有雲端語音
      S->>G: POST /guide/tts（一句）
      G->>T: SSML（zh-TW-HsiaoChenNeural）
      T-->>G: mp3
      G-->>S: audio/mpeg
    else 沒有或失敗
      S->>B: speechSynthesis
    end
    S->>S: 播放並對嘴
  end
  G-->>S: SSE event: done（全文）
  G->>DB: guide.usage 記錄 token 用量
  S->>S: 回到冥想、鏡頭拉遠
```

## 11.7 知識庫上架

```mermaid
flowchart LR
  subgraph sources["來源"]
    upload["知識管理員上傳<br/>txt・md・docx・pdf（≤ 20MB）"]
    seedf["repo 內建<br/>server/internal/guide/seed/*.md"]
  end
  upload --> extract["轉成文字"] --> draft["草稿 draft<br/>document_version v1"]
  draft --> edit["後台修改 → 新版本"]
  edit --> draft
  draft --> try["試問：看 Sunny 怎麼回答"]
  try --> publish["上架 published<br/>切成約 600 字段落 guide.chunk"]
  seedf -- "guide 啟動時，第一次看到就上架；<br/>檔案改了且後台沒人改過才更新" --> publish
  publish --> live(["Sunny 回答時使用"])
  publish --> archive["下架 draft／封存 archived"] --> gone(["不再使用，版本保留"])
```

## 11.8 發布部署流程

```mermaid
flowchart TB
  ask(["Geodown 在 thread 提出需求"]) --> branch["Claude 開功能分支、改程式"]
  branch --> local["本機檢查<br/>gofmt・go vet・go test・npm build・npm test"]
  local --> pr["推上 GitHub，開 draft PR"]
  pr --> ci{"CI 全綠？"}
  ci -- "否" --> fix["修正後再推"] --> ci
  ci -- "是" --> show["回報 thread（截圖・說明）"]
  show --> ok{"Geodown 打「merge」？"}
  ok -- "要修改" --> branch
  ok -- "是" --> merge["合併到 main"]
  merge --> watch{"Railway 依 Watch Paths<br/>判斷受影響的服務"}
  watch -- "server/**" --> goimg["重建 api、guide 映像"]
  watch -- "apps/**" --> feimg["重建 web、admin 映像<br/>VITE_* 寫入網頁"]
  goimg --> migrate["啟動時自動遷移<br/>advisory lock，只跑新檔"]
  migrate --> health["/healthz 健康檢查通過才切換"]
  feimg --> live2(["新版上線"])
  health --> live2
  live2 --> verify["Claude 查 Railway 日誌確認"]
```

## 11.9 登入與權限檢查

```mermaid
sequenceDiagram
  autonumber
  participant F as 前端（官網或後台）
  participant H as api／guide 路由
  participant Mw as 權限中介層
  participant DB as PostgreSQL

  F->>H: 請求＋Authorization: Bearer token
  H->>Mw: Require(角色, 範圍)／RequireUser／RequireCenterStaff
  Mw->>DB: SHA-256(token) 查 user_session（未過期）
  DB-->>Mw: 使用者與 role_assignment
  alt 沒有 token 或已過期
    Mw-->>F: 401 請先登入
  else 角色不符
    Mw-->>F: 403 沒有權限
  else 通過（超級管理員視為擁有所有角色）
    Mw->>H: 繼續處理
    H-->>F: 200 JSON
  end
```
