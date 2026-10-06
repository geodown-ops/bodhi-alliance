<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { Notify } from 'quasar'
import { api } from '../session'

// 覺行共修活動的菩提幣審核：發起人在活動結束後送審，超級管理員（代表菩提幣決策小組）核准後入帳。
type Recipient = { volunteer_id: string; legal_name: string; display_name: string; email: string; verified: boolean; role: 'organizer' | 'helper'; amount: number }
type Claim = {
  id: string
  status: 'submitted' | 'approved' | 'rejected'
  attendance: number
  report: string
  review_note: string
  created_at: string
  reviewed_at: string | null
  participants: number
  organizer_legal_name: string
  event: { title: string; is_online: boolean; location: string; starts_at: string; ends_at: string; capacity: number; organizer_name: string }
  recipients: Recipient[]
}

const claims = ref<Claim[]>([])
const status = ref('submitted')
const reviewing = ref<(Claim & { next: 'approved' | 'rejected'; note: string }) | null>(null)

const statusLabel: Record<string, string> = { submitted: '待審核', approved: '已核准', rejected: '已退回' }
const statuses = [{ label: '全部', value: '' }, ...Object.entries(statusLabel).map(([value, label]) => ({ label, value }))]
const fmt = new Intl.DateTimeFormat('zh-TW', { dateStyle: 'medium', timeStyle: 'short', hour12: false })
const hours = (c: Claim) => ((new Date(c.event.ends_at).getTime() - new Date(c.event.starts_at).getTime()) / 3600000).toFixed(1)
const total = (c: Claim) => c.recipients.reduce((n, r) => n + Number(r.amount || 0), 0)

async function load() {
  try {
    claims.value = await api.get<Claim[]>(`/api/admin/claims?status=${status.value}`)
  } catch (e) {
    Notify.create({ type: 'negative', message: (e as Error).message })
  }
}
onMounted(load)
watch(status, load)

function open(c: Claim) {
  reviewing.value = { ...JSON.parse(JSON.stringify(c)), next: 'approved', note: '' }
}

async function save() {
  const r = reviewing.value!
  try {
    await api.send('PATCH', `/api/admin/claims/${r.id}`, {
      status: r.next,
      review_note: r.note,
      amounts: Object.fromEntries(r.recipients.map((x) => [x.volunteer_id, Number(x.amount)])),
    })
    Notify.create({ type: 'positive', message: r.next === 'approved' ? '已核准，菩提幣已入帳' : '已退回，發起人會看到原因' })
    reviewing.value = null
    await load()
  } catch (e) {
    Notify.create({ type: 'negative', message: (e as Error).message })
  }
}
</script>

<template>
  <q-page class="admin-page">
    <h1>菩提幣審核</h1>
    <p class="text-grey-8">覺行共修活動結束後由發起人送審。核准後，發起人與協辦志工、減壓教練依下面的數量入帳；建議數量依梯級表（未滿四小時每小時 300、半日 1,500、全日 3,000），可以調整。</p>
    <q-btn-toggle v-model="status" no-caps unelevated toggle-color="secondary" :options="statuses" class="q-mb-md" />

    <p v-if="!claims.length" class="text-grey-8">沒有{{ statusLabel[status] ?? '' }}的申請。</p>
    <div class="column q-gutter-md">
      <q-card v-for="c in claims" :key="c.id" flat bordered>
        <q-card-section>
          <div class="row items-start">
            <div class="col">
              <div class="text-h6">{{ c.event.title }}</div>
              <div class="text-grey-8">
                {{ fmt.format(new Date(c.event.starts_at)) }} · {{ hours(c) }} 小時 · {{ c.event.is_online ? '線上' : c.event.location }}
              </div>
              <div class="text-grey-8">
                發起人 {{ c.organizer_legal_name }}（{{ c.event.organizer_name }}） · 報名 {{ c.participants }} 人 · 回報出席 {{ c.attendance }} 人
              </div>
            </div>
            <q-chip :color="c.status === 'submitted' ? 'orange-2' : c.status === 'approved' ? 'green-2' : 'grey-3'">{{ statusLabel[c.status] }}</q-chip>
          </div>
          <p v-if="c.report" class="q-mt-sm q-mb-none">活動紀錄：{{ c.report }}</p>
          <p v-if="c.review_note" class="q-mt-sm q-mb-none text-grey-8">審核說明：{{ c.review_note }}</p>
        </q-card-section>
        <q-markup-table flat dense>
          <tbody>
            <tr v-for="r in c.recipients" :key="r.volunteer_id">
              <td>{{ r.role === 'organizer' ? '發起人' : '協辦' }}</td>
              <td>{{ r.legal_name }}<span class="text-grey-7">（{{ r.display_name }}）</span></td>
              <td class="text-grey-8">{{ r.email }}</td>
              <td>{{ r.verified ? '已核對身分' : '身分未核對' }}</td>
              <td class="text-right">{{ r.amount.toLocaleString() }}</td>
            </tr>
          </tbody>
        </q-markup-table>
        <q-card-actions v-if="c.status === 'submitted'" align="right">
          <q-btn color="secondary" unelevated no-caps label="審核" @click="open(c)" />
        </q-card-actions>
      </q-card>
    </div>

    <q-dialog :model-value="!!reviewing" @update:model-value="reviewing = null">
      <q-card v-if="reviewing" style="width: 560px; max-width: 95vw">
        <q-card-section>
          <div class="text-h6">{{ reviewing.event.title }}</div>
          <div class="text-grey-8">回報出席 {{ reviewing.attendance }} 人 · {{ hours(reviewing) }} 小時</div>
        </q-card-section>
        <q-card-section class="q-gutter-md">
          <q-btn-toggle
            v-model="reviewing.next"
            no-caps
            unelevated
            toggle-color="secondary"
            :options="[
              { label: '核准並核發', value: 'approved' },
              { label: '退回', value: 'rejected' },
            ]"
          />
          <template v-if="reviewing.next === 'approved'">
            <q-input
              v-for="r in reviewing.recipients"
              :key="r.volunteer_id"
              v-model.number="r.amount"
              type="number"
              min="0"
              outlined
              dense
              :label="`${r.role === 'organizer' ? '發起人' : '協辦'} ${r.legal_name}`"
              suffix="菩提幣"
            />
            <div class="text-right">合計 {{ total(reviewing).toLocaleString() }} 菩提幣</div>
          </template>
          <q-input
            v-model="reviewing.note"
            type="textarea"
            autogrow
            outlined
            :label="reviewing.next === 'rejected' ? '退回原因（發起人會看到） *' : '備註（發起人會看到）'"
          />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat no-caps label="取消" @click="reviewing = null" />
          <q-btn color="secondary" unelevated no-caps label="送出" :disable="reviewing.next === 'rejected' && !reviewing.note.trim()" @click="save" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>
