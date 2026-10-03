<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { Dialog, Notify } from 'quasar'
import { api } from '../session'

type Center = {
  id?: string
  name: string
  region: string
  address: string
  contact_name: string
  email: string
  phone: string
  status: string
  note: string
  venue_count?: number
}
type Venue = {
  id?: string
  center_id: string
  center_name?: string
  name: string
  address: string
  description: string
  charges_public: boolean
  merchant_id: string | null
  merchant_name?: string | null
  status: string
}
type Merchant = { id: string; name: string; kind: string; center_id: string | null }

const centers = ref<Center[]>([])
const venues = ref<Venue[]>([])
const merchants = ref<Merchant[]>([])
const filter = ref<string | null>(null)
const editCenter = ref<Center | null>(null)
const editVenue = ref<Venue | null>(null)

const centerOptions = computed(() => centers.value.map((c) => ({ label: c.name, value: c.id })))
const shownVenues = computed(() => venues.value.filter((v) => !filter.value || v.center_id === filter.value))
// 場域只能連到同一個中心的聯盟單位
const merchantOptions = computed(() => [
  { label: '（無）', value: null },
  ...merchants.value
    .filter((m) => m.kind === 'alliance_unit' && m.center_id === editVenue.value?.center_id)
    .map((m) => ({ label: m.name, value: m.id })),
])

const centerCols = [
  { name: 'region', label: '地區', field: 'region', sortable: true },
  { name: 'name', label: '中心', field: 'name', sortable: true },
  { name: 'contact_name', label: '聯絡人', field: 'contact_name' },
  { name: 'phone', label: '電話', field: 'phone' },
  { name: 'venue_count', label: '場域數', field: 'venue_count' },
  { name: 'status', label: '狀態', field: 'status', format: (v: string) => (v === 'active' ? '運作中' : '暫停') },
]
const venueCols = [
  { name: 'center_name', label: '中心', field: 'center_name', sortable: true },
  { name: 'name', label: '場域', field: 'name', sortable: true },
  { name: 'address', label: '地址', field: 'address' },
  { name: 'charges_public', label: '對外收費', field: 'charges_public', format: (v: boolean) => (v ? '是' : '否') },
  { name: 'merchant_name', label: '對應聯盟單位', field: 'merchant_name', format: (v: string | null) => v ?? '' },
  { name: 'status', label: '狀態', field: 'status', format: (v: string) => (v === 'active' ? '使用中' : '已關閉') },
]

const toast = (e: unknown) => Notify.create({ type: 'negative', message: (e as Error).message })

async function load() {
  try {
    ;[centers.value, venues.value, merchants.value] = await Promise.all([
      api.get<Center[]>('/api/admin/centers'),
      api.get<Venue[]>('/api/admin/venues'),
      api.get<Merchant[]>('/api/admin/merchants'),
    ])
  } catch (e) {
    toast(e)
  }
}
onMounted(load)

async function saveCenter() {
  const c = editCenter.value
  if (!c) return
  try {
    await api.send(c.id ? 'PUT' : 'POST', c.id ? `/api/admin/centers/${c.id}` : '/api/admin/centers', c)
    editCenter.value = null
    load()
  } catch (e) {
    toast(e)
  }
}

async function saveVenue() {
  const v = editVenue.value
  if (!v) return
  try {
    await api.send(v.id ? 'PUT' : 'POST', v.id ? `/api/admin/venues/${v.id}` : '/api/admin/venues', v)
    editVenue.value = null
    load()
  } catch (e) {
    toast(e)
  }
}

function remove(kind: 'centers' | 'venues', id: string, name: string) {
  Dialog.create({ title: '刪除', message: `確定刪除「${name}」？`, cancel: true }).onOk(async () => {
    try {
      await api.send('DELETE', `/api/admin/${kind}/${id}`)
      editCenter.value = null
      editVenue.value = null
      load()
    } catch (e) {
      toast(e)
    }
  })
}

const newCenter = (): Center => ({ name: '', region: '', address: '', contact_name: '', email: '', phone: '', status: 'active', note: '' })
const newVenue = (): Venue => ({
  center_id: filter.value ?? centers.value[0]?.id ?? '',
  name: '',
  address: '',
  description: '',
  charges_public: false,
  merchant_id: null,
  status: 'active',
})
</script>

