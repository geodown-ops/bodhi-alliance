<script setup lang="ts">
// 彌勒心流：協會的實修。唯識三轉原本在首頁「宗旨」段，移到這裡；下半頁說明正念減壓的每一步怎麼對應三轉。
// 修學次第依《正念課程主軸與應用》《課程特色與核心經典》課程藍圖，法源法師審定編輯。
import RefList from '../components/RefList.vue'
const turns = [
  { name: '轉念', romaji: 'zhuǎn niàn', text: '放下壓力，回歸中道與平衡。有持續力的服務，是流動的是可以回饋自己的。' },
  { name: '轉識', romaji: 'zhuǎn shí', text: '透過呼吸覺察，安定身心。看清資源與付出真正流向哪裡。' },
  { name: '轉依', romaji: 'zhuǎn yī', text: '觀想與身心覺知統合，進入與法界合一。一間中心的量能，成為所有共好企業的量能。' },
]

const practice = [
  {
    turn: '轉念',
    mbsr: '停下來，放下評判',
    text: '壓力來時，先停下自動反應，承認「我現在很緊繃」，不急著批評自己或別人。念頭一轉，就從被壓力推著走，回到中道與平衡。',
  },
  {
    turn: '轉識',
    mbsr: '安住呼吸，看清當下',
    text: '以呼吸為錨，做呼吸覺察與身體掃描。心安定下來，才看得清身體的感受、情緒與念頭，而不是被它們帶走。',
  },
  {
    turn: '轉依',
    mbsr: '把覺知帶回生活與團體',
    text: '把練習帶進吃飯、走路、工作與服務，也和覺行小組的夥伴一起練。個人的安定，成為身邊的人、所有共好企業的安定。',
  },
]

// 三階課程 × 四法（課程藍圖）
const methods = [
  { name: '對治法', sutra: '《阿含經》' },
  { name: '轉念法', sutra: '《深密經》' },
  { name: '破相法', sutra: '《金剛經》' },
  { name: '顯性法', sutra: '《華嚴經》' },
]
const stages = [
  { name: '正念課程', link: '/mindfulness', cells: ['正念呼吸', '情緒座標', '觀想呼吸', '能量觀想'] },
  { name: '正定課程', cells: ['十六勝行', '瑜伽止觀', '金剛三昧', '華嚴三昧'] },
  { name: '正知課程', cells: ['十六觀智', '轉識成智', '破相見性', '圓顯佛性'] },
]

// 四念處與安般念、六妙門、止觀階段的對照
const foundations = [
  { name: '身念住', breath: '覺知息長、息短、全身息，安息身行', gate: '數、隨', stage: '念而無念', effect: '紓壓，對治粗心與昏沉散亂', weeks: [1, 2] },
  { name: '受念住', breath: '覺知喜、覺知樂、覺知心行，安息心行', gate: '止', stage: '入定，初禪、二禪', effect: '心安住', weeks: [3, 4] },
  { name: '心念住', breath: '覺知心、令心喜悅、令心等持、令心解脫', gate: '觀', stage: '三禪、四禪', effect: '對治細心', weeks: [5, 6] },
  { name: '法念住', breath: '觀無常、觀離欲、觀滅、觀捨遣', gate: '還、淨', stage: '止觀雙運', effect: '定慧等持，增上智慧', weeks: [7, 8] },
]
</script>

