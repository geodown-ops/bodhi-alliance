<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { Dialog, Notify } from 'quasar'
import { api, session, type User } from '../session'

const users = ref<User[]>([])
const centers = ref<{ id: string; name: string }[]>([])
const centerName = (scope: string) => centers.value.find((c) => `center:${c.id}` === scope)?.name ?? ''
const roleLabel: Record<string, string> = { alliance_admin: '超級管理員', knowledge_manager: '知識管理員', center_admin: '中心管理員' }
const describe = (r: { role: string; scope: string }) =>
  r.role === 'center_admin' ? `${roleLabel[r.role]}（${centerName(r.scope)}）` : (roleLabel[r.role] ?? r.role)
const columns = [
  { name: 'display_name', label: '名稱', field: 'display_name' },
  { name: 'email', label: '電子郵件', field: 'email' },
  { name: 'roles', label: '角色', field: (u: User) => (u.roles.length ? u.roles.map(describe).join('、') : '志工') },
  { name: 'actions', label: '', field: 'id' },
]

async function load() {
  try {
    ;[users.value, centers.value] = await Promise.all([api.get<User[]>('/api/auth/users'), api.get<{ id: string; name: string }[]>('/api/admin/centers')])
  } catch (e) {
    Notify.create({ type: 'negative', message: (e as Error).message })
  }
}
onMounted(load)

const show = ref(false)
const form = reactive({ email: '', display_name: '', password: '', role: 'knowledge_manager', center_id: '' })
const centerOptions = computed(() => centers.value.map((c) => ({ label: c.name, value: c.id })))
async function create() {
  try {
    await api.send('POST', '/api/auth/users', { ...form, center_id: form.role === 'center_admin' ? form.center_id : undefined })
    show.value = false
    Object.assign(form, { email: '', display_name: '', password: '' })
    load()
  } catch (e) {
    Notify.create({ type: 'negative', message: (e as Error).message })
  }
}

function disable(u: User) {
  Dialog.create({ title: '停用帳號', message: `停用 ${u.display_name}（${u.email}）？對方會立刻被登出。`, cancel: true }).onOk(async () => {
    await api.send('DELETE', `/api/auth/users/${u.id}`)
    load()
  })
}
</script>

<template>
  <q-page class="admin-page">
    <div class="row items-center q-mb-md">
      <h1 class="q-mb-none">帳號</h1>
      <q-space />
      <q-btn color="secondary" unelevated no-caps icon="person_add" label="新增帳號" @click="show = true" />
    </div>
    <q-table :rows="users" :columns="columns" row-key="id" flat bordered hide-pagination :rows-per-page-options="[0]">
      <template #body-cell-actions="props">
        <q-td :props="props" class="text-right">
          <q-btn v-if="props.row.id !== session.user?.id" flat dense no-caps color="negative" label="停用" @click="disable(props.row)" />
        </q-td>
      </template>
    </q-table>

    <q-dialog v-model="show">
      <q-card style="width: 460px; max-width: 95vw">
        <q-card-section class="text-h6">新增帳號</q-card-section>
        <q-card-section class="q-gutter-md">
          <q-input v-model="form.display_name" label="名稱" outlined dense />
          <q-input v-model="form.email" type="email" label="電子郵件" outlined dense />
          <q-input v-model="form.password" type="password" label="初始密碼（至少 10 個字元）" outlined dense />
          <q-select
            v-model="form.role"
            :options="[
              { label: '知識管理員：只能管理 AI 組長與知識庫', value: 'knowledge_manager' },
              { label: '中心管理員：核可自己中心的志工、管理小組成員', value: 'center_admin' },
              { label: '超級管理員：全部功能', value: 'alliance_admin' },
            ]"
            emit-value
            map-options
            label="角色"
            outlined
            dense
          />
          <q-select
            v-if="form.role === 'center_admin'"
            v-model="form.center_id"
            :options="centerOptions"
            emit-value
            map-options
            label="管理哪個中心"
            outlined
            dense
          />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat no-caps label="取消" v-close-popup />
          <q-btn
            color="secondary"
            unelevated
            no-caps
            label="建立"
            :disable="!form.email || !form.display_name || form.password.length < 10 || (form.role === 'center_admin' && !form.center_id)"
            @click="create"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>
