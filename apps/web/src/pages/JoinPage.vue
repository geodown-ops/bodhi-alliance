<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ApiError } from '../api'
import { me, register } from '../account'

// 覺行小組報名：建立帳號後直接登入個人頁。從小組或活動的「報名」按鈕進來時，順便加入那個小組或活動。
const route = useRoute()
const router = useRouter()
const form = reactive({ legal_name: '', email: '', password: '', display_name: '', line_id: '', website: '' })
const sending = ref(false)
const error = ref('')

async function submit() {
  sending.value = true
  error.value = ''
  try {
    await register({ ...form })
  } catch (e) {
    error.value = e instanceof ApiError ? e.message : '報名失敗'
    sending.value = false
    return
  }
  // 帳號已經建好；加入小組或活動失敗時，到個人頁再加一次就好
  try {
    if (typeof route.query.group === 'string') await me.join(route.query.group)
    if (typeof route.query.event === 'string') await me.joinEvent(route.query.event, route.query.role === 'helper' ? 'helper' : 'participant')
  } catch {
    /* ignore */
  }
  router.push('/me')
}
</script>

<template>
  <q-page class="page narrow">
    <h1>報名參加覺行小組</h1>
    <p class="lead">留下資料就完成報名，之後用電子郵件和密碼登入個人頁：看錢包餘額、參加線上或線下的共修，也可以自己發起一場。</p>
    <p>已經報名過？<router-link :to="{ path: '/login', query: { next: '/me' } }">直接登入</router-link>。</p>

    <q-form class="card q-gutter-md" @submit.prevent="submit">
      <q-input v-model="form.legal_name" label="真實姓名 *" hint="核發菩提幣時核對身分用，不會公開" outlined :rules="[(v) => !!v.trim() || '請填寫真實姓名']" />
      <q-input v-model="form.email" type="email" label="電子郵件 *" hint="登入帳號" autocomplete="username" outlined :rules="[(v) => /.+@.+\..+/.test(v) || '請填寫正確的電子郵件']" />
      <q-input v-model="form.password" type="password" label="密碼 *" autocomplete="new-password" outlined :rules="[(v) => v.length >= 10 || '至少 10 個字元']" />
      <q-input v-model="form.display_name" label="暱稱" hint="活動頁上顯示的名字；不填就用真實姓名" outlined />
      <q-input v-model="form.line_id" label="LINE ID" hint="方便我們用 LINE 聯絡你" outlined />
      <input v-model="form.website" class="hp" tabindex="-1" autocomplete="off" aria-hidden="true" />
      <p v-if="error" class="text-negative q-mb-none">{{ error }}</p>
      <q-btn type="submit" color="secondary" unelevated no-caps size="lg" :loading="sending" label="送出報名" />
      <p class="text-caption q-mb-none">真實姓名、電子郵件與 LINE ID 只有管理員看得到。</p>
    </q-form>
  </q-page>
</template>

<style scoped>
.narrow {
  max-width: 640px;
}
.hp {
  position: absolute;
  left: -9999px;
}
</style>
