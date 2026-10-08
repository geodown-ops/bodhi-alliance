<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { Notify } from 'quasar'
import { api, isAdmin } from '../session'

type Volunteer = {
  id: string
  email: string
  display_name: string
  legal_name: string
  phone: string
  line_id: string
  home_center_id: string
  center_name: string
  wants_coach: boolean
  is_coach: boolean
  is_group_leader: boolean
  status: 'pending' | 'verified' | 'rejected'
  review_note: string
  verified_at: string | null
  frozen: boolean
  created_at: string
  groups: { group_id: string; name: string; role: string }[]
}

const volunteers = ref<Volunteer[]>([])
const centers = ref<{ id: string; name: string }[]>([])
const status = ref('pending')
const centerId = ref('')
const editing = ref<(Volunteer & { next: Volunteer['status'] }) | null>(null)

const statusLabel: Record<string, string> = { pending: '待核可', verified: '已核可', rejected: '已退回' }
const statuses = [{ label: '全部', value: '' }, ...Object.entries(statusLabel).map(([value, label]) => ({ label, value }))]
const centerOptions = computed(() => [{ label: '全部中心', value: '' }, ...centers.value.map((c) => ({ label: c.name, value: c.id }))])

const columns = [
  { name: 'legal_name', label: '真實姓名', field: 'legal_name', sortable: true },
  { name: 'display_name', label: '暱稱', field: 'display_name' },
  { name: 'center_name', label: '所屬中心', field: (v: Volunteer) => v.center_name || '—', sortable: true },
  { name: 'line_id', label: 'LINE ID', field: 'line_id' },
  { name: 'is_coach', label: '教練', field: 'is_coach', align: 'center' as const },
  { name: 'is_group_leader', label: '覺行小組長', field: 'is_group_leader', align: 'center' as const },
  { name: 'status', label: '狀態', field: (v: Volunteer) => (v.frozen ? '已凍結' : statusLabel[v.status]) },
  { name: 'created_at', label: '註冊日', field: 'created_at', format: (v: string) => new Date(v).toLocaleDateString('zh-TW'), sortable: true },
]

const toast = (e: unknown) => Notify.create({ type: 'negative', message: (e as Error).message })

async function load() {
  try {
    volunteers.value = await api.get<Volunteer[]>(`/api/admin/volunteers?status=${status.value}&center_id=${centerId.value}`)
  } catch (e) {
    toast(e)
  }
}
onMounted(async () => {
  load()
  // 公開的中心清單：超級管理員用來篩選；中心管理員只會看到自己中心的志工
  try {
    centers.value = await api.get('/api/centers')
  } catch {
    /* 篩選用，失敗就不顯示 */
  }
})
watch([status, centerId], load)

function open(v: Volunteer) {
  editing.value = { ...v, next: v.status === 'pending' ? 'verified' : v.status }
}

// 名冊上直接打勾：教練、覺行小組長，打勾就表示有這個身分
async function setFlag(v: Volunteer, key: 'is_coach' | 'is_group_leader', on: boolean) {
  v[key] = on
  try {
    await api.send('PATCH', `/api/admin/volunteers/${v.id}`, { status: v.status, is_coach: v.is_coach, is_group_leader: v.is_group_leader, review_note: v.review_note, frozen: v.frozen })
  } catch (e) {
    v[key] = !on
    toast(e)
  }
}

async function save() {
  const v = editing.value
  if (!v) return
  try {
    await api.send('PATCH', `/api/admin/volunteers/${v.id}`, { status: v.next, is_coach: v.is_coach, is_group_leader: v.is_group_leader, review_note: v.review_note, frozen: v.frozen })
    editing.value = null
    load()
  } catch (e) {
    toast(e)
  }
}
</script>

<template>
  <q-page class="admin-page">
    <h1>會員名冊</h1>
    <p class="text-grey-8">
      會員在官網自己註冊並選所屬中心。請核對真實姓名與身分後按「核可」，核可後才能列入核發名單；核可後會員不能自己改姓名。教練、覺行小組長直接在表格打勾設定。
    </p>
    <div class="row q-gutter-sm q-mb-md">
      <q-select v-model="status" :options="statuses" emit-value map-options outlined dense label="狀態" style="min-width: 140px" />
      <q-select v-if="isAdmin()" v-model="centerId" :options="centerOptions" emit-value map-options outlined dense label="中心" style="min-width: 180px" />
    </div>
    <q-table :rows="volunteers" :columns="columns" row-key="id" flat bordered no-data-label="沒有符合的會員" @row-click="(_e: Event, row: Volunteer) => open(row)">
      <template #body-cell-is_coach="props">
        <q-td :props="props" @click.stop>
          <q-checkbox :model-value="props.row.is_coach" color="secondary" dense @update:model-value="(on: boolean) => setFlag(props.row, 'is_coach', on)" />
          <q-badge v-if="props.row.wants_coach && !props.row.is_coach" color="orange-2" text-color="brown-9" class="q-ml-xs">申請中</q-badge>
        </q-td>
      </template>
      <template #body-cell-is_group_leader="props">
        <q-td :props="props" @click.stop>
          <q-checkbox :model-value="props.row.is_group_leader" color="secondary" dense @update:model-value="(on: boolean) => setFlag(props.row, 'is_group_leader', on)" />
        </q-td>
      </template>
    </q-table>

    <q-dialog :model-value="!!editing" @update:model-value="editing = null">
      <q-card v-if="editing" style="width: 520px; max-width: 95vw">
        <q-card-section>
          <div class="text-h6">{{ editing.legal_name }}</div>
          <div class="text-grey-8">{{ editing.center_name }} · {{ editing.email }}{{ editing.phone ? ` · ${editing.phone}` : '' }}</div>
        </q-card-section>
        <q-card-section class="q-gutter-md">
          <q-btn-toggle
            v-model="editing.next"
            no-caps
            unelevated
            toggle-color="secondary"
            :options="[
              { label: '核可', value: 'verified' },
              { label: '待核可', value: 'pending' },
              { label: '退回', value: 'rejected' },
            ]"
          />
          <q-input
            v-model="editing.review_note"
            type="textarea"
            autogrow
            outlined
            :label="editing.next === 'rejected' ? '退回原因（會員會看到） *' : '備註（會員會看到）'"
          />
          <q-toggle v-model="editing.is_coach" :label="editing.wants_coach ? '禪修教練（本人有申請）' : '禪修教練'" />
          <q-toggle v-model="editing.is_group_leader" label="覺行小組長" />
          <q-toggle v-model="editing.frozen" color="negative" label="凍結（遺失手機或疑似冒用時，暫停兌換）" />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat no-caps label="取消" @click="editing = null" />
          <q-btn color="secondary" unelevated no-caps label="儲存" :disable="editing.next === 'rejected' && !editing.review_note.trim()" @click="save" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>
