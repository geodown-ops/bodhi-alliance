<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ApiError } from '../api'
import { me, register, type Center } from '../account'

const router = useRouter()
const centers = ref<Center[]>([])
const form = reactive({ legal_name: '', display_name: '', email: '', password: '', phone: '', home_center_id: '', wants_coach: false, website: '' })
const sending = ref(false)
const error = ref('')

onMounted(async () => {
  try {
    centers.value = await me.centers()
  } catch (e) {
    error.value = e instanceof ApiError ? e.message : '讀取中心清單失敗'
  }
})

async function submit() {
  sending.value = true
  error.value = ''
  try {
    await register({ ...form })
    router.push('/me')
  } catch (e) {
    error.value = e instanceof ApiError ? e.message : '註冊失敗'
  } finally {
    sending.value = false
  }
}
</script>

<template>
  <q-page class="page narrow">
    <h1>志工註冊</h1>
    <p class="lead">註冊後選你常去服務的中心，由中心的管理員核對身分。核可後，你在中心活動的服務就能依梯級表核發菩提幣。</p>
    <p>已經有帳號？<router-link to="/login">直接登入</router-link>。只想先參加共修，不用註冊，到<router-link to="/groups">覺行小組</router-link>頁面報名就好。</p>

    <q-form class="card q-gutter-md" @submit.prevent="submit">
      <q-input v-model="form.legal_name" label="真實姓名 *" hint="和身分證件相同，核發與到店核對時使用" outlined :rules="[(v) => !!v.trim() || '請填寫真實姓名']" />
      <q-input v-model="form.display_name" label="暱稱" hint="顯示在官網上的名字，可以不填" outlined />
      <q-input v-model="form.email" type="email" label="電子郵件 *" autocomplete="username" outlined :rules="[(v) => /.+@.+\..+/.test(v) || '請填寫正確的電子郵件']" />
      <q-input v-model="form.password" type="password" label="密碼 *" autocomplete="new-password" outlined :rules="[(v) => v.length >= 10 || '至少 10 個字元']" />
      <q-input v-model="form.phone" label="手機" outlined />
      <q-select
        v-model="form.home_center_id"
        :options="centers.map((c) => ({ label: c.region ? `${c.region}・${c.name}` : c.name, value: c.id }))"
        emit-value
        map-options
        label="所屬中心 *"
        hint="之後要換中心，請聯絡中心管理員"
        outlined
        :rules="[(v) => !!v || '請選所屬中心']"
      />
      <q-checkbox v-model="form.wants_coach" label="我有帶領禪修的經驗，想擔任禪修教練" />
      <input v-model="form.website" class="hp" tabindex="-1" autocomplete="off" aria-hidden="true" />
      <p v-if="error" class="text-negative q-mb-none">{{ error }}</p>
      <q-btn type="submit" color="secondary" unelevated no-caps size="lg" :loading="sending" label="註冊" />
      <p class="text-caption q-mb-none">真實姓名與手機只有你所屬中心的管理員看得到。</p>
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
