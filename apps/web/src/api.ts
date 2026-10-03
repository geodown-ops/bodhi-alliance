const apiBase = import.meta.env.VITE_API_BASE ?? ''
const guideBase = import.meta.env.VITE_GUIDE_BASE ?? ''

export class ApiError extends Error {}

async function request<T>(url: string, init?: RequestInit): Promise<T> {
  let res: Response
  try {
    res = await fetch(url, { ...init, headers: { 'Content-Type': 'application/json', ...init?.headers } })
  } catch {
    throw new ApiError('連線失敗，請檢查網路後再試一次。')
  }
  const text = await res.text()
  const data = text ? JSON.parse(text) : undefined
  if (!res.ok) throw new ApiError(data?.error ?? '發生錯誤，請稍後再試。')
  return data as T
}

export type Group = {
  id: string
  name: string
  region: string
  center_name: string
  schedule: string
  description: string
  is_online: boolean
}

export type ChatMessage = { role: 'user' | 'assistant'; content: string }

export const api = {
  groups: () => request<Group[]>(`${apiBase}/api/groups`),
  joinGroup: (body: Record<string, unknown>) =>
    request(`${apiBase}/api/group-applications`, { method: 'POST', body: JSON.stringify(body) }),
  registerPartner: (body: Record<string, unknown>) =>
    request(`${apiBase}/api/merchant-applications`, { method: 'POST', body: JSON.stringify(body) }),

  guideInfo: () => request<{ name: string; available: boolean }>(`${guideBase}/guide/info`),
  scripts: () => request<{ id: string; title: string }[]>(`${guideBase}/guide/scripts`),
  chat: (body: { messages: ChatMessage[]; mode: 'chat' | 'practice'; script_id?: string }) =>
    request<{ reply: string }>(`${guideBase}/guide/chat`, { method: 'POST', body: JSON.stringify(body) }),
}
