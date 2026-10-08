// 行事曆共用工具（官網與後台共用）。版面與操作照 dengo 專案的行程管理：清單／日／週三種檢視、
// 手機一次看 4 天、左右拖曳換週、滑過行程時左側時間欄標出時段。
import { onBeforeUnmount, onMounted, ref, type Ref } from 'vue'

export type CalView = 'list' | 'day' | 'week'

// 畫在行事曆上的一筆行程；顏色與標籤由各頁自己決定
export type CalEvent = {
  id: string
  title: string
  start: string // ISO 時間
  end: string | null // 沒有結束時間的，畫成一小時
  location?: string
  description?: string
  color: string // 左側色條與時間字色
  bg?: string // 底色，預設白
  dashed?: boolean // 虛線框：例如「可以報名」的活動
  prefix?: string // 時間前的小字，例如「協會」
  badge?: string // 右上角小標，例如「#講座」
}

export const pad = (n: number) => String(n).padStart(2, '0')

export function startOfWeek(d: Date) {
  const x = new Date(d)
  x.setDate(x.getDate() - ((x.getDay() + 6) % 7)) // 週一開始
  x.setHours(0, 0, 0, 0)
  return x
}

export function startOfDay(d: Date) {
  const x = new Date(d)
  x.setHours(0, 0, 0, 0)
  return x
}

export function addDays(d: Date, n: number) {
  const x = new Date(d)
  x.setDate(x.getDate() + n)
  return x
}

export const sameDay = (a: Date, b: Date) => a.toDateString() === b.toDateString()
export const hm = (d: Date) => `${pad(d.getHours())}:${pad(d.getMinutes())}`
export const endOf = (e: CalEvent) => (e.end ? new Date(e.end) : new Date(new Date(e.start).getTime() + 3600_000))

// 「10/8（四）19:00–21:00」；跨日就把結束日期也寫出來
const dayFmt = new Intl.DateTimeFormat('zh-TW', { month: 'numeric', day: 'numeric', weekday: 'short' })
export function rangeText(start: string, end: string | null) {
  const s = new Date(start)
  if (!end) return `${dayFmt.format(s)} ${hm(s)}`
  const e = new Date(end)
  return sameDay(s, e) ? `${dayFmt.format(s)} ${hm(s)}–${hm(e)}` : `${dayFmt.format(s)} ${hm(s)} – ${dayFmt.format(e)} ${hm(e)}`
}

// 某一天裡這筆行程佔的分鐘（跨日的行程在每一天各畫一段）
export function minutesOn(e: CalEvent, day: Date): [number, number] | null {
  const d0 = startOfDay(day).getTime()
  const d1 = d0 + 24 * 3600_000
  const s = new Date(e.start).getTime()
  const t = endOf(e).getTime()
  if (t <= d0 || s >= d1) return null
  return [Math.max(s, d0) - d0, Math.min(t, d1) - d0].map((ms) => Math.round(ms / 60000)) as [number, number]
}

// 時間軸範圍：預設 8–22 點，有更早或更晚的行程就往外延伸
export function hourRange(events: CalEvent[], days: Date[]): [number, number] {
  let lo = 8
  let hi = 22
  for (const d of days)
    for (const e of events) {
      const m = minutesOn(e, d)
      if (!m) continue
      lo = Math.min(lo, Math.floor(m[0] / 60))
      hi = Math.max(hi, Math.ceil(m[1] / 60))
    }
  return [lo, Math.min(hi, 24)]
}

// 手機（≤640px）一次看 4 天，電腦看 7 天
export function useDayCount(): Ref<number> {
  const count = ref(7)
  let mq: MediaQueryList | null = null
  const update = () => (count.value = mq?.matches ? 4 : 7)
  onMounted(() => {
    mq = window.matchMedia('(max-width: 640px)')
    update()
    mq.addEventListener('change', update)
  })
  onBeforeUnmount(() => mq?.removeEventListener('change', update))
  return count
}

// 左右拖曳換頁：拖超過 60px 就換；拖動過的那一下不算點擊
export function useHorizontalDrag(onPrev: () => void, onNext: () => void) {
  const dx = ref(0)
  let x0: number | null = null
  let moved = false
  function down(e: PointerEvent) {
    if (e.pointerType === 'mouse' && e.button !== 0) return
    x0 = e.clientX
    moved = false
  }
  function move(e: PointerEvent) {
    if (x0 === null) return
    dx.value = e.clientX - x0
    if (Math.abs(dx.value) > 8) moved = true
  }
  function up() {
    if (x0 === null) return
    if (dx.value > 60) onPrev()
    else if (dx.value < -60) onNext()
    x0 = null
    dx.value = 0
  }
  // 拖曳結束後緊接的 click 不要打開行程
  const clicked = () => {
    const was = moved
    moved = false
    return !was
  }
  return { dx, down, move, up, clicked }
}

// 行程類型（dengo 的 General／Training／Meeting／Activity／Ceremony），顏色取自站上配色
export const KINDS = [
  { value: 'general', label: '一般', color: '#426631' },
  { value: 'training', label: '培訓', color: '#3f6f8f' },
  { value: 'meeting', label: '會議', color: '#b9852f' },
  { value: 'activity', label: '活動', color: '#6f9a5a' },
  { value: 'ceremony', label: '典禮', color: '#ba5854' },
] as const
export const kindOf = (v: string | undefined) => KINDS.find((k) => k.value === v) ?? KINDS[0]

// 表單用的本地時間字串 YYYY-MM-DDTHH:mm（牆上時間）
export function localInput(d: Date | string | null | undefined) {
  if (!d) return ''
  const x = new Date(d)
  return `${x.getFullYear()}-${pad(x.getMonth() + 1)}-${pad(x.getDate())}T${pad(x.getHours())}:${pad(x.getMinutes())}`
}
export function shiftLocal(v: string, minutes: number) {
  if (!v) return ''
  const d = new Date(v)
  if (isNaN(d.getTime())) return ''
  return localInput(new Date(d.getTime() + minutes * 60000))
}

// 說明開頭自動寫「共計 X 小時」，保留其餘文字（dengo 的 applyDuration）
export function withDuration(description: string, start: string, end: string) {
  const base = description.replace(/^共計 [\d.]+ 小時\n?/, '')
  const ms = start && end ? new Date(end).getTime() - new Date(start).getTime() : NaN
  if (!(ms > 0)) return base
  const text = `共計 ${Math.round((ms / 3600000) * 10) / 10} 小時`
  return base ? `${text}\n${base}` : text
}
