<script setup lang="ts">
// 彌勒心流：協會的實修。唯識三轉原本在首頁「宗旨」段，移到這裡；下半頁說明正念減壓的每一步怎麼對應三轉。
// 彌勒心流本身的正式介紹（傳承、指導法師、課程）待協會提供。
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
</script>

<template>
  <q-page class="page">
    <h1>彌勒心流</h1>
    <p class="lead">彌勒心流（Maitreya Flow）是世界佛教教育協會的靜坐實修。覺行小組所練的正念減壓，實作原理就來自這裡。</p>

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

    <h2>一起練</h2>
    <p>
      三個人以上，就可以發起一次<router-link to="/groups">覺行小組</router-link>的正念減壓實作。對練習有疑問，可以問問<router-link to="/guide">線上問答</router-link>的 Sunny。
    </p>
    <p class="text-caption">彌勒心流的傳承、指導法師與課程介紹，將依協會提供的資料補上。</p>
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
@media (max-width: 700px) {
  .triad {
    grid-template-columns: 1fr;
    gap: 24px 0;
  }
}
</style>
