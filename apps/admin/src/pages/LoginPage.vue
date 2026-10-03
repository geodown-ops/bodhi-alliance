<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { login } from '../session'

const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)
const route = useRoute()
const router = useRouter()

async function submit() {
  loading.value = true
  error.value = ''
  try {
    await login(email.value, password.value)
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
          <q-input v-model="password" type="password" label="密碼" outlined autocomplete="current-password" />
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
