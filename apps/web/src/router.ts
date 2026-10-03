import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { account, restore } from './account'

// 主選單七項，順序即選單順序
export const menu = [
  { path: '/coin', label: '菩提幣介紹', icon: 'eco' },
  { path: '/groups', label: '覺行小組', icon: 'self_improvement' },
  { path: '/partners', label: '共好企業', icon: 'storefront' },
  { path: '/wallet', label: '菩提幣錢包', icon: 'account_balance_wallet' },
  { path: '/committee', label: '主辦審核小組', icon: 'gavel' },
  { path: '/association', label: '世界佛教教育協會', icon: 'temple_buddhist' },
  { path: '/guide', label: 'AI 組長', icon: 'forum' },
] as const

const routes: RouteRecordRaw[] = [
  { path: '/', component: () => import('./pages/HomePage.vue'), meta: { title: '' } },
  { path: '/coin', component: () => import('./pages/CoinPage.vue'), meta: { title: '菩提幣介紹' } },
  { path: '/groups', component: () => import('./pages/GroupsPage.vue'), meta: { title: '覺行小組介紹及參加' } },
  { path: '/partners', component: () => import('./pages/PartnersPage.vue'), meta: { title: '共好企業登記及管理' } },
  { path: '/wallet', component: () => import('./pages/WalletPage.vue'), meta: { title: '菩提幣錢包' } },
  { path: '/committee', component: () => import('./pages/CommitteePage.vue'), meta: { title: '主辦審核小組' } },
  { path: '/association', component: () => import('./pages/AssociationPage.vue'), meta: { title: '世界佛教教育協會介紹' } },
  { path: '/guide', component: () => import('./pages/GuidePage.vue'), meta: { title: '線上覺行小組 AI 組長' } },
  { path: '/join', component: () => import('./pages/JoinPage.vue'), meta: { title: '志工註冊' } },
  { path: '/login', component: () => import('./pages/LoginPage.vue'), meta: { title: '志工登入' } },
  { path: '/me', component: () => import('./pages/MePage.vue'), meta: { title: '我的志工資料', signedIn: true } },
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
  document.title = t ? `${t}｜菩提幣聯盟` : '菩提幣聯盟｜世界佛教教育協會'
})
