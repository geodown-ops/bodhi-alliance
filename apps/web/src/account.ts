import { reactive } from 'vue'
import { ApiError } from './api'

// 志工帳號：登入狀態存在 localStorage，與管理後台分開。

export type User = { id: string; email: string; display_name: string; roles: { role: string; scope: string }[] }
export type Volunteer = {
  id: string
  email: string
  display_name: string
  legal_name: string
  phone: string
  home_center_id: string
  center_name: string
  wants_coach: boolean
  is_coach: boolean
  status: 'pending' | 'verified' | 'rejected'
  review_note: string
  frozen: boolean
  groups: { group_id: string; name: string; role: 'member' | 'leader' }[]
}
export type Center = { id: string; name: string; region: string }

const apiBase = import.meta.env.VITE_API_BASE ?? ''
const tokenKey = 'bodhi.web.token'

function readToken() {
  try {
    return localStorage.getItem(tokenKey) ?? ''
  } catch {
    return ''
  }
}

export const account = reactive({ token: readToken(), user: null as User | null, ready: false })

function setToken(t: string) {
  account.token = t
  try {
    if (t) localStorage.setItem(tokenKey, t)
    else localStorage.removeItem(tokenKey)
  } catch {
    /* 無法使用 localStorage 時只在這個分頁保持登入 */
  }
}

async function call<T>(method: string, path: string, body?: unknown): Promise<T> {
  const headers: Record<string, string> = { 'Content-Type': 'application/json' }
  if (account.token) headers.Authorization = `Bearer ${account.token}`
  let res: Response
  try {
    res = await fetch(apiBase + path, { method, headers, body: body === undefined ? undefined : JSON.stringify(body) })
  } catch {
    throw new ApiError('連線失敗，請檢查網路後再試一次。')
  }
  const text = await res.text()
  const data = text ? JSON.parse(text) : undefined
  if (res.status === 401 && account.token) {
    setToken('')
    account.user = null
  }
  if (!res.ok) throw new ApiError(data?.error ?? '發生錯誤，請稍後再試。')
  return data as T
}

export async function restore() {
  if (account.token && !account.user) {
    try {
      account.user = await call<User>('GET', '/api/auth/me')
    } catch {
      account.user = null
    }
  }
  account.ready = true
}

export async function login(email: string, password: string) {
  const res = await call<{ token: string; user: User }>('POST', '/api/auth/login', { email, password })
  setToken(res.token)
  account.user = res.user
}

export async function register(body: Record<string, unknown>) {
  const res = await call<{ token?: string; user?: User }>('POST', '/api/volunteers', body)
  if (res.token) {
    setToken(res.token)
    account.user = res.user ?? null
  }
}

export async function logout() {
  try {
    await call('POST', '/api/auth/logout')
  } finally {
    setToken('')
    account.user = null
  }
}

export const me = {
  centers: () => call<Center[]>('GET', '/api/centers'),
  profile: () => call<Volunteer>('GET', '/api/me/volunteer'),
  update: (body: Record<string, unknown>) => call<Volunteer>('PUT', '/api/me/volunteer', body),
  join: (groupId: string) => call('POST', `/api/me/groups/${groupId}`),
  leave: (groupId: string) => call('DELETE', `/api/me/groups/${groupId}`),
}
