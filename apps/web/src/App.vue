<script setup lang="ts">
import { ref } from 'vue'
import { menu } from './router'
import { account } from './account'
import { wordmark } from './wordmark'

const drawer = ref(false)
</script>

<template>
  <q-layout view="hHh lpR fff">
    <q-header class="bg-primary text-white" bordered>
      <q-toolbar>
        <q-btn flat round dense icon="menu" class="lt-md" aria-label="選單" @click="drawer = !drawer" />
        <router-link to="/" class="brand">
          <img src="/favicon.svg" alt="" width="32" height="32" />
          <svg class="wordmark" :viewBox="wordmark.viewBox" role="img" aria-label="Sunny life"><path fill="currentColor" :d="wordmark.d" /></svg>
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
          :aria-label="account.user ? '我的個人頁' : '登入'"
        >
          <q-tooltip>{{ account.user ? '我的個人頁' : '登入' }}</q-tooltip>
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
          <q-item-section>{{ account.user ? '我的個人頁' : '登入／報名' }}</q-item-section>
        </q-item>
      </q-list>
    </q-drawer>

    <q-page-container>
      <router-view />
    </q-page-container>

    <q-footer class="footer">
      <div>一即一切，一切即一</div>
      <router-link to="/privacy" class="footer-link">隱私權保護聲明</router-link>
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
  margin-left: 4px;
}
.wordmark {
  height: 34px;
  width: auto;
  display: block;
}
.footer {
  background: var(--ground-sunk);
  color: var(--ink-soft);
  text-align: center;
  padding: 20px 16px;
  line-height: 1.9;
  font-size: 0.9rem;
}
.footer-link {
  color: inherit;
  font-size: 0.8rem;
  opacity: 0.8;
}
.footer-link:hover {
  opacity: 1;
}
</style>
