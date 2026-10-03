<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { Dialog, Notify } from 'quasar'
import { api, ApiError, type Group } from '../api'
import { logout, me, type Volunteer } from '../account'

const router = useRouter()
const profile = ref<Volunteer | null>(null)
const groups = ref<Group[]>([])
const notVolunteer = ref(false)
const loading = ref(true)
const editing = ref(false)
const form = reactive({ display_name: '', legal_name: '', phone: '', wants_coach: false })

const toast = (e: unknown) => Notify.create({ type: 'negative', message: e instanceof ApiError ? e.message : '發生錯誤' })

async function load() {
  try {
    profile.value = await me.profile()
  } catch (e) {
    if (e instanceof ApiError && e.message.includes('還不是志工')) notVolunteer.value = true
    else toast(e)
  } finally {
    loading.value = false
  }
}
onMounted(async () => {
  await load()
  try {
    groups.value = await api.groups()
  } catch {
    /* 小組清單讀不到時只是不能加入 */
  }
})

const joined = computed(() => new Set(profile.value?.groups.map((g) => g.group_id)))
const statusText = computed(() => {
  const p = profile.value
  if (!p) return ''
  if (p.frozen) return '已凍結。如果不是你自己申請的，請聯絡中心管理員。'
  return { pending: `等待${p.center_name}的管理員核對身分。`, verified: '已核可。你在中心活動的服務可以列入核發名單。', rejected: '中心管理員退回了你的資料，請依說明修改後儲存，會重新送審。' }[p.status]
})

function startEdit() {
  const p = profile.value!
  Object.assign(form, { display_name: p.display_name, legal_name: p.legal_name, phone: p.phone, wants_coach: p.wants_coach })
  editing.value = true
}

async function save() {
  try {
    profile.value = await me.update({ ...form })
    editing.value = false
  } catch (e) {
    toast(e)
  }
}

async function join(g: Group) {
  try {
    await me.join(g.id)
    await load()
  } catch (e) {
    toast(e)
  }
}

function leave(g: { group_id: string; name: string }) {
  Dialog.create({ title: '退出小組', message: `確定退出「${g.name}」？`, cancel: true }).onOk(async () => {
    try {
      await me.leave(g.group_id)
      await load()
    } catch (e) {
      toast(e)
    }
  })
}

async function signOut() {
  await logout()
  router.push('/')
}
</script>

<template>
  <q-page class="page">
    <div class="row items-center q-mb-md">
      <h1 class="q-mb-none">我的志工資料</h1>
      <q-space />
      <q-btn flat no-caps icon="logout" label="登出" @click="signOut" />
    </div>
    <q-spinner v-if="loading" color="secondary" size="32px" />

    <div v-else-if="notVolunteer" class="note">
      這個帳號還沒有志工資料。請用<router-link to="/join">志工註冊</router-link>另外建立帳號，或聯絡中心管理員。
    </div>

    <template v-else-if="profile">
      <div class="card">
        <div class="row items-start">
          <div class="col">
            <h2 class="q-mt-none q-mb-xs">{{ profile.display_name }}</h2>
            <div>{{ profile.legal_name }} · {{ profile.center_name }}<span v-if="profile.is_coach"> · 禪修教練</span></div>
            <div class="text-caption">{{ profile.email }}{{ profile.phone ? ` · ${profile.phone}` : '' }}</div>
          </div>
          <q-btn v-if="!editing" flat no-caps color="secondary" icon="edit" label="修改" @click="startEdit" />
        </div>
        <p :class="['q-mt-md', 'q-mb-none', profile.status === 'rejected' ? 'text-negative' : '']">
          <span class="status-chip q-mr-sm">{{ { pending: '待核可', verified: '已核可', rejected: '已退回' }[profile.status] }}</span>{{ statusText }}
        </p>
        <p v-if="profile.review_note" class="q-mt-sm q-mb-none">中心的說明：{{ profile.review_note }}</p>

        <q-form v-if="editing" class="q-gutter-md q-mt-md" @submit.prevent="save">
          <q-input v-model="form.display_name" label="暱稱" outlined dense />
          <q-input v-model="form.legal_name" label="真實姓名" outlined dense :disable="profile.status === 'verified'" :hint="profile.status === 'verified' ? '核可後要改姓名，請聯絡中心管理員' : ''" />
          <q-input v-model="form.phone" label="手機" outlined dense />
          <q-checkbox v-model="form.wants_coach" label="我想擔任禪修教練" />
          <div class="row q-gutter-sm">
            <q-btn type="submit" color="secondary" unelevated no-caps label="儲存" />
            <q-btn flat no-caps label="取消" @click="editing = false" />
          </div>
        </q-form>
      </div>

      <h2>我的覺行小組</h2>
      <p v-if="!profile.groups.length">還沒加入小組。從下面挑一個加入，小組組長會和你聯絡。</p>
      <div v-else class="grid q-mb-md">
        <div v-for="g in profile.groups" :key="g.group_id" class="card">
          <h3 class="q-my-none">{{ g.name }}</h3>
          <p v-if="g.role === 'leader'" class="q-my-xs"><span class="status-chip">組長</span></p>
          <q-btn flat dense no-caps color="grey-8" label="退出" @click="leave(g)" />
        </div>
      </div>
      <template v-if="groups.some((g) => !joined.has(g.id))">
        <h3>可以加入的小組</h3>
        <div class="grid">
          <div v-for="g in groups.filter((g) => !joined.has(g.id))" :key="g.id" class="card">
            <div class="status-chip">{{ g.is_online ? '線上' : g.region }}</div>
            <h3 class="q-my-sm">{{ g.name }}</h3>
            <p v-if="g.schedule" class="q-mb-sm">{{ g.schedule }}</p>
            <q-btn outline color="secondary" no-caps label="加入" @click="join(g)" />
          </div>
        </div>
      </template>

      <h2>服務紀錄與菩提幣</h2>
      <p>服務時數登錄與核發開放後，你的服務紀錄、餘額與券會顯示在這裡。<router-link to="/wallet">錢包說明</router-link></p>
    </template>
  </q-page>
</template>
