<script setup lang="ts">
import PasswordInput from '../components/PasswordInput.vue'
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { homePath, login, logout } from '../session'

const email = ref('')
const password = ref('')
const route = useRoute()
const router = useRouter()
const denied = '這個帳號沒有後台權限。志工請到官網的「我的志工資料」登入。'
const error = ref(route.query.denied ? denied : '')
const loading = ref(false)
if (route.query.denied) logout()

async function submit() {
  loading.value = true
  error.value = ''
  try {
    await login(email.value, password.value)
    if (!homePath()) {
      await logout()
      error.value = denied
      return
    }
    router.replace((route.query.next as string) || '/')
  } catch (e) {
    error.value = (e as Error).message
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="wrap">
    <q-card class="login" flat bordered>
      <q-card-section class="row items-center q-gutter-sm">
        <img src="/favicon.svg" alt="" width="36" height="36" />
        <div class="text-h6">菩提幣管理後台</div>
      </q-card-section>
      <q-card-section>
        <q-form class="q-gutter-md" @submit.prevent="submit">
          <q-input v-model="email" type="email" label="電子郵件" outlined autocomplete="username" />
          <PasswordInput v-model="password" label="密碼" outlined autocomplete="current-password" />
          <p v-if="error" class="text-negative q-mb-none">{{ error }}</p>
          <q-btn type="submit" color="primary" unelevated no-caps class="full-width" label="登入" :loading="loading" />
        </q-form>
      </q-card-section>
    </q-card>
  </div>
</template>

<style scoped>
.wrap {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
}
.login {
  width: 100%;
  max-width: 380px;
}
</style>
