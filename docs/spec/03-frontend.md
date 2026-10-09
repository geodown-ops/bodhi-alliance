# 3. 前台與管理後台

兩個前端都是 Vue 3 + Quasar 單頁應用，以 Vite 建置成靜態檔，由 Caddy 提供（第 8 章）。api 與 guide 的網址在建置時由 `VITE_API_BASE`、`VITE_GUIDE_BASE` 寫入。

## 3.1 官網 `apps/web`

### 路由

| 路徑 | 頁面 | 說明 |
| --- | --- | --- |
| `/` | `HomePage.vue` | 第一屏是 Sunny 場景＋單一輸入框（`GuideStage` compact），下方是宗旨（正念減壓、彌勒心流兩段摘要）與「加入共好企業」 |
| `/coin` | `CoinPage.vue` | 菩提幣介紹、梯級表、「在區塊鏈上查詢菩提幣」（呼叫 `GET /api/chain`） |
| `/mindfulness` | `MindfulnessPage.vue` | 正念減壓八週課程總覽 |
| `/mindfulness/lesson/:n` | `MindfulnessLessonPage.vue` | 各週講義（`/mindfulness/week/:n` 轉址至此） |
| `/mindfulness/sitting` | `MindfulnessSittingPage.vue` | 上座與下座 |
| `/groups` | `GroupsPage.vue` | 覺行小組定義、活動場域、近期活動，各場域可報名或發起 |
| `/partners` | `PartnersPage.vue` | 共好企業招募與登記表 |
| `/wallet` | `WalletPage.vue` | 我的錢包說明，登入後導向個人頁 |
| `/maitreya` | `MaitreyaPage.vue` | 彌勒心流與唯識三轉 |
| `/association` | `AssociationPage.vue` | 世界佛教教育協會介紹（`components/association/Assoc*.vue`），結尾「加入會員」 |
| `/guide` | `GuidePage.vue` | 線上問答：全畫面 Sunny 場景、對話、帶我共修 |
| `/join` | `JoinPage.vue` | 加入會員（實名、email、密碼必填；暱稱、電話、LINE 選填；選覺行小組／協會） |
| `/login` | `LoginPage.vue` | 會員登入 |
| `/me` | `MePage.vue` | 個人頁（需登入）：基本資料、身分切換、菩提幣錢包、鏈上菩提幣、我的活動、發起共修、送審、協會通知／我的行事曆／會刊 |
| `/privacy` | `PrivacyPage.vue` | 隱私權保護聲明 |
| `/committee` | `CommitteePage.vue` | 主辦審核小組（不在選單） |

### 視覺規範

- 站名「Sunny life」：Whisper 字型（OFL）外框 SVG，橫向拉長 1.4 倍、筆畫加粗，高 20px，與菩提葉 logo 同高（`src/wordmark.ts`、`public/images/sunny-life-wordmark.svg`）。
- 全站字型：黑體（`--sans`：Noto Sans TC／PingFang TC／微軟正黑體）。
- 協會頁色票：粉、綠、米、logo 紅 `#ba5854`，頁底淺米 `#f9f3e4`，加白色與主文字深咖啡。
- 頁尾：「一即一切，一切即一」＋隱私權保護聲明連結。
- 手機版：抽屜式選單；表單卡片用 `.card-form`；密碼欄有顯示／隱藏的眼睛圖示。

### 會員狀態

`src/account.ts` 把登入 token 存在 `localStorage`（與後台分開）。需要登入的路由（`meta.signedIn`）未登入時導向 `/login?next=…`。頁首「登入」是外框文字按鈕，登入後顯示「個人頁」。

## 3.2 Sunny 3D 場景

| 檔案 | 內容 |
| --- | --- |
| `components/GuideStage.vue` | 場景＋對話區，首頁（compact）與 `/guide` 共用；對話存 `sessionStorage`，靜音設定存 `localStorage` 的 `bodhi.voice` |
| `components/GuideScene.vue` | 載入 3D 角色與場景、進度條 |
| `guide3d/avatar.js` | VRM 角色 `public/models/bodhi.vrm`：站在水中冥想、雙手捧熱茶；被提問時睜眼、鏡頭推近臉部，回答時依語音對嘴 |
| `guide3d/scene-panorama.js` | 全景圖場景共用做法（天空盒圖＋河流遮罩＋水面反射） |
| `guide3d/scene-meadow.js` | 晨霧草原：菩提樹隨風擺動與落葉、3 隻梅花鹿從河中走到草地、河中央 3 朵蓮花 |
| `guide3d/scene-pines.js` | 松林雪山 |
| `guide3d/scene-dusk.js` | 黃昏湖景與蘆葦（程序化生成） |
| `guide3d/scene.js` | 原本的藍天湖景（`?scene=day`） |
| `sceneTime.ts` | 依觀看者裝置時間選場景 |

