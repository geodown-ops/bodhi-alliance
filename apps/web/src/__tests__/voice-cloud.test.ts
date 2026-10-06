import { afterEach, describe, expect, it, vi } from 'vitest'
import { createVoice } from '../voice'

// 假的 Web Audio：每段聲音 start 後下一個 tick 就播完
class FakeSource {
  buffer: unknown = null
  onended: (() => void) | null = null
  connect() {}
  start() {
    setTimeout(() => this.onended?.(), 0)
  }
  stop() {
    this.onended?.()
  }
}
class FakeAudioContext {
  state = 'running'
  destination = {}
  resume() {
    return Promise.resolve()
  }
  createBuffer() {
    return {}
  }
  createBufferSource() {
    return new FakeSource()
  }
  decodeAudioData(b: ArrayBuffer) {
    return Promise.resolve({ text: new TextDecoder().decode(b) })
  }
}

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('cloud voice', () => {
  it('plays sentences in order even when later audio arrives first, then finishes', async () => {
    vi.stubGlobal('AudioContext', FakeAudioContext)
    const log: string[] = []
    const done = new Promise<void>((resolve) => {
      const voice = createVoice({
        onStart: () => log.push('start'),
        onEnd: () => log.push('end'),
        onDone: resolve,
        onFail: (t) => log.push('fail:' + t),
      })
      const delays: Record<string, number> = { '你好。': 20, '一起靜坐吧？': 0 }
      voice.useCloud(async (text) => {
        log.push('fetch:' + text)
        await new Promise((r) => setTimeout(r, delays[text] ?? 0))
        return new TextEncoder().encode(text).buffer as ArrayBuffer
      })
      voice.unlock()
      voice.feed('你好。一起靜坐吧？')
      voice.end()
    })
    await done
    expect(log).toEqual(['fetch:你好。', 'fetch:一起靜坐吧？', 'start', 'end', 'start', 'end'])
  })

  it('falls back when the cloud voice fails', async () => {
    vi.stubGlobal('AudioContext', FakeAudioContext)
    const failed: string[] = []
    await new Promise<void>((resolve) => {
      const voice = createVoice({ onStart() {}, onEnd() {}, onDone: resolve, onFail: (t) => failed.push(t) })
      voice.useCloud(() => Promise.reject(new Error('503')))
      voice.unlock()
      voice.feed('你好。')
      voice.end()
    })
    // jsdom 沒有瀏覽器語音，最後改用文字對嘴
    expect(failed).toEqual(['你好。'])
  })
})
