<script setup lang="ts">
import { computed, nextTick, onMounted, ref, watch } from 'vue'
import { api, ApiError, type ChatMessage } from '../api'
import GuideScene from '../components/GuideScene.vue'

const name = ref('Sunny')
const available = ref(true)
const mode = ref<'chat' | 'practice'>('chat')
const scripts = ref<{ id: string; title: string }[]>([])
const scriptId = ref<string | null>(null)
const messages = ref<ChatMessage[]>([])
const input = ref('')
const thinking = ref(false)
const streaming = ref('')
const error = ref('')
const subtitle = ref<HTMLElement>()

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

// 單一對話框只顯示最近一問一答，完整歷史仍送給組長當上下文。
const question = computed(() => [...messages.value].reverse().find((m) => m.role === 'user')?.content ?? '')
const answer = computed(() => {
  if (thinking.value) return streaming.value
  const last = messages.value[messages.value.length - 1]
  return last?.role === 'assistant' ? last.content : ''
})
const sceneState = computed(() => (thinking.value ? (streaming.value ? 'speaking' : 'listening') : 'idle'))

async function follow() {
  await nextTick()
  subtitle.value?.scrollTo({ top: subtitle.value.scrollHeight })
}

async function send(text = input.value) {
  text = text.trim()
  if (!text || thinking.value || !available.value) return
  if (mode.value === 'practice' && !scriptId.value) {
    error.value = '請先選一套共修流程'
    return
  }
  error.value = ''
  messages.value.push({ role: 'user', content: text })
  input.value = ''
  streaming.value = ''
  thinking.value = true
  try {
    const res = await api.chatStream(
      { messages: messages.value, mode: mode.value, script_id: scriptId.value ?? undefined },
      (t) => {
        streaming.value += t
        follow()
      },
    )
    messages.value.push({ role: 'assistant', content: res.reply })
  } catch (e) {
    messages.value.pop()
    input.value = text
    error.value = e instanceof ApiError ? e.message : '組長一時沒有回應，請再試一次。'
  } finally {
    thinking.value = false
    streaming.value = ''
    follow()
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
// 畫境佔滿頁首以下的整個畫面，頁尾留在捲動下方。
const pageHeight = () => ({ height: "calc(100svh - 50px)" })
</script>

<template>
  <q-page class="stage" :style-fn="pageHeight">
    <GuideScene :state="sceneState" />
    <h1 class="sr-only">線上覺行小組 AI 組長</h1>

    <section class="dialog" aria-label="和組長對話">
      <div class="modes">
        <button :class="{ on: mode === 'chat' }" type="button" @click="mode = 'chat'">問問組長</button>
        <button :class="{ on: mode === 'practice' }" type="button" @click="mode = 'practice'">帶我共修</button>
        <q-space />
        <button v-if="messages.length && !thinking" type="button" class="plain" @click="reset">重新開始</button>
      </div>

      <div ref="subtitle" class="subtitle" aria-live="polite">
        <template v-if="question">
          <div class="asked"><span class="label">你問</span>{{ question }}</div>
          <p class="answer">
            {{ answer }}<span v-if="thinking && !streaming" class="dots" aria-label="組長思考中"><i>．</i><i>．</i><i>．</i></span>
          </p>
        </template>
        <template v-else-if="!available">
          <p class="answer">線上組長目前休息中，請稍後再來，或到<router-link to="/groups">覺行小組</router-link>頁面聯絡真人組長。</p>
        </template>
        <template v-else-if="mode === 'practice'">
          <p class="answer">選一套共修流程，我一步一步帶你。每一步完成後回我「好了」，我們再往下走。</p>
          <div class="row items-center q-gutter-sm">
            <q-select
              v-model="scriptId"
              :options="scripts"
              option-value="id"
              option-label="title"
              emit-value
              map-options
              dense
              filled
              dark
              label="共修流程"
              class="col"
              :disable="!scripts.length"
            />
            <q-btn color="light-green-3" text-color="brown-10" unelevated rounded no-caps label="開始共修" :disable="!scriptId" @click="startPractice" />
          </div>
          <p v-if="!scripts.length" class="hint q-mt-sm q-mb-none">共修流程還在準備中。</p>
        </template>
        <template v-else>
          <p class="answer">我是線上覺行小組組長{{ name }}。可以問我覺行與共修的問題，也可以請我帶你共修一段。</p>
          <div class="chips">
            <button v-for="s in suggestions" :key="s" type="button" @click="send(s)">{{ s }}</button>
          </div>
        </template>
      </div>

      <form class="ask" @submit.prevent="send()">
        <input
          v-model="input"
          :maxlength="2000"
          :placeholder="mode === 'practice' ? '回應組長，例如「好了」「下一步」' : '輸入你的問題'"
          :disabled="!available"
          aria-label="輸入訊息"
        />
        <button type="submit" :disabled="!input.trim() || thinking || !available" aria-label="送出">
          <q-icon name="arrow_upward" size="22px" />
        </button>
      </form>
      <p v-if="error" class="error" role="alert">{{ error }}</p>
      <p class="hint">AI 組長只根據協會提供的資料回答，可能會出錯。重要的事請再向真人組長確認。</p>
    </section>
  </q-page>
</template>

<style scoped>
.stage {
  position: relative;
  overflow: hidden;
  min-height: 560px;
  display: flex;
  align-items: flex-end;
  justify-content: center;
  padding: 0 16px 20px;
}
.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  overflow: hidden;
  clip: rect(0 0 0 0);
}
.dialog {
  position: relative;
  width: 100%;
  max-width: 640px;
  padding: 14px 16px 10px;
  border-radius: 18px;
  background: rgba(59, 42, 32, 0.62);
  backdrop-filter: blur(14px);
  -webkit-backdrop-filter: blur(14px);
  color: #f6f2e8;
  box-shadow: 0 12px 40px rgba(59, 42, 32, 0.25);
}
.modes {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 8px;
}
.modes button {
  border: 1px solid rgba(246, 242, 232, 0.35);
  background: transparent;
  color: inherit;
  border-radius: 999px;
  padding: 3px 12px;
  font: inherit;
  font-size: 0.85rem;
  cursor: pointer;
}
.modes button.on {
  background: #b8d8a0;
  border-color: #b8d8a0;
  color: #3b2a20;
}
.modes button.plain {
  border: none;
  opacity: 0.8;
}
.subtitle {
  max-height: 38svh;
  overflow-y: auto;
  margin-bottom: 10px;
}
.asked {
  font-size: 0.9rem;
  opacity: 0.85;
  margin-bottom: 6px;
}
.label {
  color: #b8d8a0;
  font-weight: 600;
  margin-right: 8px;
}
.answer {
  font-size: 16px;
  line-height: 1.85;
  white-space: pre-wrap;
  margin: 0 0 8px;
}
.answer a {
  color: #b8d8a0;
}
.chips {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}
.chips button {
  border: 1px solid rgba(184, 216, 160, 0.6);
  background: rgba(184, 216, 160, 0.12);
  color: #e4f0da;
  border-radius: 999px;
  padding: 4px 12px;
  font: inherit;
  font-size: 0.9rem;
  cursor: pointer;
}
.dots i {
  display: inline-block;
  font-style: normal;
  animation: bounce 1.2s ease-in-out infinite;
}
.dots i:nth-child(2) {
  animation-delay: 0.15s;
}
.dots i:nth-child(3) {
  animation-delay: 0.3s;
}
.ask {
  display: flex;
  align-items: center;
  gap: 8px;
  background: #f6f2e8;
  border-radius: 999px;
  padding: 5px 5px 5px 18px;
}
.ask input {
  flex: 1;
  min-width: 0;
  border: none;
  outline: none;
  background: transparent;
  font: inherit;
  font-size: 16px;
  color: #3b2a20;
}
.ask button {
  width: 40px;
  height: 40px;
  flex: none;
  border: none;
  border-radius: 50%;
  background: #3b2a20;
  color: #b8d8a0;
  cursor: pointer;
  display: grid;
  place-items: center;
}
.ask button:disabled {
  opacity: 0.45;
  cursor: default;
}
.error {
  color: #ffc9b8;
  margin: 8px 4px 0;
}
.hint {
  font-size: 0.75rem;
  opacity: 0.7;
  margin: 8px 4px 0;
}
@media (max-width: 599px) {
  .subtitle {
    max-height: 30svh;
  }
}
@keyframes bounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-4px); }
}
@media (prefers-reduced-motion: reduce) {
  .dots i {
    animation: none;
  }
}
</style>
