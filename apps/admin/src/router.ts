import { createRouter, createWebHistory } from 'vue-router'
import { homePath, isAdmin, isCenterStaff, isKnowledgeManager, restore, session } from './session'

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login', component: () => import('./pages/LoginPage.vue'), meta: { public: true } },
    { path: '/', component: { render: () => null } },
    { path: '/applications', component: () => import('./pages/ApplicationsPage.vue'), meta: { admin: true, title: '報名與登記' } },
    { path: '/venues', component: () => import('./pages/VenuesPage.vue'), meta: { admin: true, title: '場域管理' } },
    { path: '/merchants', component: () => import('./pages/MerchantsPage.vue'), meta: { admin: true, title: '共好企業管理' } },
    { path: '/groups', component: () => import('./pages/GroupsPage.vue'), meta: { admin: true, title: '覺行小組' } },
    { path: '/claims', component: () => import('./pages/ClaimsPage.vue'), meta: { admin: true, title: '菩提幣審核' } },
    { path: '/association', component: () => import('./pages/AssociationPage.vue'), meta: { admin: true, title: '世界佛教教育協會' } },
    { path: '/volunteers', component: () => import('./pages/VolunteersPage.vue'), meta: { center: true, title: '會員名冊' } },
    { path: '/knowledge', component: () => import('./pages/KnowledgePage.vue'), meta: { knowledge: true, title: '知識庫' } },
    { path: '/knowledge/:id', component: () => import('./pages/DocumentPage.vue'), meta: { knowledge: true, title: '知識文件' } },
    { path: '/guide-settings', component: () => import('./pages/GuideSettingsPage.vue'), meta: { knowledge: true, title: 'AI 組長設定' } },
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
  // 依角色導向：放在 restore() 之後，才知道登入者是誰
  const home = homePath()
  if (!home) return { path: '/login', query: { denied: '1' } }
  if (to.path === '/') return home
  if (to.meta.admin && !isAdmin()) return home
  if (to.meta.center && !isCenterStaff()) return home
  if (to.meta.knowledge && !isKnowledgeManager()) return home
  return true
})

router.afterEach((to) => {
  document.title = to.meta.title ? `${to.meta.title}｜菩提幣管理後台` : '菩提幣管理後台'
})
