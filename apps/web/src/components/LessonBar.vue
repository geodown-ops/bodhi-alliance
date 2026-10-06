<script setup lang="ts">
// 正念減壓頁面底部固定浮動的課堂選單，方便在各堂課之間快速切換
import { nextTick, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { lessons } from '../mindfulness'

// 手機上選單會橫向捲動，切換課堂時把目前這一堂捲到中間
const inner = ref<HTMLElement>()
const route = useRoute()
const center = () =>
  nextTick(() => {
    const box = inner.value
    const on = box?.querySelector<HTMLElement>('.on')
    if (box && on) box.scrollLeft = on.offsetLeft - (box.clientWidth - on.offsetWidth) / 2
  })
onMounted(center)
watch(() => route.path, center)
</script>

<template>
  <nav class="lesson-bar" aria-label="課堂選單">
    <div ref="inner" class="inner">
      <router-link to="/mindfulness" exact-active-class="on" class="item">總覽</router-link>
      <router-link
        v-for="l in lessons"
        :key="l.n"
        :to="`/mindfulness/lesson/${l.n}`"
        active-class="on"
        class="item"
        :title="`第 ${l.n} 堂 ${l.title}`"
      >
        <span class="num">第 {{ l.n }} 堂</span>
        <span class="name">{{ l.title }}</span>
      </router-link>
      <router-link to="/mindfulness/sitting" active-class="on" class="item">上座與下座</router-link>
    </div>
  </nav>
  <div class="spacer" />
</template>

<style scoped>
.lesson-bar {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 12px;
  z-index: 2000;
  display: flex;
  justify-content: center;
  padding: 0 12px;
  pointer-events: none;
}
.inner {
  position: relative;
  pointer-events: auto;
  display: flex;
  gap: 4px;
  max-width: 100%;
  overflow-x: auto;
  scrollbar-width: none;
  padding: 6px;
  background: color-mix(in srgb, var(--ground-raised) 92%, transparent);
  border: 1px solid var(--rule);
  border-radius: 999px;
  box-shadow: 0 6px 20px rgb(59 42 32 / 0.16);
  backdrop-filter: blur(8px);
}
.inner::-webkit-scrollbar {
  display: none;
}
.item {
  flex: none;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 40px;
  padding: 4px 12px;
  border-radius: 999px;
  color: var(--ink-soft);
  text-decoration: none;
  font-size: 0.9rem;
  line-height: 1.25;
  white-space: nowrap;
}
.item .name {
  font-size: 0.75rem;
  color: var(--ink-faint);
}
.item:hover {
  background: var(--ground-sunk);
}
.item.on {
  background: var(--leaf);
  color: #fff;
}
.item.on .name {
  color: rgb(255 255 255 / 0.85);
}
.spacer {
  height: 80px;
}
@media (max-width: 1100px) {
  .item .name {
    display: none;
  }
}
</style>
