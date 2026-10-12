<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import VolunteersPage from './VolunteersPage.vue'
import UsersPage from './UsersPage.vue'
import { isSuperAdmin } from '../session'

// 帳號：會員（原會員名冊）＋後台帳號（只有超級管理員看得到）
const route = useRoute()
const router = useRouter()
const tab = computed({
  get: () => (route.query.tab === 'staff' && isSuperAdmin() ? 'staff' : 'members'),
  set: (t: string) => router.replace({ query: t === 'staff' ? { tab: 'staff' } : {} }),
})
</script>

<template>
  <q-page class="admin-page">
    <h1>帳號</h1>
    <q-tabs v-if="isSuperAdmin()" v-model="tab" align="left" no-caps active-color="secondary" indicator-color="secondary" class="q-mb-md">
      <q-tab name="members" label="會員" />
      <q-tab name="staff" label="後台帳號" />
    </q-tabs>
    <VolunteersPage v-if="tab === 'members'" />
    <UsersPage v-else />
  </q-page>
</template>
