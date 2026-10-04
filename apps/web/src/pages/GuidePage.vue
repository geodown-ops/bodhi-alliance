<script setup lang="ts">
import GuideStage from '../components/GuideStage.vue'

// 畫境佔滿頁首以下的整個畫面，頁尾留在捲動下方。
const pageHeight = () => ({ height: 'calc(100svh - 50px)' })
</script>

<template>
  <q-page :style-fn="pageHeight">
    <h1 class="sr-only">線上問答</h1>
    <GuideStage />
  </q-page>
</template>

<style scoped>
.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  overflow: hidden;
  clip: rect(0 0 0 0);
}
</style>
