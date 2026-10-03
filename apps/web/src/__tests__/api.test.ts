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
})
