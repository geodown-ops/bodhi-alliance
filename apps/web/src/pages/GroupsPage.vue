<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { Notify } from 'quasar'
import { api, ApiError, type Group } from '../api'
import { account, eventTime, me, publicEvents, type PracticeEvent } from '../account'

const router = useRouter()
const groups = ref<Group[]>([])
const events = ref<PracticeEvent[]>([])
const loadError = ref('')
const loading = ref(true)
const region = ref<string | null>(null)
const kind = ref<'all' | 'online' | 'offline'>('all')

onMounted(async () => {
  const [g, e] = await Promise.allSettled([api.groups(), publicEvents()])
  if (g.status === 'fulfilled') groups.value = g.value
  if (e.status === 'fulfilled') events.value = e.value
  if (g.status === 'rejected' && e.status === 'rejected') loadError.value = g.reason instanceof ApiError ? g.reason.message : '讀取失敗'
  loading.value = false
})

const regions = computed(() => [...new Set(groups.value.map((g) => g.region))])
const shown = computed(() => groups.value.filter((g) => !region.value || g.region === region.value))
const shownEvents = computed(() => events.value.filter((e) => kind.value === 'all' || e.is_online === (kind.value === 'online')))

const toast = (e: unknown) => Notify.create({ type: 'negative', message: e instanceof ApiError ? e.message : '發生錯誤' })

// 還沒報名的人先到報名頁，報名完成後自動加入；已登入就直接加入
async function joinGroup(g: Group) {
  if (!account.user) return router.push({ path: '/join', query: { group: g.id } })
  try {
    await me.join(g.id)
    router.push('/me')
  } catch (e) {
    toast(e)
  }
}

async function joinEvent(e: PracticeEvent, role: 'participant' | 'helper') {
  if (!account.user) return router.push({ path: '/join', query: { event: e.id, role } })
  try {
    await me.joinEvent(e.id, role)
    router.push('/me')
  } catch (err) {
    toast(err)
  }
}
</script>

<template>
  <q-page class="page">
    <h1>覺行小組介紹及參加</h1>
    <p class="lead">覺行小組是隨興或定期相約，一起進行正念減壓的小組活動。任何人都可以發起，只要三人以上，就能進行一次正念減壓實作。</p>
    <p>
      活動在菩提幣網站登錄後，協助的志工或減壓教練就可以得到
      <router-link to="/coin">菩提幣</router-link>。參加不需要任何經驗；想先了解正念減壓在做什麼，可以問問
      <router-link to="/guide">線上問答</router-link>。
    </p>

    <div class="card signup q-my-lg">
      <div>
        <h2 class="q-mt-none q-mb-xs">{{ account.user ? '到個人頁發起或參加活動' : '報名參加' }}</h2>
        <p class="q-mb-none">
          {{
            account.user
              ? '在個人頁可以看錢包餘額、參加線上或線下的覺行小組，也可以自己發起一場共修。'
              : '留下真實姓名、電子郵件和密碼就完成報名；之後登入個人頁，就能參加或發起共修活動。'
          }}
        </p>
      </div>
      <q-btn
        color="secondary"
        unelevated
        no-caps
        size="lg"
        :to="account.user ? '/me' : '/join'"
        :label="account.user ? '我的個人頁' : '報名'"
      />
    </div>

    <h2>近期共修活動</h2>
    <q-btn-toggle
      v-if="events.length"
      v-model="kind"
      no-caps
      unelevated
      toggle-color="secondary"
      class="q-mb-md"
      :options="[
        { label: '全部', value: 'all' },
        { label: '線下', value: 'offline' },
        { label: '線上', value: 'online' },
      ]"
    />
    <q-spinner v-if="loading" color="secondary" size="32px" />
    <p v-else-if="!shownEvents.length">目前沒有排定的活動。報名後就可以自己發起一場，找兩位以上的朋友一起練習。</p>
    <div v-else class="grid">
      <div v-for="e in shownEvents" :key="e.id" class="card">
        <div class="status-chip">{{ e.is_online ? '線上' : '線下' }}</div>
        <h3 class="q-my-sm">{{ e.title }}</h3>
        <p class="q-mb-xs">{{ eventTime(e) }}</p>
        <p class="q-mb-xs">{{ e.is_online ? '線上活動，報名後在個人頁看連結' : e.location }}</p>
        <p class="q-mb-xs text-caption">發起人 {{ e.organizer_name }} · 已報名 {{ e.joined }}／{{ e.capacity }} 人</p>
        <p v-if="e.description" class="q-mb-sm">{{ e.description }}</p>
        <div v-if="e.joined < e.capacity" class="row q-gutter-sm">
          <q-btn outline color="secondary" no-caps label="報名參加" @click="joinEvent(e, 'participant')" />
          <q-btn flat color="secondary" no-caps label="我來協辦" @click="joinEvent(e, 'helper')" />
        </div>
        <span v-else class="status-chip">已額滿</span>
      </div>
    </div>

    <h2>定期聚會的小組</h2>
    <q-select
      v-if="regions.length > 1"
      v-model="region"
      :options="regions"
      label="依地區篩選"
      clearable
      outlined
      dense
      class="q-mb-md"
      style="max-width: 240px"
    />
    <p v-if="loadError" class="text-negative">{{ loadError }}</p>
    <p v-else-if="!loading && !groups.length">小組名單整理中。</p>
    <div v-else class="grid">
      <div v-for="g in shown" :key="g.id" class="card">
        <div class="status-chip">{{ g.is_online ? '線上' : g.region }}</div>
        <h3 class="q-my-sm">{{ g.name }}</h3>
        <p v-if="g.center_name" class="q-mb-xs">{{ g.center_name }}</p>
        <p v-if="g.schedule" class="q-mb-xs">{{ g.schedule }}</p>
        <p v-if="g.description" class="q-mb-sm">{{ g.description }}</p>
        <q-btn flat color="secondary" no-caps label="報名這個小組" @click="joinGroup(g)" />
      </div>
    </div>
  </q-page>
</template>

<style scoped>
.signup {
  display: flex;
  gap: 16px;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
}
</style>
