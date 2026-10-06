<script setup lang="ts">
// 世界佛教教育協會介紹：完整呈現協會簡報《世界佛教教育協會介紹》（36 頁）的內容，
// 簡介與任務的文字以協會提供的《協會簡介》為準。簡報裡的圖表與插圖全部改用網頁原生的 SVG／HTML 重新繪製，
// 不使用 AI 產生的圖片；簡報中的系統截圖與照片改成示意圖。

// ---------- 協會總則 ----------
const pillars = [
  { name: '慈悲關懷', icon: 'volunteer_activism', text: '人文關懷、社會救助、與護生及生態環保。' },
  { name: '智慧無礙', icon: 'auto_stories', text: '雲端教學、佛學研究、生命教育、終身學習。' },
  { name: '善巧方便', icon: 'diversity_3', text: '佛學資訊、音像紀錄，養生及身心健康、食品安全，以及社會型企業。' },
]

const tasks = [
  { icon: 'volunteer_activism', text: '推動慈悲關懷工作，包括人文關懷、社會救助、與護生及生態環保等。' },
  { icon: 'auto_stories', text: '推動智慧無礙工作，包括雲端教學、佛學研究、生命教育、終身學習等。' },
  { icon: 'diversity_3', text: '推動善巧方便工作，包括佛學資訊、音像紀錄，養生及身心健康、食品安全，以及社會型企業等。' },
  { icon: 'hub', text: '結合科技、醫藥、宗教、政府相關單位及民間社團組織等，推動終身學習及社會型企業策略聯盟。' },
  { icon: 'groups', text: '舉辦研討會、座談會、職業訓練、組織學習、關懷訪視等各項相關活動，以促進佛教普及弘化及世界永續發展之目的。' },
]

const aims = [
  { key: 'exchange', html: '推廣並促進國內外<b>產官學研</b>界交流合作，讓<b>佛教</b>能普及弘化及永續發展' },
  { key: 'healing', html: '透過佛教<b>八正道</b>傳統精神，讓現代人面臨的<b>身心亞健康</b>得以舒緩修復' },
  { key: 'enterprise', html: '藉由傳統<b>佛法教育</b>，讓企業走向善管理的善企業，進而善盡<b>社會責任</b>' },
]

// ---------- 協會起源：世代變遷 ----------
const generations = [
  { name: '嬰兒潮世代', years: '1946–1965', people: '11.7 億', pct: 15 },
  { name: 'X 世代', years: '1966–1980', people: '14.2 億', pct: 18 },
  { name: 'Y 世代', years: '1981–1994', people: '17.4 億', pct: 22 },
  { name: 'Z 世代', years: '1995–2009', people: '18.5 億', pct: 24 },
]

const faith = [
  { value: '三分之一', icon: 'psychology_alt', text: 'Z 世代感覺大部分的時間都處於壓力和孤獨狀況' },
  { value: '四分之三', icon: 'person_off', text: 'Z 世代屬於無宗教信仰傾向，但只有 16% 面對無助時，會尋求信仰的求助' },
  { value: '52%', icon: 'domain_disabled', text: 'Z 世代幾乎不信任宗教組織' },
]

// ---------- 協會起源：全人健康 ----------
const stress = [
  { value: 68, text: 'Z 世代認為壓力是他們幸福的障礙，相對於 2021 年 65% 的比例成長了 3%' },
  { value: 48, text: 'Z 世代認為他們自己是無法或沒有能力有效管理和解決他們自己的壓力' },
  { value: 53, text: 'Z 世代在疫情大流行期間，無法有效處理壓力感到更加孤獨、寂寞與空虛' },
]

const symptoms = [
  { name: '嚴重症狀', note: '積極醫療', range: '11–15%', inner: false },
  { name: '重感症狀', note: '需接受建議', range: '32–38%', inner: true },
  { name: '中度症狀', note: '積極找尋治療', range: '19–26%', inner: true },
  { name: '輕度症狀', note: '以前有症狀', range: '18–22%', inner: true },
  { name: '輕微症狀', note: '以前無症狀', range: '6–8%', inner: false },
]
const bodyCerts = ['健康管理師（中國、國際、進階…）', '保健營養規畫師', '長期照顧健康管理師', '健康餐飲調配師', '體重管理師']
const beautyCerts = ['抗衰老健康管理師', '醫學美容諮詢師', '美容醫學美容師']

// 全人健康六面向：六角形依簡報的位置排列，外圈是六種商數
const hexes = [
  { name: '環境', x: 0, y: -1 },
  { name: '理智', x: 0.87, y: -0.5 },
  { name: '感情', x: 0.87, y: 0.5 },
  { name: '社會', x: 0, y: 1 },
  { name: '身體', x: -0.87, y: 0.5 },
  { name: '精神', x: -0.87, y: -0.5 },
]
const quotients = [
  { name: 'SQ', x: 0.58, y: -1.0 },
  { name: 'VQ', x: 1.18, y: 0 },
  { name: 'EQ', x: 0.58, y: 1.0 },
  { name: 'CQ', x: -0.58, y: 1.0 },
  { name: 'AQ', x: -1.18, y: 0 },
  { name: 'PQ', x: -0.58, y: -1.0 },
]
const hexPoints = (cx: number, cy: number, r: number) =>
  Array.from({ length: 6 }, (_, i) => {
    const a = (Math.PI / 3) * i
    return `${(cx + r * Math.cos(a)).toFixed(1)},${(cy + r * Math.sin(a)).toFixed(1)}`
  }).join(' ')

const wholeHealthFor = [
  '希望改善緊張、焦慮、失眠、長期疲勞或長期疼痛等困擾者',
  '處於高壓競爭環境的上班族、專業經理人、主管或經營者',
  '醫護、心理、社工、教育、人資等相關領域工作者及系所學生',
  '需要在工作挑戰與家庭責任中找到平衡點者',
  '希望增進情緒管理能力與改善人際關係者',
  '希望成長自我並提昇生活品質與幸福感者',
]

// ---------- 協會起源：社會責任 ----------
const csrTotal = [
  { name: '公司治理', score: 66.1 },
  { name: '員工、供應鏈、顧客', score: 85.5 },
  { name: '環境績效', score: 67.4 },
  { name: '社會關懷', score: 76.5 },
]
const csrIndustries = [
  { name: '電信業', icon: 'public', total: 85.8, s: [72.6, 89.8, 89.7, 91.0] },
  { name: '金融保險業', icon: 'shield', total: 73.5, s: [67.3, 86.9, 63.1, 81.3] },
  { name: '電子科技業', icon: 'memory', total: 73.1, s: [65.7, 91.7, 67.6, 71.3] },
  { name: '傳統產業', icon: 'settings', total: 70.5, s: [65.2, 79.8, 67.7, 71.3] },
  { name: '服務業', icon: 'room_service', total: 70.4, s: [62.4, 84.8, 63.2, 75.2] },
]
const issues = [
  ['法令遵循', 64.9], ['經營績效', 62.2], ['內稽內控與風險品質管理（涵蓋產品品質追蹤、危機處理）', 59.5], ['能源管理', 56.8],
  ['公司治理與誠信經營', 51.4], ['職場安全衛生', 51.4], ['員工權益與人權（含友善職場、溝通管道）', 51.4], ['永續發展策略', 48.6],
  ['人才管理與培育（教育訓練）', 45.9], ['廢物、污染管理', 45.9],
] as const
const sdgsTop = [
  ['教育品質', 86.5], ['就業與經濟成長', 83.8], ['責任消費與生產', 83.8], ['氣候行動', 83.8], ['可負擔能源', 78.4],
  ['健康與福祉', 75.7], ['工作、創新和基礎設施', 73.0], ['性別平等', 70.3], ['永續城市', 70.3], ['和平與正義制度', 67.6],
] as const

