import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { account, restore } from './account'

// 主選單八項，順序即選單順序
export const menu = [
  { path: '/coin', label: '菩提幣介紹', icon: 'eco' },
  { path: '/mindfulness', label: '正念減壓', icon: 'air' },
  { path: '/groups', label: '覺行小組', icon: 'self_improvement' },
  { path: '/partners', label: '共好企業', icon: 'storefront' },
  { path: '/wallet', label: '我的錢包', icon: 'account_balance_wallet' },
  { path: '/maitreya', label: '彌勒心流', icon: 'spa' },
  { path: '/association', label: '世界佛教教育協會', icon: 'temple_buddhist' },
  { path: '/guide', label: '線上問答', icon: 'forum' },
] as const

const routes: RouteRecordRaw[] = [
  { path: '/', component: () => import('./pages/HomePage.vue'), meta: { title: '' } },
  { path: '/coin', component: () => import('./pages/CoinPage.vue'), meta: { title: '菩提幣介紹' } },
  { path: '/mindfulness', component: () => import('./pages/MindfulnessPage.vue'), meta: { title: '正念減壓' } },
  { path: '/mindfulness/lesson/:n(\\d)', component: () => import('./pages/MindfulnessLessonPage.vue'), meta: { title: '正念減壓' } },
  { path: '/mindfulness/week/:n(\\d)', redirect: (to) => `/mindfulness/lesson/${to.params.n}` },
  { path: '/mindfulness/sitting', component: () => import('./pages/MindfulnessSittingPage.vue'), meta: { title: '上座與下座' } },
  { path: '/groups', component: () => import('./pages/GroupsPage.vue'), meta: { title: '覺行小組介紹及參加' } },
  { path: '/partners', component: () => import('./pages/PartnersPage.vue'), meta: { title: '共好企業登記及管理' } },
  { path: '/wallet', component: () => import('./pages/WalletPage.vue'), meta: { title: '我的錢包' } },
  { path: '/maitreya', component: () => import('./pages/MaitreyaPage.vue'), meta: { title: '彌勒心流' } },
  { path: '/committee', component: () => import('./pages/CommitteePage.vue'), meta: { title: '主辦審核小組' } },
  { path: '/association', component: () => import('./pages/AssociationPage.vue'), meta: { title: '世界佛教教育協會介紹' } },
  { path: '/guide', component: () => import('./pages/GuidePage.vue'), meta: { title: '線上問答' } },
  { path: '/join', component: () => import('./pages/JoinPage.vue'), meta: { title: '報名參加覺行小組' } },
  { path: '/login', component: () => import('./pages/LoginPage.vue'), meta: { title: '登入' } },
  { path: '/me', component: () => import('./pages/MePage.vue'), meta: { title: '我的個人頁', signedIn: true } },
  { path: '/:pathMatch(.*)*', redirect: '/' },
]

export const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior: (_to, _from, saved) => saved ?? { top: 0 },
})

router.beforeEach(async (to) => {
  if (!account.ready) await restore()
  if (to.meta.signedIn && !account.user) return { path: '/login', query: { next: to.fullPath } }
  if ((to.path === '/login' || to.path === '/join') && account.user) return '/me'
  return true
})

router.afterEach((to) => {
  const t = to.meta.title as string
  document.title = t ? `${t}｜Sunny life` : 'Sunny life｜世界佛教教育協會'
})
