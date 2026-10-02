<script setup lang="ts">
import { nextTick, onMounted, ref, watch } from 'vue'
import { api, ApiError, type ChatMessage } from '../api'

const name = ref('Sunny')
const available = ref(true)
const mode = ref<'chat' | 'practice'>('chat')
const scripts = ref<{ id: string; title: string }[]>([])
const scriptId = ref<string | null>(null)
const messages = ref<ChatMessage[]>([])
const input = ref('')
const thinking = ref(false)
const error = ref('')
const log = ref<HTMLElement>()

const storeKey = 'bodhi.guide.v1'
try {
  const saved = JSON.parse(sessionStorage.getItem(storeKey) ?? 'null')
  if (saved?.messages) {
    messages.value = saved.messages
    mode.value = saved.mode ?? 'chat'
    scriptId.value = saved.scriptId ?? null
  }
} catch {
  /* 無法使用 sessionStorage 時就不保留對話 */
}
watch([messages, mode, scriptId], () => {
  try {
    sessionStorage.setItem(storeKey, JSON.stringify({ messages: messages.value, mode: mode.value, scriptId: scriptId.value }))
  } catch {
    /* ignore */
  }
}, { deep: true })

onMounted(async () => {
  try {
    const info = await api.guideInfo()
    name.value = info.name
    available.value = info.available
    scripts.value = await api.scripts()
  } catch {
    available.value = false
  }
})

async function scrollDown() {
  await nextTick()
  log.value?.scrollTo({ top: log.value.scrollHeight, behavior: 'smooth' })
}

async function send(text = input.value) {
  text = text.trim()
  if (!text || thinking.value) return
  if (mode.value === 'practice' && !scriptId.value) {
    error.value = '請先選一套共修流程'
    return
  }
  error.value = ''
  messages.value.push({ role: 'user', content: text })
  input.value = ''
  thinking.value = true
  scrollDown()
  try {
    const res = await api.chat({ messages: messages.value, mode: mode.value, script_id: scriptId.value ?? undefined })
    messages.value.push({ role: 'assistant', content: res.reply })
  } catch (e) {
    messages.value.pop()
    input.value = text
    error.value = e instanceof ApiError ? e.message : '組長一時沒有回應，請再試一次。'
  } finally {
    thinking.value = false
    scrollDown()
  }
}

function startPractice() {
  messages.value = []
  send('我準備好了，請帶我開始共修。')
}

function reset() {
  messages.value = []
  error.value = ''
}

const suggestions = ['你是誰？', '覺行小組在做什麼？', '第一次靜坐要注意什麼？', '菩提幣怎麼拿到？']
</script>

<template>
  <q-page class="page guide">
    <h1>線上覺行小組 AI 組長</h1>
    <p class="lead">我是線上覺行小組組長{{ name }}。可以問我覺行與共修的問題，也可以請我帶你共修一段。</p>

    <div v-if="!available" class="note q-mb-md">
      線上組長目前休息中，請稍後再來，或到<router-link to="/groups">覺行小組</router-link>頁面聯絡真人組長。
    </div>

    <q-tabs v-model="mode" dense align="left" active-color="secondary" indicator-color="secondary" no-caps class="q-mb-md">
      <q-tab name="chat" label="問問組長" />
      <q-tab name="practice" label="帶我共修" />
    </q-tabs>

    <div v-if="mode === 'practice'" class="row items-center q-gutter-sm q-mb-md">
      <q-select
        v-model="scriptId"
        :options="scripts"
        option-value="id"
        option-label="title"
        emit-value
        map-options
        outlined
        dense
        label="選一套共修流程"
        style="min-width: 240px"
        :disable="!scripts.length"
      />
      <q-btn color="secondary" unelevated no-caps label="開始共修" :disable="!scriptId || thinking || !available" @click="startPractice" />
      <span v-if="!scripts.length" class="text-caption">共修流程還在準備中。</span>
    </div>

    <div ref="log" class="log card" aria-live="polite">
      <p v-if="!messages.length" class="text-caption q-mb-sm">試著問問看：</p>
      <div v-if="!messages.length" class="row q-gutter-sm">
        <q-chip v-for="s in suggestions" :key="s" clickable outline color="secondary" :disable="!available" @click="send(s)">{{ s }}</q-chip>
      </div>
      <div v-for="(m, i) in messages" :key="i" :class="['bubble', m.role]">
        <div class="who">{{ m.role === 'user' ? '我' : name }}</div>
        <div class="text">{{ m.content }}</div>
      </div>
      <div v-if="thinking" class="bubble assistant">
        <div class="who">{{ name }}</div>
        <q-spinner-dots color="secondary" size="28px" />
      </div>
    </div>

    <q-form class="row no-wrap items-end q-gutter-sm q-mt-md" @submit.prevent="send()">
      <q-input
        v-model="input"
        class="col"
        type="textarea"
        autogrow
        outlined
        :maxlength="2000"
        :placeholder="mode === 'practice' ? '回應組長，例如「好了」「下一步」' : '輸入你的問題'"
        :disable="!available"
        @keydown.enter.exact.prevent="send()"
      />
      <q-btn type="submit" color="secondary" unelevated round icon="send" aria-label="送出" :loading="thinking" :disable="!available" />
    </q-form>
    <div class="row items-center justify-between q-mt-sm">
      <p v-if="error" class="text-negative q-mb-none">{{ error }}</p>
      <q-space />
      <q-btn v-if="messages.length" flat dense no-caps color="grey-8" label="重新開始" @click="reset" />
    </div>
    <p class="text-caption q-mt-md">AI 組長只根據協會提供的資料回答，可能會出錯。重要的事請再向真人組長確認。</p>
  </q-page>
</template>

<style scoped>
.log {
  min-height: 320px;
  max-height: 60vh;
  overflow-y: auto;
}
.bubble {
  margin: 12px 0;
  max-width: 85%;
}
.bubble.user {
  margin-left: auto;
  text-align: right;
}
.who {
  font-size: 0.8rem;
  color: var(--ink-faint);
  margin-bottom: 2px;
}
.text {
  display: inline-block;
  text-align: left;
  white-space: pre-wrap;
  line-height: 1.8;
  padding: 10px 14px;
  border-radius: 12px;
  background: var(--ground-sunk);
}
.bubble.user .text {
  background: var(--leaf);
  color: #fff;
}
</style>
