<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Dialog, Notify } from 'quasar'
import { api } from '../session'

type Doc = { id: string; title: string; category: string; status: string; current_version: number; source_name: string; body: string }
type Version = { version: number; source_name: string; chars: number; edited_by: string | null; created_at: string }

const route = useRoute()
const router = useRouter()
const id = route.params.id as string
const doc = ref<Doc | null>(null)
const original = ref('')
const versions = ref<Version[]>([])
const categories = ref<Record<string, string>>({})
const saving = ref(false)

const dirty = computed(() => !!doc.value && JSON.stringify([doc.value.title, doc.value.category, doc.value.body]) !== original.value)

async function load() {
  try {
    const [d, v, list] = await Promise.all([
      api.guideGet<Doc>(`/guide/admin/documents/${id}`),
      api.guideGet<Version[]>(`/guide/admin/documents/${id}/versions`),
      api.guideGet<{ categories: Record<string, string> }>('/guide/admin/documents'),
    ])
    doc.value = d
    original.value = JSON.stringify([d.title, d.category, d.body])
    versions.value = v
    categories.value = list.categories
  } catch (e) {
    Notify.create({ type: 'negative', message: (e as Error).message })
  }
}
onMounted(load)

async function save() {
  if (!doc.value) return
  saving.value = true
  try {
    await api.guideSend('PUT', `/guide/admin/documents/${id}`, { title: doc.value.title, category: doc.value.category, body: doc.value.body })
    Notify.create({ type: 'positive', message: doc.value.status === 'published' ? '已儲存，AI 組長立即使用新內容' : '已儲存' })
    load()
  } catch (e) {
    Notify.create({ type: 'negative', message: (e as Error).message })
  } finally {
    saving.value = false
  }
}

async function setStatus(status: 'draft' | 'published' | 'archived') {
  try {
    await api.guideSend('POST', `/guide/admin/documents/${id}/status`, { status })
    if (status === 'archived') router.push('/knowledge')
    else load()
  } catch (e) {
    Notify.create({ type: 'negative', message: (e as Error).message })
  }
}

function archive() {
  Dialog.create({ title: '移除文件', message: 'AI 組長將不再使用這份文件。內容與版本紀錄會保留供稽核。', cancel: true }).onOk(() => setStatus('archived'))
}

const question = ref('')
const answer = ref('')
const asking = ref(false)
async function ask() {
  asking.value = true
  answer.value = ''
  try {
    const res = await api.guideSend<{ reply: string }>('POST', '/guide/admin/try', {
      include_document_id: id,
      messages: [{ role: 'user', content: question.value }],
    })
    answer.value = res.reply
  } catch (e) {
    answer.value = ''
    Notify.create({ type: 'negative', message: (e as Error).message })
  } finally {
    asking.value = false
  }
}
</script>

<template>
  <q-page v-if="doc" class="admin-page">
    <q-btn flat dense no-caps icon="arrow_back" label="知識庫" to="/knowledge" class="q-mb-sm" />
    <div class="row items-center q-gutter-sm q-mb-md">
      <h1 class="q-mb-none">{{ doc.title }}</h1>
      <q-badge :color="doc.status === 'published' ? 'secondary' : 'grey-6'">{{ doc.status === 'published' ? '已上架' : '草稿' }}</q-badge>
      <q-space />
      <q-btn v-if="doc.status !== 'published'" color="secondary" unelevated no-caps label="上架" :disable="dirty" @click="setStatus('published')" />
      <q-btn v-else outline color="secondary" no-caps label="下架" @click="setStatus('draft')" />
      <q-btn flat color="negative" no-caps label="移除" @click="archive" />
    </div>

    <div class="row q-col-gutter-lg">
      <div class="col-12 col-md-8 q-gutter-md">
        <q-input v-model="doc.title" label="標題" outlined dense />
        <q-select
          v-model="doc.category"
          :options="Object.entries(categories).map(([value, label]) => ({ value, label }))"
          emit-value
          map-options
          label="分類"
          outlined
          dense
        />
        <q-input v-model="doc.body" type="textarea" label="內容（可修正轉檔錯字；用 # 開頭的行當小標題）" outlined input-style="min-height: 420px; font-family: inherit; line-height: 1.8" />
        <div class="row items-center">
          <span class="text-grey-7">第 {{ doc.current_version }} 版 · 來源：{{ doc.source_name }}</span>
          <q-space />
          <q-btn color="primary" unelevated no-caps label="儲存" :disable="!dirty" :loading="saving" @click="save" />
        </div>
      </div>

      <div class="col-12 col-md-4">
        <q-card flat bordered class="q-mb-md">
          <q-card-section>
            <div class="text-subtitle1 q-mb-sm">試問</div>
            <p class="text-caption text-grey-8">用目前已上架的知識，加上這份文件（即使還是草稿），看看 AI 組長會怎麼回答。共修腳本會用「帶我共修」的方式試，可以輸入「開始」。</p>
            <q-input v-model="question" type="textarea" autogrow outlined dense placeholder="輸入一個問題" />
            <q-btn class="q-mt-sm" color="secondary" unelevated no-caps label="問問看" :loading="asking" :disable="!question.trim() || dirty" @click="ask" />
            <p v-if="dirty" class="text-caption text-grey-7 q-mt-sm">請先儲存再試問。</p>
            <div v-if="answer" class="answer q-mt-md">{{ answer }}</div>
          </q-card-section>
        </q-card>
        <q-card flat bordered>
          <q-card-section>
            <div class="text-subtitle1 q-mb-sm">版本紀錄</div>
            <q-list dense separator>
              <q-item v-for="v in versions" :key="v.version">
                <q-item-section>
                  <q-item-label>第 {{ v.version }} 版 · {{ v.chars }} 字</q-item-label>
                  <q-item-label caption>{{ v.edited_by ?? '—' }} · {{ new Date(v.created_at).toLocaleString('zh-TW') }} · {{ v.source_name }}</q-item-label>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<style scoped>
.answer {
  white-space: pre-wrap;
  line-height: 1.8;
  background: #ece5d6;
  border-radius: 8px;
  padding: 10px 12px;
}
</style>
