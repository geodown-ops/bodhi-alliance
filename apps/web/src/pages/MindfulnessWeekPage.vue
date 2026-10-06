<script setup lang="ts">
// 八週課程的單週頁：/mindfulness/week/:n
import { computed, watchEffect } from 'vue'
import { useRoute } from 'vue-router'
import RefList from '../components/RefList.vue'
import { dedication, weeks } from '../mindfulness'

const route = useRoute()
const n = computed(() => Math.min(8, Math.max(1, Number(route.params.n) || 1)))
const week = computed(() => weeks[n.value - 1])
const refs = computed(() => week.value.sections.flatMap((s) => s.refs ?? []))

watchEffect(() => {
  document.title = `第 ${n.value} 週 ${week.value.title}｜正念減壓｜Sunny life`
})
</script>

<template>
  <q-page class="page">
    <router-link to="/mindfulness" class="back">← 正念減壓八週課程</router-link>
    <div class="status-chip q-mt-md">第 {{ week.n }} 週 · {{ week.part }}</div>
    <h1 class="q-mt-sm">{{ week.title }}</h1>
    <p class="lead">{{ week.summary }}</p>

    <section v-for="s in week.sections" :key="s.title">
      <h2>{{ s.title }}</h2>
      <p v-for="p in s.paras ?? []" :key="p">{{ p }}</p>
      <ul v-if="s.items">
        <li v-for="i in s.items" :key="i">{{ i }}</li>
      </ul>
      <p v-if="s.note" class="small">{{ s.note }}</p>
    </section>

    <h2>跟著練</h2>
    <ol class="steps">
      <li v-for="s in week.steps" :key="s">{{ s }}</li>
    </ol>
    <p v-if="week.n === 4">
      坐姿與收功的完整說明，請看<router-link to="/mindfulness/sitting">上座與下座</router-link>。
    </p>

    <h2>這一週的生活練習</h2>
    <ul>
      <li>{{ week.life }}</li>
      <li>每天練習完，記下今天練了什麼、有什麼體會。</li>
    </ul>

    <h2>祈願與回向</h2>
    <div class="note dedication">
      <div v-for="line in dedication" :key="line">{{ line }}</div>
    </div>

    <div class="actions">
      <q-btn unelevated no-caps color="secondary" to="/groups" label="找覺行小組一起練" />
      <q-btn outline no-caps color="secondary" to="/guide" label="問問 Sunny" />
    </div>

    <nav class="pager">
      <router-link v-if="week.n > 1" :to="`/mindfulness/week/${week.n - 1}`">← 第 {{ week.n - 1 }} 週 {{ weeks[week.n - 2].title }}</router-link>
      <span v-else />
      <router-link v-if="week.n < 8" :to="`/mindfulness/week/${week.n + 1}`">第 {{ week.n + 1 }} 週 {{ weeks[week.n].title }} →</router-link>
    </nav>

    <template v-if="refs.length">
      <h2>出處</h2>
      <RefList :keys="refs" />
    </template>
    <p class="small q-mt-lg">本課程內容由法源法師審定編輯。</p>
  </q-page>
</template>

<style scoped>
.back,
.pager a {
  color: var(--leaf);
  text-decoration: none;
}
.small {
  font-size: 0.9rem !important;
  color: var(--ink-faint) !important;
}
.steps li {
  margin-bottom: 8px;
}
.dedication {
  font-family: var(--wenkai);
  font-size: 1.15rem;
  line-height: 2;
  color: var(--ink);
}
.actions {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 32px;
}
.pager {
  display: flex;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 32px;
  padding-top: 16px;
  border-top: 1px solid var(--rule);
}
</style>
