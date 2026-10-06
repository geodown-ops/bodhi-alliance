<script setup lang="ts">
// 上座與下座：依法源法師〈止觀禪修上下坐引導〉整理
import LessonBar from '../components/LessonBar.vue'
import { sitting } from '../mindfulness'
</script>

<template>
  <q-page class="page">
    <router-link to="/mindfulness" class="back">← 正念減壓八堂課程</router-link>
    <h1 class="q-mt-md">上座與下座</h1>
    <p class="lead">
      靜坐開始時調身、調息、調心，坐中保持住；結束時倒過來，調心、調息、調身，做完收功動作才起身。
    </p>
    <p v-for="b in sitting.before" :key="b">{{ b }}</p>

    <div class="table-wrap">
      <table class="table">
        <thead>
          <tr><th /><th>上座（入）</th><th>坐中（住）</th><th>下座（出）</th></tr>
        </thead>
        <tbody>
          <tr v-for="r in sitting.table" :key="r.k">
            <th>{{ r.k }}</th><td>{{ r.up }}</td><td>{{ r.mid }}</td><td>{{ r.down }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <h2>上座</h2>
    <section v-for="s in sitting.up" :key="s.title">
      <h3>{{ s.title }}</h3>
      <ul>
        <li v-for="i in s.items" :key="i">{{ i }}</li>
      </ul>
    </section>

    <h2>下座</h2>
    <section v-for="s in sitting.down" :key="s.title">
      <h3>{{ s.title }}</h3>
      <ul>
        <li v-for="i in s.items" :key="i">{{ i }}</li>
      </ul>
    </section>

    <div class="note q-mt-lg">
      身體有舊傷或不適，動作做到舒服的程度就好，不需要勉強。
    </div>
    <p class="small q-mt-lg">本頁依法源法師的上下座引導整理，法源法師審定編輯。</p>
    <LessonBar />
  </q-page>
</template>

<style scoped>
.back {
  color: var(--leaf);
  text-decoration: none;
}
.small {
  font-size: 0.9rem !important;
  color: var(--ink-faint) !important;
}
.table-wrap {
  overflow-x: auto;
  margin-top: 16px;
}
.table th {
  white-space: nowrap;
}
</style>
