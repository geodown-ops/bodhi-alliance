<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { Dialog, Notify } from 'quasar'
import { api } from '../session'

// 世界佛教教育協會：發布協會會刊、會員行事曆、協會通知，並查看協會會員名單。
type Issue = { id?: string; title: string; issued_on: string; summary: string; url: string }
type AssocEvent = { id?: string; title: string; starts_at: string; ends_at: string | null; location: string; description: string }
type Notice = { id?: string; title: string; body: string; created_at?: string }
type Member = { id: string; legal_name: string; display_name: string; email: string; phone: string; line_id: string; in_groups: boolean; joined_at: string | null }

const tab = ref<'notices' | 'events' | 'issues' | 'members'>('notices')
const issues = ref<Issue[]>([])
const events = ref<AssocEvent[]>([])
const notices = ref<Notice[]>([])
const members = ref<Member[]>([])
const editIssue = ref<Issue | null>(null)
const editEvent = ref<(AssocEvent & { start: string; end: string }) | null>(null)
const editNotice = ref<Notice | null>(null)

const toast = (e: unknown) => Notify.create({ type: 'negative', message: (e as Error).message })
const dt = (s: string | null | undefined) => (s ? new Date(s).toLocaleString('zh-TW', { hour12: false }) : '')
const d = (s: string | null | undefined) => (s ? new Date(s).toLocaleDateString('zh-TW') : '')

const pad = (n: number) => String(n).padStart(2, '0')
function localInput(s: string | null) {
  if (!s) return ''
  const x = new Date(s)
  return `${x.getFullYear()}-${pad(x.getMonth() + 1)}-${pad(x.getDate())}T${pad(x.getHours())}:${pad(x.getMinutes())}`
}

const issueCols = [
  { name: 'issued_on', label: '出刊日期', field: 'issued_on', sortable: true },
  { name: 'title', label: '會刊', field: 'title' },
  { name: 'url', label: '連結', field: 'url' },
]
const eventCols = [
  { name: 'starts_at', label: '開始', field: 'starts_at', format: dt, sortable: true },
  { name: 'title', label: '活動', field: 'title' },
  { name: 'location', label: '地點', field: 'location' },
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
    ;[issues.value, events.value, notices.value, members.value] = await Promise.all([
      api.get<Issue[]>('/api/admin/association/issues'),
      api.get<AssocEvent[]>('/api/admin/association/events'),
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

function saveEvent() {
  const e = editEvent.value!
  const { start, end, ...rest } = e
  save(
    'events',
    { ...rest, starts_at: new Date(start).toISOString(), ends_at: end ? new Date(end).toISOString() : null },
    () => (editEvent.value = null),
  )
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
const openEvent = (e?: AssocEvent) =>
  (editEvent.value = e
    ? { ...e, start: localInput(e.starts_at), end: localInput(e.ends_at) }
    : { title: '', starts_at: '', ends_at: null, location: '', description: '', start: '', end: '' })
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
      <div class="row q-mb-sm">
        <q-space />
        <q-btn color="secondary" unelevated no-caps icon="add" label="新增行事曆" @click="openEvent()" />
      </div>
      <q-table :rows="events" :columns="eventCols" row-key="id" flat bordered no-data-label="還沒有行事曆" :rows-per-page-options="[25, 50, 0]" @row-click="(_e: Event, row: AssocEvent) => openEvent(row)" />
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
      <q-card v-if="editEvent" style="width: 560px; max-width: 95vw">
        <q-card-section class="text-h6">{{ editEvent.id ? '編輯行事曆' : '新增行事曆' }}</q-card-section>
        <q-card-section class="q-gutter-md">
          <q-input v-model="editEvent.title" label="活動名稱 *" outlined dense />
          <div class="row q-col-gutter-sm">
            <q-input v-model="editEvent.start" type="datetime-local" label="開始 *" stack-label outlined dense class="col-12 col-sm-6" />
            <q-input v-model="editEvent.end" type="datetime-local" label="結束" stack-label outlined dense class="col-12 col-sm-6" />
          </div>
          <q-input v-model="editEvent.location" label="地點" outlined dense />
          <q-input v-model="editEvent.description" type="textarea" autogrow label="說明" outlined />
        </q-card-section>
        <q-card-actions>
          <q-btn v-if="editEvent.id" flat color="negative" no-caps label="刪除" @click="remove('events', editEvent.id!, editEvent.title)" />
          <q-space />
          <q-btn flat no-caps label="取消" @click="editEvent = null" />
          <q-btn color="secondary" unelevated no-caps label="儲存" :disable="!editEvent.title.trim() || !editEvent.start" @click="saveEvent" />
        </q-card-actions>
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
