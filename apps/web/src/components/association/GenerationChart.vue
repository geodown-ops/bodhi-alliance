<script setup lang="ts">
// 世代人口與占比（簡報資料來源：CBNData《2020 Z世代消費態度洞察報告》），半圓大小依人口數。
const base = 210
const generations = [
  { name: '嬰兒潮世代', years: '1946–1965', people: 11.7, share: '15%', color: '#a02a27', x: 105 },
  { name: 'X 世代', years: '1966–1980', people: 14.2, share: '18%', color: '#a02a27', x: 290 },
  { name: 'Y 世代', years: '1981–1994', people: 17.4, share: '22%', color: '#a02a27', x: 490 },
  { name: 'Z 世代', years: '1995–2009', people: 18.5, share: '24%', color: '#a02a27', x: 705 },
].map((g) => ({ ...g, r: Math.round((100 * g.people) / 18.5) }))

const arc = (x: number, r: number) => `M${x - r},${base} A${r},${r} 0 0 1 ${x + r},${base} Z`
</script>

<template>
  <div class="scroll">
    <svg viewBox="0 70 820 230" class="chart" role="img" aria-label="嬰兒潮世代 11.7 億人占 15%，X 世代 14.2 億人占 18%，Y 世代 17.4 億人占 22%，Z 世代 18.5 億人占 24%">
      <line x1="16" :y1="base" x2="804" :y2="base" stroke="#a02a27" stroke-width="3" />
      <path :d="`M796,${base - 6} L808,${base} L796,${base + 6}`" fill="none" stroke="#a02a27" stroke-width="3" />
      <g v-for="g in generations" :key="g.name" text-anchor="middle">
        <path :d="arc(g.x, g.r)" :fill="g.color" />
        <text :x="g.x" :y="base - g.r - 14" class="people">{{ g.people }} 億</text>
        <text :x="g.x" :y="base - g.r * 0.32" class="share">{{ g.share }}</text>
        <rect :x="g.x - 56" :y="base + 12" width="112" height="28" rx="6" fill="#f3e6c8" />
        <text :x="g.x" :y="base + 32" class="years">{{ g.years }}</text>
        <text :x="g.x" :y="base + 68" class="name" :fill="g.color">{{ g.name }}</text>
      </g>
    </svg>
  </div>
</template>

<style scoped>
.scroll {
  overflow-x: auto;
}
.chart {
  display: block;
  width: 100%;
  min-width: 520px;
}
.people {
  font-size: 22px;
  font-weight: 700;
  fill: #a02a27;
}
.share {
  font-size: 26px;
  fill: #f3e6c8;
  font-weight: 600;
}
.years {
  font-size: 16px;
  font-weight: 600;
  fill: #a02a27;
}
.name {
  font-size: 22px;
  font-weight: 700;
}
</style>
