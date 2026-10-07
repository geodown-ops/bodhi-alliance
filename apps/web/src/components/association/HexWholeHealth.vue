<script setup lang="ts">
// 全人健康六個面向（環境、理智、感情、社會、身體、精神），以及向外的六種智商。
const R = 46
const cx = 220
const cy = 170
const d = Math.sqrt(3) * R + 6
const hex = (x: number, y: number) =>
  [0, 60, 120, 180, 240, 300]
    .map((a) => `${(x + R * Math.cos((a * Math.PI) / 180)).toFixed(1)},${(y + R * Math.sin((a * Math.PI) / 180)).toFixed(1)}`)
    .join(' ')

const facets = [
  { name: '環境', a: 90 },
  { name: '理智', a: 30 },
  { name: '感情', a: -30 },
  { name: '社會', a: -90 },
  { name: '身體', a: -150 },
  { name: '精神', a: 150 },
].map((f) => ({
  ...f,
  x: cx + d * Math.cos((f.a * Math.PI) / 180),
  y: cy - d * Math.sin((f.a * Math.PI) / 180),
}))

const axes = [
  { name: 'SQ', a: 60 },
  { name: 'VQ', a: 0 },
  { name: 'EQ', a: -60 },
  { name: 'CQ', a: -120 },
  { name: 'AQ', a: 180 },
  { name: 'PQ', a: 120 },
].map((q) => {
  const ux = Math.cos((q.a * Math.PI) / 180)
  const uy = -Math.sin((q.a * Math.PI) / 180)
  return { ...q, x2: cx + 148 * ux, y2: cy + 148 * uy, lx: cx + 172 * ux, ly: cy + 172 * uy + 5, ux, uy }
})
const head = (q: (typeof axes)[number]) => {
  const px = -q.uy
  const py = q.ux
  const bx = q.x2 - 10 * q.ux
  const by = q.y2 - 10 * q.uy
  return `${q.x2},${q.y2} ${bx + 5 * px},${by + 5 * py} ${bx - 5 * px},${by - 5 * py}`
}
</script>

<template>
  <svg viewBox="0 0 440 340" class="hex" role="img" aria-label="全人健康：環境、理智、感情、社會、身體、精神六個面向">
    <g v-for="q in axes" :key="q.name">
      <line :x1="cx" :y1="cy" :x2="q.x2" :y2="q.y2" stroke="#ba5854" stroke-width="2" />
      <polygon :points="head(q)" fill="#ba5854" />
      <text :x="q.lx" :y="q.ly" text-anchor="middle" class="axis">{{ q.name }}</text>
    </g>
    <g v-for="f in facets" :key="f.name">
      <polygon :points="hex(f.x, f.y)" fill="#f3e6c8" stroke="#f3e6c8" stroke-width="1.5" />
      <text :x="f.x" :y="f.y + 6" text-anchor="middle" class="facet">{{ f.name }}</text>
    </g>
    <polygon :points="hex(cx, cy)" fill="#f1d9d4" stroke="#ba5854" stroke-width="2" />
    <text :x="cx" :y="cy - 4" text-anchor="middle" class="core">全人</text>
    <text :x="cx" :y="cy + 20" text-anchor="middle" class="core">健康</text>
  </svg>
</template>

<style scoped>
.hex {
  display: block;
  width: 100%;
  max-width: 420px;
  margin: 0 auto;
}
.axis {
  font-size: 15px;
  font-weight: 700;
  fill: #ba5854;
}
.facet {
  font-size: 16px;
  fill: #ba5854;
}
.core {
  font-size: 18px;
  font-weight: 700;
  fill: #ba5854;
}
</style>