**場景時刻表**（頁面打開時決定一次；網址加 `?scene=meadow|pines|dusk|day` 可固定）：

| 頁面 | 時段內 | 時段外 |
| --- | --- | --- |
| 首頁 | 06:00–22:00 晨霧草原（菩提樹） | 黃昏蘆葦 |
| 線上問答 | 08:00–16:00 松林雪山 | 晨霧草原（菩提樹） |

場景支援滑鼠滾輪縮放。

**素材與授權**

| 素材 | 來源 | 授權 |
| --- | --- | --- |
| 全景天空圖 `public/scenes/*.jpg`（含 8K 版） | Skybox AI，Geodown 的付費帳號生成 | 付費帳號生成物；Skybox 範例包是 CC BY-NC，**不得使用** |
| 菩提樹 `public/models/bodhi-tree` | Sketchfab「Bodhi Tree」by Ashim Shakya | CC BY 4.0，場景左下角標示作者 |
| 梅花鹿、蓮花 `public/models/deer`、`lotus` | Geodown 的 Meshy 帳號生成 | 帳號生成物；處理過程見各目錄 `SOURCE.txt` |
| Sunny 角色 `bodhi.vrm` | 原 bodhi-guide 本機原型 | 專案自有 |

## 3.3 語音

`src/voice.ts`：回答串流進來時，每湊滿一句（遇到。！？；或超過 60 字在逗號處斷開）就先唸。

1. 先呼叫 `POST {guide}/guide/tts` 取得 Azure 神經語音的音檔（伺服器有設定 `AZURE_SPEECH_KEY` 時）。
2. 取不到時改用瀏覽器 `speechSynthesis`。
3. 唸之前移除網址、Markdown 符號與表情符號；讀音修正（例如「覺行」唸 jué xíng）只改唸的文字，畫面文字不變。

右上角有靜音按鈕。

## 3.4 行事曆

`apps/shared/calendar`（`CalendarBoard.vue`、`CalGrid.vue`、`DateTimeField.vue`、`cal.ts`）移植自 dengo 專案的行程管理：清單／日／週三種檢視，手機一次看 4 天、左右拖曳換週，行程類型決定顏色（一般、培訓、會議、活動、典禮），可掛一個活動標籤篩選。官網個人頁「我的行事曆」唯讀，後台「世界佛教教育協會」頁可編輯。

## 3.5 管理後台 `apps/admin`

登入後依角色決定可見選單；沒有任何可用頁面的帳號會被擋在登入頁（`?denied=1`）。

| 路徑 | 頁面 | 角色 | 功能 |
| --- | --- | --- | --- |
| `/applications` | 報名與登記 | 超級管理員 | 覺行小組報名、共好企業登記的審核；登記一鍵轉成共好企業 |
| `/venues` | 場域管理 | 超級管理員 | 中心與活動場域 |
| `/merchants` | 共好企業管理 | 超級管理員 | 共好企業資料與狀態 |
| `/groups` | 覺行小組 | 超級管理員 | 小組資料、成員、各小組組長 |
| `/claims` | 菩提幣審核 | 超級管理員 | 活動送審核准／退回，可調整每人金額 |
| `/association` | 世界佛教教育協會 | 超級管理員 | 協會通知、會員行事曆（含標籤）、會刊、協會會員名單 |
| `/volunteers` | 會員名冊 | 超級管理員、中心管理員 | 會員核可；教練、覺行小組長、菩提幣決策小組的勾選；每欄可點標題排序 |
| `/knowledge`、`/knowledge/:id` | 知識庫 | 知識管理員 | 上傳、編輯、版本、上下架、試問 |
| `/guide-settings` | Sunny訓練設定 | 知識管理員 | 角色設定（名字、提示）、開關對話、每月預算、本月用量 |
| `/users` | 帳號 | 超級管理員 | 新增／刪除後台帳號與角色 |

## 3.6 測試

`npm test` 執行 vitest（`apps/web/src/__tests__`），`npm run build` 含 `vue-tsc` 型別檢查；CI 兩者都跑。
