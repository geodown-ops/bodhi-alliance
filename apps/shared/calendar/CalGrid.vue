<script setup lang="ts">
// 日／週時間格（照 dengo 的週曆）：表頭週六綠、週日紅、平日米色，今天的日期圈白；
// 單數日淡灰底；滑過行程時左側時間欄標出時段；左右拖曳換頁。
import { computed, ref } from 'vue'
import { hm, hourRange, minutesOn, pad, sameDay, useHorizontalDrag, type CalEvent } from './cal'

const props = defineProps<{ days: Date[]; events: CalEvent[]; hourHeight?: number }>()
const emit = defineEmits<{ select: [id: string]; prev: []; next: [] }>()

const H = computed(() => props.hourHeight ?? (props.days.length === 1 ? 64 : 52))
const range = computed(() => hourRange(props.events, props.days))
const hours = computed(() => Array.from({ length: range.value[1] - range.value[0] + 1 }, (_, i) => range.value[0] + i))
const totalH = computed(() => (range.value[1] - range.value[0]) * H.value + 16)
const single = computed(() => props.days.length === 1)
const today = new Date()
const weekday = new Intl.DateTimeFormat('zh-TW', { weekday: 'short' })

type Placed = { e: CalEvent; top: number; height: number; left: number; width: number; short: boolean }

// 重疊的行程並排（同一串重疊的平分寬度）
function placed(day: Date): Placed[] {
  const lo = range.value[0] * 60
  const hi = range.value[1] * 60
  const items = props.events
    .map((e) => ({ e, m: minutesOn(e, day) }))
    .filter((x): x is { e: CalEvent; m: [number, number] } => !!x.m)
    .sort((a, b) => a.m[0] - b.m[0] || b.m[1] - a.m[1])
  const out: Placed[] = []
  let cluster: { x: (typeof items)[number]; lane: number }[] = []
  let lanesEnd: number[] = []
  let clusterEnd = -1
  const flush = () => {
    const n = lanesEnd.length
    for (const { x, lane } of cluster) {
      const s = Math.max(x.m[0], lo)
      const t = Math.min(x.m[1], hi)
      out.push({
        e: x.e,
        top: ((s - lo) / 60) * H.value,
        height: Math.max(((t - s) / 60) * H.value - 2, 22),
        left: lane / n,
        width: 1 / n,
        short: x.m[1] - x.m[0] < 120,
      })
    }
    cluster = []
    lanesEnd = []
  }
  for (const x of items) {
    if (x.m[0] >= clusterEnd) flush()
    let lane = lanesEnd.findIndex((end) => end <= x.m[0])
    if (lane < 0) lane = lanesEnd.push(0) - 1
    lanesEnd[lane] = x.m[1]
    clusterEnd = Math.max(clusterEnd, x.m[1])
    cluster.push({ x, lane })
  }
  flush()
  return out
}

const hover = ref<{ top: number; height: number } | null>(null)
const drag = useHorizontalDrag(
  () => emit('prev'),
  () => emit('next'),
)
function open(id: string) {
  if (drag.clicked()) emit('select', id)
}
const headBg = (d: Date) => (d.getDay() === 6 ? '#8dae78' : d.getDay() === 0 ? '#ba5854' : '#c9a86a')
const cut = (s: string, n: number) => (s.length > n ? s.slice(0, n) + '…' : s)
const timeOf = (e: CalEvent, d: Date) => {
  const m = minutesOn(e, d)!
  const t = (x: number) => (x >= 1440 ? '24:00' : `${pad(Math.floor(x / 60))}:${pad(x % 60)}`)
  return e.end ? `${t(m[0])}–${t(m[1])}` : hm(new Date(e.start))
}
</script>

