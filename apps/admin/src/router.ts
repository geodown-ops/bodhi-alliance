import { createRouter, createWebHistory } from 'vue-router'
import { isAdmin, isKnowledgeManager, restore, session } from './session'

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login', component: () => import('./pages/LoginPage.vue'), meta: { public: true } },
    { path: '/', component: { render: () => null } },
    { path: '/applications', component: () => import('./pages/ApplicationsPage.vue'), meta: { admin: true, title: '報名與登記' } },
    { path: '/venues', component: () => import('./pages/VenuesPage.vue'), meta: { admin: true, title: '場域管理' } },
    { path: '/merchants', component: () => import('./pages/MerchantsPage.vue'), meta: { admin: true, title: '共好企業管理' } },
    { path: '/groups', component: () => import('./pages/GroupsPage.vue'), meta: { admin: true, title: '覺行小組' } },
    { path: '/knowledge', component: () => import('./pages/KnowledgePage.vue'), meta: { title: '知識庫' } },
    { path: '/knowledge/:id', component: () => import('./pages/DocumentPage.vue'), meta: { title: '知識文件' } },
    { path: '/guide-settings', component: () => import('./pages/GuideSettingsPage.vue'), meta: { title: 'AI 組長設定' } },
    { path: '/users', component: () => import('./pages/UsersPage.vue'), meta: { admin: true, title: '帳號' } },
    { path: '/:pathMatch(.*)*', redirect: '/' },
  ],
})

let restored = false
router.beforeEach(async (to) => {
  if (!restored) {
    await restore()
    restored = true
  }
  if (to.meta.public) return true
  if (!session.user) return { path: '/login', query: to.path === '/' ? {} : { next: to.fullPath } }
  // 首頁依角色導向：放在 restore() 之後，才知道登入者是誰
  if (to.path === '/') return isAdmin() ? '/applications' : '/knowledge'
  if (to.meta.admin && !isAdmin()) return '/knowledge'
  if (!isKnowledgeManager()) return '/login'
  return true
})

router.afterEach((to) => {
  document.title = to.meta.title ? `${to.meta.title}｜菩提幣管理後台` : '菩提幣管理後台'
})
