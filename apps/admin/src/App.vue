<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ChangePasswordDialog from './components/ChangePasswordDialog.vue'
import { api, isAdmin, isCenterStaff, isKnowledgeManager, isSuperAdmin, logout, session } from './session'

const route = useRoute()
const router = useRouter()
const drawer = ref(false)
const changingPassword = ref(false)
const changePassword = (current: string, next: string) =>
  api.send('PUT', '/api/auth/password', { current_password: current, new_password: next })

const items = computed(() => [
  ...(isAdmin()
    ? [
        { to: '/applications', label: '報名與登記', icon: 'inbox' },
        { to: '/groups', label: '覺行小組', icon: 'groups' },
        ...(isSuperAdmin()
          ? [
              { to: '/claims', label: '菩提幣審核', icon: 'paid' },
              { to: '/site-issues', label: '使用問題', icon: 'support_agent' },
            ]
          : []),
        { to: '/venues', label: '場域管理', icon: 'place' },
        { to: '/merchants', label: '共好企業管理', icon: 'storefront' },
        { to: '/association', label: '世界佛教教育協會', icon: 'img:/images/association/logo.webp' },
      ]
    : []),
  ...(isCenterStaff() ? [{ to: '/volunteers', label: '帳號', icon: 'badge' }] : []),
  ...(isKnowledgeManager()
    ? [
        { to: '/knowledge', label: '知識庫', icon: 'menu_book' },
        { to: '/guide-settings', label: 'Sunny訓練設定', icon: 'tune' },
      ]
    : []),
  ...(isSuperAdmin() ? [{ to: '/users', label: '管理員', icon: 'manage_accounts' }] : []),
])

async function signOut() {
  await logout()
  router.push('/login')
}
</script>

<template>
  <router-view v-if="route.meta.public" />
  <q-layout v-else view="hHh Lpr lFf">
    <q-header class="bg-primary text-white">
      <q-toolbar>
        <q-btn flat round dense icon="menu" aria-label="選單" @click="drawer = !drawer" />
        <q-toolbar-title>菩提幣管理後台</q-toolbar-title>
        <span class="gt-xs q-mr-sm">{{ session.user?.display_name }}</span>
        <q-btn flat dense no-caps icon="key" label="修改密碼" class="q-mr-sm" @click="changingPassword = true" />
        <q-btn flat dense no-caps icon="logout" label="登出" @click="signOut" />
      </q-toolbar>
    </q-header>
    <q-drawer v-model="drawer" show-if-above bordered :width="220">
      <q-list padding>
        <q-item v-for="i in items" :key="i.to" :to="i.to" clickable active-class="text-secondary">
          <q-item-section avatar><q-icon :name="i.icon" /></q-item-section>
          <q-item-section>{{ i.label }}</q-item-section>
        </q-item>
      </q-list>
    </q-drawer>
    <q-page-container>
      <router-view />
    </q-page-container>
    <ChangePasswordDialog v-model="changingPassword" :save="changePassword" />
  </q-layout>
</template>

<style>
body {
  background: #f6f2e8;
  color: #3b2a20;
  font-family: 'Noto Sans TC', 'PingFang TC', 'Microsoft JhengHei', system-ui, sans-serif;
}
.admin-page {
  max-width: 1100px;
  margin: 0 auto;
  padding: 24px 16px 48px;
}
.admin-page h1 {
  font-size: 1.6rem;
  margin: 0 0 16px;
  line-height: 1.3;
}
</style>
