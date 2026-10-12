<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { Dialog, Notify } from 'quasar'
import { api } from '../session'
import CalendarBoard from '../../../shared/calendar/CalendarBoard.vue'
import DateTimeField from '../../../shared/calendar/DateTimeField.vue'
import { KINDS, kindOf, localInput, rangeText, shiftLocal, withDuration, type CalEvent, type CalView } from '../../../shared/calendar/cal'

// 世界佛教教育協會：發布協會會刊、會員行事曆、協會通知，並查看協會會員名單。
type Issue = { id?: string; title: string; issued_on: string; summary: string; url: string }
type AssocEvent = {
  id?: string
  title: string
  starts_at: string
  ends_at: string | null
  location: string
  description: string
  kind: string
  tag_id: string | null
  tag_name?: string
}
type Tag = { id: string; name: string }
type Notice = { id?: string; title: string; body: string; created_at?: string }
type Member = { id: string; legal_name: string; display_name: string; email: string; phone: string; line_id: string; in_groups: boolean; joined_at: string | null }

const tab = ref<'notices' | 'events' | 'issues' | 'members'>('notices')
const issues = ref<Issue[]>([])
const events = ref<AssocEvent[]>([])
const notices = ref<Notice[]>([])
const members = ref<Member[]>([])
const editIssue = ref<Issue | null>(null)
const tags = ref<Tag[]>([])
type EventForm = AssocEvent & { start: string; end: string; allDay: boolean }
const editEvent = ref<EventForm | null>(null)
const editNotice = ref<Notice | null>(null)

const toast = (e: unknown) => Notify.create({ type: 'negative', message: (e as Error).message })
const dt = (s: string | null | undefined) => (s ? new Date(s).toLocaleString('zh-TW', { hour12: false }) : '')
const d = (s: string | null | undefined) => (s ? new Date(s).toLocaleDateString('zh-TW') : '')


const issueCols = [
  { name: 'issued_on', label: '出刊日期', field: 'issued_on', sortable: true },
  { name: 'title', label: '會刊', field: 'title' },
  { name: 'url', label: '連結', field: 'url' },
]
const noticeCols = [
  { name: 'created_at', label: '發布時間', field: 'created_at', format: dt, sortable: true },
  { name: 'title', label: '標題', field: 'title' },
]
const memberCols = [
  { name: 'legal_name', label: '真實姓名', field: 'legal_name', sortable: true },
  { name: 'display_name', label: '暱稱', field: 'display_name' },
  { name: 'email', label: '電子郵件', field: 'email' },
  { name: 'phone', label: '手機', field: 'phone' },
  { name: 'line_id', label: 'LINE ID', field: 'line_id' },
  { name: 'in_groups', label: '也是覺行小組會員', field: 'in_groups', format: (v: boolean) => (v ? '是' : '') },
  { name: 'joined_at', label: '加入協會', field: 'joined_at', format: d, sortable: true },
]

async function load() {
  try {
    ;[issues.value, events.value, tags.value, notices.value, members.value] = await Promise.all([
      api.get<Issue[]>('/api/admin/association/issues'),
      api.get<AssocEvent[]>('/api/admin/association/events'),
      api.get<Tag[]>('/api/admin/association/tags'),
      api.get<Notice[]>('/api/admin/association/notices'),
      api.get<Member[]>('/api/admin/association/members'),
    ])
  } catch (e) {
    toast(e)
  }
}
onMounted(load)

async function save(kind: 'issues' | 'events' | 'notices', body: { id?: string } & Record<string, unknown>, done: () => void) {
  try {
    await api.send(body.id ? 'PUT' : 'POST', body.id ? `/api/admin/association/${kind}/${body.id}` : `/api/admin/association/${kind}`, body)
    done()
    load()
  } catch (e) {
    toast(e)
  }
}

function remove(kind: 'issues' | 'events' | 'notices', id: string, name: string) {
  Dialog.create({ title: '刪除', message: `確定刪除「${name}」？`, cancel: true }).onOk(async () => {
    try {
      await api.send('DELETE', `/api/admin/association/${kind}/${id}`)
      editIssue.value = null
      editEvent.value = null
      editNotice.value = null
      load()
    } catch (e) {
      toast(e)
    }
  })
}

