import { afterEach, describe, expect, it, vi } from 'vitest'
import { api, ApiError } from '../api'

afterEach(() => vi.unstubAllGlobals())

describe('api', () => {
  it('surfaces the server error message', async () => {
    vi.stubGlobal('fetch', vi.fn(async () => new Response(JSON.stringify({ error: '請填寫姓名與正確的電子郵件' }), { status: 400 })))
    await expect(api.joinGroup({})).rejects.toThrow('請填寫姓名與正確的電子郵件')
  })

  it('reports a network failure in Chinese', async () => {
    vi.stubGlobal('fetch', vi.fn(async () => { throw new TypeError('Failed to fetch') }))
    await expect(api.groups()).rejects.toBeInstanceOf(ApiError)
  })

  it('posts the chat history', async () => {
    const fetch = vi.fn(async () => new Response(JSON.stringify({ reply: '我是線上覺行小組組長Sunny' })))
    vi.stubGlobal('fetch', fetch)
    const res = await api.chat({ mode: 'chat', messages: [{ role: 'user', content: '你是誰' }] })
    expect(res.reply).toBe('我是線上覺行小組組長Sunny')
    const [url, init] = fetch.mock.calls[0] as unknown as [string, RequestInit]
    expect(url).toBe('/guide/chat')
    expect(JSON.parse(init.body as string).messages[0].content).toBe('你是誰')
  })

  it('streams the answer piece by piece', async () => {
    // Split mid-event to check buffering across reads.
    const parts = ['event:delta\ndata:{"text":"我是"}\n\nevent:del', 'ta\ndata:{"text":"Sunny"}\n\nevent:done\ndata:{"reply":"我是Sunny"}\n\n']
    const body = new ReadableStream({
      start(c) {
        for (const p of parts) c.enqueue(new TextEncoder().encode(p))
        c.close()
      },
    })
    const fetch = vi.fn(async () => new Response(body, { headers: { 'Content-Type': 'text/event-stream' } }))
    vi.stubGlobal('fetch', fetch)
    const seen: string[] = []
    const res = await api.chatStream({ mode: 'chat', messages: [{ role: 'user', content: '你是誰' }] }, (t) => seen.push(t))
    expect(seen).toEqual(['我是', 'Sunny'])
    expect(res.reply).toBe('我是Sunny')
    const [, init] = fetch.mock.calls[0] as unknown as [string, RequestInit]
    expect((init.headers as Record<string, string>).Accept).toBe('text/event-stream')
  })

  it('reports a stream error', async () => {
    vi.stubGlobal('fetch', vi.fn(async () => new Response('event:error\ndata:{"error":"線上組長一時沒有回應，請再試一次。"}\n\n')))
    await expect(api.chatStream({ mode: 'chat', messages: [] }, () => {})).rejects.toThrow('線上組長一時沒有回應')
  })
})
