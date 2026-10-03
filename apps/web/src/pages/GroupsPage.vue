<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { api, ApiError, type Group } from '../api'

const groups = ref<Group[]>([])
const loadError = ref('')
const loading = ref(true)
const region = ref<string | null>(null)

onMounted(async () => {
  try {
    groups.value = await api.groups()
  } catch (e) {
    loadError.value = e instanceof ApiError ? e.message : '讀取失敗'
  } finally {
    loading.value = false
  }
})

const regions = computed(() => [...new Set(groups.value.map((g) => g.region))])
const shown = computed(() => groups.value.filter((g) => !region.value || g.region === region.value))
const groupOptions = computed(() => [
  { label: '還沒決定，請幫我安排', value: '' },
  ...groups.value.map((g) => ({ label: `${g.region}・${g.name}`, value: g.id })),
])

const form = reactive({ group_id: '', name: '', email: '', phone: '', region: '', wants_coach: false, message: '', website: '' })
const sending = ref(false)
const sent = ref(false)
const error = ref('')

function choose(g: Group) {
  form.group_id = g.id
  document.getElementById('join')?.scrollIntoView({ behavior: 'smooth' })
}

async function submit() {
  sending.value = true
  error.value = ''
  try {
    await api.joinGroup({ ...form })
    sent.value = true
  } catch (e) {
    error.value = e instanceof ApiError ? e.message : '送出失敗'
  } finally {
    sending.value = false
  }
}
</script>

<template>
  <q-page class="page">
    <h1>覺行小組介紹及參加</h1>
    <!-- 小組介紹文字待協會提供定稿 -->
    <p class="lead">覺行小組是在各中心與線上定期共修的小組。大家一起靜坐、讀經、分享，也一起在中心的活動中服務。</p>
    <p>
      參加小組不需要任何經驗。成為志工或禪修教練後，在中心活動中的服務會依梯級表核發
      <router-link to="/coin">菩提幣</router-link>。想先了解共修在做什麼，可以問問
      <router-link to="/guide">線上覺行小組 AI 組長</router-link>。
    </p>

    <div class="note q-mb-md">
      想在中心服務、領取菩提幣？先<router-link to="/join">註冊志工</router-link>，登入後可以直接在「<router-link to="/me">我的志工資料</router-link>」加入小組。
      只想參加共修，用下面的表單報名就好。
    </div>

    <h2>找一個小組</h2>
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
    <q-spinner v-if="loading" color="secondary" size="32px" />
    <p v-else-if="loadError" class="text-negative">{{ loadError }}</p>
    <p v-else-if="!groups.length">小組名單整理中。先留下資料，我們會為你安排離你最近的小組。</p>
    <div v-else class="grid">
      <div v-for="g in shown" :key="g.id" class="card">
        <div class="status-chip">{{ g.is_online ? '線上' : g.region }}</div>
        <h3 class="q-my-sm">{{ g.name }}</h3>
        <p v-if="g.center_name" class="q-mb-xs">{{ g.center_name }}</p>
        <p v-if="g.schedule" class="q-mb-xs">{{ g.schedule }}</p>
        <p v-if="g.description" class="q-mb-sm">{{ g.description }}</p>
        <q-btn flat color="secondary" no-caps label="報名這個小組" @click="choose(g)" />
      </div>
    </div>

    <h2 id="join">報名參加</h2>
    <div v-if="sent" class="note">
      <strong>已收到你的報名。</strong>小組組長會用電子郵件或電話和你聯絡。
    </div>
    <q-form v-else class="card q-gutter-md" @submit.prevent="submit">
      <q-select v-model="form.group_id" :options="groupOptions" emit-value map-options label="想參加的小組" outlined />
      <q-input v-model="form.name" label="姓名 *" outlined :rules="[(v) => !!v.trim() || '請填寫姓名']" />
      <q-input v-model="form.email" type="email" label="電子郵件 *" outlined :rules="[(v) => /.+@.+\..+/.test(v) || '請填寫正確的電子郵件']" />
      <q-input v-model="form.phone" label="聯絡電話" outlined />
      <q-input v-model="form.region" label="所在地區" outlined />
      <q-checkbox v-model="form.wants_coach" label="我有帶領禪修的經驗，想擔任禪修教練" />
      <q-input v-model="form.message" type="textarea" label="想說的話" outlined autogrow />
      <input v-model="form.website" class="hp" tabindex="-1" autocomplete="off" aria-hidden="true" />
      <p v-if="error" class="text-negative q-mb-none">{{ error }}</p>
      <q-btn type="submit" color="secondary" unelevated no-caps size="lg" :loading="sending" label="送出報名" />
      <p class="text-caption q-mb-none">資料僅供小組組長聯繫使用。</p>
    </q-form>
  </q-page>
</template>

<style scoped>
.hp {
  position: absolute;
  left: -9999px;
}
</style>
