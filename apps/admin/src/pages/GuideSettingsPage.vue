<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { Notify } from 'quasar'
import { api } from '../session'

type Persona = { version: number; name: string; prompt: string; edited_by: string | null; created_at: string }
type Usage = { requests: number; cost_usd: number; budget_usd: number; enabled: boolean }

const history = ref<Persona[]>([])
const name = ref('')
const prompt = ref('')
const usage = ref<Usage | null>(null)
const configured = ref(true)
const budget = ref(0)
const enabled = ref(true)

async function load() {
  try {
    history.value = await api.guideGet<Persona[]>('/guide/admin/persona')
    if (history.value[0]) {
      name.value = history.value[0].name
      prompt.value = history.value[0].prompt
    }
    const u = await api.guideGet<{ usage: Usage; configured: boolean }>('/guide/admin/usage')
    usage.value = u.usage
    configured.value = u.configured
    budget.value = u.usage.budget_usd
    enabled.value = u.usage.enabled
  } catch (e) {
    Notify.create({ type: 'negative', message: (e as Error).message })
  }
}
onMounted(load)

async function savePersona() {
  try {
    await api.guideSend('PUT', '/guide/admin/persona', { name: name.value, prompt: prompt.value })
    Notify.create({ type: 'positive', message: '已儲存，下一則對話開始生效' })
    load()
  } catch (e) {
    Notify.create({ type: 'negative', message: (e as Error).message })
  }
}

async function saveSettings() {
  try {
    await api.guideSend('PUT', '/guide/admin/settings', { monthly_budget_usd: Number(budget.value), enabled: enabled.value })
    Notify.create({ type: 'positive', message: '已儲存' })
    load()
  } catch (e) {
    Notify.create({ type: 'negative', message: (e as Error).message })
  }
}

function restore(p: Persona) {
  name.value = p.name
  prompt.value = p.prompt
}
</script>

<template>
  <q-page class="admin-page">
    <h1>AI 組長設定</h1>

    <q-banner v-if="!configured" class="bg-orange-1 q-mb-md" rounded>
      伺服器還沒有設定 ANTHROPIC_API_KEY，官網上的 AI 組長目前不會回答。知識庫仍可以先整理。
    </q-banner>

    <div class="row q-col-gutter-lg">
      <div class="col-12 col-md-8">
        <q-card flat bordered>
          <q-card-section class="q-gutter-md">
            <div class="text-subtitle1">角色設定</div>
            <q-input v-model="name" label="名字" outlined dense />
            <q-input
              v-model="prompt"
              type="textarea"
              label="角色設定（自我介紹、語氣、界線）"
              outlined
              input-style="min-height: 360px; line-height: 1.8"
            />
            <div class="row">
              <q-space />
              <q-btn color="primary" unelevated no-caps label="儲存為新版本" :disable="!name.trim() || !prompt.trim()" @click="savePersona" />
            </div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-md-4">
        <q-card flat bordered class="q-mb-md">
          <q-card-section class="q-gutter-sm">
            <div class="text-subtitle1">本月用量</div>
            <div v-if="usage">{{ usage.requests }} 次對話 · 約 US${{ usage.cost_usd.toFixed(2) }}</div>
            <q-input v-model.number="budget" type="number" min="0" step="10" label="每月費用上限（美元）" outlined dense hint="超過上限時，官網會顯示組長休息中" />
            <q-toggle v-model="enabled" label="開放官網對話" />
            <q-btn color="secondary" unelevated no-caps label="儲存" @click="saveSettings" />
          </q-card-section>
        </q-card>
        <q-card flat bordered>
          <q-card-section>
            <div class="text-subtitle1 q-mb-sm">角色設定版本</div>
            <q-list dense separator>
              <q-item v-for="p in history" :key="p.version" clickable @click="restore(p)">
                <q-item-section>
                  <q-item-label>第 {{ p.version }} 版 · {{ p.name }}</q-item-label>
                  <q-item-label caption>{{ p.edited_by ?? '系統預設' }} · {{ new Date(p.created_at).toLocaleString('zh-TW') }}</q-item-label>
                </q-item-section>
              </q-item>
            </q-list>
            <p class="text-caption text-grey-7 q-mt-sm q-mb-none">點一個版本可載入編輯框，儲存後才會生效。</p>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>
