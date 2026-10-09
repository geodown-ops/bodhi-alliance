# 5. HTTP 介面

兩個服務都回 JSON，錯誤格式為 `{"error": "中文訊息"}`。需要登入的路由用 `Authorization: Bearer <token>`。兩個服務都有 `GET /healthz`（Railway 健康檢查）。

圖例 — 權限：**公開**、**會員**（任何登入者）、**中心**（超級管理員或中心管理員）、**超管**（超級管理員）、**知識**（知識管理員或超級管理員）。

## 5.1 核心 API（`api` 服務，前綴 `/api`）

### 登入與帳號 `internal/auth`

| 方法 | 路徑 | 權限 | 說明 |
| --- | --- | --- | --- |
| POST | `/auth/login` | 公開（限流 10/30 秒） | email＋密碼，回 token 與使用者 |
| POST | `/auth/logout` | 會員 | 刪除目前 session |
| GET | `/auth/me` | 會員 | 目前使用者與角色 |
| GET | `/auth/users` | 超管 | 後台帳號列表 |
| POST | `/auth/users` | 超管 | 新增後台帳號與角色 |
| DELETE | `/auth/users/:id` | 超管 | 刪除帳號並登出其所有 session |

### 會員 `internal/members`

| 方法 | 路徑 | 權限 | 說明 |
| --- | --- | --- | --- |
| GET | `/centers` | 公開 | 中心清單（註冊選所屬中心） |
| POST | `/volunteers` | 公開（限流 5/分） | 加入會員：email、password、legal_name 必填；display_name、phone、line_id、home_center_id、wants_coach、in_groups、in_association 選填 |
| GET | `/me/volunteer` | 會員 | 我的會員資料與所屬小組 |
| PUT | `/me/volunteer` | 會員 | 修改個人資料 |
| PUT | `/me/memberships` | 會員 | 切換覺行小組／協會身分（至少保留一個） |
| POST、DELETE | `/me/groups/:id` | 會員 | 加入／退出固定小組 |
| GET | `/admin/volunteers` | 中心 | 會員名冊（中心管理員只看自己中心） |
| PATCH | `/admin/volunteers/:id` | 中心 | 核可／退回；教練、覺行小組長、決策小組勾選 |
| GET | `/admin/groups/:id/members` | 中心 | 小組成員 |
| PUT、DELETE | `/admin/groups/:id/members/:vid` | 中心 | 設成員角色（組長）／移除 |

### 覺行小組、報名與登記 `internal/apply`

| 方法 | 路徑 | 權限 | 說明 |
| --- | --- | --- | --- |
| GET | `/groups` | 公開 | 公開的小組 |
| POST | `/group-applications` | 公開（限流 5/分） | 小組報名 |
| POST | `/merchant-applications` | 公開（限流 5/分） | 共好企業登記 |
| GET、POST、PUT、DELETE | `/admin/groups[/:id]` | 超管 | 小組維護 |
| GET、PATCH | `/admin/group-applications[/:id]` | 超管 | 報名審核 |
| GET、PATCH | `/admin/merchant-applications[/:id]` | 超管 | 登記審核 |

### 組織 `internal/org`

| 方法 | 路徑 | 權限 | 說明 |
| --- | --- | --- | --- |
| GET | `/venues` | 公開 | 活動場域（覺行小組頁） |
| GET、POST、PUT、DELETE | `/admin/centers[/:id]` | 超管 | 中心（還有場域或共好企業時不能刪） |
| GET、POST、PUT、DELETE | `/admin/venues[/:id]` | 超管 | 場域 |
| GET、POST、PUT、DELETE | `/admin/merchants[/:id]` | 超管 | 共好企業 |
| POST | `/admin/merchant-applications/:id/merchant` | 超管 | 登記一鍵轉成共好企業 |

### 共修活動與菩提幣 `internal/events`

