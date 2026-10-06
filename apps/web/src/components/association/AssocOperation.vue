<script setup lang="ts">
// 第三章「協會營運特色分享」：品質系統、學習標章、總結。
import ScreenShot from './ScreenShot.vue'

const img = (name: string) => `/images/association/${name}.webp`

const addie = [
  {
    letter: 'A',
    stage: '分析階段',
    task: '進行課程需求分析',
    items: ['進行市場需求分析', '界定主要學習對象', '確定教學／訓練目的', '重組職能內涵'],
    outputs: ['先備條件', '課程架構或學習路徑'],
  },
  {
    letter: 'D',
    stage: '設計階段',
    task: '訂定教學／訓練目標與課程內容',
    items: ['訂定教學／訓練目標', '規劃課程內容'],
    outputs: ['教學／訓練目標', '課程內容（授課大綱）'],
  },
  {
    letter: 'D',
    stage: '發展階段',
    task: '選擇教學方法，準備教材與教學資源',
    items: ['選擇教學方法', '撰寫教材與講義', '準備教學資源', '發展評量工具'],
    outputs: ['教學方法', '教學資源'],
  },
  {
    letter: 'I',
    stage: '實施階段',
    task: '辦理課程並進行評量',
    items: ['辦理課程', '實施評量活動'],
    outputs: ['課程辦理相關文件記錄'],
  },
  {
    letter: 'E',
    stage: '評估階段',
    task: '評量學習成果，監控評估課程成效及回饋修正',
    items: ['評量學習成果', '收集學習成果證據', '監控評估與回饋修正'],
    outputs: ['學習成果評量說明文件', '學習成果證據相關文件', '監控評估機制說明與記錄文件'],
  },
]

const steps = [
  { name: '問卷調查', lines: ['填寫數量需', '達到回收標準'], color: '#3f7a3a' },
  { name: '課程規劃', lines: ['課程企劃', '給予課程規畫建議'], color: '#22406b' },
  { name: '課程製作', lines: ['老師正式製作課程', '考核老師上課教學品質'], color: '#c8442b' },
  { name: '課程說明', lines: ['透過說明會銷講方式', '學生互動回饋達成交易'], color: '#76761f' },
  { name: '正式開課', lines: ['通知購課同學上課', '並於課後收集教學回饋'], color: '#93641a' },
]

</script>

<template>
  <section id="operation" class="chapter">
    <p class="kicker">協會營運特色分享</p>
    <h2>品質系統・學習標章</h2>
  </section>

  <h2>品質系統</h2>
  <h3>教學系統化</h3>
  <div class="shots two">
    <ScreenShot :src="img('teach-teacher')" alt="教師管理系統登入畫面" caption="教學教案・課件與習題・備課參考資料" />
    <ScreenShot :src="img('teach-student')" alt="學生實訓系統登入畫面" caption="課件內容・習題問答・實例視頻" />
  </div>

  <h3>數位典藏</h3>
  <p>課程內容以數位方式保存，依課程與單元整理，可線上觀看與下載。</p>
  <div class="shots two">
    <ScreenShot :src="img('archive-courses')" alt="數位典藏的課程列表畫面" />
    <ScreenShot :src="img('archive-lessons')" alt="數位典藏的課程單元列表，每個單元可線上觀看或下載" />
  </div>

  <h3>知識管理暨線上學習</h3>
  <div class="shots">
    <ScreenShot :src="img('elearning')" alt="線上學習平台依主題分節的課程畫面" />
  </div>

  <h3>開發模組化</h3>
  <div class="addie">
    <div v-for="(s, i) in addie" :key="i" class="phase">
      <div class="phase-head">
        <span class="letter">{{ s.letter }}</span>
        <span>{{ s.stage }}</span>
      </div>
      <div class="phase-body">
        <p class="phase-task">{{ s.task }}</p>
        <ul>
          <li v-for="it in s.items" :key="it">{{ it }}</li>
        </ul>
      </div>
      <div class="phase-out">
        <p class="out-label">產出</p>
        <ul>
          <li v-for="o in s.outputs" :key="o">{{ o }}</li>
        </ul>
      </div>
    </div>
  </div>
  <p class="loop">↺ 每一階段都經過<b>檢核、回饋與修正</b>（形成性評鑑）</p>

  <h3>執行標準化</h3>
  <ol class="steps">
    <li v-for="(s, i) in steps" :key="s.name" :style="{ '--c': s.color }">
      <p class="step-head">{{ i + 1 }}. {{ s.name }}</p>
      <p v-for="l in s.lines" :key="l" class="step-line">{{ l }}</p>
    </li>
  </ol>

  <h2>學習標章：數位神經系統一環</h2>
  <div class="figure dark poap">
    <div>
      <img :src="img('poap-scene')" alt="眾人拿著手機領取 POAP 數位學習標章的插圖" class="poap-scene" loading="lazy" />
      <ul class="poap-list">
        <li>學員／老師<b class="y">未來課程</b>身份認證使用</li>
        <li>與學員活動訊息的溝通管道</li>
        <li>結業或是證書的<b class="y">數位證明</b></li>
        <li><b class="y">學員／老師社群互動的應用</b></li>
      </ul>
    </div>
    <div class="poap-shots">
      <a :href="img('poap-phone')" target="_blank" rel="noopener"><img :src="img('poap-phone')" alt="POAP 錢包裡收藏的數位標章" class="phone" loading="lazy" /></a>
      <a :href="img('poap-badges')" target="_blank" rel="noopener"><img :src="img('poap-badges')" alt="POAP 動章簿：學習動章、活動動章與動章內容" class="badges" loading="lazy" /></a>
    </div>
  </div>

  <h2>總結</h2>
  <div class="quote">
    <p class="quote-title">不辦教育，佛教就沒有明天</p>
    <p>
      <b>聖嚴法師</b>有感於當時社會對佛教的誤解，甚至視之為迷信，雖然僧尼無數，寺院儼然，但是<span class="hl">佛教的形象卻普遍地落於民間信仰的層次</span>。
      <b>聖嚴法師</b>懇切呼籲<span class="hl">推動佛教教育事業為當務之急</span>。
    </p>
    <p class="by">1985 年</p>
  </div>
  <div class="lifelong">
    <p class="flow">教育 <span>→</span> 學習 <span>→</span> 終身</p>
    <div class="flow-notes">
      <span>未來生活之準備</span>
      <span>提升生活水平價值</span>
      <span>活到老學到老</span>
    </div>
  </div>
