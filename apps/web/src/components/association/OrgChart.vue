<script setup lang="ts">
// 協會組織架構圖（依簡報「標準組織」第二頁重畫）。直排的方塊一字一行。
type Box = { label: string; x: number; y: number; w: number; h: number; tone?: 'red' | 'plum'; vertical?: boolean }

const boxes: Box[] = [
  { label: '會員大會', x: 310, y: 15, w: 270, h: 52 },
  { label: '監事會', x: 20, y: 72, w: 265, h: 52 },
  { label: '理事會', x: 605, y: 72, w: 265, h: 52 },
  { label: '理事長', x: 310, y: 130, w: 270, h: 52 },
  { label: '秘書處', x: 20, y: 217, w: 210, h: 42, tone: 'red' },
  { label: '工作委員會', x: 263, y: 217, w: 364, h: 42, tone: 'red' },
  ...['行政', '總務', '財務', '秘書'].map((label, i) => ({ label, x: 25 + i * 52, y: 288, w: 44, h: 112, vertical: true })),
  ...['公關委員會', '業務拓展委員會', '募款委員會', '資訊委員會', '寺廟規畫委員會', '課程發展中心'].map((label, i) => ({
    label,
    x: 263 + i * 62.8,
    y: 288,
    w: 50,
    h: 252,
    vertical: true,
  })),
  ...['課程一', '課程二', '課程三', '課程四', '課程五', '課程六'].map((label, i) => ({
    label,
    x: 652,
    y: 280 + i * 44,
    w: 95,
    h: 40,
    tone: 'plum' as const,
  })),
  { label: '審查', x: 770, y: 280, w: 52, h: 128, vertical: true },
  { label: '紀律', x: 770, y: 412, w: 52, h: 128, vertical: true },
  { label: '指導委員會', x: 845, y: 280, w: 52, h: 260, vertical: true },
  { label: '認證暨品質中心', x: 915, y: 280, w: 52, h: 260, vertical: true, tone: 'red' },
]

const courseY = [300, 344, 388, 432, 476, 520]
const lines = [
  'M445,67 V203',
  'M285,98 H605',
  'M125,203 H445',
  'M125,203 V217',
  'M445,203 V217',
  'M125,259 V274 M47,274 H203 M47,274 V288 M99,274 V288 M151,274 V288 M203,274 V288',
  'M445,259 V274 M288,274 H602 ' + [0, 1, 2, 3, 4, 5].map((i) => `M${288 + i * 62.8},274 V288`).join(' '),
  ...courseY.map((y) => `M627,${y} H652`),
  ...courseY.map((y) => `M747,${y} H758`),
  'M758,300 V388 M758,344 H770',
  'M758,432 V520 M758,476 H770',
  'M822,344 H833 M822,476 H833 M833,344 V476 M833,410 H845',
  'M897,410 H915',
]

const fill = (b: Box) => (b.tone === 'red' ? '#ba5854' : b.tone === 'plum' ? '#ba5854' : '#f3e6c8')
const ink = (b: Box) => (b.tone ? '#f3e6c8' : '#ba5854')
const chars = (b: Box) => {
  const step = 21
  const top = b.y + (b.h - b.label.length * step) / 2 + 16
  return [...b.label].map((c, i) => ({ c, y: top + i * step }))
}
</script>

<template>
  <div class="scroll">
    <svg viewBox="0 0 985 560" class="org" role="img" aria-label="協會組織架構圖">
      <path v-for="(d, i) in lines" :key="i" :d="d" fill="none" stroke="#ba5854" stroke-width="1.5" />
      <g v-for="b in boxes" :key="b.label">
        <rect :x="b.x" :y="b.y" :width="b.w" :height="b.h" rx="6" :fill="fill(b)" :stroke="b.tone ? 'none' : '#ba5854'" />
        <template v-if="b.vertical">
          <text v-for="t in chars(b)" :key="t.y" :x="b.x + b.w / 2" :y="t.y" text-anchor="middle" class="t" :fill="ink(b)">{{ t.c }}</text>
        </template>
        <text v-else :x="b.x + b.w / 2" :y="b.y + b.h / 2 + 6" text-anchor="middle" class="t" :fill="ink(b)">{{ b.label }}</text>
      </g>
    </svg>
  </div>
</template>

<style scoped>
.scroll {
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}
.org {
  display: block;
  width: 100%;
  min-width: 720px;
}
.t {
  font-size: 17px;
}
</style>