const today = () => new Date().toISOString().slice(0, 10)

// ── 會員行事曆：照 dengo 的行程管理（清單／日／週、活動標籤、複製到隔天、往下複製）──
const calView = ref<CalView>('list')
const filterTag = ref('')
const visibleEvents = computed(() => events.value.filter((e) => !filterTag.value || e.tag_id === filterTag.value))
const calEvents = computed<CalEvent[]>(() =>
  visibleEvents.value.map((e) => ({
    id: e.id!,
    title: e.title,
    start: e.starts_at,
    end: e.ends_at,
    location: e.location,
    description: e.description,
    color: kindOf(e.kind).color,
    badge: e.tag_name ? `#${e.tag_name}` : undefined,
  })),
)
const byId = (id: string) => events.value.find((e) => e.id === id)!

// 整天＝08:00–20:00（同 dengo）
const isAllDay = (s: string, e: string) => !!s && !!e && s.slice(0, 10) === e.slice(0, 10) && s.slice(11) === '08:00' && e.slice(11) === '20:00'
function openEvent(e?: AssocEvent) {
  if (e) {
    const start = localInput(e.starts_at)
    const end = localInput(e.ends_at)
    editEvent.value = { ...e, start, end, allDay: isAllDay(start, end) }
    return
  }
  const s = new Date()
  s.setDate(s.getDate() + 1)
  s.setHours(9, 0, 0, 0)
  const start = localInput(s)
  const end = shiftLocal(start, 120)
  editEvent.value = {
    title: '',
    starts_at: '',
    ends_at: null,
    location: '',
    description: withDuration('', start, end),
    kind: 'general',
    tag_id: filterTag.value || null,
    start,
    end,
    allDay: false,
  }
}
// 改開始時間：結束沒填或早於開始，就帶開始 +2 小時；整天的行程跟著換日期
function setStart(v: string) {
  const f = editEvent.value!
  if (f.allDay) {
    f.start = `${v.slice(0, 10)}T08:00`
    f.end = `${v.slice(0, 10)}T20:00`
  } else {
    f.start = v
    if (!f.end || f.end <= v) f.end = shiftLocal(v, 120)
  }
  f.description = withDuration(f.description, f.start, f.end)
}
function setEnd(v: string) {
  const f = editEvent.value!
  f.end = v
  f.description = withDuration(f.description, f.start, f.end)
}
function setAllDay(on: boolean) {
  const f = editEvent.value!
  f.allDay = on
  if (on) setStart(f.start || `${today()}T08:00`)
}
function payload(f: EventForm, start = f.start, end = f.end) {
  return {
    title: f.title.trim(),
    kind: f.kind,
    tag_id: f.tag_id || null,
    location: f.location,
    description: f.description,
    starts_at: new Date(start).toISOString(),
    ends_at: end ? new Date(end).toISOString() : null,
  }
}
const eventReady = computed(() => {
  const f = editEvent.value
  return !!f && !!f.title.trim() && !!f.start && (!f.end || f.end > f.start)
})
function saveEvent() {
  const f = editEvent.value!
  save('events', { id: f.id, ...payload(f) }, () => (editEvent.value = null))
}
// 複製到隔天：用目前表單內容、日期 +1 天另存新行程
function copyNextDay() {
  const f = editEvent.value!
  save('events', payload(f, shiftLocal(f.start, 1440), f.end ? shiftLocal(f.end, 1440) : ''), () => (editEvent.value = null))
}
// 往下複製：新行程從目前的結束時間開始，時長相同
function copyBelow() {
  const f = editEvent.value!
  const mins = Math.round((new Date(f.end).getTime() - new Date(f.start).getTime()) / 60000)
  if (!f.end || mins <= 0) return Notify.create({ type: 'warning', message: '請先設定結束時間' })
  save('events', payload(f, f.end, shiftLocal(f.end, mins)), () => (editEvent.value = null))
}
// 清單上的複製鈕：時間往後 2 小時另存一筆
function duplicate(e: AssocEvent) {
  const f = { ...e, start: localInput(e.starts_at), end: localInput(e.ends_at), allDay: false }
  save('events', payload(f, shiftLocal(f.start, 120), f.end ? shiftLocal(f.end, 120) : ''), () => {})
}

