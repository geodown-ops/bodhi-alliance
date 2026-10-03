<script setup lang="ts">
import { ref } from 'vue'
import { menu } from './router'
import { account } from './account'

const drawer = ref(false)
</script>

<template>
  <q-layout view="hHh lpR fff">
    <q-header class="bg-primary text-white" bordered>
      <q-toolbar>
        <q-btn flat round dense icon="menu" class="lt-md" aria-label="選單" @click="drawer = !drawer" />
        <router-link to="/" class="brand">
          <img src="/favicon.svg" alt="" width="32" height="32" />
          <span>Sunny life</span>
        </router-link>
        <q-space />
        <nav class="gt-sm row no-wrap">
          <q-btn v-for="item in menu" :key="item.path" :to="item.path" flat no-caps :label="item.label" />
        </nav>
        <q-btn
          :to="account.user ? '/me' : '/login'"
          flat
          round
          dense
          icon="account_circle"
          class="q-ml-xs"
          :aria-label="account.user ? '我的志工資料' : '志工登入'"
        >
          <q-tooltip>{{ account.user ? '我的志工資料' : '志工登入' }}</q-tooltip>
        </q-btn>
      </q-toolbar>
    </q-header>

    <q-drawer v-model="drawer" side="left" overlay bordered behavior="mobile">
      <q-list>
        <q-item clickable to="/" exact @click="drawer = false">
          <q-item-section avatar><q-icon name="home" /></q-item-section>
          <q-item-section>首頁</q-item-section>
        </q-item>
        <q-item v-for="item in menu" :key="item.path" clickable :to="item.path" @click="drawer = false">
          <q-item-section avatar><q-icon :name="item.icon" /></q-item-section>
          <q-item-section>{{ item.label }}</q-item-section>
        </q-item>
        <q-separator />
        <q-item clickable :to="account.user ? '/me' : '/login'" @click="drawer = false">
          <q-item-section avatar><q-icon name="account_circle" /></q-item-section>
          <q-item-section>{{ account.user ? '我的志工資料' : '志工登入／註冊' }}</q-item-section>
        </q-item>
      </q-list>
    </q-drawer>

    <q-page-container>
      <router-view />
    </q-page-container>

    <q-footer class="footer">
      <div>Sunny life · 世界佛教教育協會</div>
      <div>幣不販售 · 不提領 · 不可兌現 · 非投資標的</div>
      <a href="https://www.sunnylife.world" target="_blank" rel="noopener">www.sunnylife.world</a>
    </q-footer>
  </q-layout>
</template>

<style scoped>
.brand {
  display: flex;
  align-items: center;
  gap: 10px;
  color: inherit;
  text-decoration: none;
  font-family: var(--wenkai);
  font-size: 1.3rem;
  margin-left: 4px;
}
.footer {
  background: var(--ground-sunk);
  color: var(--ink-soft);
  text-align: center;
  padding: 20px 16px;
  line-height: 1.9;
  font-size: 0.9rem;
}
.footer a {
  color: var(--leaf);
}
</style>
