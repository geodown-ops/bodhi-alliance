<script setup lang="ts">
// 第一章「協會總則」：簡介、任務、宗旨（簡報第 2–5 頁；簡介與任務文字以《協會簡介》為準）。
import PillarWheel from './PillarWheel.vue'
import { aims } from './aims'

const numerals = ['一', '二', '三', '四', '五']
// 每項任務：[前文, 重點, 後文]
const tasks = [
  ['推動', '慈悲關懷', '工作，包括人文關懷、社會救助、與護生及生態環保等。'],
  ['推動', '智慧無礙', '工作，包括雲端教學、佛學研究、生命教育、終身學習等。'],
  ['推動', '善巧方便', '工作，包括佛學資訊、音像紀錄，養生及身心健康、食品安全，以及社會型企業等。'],
  ['結合科技、醫藥、宗教、政府相關單位及民間社團組織等，推動', '終身學習及社會型企業策略聯盟', '。'],
  ['舉辦研討會、座談會、職業訓練、組織學習、關懷訪視等各項相關活動，以促進', '佛教普及弘化及世界永續發展', '之目的。'],
]
</script>

<template>
  <section id="general" class="chapter">
    <p class="kicker">協會總則</p>
    <h2>簡介・任務・宗旨</h2>
  </section>

  <h2>協會簡介</h2>
  <div class="intro-grid">
    <PillarWheel />
    <ul class="facts">
      <li>本會為依法設立、非以營利為目的之社會團體。</li>
      <li>本會成立於民國 111 年 12 月 30 日（送件申請）。</li>
      <li>
        本會以（1）<span class="hl">慈悲關懷</span>、（2）<span class="hl">智慧無礙</span>、（3）<span class="hl">善巧方便</span>為主軸，
        推展世界佛教雲端教學、佛學研究、社會教育、人文關懷、永續發展及社會型企業為宗旨。
      </li>
    </ul>
  </div>

  <h2>協會任務</h2>
  <p>本會之任務如下，並依相關法令規定推動及執行：</p>
  <ol class="tasks">
    <li v-for="(t, i) in tasks" :key="i">
      <span class="numeral">{{ numerals[i] }}</span>
      <span>{{ t[0] }}<span class="hl-sky task-key">{{ t[1] }}</span>{{ t[2] }}</span>
    </li>
  </ol>

  <h2>協會宗旨</h2>
  <div class="aim-bars">
    <p v-for="(a, i) in aims" :key="i" class="aim-bar" :class="a.tone">
      <template v-for="(part, j) in a.parts" :key="j">
        <b v-if="j % 2">{{ part }}</b>
        <span v-else>{{ part }}</span>
      </template>
    </p>
  </div>
</template>

<style scoped>
.intro-grid {
  display: grid;
  grid-template-columns: minmax(0, 2fr) minmax(0, 3fr);
  gap: 28px;
  align-items: center;
}
.facts {
  margin: 0;
  padding-left: 1.2em;
}
.facts li {
  margin: 6px 0;
}
.facts li::marker {
  color: var(--weba);
}
.tasks {
  list-style: none;
  margin: 0;
  padding: 0;
}
.tasks li {
  display: flex;
  gap: 14px;
  padding: 12px 4px;
  border-bottom: 1px dashed var(--rule);
}
.numeral {
  flex: none;
  display: grid;
  place-items: center;
  width: 32px;
  height: 32px;
  margin-top: 2px;
  border-radius: 50%;
  background: var(--sky-soft);
  color: var(--sky);
  font-weight: 700;
  font-size: 0.95rem;
}
.task-key {
  text-decoration: underline;
  text-underline-offset: 4px;
}
.aim-bars {
  display: grid;
  gap: 12px;
}
.aim-bar {
  margin: 0;
  padding: 18px 22px;
  border-radius: 12px;
  background: var(--night);
  color: #f2df8a !important;
  text-align: center;
  font-size: 1.08rem;
  line-height: 1.8;
}
.aim-bar.red {
  background: var(--weba);
}
.aim-bar b {
  color: #fff;
  font-size: 1.22rem;
}
.aim-bar.red b {
  color: #ffe36e;
}
@media (max-width: 680px) {
  .intro-grid {
    grid-template-columns: 1fr;
  }
  .aim-bar {
    text-align: left;
    padding: 14px 16px;
  }
}
</style>