// 活動標籤：像 hashtag，用來分類與篩選；刪掉標籤不會刪行程
const tagDialog = ref(false)
const tagName = ref('')
async function addTag() {
  if (!tagName.value.trim()) return
  try {
    await api.send('POST', '/api/admin/association/tags', { name: tagName.value })
    tagName.value = ''
    load()
  } catch (e) {
    toast(e)
  }
}
function removeTag(t: Tag) {
  Dialog.create({ title: '刪除標籤', message: `確定刪除「#${t.name}」？掛這個標籤的行程會保留，只是拿掉標籤。`, cancel: true }).onOk(async () => {
    try {
      await api.send('DELETE', `/api/admin/association/tags/${t.id}`)
      if (filterTag.value === t.id) filterTag.value = ''
      load()
    } catch (e) {
      toast(e)
    }
  })
}
</script>

<template>
  <q-page class="admin-page">
    <h1>世界佛教教育協會</h1>
    <p class="text-grey-8">這裡發布的通知、行事曆與會刊，協會會員登入官網後會在個人頁看到。會員和覺行小組共用同一個帳號，自己在個人頁選擇加入哪一邊。</p>

    <q-tabs v-model="tab" align="left" no-caps active-color="secondary" class="q-mb-md">
      <q-tab name="notices" label="協會通知" />
      <q-tab name="events" label="會員行事曆" />
      <q-tab name="issues" label="協會會刊" />
      <q-tab name="members" :label="`協會會員（${members.length}）`" />
    </q-tabs>

    <template v-if="tab === 'notices'">
      <div class="row q-mb-sm">
        <q-space />
        <q-btn color="secondary" unelevated no-caps icon="add" label="發布通知" @click="editNotice = { title: '', body: '' }" />
      </div>
      <q-table :rows="notices" :columns="noticeCols" row-key="id" flat bordered no-data-label="還沒有通知" :rows-per-page-options="[25, 50, 0]" @row-click="(_e: Event, row: Notice) => (editNotice = { ...row })" />
    </template>

    <template v-if="tab === 'events'">
      <div class="row q-gutter-sm q-mb-md">
        <q-btn color="secondary" unelevated no-caps icon="add" label="新增行程" @click="openEvent()" />
        <q-btn outline color="secondary" no-caps icon="sell" label="活動標籤" @click="tagDialog = true" />
      </div>
      <CalendarBoard v-model:view="calView" :events="calEvents" empty-text="還沒有行程" @select="(id) => openEvent(byId(id))">
        <template v-if="tags.length" #filters>
          <div class="row q-gutter-xs">
            <q-chip clickable :outline="filterTag !== ''" color="secondary" :text-color="filterTag === '' ? 'white' : 'secondary'" label="全部" @click="filterTag = ''" />
            <q-chip
              v-for="t in tags"
              :key="t.id"
              clickable
              :outline="filterTag !== t.id"
              color="secondary"
              :text-color="filterTag === t.id ? 'white' : 'secondary'"
              :label="`#${t.name}`"
              @click="filterTag = t.id"
            />
          </div>
        </template>
        <template #item="{ event }">
          <div class="row no-wrap items-start q-gutter-sm">
            <div class="col">
              <div class="row items-center q-gutter-xs">
                <span class="kind" :style="{ background: kindOf(byId(event.id).kind).color }">{{ kindOf(byId(event.id).kind).label }}</span>
                <span class="text-weight-bold">{{ event.title }}</span>
                <span v-if="event.badge" class="tag">{{ event.badge }}</span>
              </div>
              <div v-if="event.description" class="text-grey-8 pre q-mt-xs">{{ event.description }}</div>
              <div class="text-caption text-grey-7 q-mt-xs">
                {{ rangeText(event.start, event.end) }}<span v-if="event.location"> · <q-icon name="place" size="12px" />{{ event.location }}</span>
              </div>
            </div>
            <div class="row no-wrap" @click.stop>
              <q-btn flat dense round size="sm" icon="edit" aria-label="編輯" @click="openEvent(byId(event.id))" />
              <q-btn flat dense round size="sm" icon="content_copy" aria-label="複製（時間 +2 小時）" @click="duplicate(byId(event.id))"><q-tooltip>複製（時間 +2 小時）</q-tooltip></q-btn>
              <q-btn flat dense round size="sm" icon="delete" aria-label="刪除" @click="remove('events', event.id, event.title)" />
            </div>
          </div>
        </template>
      </CalendarBoard>
    </template>

    <template v-if="tab === 'issues'">
      <div class="row q-mb-sm">
        <q-space />
        <q-btn color="secondary" unelevated no-caps icon="add" label="新增會刊" @click="editIssue = { title: '', issued_on: today(), summary: '', url: '' }" />
      </div>
      <q-table :rows="issues" :columns="issueCols" row-key="id" flat bordered no-data-label="還沒有會刊" :rows-per-page-options="[25, 50, 0]" @row-click="(_e: Event, row: Issue) => (editIssue = { ...row })" />
    </template>

    <template v-if="tab === 'members'">
      <q-table :rows="members" :columns="memberCols" row-key="id" flat bordered no-data-label="還沒有協會會員" :rows-per-page-options="[50, 100, 0]" />
    </template>

    <q-dialog :model-value="!!editNotice" @update:model-value="editNotice = null">
      <q-card v-if="editNotice" style="width: 560px; max-width: 95vw">
        <q-card-section class="text-h6">{{ editNotice.id ? '編輯通知' : '發布通知' }}</q-card-section>
        <q-card-section class="q-gutter-md">
          <q-input v-model="editNotice.title" label="標題 *" outlined dense />
          <q-input v-model="editNotice.body" type="textarea" autogrow label="內容" outlined />
        </q-card-section>
        <q-card-actions>
          <q-btn v-if="editNotice.id" flat color="negative" no-caps label="刪除" @click="remove('notices', editNotice.id!, editNotice.title)" />
          <q-space />
          <q-btn flat no-caps label="取消" @click="editNotice = null" />
          <q-btn color="secondary" unelevated no-caps label="儲存" :disable="!editNotice.title.trim()" @click="save('notices', editNotice, () => (editNotice = null))" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <q-dialog :model-value="!!editEvent" @update:model-value="editEvent = null">
      <q-card v-if="editEvent" style="width: 620px; max-width: 95vw">
        <q-card-section class="text-h6">{{ editEvent.id ? '編輯行程' : '新增行程' }}</q-card-section>
        <q-card-section class="event-form">
          <div class="row q-gutter-sm">
            <q-btn-toggle v-model="editEvent.kind" no-caps unelevated dense toggle-color="secondary" class="kind-toggle" :options="KINDS.map((k) => ({ label: k.label, value: k.value }))" />
          </div>
          <div v-if="tags.length" class="row q-gutter-xs">
            <q-chip
              v-for="t in tags"
              :key="t.id"
              clickable
              :outline="editEvent.tag_id !== t.id"
              color="secondary"
              :text-color="editEvent.tag_id === t.id ? 'white' : 'secondary'"
              :icon="editEvent.tag_id === t.id ? 'check' : undefined"
              :label="`#${t.name}`"
              @click="editEvent.tag_id = editEvent.tag_id === t.id ? null : t.id"
            />
          </div>
          <q-input v-model="editEvent.title" label="行程名稱 *" outlined dense />
          <div class="time-grid">
            <DateTimeField :model-value="editEvent.start" :label="editEvent.allDay ? '日期 *' : '開始 *'" :hide-time="editEvent.allDay" @update:model-value="setStart" />
            <DateTimeField v-if="!editEvent.allDay" :model-value="editEvent.end" label="結束" @update:model-value="setEnd" />
          </div>
          <div v-if="editEvent.end && editEvent.end <= editEvent.start" class="text-negative text-caption">結束時間要晚於開始時間</div>
          <q-checkbox :model-value="editEvent.allDay" label="整天（08:00–20:00）" dense @update:model-value="setAllDay" />
          <q-input v-model="editEvent.location" label="地點" outlined dense />
          <q-input v-model="editEvent.description" type="textarea" autogrow label="說明" outlined />
        </q-card-section>
        <q-card-actions class="event-actions">
          <q-btn v-if="editEvent.id" flat color="negative" no-caps icon="delete" label="刪除" @click="remove('events', editEvent.id!, editEvent.title)" />
          <q-space />
          <q-btn flat no-caps label="取消" @click="editEvent = null" />
          <template v-if="editEvent.id">
            <q-btn outline color="secondary" no-caps icon="content_copy" :label="$q.screen.lt.sm ? '隔天' : '複製到隔天'" :disable="!eventReady" @click="copyNextDay"><q-tooltip>日期 +1 天另存一筆</q-tooltip></q-btn>
            <q-btn outline color="secondary" no-caps icon="content_copy" :label="$q.screen.lt.sm ? '往下' : '往下複製'" :disable="!eventReady" @click="copyBelow"><q-tooltip>從這筆的結束時間接著排一筆，時長相同</q-tooltip></q-btn>
          </template>
          <q-btn color="secondary" unelevated no-caps :label="editEvent.id ? '儲存' : '新增'" :disable="!eventReady" @click="saveEvent" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <q-dialog v-model="tagDialog">
      <q-card style="width: 520px; max-width: 95vw">
        <q-card-section class="text-h6">活動標籤</q-card-section>
        <q-card-section class="q-gutter-md">
          <p class="text-grey-8 q-mb-none">標籤像 hashtag，用來分類行程；新增行程時可以點選一個，行事曆上方也能用標籤篩選。</p>
          <q-form class="row no-wrap q-gutter-sm items-center" @submit.prevent="addTag">
            <q-input v-model="tagName" label="標籤名稱" outlined dense class="col" maxlength="40" />
            <q-btn type="submit" color="secondary" unelevated round icon="add" aria-label="新增標籤" :disable="!tagName.trim()" />
          </q-form>
          <p v-if="!tags.length" class="text-grey-7 text-center q-py-sm q-mb-none">還沒有標籤</p>
          <div v-else class="row q-gutter-xs">
            <q-chip v-for="t in tags" :key="t.id" removable color="secondary" text-color="white" :label="`#${t.name}`" @remove="removeTag(t)" />
          </div>
        </q-card-section>
        <q-card-actions align="right"><q-btn flat no-caps label="完成" @click="tagDialog = false" /></q-card-actions>
      </q-card>
    </q-dialog>

    <q-dialog :model-value="!!editIssue" @update:model-value="editIssue = null">
      <q-card v-if="editIssue" style="width: 560px; max-width: 95vw">
        <q-card-section class="text-h6">{{ editIssue.id ? '編輯會刊' : '新增會刊' }}</q-card-section>
        <q-card-section class="q-gutter-md">
          <q-input v-model="editIssue.title" label="會刊名稱 *" hint="例如：協會會刊 2026 年 10 月號" outlined dense />
          <q-input v-model="editIssue.issued_on" type="date" label="出刊日期 *" stack-label outlined dense />
          <q-input v-model="editIssue.url" label="連結" hint="線上版或 PDF 的完整網址" outlined dense />
          <q-input v-model="editIssue.summary" type="textarea" autogrow label="本期內容摘要" outlined />
        </q-card-section>
        <q-card-actions>
          <q-btn v-if="editIssue.id" flat color="negative" no-caps label="刪除" @click="remove('issues', editIssue.id!, editIssue.title)" />
          <q-space />
          <q-btn flat no-caps label="取消" @click="editIssue = null" />
          <q-btn color="secondary" unelevated no-caps label="儲存" :disable="!editIssue.title.trim() || !editIssue.issued_on" @click="save('issues', editIssue, () => (editIssue = null))" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<style scoped>
.kind {
  color: #fff;
  font-size: 0.75rem;
  border-radius: 4px;
  padding: 1px 6px;
}
.tag {
  font-size: 0.75rem;
  border: 1px solid #c9a86a;
  color: #8a6a32;
  background: #faf5ec;
  border-radius: 10px;
  padding: 0 6px;
}
.pre {
  white-space: pre-line;
}
.event-form {
  display: flex;
  flex-direction: column;
  gap: 14px;
}
.kind-toggle {
  border: 1px solid rgba(0, 0, 0, 0.18);
}
.time-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}
.event-actions {
  flex-wrap: wrap;
  gap: 6px;
}
@media (max-width: 599px) {
  .time-grid {
    grid-template-columns: 1fr;
  }
}
</style>
