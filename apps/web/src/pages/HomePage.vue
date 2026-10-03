<script setup lang="ts">
const sections = [
  { to: '/coin', title: '菩提幣介紹', text: '服務時數換成聯盟共通的記帳單位，在共好企業兌換住宿、餐飲與課程。' },
  { to: '/groups', title: '覺行小組介紹及參加', text: '找一個離你近的覺行小組，一起共修，也一起服務。' },
  { to: '/partners', title: '共好企業登記及管理', text: '把淡季的空房、餐飲與課程名額，回饋給服務的人。' },
  { to: '/wallet', title: '菩提幣錢包', text: '查看餘額、兌換券、出示身份 QR 核銷。' },
  { to: '/committee', title: '主辦審核小組', text: '制定梯級表、審核核發名單，所有標準都有決議可查。' },
  { to: '/association', title: '世界佛教教育協會', text: '實修唯識三轉：轉念、轉識、轉依。' },
  { to: '/guide', title: '線上覺行小組 AI 組長', text: '隨時問覺行與共修的問題，或請組長帶你共修一段。' },
]
</script>

<template>
  <q-page>
    <section class="hero">
      <div class="page">
        <p class="kicker">世界佛教教育協會 · Sunny life</p>
        <h1>一即一切，一切即一</h1>
        <p class="lead">一個人的服務，回到所有人身上；所有人的供養，回到每一個人身上。</p>
        <div class="row q-gutter-sm q-mt-md">
          <q-btn color="secondary" unelevated no-caps size="lg" to="/groups" label="參加覺行小組" />
          <q-btn outline color="primary" no-caps size="lg" to="/coin" label="認識菩提幣" />
        </div>
      </div>
    </section>
    <div class="page">
      <div class="grid">
        <router-link v-for="s in sections" :key="s.to" :to="s.to" class="card tile">
          <h3>{{ s.title }}</h3>
          <p>{{ s.text }}</p>
        </router-link>
      </div>
    </div>
  </q-page>
</template>

<style scoped>
.hero {
  background: var(--ground-raised);
  border-bottom: 1px solid var(--rule);
}
.hero h1 {
  font-size: clamp(2rem, 6vw, 3.2rem);
}
.kicker {
  color: var(--leaf) !important;
  letter-spacing: 0.1em;
  margin: 0;
}
.tile {
  text-decoration: none;
  color: inherit;
  transition: border-color 0.15s;
}
.tile:hover {
  border-color: var(--leaf);
}
.tile h3 {
  font-family: var(--wenkai);
  font-size: 1.3rem;
  margin: 0 0 6px;
  line-height: 1.4;
}
.tile p {
  margin: 0;
  font-size: 0.98rem;
}
</style>