</template>

<style scoped>
.shots {
  margin: 12px 0 8px;
}
.shots.two {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.addie {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 12px;
  margin-top: 12px;
}
.phase {
  display: flex;
  flex-direction: column;
  border-radius: 12px;
  overflow: hidden;
  background: var(--ground-raised);
  border: 1px solid var(--rule);
}
.phase-head {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  background: linear-gradient(135deg, #2c6e9f, #3f8cc4);
  color: #fff;
  font-weight: 700;
}
.letter {
  display: grid;
  place-items: center;
  width: 30px;
  height: 30px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.2);
  font-size: 1.1rem;
}
.phase-body {
  flex: 1;
  padding: 12px;
  background: var(--sky-soft);
}
.phase-task {
  margin: 0 0 6px;
  font-weight: 700;
  font-size: 0.98rem !important;
  line-height: 1.5 !important;
  color: var(--ink) !important;
}
.phase ul {
  margin: 0;
  padding-left: 1.1em;
}
.phase li {
  font-size: 0.9rem;
  line-height: 1.6;
}
.phase-out {
  padding: 10px 12px 12px;
}
.out-label {
  margin: 0 0 2px;
  font-size: 0.8rem !important;
  font-weight: 700;
  letter-spacing: 0.1em;
  color: var(--weba) !important;
}
.phase-out li {
  color: var(--weba);
}
.loop {
  margin: 10px 0 0;
  padding: 10px 16px;
  border-radius: 10px;
  border: 1.5px dashed var(--weba);
  text-align: center;
}

.steps {
  list-style: none;
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 22px;
  margin: 12px 0 0;
  padding: 0;
  counter-reset: none;
}
.steps li {
  position: relative;
  border-radius: 12px;
  border: 1.5px solid var(--c);
  background: var(--ground-raised);
  overflow: hidden;
}
.steps li:not(:last-child)::after {
  content: '➜';
  position: absolute;
  right: -20px;
  top: 14px;
  color: var(--ink-faint);
  font-size: 1rem;
}
.steps li {
  overflow: visible;
}
.step-head {
  margin: 0 0 8px;
  padding: 10px;
  border-radius: 10px 10px 0 0;
  background: var(--c);
  color: #fff !important;
  font-weight: 700;
  text-align: center;
}
.step-line {
  margin: 0;
  padding: 0 10px;
  font-size: 0.92rem !important;
  line-height: 1.6 !important;
  text-align: center;
}
.step-line:last-child {
  padding-bottom: 12px;
}

.poap {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(0, 1.25fr);
  gap: 20px;
  align-items: center;
}
.poap-scene {
  display: block;
  width: 100%;
  border-radius: 12px;
}
.poap-list {
  margin: 14px 0 0;
  padding-left: 1.2em;
}
.poap-list li {
  color: #f3ede4;
}
.y {
  color: #ffe14d;
}
.poap-shots {
  display: grid;
  grid-template-columns: 0.42fr 1fr;
  gap: 12px;
  align-items: center;
}
.poap-shots img {
  display: block;
  width: 100%;
  border-radius: 12px;
}
.poap-shots .phone {
  border-radius: 18px;
  border: 4px solid #48413c;
}

.lifelong {
  text-align: center;
}
.flow {
  margin: 8px 0 4px;
  font-family: var(--wenkai);
  font-size: 2.2rem !important;
  font-weight: 700;
  color: var(--plum) !important;
}
.flow span {
  margin: 0 8px;
  color: var(--ink-faint);
}
.flow-notes {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 10px;
}
.flow-notes span {
  padding: 4px 14px;
  border-radius: 999px;
  background: var(--plum-soft);
  color: var(--plum);
  font-size: 0.95rem;
}

@media (max-width: 860px) {
  .addie {
    grid-template-columns: 1fr;
  }
  .steps {
    grid-template-columns: 1fr;
    gap: 26px;
  }
  .steps li:not(:last-child)::after {
    content: '↓';
    right: auto;
    left: 50%;
    top: auto;
    bottom: -24px;
    transform: translateX(-50%);
  }
}
@media (max-width: 680px) {
  .shots.two,
  .poap {
    grid-template-columns: 1fr;
  }
}
</style>
