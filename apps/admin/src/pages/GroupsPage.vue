<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { Dialog, Notify } from 'quasar'
import { api } from '../session'

type Group = {
  id?: string
  name: string
  region: string
  center_name: string
  schedule: string
  description: string
  is_online: boolean
  is_listed: boolean
  sort_order: number
}
const blank = (): Group => ({ name: '', region: '', center_name: '', schedule: '', description: '', is_online: false, is_listed: true, sort_order: 0 })

const groups = ref<Group[]>([])
const editing = ref<Group | null>(null)

const columns = [
  { name: 'region', label: '地區', field: 'region', sortable: true },
  { name: 'name', label: '小組', field: 'name', sortable: true },
  { name: 'center_name', label: '所屬中心', field: 'center_name' },
  { name: 'schedule', label: '共修時間', field: 'schedule' },
  { name: 'is_listed', label: '公開', field: 'is_listed', format: (v: boolean) => (v ? '是' : '否') },
]

async function load() {
  try {
    groups.value = await api.get<Group[]>('/api/admin/groups')
  } catch (e) {
    Notify.create({ type: 'negative', message: (e as Error).message })
  }
}
onMounted(load)

async function save() {
  const g = editing.value
  if (!g) return
  try {
    await api.send(g.id ? 'PUT' : 'POST', g.id ? `/api/admin/groups/${g.id}` : '/api/admin/groups', g)
    editing.value = null
    load()
  } catch (e) {
    Notify.create({ type: 'negative', message: (e as Error).message })
  }
}

function remove(g: Group) {
  Dialog.create({ title: '刪除小組', message: `確定刪除「${g.name}」？已報名的資料會保留，但不再對應到這個小組。`, cancel: true }).onOk(async () => {
    await api.send('DELETE', `/api/admin/groups/${g.id}`)
    editing.value = null
    load()
  })
}
</script>

<template>
  <q-page class="admin-page">
    <div class="row items-center q-mb-md">
      <h1 class="q-mb-none">覺行小組</h1>
      <q-space />
      <q-btn color="secondary" unelevated no-caps icon="add" label="新增小組" @click="editing = blank()" />
    </div>
    <q-table :rows="groups" :columns="columns" row-key="id" flat bordered no-data-label="還沒有小組" @row-click="(_e: Event, row: Group) => (editing = { ...row })" />

    <q-dialog :model-value="!!editing" @update:model-value="editing = null">
      <q-card v-if="editing" style="width: 560px; max-width: 95vw">
        <q-card-section class="text-h6">{{ editing.id ? '編輯小組' : '新增小組' }}</q-card-section>
        <q-card-section class="q-gutter-md">
          <q-input v-model="editing.name" label="小組名稱 *" outlined dense />
          <q-input v-model="editing.region" label="地區 *" outlined dense />
          <q-input v-model="editing.center_name" label="所屬中心" outlined dense />
          <q-input v-model="editing.schedule" label="共修時間（例如：每週三 19:30）" outlined dense />
          <q-input v-model="editing.description" type="textarea" autogrow label="介紹" outlined />
          <q-input v-model.number="editing.sort_order" type="number" label="排序（小的在前）" outlined dense />
          <q-toggle v-model="editing.is_online" label="線上小組" />
          <q-toggle v-model="editing.is_listed" label="在官網公開" />
        </q-card-section>
        <q-card-actions>
          <q-btn v-if="editing.id" flat color="negative" no-caps label="刪除" @click="remove(editing)" />
          <q-space />
          <q-btn flat no-caps label="取消" @click="editing = null" />
          <q-btn color="secondary" unelevated no-caps label="儲存" :disable="!editing.name || !editing.region" @click="save" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>
