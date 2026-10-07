<script setup lang="ts">
// 聯合國永續發展目標（SDGs）分組圓盤：顏色用 SDGs 官方配色，名稱沿用簡報用字。
const props = defineProps<{ label: string; goals: number[] }>()

const sdg: Record<number, { name: string; color: string }> = {
  1: { name: '消除貧窮', color: '#a02a27' },
  2: { name: '消除飢餓', color: '#a02a27' },
  3: { name: '健康與福祉', color: '#a02a27' },
  4: { name: '教育品質', color: '#a02a27' },
  5: { name: '性別平等', color: '#a02a27' },
  6: { name: '淨水與衛生', color: '#a02a27' },
  7: { name: '可負擔能源', color: '#a02a27' },
  8: { name: '就業與經濟成長', color: '#a02a27' },
  9: { name: '工業、創新基礎建設', color: '#a02a27' },
  11: { name: '永續城市', color: '#a02a27' },
  12: { name: '責任消費與生產', color: '#a02a27' },
  13: { name: '氣候行動', color: '#a02a27' },
  15: { name: '陸地生態', color: '#a02a27' },
  16: { name: '和平與正義制度', color: '#a02a27' },
}

const R = 120
const r = 50
const n = props.goals.length
const pt = (radius: number, a: number) => [150 + radius * Math.sin(a), 150 - radius * Math.cos(a)]
const segments = props.goals.map((g, i) => {
  const a0 = (2 * Math.PI * i) / n
  const a1 = (2 * Math.PI * (i + 1)) / n
  const [x0, y0] = pt(R, a0)
  const [x1, y1] = pt(R, a1)
  const [x2, y2] = pt(r, a1)
  const [x3, y3] = pt(r, a0)
  const large = a1 - a0 > Math.PI ? 1 : 0
  const [lx, ly] = pt((R + r) / 2, (a0 + a1) / 2)
  return {
    goal: g,
    ...sdg[g],
    d: `M${x0},${y0} A${R},${R} 0 ${large} 1 ${x1},${y1} L${x2},${y2} A${r},${r} 0 ${large} 0 ${x3},${y3} Z`,
    lx,
    ly,
  }
})
</script>

<template>
  <figure class="sdg">
    <svg viewBox="0 0 300 300" role="img" :aria-label="`${label}：${segments.map((s) => `${s.goal} ${s.name}`).join('、')}`">
      <path v-for="s in segments" :key="s.goal" :d="s.d" :fill="s.color" stroke="#f3e6c8" stroke-width="3" />
      <text v-for="s in segments" :key="`n-${s.goal}`" :x="s.lx" :y="s.ly + 8" text-anchor="middle" class="num">
        {{ String(s.goal).padStart(2, '0') }}
      </text>
      <circle cx="150" cy="150" :r="r - 4" fill="#f3e6c8" />
      <text x="150" y="158" text-anchor="middle" class="label">{{ label }}</text>
    </svg>
    <figcaption>
      <span v-for="s in segments" :key="`c-${s.goal}`" class="chip">
        <i :style="{ background: s.color }" />{{ s.goal }} {{ s.name }}
      </span>
    </figcaption>
  </figure>
</template>

<style scoped>
.sdg {
  margin: 0;
}
.sdg svg {
  display: block;
  width: 100%;
  max-width: 240px;
  margin: 0 auto;
}
.num {
  font-size: 22px;
  font-weight: 700;
  fill: #f3e6c8;
}
.label {
  font-size: 22px;
  font-weight: 700;
  fill: #a02a27;
}
figcaption {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  justify-content: center;
  margin-top: 12px;
}
.chip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 2px 10px 2px 6px;
  border-radius: 999px;
  background: var(--ground-sunk);
  font-size: 0.82rem;
  color: var(--ink-soft);
}
.chip i {
  width: 10px;
  height: 10px;
  border-radius: 50%;
}
</style>
