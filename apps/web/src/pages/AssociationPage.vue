<script setup lang="ts">
// 世界佛教教育協會介紹：完整呈現協會簡報《世界佛教教育協會介紹》36 頁的內容，簡介與任務的文字以協會提供的《協會簡介》為準。
// 圖表依簡報數據重畫（components/association/）；照片與標誌取自簡報原檔，放在 public/images/association/。
import '../components/association/association.css'
import AssocGeneral from '../components/association/AssocGeneral.vue'
import AssocMembers from '../components/association/AssocMembers.vue'
import AssocOperation from '../components/association/AssocOperation.vue'
import AssocOrigin from '../components/association/AssocOrigin.vue'

const chapters = [
  { id: 'general', label: '協會總則' },
  { id: 'origin', label: '協會起源' },
  { id: 'operation', label: '核心共識' },
  { id: 'members', label: '會員招募' },
]
const go = (id: string) => document.getElementById(id)?.scrollIntoView({ behavior: 'smooth' })
</script>

<template>
  <q-page class="page page-wide assoc">
    <header class="hero">
      <img src="/images/association/logo.webp" alt="世界佛教教育協會標誌" class="logo" />
      <div>
        <h1>世界佛教教育協會</h1>
        <p class="lead q-mb-none">以慈悲關懷、智慧無礙、善巧方便為主軸，推展世界佛教雲端教學、佛學研究、社會教育、人文關懷、永續發展及社會型企業。</p>
      </div>
    </header>
    <nav class="toc" aria-label="章節">
      <button v-for="c in chapters" :key="c.id" type="button" @click="go(c.id)">{{ c.label }}</button>
    </nav>

    <AssocGeneral />
    <AssocOrigin />
    <AssocOperation />
    <AssocMembers />
  </q-page>
</template>

<style scoped>
.hero {
  display: flex;
  gap: 22px;
  align-items: center;
}
.logo {
  flex: none;
  width: 104px;
  height: 104px;
  border-radius: 50%;
  box-shadow: 0 8px 22px rgba(186, 88, 84, 0.25);
}
.toc {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin: 20px 0 0;
}
.toc button {
  padding: 6px 16px;
  border: 1px solid var(--rule);
  border-radius: 999px;
  background: var(--ground-raised);
  color: var(--ink-soft);
  font: inherit;
  cursor: pointer;
}
.toc button:hover {
  border-color: var(--weba);
  color: var(--weba);
}
@media (max-width: 600px) {
  .hero {
    flex-direction: column;
    align-items: flex-start;
  }
  .logo {
    width: 76px;
    height: 76px;
  }
}
</style>
