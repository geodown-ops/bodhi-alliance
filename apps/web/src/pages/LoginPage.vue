<script setup lang="ts">
import PasswordInput from '../components/PasswordInput.vue'
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ApiError } from '../api'
import { login } from '../account'

const route = useRoute()
const router = useRouter()
const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

async function submit() {
  loading.value = true
  error.value = ''
  try {
    await login(email.value, password.value)
    router.replace((route.query.next as string) || '/me')
  } catch (e) {
    error.value = e instanceof ApiError ? e.message : '登入失敗'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <q-page class="page narrow">
    <h1>登入</h1>
    <q-form class="card card-form" @submit.prevent="submit">
      <q-input v-model="email" type="email" label="電子郵件" autocomplete="username" outlined />
      <PasswordInput v-model="password" label="密碼" autocomplete="current-password" outlined />
      <p v-if="error" class="text-negative q-mb-none">{{ error }}</p>
      <q-btn type="submit" color="secondary" unelevated no-caps size="lg" :loading="loading" label="登入" />
    </q-form>
    <p class="q-mt-md">還沒有帳號？<router-link to="/join">加入會員</router-link>。忘記密碼請聯絡管理員。</p>
  </q-page>
</template>

<style scoped>
.narrow {
  max-width: 480px;
}
</style>
