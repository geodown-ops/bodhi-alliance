<script setup lang="ts">
// 八堂課程的單堂頁：/mindfulness/lesson/:n
import { computed, watchEffect } from 'vue'
import { useRoute } from 'vue-router'
import LessonBar from '../components/LessonBar.vue'
import RefList from '../components/RefList.vue'
import { dedication, lessons } from '../mindfulness'

const route = useRoute()
const n = computed(() => Math.min(8, Math.max(1, Number(route.params.n) || 1)))
const lesson = computed(() => lessons[n.value - 1])
const refs = computed(() => lesson.value.sections.flatMap((s) => s.refs ?? []))

watchEffect(() => {
  document.title = `第 ${n.value} 堂 ${lesson.value.title}｜正念減壓｜Sunny life`
})
</script>

<template>
  <q-page class="page">
    <router-link to="/mindfulness" class="back">← 正念減壓八堂課程</router-link>
    <div class="q-mt-md"><span class="status-chip">第 {{ lesson.n }} 堂 · {{ lesson.part }}</span></div>
    <h1 class="q-mt-sm">{{ lesson.title }}</h1>
    <p class="lead">{{ lesson.summary }}</p>

    <section v-for="s in lesson.sections" :key="s.title">
      <h2>{{ s.title }}</h2>
      <p v-for="p in s.paras ?? []" :key="p">{{ p }}</p>
      <ul v-if="s.items">
        <li v-for="i in s.items" :key="i">{{ i }}</li>
      </ul>
      <p v-if="s.note" class="small">{{ s.note }}</p>
    </section>

    <h2>跟著練</h2>
    <ol class="steps">
      <li v-for="s in lesson.steps" :key="s">{{ s }}</li>
    </ol>
    <p v-if="lesson.n === 4">
      坐姿與收功的完整說明，請看<router-link to="/mindfulness/sitting">上座與下座</router-link>。
    </p>

    <h2>這一堂的生活練習</h2>
    <ul>
      <li>{{ lesson.life }}</li>
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
      <router-link v-if="lesson.n > 1" :to="`/mindfulness/lesson/${lesson.n - 1}`">← 第 {{ lesson.n - 1 }} 堂 {{ lessons[lesson.n - 2].title }}</router-link>
      <span v-else />
      <router-link v-if="lesson.n < 8" :to="`/mindfulness/lesson/${lesson.n + 1}`">第 {{ lesson.n + 1 }} 堂 {{ lessons[lesson.n].title }} →</router-link>
    </nav>

    <template v-if="refs.length">
      <h2>出處</h2>
      <RefList :keys="refs" />
    </template>
    <p class="small q-mt-lg">本課程內容由法源法師審定編輯。</p>
    <LessonBar />
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
  font-family: var(--sans);
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