| 方法 | 路徑 | 權限 | 說明 |
| --- | --- | --- | --- |
| GET | `/events` | 公開 | 近期開放的活動 |
| GET | `/me/events` | 會員 | 我發起／參加／協助的活動與送審狀態 |
| POST | `/me/events` | 會員 | 發起活動：title、is_online、location、starts_at、ends_at、capacity（≥3）、description、venue_id |
| POST、DELETE | `/me/events/:id/join` | 會員 | 參加（role `participant` 或 `helper`）／退出 |
| POST | `/me/events/:id/cancel` | 會員 | 發起人取消（未送審） |
| POST | `/me/events/:id/claim` | 會員 | 活動結束後送審：attendance（≥3）、report |
| GET | `/me/wallet` | 會員 | 餘額與流水 |
| GET | `/admin/claims` | 超管 | 送審清單，附受領人與建議金額 |
| PATCH | `/admin/claims/:id` | 超管 | status `approved/rejected`、review_note、amounts（每位受領人金額）；核准即寫入帳本 |

### 世界佛教教育協會 `internal/association`

| 方法 | 路徑 | 權限 | 說明 |
| --- | --- | --- | --- |
| GET | `/me/association` | 會員（協會會員才有內容） | 通知、行事曆、會刊 |
| GET、POST、PUT、DELETE | `/admin/association/notices[/:id]` | 超管 | 協會通知 |
| GET、POST、PUT、DELETE | `/admin/association/events[/:id]` | 超管 | 會員行事曆 |
| GET、POST、DELETE | `/admin/association/tags[/:id]` | 超管 | 行事曆標籤 |
| GET、POST、PUT、DELETE | `/admin/association/issues[/:id]` | 超管 | 協會會刊 |
| GET | `/admin/association/members` | 超管 | 協會會員名單 |

### 區塊鏈 `internal/chain`

| 方法 | 路徑 | 權限 | 說明 |
| --- | --- | --- | --- |
| GET | `/chain` | 公開 | `enabled`、chain_id、network、explorer、合約地址與連結（菩提幣介紹頁） |
| GET | `/me/chain` | 會員 | 我的鏈上地址、餘額、入會贈幣狀態與交易連結 |

## 5.2 Sunny 服務（`guide` 服務，前綴 `/guide`）

| 方法 | 路徑 | 權限 | 說明 |
| --- | --- | --- | --- |
| GET | `/info` | 公開 | name、available（有金鑰、已開啟、未超預算）、tts（是否有雲端語音） |
| GET | `/scripts` | 公開 | 已上架的共修腳本 |
| POST | `/chat` | 公開（限流 20/30 秒） | 見下 |
| POST | `/tts` | 公開（限流 40/3 秒） | `{"text"}`（最多 300 字）→ `audio/mpeg`（Azure 24kHz mp3）；未設定雲端語音時回 404，前端改用瀏覽器語音 |
| GET、POST | `/admin/documents` | 知識 | 列表／上傳（txt、md、docx、pdf，上限 20MB） |
| GET、PUT | `/admin/documents/:id` | 知識 | 讀取／修改（存成新版本） |
| GET | `/admin/documents/:id/versions` | 知識 | 版本歷史 |
| POST | `/admin/documents/:id/status` | 知識 | `draft/published/archived` |
| POST | `/admin/try` | 知識 | 試問，可帶 include_document_id 試未上架的文件 |
| GET、PUT | `/admin/persona` | 知識 | 角色設定（name ≤50 字、prompt ≤20,000 字） |
| GET | `/admin/usage` | 知識 | 本月用量、估計費用、預算 |
| PUT | `/admin/settings` | 知識 | monthly_budget_usd、enabled |

### `POST /guide/chat`

請求：

```json
{ "messages": [{ "role": "user", "content": "什麼是覺行小組？" }], "mode": "chat", "script_id": "" }
```

- `mode`：`chat`（問答）或 `practice`（依 `script_id` 的共修腳本帶領）。
- 歷史最多 20 則，每則最多 2,000 字。
- `Accept: text/event-stream` 時串流：多個 `event: delta`（`{"text": "…"}`），最後 `event: done`（`{"reply": "全文", "refused": false}`），失敗時 `event: error`。否則回一般 JSON `{"reply": …}`。