<template>
  <q-page class="page">
    <h1>彌勒心流</h1>
    <p class="lead">彌勒心流（Maitreya Flow）是世界佛教教育協會的止觀靜坐實修，由法源法師指導，依正念、正定、正知三階次第修學。覺行小組所練的正念減壓，就是第一階的正念課程。</p>

    <h2>實修唯識三轉</h2>
    <p>世界佛教教育協會的修行核心，是把覺察落在日常裡：轉念、轉識、轉依。Sunny life 是這三轉在「共同生活」這一層的實作——讓服務與受用在同一個善行裡循環再擴大。</p>
    <div class="triad">
      <div v-for="t in turns" :key="t.name">
        <div class="glyph">{{ t.name }}</div>
        <div class="romaji">{{ t.romaji }}</div>
        <p>{{ t.text }}</p>
      </div>
    </div>

    <h2>正念減壓的實作原理</h2>
    <p>
      <router-link to="/mindfulness">正念減壓</router-link>的每一步，都可以在三轉裡找到它的位置。練的是呼吸與身體，轉的是看待壓力與生活的方式。
    </p>
    <div class="steps">
      <div v-for="p in practice" :key="p.turn" class="card">
        <div class="status-chip">{{ p.turn }}</div>
        <h3 class="q-my-sm">{{ p.mbsr }}</h3>
        <p class="q-mb-none">{{ p.text }}</p>
      </div>
    </div>

    <h2>修學次第</h2>
    <p>
      彌勒心流分三階：先以正念課程紓解壓力、安定身心，再以正定課程深入禪定，最後以正知課程開啟智慧。每一階都以四部經典為根據，依序用對治、轉念、破相、顯性四種方法練習。其中「轉念法」依《深密經》，正知課程的「轉識成智」，正是唯識三轉的核心。
    </p>
    <div class="table-wrap">
      <table class="table">
        <thead>
          <tr>
            <th />
            <th v-for="m in methods" :key="m.name">{{ m.name }}<div class="sutra">{{ m.sutra }}</div></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="s in stages" :key="s.name">
            <th>
              <router-link v-if="s.link" :to="s.link">{{ s.name }}</router-link>
              <span v-else>{{ s.name }}</span>
            </th>
            <td v-for="c in s.cells" :key="c">{{ c }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <h3>四念處：正念課程的骨架</h3>
    <p>
      正念課程八週依四念處展開，並對照南傳《清淨道論》的安般念、天台智者大師的六妙門（數、隨、止、觀、還、淨）與止觀階段：
    </p>
    <div class="steps four">
      <div v-for="f in foundations" :key="f.name" class="card">
        <div class="status-chip">六妙門：{{ f.gate }}</div>
        <h3 class="q-my-sm">{{ f.name }}</h3>
        <p class="q-mb-xs">{{ f.breath }}</p>
        <p class="q-mb-xs small">止觀：{{ f.stage }}　效果：{{ f.effect }}</p>
        <p class="q-mb-none small">
          <router-link v-for="n in f.weeks" :key="n" :to="`/mindfulness/week/${n}`" class="wk">第 {{ n }} 週</router-link>
        </p>
      </div>
    </div>

    <h2>一起練</h2>
    <p>
      三個人以上，就可以發起一次<router-link to="/groups">覺行小組</router-link>的正念減壓實作。對練習有疑問，可以問問<router-link to="/guide">線上問答</router-link>的 Sunny。
    </p>
    <p class="text-caption">正定課程與正知課程的詳細內容，將陸續補上。修學次第由法源法師審定編輯。</p>

    <h2>出處</h2>
    <RefList :keys="['satipatthana', 'anapana', 'visuddhimagga', 'zhiyi']" />
  </q-page>
</template>

<style scoped>
.triad {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 0 32px;
  margin-top: 28px;
}
.triad > div {
  padding-top: 18px;
  border-top: 2px solid var(--leaf);
}
.glyph {
  font-family: var(--wenkai);
  font-weight: 700;
  font-size: 2.4rem;
  line-height: 1.1;
  letter-spacing: 0.06em;
  color: var(--ink);
}
.romaji {
  font-size: 0.75rem;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: var(--ink-faint);
  margin-top: 8px;
}
.triad p {
  margin-top: 10px;
}
.steps {
  display: grid;
  gap: 16px;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  margin-top: 20px;
}
.steps.four {
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
}
.table-wrap {
  overflow-x: auto;
  margin-top: 16px;
}
.table th {
  white-space: nowrap;
  vertical-align: top;
}
.sutra {
  font-weight: 400;
  font-size: 0.85rem;
  color: var(--ink-faint);
}
.small {
  font-size: 0.92rem !important;
}
.wk {
  color: var(--leaf);
  margin-right: 12px;
}
@media (max-width: 700px) {
  .triad {
    grid-template-columns: 1fr;
    gap: 24px 0;
  }
}
</style>
