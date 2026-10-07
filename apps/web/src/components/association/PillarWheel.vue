<script setup lang="ts">
// 協會三主軸：慈悲關懷、智慧無礙、善巧方便，中間是協會標誌。
const sectors = [
  { name: '慈悲關懷', no: '1', d: 'M150,150 L28.76,220 A140,140 0 0 1 150,10 Z', fill: '#f1d9d4', x: 73.8, y: 104 },
  { name: '智慧無礙', no: '2', d: 'M150,150 L271.24,220 A140,140 0 0 1 28.76,220 Z', fill: '#dfe9d8', x: 150, y: 240 },
  { name: '善巧方便', no: '3', d: 'M150,150 L150,10 A140,140 0 0 1 271.24,220 Z', fill: '#f3e6c8', x: 226.2, y: 104 },
]
</script>

<template>
  <svg viewBox="0 0 300 300" class="wheel" role="img" aria-label="協會三主軸：慈悲關懷、智慧無礙、善巧方便">
    <path v-for="s in sectors" :key="s.name" :d="s.d" :fill="s.fill" stroke="#a02a27" stroke-width="3" />
    <g v-for="s in sectors" :key="`t-${s.name}`" text-anchor="middle">
      <text :x="s.x" :y="s.y - 18" class="no">{{ s.no }}</text>
      <text :x="s.x" :y="s.y + 6" class="name">{{ s.name }}</text>
    </g>
    <circle cx="150" cy="150" r="47" fill="#f3e6c8" stroke="#a02a27" stroke-width="3" />
    <image href="/images/association/logo.webp" x="108" y="108" width="84" height="84" />
  </svg>
</template>

<style scoped>
.wheel {
  display: block;
  width: 100%;
  max-width: 300px;
  overflow: visible;
  margin: 0 auto;
}
.no {
  font-size: 13px;
  fill: #a02a27;
  font-weight: 700;
}
.name {
  font-size: 19px;
  font-weight: 600;
  fill: #a02a27;
  letter-spacing: 1px;
}
</style>
