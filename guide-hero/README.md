# guide-hero（覺行小組線上組長首頁開頭）

從 bodhi-guide 複製來的首頁 #guide 區塊，供官網重用。

- `index.html`：單獨可開的預覽頁（hero HTML + CSS + importmap + 字型，`{{LEAF}}` 已替換）。需用 http 伺服器開啟（ES module）。
- `guide-hero.html` / `guide-hero.css`：原始片段（`{{LEAF}}` 為葉子 SVG path，由 bodhi-guide 的 tools/build-site.mjs 替換）。
- `guide.js`：單一輸入框問答＋字幕，POST `${BODHI_GUIDE_API}/api/chat`（SSE：event delta / error）。
- `avatar.js`：VRM 角色（冥想、對嘴、走動、鏡頭推近）。`scene.js`：程序化黃昏湖景（天空、湖面、蘆葦、麥田），無外部貼圖。
- `models/bodhi.vrm`：角色模型。
- `persona-knowledge/`：角色設定與知識庫文字（不含任何金鑰）。
- 外部依賴：three@0.186.1、@pixiv/three-vrm@3.5.5（jsDelivr importmap），Google Fonts（Noto Sans TC、Spline Sans Mono 等）。