<template>
  <div class="cal-scroll">
    <div
      class="cal-grid"
      :class="{ single }"
      :style="{ '--cols': days.length, transform: drag.dx.value ? `translateX(${drag.dx.value / 3}px)` : undefined }"
      @pointerdown="drag.down"
      @pointermove="drag.move"
      @pointerup="drag.up"
      @pointercancel="drag.up"
      @pointerleave="drag.up"
    >
      <template v-if="!single">
        <div />
        <div v-for="d in days" :key="'h' + d.toDateString()" class="head" :style="{ background: headBg(d) }">
          <span>{{ weekday.format(d) }}</span>
          <span class="num" :class="{ today: sameDay(d, today) }">{{ d.getDate() }}</span>
        </div>
      </template>

      <div class="gutter" :style="{ height: totalH + 'px' }">
        <div v-if="hover" class="hover" :style="{ top: hover.top + 'px', height: hover.height + 'px' }" />
        <div v-for="h in hours" :key="h" class="hour" :style="{ top: (h - range[0]) * H + 'px' }">{{ single ? `${pad(h)}:00` : h }}</div>
      </div>

      <div v-for="d in days" :key="d.toDateString()" class="col" :class="{ odd: !single && d.getDate() % 2 === 1 }" :style="{ height: totalH + 'px' }">
        <div v-for="h in hours" :key="h" class="line" :style="{ top: (h - range[0]) * H + 'px' }" />
        <div
          v-for="p in placed(d)"
          :key="p.e.id"
          class="evt"
          :class="{ dashed: p.e.dashed }"
          :style="{
            top: p.top + 'px',
            height: p.height + 'px',
            left: `calc(${p.left * 100}% + 2px)`,
            width: `calc(${p.width * 100}% - 4px)`,
            borderLeftColor: p.e.color,
            background: p.e.bg ?? '#fff',
          }"
          @click="open(p.e.id)"
          @mouseenter="hover = { top: p.top, height: p.height }"
          @mouseleave="hover = null"
        >
          <span v-if="p.e.badge" class="badge">{{ p.e.badge }}</span>
          <!-- 和別的行程並排時欄位窄，只留時間與名稱 -->
          <p class="t" :style="{ color: p.e.color }">{{ p.e.prefix && p.width === 1 ? p.e.prefix + ' ' : '' }}{{ p.width === 1 ? timeOf(p.e, d) : timeOf(p.e, d).slice(0, 5) }}</p>
          <p class="name">{{ p.short && !single ? cut(p.e.title, 8) : p.e.title }}</p>
          <template v-if="p.width === 1 || single">
            <p v-if="p.e.location" class="sub">{{ p.e.location }}</p>
            <p v-if="p.e.description && !p.short" class="sub desc">{{ single ? p.e.description : cut(p.e.description, 30) }}</p>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.cal-scroll {
  overflow-x: auto;
  padding-bottom: 16px;
}
.cal-grid {
  display: grid;
  grid-template-columns: 40px repeat(var(--cols), minmax(0, 1fr));
  user-select: none;
  touch-action: pan-y;
  min-width: 0;
}
.cal-grid.single {
  grid-template-columns: 52px 1fr;
}
.head {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  padding: 2px 0;
  color: #fff;
  font-weight: 700;
  font-size: 0.85rem;
  border-left: 1px solid #d6cdbd;
}
.head .num {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  font-size: 1rem;
}
.head .num.today {
  background: #fff;
  color: #3b2a20;
}
.gutter {
  position: relative;
}
.gutter .hover {
  position: absolute;
  left: 0;
  right: 0;
  border-radius: 4px;
  background: rgba(201, 168, 106, 0.45);
  pointer-events: none;
}
.hour {
  position: absolute;
  left: 0;
  right: 4px;
  text-align: center;
  transform: translateY(-50%);
  font-size: 0.85rem;
  font-weight: 700;
  color: #6b5c4f;
}
.col {
  position: relative;
  border-left: 1px solid #d6cdbd;
}
.col.odd {
  background: rgba(0, 0, 0, 0.03);
}
.line {
  position: absolute;
  left: 0;
  right: 0;
  border-top: 1px dashed #d6cdbd;
}
.evt {
  position: absolute;
  overflow: hidden;
  border: 1px solid #e3dccf;
  border-left: 3px solid;
  border-radius: 6px;
  padding: 2px 5px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.06);
  cursor: pointer;
  transition: box-shadow 0.15s;
}
.evt:hover {
  box-shadow: 0 0 0 2px #3b2a20;
  z-index: 5;
}
.evt.dashed {
  border-style: dashed;
  border-left-style: solid;
}
.evt p {
  margin: 0;
  line-height: 1.25;
  word-break: break-word;
}
.evt .t {
  font-size: 10px;
}
.evt .name {
  font-size: 12px;
  font-weight: 600;
  color: #3b2a20;
}
.evt .sub {
  font-size: 10px;
  color: #6b5c4f;
}
.evt .desc {
  white-space: pre-line;
}
.evt .badge {
  float: right;
  margin-left: 2px;
  font-size: 9px;
  font-weight: 700;
  border-radius: 4px;
  padding: 0 3px;
  background: #f3ead8;
  color: #8a6a32;
}
.single .evt {
  padding: 4px 10px;
  border-left-width: 4px;
}
.single .evt .t {
  font-size: 12px;
}
.single .evt .name {
  font-size: 14px;
}
.single .evt .sub {
  font-size: 12px;
}
</style>
