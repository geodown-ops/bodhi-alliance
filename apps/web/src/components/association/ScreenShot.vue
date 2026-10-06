<script setup lang="ts">
// 系統畫面截圖：加上瀏覽器外框，點一下開原圖。
defineProps<{ src: string; alt: string; caption?: string }>()
</script>

<template>
  <figure class="shot">
    <a :href="src" target="_blank" rel="noopener" :aria-label="`開啟大圖：${alt}`">
      <div class="bar"><i /><i /><i /></div>
      <img :src="src" :alt="alt" loading="lazy" />
    </a>
    <figcaption v-if="caption">{{ caption }}</figcaption>
  </figure>
</template>

<style scoped>
.shot {
  margin: 0;
}
.shot a {
  display: block;
  overflow: hidden;
  border-radius: 12px;
  border: 1px solid var(--rule);
  background: #fff;
  box-shadow: 0 10px 28px rgba(59, 42, 32, 0.12);
  transition: transform 0.2s ease;
}
.shot a:hover {
  transform: translateY(-2px);
}
.bar {
  display: flex;
  gap: 6px;
  padding: 8px 12px;
  background: var(--ground-sunk);
}
.bar i {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: #cfc4b3;
}
.bar i:first-child {
  background: #e08a7f;
}
.shot img {
  display: block;
  width: 100%;
  height: auto;
}
figcaption {
  margin-top: 10px;
  text-align: center;
  font-weight: 600;
  color: var(--plum);
}
</style>
