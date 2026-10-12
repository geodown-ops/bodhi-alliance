<script setup lang="ts">
import { reactive, ref } from 'vue'
import { Notify } from 'quasar'
import PasswordInput from './PasswordInput.vue'

// 修改密碼：管理員重設後拿到的臨時密碼，登入後在這裡換成自己的
const props = defineProps<{ save: (current: string, next: string) => Promise<unknown> }>()
const open = defineModel<boolean>({ required: true })
const form = reactive({ current: '', next: '', again: '' })
const saving = ref(false)
const error = ref('')

async function submit() {
  error.value = ''
  if (form.next.length < 10) return (error.value = '新密碼至少 10 個字元')
  if (form.next !== form.again) return (error.value = '兩次輸入的新密碼不一樣')
  saving.value = true
  try {
    await props.save(form.current, form.next)
    Object.assign(form, { current: '', next: '', again: '' })
    open.value = false
    Notify.create({ type: 'positive', message: '密碼已更新' })
  } catch (e) {
    error.value = (e as Error).message
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <q-dialog v-model="open">
    <q-card style="width: 420px; max-width: 95vw">
      <q-card-section class="text-h6">修改密碼</q-card-section>
      <q-form @submit.prevent="submit">
        <q-card-section class="column q-gutter-md">
          <PasswordInput v-model="form.current" label="目前的密碼（或臨時密碼）" outlined autocomplete="current-password" />
          <PasswordInput v-model="form.next" label="新密碼（至少 10 個字元）" outlined autocomplete="new-password" />
          <PasswordInput v-model="form.again" label="再輸入一次新密碼" outlined autocomplete="new-password" />
          <p v-if="error" class="text-negative q-mb-none">{{ error }}</p>
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat no-caps label="取消" @click="open = false" />
          <q-btn type="submit" color="secondary" unelevated no-caps label="儲存" :loading="saving" :disable="!form.current || !form.next || !form.again" />
        </q-card-actions>
      </q-form>
    </q-card>
  </q-dialog>
</template>
