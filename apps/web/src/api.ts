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
export type ChatRequest = { messages: ChatMessage[]; mode: 'chat' | 'practice'; script_id?: string }
export type ChatReply = { reply: string; refused?: boolean }

// chatStream asks the guide and calls onText with each piece of the answer as it is
// written. The resolved reply is the whole answer; on a refusal it replaces what was
// streamed.
async function chatStream(body: ChatRequest, onText: (text: string) => void): Promise<ChatReply> {
  let res: Response
  try {
    res = await fetch(`${guideBase}/guide/chat`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', Accept: 'text/event-stream' },
      body: JSON.stringify(body),
    })
  } catch {
    throw new ApiError('連線失敗，請檢查網路後再試一次。')
  }
  if (!res.ok || !res.body) {
    const data = await res.json().catch(() => undefined)
    throw new ApiError(data?.error ?? '組長一時沒有回應，請再試一次。')
  }
  const reader = res.body.pipeThrough(new TextDecoderStream()).getReader()
  let buf = ''
  for (;;) {
    const { value, done } = await reader.read()
    if (value) buf += value
    let end
    while ((end = buf.indexOf('\n\n')) >= 0) {
      const raw = buf.slice(0, end)
      buf = buf.slice(end + 2)
      let event = 'message'
      let data = ''
      for (const line of raw.split('\n')) {
        if (line.startsWith('event:')) event = line.slice(6).trim()
        else if (line.startsWith('data:')) data += line.slice(5).trimStart()
      }
      const payload = data ? JSON.parse(data) : {}
      if (event === 'delta') onText(payload.text)
      else if (event === 'done') return payload as ChatReply
      else if (event === 'error') throw new ApiError(payload.error ?? '組長一時沒有回應，請再試一次。')
    }
    if (done) throw new ApiError('組長的回答中斷了，請再試一次。')
  }
}

export type ChainInfo = {
  enabled: boolean
  chain_id?: number
  network?: string
  explorer?: string
  contract?: string
  contract_url?: string
}

export const api = {
  groups: () => request<Group[]>(`${apiBase}/api/groups`),
  chainInfo: () => request<ChainInfo>(`${apiBase}/api/chain`),
  joinGroup: (body: Record<string, unknown>) =>
    request(`${apiBase}/api/group-applications`, { method: 'POST', body: JSON.stringify(body) }),
  registerPartner: (body: Record<string, unknown>) =>
    request(`${apiBase}/api/merchant-applications`, { method: 'POST', body: JSON.stringify(body) }),

  guideInfo: () => request<{ name: string; available: boolean; tts?: boolean }>(`${guideBase}/guide/info`),
  /** Sunny 的雲端語音，一次一句 */
  async tts(text: string): Promise<ArrayBuffer> {
    const res = await fetch(`${guideBase}/guide/tts`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ text }),
    })
    if (!res.ok) throw new ApiError('語音暫時無法使用')
    return res.arrayBuffer()
  },
  scripts: () => request<{ id: string; title: string }[]>(`${guideBase}/guide/scripts`),
  chat: (body: ChatRequest) =>
    request<ChatReply>(`${guideBase}/guide/chat`, { method: 'POST', body: JSON.stringify(body) }),
  chatStream,
}
