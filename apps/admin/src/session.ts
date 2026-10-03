import { reactive } from 'vue'

export type Role = { role: string; scope: string }
export type User = { id: string; email: string; display_name: string; roles: Role[] }

const apiBase = import.meta.env.VITE_API_BASE ?? ''
const guideBase = import.meta.env.VITE_GUIDE_BASE ?? ''
const tokenKey = 'bodhi.admin.token'

function readToken() {
  try {
    return localStorage.getItem(tokenKey) ?? ''
  } catch {
    return ''
  }
}

export const session = reactive({ token: readToken(), user: null as User | null })

export function isAdmin(u = session.user) {
  return !!u?.roles.some((r) => r.role === 'alliance_admin' && r.scope === 'alliance')
}
export function isKnowledgeManager(u = session.user) {
  return isAdmin(u) || !!u?.roles.some((r) => r.role === 'knowledge_manager' && r.scope === 'guide')
}

function setToken(t: string) {
  session.token = t
  try {
    if (t) localStorage.setItem(tokenKey, t)
    else localStorage.removeItem(tokenKey)
  } catch {
    /* ignore */
  }
}

export class ApiError extends Error {
  constructor(message: string, public status: number) {
    super(message)
  }
}

async function call<T>(url: string, init: RequestInit = {}): Promise<T> {
  const headers: Record<string, string> = { ...(init.headers as Record<string, string>) }
  if (!(init.body instanceof FormData)) headers['Content-Type'] = 'application/json'
  if (session.token) headers.Authorization = `Bearer ${session.token}`
  let res: Response
  try {
    res = await fetch(url, { ...init, headers })
  } catch {
    throw new ApiError('連線失敗', 0)
  }
  const text = await res.text()
  const data = text ? JSON.parse(text) : undefined
  if (res.status === 401 && session.token) {
    setToken('')
    session.user = null
  }
  if (!res.ok) throw new ApiError(data?.error ?? `錯誤 ${res.status}`, res.status)
  return data as T
}

const json = (method: string, body?: unknown): RequestInit => ({ method, body: body === undefined ? undefined : JSON.stringify(body) })

export const api = {
  get: <T>(path: string) => call<T>(apiBase + path),
  send: <T = void>(method: string, path: string, body?: unknown) => call<T>(apiBase + path, json(method, body)),
  guideGet: <T>(path: string) => call<T>(guideBase + path),
  guideSend: <T = void>(method: string, path: string, body?: unknown) => call<T>(guideBase + path, json(method, body)),
  guideUpload: <T>(path: string, form: FormData) => call<T>(guideBase + path, { method: 'POST', body: form }),
}

export async function login(email: string, password: string) {
  const res = await call<{ token: string; user: User }>(`${apiBase}/api/auth/login`, json('POST', { email, password }))
  setToken(res.token)
  session.user = res.user
}

export async function restore() {
  if (!session.token) return
  try {
    session.user = await api.get<User>('/api/auth/me')
  } catch {
    session.user = null
  }
}

export async function logout() {
  try {
    await api.send('POST', '/api/auth/logout')
  } finally {
    setToken('')
    session.user = null
  }
}
