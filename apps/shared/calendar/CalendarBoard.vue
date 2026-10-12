<script setup lang="ts">
// 行事曆看板（照 dengo 的行程管理）：右上角切換清單／日／週，日／週有前後翻頁、點日期跳到某天、
// 「今天」「本週」回到現在；清單只列還沒結束的，過去的收在最下面。
import { computed, ref } from 'vue'
import CalGrid from './CalGrid.vue'
import { addDays, endOf, pad, rangeText, startOfDay, startOfWeek, useDayCount, type CalEvent, type CalView } from './cal'

const props = defineProps<{ events: CalEvent[]; emptyText?: string }>()
const view = defineModel<CalView>('view', { default: 'list' })
const emit = defineEmits<{ select: [id: string] }>()

const dayCount = useDayCount()
const weekStart = ref(startOfWeek(new Date()))
const day = ref(startOfDay(new Date()))

const days = computed(() =>
  view.value === 'day' ? [day.value] : Array.from({ length: dayCount.value }, (_, i) => addDays(weekStart.value, i)),
)
function move(n: number) {
  if (view.value === 'day') day.value = addDays(day.value, n)
  else weekStart.value = addDays(weekStart.value, n * dayCount.value)
}
function goToday() {
  day.value = startOfDay(new Date())
  weekStart.value = startOfWeek(new Date())
}
// q-date 用 YYYY/MM/DD
const pickValue = computed(() => {
  const d = view.value === 'day' ? day.value : weekStart.value
  return `${d.getFullYear()}/${pad(d.getMonth() + 1)}/${pad(d.getDate())}`
})
function pick(v: string | null) {
  if (!v) return
  const d = startOfDay(new Date(v.replace(/\//g, '-') + 'T00:00'))
  day.value = d
  weekStart.value = dayCount.value === 7 ? startOfWeek(d) : d
}
const md = (d: Date) => `${d.getMonth() + 1}月${d.getDate()}日`
const label = computed(() =>
  view.value === 'day'
    ? day.value.toLocaleDateString('zh-TW', { year: 'numeric', month: 'long', day: 'numeric', weekday: 'long' })
    : `${md(days.value[0])} – ${md(days.value[days.value.length - 1])}`,
)
const dayEmpty = computed(() => view.value === 'day' && !props.events.some((e) => new Date(e.start) < addDays(day.value, 1) && endOf(e) > day.value))

// 清單：還沒結束的照時間排；過去的收起來
const showPast = ref(false)
const sorted = computed(() => [...props.events].sort((a, b) => a.start.localeCompare(b.start)))
const upcoming = computed(() => sorted.value.filter((e) => endOf(e) >= startOfDay(new Date())))
const past = computed(() => sorted.value.filter((e) => endOf(e) < startOfDay(new Date())).reverse())
const views = computed(() => [
  { label: '清單', value: 'list' },
  { label: '日', value: 'day' },
  { label: dayCount.value === 7 ? '週' : '4日', value: 'week' },
])
</script>

<template>
  <div class="cal-board">
    <div class="cal-bar">
      <div v-if="view !== 'list'" class="cal-nav">
        <q-btn flat dense round icon="chevron_left" aria-label="上一頁" @click="move(-1)" />
        <q-btn flat dense no-caps class="cal-label" :label="label">
          <q-popup-proxy cover transition-show="scale" transition-hide="scale">
            <q-date :model-value="pickValue" mask="YYYY/MM/DD" minimal first-day-of-week="1" @update:model-value="pick">
              <div class="row justify-end"><q-btn v-close-popup flat no-caps label="確定" /></div>
            </q-date>
          </q-popup-proxy>
        </q-btn>
        <q-btn flat dense round icon="chevron_right" aria-label="下一頁" @click="move(1)" />
        <q-btn outline dense no-caps size="sm" class="q-px-sm" :label="view === 'day' ? '今天' : '本週'" @click="goToday" />
      </div>
      <slot name="filters" />
      <q-btn-toggle v-model="view" class="cal-views" no-caps unelevated dense toggle-color="secondary" :options="views" />
    </div>

    <template v-if="view === 'list'">
      <p v-if="!upcoming.length && !past.length" class="cal-empty">{{ emptyText ?? '還沒有行程' }}</p>
      <p v-else-if="!upcoming.length" class="cal-empty">近期沒有行程</p>
      <div class="cal-list">
        <div v-for="e in upcoming" :key="e.id" class="cal-item" :style="{ borderLeftColor: e.color, background: e.bg }" @click="emit('select', e.id)">
          <slot name="item" :event="e">
            <div class="text-caption" :style="{ color: e.color }">{{ e.prefix ? e.prefix + ' · ' : '' }}{{ rangeText(e.start, e.end) }}</div>
            <div class="text-weight-bold">{{ e.title }}<span v-if="e.badge" class="cal-badge">{{ e.badge }}</span></div>
            <div v-if="e.location" class="text-caption">{{ e.location }}</div>
          </slot>
        </div>
      </div>
      <template v-if="past.length">
        <q-btn flat dense no-caps color="grey-8" class="q-mt-sm" :icon="showPast ? 'expand_less' : 'expand_more'" :label="`過去的行程（${past.length}）`" @click="showPast = !showPast" />
        <div v-if="showPast" class="cal-list past">
          <div v-for="e in past" :key="e.id" class="cal-item" :style="{ borderLeftColor: e.color }" @click="emit('select', e.id)">
            <slot name="item" :event="e">
              <div class="text-caption">{{ rangeText(e.start, e.end) }}</div>
              <div class="text-weight-bold">{{ e.title }}</div>
            </slot>
          </div>
        </div>
      </template>
    </template>
    <template v-else>
      <p v-if="dayEmpty" class="cal-empty">這天沒有行程<br /><span class="text-caption">按左右箭頭或左右拖曳換日期</span></p>
      <CalGrid v-else :days="days" :events="events" @select="(id) => emit('select', id)" @prev="move(-1)" @next="move(1)" />
    </template>
  </div>
</template>

<style scoped>
.cal-bar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
}
.cal-nav {
  display: flex;
  align-items: center;
  gap: 2px;
}
.cal-label {
  font-weight: 600;
}
.cal-views {
  margin-left: auto;
  border: 1px solid rgba(0, 0, 0, 0.18);
}
.cal-empty {
  text-align: center;
  padding: 48px 0;
  color: #8a7c70;
}
.cal-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.cal-item {
  background: #fff;
  border: 1px solid #e3dccf;
  border-left: 4px solid;
  border-radius: 8px;
  padding: 10px 12px;
  cursor: pointer;
}
.cal-item:hover {
  box-shadow: 0 1px 6px rgba(0, 0, 0, 0.1);
}
.past .cal-item {
  opacity: 0.75;
}
.cal-badge {
  margin-left: 8px;
  font-size: 0.75rem;
  font-weight: 600;
  border: 1px solid #c9a86a;
  color: #8a6a32;
  background: #faf5ec;
  border-radius: 10px;
  padding: 0 6px;
}
</style>
