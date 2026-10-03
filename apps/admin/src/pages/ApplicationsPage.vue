<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { Notify } from 'quasar'
import { api } from '../session'

type App = Record<string, unknown> & { id: string; status: string; admin_note: string; created_at: string }

const tab = ref<'group' | 'merchant'>('group')
const status = ref('')
const rows = ref<App[]>([])
const loading = ref(false)
const editing = ref<App | null>(null)

const statuses = [
  { label: '全部', value: '' },
  { label: '新進', value: 'new' },
  { label: '已聯絡', value: 'contacted' },
  { label: '已接受', value: 'accepted' },
  { label: '婉拒', value: 'declined' },
]
const statusLabel = (s: string) => statuses.find((x) => x.value === s)?.label ?? s
const kindLabel: Record<string, string> = { center: '禪修中心／道場', sponsor: '贊助商家', other: '其他' }

const groupCols = [
  { name: 'created_at', label: '時間', field: 'created_at', format: (v: string) => new Date(v).toLocaleString('zh-TW') },
  { name: 'name', label: '姓名', field: 'name' },
  { name: 'group_name', label: '小組', field: 'group_name', format: (v: string | null) => v ?? '待安排' },
  { name: 'email', label: '電子郵件', field: 'email' },
  { name: 'phone', label: '電話', field: 'phone' },
  { name: 'wants_coach', label: '教練', field: 'wants_coach', format: (v: boolean) => (v ? '是' : '') },
  { name: 'status', label: '狀態', field: 'status', format: statusLabel },
]
const merchantCols = [
  { name: 'created_at', label: '時間', field: 'created_at', format: (v: string) => new Date(v).toLocaleString('zh-TW') },
  { name: 'org_name', label: '單位', field: 'org_name' },
  { name: 'kind', label: '身份', field: 'kind', format: (v: string) => kindLabel[v] ?? v },
  { name: 'contact_name', label: '聯絡人', field: 'contact_name' },
  { name: 'email', label: '電子郵件', field: 'email' },
  { name: 'monthly_scale', label: '每月規模', field: 'monthly_scale' },
  { name: 'status', label: '狀態', field: 'status', format: statusLabel },
]

async function load() {
  loading.value = true
  try {
    const path = tab.value === 'group' ? 'group-applications' : 'merchant-applications'
    rows.value = await api.get<App[]>(`/api/admin/${path}?status=${status.value}`)
  } catch (e) {
    Notify.create({ type: 'negative', message: (e as Error).message })
  } finally {
    loading.value = false
  }
}
onMounted(load)
watch([tab, status], load)

async function save() {
  const a = editing.value
  if (!a) return
  const path = tab.value === 'group' ? 'group-applications' : 'merchant-applications'
  try {
    await api.send('PATCH', `/api/admin/${path}/${a.id}`, { status: a.status, admin_note: a.admin_note })
    editing.value = null
    load()
  } catch (e) {
    Notify.create({ type: 'negative', message: (e as Error).message })
  }
}

const centers = ref<{ id: string; name: string }[]>([])
const convertCenter = ref<string | null>(null)
onMounted(async () => {
  try {
    centers.value = await api.get('/api/admin/centers')
  } catch {
    /* 中心清單只影響「建立共好企業」的選項 */
  }
})

async function toMerchant() {
  const a = editing.value
  if (!a) return
  try {
    await api.send('POST', `/api/admin/merchant-applications/${a.id}/merchant`, convertCenter.value ? { center_id: convertCenter.value } : {})
    Notify.create({ type: 'positive', message: '已建立共好企業，狀態為審核中' })
    editing.value = null
    convertCenter.value = null
    load()
  } catch (e) {
    Notify.create({ type: 'negative', message: (e as Error).message })
  }
}

const details = (a: App) =>
  Object.entries(a).filter(([k, v]) => !['id', 'group_id', 'status', 'admin_note'].includes(k) && v !== '' && v !== null)
const fieldLabel: Record<string, string> = {
  created_at: '時間', name: '姓名', group_name: '小組', email: '電子郵件', phone: '電話', region: '地區',
  wants_coach: '想擔任教練', message: '留言', kind: '身份', org_name: '單位', contact_name: '聯絡人',
  offerings: '可提供品項', monthly_scale: '每月規模',
}
</script>

<template>
  <q-page class="admin-page">
    <h1>報名與登記</h1>
    <q-tabs v-model="tab" dense align="left" no-caps active-color="secondary" indicator-color="secondary">
      <q-tab name="group" label="覺行小組報名" />
      <q-tab name="merchant" label="共好企業登記" />
    </q-tabs>
    <q-btn-toggle v-model="status" :options="statuses" no-caps unelevated toggle-color="secondary" class="q-my-md" />
    <q-table
      :rows="rows"
      :columns="tab === 'group' ? groupCols : merchantCols"
      row-key="id"
      :loading="loading"
      flat
      bordered
      no-data-label="目前沒有資料"
      :rows-per-page-options="[25, 50, 0]"
      @row-click="(_e: Event, row: App) => (editing = { ...row })"
    />

    <q-dialog :model-value="!!editing" @update:model-value="editing = null">
      <q-card v-if="editing" style="width: 560px; max-width: 95vw">
        <q-card-section class="text-h6">{{ editing.name ?? editing.org_name }}</q-card-section>
        <q-card-section>
          <div v-for="[k, v] in details(editing)" :key="k" class="row q-mb-xs">
            <div class="col-4 text-grey-8">{{ fieldLabel[k] ?? k }}</div>
            <div class="col-8" style="white-space: pre-wrap">
              {{ k === 'created_at' ? new Date(v as string).toLocaleString('zh-TW') : k === 'kind' ? kindLabel[v as string] : v === true ? '是' : v }}
            </div>
          </div>
        </q-card-section>
        <q-card-section class="q-gutter-md">
          <q-select v-model="editing.status" :options="statuses.slice(1)" emit-value map-options label="狀態" outlined dense />
          <q-input v-model="editing.admin_note" type="textarea" autogrow label="內部備註" outlined />
        </q-card-section>
        <q-card-section v-if="tab === 'merchant'" class="q-gutter-sm">
          <div class="text-subtitle2">建立成共好企業</div>
          <q-select
            v-if="editing.kind === 'center'"
            v-model="convertCenter"
            :options="centers.map((c) => ({ label: c.name, value: c.id }))"
            emit-value
            map-options
            clearable
            outlined
            dense
            label="所屬中心（選了就建成聯盟單位）"
          />
          <q-btn outline color="secondary" no-caps icon="storefront" label="建立共好企業" @click="toMerchant" />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat no-caps label="取消" @click="editing = null" />
          <q-btn color="secondary" unelevated no-caps label="儲存" @click="save" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>
