# Sunny life 菩提幣系統技術規格書

版本 1.0 · 2026-10-09 · 依 `main`（PR #59 合併後）實際上線的設計撰寫

本規格書取代 2026-10-02～10-03 的草案文件（舊的 `docs/architecture.md`、`docs/deploy-railway.md` 與 Claude Doc〈菩提幣官網實作計畫〉）。內容以程式碼為準；程式與本文件不一致時，以程式碼為準並回頭修正本文件。

## 文件目錄

| 章 | 檔案 | 內容 |
| --- | --- | --- |
| 1 | [01-overview.md](01-overview.md) | 系統目的、名詞、功能地圖、角色、目前範圍 |
| 2 | [02-architecture.md](02-architecture.md) | 全系統架構、元件、技術棧、repo 結構、安全設計 |
| 3 | [03-frontend.md](03-frontend.md) | 官網與管理後台：路由、頁面、Sunny 3D 場景、語音、行事曆 |
| 4 | [04-data-model.md](04-data-model.md) | PostgreSQL 資料模型（遷移 001–014） |
| 5 | [05-api.md](05-api.md) | 核心 API 與 Sunny 服務的 HTTP 介面 |
| 6 | [06-blockchain.md](06-blockchain.md) | 菩提幣合約與 Polygon 上鏈 |
| 7 | [07-sunny-ai.md](07-sunny-ai.md) | 線上覺行小組組長 Sunny：對話、知識庫、語音 |
| 8 | [08-deployment.md](08-deployment.md) | 雲端部署設計（Railway）、網域、CI、維運 |
| 9 | [09-accounts.md](09-accounts.md) | 關鍵帳號、金鑰存放位置與存取權 |
| 10 | [10-roadmap.md](10-roadmap.md) | 尚未實作的部分與待決事項 |

## 一句話

禪修志工與減壓教練協助「覺行小組」共修，活動在官網登錄、經超級管理員審核後得到菩提幣；系統會員入會時另由金庫在 Polygon 主網上轉 1 枚菩提幣到平台保管的個人地址。官網同時介紹世界佛教教育協會、正念減壓與彌勒心流，並由 3D 場景中的「線上覺行小組組長 Sunny」以 AI 對話與語音回答問題、帶領共修。

## 線上服務

| 用途 | 網址 |
| --- | --- |
| 官網 | https://www.sunnylife.world（`sunnylife.world` 由 Wix 轉址到 www） |
| 管理後台 | https://admin.sunnylife.world |
| 核心 API | `api` 服務的 Railway 網域（`api-production-9e87.up.railway.app`） |
| Sunny 服務 | `guide` 服務的 Railway 網域（`guide-production-42bb.up.railway.app`） |
| 原始碼 | https://github.com/geodown-ops/bodhi-alliance（`main` 自動部署） |

## 安全原則

本規格書與整個 repo **不含任何金鑰、密碼、私鑰或助記詞的實際值**。第 9 章只記錄帳號名稱、用途、擁有者、存放位置與環境變數名稱。

## 與舊文件的關係

- `docs/data-model.md`：SPEC v2.0 的**目標**資料模型（決議、核發名單、券、帳本分錄等），屬於尚未實作的 M3–M5，保留作設計依據；目前實際的表見第 4 章。
- `docs/architecture.md`、`docs/deploy-railway.md`：已改為指向本規格書。