const SDG: Record<string, [string, string]> = {
  '01': ['消除貧窮', '#e5243b'], '02': ['消除飢餓', '#dda63a'], '03': ['健康與福祉', '#4c9f38'], '04': ['教育品質', '#c5192d'],
  '05': ['性別平等', '#ff3a21'], '06': ['淨水與衛生', '#26bde2'], '07': ['可負擔能源', '#fcc30b'], '08': ['就業與經濟成長', '#a21942'],
  '09': ['工業、創新與基礎建設', '#fd6925'], '11': ['永續城市', '#fd9d24'], '12': ['責任消費與生產', '#bf8b2e'], '13': ['氣候行動', '#3f7e44'],
  '15': ['陸域生態', '#56c02b'], '16': ['和平與正義制度', '#00689d'],
}
const sdgGroups = [
  { name: '環境', icon: 'public', goals: ['06', '11', '12', '13'] },
  { name: '社會', icon: 'groups', goals: ['01', '03', '04', '05', '08', '12', '13', '15'] },
  { name: '治理', icon: 'gavel', goals: ['01', '02', '03', '07', '08', '09', '12', '13', '16'] },
]

const dharmaEdu = [
  ['一個以', '人為中心', '的管理模式'],
  ['一種', '圓融', '的辯證法思想'],
  ['', '諸行無常', '：永恆發展變化規律'],
  ['諸法無我的', '大智慧', '和大境界'],
  ['對', '人性細緻剖析', '能為管理者自我提升'],
  ['對', '僧團的管理', '揭示了企業管理的關鍵'],
]

// ---------- 營運特色 ----------
const archive = [
  { title: '店鋪註冊', icon: 'storefront', tone: '#3d7fd0' },
  { title: '跨境選品', icon: 'travel_explore', tone: '#2f5fbf' },
  { title: '產品編輯', icon: 'edit_note', tone: '#5d8ad6' },
  { title: '產品管理（二次優化）', icon: 'tune', tone: '#8c8f99' },
  { title: '店鋪運營', icon: 'insights', tone: '#e5a24a' },
  { title: '訂單物流', icon: 'local_shipping', tone: '#24307a' },
  { title: '客戶服務', icon: 'support_agent', tone: '#d77aa8' },
  { title: '測試課程', icon: 'quiz', tone: '#3d7fd0' },
]
const archiveCourses = [
  { task: '任務 2：wish 跨境電商實訓課程', items: ['2-1 認識 wish 平台', '2-2 企業賣家註冊流程', '2-3 個人賣家註冊流程', '2-4 後台產品添加流程'] },
  { task: '任務 3：速賣通實訓課程', items: ['3-1 全球速賣通介紹', '3-2 速賣通註冊流程', '3-3 速賣通資費及物流'] },
]
const kmNav = ['總覽', '店鋪註冊', '跨境選品', '產品編輯', '產品管理', '訂單物流', '產品展示', '客戶服務', '我的實訓', '幫助中心']
const kmLessons = [
  'Amazon 亞馬遜圖片要求與處理技巧', 'wish 商品標題編輯技巧', 'wish 平台上那些被用爛了的標籤', '亞馬遜 listing 編寫技巧', '修改或編輯 eBay 促銷刊登活動',
  '商品編輯中的 5 個銷量', '打造完美的 eBay 商品刊登', '跨境電商的文字編輯', '速賣通標題編寫技巧',
]

const addie = [
  { letter: 'A', stage: '分析階段', work: '進行課程需求分析', items: ['進行市場需求分析', '界定主要學習對象', '確定教學／訓練目的', '重組職能內涵'], out: ['先備條件', '課程架構或學習路徑'] },
  { letter: 'D', stage: '設計階段', work: '訂定教學／訓練目標與課程內容', items: ['訂定教學／訓練目標', '規劃課程內容'], out: ['教學／訓練目標', '課程內容（授課大綱）'] },
  { letter: 'D', stage: '發展階段', work: '選擇教學方法，準備教材與教學資源', items: ['選擇教學方法', '撰寫教材與講義', '準備教學資源', '發展評量工具'], out: ['教學方法', '教學資源'] },
  { letter: 'I', stage: '實施階段', work: '辦理課程並進行評量', items: ['辦理課程', '實施評量活動'], out: ['課程辦理相關文件記錄'] },
  { letter: 'E', stage: '評估階段', work: '評量學習成果，監控評估課程成效及回饋修正', items: ['評量學習成果', '收集學習成果證據', '監控評估與回饋修正'], out: ['學習成果評量說明文件', '學習成果證據相關文件', '監控評估機制說明與記錄文件'] },
]
const steps = [
  { name: '問卷調查', text: '填寫數量需達到回收標準', tone: '#2f6b2a' },
  { name: '課程規劃', text: '課程企劃，給予課程規畫建議', tone: '#24406b' },
  { name: '課程製作', text: '老師正式製作課程，考核老師上課教學品質', tone: '#c2462a' },
  { name: '課程說明', text: '透過說明會銷講方式，學生互動回饋達成交易', tone: '#6f7a1f' },
  { name: '正式開課', text: '通知購課同學上課，並於課後收集教學回饋', tone: '#8a5a12' },
]
const badgeUses = ['學員／老師未來課程身份認證使用', '與學員活動訊息的溝通管道', '結業或是證書的數位證明', '學員／老師社群互動的應用']

const units = [
  { name: '指導委員會', tone: '#6a3a8c', items: ['提供市場與產業課程需求趨勢', '審核課程教材適用性與適法性', '評核課程師資的教學專業性', '稽核課程上架後的發展績效', '決議課程於市場發展生命週期'] },
  { name: '課程發展中心', tone: '#3c6a24', items: ['負責課程市場需求的調研', '開發課程的教材與教具', '舉辦課程說明或體驗會活動', '負責整體課程招生計畫與執行', '評量課程上課的實際狀況與回饋'] },
  { name: '認證中心', tone: '#1f6db5', items: ['負責課程內容產出的標準化制定、更新及規範實施流程', '依指導委員會的評核結果，給予課程師資證照', '評測後給予學員結業證明或證照'] },
]
const secretariat = ['行政', '總務', '財務', '秘書']
const committees = ['公關委員會', '業務拓展委員會', '募款委員會', '資訊委員會', '寺廟規劃委員會']
const courses = ['課程一', '課程二', '課程三', '課程四', '課程五', '課程六']

const temple = [
  { title: '林間光之樹屋', icon: 'cabin', from: '#2d3a24', to: '#c98a3c' },
  { title: '沉浸式森林光影廳', icon: 'forest', from: '#14301f', to: '#4f9a52' },
  { title: '林中禪修小屋', icon: 'house_siding', from: '#3b3a33', to: '#a4865c' },
  { title: '夜間林間膠囊屋', icon: 'nights_stay', from: '#1d2232', to: '#c48a4a' },
  { title: '蛋形樹屋', icon: 'egg', from: '#2f5a3a', to: '#d7c38f' },
  { title: '林間餐廳', icon: 'restaurant', from: '#24402a', to: '#b8853f' },
  { title: '頌缽聲浴療癒', icon: 'graphic_eq', from: '#5a3c22', to: '#d8b66a' },
  { title: '鑼浴療癒空間', icon: 'surround_sound', from: '#3a3a3a', to: '#d9d4c7' },
  { title: '光影冥想靜坐', icon: 'self_improvement', from: '#0d2a24', to: '#3fbfa0' },
]

const members = [
  { role: '策略夥伴', icon: 'handshake', html: '以策略投資的方式，支援協會在初始階段及未來營運過程中的<b>財務需求</b>。' },
  { role: '資源協助者', icon: 'inventory_2', html: '提供協會<b>多面向的相關必要產業資源（例如：企業培訓需求等）</b>，以完成協會的設定任務。' },
  { role: '推廣協同者', icon: 'campaign', html: '提供協會在<b>會員招募、課程招生與必要輔助活動的執行</b>，並藉由會員力量協助協會推廣。' },
]

const pct = (v: number, max = 100) => `${(v / max) * 100}%`
</script>