<template>
  <q-page class="admin-page">
    <h1>場域管理</h1>
    <p class="text-grey-8">
      禪修中心底下的場域（禪堂、營地、齋堂…）是舉辦活動、核發菩提幣的單位。場域如果也對外營業，可以連到同一個中心的聯盟單位，帳目仍分開。
    </p>

    <div class="row items-center q-mb-sm">
      <div class="text-h6">禪修中心</div>
      <q-space />
      <q-btn color="secondary" unelevated no-caps icon="add" label="新增中心" @click="editCenter = newCenter()" />
    </div>
    <q-table
      :rows="centers"
      :columns="centerCols"
      row-key="id"
      flat
      bordered
      no-data-label="還沒有中心"
      :rows-per-page-options="[0]"
      hide-pagination
      @row-click="(_e: Event, row: Center) => (editCenter = { ...row })"
    />

    <div class="row items-center q-mt-xl q-mb-sm q-gutter-sm">
      <div class="text-h6">場域</div>
      <q-select v-model="filter" :options="centerOptions" emit-value map-options clearable dense outlined label="依中心篩選" style="min-width: 200px" />
      <q-space />
      <q-btn color="secondary" unelevated no-caps icon="add" label="新增場域" :disable="!centers.length" @click="editVenue = newVenue()" />
    </div>
    <q-table
      :rows="shownVenues"
      :columns="venueCols"
      row-key="id"
      flat
      bordered
      :no-data-label="centers.length ? '還沒有場域' : '請先新增中心'"
      :rows-per-page-options="[25, 50, 0]"
      @row-click="(_e: Event, row: Venue) => (editVenue = { ...row })"
    />

    <q-dialog :model-value="!!editCenter" @update:model-value="editCenter = null">
      <q-card v-if="editCenter" style="width: 560px; max-width: 95vw">
        <q-card-section class="text-h6">{{ editCenter.id ? '編輯中心' : '新增中心' }}</q-card-section>
        <q-card-section class="q-gutter-md">
          <q-input v-model="editCenter.name" label="中心名稱 *" outlined dense />
          <q-input v-model="editCenter.region" label="地區" outlined dense />
          <q-input v-model="editCenter.address" label="地址" outlined dense />
          <div class="row q-col-gutter-sm">
            <q-input v-model="editCenter.contact_name" class="col-12 col-sm-4" label="聯絡人" outlined dense />
            <q-input v-model="editCenter.phone" class="col-12 col-sm-4" label="電話" outlined dense />
            <q-input v-model="editCenter.email" class="col-12 col-sm-4" label="電子郵件" outlined dense />
          </div>
          <q-input v-model="editCenter.note" type="textarea" autogrow label="備註" outlined />
          <q-toggle :model-value="editCenter.status === 'active'" label="運作中" @update:model-value="(v: boolean) => (editCenter!.status = v ? 'active' : 'suspended')" />
        </q-card-section>
        <q-card-actions>
          <q-btn v-if="editCenter.id" flat color="negative" no-caps label="刪除" @click="remove('centers', editCenter.id!, editCenter.name)" />
          <q-space />
          <q-btn flat no-caps label="取消" @click="editCenter = null" />
          <q-btn color="secondary" unelevated no-caps label="儲存" :disable="!editCenter.name.trim()" @click="saveCenter" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <q-dialog :model-value="!!editVenue" @update:model-value="editVenue = null">
      <q-card v-if="editVenue" style="width: 560px; max-width: 95vw">
        <q-card-section class="text-h6">{{ editVenue.id ? '編輯場域' : '新增場域' }}</q-card-section>
        <q-card-section class="q-gutter-md">
          <q-select v-model="editVenue.center_id" :options="centerOptions" emit-value map-options label="所屬中心 *" outlined dense @update:model-value="editVenue!.merchant_id = null" />
          <q-input v-model="editVenue.name" label="場域名稱 *" outlined dense />
          <q-input v-model="editVenue.address" label="地址" outlined dense />
          <q-input v-model="editVenue.description" type="textarea" autogrow label="說明" outlined />
          <q-select v-model="editVenue.merchant_id" :options="merchantOptions" emit-value map-options label="對應的聯盟單位（場域也對外營業時）" outlined dense />
          <q-toggle v-model="editVenue.charges_public" label="這個場域的活動對外收費" />
          <q-toggle :model-value="editVenue.status === 'active'" label="使用中" @update:model-value="(v: boolean) => (editVenue!.status = v ? 'active' : 'closed')" />
        </q-card-section>
        <q-card-actions>
          <q-btn v-if="editVenue.id" flat color="negative" no-caps label="刪除" @click="remove('venues', editVenue.id!, editVenue.name)" />
          <q-space />
          <q-btn flat no-caps label="取消" @click="editVenue = null" />
          <q-btn color="secondary" unelevated no-caps label="儲存" :disable="!editVenue.name.trim() || !editVenue.center_id" @click="saveVenue" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>
