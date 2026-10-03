<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { Dialog, Notify } from 'quasar'
import { api } from '../session'

type Merchant = {
  id?: string
  kind: 'alliance_unit' | 'sponsor'
  center_id: string | null
  center_name?: string | null
  name: string
  contact_name: string
  email: string
  phone: string
  region: string
  address: string
  offerings: string
  proposed_monthly_cap: number | null
  status: string
  application_id?: string | null
  note: string
}

const merchants = ref<Merchant[]>([])
const centers = ref<{ id: string; name: string }[]>([])
const status = ref('')
const editing = ref<Merchant | null>(null)

const kindLabel = { alliance_unit: '聯盟單位', sponsor: '贊助商家' }
const statusLabel: Record<string, string> = { pending: '審核中', active: '已上架', suspended: '暫停' }
const statuses = [{ label: '全部', value: '' }, ...Object.entries(statusLabel).map(([value, label]) => ({ label, value }))]
const centerOptions = computed(() => centers.value.map((c) => ({ label: c.name, value: c.id })))

const columns = [
  { name: 'name', label: '企業', field: 'name', sortable: true },
  { name: 'kind', label: '類型', field: 'kind', format: (v: 'alliance_unit' | 'sponsor') => kindLabel[v] },
  { name: 'center_name', label: '所屬中心', field: 'center_name', format: (v: string | null) => v ?? '' },
  { name: 'contact_name', label: '聯絡人', field: 'contact_name' },
  { name: 'region', label: '地區', field: 'region' },
  { name: 'proposed_monthly_cap', label: '提出的每月規模（幣）', field: 'proposed_monthly_cap', format: (v: number | null) => (v == null ? '' : v.toLocaleString()) },
  { name: 'status', label: '狀態', field: 'status', format: (v: string) => statusLabel[v] ?? v, sortable: true },
]

const toast = (e: unknown) => Notify.create({ type: 'negative', message: (e as Error).message })

async function load() {
  try {
    ;[merchants.value, centers.value] = await Promise.all([
      api.get<Merchant[]>(`/api/admin/merchants?status=${status.value}`),
      api.get<{ id: string; name: string }[]>('/api/admin/centers'),
    ])
  } catch (e) {
    toast(e)
  }
}
onMounted(load)
watch(status, load)

async function save() {
  const m = editing.value
  if (!m) return
  try {
    await api.send(m.id ? 'PUT' : 'POST', m.id ? `/api/admin/merchants/${m.id}` : '/api/admin/merchants', {
      ...m,
      center_id: m.kind === 'alliance_unit' ? m.center_id : null,
      proposed_monthly_cap: m.proposed_monthly_cap === null || (m.proposed_monthly_cap as unknown) === '' ? null : Number(m.proposed_monthly_cap),
    })
    editing.value = null
    load()
  } catch (e) {
    toast(e)
  }
}

function remove(m: Merchant) {
  Dialog.create({ title: '刪除共好企業', message: `確定刪除「${m.name}」？`, cancel: true }).onOk(async () => {
    try {
      await api.send('DELETE', `/api/admin/merchants/${m.id}`)
      editing.value = null
      load()
    } catch (e) {
      toast(e)
    }
  })
}

const blank = (): Merchant => ({
  kind: 'sponsor',
  center_id: null,
  name: '',
  contact_name: '',
  email: '',
  phone: '',
  region: '',
  address: '',
  offerings: '',
  proposed_monthly_cap: null,
  status: 'pending',
  note: '',
})
</script>

<template>
  <q-page class="admin-page">
    <div class="row items-center q-mb-sm">
      <h1 class="q-mb-none">共好企業管理</h1>
      <q-space />
      <q-btn color="secondary" unelevated no-caps icon="add" label="新增共好企業" @click="editing = blank()" />
    </div>
    <p class="text-grey-8">
      官網的共好企業登記可以在「報名與登記」一鍵建立成這裡的資料。正式上架與每月額度要等主辦審核小組決議；在那之前，這裡記的是企業自己提出的規模。
    </p>
    <q-btn-toggle v-model="status" :options="statuses" no-caps unelevated toggle-color="secondary" class="q-mb-md" />
    <q-table
      :rows="merchants"
      :columns="columns"
      row-key="id"
      flat
      bordered
      no-data-label="還沒有共好企業"
      :rows-per-page-options="[25, 50, 0]"
      @row-click="(_e: Event, row: Merchant) => (editing = { ...row })"
    />

    <q-dialog :model-value="!!editing" @update:model-value="editing = null">
      <q-card v-if="editing" style="width: 600px; max-width: 95vw">
        <q-card-section class="text-h6">{{ editing.id ? '編輯共好企業' : '新增共好企業' }}</q-card-section>
        <q-card-section class="q-gutter-md">
          <q-btn-toggle
            v-model="editing.kind"
            :options="[{ label: '贊助商家', value: 'sponsor' }, { label: '聯盟單位（屬於某中心）', value: 'alliance_unit' }]"
            no-caps
            unelevated
            toggle-color="secondary"
          />
          <q-select
            v-if="editing.kind === 'alliance_unit'"
            v-model="editing.center_id"
            :options="centerOptions"
            emit-value
            map-options
            label="所屬中心 *"
            outlined
            dense
            :hint="centers.length ? '' : '請先到「場域管理」新增中心'"
          />
          <q-input v-model="editing.name" label="企業名稱 *" outlined dense />
          <div class="row q-col-gutter-sm">
            <q-input v-model="editing.contact_name" class="col-12 col-sm-4" label="聯絡人" outlined dense />
            <q-input v-model="editing.phone" class="col-12 col-sm-4" label="電話" outlined dense />
            <q-input v-model="editing.email" class="col-12 col-sm-4" label="電子郵件" outlined dense />
          </div>
          <div class="row q-col-gutter-sm">
            <q-input v-model="editing.region" class="col-12 col-sm-4" label="地區" outlined dense />
            <q-input v-model="editing.address" class="col-12 col-sm-8" label="地址" outlined dense />
          </div>
          <q-input v-model="editing.offerings" type="textarea" autogrow label="可提供的品項" outlined />
          <q-input v-model.number="editing.proposed_monthly_cap" type="number" min="0" label="提出的每月贊助規模（幣）" outlined dense clearable />
          <q-select v-model="editing.status" :options="statuses.slice(1)" emit-value map-options label="狀態" outlined dense />
          <q-input v-model="editing.note" type="textarea" autogrow label="內部備註" outlined />
        </q-card-section>
        <q-card-actions>
          <q-btn v-if="editing.id" flat color="negative" no-caps label="刪除" @click="remove(editing)" />
          <q-space />
          <q-btn flat no-caps label="取消" @click="editing = null" />
          <q-btn
            color="secondary"
            unelevated
            no-caps
            label="儲存"
            :disable="!editing.name.trim() || (editing.kind === 'alliance_unit' && !editing.center_id)"
            @click="save"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>
