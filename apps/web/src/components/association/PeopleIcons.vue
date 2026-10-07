<script setup lang="ts">
// 人形圖示列：total 個人裡有 filled 個上色，用來畫「三分之一」「四分之三」這類比例。
defineProps<{ total: number; filled: number; color: string }>()
</script>

<template>
  <svg :viewBox="`0 0 ${total * 44} 60`" class="people" aria-hidden="true">
    <g v-for="i in total" :key="i" :transform="`translate(${(i - 1) * 44 + 4},0)`" :fill="i <= filled ? color : '#f1d9d4'">
      <circle cx="18" cy="11" r="10" />
      <path d="M2,58 C2,36 8,25 18,25 C28,25 34,36 34,58 Z" />
    </g>
  </svg>
</template>

<style scoped>
.people {
  display: block;
  height: 48px;
  margin: 0 auto 10px;
}
</style>
