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
  line_id: string
  home_center_id: string
  center_name: string
  wants_coach: boolean
  is_coach: boolean
  status: 'pending' | 'verified' | 'rejected'
  review_note: string
  frozen: boolean
  // 系統會員可以加入覺行小組、世界佛教教育協會，或兩者都加入
  in_groups: boolean
  in_association: boolean
  association_joined_at: string | null
  groups: { group_id: string; name: string; role: 'member' | 'leader' }[]
}
export type Center = { id: string; name: string; region: string }

// 覺行共修活動：my_role 與 claim_* 只出現在「我的活動」
export type PracticeEvent = {
  id: string
  title: string
  is_online: boolean
  location: string
  starts_at: string
  ends_at: string
  capacity: number
  description: string
  status: 'open' | 'cancelled'
  organizer_name: string
  joined: number
  venue_id: string | null
  venue_name: string
  my_role?: 'organizer' | 'helper' | 'participant'
  claim_status?: '' | 'submitted' | 'approved' | 'rejected'
  claim_note?: string
}
// 活動場域：upcoming 是還沒結束的活動數
export type Venue = { id: string; name: string; center_name: string; region: string; address: string; description: string; upcoming: number }
// 協會會員看得到的會刊、行事曆與通知
export type AssociationFeed = {
  issues: { id: string; title: string; issued_on: string; summary: string; url: string }[]
  events: { id: string; title: string; starts_at: string; ends_at: string | null; location: string; description: string; kind: string; tag_name: string }[]
  notices: { id: string; title: string; body: string; created_at: string }[]
}
export type ChainGrant = { kind: string; amount: string; status: 'pending' | 'sent' | 'confirmed' | 'failed'; tx_hash: string; tx_url: string; confirmed_at: string | null }
export type ChainAccount = {
  enabled: boolean
  network?: string
  address?: string
  address_url?: string
  token_url?: string
  contract_url?: string
  balance?: string
  grants?: ChainGrant[]
}
export type Wallet = { balance: number; entries: { amount: number; kind: string; memo: string; created_at: string }[] }

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
  memberships: (body: { in_groups: boolean; in_association: boolean }) => call<Volunteer>('PUT', '/api/me/memberships', body),
  association: () => call<AssociationFeed>('GET', '/api/me/association'),
  join: (groupId: string) => call('POST', `/api/me/groups/${groupId}`),
  leave: (groupId: string) => call('DELETE', `/api/me/groups/${groupId}`),
  wallet: () => call<Wallet>('GET', '/api/me/wallet'),
  chain: () => call<ChainAccount>('GET', '/api/me/chain'),
  events: () => call<PracticeEvent[]>('GET', '/api/me/events'),
  createEvent: (body: Record<string, unknown>) => call<{ id: string }>('POST', '/api/me/events', body),
  joinEvent: (id: string, role: 'participant' | 'helper') => call('POST', `/api/me/events/${id}/join`, { role }),
  leaveEvent: (id: string) => call('DELETE', `/api/me/events/${id}/join`),
  cancelEvent: (id: string) => call('POST', `/api/me/events/${id}/cancel`),
  claim: (id: string, body: { attendance: number; report: string }) => call('POST', `/api/me/events/${id}/claim`, body),
}

export const publicEvents = () => call<PracticeEvent[]>('GET', '/api/events')
export const publicVenues = () => call<Venue[]>('GET', '/api/venues')

const dateFmt = new Intl.DateTimeFormat('zh-TW', { month: 'numeric', day: 'numeric', weekday: 'short', hour: '2-digit', minute: '2-digit', hour12: false })
const timeFmt = new Intl.DateTimeFormat('zh-TW', { hour: '2-digit', minute: '2-digit', hour12: false })

// 活動時間：同一天只寫一次日期
export function eventTime(e: Pick<PracticeEvent, 'starts_at' | 'ends_at'>) {
  const s = new Date(e.starts_at)
  const t = new Date(e.ends_at)
  const sameDay = s.toDateString() === t.toDateString()
  return `${dateFmt.format(s)} – ${sameDay ? timeFmt.format(t) : dateFmt.format(t)}`
}
