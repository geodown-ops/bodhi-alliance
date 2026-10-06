<script setup lang="ts">
// 正念減壓：八週科技正念減壓課程總覽。內容見 ../mindfulness.ts，法源法師審定編輯。
import RefList from '../components/RefList.vue'
import { attitudes, dedication, parts, weeks } from '../mindfulness'

const weekOf = (n: number) => weeks[n - 1]

const flow = [
  { title: '安頓', text: '簽到、入座，用幾分鐘從忙碌裡停下來。' },
  { title: '上座', text: '調身、調息、調心，坐穩、呼吸放慢。' },
  { title: '本週練習', text: '由減壓教練或志工帶這一週的練習。' },
  { title: '分享', text: '每個人說說剛才的體會，只說自己的經驗，不評論、不給建議。' },
  { title: '下座收功', text: '出靜、深呼吸，做收功動作再起身。' },
  { title: '祈願回向', text: '一起念祈願文，約好這一週的生活練習。' },
  { title: '登錄活動', text: '活動在 Sunny life 官網登錄後，協助的志工或減壓教練可以得到菩提幣。' },
]
</script>

<template>
  <q-page class="page">
    <h1>正念減壓</h1>
    <p class="lead">
      正念減壓是用正念的練習，學會與壓力相處。Sunny life 的正念減壓是一套八週的科技正念減壓課程，依佛法的四念處，從身體、感受、心念一路練到覺知無常；覺行小組一起練的，就是這套課程。
    </p>
    <p class="credit">本課程內容由法源法師審定編輯。</p>

    <h2>什麼是正念</h2>
    <p>
      「念」拆開是「今＋心」，是此刻心中的念頭、想法、感受與情緒；「正」是當下、不偏不倚，不是「正向」的意思。正念，就是刻意專注於當下的覺察，暫時收起判斷心，帶著好奇與接納，一刻接著一刻。
    </p>
    <p>
      正念減壓（MBSR）由卡巴金博士（Jon Kabat-Zinn）於 1979 年在美國麻州大學醫學院創立，原本是協助慢性疼痛與壓力病人的八週課程。許多研究整理顯示，它能明顯減輕壓力，也有助於緩解焦慮、憂鬱與提升生活品質。
    </p>
    <p>
      這套課程把正念減壓和佛法的安般念、四念處與六妙門對照起來，並用心律與腦波量測看見練習的效果，所以叫「科技正念」。它是<router-link to="/maitreya">彌勒心流</router-link>正念、正定、正知三階修學的第一階。
    </p>

    <h2>九個正念態度</h2>
    <div class="grid">
      <div v-for="a in attitudes" :key="a.name" class="card">
        <h3 class="q-mt-none q-mb-xs">{{ a.name }}</h3>
        <div class="en">{{ a.en }}</div>
        <p class="q-mb-none">{{ a.text }}</p>
      </div>
    </div>

    <h2>八週課程</h2>
    <p>八週分成四段，每段兩週。每週有一個主題、一套練習，以及每天在生活裡做的小練習。</p>
    <div v-for="p in parts" :key="p.name" class="part">
      <div class="part-head">
        <h3 class="q-my-none">{{ p.name }}</h3>
        <span>{{ p.text }}</span>
      </div>
      <div class="grid">
        <router-link v-for="n in p.weeks" :key="n" :to="`/mindfulness/week/${n}`" class="card week">
          <div class="status-chip">第 {{ n }} 週</div>
          <h3 class="q-my-sm">{{ weekOf(n).title }}</h3>
          <p class="q-mb-none">{{ weekOf(n).summary }}</p>
        </router-link>
      </div>
    </div>

    <h2>上座與下座</h2>
    <p>
      每次靜坐都從調身、調息、調心開始，結束時再反過來調心、調息、調身，做一套收功動作才起身。詳細步驟請看<router-link to="/mindfulness/sitting">上座與下座</router-link>。
    </p>

    <h2>在覺行小組一起練</h2>
    <p>
      <router-link to="/groups">覺行小組</router-link>是隨興或定期相約，一起進行正念減壓的小組活動。任何人都可以發起，三個人以上就能進行一次，在中心、在家或線上都可以。可以照著八週的順序一週練一個主題，一次聚會大約一到兩小時：
    </p>
    <ol class="flow">
      <li v-for="f in flow" :key="f.title">
        <strong>{{ f.title }}</strong>
        <span>{{ f.text }}</span>
      </li>
    </ol>

    <h2>祈願與回向</h2>
    <p>每次練習結束，一起念：</p>
    <div class="note dedication">
      <div v-for="line in dedication" :key="line">{{ line }}</div>
    </div>

    <h2>練習時請留意</h2>
    <ul>
      <li>不需要任何經驗，也不需要特定的坐姿；坐椅子、躺著都可以，舒服而清醒就好。</li>
      <li>練習中感到不適，可以睜開眼睛、調整姿勢或停下來，照顧好自己比完成練習重要。</li>
      <li>心律與腦波量測只用來觀察自己的練習，不是醫療檢查。</li>
      <li>正念減壓不能取代醫療或心理治療。正在接受治療、或有創傷經驗與嚴重身心狀況的人，請先和醫師或專業人員討論。</li>
    </ul>

    <div class="note q-mt-lg">
      想試試看？從<router-link to="/mindfulness/week/1">第一週</router-link>開始，找一個<router-link to="/groups">覺行小組</router-link>，或先問問<router-link to="/guide">線上問答</router-link>的 Sunny。
    </div>

    <h2>出處</h2>
    <RefList :keys="['kz1982', 'kz1990', 'kz2003', 'goyal2014', 'khoury2015', 'satipatthana', 'anapana']" />
  </q-page>
</template>

<style scoped>
.credit {
  font-size: 0.95rem !important;
  color: var(--ink-faint) !important;
}
.en {
  font-size: 0.85rem;
  color: var(--ink-faint);
  margin-bottom: 6px;
}
.part {
  margin-top: 24px;
}
.part-head {
  display: flex;
  flex-wrap: wrap;
  align-items: baseline;
  gap: 4px 12px;
  margin-bottom: 12px;
}
.part-head h3 {
  font-family: var(--wenkai);
  color: var(--leaf);
}
.part-head span {
  color: var(--ink-soft);
}
.week {
  display: block;
  color: inherit;
  text-decoration: none;
  transition: border-color 0.15s;
}
.week:hover,
.week:focus-visible {
  border-color: var(--leaf);
}
.dedication {
  font-family: var(--wenkai);
  font-size: 1.15rem;
  line-height: 2;
  color: var(--ink);
}
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
