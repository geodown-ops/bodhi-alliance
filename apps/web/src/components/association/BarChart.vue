<script setup lang="ts">
// 橫條圖：第一項是最高值，以強調色顯示。
const props = defineProps<{
  title: string
  question?: string
  items: { label: string; value: number }[]
  note?: string
  color?: string
}>()
const max = Math.max(...props.items.map((i) => i.value))
</script>

<template>
  <div class="bars">
    <h4>{{ title }}</h4>
    <p v-if="question" class="question">{{ question }}</p>
    <div v-for="(item, i) in items" :key="item.label" class="row" :class="{ top: i === 0 }">
      <span class="label">{{ item.label }}</span>
      <span class="track">
        <span class="fill" :style="{ width: `${(item.value / max) * 100}%`, background: i === 0 ? color ?? '#5d7f52' : undefined }" />
      </span>
      <span class="value" :style="i === 0 ? { color: color ?? '#5d7f52' } : undefined">{{ item.value.toFixed(1) }}</span>
    </div>
    <p v-if="note" class="note-line">{{ note }}</p>
  </div>
</template>

<style scoped>
.bars h4 {
  margin: 0 0 4px;
  font-size: 1.1rem;
  font-weight: 700;
  line-height: 1.5;
  color: var(--ink);
}
.question {
  margin: 0 0 14px;
  font-size: 0.9rem !important;
  line-height: 1.6 !important;
  color: var(--ink-faint) !important;
}
.row {
  display: grid;
  grid-template-columns: minmax(0, 11em) minmax(0, 1fr) 3em;
  gap: 10px;
  align-items: center;
  margin: 7px 0;
}
.label {
  font-size: 0.92rem;
  line-height: 1.35;
  text-align: right;
  color: var(--ink-soft);
}
.track {
  height: 14px;
  border-radius: 7px;
  background: var(--ground-sunk);
  overflow: hidden;
}
.fill {
  display: block;
  height: 100%;
  border-radius: 7px;
  background: #c9bfae;
}
.value {
  font-variant-numeric: tabular-nums;
  font-size: 0.95rem;
  color: var(--ink-soft);
}
.row.top .label,
.row.top .value {
  font-weight: 700;
  font-size: 1.02rem;
}
.note-line {
  margin: 10px 0 0;
  text-align: right;
  font-size: 0.85rem !important;
  color: var(--ink-faint) !important;
}
@media (max-width: 600px) {
  .row {
    grid-template-columns: minmax(0, 8em) minmax(0, 1fr) 2.8em;
  }
  .label {
    font-size: 0.85rem;
  }
}
</style>
