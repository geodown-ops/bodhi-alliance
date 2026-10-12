<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { Notify } from 'quasar'
import { api } from '../session'

// 使用問題：訪客或會員問 Sunny 網站怎麼用，Sunny 答不了或對方遇到狀況時記下來，由超級管理員回答。
// 會員在個人頁「我的使用問題」看到回答；沒登入的訪客若留了 email，請自行用 email 聯絡（網站目前不寄信）。
type Issue = {
  id: string
  user_id: string | null
  member_name: string
  member_email: string
  contact_email: string
  summary: string
  transcript: { role: 'user' | 'assistant'; content: string }[]
  status: 'open' | 'answered'
  answer: string
  answered_by: string
  answered_at: string | null
  created_at: string
}

const issues = ref<Issue[]>([])
const status = ref<'open' | 'answered' | ''>('open')
const answering = ref<(Issue & { draft: string }) | null>(null)
const statusLabel = { open: '待回答', answered: '已回答' }
const statuses = [
  { label: '待回答', value: 'open' },
  { label: '已回答', value: 'answered' },
  { label: '全部', value: '' },
]
const fmt = new Intl.DateTimeFormat('zh-TW', { dateStyle: 'medium', timeStyle: 'short', hour12: false })
const shown = computed(() => issues.value.filter((i) => !status.value || i.status === status.value))
const openCount = computed(() => issues.value.filter((i) => i.status === 'open').length)

async function load() {
  try {
    issues.value = await api.guideGet<Issue[]>('/guide/admin/issues')
  } catch (e) {
    Notify.create({ type: 'negative', message: (e as Error).message })
  }
}
onMounted(load)

function open(i: Issue) {
  answering.value = { ...i, draft: i.answer }
}

async function save() {
  const a = answering.value!
  try {
    await api.guideSend('PUT', `/guide/admin/issues/${a.id}/answer`, { answer: a.draft })
    Notify.create({ type: 'positive', message: a.user_id ? '已回答，會員在個人頁看得到' : '已回答' })
    answering.value = null
    await load()
  } catch (e) {
    Notify.create({ type: 'negative', message: (e as Error).message })
  }
}
</script>

<template>
  <q-page class="admin-page">
    <h1>使用問題</h1>
    <p class="text-grey-8">
      訪客或會員問 Sunny 網站怎麼用時，Sunny 答不了、或對方說遇到狀況，就會記在這裡。登入的會員會在個人頁「我的使用問題」看到你的回答；沒登入的訪客如果留了 email，請用 email 聯絡（網站目前不會自動寄信）。
    </p>
    <q-btn-toggle v-model="status" no-caps unelevated toggle-color="secondary" :options="statuses" class="q-mb-md" />
    <span v-if="openCount" class="q-ml-md text-grey-8">待回答 {{ openCount }} 則</span>

    <p v-if="!shown.length" class="text-grey-8">沒有{{ status ? statusLabel[status] : '' }}的問題。</p>
    <div class="column q-gutter-md">
      <q-card v-for="i in shown" :key="i.id" flat bordered>
        <q-card-section>
          <div class="row items-start">
            <div class="col">
              <div class="text-grey-8">
                {{ fmt.format(new Date(i.created_at)) }} ·
                <template v-if="i.user_id">會員 {{ i.member_name }}（{{ i.member_email }}）</template>
                <template v-else>訪客（未登入）</template>
                <template v-if="i.contact_email"> · 聯絡 email：<a :href="`mailto:${i.contact_email}`">{{ i.contact_email }}</a></template>
              </div>
              <div class="text-subtitle1 q-mt-xs pre">{{ i.summary }}</div>
            </div>
            <q-chip :color="i.status === 'open' ? 'orange-2' : 'green-2'">{{ statusLabel[i.status] }}</q-chip>
          </div>
          <q-expansion-item v-if="i.transcript.length" dense dense-toggle label="當時的對話" header-class="text-grey-8 q-px-none" class="q-mt-xs">
            <div v-for="(m, n) in i.transcript" :key="n" class="q-mb-xs pre">
              <b>{{ m.role === 'user' ? '對方' : 'Sunny' }}：</b>{{ m.content }}
            </div>
          </q-expansion-item>
          <div v-if="i.answer" class="answer q-mt-sm pre">
            <div class="text-caption text-grey-8">{{ i.answered_by }} 回答於 {{ i.answered_at ? fmt.format(new Date(i.answered_at)) : '' }}</div>
            {{ i.answer }}
          </div>
        </q-card-section>
        <q-card-actions align="right">
          <q-btn color="secondary" :outline="i.status === 'answered'" unelevated no-caps :label="i.status === 'answered' ? '修改回答' : '回答'" @click="open(i)" />
        </q-card-actions>
      </q-card>
    </div>

    <q-dialog :model-value="!!answering" @update:model-value="answering = null">
      <q-card v-if="answering" style="width: 560px; max-width: 95vw">
        <q-card-section>
          <div class="text-h6">回答使用問題</div>
          <div class="text-grey-8 pre q-mt-xs">{{ answering.summary }}</div>
        </q-card-section>
        <q-card-section>
          <q-input
            v-model="answering.draft"
            type="textarea"
            autogrow
            outlined
            autofocus
            :label="answering.user_id ? '回答（會員在個人頁看得到） *' : '回答（記錄用，請另外用 email 聯絡訪客） *'"
          />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat no-caps label="取消" @click="answering = null" />
          <q-btn color="secondary" unelevated no-caps label="送出" :disable="!answering.draft.trim()" @click="save" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<style scoped>
.pre {
  white-space: pre-line;
}
.answer {
  padding: 8px 12px;
  border-left: 3px solid var(--q-secondary);
  background: rgba(0, 0, 0, 0.03);
}
</style>
