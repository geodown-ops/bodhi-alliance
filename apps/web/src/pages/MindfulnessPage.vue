<script setup lang="ts">
// 正念減壓介紹。內容依正念減壓（MBSR）的一般說明與覺行小組定義整理，協會的正念課程講義定稿後再充實。
const attitudes = [
  { name: '不評判', text: '只是看見，不急著貼上好壞、對錯的標籤。' },
  { name: '耐心', text: '讓事情按它自己的速度展開，身心也一樣。' },
  { name: '初心', text: '把每一次呼吸都當作第一次，帶著好奇去看。' },
  { name: '信任', text: '相信自己的感受與判斷，練習是自己的。' },
  { name: '不強求', text: '不為了「放鬆」而練，放鬆是自然的結果。' },
  { name: '接納', text: '承認事情此刻就是這個樣子，才有轉圜的空間。' },
  { name: '放下', text: '不緊抓喜歡的，也不推開不喜歡的。' },
]

const practices = [
  { name: '呼吸覺察', time: '5–20 分鐘', text: '坐穩，把注意力放在呼吸進出的感覺上。心跑掉了就溫和地帶回來，帶回來本身就是練習。' },
  { name: '身體掃描', time: '20–40 分鐘', text: '躺下或坐著，從腳趾到頭頂，依序覺察身體各部位的感受，不需要改變什麼。' },
  { name: '正念伸展', time: '10–30 分鐘', text: '緩慢、溫和的伸展與瑜伽動作，覺察肌肉的拉伸與界限，不勉強。' },
  { name: '正念行走', time: '10–15 分鐘', text: '放慢腳步，感受腳掌抬起、移動、落地，把走路本身當作練習。' },
  { name: '三分鐘呼吸空間', time: '3 分鐘', text: '覺察此刻的念頭與感受、收攝到呼吸、再擴展到整個身體。壓力上來時隨時可用。' },
  { name: '日常正念', time: '隨時', text: '吃飯時好好吃飯、聽人說話時好好聽，把覺察帶進每天重複的小事。' },
]

const flow = [
  { title: '安頓', text: '簽到、入座，用幾分鐘從忙碌裡停下來，說說今天的狀態。' },
  { title: '引導練習', text: '由減壓教練或志工引導，一次做一到兩個練習，例如呼吸覺察與身體掃描。' },
  { title: '分享', text: '每個人說說剛才的體會，只說自己的經驗，不評論、不給建議。' },
  { title: '帶回日常', text: '約好這一週要練的一個小練習，下次聚會時再分享。' },
  { title: '登錄活動', text: '活動在 Sunny life 官網登錄後，協助的志工或減壓教練可以得到菩提幣。' },
]
</script>

<template>
  <q-page class="page">
    <h1>正念減壓</h1>
    <p class="lead">
      正念減壓是用正念的練習，學會與壓力相處。練的是呼吸、身體與注意力，改變的是我們面對壓力的方式。覺行小組一起練的，就是正念減壓。
    </p>

    <h2>什麼是正念</h2>
    <p>
      正念，是有意識地、不加評判地，把注意力放在當下這一刻。壓力常常來自對過去的反芻與對未來的擔憂；正念把我們帶回此時此地，看清身體的感受、情緒與念頭，從「自動反應」轉為「有意識的回應」。
    </p>

    <h2>正念減壓從哪裡來</h2>
    <p>
      正念減壓（Mindfulness-Based Stress Reduction，簡稱 MBSR）由卡巴金博士（Jon Kabat-Zinn）於 1979 年在美國麻州大學醫學院創立，原本是協助慢性疼痛與壓力病人的八週課程。四十多年來，它在醫院、學校與企業裡廣泛運用，許多研究顯示，持續練習有助於減輕壓力與焦慮、改善睡眠、提升專注與情緒的穩定。
    </p>
    <p>
      在 Sunny life，正念減壓的實作原理來自世界佛教教育協會的<router-link to="/maitreya">彌勒心流</router-link>實修：放下評判是轉念，安住呼吸是轉識，把覺知帶回生活與團體是轉依。
    </p>

    <h2>七個基本態度</h2>
    <div class="grid">
      <div v-for="a in attitudes" :key="a.name" class="card">
        <h3 class="q-mt-none q-mb-xs">{{ a.name }}</h3>
        <p class="q-mb-none">{{ a.text }}</p>
      </div>
    </div>

    <h2>常用的練習</h2>
    <div class="grid">
      <div v-for="p in practices" :key="p.name" class="card">
        <div class="status-chip">{{ p.time }}</div>
        <h3 class="q-my-sm">{{ p.name }}</h3>
        <p class="q-mb-none">{{ p.text }}</p>
      </div>
    </div>

    <h2>在覺行小組一起練</h2>
    <p>
      <router-link to="/groups">覺行小組</router-link>是隨興或定期相約，一起進行正念減壓的小組活動。任何人都可以發起，三個人以上就能進行一次，在中心、在家或線上都可以。一次聚會大約一到兩小時：
    </p>
    <ol class="flow">
      <li v-for="f in flow" :key="f.title">
        <strong>{{ f.title }}</strong>
        <span>{{ f.text }}</span>
      </li>
    </ol>

    <h2>練習時請留意</h2>
    <ul>
      <li>不需要任何經驗，也不需要特定的坐姿；坐椅子、躺著都可以，舒服而清醒就好。</li>
      <li>練習中感到不適，可以睜開眼睛、調整姿勢或停下來，照顧好自己比完成練習重要。</li>
      <li>正念減壓不能取代醫療或心理治療。正在接受治療、或有創傷經驗與嚴重身心狀況的人，請先和醫師或專業人員討論。</li>
    </ul>

    <div class="note q-mt-lg">
      想試試看？找一個<router-link to="/groups">覺行小組</router-link>，或先問問<router-link to="/guide">線上問答</router-link>的 Sunny。
    </div>
  </q-page>
</template>

<style scoped>
.flow {
  list-style: none;
  counter-reset: step;
  padding: 0;
  margin: 16px 0 0;
}
.flow li {
  counter-increment: step;
  display: grid;
  grid-template-columns: 2.2rem 7rem minmax(0, 1fr);
  gap: 4px 12px;
  padding: 12px 0;
  border-bottom: 1px solid var(--rule);
}
.flow li::before {
  content: counter(step);
  font-family: var(--wenkai);
  font-weight: 700;
  font-size: 1.3rem;
  color: var(--leaf);
}
.flow strong {
  color: var(--ink);
  font-weight: 600;
}
@media (max-width: 600px) {
  .flow li {
    grid-template-columns: 2.2rem minmax(0, 1fr);
  }
  .flow span {
    grid-column: 2;
  }
}
</style>
