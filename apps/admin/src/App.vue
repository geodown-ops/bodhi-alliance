<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { isAdmin, logout, session } from './session'

const route = useRoute()
const router = useRouter()
const drawer = ref(false)

const items = computed(() => [
  ...(isAdmin()
    ? [
        { to: '/applications', label: '報名與登記', icon: 'inbox' },
        { to: '/groups', label: '覺行小組', icon: 'groups' },
        { to: '/venues', label: '場域管理', icon: 'place' },
        { to: '/merchants', label: '共好企業管理', icon: 'storefront' },
      ]
    : []),
  { to: '/knowledge', label: '知識庫', icon: 'menu_book' },
  { to: '/guide-settings', label: 'AI 組長設定', icon: 'tune' },
  ...(isAdmin() ? [{ to: '/users', label: '帳號', icon: 'manage_accounts' }] : []),
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
