<script setup lang="ts">
// 圓環百分比：value 是 0–100。
const props = withDefaults(defineProps<{ value: number; color: string; track?: string }>(), {
  track: '#f1d9d4',
})
const r = 52
const c = 2 * Math.PI * r
const dash = `${(c * props.value) / 100} ${c}`
</script>

<template>
  <svg viewBox="0 0 128 128" class="ring" role="img" :aria-label="`${value}%`">
    <circle cx="64" cy="64" :r="r" fill="none" :stroke="track" stroke-width="12" />
    <circle
      cx="64"
      cy="64"
      :r="r"
      fill="none"
      :stroke="color"
      stroke-width="12"
      stroke-linecap="round"
      :stroke-dasharray="dash"
      transform="rotate(-90 64 64)"
    />
    <text x="64" y="74" text-anchor="middle" class="value" fill="#a02a27">
      {{ value }}<tspan class="pct">%</tspan>
    </text>
  </svg>
</template>

<style scoped>
.ring {
  display: block;
  width: 128px;
  height: 128px;
  margin: 0 auto 8px;
}
.value {
  font-size: 32px;
  font-weight: 700;
}
.pct {
  font-size: 16px;
}
</style>