<template>
  <q-page class="page assoc">
    <!-- ===== 封面 ===== -->
    <header class="cover">
      <svg class="cover-mark" viewBox="0 0 120 120" aria-hidden="true">
        <circle cx="60" cy="42" r="30" />
        <circle cx="42" cy="74" r="30" />
        <circle cx="78" cy="74" r="30" />
      </svg>
      <div>
        <p class="eyebrow">World Buddhist Education Association</p>
        <h1>世界佛教教育協會</h1>
        <p class="lead q-mb-none">慈悲關懷・智慧無礙・善巧方便</p>
      </div>
    </header>

    <nav class="toc" aria-label="本頁章節">
      <a href="#rules">協會總則</a>
      <a href="#origin">協會起源</a>
      <a href="#operation">營運特色</a>
      <a href="#summary">總結</a>
      <a href="#members">會員招募</a>
    </nav>

    <!-- ===== 協會總則 ===== -->
    <section id="rules" class="chapter">
      <p class="chapter-kicker">協會總則</p>
      <h2 class="chapter-title">簡介・任務・宗旨</h2>
    </section>

    <h3 class="sec">協會簡介</h3>
    <div class="intro-grid">
      <svg class="pillar-ring" viewBox="-110 -110 220 220" role="img" aria-label="三大主軸：慈悲關懷、善巧方便、智慧無礙">
        <path class="ring-a" d="M0,0 L0,-100 A100,100 0 0,1 86.6,50 Z" />
        <path class="ring-b" d="M0,0 L86.6,50 A100,100 0 0,1 -86.6,50 Z" />
        <path class="ring-c" d="M0,0 L-86.6,50 A100,100 0 0,1 0,-100 Z" />
        <circle r="34" class="ring-core" />
        <text x="-44" y="-30" class="ring-label">慈悲關懷</text>
        <text x="44" y="-30" class="ring-label">善巧方便</text>
        <text x="0" y="72" class="ring-label">智慧無礙</text>
        <text x="0" y="6" class="ring-core-label">三主軸</text>
      </svg>
      <ul class="dots">
        <li>本會為依法設立、非以營利為目的之社會團體</li>
        <li>本會成立於民國 111 年 12 月 30 日（送件申請）</li>
        <li>
          本會以（1）慈悲關懷、（2）智慧無礙、（3）善巧方便為主軸，推展世界佛教雲端教學、佛學研究、社會教育、人文關懷、永續發展及社會型企業為宗旨。
        </li>
      </ul>
    </div>
    <div class="grid q-mt-md">
      <div v-for="p in pillars" :key="p.name" class="card pillar">
        <q-icon :name="p.icon" size="28px" class="pillar-icon" />
        <h4>{{ p.name }}</h4>
        <p class="q-mb-none">{{ p.text }}</p>
      </div>
    </div>

    <h3 class="sec">協會任務</h3>
    <p>本會之任務如下，並依相關法令規定推動及執行：</p>
    <ol class="tasks">
      <li v-for="(t, i) in tasks" :key="i">
        <span class="task-no">{{ ['一', '二', '三', '四', '五'][i] }}</span>
        <q-icon :name="t.icon" size="22px" class="task-icon" />
        <span>{{ t.text }}</span>
      </li>
    </ol>

    <h3 class="sec">協會宗旨</h3>
    <div class="aims">
      <p v-for="(a, i) in aims" :key="a.key" :class="['aim', i === 1 ? 'aim-red' : 'aim-dark']" v-html="a.html" />
    </div>

    <!-- ===== 協會起源 ===== -->
    <section id="origin" class="chapter">
      <p class="chapter-kicker">協會起源</p>
      <h2 class="chapter-title">世代變遷・全人健康・社會責任</h2>
    </section>

    <h3 class="sec">世代變遷：Z 世代</h3>
    <div class="panel dark">
      <div class="domes">
        <div v-for="g in generations" :key="g.name" class="dome-col">
          <div class="dome-people">{{ g.people }}</div>
          <div class="dome" :style="{ width: `${g.pct * 6}px`, height: `${g.pct * 3}px` }">
            <span>{{ g.pct }}%</span>
          </div>
          <div class="dome-years">{{ g.years }}</div>
          <div class="dome-name">{{ g.name }}</div>
        </div>
      </div>
      <p class="callout">Z 世代受到科技產物影響很大；可說是<b>第一個</b>自小同時生活在<b>電子虛擬與現實世界的原生世代</b>。</p>
      <p class="source">資料來源：CBNData 發布的《2020 Z 世代消費態度洞察報告》</p>
    </div>

    <h3 class="sec">世代變遷：信仰衰退</h3>
    <div class="panel dark">
      <div class="stats3">
        <div v-for="f in faith" :key="f.value" class="stat">
          <q-icon :name="f.icon" size="40px" class="stat-icon" />
          <div class="stat-value">{{ f.value }}</div>
          <p class="q-mb-none">{{ f.text }}</p>
        </div>
      </div>
      <p class="callout">Z 世代無宗教信仰也不信任宗教組織，但是<b>大部分的時間</b>都處於<b>壓力和孤獨的自處狀況</b>。</p>
      <p class="source">資料來源：Springtide Research Institute 發布的《2020 年宗教與年輕人現況》</p>
    </div>
    <blockquote class="quote">
      <p>佛教衰退的問題在於<b>教育</b>和<b>人才</b>上，而未必是<b>思想</b>的問題。</p>
      <cite>— 聖嚴法師</cite>
    </blockquote>
    <p class="aim-recall"><span class="status-chip">對應宗旨</span> 推廣並促進國內外產官學研界交流合作，讓佛教能普及弘化及永續發展。</p>

    <h3 class="sec">全人健康：壓力</h3>
    <div class="panel gold">
      <div class="stats3">
        <div v-for="(s, i) in stress" :key="s.value" :class="['drop', `drop-${i}`]">
          <div class="drop-value">{{ s.value }}<small>%</small></div>
          <p class="q-mb-none">{{ s.text }}</p>
        </div>
      </div>
      <p class="callout">研究顯示在所有世代中，<b>Z 世代</b>在<b>心理及社會幸福感</b>方面的<b>指數最低</b>，自我感受<b>壓力最大</b>。</p>
      <p class="source">資料來源：lululemon 發布的《2022 年全球幸福感報告》</p>
    </div>

    <h3 class="sec">全人健康：更大壓力</h3>
    <div class="panel gold big-stat">
      <svg class="stress-art" viewBox="0 0 200 160" aria-hidden="true">
        <path d="M30 120 C 20 60, 80 20, 120 40 S 190 60, 170 110 S 60 160, 30 120 Z" class="cloud" />
        <circle cx="100" cy="78" r="16" class="head" />
        <path d="M76 140 C 78 104, 122 104, 124 140 Z" class="body" />
        <g class="worry">
          <circle cx="52" cy="58" r="10" /><rect x="128" y="36" width="26" height="16" rx="3" /><circle cx="150" cy="96" r="9" />
          <path d="M60 30 l10 -6 l4 10 z" />
        </g>
      </svg>
      <div>
        <p class="q-mb-sm"><b>Z 世代</b>比<b>老一代</b>承受著<b>更大</b>的壓力</p>
        <div class="huge">91<small>%</small></div>
        <p class="q-mb-none"><b>Z 世代</b>過去一個月因<b>壓力</b>經歷<b>至少一種</b>身體或情緒<b>症狀</b></p>
        <p class="source">資料來源：SocialBee 發布的《Z 世代壓力調查》</p>
      </div>
    </div>

    <h3 class="sec">全人健康：身體健康</h3>
    <div class="panel light health">
      <div class="health-services">
        <div><b>醫療服務</b><small>（具醫療專業）</small></div>
        <span class="vs">VS</span>
        <div><b>業務服務</b><small>（醫療服務協助）</small></div>
      </div>
      <div class="health-layers">
        <span class="market market-full">全健康市場</span>
        <div class="layers">
          <div v-for="s in symptoms" :key="s.name" :class="['layer', { inner: s.inner }]">
            <b>{{ s.name }}</b><small>（{{ s.note }} {{ s.range }}）</small>
          </div>
        </div>
        <span class="market market-sub">亞健康市場</span>
      </div>
      <div class="health-certs">
        <div class="cert cert-down">
          <div class="cert-pct">85% <q-icon name="south" color="green-7" /></div>
          <ul><li v-for="c in bodyCerts" :key="c">{{ c }}</li><li>…</li></ul>
        </div>
        <div class="cert cert-up">
          <div class="cert-pct">15% <q-icon name="north" color="red-7" /></div>
          <ul><li v-for="c in beautyCerts" :key="c">{{ c }}</li><li>…</li></ul>
        </div>
      </div>
      <p class="callout callout-navy">目前市場上絕大部分的健康相關證照，還是著重在身體健康的部分，身心靈的部分是趨勢，比較以減壓療癒為主。</p>
    </div>

    <h3 class="sec">全人健康：身心靈健康</h3>
    <div class="panel light spirit">
      <svg class="hexes" viewBox="-170 -150 340 300" role="img" aria-label="全人健康六面向：環境、理智、感情、社會、身體、精神">
        <polygon v-for="h in hexes" :key="h.name" :points="hexPoints(h.x * 92, h.y * 92, 50)" class="hex" />
        <text v-for="h in hexes" :key="`t${h.name}`" :x="h.x * 92" :y="h.y * 92 + 6" class="hex-label">{{ h.name }}</text>
        <polygon :points="hexPoints(0, 0, 50)" class="hex hex-core" />
        <text x="0" y="-4" class="hex-core-label">全人</text>
        <text x="0" y="20" class="hex-core-label">健康</text>
        <text v-for="q in quotients" :key="q.name" :x="q.x * 125" :y="q.y * 125 + 5" class="q-label">{{ q.name }}</text>
      </svg>
      <div>
        <ol class="for-list">
          <li v-for="w in wholeHealthFor" :key="w">{{ w }}</li>
        </ol>
        <div class="focus">
          <span class="focus-other">其他證照：專注在<b>身體亞健康</b></span>
          <span class="focus-us">協會專注在<b>身心靈亞健康</b></span>
        </div>
      </div>
    </div>
    <p class="aim-recall"><span class="status-chip">對應宗旨</span> 透過佛教八正道傳統精神，讓現代人面臨的身心亞健康得以舒緩修復。</p>

    <h3 class="sec">社會責任：企業社會責任成績</h3>
    <div class="csr">
      <div class="panel green-soft csr-total">
        <p class="q-mb-xs"><b>2020 年 CSR 全產業成績總平均</b></p>
        <div class="huge green">73.0<small>分</small></div>
        <div v-for="c in csrTotal" :key="c.name" class="bar-row">
          <span class="bar-name">{{ c.name }}</span>
          <span class="bar"><i :style="{ width: pct(c.score) }" /></span>
          <span class="bar-num">{{ c.score.toFixed(1) }}</span>
        </div>
      </div>
      <div class="panel light">
        <p class="q-mb-sm"><b>五大產業別成績</b>：電信業居冠、金融保險業亞軍</p>
        <table class="table csr-table">
          <thead>
            <tr><th>產業</th><th class="num">總平均</th><th class="num">公司<br />治理</th><th class="num">員工、供應<br />鏈、顧客</th><th class="num">環境<br />績效</th><th class="num">社會<br />關懷</th></tr>
          </thead>
          <tbody>
            <tr v-for="r in csrIndustries" :key="r.name">
              <td><q-icon :name="r.icon" size="18px" class="q-mr-xs" />{{ r.name }}</td>
              <td class="num"><b>{{ r.total.toFixed(1) }}</b></td>
              <td v-for="(v, i) in r.s" :key="i" class="num">{{ v.toFixed(1) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <h3 class="sec">社會責任：企業重視的議題與 SDGs</h3>
    <div class="two-charts">
      <div class="panel light">
        <p class="q-mb-xs"><b>企業認為法令遵循、經營績效，是最重要兩大議題</b>（可複選，%）</p>
        <p class="chart-q">問：請以治理面、社會面、環境面，列出貴公司「前 10 大重大性議題」為何？</p>
        <div v-for="([n, v], i) in issues" :key="n" class="hbar">
          <span class="hbar-name">{{ n }}</span>
          <span class="hbar-track"><i :class="{ top: i === 0 }" :style="{ width: pct(v, 70) }" /></span>
          <span class="hbar-num" :class="{ top: i === 0 }">{{ v.toFixed(1) }}</span>
        </div>
        <p class="source">註：僅列出前 10 項</p>
      </div>
      <div class="panel light">
        <p class="q-mb-xs"><b>教育品質是國內企業呼應最高的 SDGs</b>（可複選，%）</p>
        <p class="chart-q">問：請根據貴公司永續發展藍圖，勾選出所呼應之聯合國永續發展目標（SDGs）？</p>
        <div v-for="([n, v], i) in sdgsTop" :key="n" class="hbar">
          <span class="hbar-name">{{ n }}</span>
          <span class="hbar-track"><i :class="{ top: i === 0 }" :style="{ width: pct(v, 90) }" /></span>
          <span class="hbar-num" :class="{ top: i === 0 }">{{ v.toFixed(1) }}</span>
        </div>
        <p class="source">註：僅列出前 10 項</p>
      </div>
    </div>

    <h3 class="sec">社會責任：善企業</h3>
    <div class="panel gold">
      <div class="flow">
        <span>品德</span><q-icon name="arrow_forward" /><span>善管理</span><q-icon name="arrow_forward" /><span>善企業</span><q-icon name="arrow_forward" /><span>善經濟</span>
      </div>
      <p class="q-mb-xs"><b>企業家從自我的品德要求，進而走向能利他的善管理的善企業</b></p>
      <p class="q-mb-none">
        「品善德澤利群生，善德者常住真心；企業家從事企業必須為善，創造員工與社會的福祉，創造環境的永續；企業只要是善企業，就能很容易達成理想的市場法則。」
      </p>
    </div>
    <div class="sdg-groups">
      <div v-for="g in sdgGroups" :key="g.name" class="sdg-group">
        <div class="sdg-head"><q-icon :name="g.icon" size="20px" /> {{ g.name }}</div>
        <div class="sdg-chips">
          <span v-for="n in g.goals" :key="n" class="sdg" :style="{ background: SDG[n][1] }"><b>{{ n }}</b>{{ SDG[n][0] }}</span>
        </div>
      </div>
    </div>

    <h3 class="sec">社會責任：佛法教育</h3>
    <div class="dharma">
      <div class="dharma-core">佛法<br />教育</div>
      <ul class="dharma-list">
        <li v-for="d in dharmaEdu" :key="d[1]">{{ d[0] }}<b>{{ d[1] }}</b>{{ d[2] }}</li>
      </ul>
    </div>
    <p class="aim-recall"><span class="status-chip">對應宗旨</span> 藉由傳統佛法教育，讓企業走向善管理的善企業，進而善盡社會責任。</p>

    <!-- ===== 營運特色 ===== -->
    <section id="operation" class="chapter">
      <p class="chapter-kicker">協會營運特色分享</p>
      <h2 class="chapter-title">品質系統・學習標章・標準組織</h2>
    </section>

    <h3 class="sec">營運方針：21 世紀佛教教育「現代化」的方向</h3>
    <blockquote class="quote">
      <p>
        在科技資訊化社會潮流，佛教教育若能<b>培育兼具佛學及資訊知能人才</b>，運用<b>資訊、傳播以及教育科技</b>等理論與方法，建立佛學「<b>數位典藏</b>」與「<b>知識管理</b>」系統以及「<b>線上學習</b>」環境，建構成「<b>數位神經系統</b>」，並且結合人文與藝術的資源，發展「<b>文化創意產業</b>」，達成真、善、美的人生目標。
      </p>
      <cite>— 惠敏法師，2003 年</cite>
    </blockquote>

    <h3 class="sec">品質系統：教學系統化</h3>
    <div class="mock-pair">
      <div class="mock">
        <div class="mock-tabs"><span>校方登入</span><span class="on">教師登入</span><span>學生登入</span></div>
        <div class="mock-body">
          <div class="mock-title">歡迎登入<br /><small>教師管理系統</small></div>
          <div class="mock-field"><q-icon name="person" /></div>
          <div class="mock-field"><q-icon name="lock" /></div>
          <div class="mock-btn">登入</div>
        </div>
        <q-icon name="co_present" size="64px" class="mock-art" />
        <p class="mock-cap">教學教案・課件與習題・備課參考資料</p>
      </div>
      <div class="mock">
        <div class="mock-tabs"><span>校方登入</span><span>教師登入</span><span class="on">學生登入</span></div>
        <div class="mock-body">
          <div class="mock-title">歡迎登入<br /><small>學生實訓系統</small></div>
          <div class="mock-field"><q-icon name="person" /></div>
          <div class="mock-field"><q-icon name="lock" /></div>
          <div class="mock-btn">登入</div>
        </div>
        <q-icon name="dashboard" size="64px" class="mock-art" />
        <p class="mock-cap">課件內容・習題問答・實例視頻</p>
      </div>
    </div>

    <h3 class="sec">品質系統：數位典藏</h3>
    <div class="tiles">
      <div v-for="a in archive" :key="a.title" class="tile">
        <div class="tile-art" :style="{ background: `linear-gradient(135deg, ${a.tone}, #f2f5fb)` }"><q-icon :name="a.icon" size="44px" /></div>
        <div class="tile-title">{{ a.title }}</div>
      </div>
    </div>
    <div class="courses">
      <div v-for="c in archiveCourses" :key="c.task" class="course">
        <div class="course-head">{{ c.task }}</div>
        <div v-for="it in c.items" :key="it" class="course-row"><span>{{ it }}</span><span class="course-act">線上・下載</span></div>
      </div>
    </div>

    <h3 class="sec">品質系統：知識管理暨線上學習</h3>
    <div class="km">
      <div class="km-nav"><span v-for="n in kmNav" :key="n">{{ n }}</span></div>
      <div class="km-section">產／品／編／輯</div>
      <div class="km-grid">
        <div v-for="(l, i) in kmLessons" :key="l" class="km-card">
          <span class="km-no">第 {{ i + 1 }} 節</span>
          <q-icon name="laptop_chromebook" size="34px" class="km-icon" />
          <p>{{ l }}</p>
        </div>
      </div>
      <div class="km-section">訂／單／管／理</div>
    </div>

    <h3 class="sec">品質系統：開發模組化</h3>
    <div class="addie">
      <div class="addie-labels"><span>階段</span><span>工作項目</span><span>產出</span></div>
      <div v-for="(a, i) in addie" :key="i" class="addie-col">
        <div class="addie-stage"><b>{{ a.letter }}</b>{{ a.stage }}</div>
        <div class="addie-work"><b>{{ a.work }}</b><ul><li v-for="x in a.items" :key="x">{{ x }}</li></ul></div>
        <div class="addie-out"><ul><li v-for="x in a.out" :key="x">{{ x }}</li></ul></div>
      </div>
      <div class="addie-loop"><q-icon name="sync" /> 檢核、回饋與修正（形成性評鑑）</div>
    </div>

    <h3 class="sec">品質系統：執行標準化</h3>
    <div class="steps">
      <div v-for="(s, i) in steps" :key="s.name" class="step">
        <div class="step-head" :style="{ background: s.tone }">{{ i + 1 }}. {{ s.name }}</div>
        <p class="q-mb-none">{{ s.text }}</p>
      </div>
    </div>

    <h3 class="sec">學習標章：數位神經系統一環</h3>
    <div class="panel night badges">
      <svg class="medal" viewBox="0 0 160 200" aria-hidden="true">
        <path d="M48 0 L80 60 L112 0 Z" class="ribbon" />
        <circle cx="80" cy="118" r="62" class="medal-ring" />
        <circle cx="80" cy="118" r="48" class="medal-face" />
        <path d="M80 86 l9 19 21 3 -15 15 4 21 -19 -10 -19 10 4 -21 -15 -15 21 -3 z" class="medal-star" />
      </svg>
      <div class="phone">
        <div class="phone-bar">我的學習勳章簿</div>
        <div class="phone-row">學習勳章 (1) <q-icon name="expand_less" /></div>
        <div class="phone-badge"><q-icon name="workspace_premium" size="40px" /><small>2022 課程結業勳章</small></div>
        <div class="phone-row">活動勳章 (1) <q-icon name="expand_more" /></div>
        <div class="phone-row">藝術勳章 (0) <q-icon name="expand_more" /></div>
      </div>
      <ul class="dots light-dots">
        <li v-for="(b, i) in badgeUses" :key="b" :class="{ hl: i === 3 }">{{ b }}</li>
      </ul>
    </div>

    <h3 class="sec">標準組織：標準組織規劃</h3>
    <div class="units">
      <div v-for="u in units" :key="u.name" class="unit">
        <div class="unit-head" :style="{ background: u.tone }">{{ u.name }}</div>
        <ul :style="{ color: u.tone }"><li v-for="x in u.items" :key="x"><span>{{ x }}</span></li></ul>
      </div>
    </div>

    <h3 class="sec">標準組織：組織架構</h3>
    <div class="org">
      <div class="org-top">
        <div class="org-box">監事會</div>
        <div class="org-box main">會員大會</div>
        <div class="org-box">理事會</div>
      </div>
      <div class="org-line" />
      <div class="org-box main narrow">理事長</div>
      <div class="org-line" />
      <div class="org-branches">
        <div class="org-branch">
          <div class="org-box red">秘書處</div>
          <div class="org-leaves"><span v-for="s in secretariat" :key="s">{{ s }}</span></div>
        </div>
        <div class="org-branch wide">
          <div class="org-box red">工作委員會</div>
          <div class="org-leaves"><span v-for="c in committees" :key="c">{{ c }}</span><span class="cdc">課程發展中心</span></div>
          <div class="org-courses">
            <span v-for="c in courses" :key="c">{{ c }}</span>
            <q-icon name="arrow_forward" />
            <span class="gate">審查・紀律</span>
            <q-icon name="arrow_forward" />
            <span class="gate">指導委員會</span>
            <q-icon name="arrow_forward" />
            <span class="gate red">認證暨品質中心</span>
          </div>
        </div>
      </div>
    </div>

    <h3 class="sec">寺廟教育文創</h3>
    <div class="temple">
      <div v-for="t in temple" :key="t.title" class="temple-tile" :style="{ background: `linear-gradient(160deg, ${t.from}, ${t.to})` }">
        <q-icon :name="t.icon" size="46px" />
        <span>{{ t.title }}</span>
      </div>
    </div>

    <!-- ===== 總結 ===== -->
    <section id="summary" class="chapter">
      <p class="chapter-kicker">總結</p>
      <h2 class="chapter-title">不辦教育，佛教就沒有明天</h2>
    </section>
    <blockquote class="quote">
      <p>
        聖嚴法師有感於當時社會對佛教的誤解，甚至視之為迷信，雖然僧尼無數，寺院儼然，但是<b>佛教的形象卻普遍地落於民間信仰的層次</b>。聖嚴法師懇切呼籲<b>推動佛教教育事業為當務之急</b>。
      </p>
      <cite>— 1985 年</cite>
    </blockquote>
    <div class="lifelong">
      <div><b>教育</b><small>未來生活之準備</small></div>
      <q-icon name="arrow_forward" size="28px" />
      <div><b>學習</b><small>提升生活水平價值</small></div>
      <q-icon name="arrow_forward" size="28px" />
      <div><b>終身</b><small>活到老學到老</small></div>
    </div>

    <!-- ===== 會員招募 ===== -->
    <section id="members" class="chapter green">
      <p class="chapter-kicker">協會會員招募中</p>
      <h2 class="chapter-title">產業的專業會員</h2>
    </section>
    <h3 class="sec">會員人設</h3>
    <div class="grid">
      <div v-for="m in members" :key="m.role" class="card member">
        <q-icon :name="m.icon" size="34px" class="pillar-icon" />
        <h4>{{ m.role }}</h4>
        <p class="q-mb-none" v-html="m.html" />
      </div>
    </div>

    <div class="closing">
      <svg class="closing-tree" viewBox="0 0 400 200" preserveAspectRatio="xMidYMax slice" aria-hidden="true">
        <path d="M200 200 C 196 150, 190 120, 170 96 M200 200 C 204 150, 214 118, 238 92 M188 140 C 160 128, 132 112, 104 110 M212 136 C 244 124, 276 110, 306 112" />
        <circle cx="170" cy="80" r="40" /><circle cx="236" cy="76" r="44" /><circle cx="110" cy="100" r="34" /><circle cx="300" cy="104" r="34" /><circle cx="204" cy="52" r="36" />
      </svg>
      <p>我們有決心與行動　讓世界佛教教育協會成長茁壯　<b>歡迎您的加入</b></p>
      <small>敬請指教</small>
    </div>
  </q-page>
</template>

<style scoped>
.assoc {
  --red: #9e2b25;
  --gold-bg: #f7d36a;
  --dark: #2b2421;
  --navy: #1d2b5a;
}
.assoc h4 {
  margin: 6px 0;
  font-size: 1.1rem;
}

/* 封面與目錄 */
.cover {
  display: flex;
  align-items: center;
  gap: 24px;
  padding: 28px 24px;
  border-radius: 18px;
  background: radial-gradient(circle at 15% 30%, #fff7e6, var(--ground-sunk));
}
.cover h1 {
  margin: 4px 0 6px;
}
.cover-mark {
  width: 96px;
  flex: none;
}
.cover-mark circle {
  fill: var(--red);
  opacity: 0.35;
  mix-blend-mode: multiply;
}
.cover-mark circle:nth-child(2) {
  fill: var(--leaf);
}
.cover-mark circle:nth-child(3) {
  fill: #c89b3c;
}
.eyebrow {
  margin: 0;
  letter-spacing: 0.12em;
  font-size: 0.8rem;
  color: var(--ink-faint);
  text-transform: uppercase;
}
.toc {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin: 16px 0 8px;
}
.toc a {
  padding: 4px 14px;
  border-radius: 999px;
  border: 1px solid var(--rule);
  color: var(--ink-soft);
  text-decoration: none;
  font-size: 0.9rem;
}

/* 章節頁與小節標題 */
.chapter {
  margin: 44px 0 8px;
  padding: 26px 24px;
  border-radius: 16px;
  background: linear-gradient(120deg, #2e6f9e, #48a3d6);
  color: #fff;
  scroll-margin-top: 70px;
}
.chapter.green {
  background: linear-gradient(120deg, #2f7a3e, #4fae5e);
}
.chapter-kicker {
  margin: 0;
  font-weight: 700;
  opacity: 0.9;
}
.chapter-title {
  margin: 4px 0 0 !important;
  padding-top: 8px;
  border-top: 2px solid rgba(255, 255, 255, 0.6);
  color: #fff !important;
}
.sec {
  margin: 32px 0 12px !important;
  color: var(--red);
  font-family: var(--wenkai);
  font-size: 1.35rem !important;
}

/* 簡介 */
.intro-grid {
  display: grid;
  grid-template-columns: minmax(180px, 260px) 1fr;
  gap: 24px;
  align-items: center;
}
.pillar-ring path {
  stroke: #fff;
  stroke-width: 2;
}
.ring-a { fill: #cfe3e6; }
.ring-b { fill: #b9d6da; }
.ring-c { fill: #a8ccd1; }
.ring-core { fill: var(--red); }
.ring-label {
  font-size: 15px;
  font-weight: 700;
  fill: #237a85;
  text-anchor: middle;
}
.ring-core-label {
  font-size: 14px;
  fill: #fff;
  text-anchor: middle;
  font-weight: 700;
}
.dots {
  padding-left: 1.2em;
}
.dots li {
  margin: 6px 0;
}
.pillar-icon {
  color: var(--leaf);
}

/* 任務 */
.tasks {
  list-style: none;
  padding: 0;
  display: grid;
  gap: 10px;
}
.tasks li {
  display: grid;
  grid-template-columns: 34px 28px 1fr;
  align-items: start;
  gap: 8px;
  padding: 12px 14px;
  background: var(--ground-raised);
  border: 1px solid var(--rule);
  border-radius: 10px;
}
.task-no {
  font-family: var(--wenkai);
  font-size: 1.3rem;
  color: var(--red);
  line-height: 1.2;
}
.task-icon {
  color: var(--leaf);
  margin-top: 2px;
}

/* 宗旨三色帶 */
.aims {
  display: grid;
  gap: 10px;
}
.aim {
  margin: 0;
  padding: 18px 20px;
  border-radius: 8px;
  color: #f3e27a;
  text-align: center;
  font-size: 1.05rem;
}
.aim :deep(b) {
  color: #fff;
  font-size: 1.15em;
}
.aim-dark { background: #3f3f3f; }
.aim-red { background: var(--red); color: #fff; }
.aim-red :deep(b) { color: #ffe14d; }
.aim-recall {
  margin-top: 14px;
  color: var(--ink-soft);
}

/* 面板 */
.panel {
  border-radius: 14px;
  padding: 22px;
}
.panel.dark {
  background: #262322;
  color: #eee;
}
.panel.gold {
  background: var(--gold-bg);
  color: var(--dark);
}
.panel.light {
  background: var(--ground-raised);
  border: 1px solid var(--rule);
}
.panel.green-soft {
  background: linear-gradient(160deg, #e9f4e3, #cfe8c4);
}
.panel.night {
  background: #15171d;
  color: #eee;
}
.callout {
  margin: 18px 0 4px;
  padding: 10px 14px;
  background: #5b2e86;
  color: #fff;
  border-radius: 6px;
  text-align: center;
}
.callout :deep(b),
.callout b {
  color: #9be56b;
}
.callout-navy {
  background: var(--navy);
}
.source {
  margin: 6px 0 0;
  font-size: 0.78rem;
  opacity: 0.75;
  text-align: right;
}

/* 世代圓頂圖 */
.domes {
  display: flex;
  justify-content: space-around;
  align-items: flex-end;
  gap: 8px;
  flex-wrap: wrap;
  border-bottom: 2px solid #e2b33c;
  padding-bottom: 0;
}
.dome-col {
  display: flex;
  flex-direction: column;
  align-items: center;
}
.dome-people {
  font-weight: 700;
  margin-bottom: 6px;
}
.dome {
  border-radius: 999px 999px 0 0;
  display: grid;
  place-items: center;
  font-size: 1.4rem;
  min-width: 70px;
  min-height: 36px;
}
.dome-col:nth-child(1) .dome { background: #0b6fbf; color: #8fe2ff; }
.dome-col:nth-child(2) .dome { background: #0f6b2a; color: #9be56b; }
.dome-col:nth-child(3) .dome { background: #f5ef3a; color: #b19a16; }
.dome-col:nth-child(4) .dome { background: #d39494; color: #9e2b45; }
.dome-years {
  margin-top: 8px;
  padding: 2px 10px;
  background: #f2f2f2;
  color: #333;
  border-radius: 6px;
  font-weight: 700;
  font-size: 0.85rem;
}
.dome-name {
  margin-top: 4px;
  font-weight: 700;
}
.dome-col:nth-child(1) .dome-name { color: #4fb3ff; }
.dome-col:nth-child(2) .dome-name { color: #38b45b; }
.dome-col:nth-child(3) .dome-name { color: #f5ef3a; }
.dome-col:nth-child(4) .dome-name { color: #e7a3a3; }

/* 統計格 */
.stats3 {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}
.stat {
  text-align: center;
}
.stat-icon {
  color: #36c5d8;
}
.stat-value {
  font-size: 2rem;
  font-weight: 800;
  color: #36d8e8;
}
.drop {
  border-radius: 50% 50% 50% 8px;
  padding: 30px 22px;
  text-align: center;
  color: #fff;
  aspect-ratio: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  max-width: 260px;
  margin: 0 auto;
}
.drop-0 { background: #7d7f84; }
.drop-1 { background: #3d3d3f; }
.drop-2 { background: #050505; }
.drop-value {
  font-size: 2.6rem;
  font-weight: 800;
  line-height: 1;
  margin-bottom: 8px;
}
.drop-value small,
.huge small {
  font-size: 0.5em;
}
.big-stat {
  display: grid;
  grid-template-columns: minmax(160px, 260px) 1fr;
  gap: 24px;
  align-items: center;
}
.huge {
  font-size: 4.4rem;
  font-weight: 800;
  line-height: 1;
  color: #b8121b;
  margin: 4px 0 8px;
}
.huge.green {
  color: #c21b1b;
  font-size: 3rem;
}
.stress-art .cloud { fill: #bccaf4; }
.stress-art .head,
.stress-art .body { fill: #23202a; }
.stress-art .worry { fill: none; stroke: #6b6f86; stroke-width: 2; }

/* 身體健康分層 */
.health {
  display: grid;
  grid-template-columns: auto 1fr auto;
  gap: 20px;
  align-items: center;
}
.health .callout {
  grid-column: 1 / -1;
}
.health-services {
  display: flex;
  flex-direction: column;
  gap: 10px;
  color: #7c7a12;
  text-align: center;
}
.health-services small {
  display: block;
  color: var(--ink-faint);
}
.vs {
  font-weight: 700;
  color: var(--ink);
}
.health-layers {
  display: grid;
  grid-template-columns: auto 1fr auto;
  align-items: center;
  gap: 8px;
}
.market {
  writing-mode: vertical-rl;
  font-weight: 700;
  font-size: 0.9rem;
}
.market-full { color: #1f6db5; }
.market-sub { color: var(--red); }
.layers {
  display: grid;
  gap: 4px;
}
.layer {
  text-align: center;
  padding: 8px 6px;
  border-radius: 6px;
  background: #fff;
  border: 1px solid var(--rule);
}
.layer.inner {
  background: #f1f4f8;
  margin: 0 14px;
}
.layer small {
  display: block;
  color: #1aa3d9;
}
.health-certs {
  display: grid;
  gap: 14px;
}
.cert-pct {
  font-size: 1.8rem;
  font-weight: 800;
  color: #6c2c9c;
}
.cert ul {
  margin: 4px 0 0;
  padding-left: 1.1em;
  color: var(--ink-soft);
  font-size: 0.9rem;
}

/* 身心靈六角 */
.spirit {
  display: grid;
  grid-template-columns: minmax(220px, 340px) 1fr;
  gap: 20px;
  align-items: center;
}
.hex {
  fill: #f1f1ef;
  stroke: #d7d7d2;
  stroke-width: 1.5;
}
.hex-core {
  fill: #fff;
}
.hex-label {
  font-size: 15px;
  fill: #777;
  text-anchor: middle;
}
.hex-core-label {
  font-size: 18px;
  font-weight: 700;
  fill: var(--red);
  text-anchor: middle;
}
.q-label {
  font-size: 15px;
  font-weight: 700;
  fill: #7a7a14;
  text-anchor: middle;
}
.for-list li {
  margin: 4px 0;
}
.focus {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 12px;
}
.focus span {
  padding: 6px 12px;
  border-radius: 8px;
}
.focus-other {
  background: var(--ground-sunk);
}
.focus-us {
  background: var(--red);
  color: #fff;
}

/* CSR */
.csr {
  display: grid;
  grid-template-columns: minmax(240px, 1fr) 2fr;
  gap: 16px;
}
.bar-row {
  display: grid;
  grid-template-columns: 8.5em 1fr 3em;
  gap: 8px;
  align-items: center;
  margin: 6px 0;
  font-size: 0.9rem;
}
.bar {
  height: 10px;
  background: rgba(255, 255, 255, 0.7);
  border-radius: 999px;
  overflow: hidden;
}
.bar i {
  display: block;
  height: 100%;
  background: var(--leaf);
}
.bar-num {
  text-align: right;
  font-variant-numeric: tabular-nums;
}
.csr-table {
  font-size: 0.9rem;
}
.csr-table .num {
  text-align: right;
  font-variant-numeric: tabular-nums;
}
.two-charts {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 16px;
}
.chart-q {
  font-size: 0.82rem;
  color: var(--ink-faint);
}
.hbar {
  display: grid;
  grid-template-columns: minmax(8em, 13em) 1fr 3em;
  gap: 8px;
  align-items: center;
  margin: 5px 0;
  font-size: 0.88rem;
}
.hbar-name {
  text-align: right;
  line-height: 1.25;
}
.hbar-track {
  height: 14px;
  background: var(--ground-sunk);
}
.hbar-track i {
  display: block;
  height: 100%;
  background: #cfcfcf;
}
.hbar-track i.top {
  background: #3aa82a;
}
.hbar-num {
  font-variant-numeric: tabular-nums;
}
.hbar-num.top {
  color: #3aa82a;
  font-weight: 800;
}

/* 善企業與 SDGs */
.flow {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
  font-size: 1.3rem;
  font-weight: 800;
  color: #b8121b;
  margin-bottom: 10px;
}
.sdg-groups {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 14px;
  margin-top: 14px;
  padding: 16px;
  background: #476f4d;
  border-radius: 14px;
}
.sdg-head {
  color: #fff;
  font-weight: 700;
  margin-bottom: 8px;
}
.sdg-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}
.sdg {
  color: #fff;
  border-radius: 6px;
  padding: 3px 8px;
  font-size: 0.82rem;
}
.sdg b {
  margin-right: 4px;
}

/* 佛法教育 */
.dharma {
  display: grid;
  grid-template-columns: 150px 1fr;
  gap: 22px;
  align-items: center;
}
.dharma-core {
  aspect-ratio: 1;
  display: grid;
  place-items: center;
  text-align: center;
  font-size: 2rem;
  font-weight: 800;
  color: #555;
  background: #c4c4c4;
  border-radius: 10px;
  line-height: 1.3;
}
.dharma-list {
  list-style: none;
  padding: 0;
  margin: 0;
  border-left: 2px solid #d4d4d4;
}
.dharma-list li {
  padding: 6px 0 6px 16px;
  font-weight: 600;
  color: #555;
}
.dharma-list b {
  color: var(--red);
  border-bottom: 2px solid var(--red);
}

/* 引言 */
.quote {
  margin: 16px 0;
  padding: 20px 24px;
  background: #e2e2e2;
  border-radius: 8px;
  text-align: center;
}
.quote p {
  margin: 0;
  font-size: 1.08rem;
  line-height: 1.9;
}
.quote b {
  color: var(--red);
}
.quote cite {
  display: block;
  margin-top: 8px;
  font-style: normal;
  color: var(--ink-faint);
}

/* 系統示意 */
.mock-pair {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 16px;
}
.mock {
  position: relative;
  padding: 20px;
  border-radius: 12px;
  background: linear-gradient(135deg, #cfdcfa, #f3f6ff);
  overflow: hidden;
}
.mock-tabs {
  display: flex;
  gap: 4px;
  font-size: 0.75rem;
}
.mock-tabs span {
  padding: 3px 8px;
  background: #fff;
  border-radius: 4px 4px 0 0;
  color: #777;
}
.mock-tabs .on {
  background: #f6a623;
  color: #fff;
}
.mock-body {
  width: 60%;
  background: #fff;
  border-radius: 0 8px 8px 8px;
  padding: 16px;
  display: grid;
  gap: 10px;
}
.mock-title {
  color: #f6a623;
  font-weight: 800;
}
.mock-title small {
  font-weight: 600;
}
.mock-field {
  border-bottom: 1px solid #ddd;
  color: #f26f5b;
  padding-bottom: 4px;
}
.mock-btn {
  background: #f6a623;
  color: #fff;
  text-align: center;
  border-radius: 999px;
  padding: 6px;
  font-weight: 700;
}
.mock-art {
  position: absolute;
  right: 22px;
  top: 70px;
  color: #6b8ff0;
}
.mock-cap {
  margin: 14px 0 0;
  color: #6233a8;
  font-weight: 700;
  text-align: center;
}
.tiles {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 12px;
}
.tile {
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.06);
}
.tile-art {
  aspect-ratio: 16 / 10;
  display: grid;
  place-items: center;
  color: #fff;
}
.tile-title {
  padding: 8px 10px;
  font-size: 0.9rem;
  color: #555;
}
.courses {
  display: grid;
  gap: 12px;
  margin-top: 14px;
}
.course {
  border: 1px solid #ccc;
  background: #fff;
}
.course-head {
  background: #4b83ff;
  color: #fff;
  padding: 8px 14px;
  font-weight: 700;
}
.course-row {
  display: flex;
  justify-content: space-between;
  padding: 8px 14px 8px 30px;
  font-size: 0.9rem;
  color: #555;
}
.course-act {
  color: #aaa;
}
.km {
  border-radius: 10px;
  overflow: hidden;
  background: #e3f0ff;
}
.km-nav {
  display: flex;
  flex-wrap: wrap;
  gap: 14px;
  padding: 10px 16px;
  background: #4b83ff;
  color: #fff;
  font-size: 0.85rem;
}
.km-section {
  text-align: center;
  letter-spacing: 0.3em;
  padding: 14px;
  color: #333;
  font-size: 1.1rem;
}
.km-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 10px;
  padding: 0 14px;
}
.km-card {
  position: relative;
  background: #fff;
  border-radius: 6px;
  overflow: hidden;
}
.km-card .km-icon {
  display: block;
  margin: 0 auto;
  padding: 18px 0;
  width: 100%;
  background: linear-gradient(135deg, #6a4fc8, #4b83ff);
  color: #fff;
}
.km-no {
  position: absolute;
  left: 6px;
  top: 6px;
  background: #2f8cff;
  color: #fff;
  font-size: 0.72rem;
  padding: 1px 6px;
  border-radius: 4px;
}
.km-card p {
  margin: 0;
  padding: 8px;
  font-size: 0.82rem;
  text-align: center;
  color: #334;
}

/* ADDIE */
.addie {
  display: grid;
  grid-template-columns: 56px repeat(5, minmax(0, 1fr));
  grid-template-rows: 64px auto auto auto;
  gap: 8px;
}
/* 每一欄跨三列並共用同一組列高（subgrid），工作項目與產出才會左右對齊 */
.addie-labels,
.addie-col {
  grid-row: 1 / span 3;
  display: grid;
  grid-template-rows: subgrid;
}
.addie-labels span {
  display: grid;
  place-items: center;
  writing-mode: vertical-rl;
  background: #d9d9d9;
  font-weight: 800;
  border-radius: 4px;
}
.addie-stage {
  background: linear-gradient(#3d7de0, #2c62c0);
  color: #fff;
  display: grid;
  place-items: center;
  text-align: center;
  font-weight: 700;
  line-height: 1.2;
  border-radius: 4px;
}
.addie-stage b {
  display: block;
  font-size: 1.2rem;
}
.addie-work {
  background: #dce6f5;
  border: 1px dashed #4a77c7;
  border-radius: 10px;
  padding: 10px;
  font-size: 0.86rem;
}
.addie-work ul,
.addie-out ul {
  margin: 6px 0 0;
  padding-left: 1.1em;
}
.addie-out {
  border: 1px solid #4a77c7;
  background: #fff;
  padding: 8px;
  font-size: 0.84rem;
  color: #c0272d;
}
.addie-loop {
  grid-column: 2 / -1;
  text-align: center;
  font-weight: 700;
  border-top: 2px solid #9b1c1c;
  padding-top: 6px;
}

/* 執行標準化 */
.steps {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(170px, 1fr));
  gap: 12px;
}
.step {
  border: 1px solid #234;
  background: #fff;
}
.step-head {
  color: #fff;
  font-weight: 700;
  padding: 8px 12px;
}
.step p {
  padding: 12px;
  color: #3b4d6b;
  text-align: center;
}

/* 學習標章 */
.badges {
  display: grid;
  grid-template-columns: 140px 220px 1fr;
  gap: 22px;
  align-items: center;
}
.ribbon { fill: #3a7bd5; }
.medal-ring { fill: #e7b93c; }
.medal-face { fill: #2b2b2b; }
.medal-star { fill: #e7b93c; }
.phone {
  background: linear-gradient(#1f7ae0, #22c3e6);
  border-radius: 16px;
  padding: 12px;
  font-size: 0.82rem;
}
.phone-bar {
  font-weight: 700;
  margin-bottom: 8px;
}
.phone-row {
  display: flex;
  justify-content: space-between;
  background: #fff;
  color: #333;
  padding: 6px 8px;
  border-bottom: 1px solid #eee;
}
.phone-badge {
  background: #fff;
  color: #555;
  padding: 8px;
  display: flex;
  align-items: center;
  gap: 8px;
}
.light-dots li.hl {
  color: #ffe14d;
  font-weight: 700;
}

/* 標準組織 */
.units {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 12px;
}
.unit {
  border: 1px solid #234;
  background: #fff;
}
.unit-head {
  color: #fff;
  text-align: center;
  font-weight: 700;
  padding: 10px;
}
.unit ul {
  margin: 0;
  padding: 12px 14px 12px 28px;
}
.unit li {
  margin: 6px 0;
}
.unit li span {
  color: var(--ink);
}
.org {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0;
  padding: 18px;
  background: var(--ground-raised);
  border: 1px solid var(--rule);
  border-radius: 14px;
  overflow-x: auto;
}
.org-top {
  display: flex;
  gap: 16px;
  align-items: center;
}
.org-box {
  padding: 8px 18px;
  background: #f1f1f1;
  border: 1px solid #777;
  color: #555;
  text-align: center;
  white-space: nowrap;
}
.org-box.main {
  font-weight: 700;
  color: #333;
}
.org-box.red {
  background: #c00000;
  color: #fff;
  border-color: #c00000;
}
.org-line {
  width: 2px;
  height: 18px;
  background: #777;
}
.org-branches {
  display: flex;
  gap: 18px;
  flex-wrap: wrap;
  justify-content: center;
  margin-top: 4px;
}
.org-branch {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}
.org-leaves {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  justify-content: center;
  max-width: 520px;
}
.org-leaves span,
.org-courses span {
  padding: 6px 10px;
  background: #f1f1f1;
  border: 1px solid #777;
  font-size: 0.85rem;
}
.org-leaves .cdc {
  background: #e9ddf5;
}
.org-courses {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  align-items: center;
  justify-content: center;
  max-width: 560px;
}
.org-courses span {
  background: #7030a0;
  color: #fff;
  border-color: #7030a0;
}
.org-courses .gate {
  background: #f1f1f1;
  color: #555;
  border-color: #777;
}
.org-courses .gate.red {
  background: #c00000;
  color: #fff;
  border-color: #c00000;
}

/* 寺廟教育文創 */
.temple {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(170px, 1fr));
  gap: 10px;
}
.temple-tile {
  aspect-ratio: 4 / 3;
  border-radius: 10px;
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  padding: 12px;
  color: #fff;
  gap: 6px;
  text-shadow: 0 1px 4px rgba(0, 0, 0, 0.4);
}
.temple-tile span {
  font-weight: 700;
}

/* 總結 */
.lifelong {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  align-items: center;
  gap: 14px;
  color: #6c2c9c;
  margin: 12px 0;
}
.lifelong div {
  text-align: center;
}
.lifelong b {
  display: block;
  font-size: 2rem;
  font-family: var(--wenkai);
}
.lifelong small {
  color: var(--ink-faint);
}

/* 會員與結尾 */
.member h4 {
  color: #2a6b1f;
}
.member :deep(b) {
  color: var(--red);
}
.closing {
  position: relative;
  margin: 36px 0 12px;
  padding: 150px 20px 26px;
  border-radius: 16px;
  overflow: hidden;
  background: linear-gradient(#3d3d3d, #111);
  color: #fff;
  text-align: center;
}
.closing-tree {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
}
.closing-tree path {
  fill: none;
  stroke: #000;
  stroke-width: 10;
  stroke-linecap: round;
}
.closing-tree circle {
  fill: rgba(255, 255, 255, 0.12);
}
.closing p {
  position: relative;
  display: inline-block;
  margin: 0;
  padding: 10px 22px;
  background: var(--red);
  border-radius: 999px;
  font-weight: 700;
}
.closing p b {
  color: #ffe14d;
}
.closing small {
  position: relative;
  display: block;
  margin-top: 10px;
  opacity: 0.8;
}

/* 站內 .page p／li 預設是深色字；在深色、彩色面板裡改為繼承面板的字色 */
.chapter p,
.panel.dark p,
.panel.night li,
.drop p,
.aim,
.closing p,
.closing small {
  color: inherit !important;
}
.chapter-kicker {
  color: #fff !important;
}
.callout,
.callout p {
  color: #fff !important;
}
.panel.dark .stat p {
  color: #d6d6d6 !important;
}
.drop p {
  color: #f3f3f3 !important;
  font-size: 0.92rem;
  line-height: 1.6;
}
.addie li,
.cert li,
.unit li,
.dharma-list li,
.course-row,
.km-card p,
.tile-title,
.step p,
.mock-cap {
  font-size: 0.9rem;
  line-height: 1.5;
}
.aim {
  color: #f3e27a !important;
}
.aim-red {
  color: #fff !important;
}

@media (max-width: 760px) {
  .intro-grid,
  .big-stat,
  .spirit,
  .csr,
  .dharma,
  .badges,
  .health {
    grid-template-columns: 1fr;
  }
  .cover {
    flex-direction: column;
    text-align: center;
  }
  .addie {
    grid-template-columns: 1fr;
  }
  .addie-labels {
    display: none;
  }
  .addie-col {
    grid-row: auto;
    grid-template-rows: auto;
  }
  .addie-loop {
    grid-column: auto;
  }
  .hbar {
    grid-template-columns: 9em 1fr 3em;
  }
  .csr-table {
    display: block;
    overflow-x: auto;
  }
}
</style>
