// Sunny 的聲音：用瀏覽器內建的語音合成（Web Speech API）把回答一句一句唸出來。
// 回答是串流進來的，湊滿一句就先唸，不必等整段寫完。音色取決於使用者裝置上的中文語音。

const SENTENCE_END = /[。！？!?；;\n]/
const SOFT_BREAK = /[，、,：:]/
const MAX_CHUNK = 60 // 太長的句子在逗號處先斷開，部分瀏覽器唸超過十幾秒會中斷

/** 把串流文字切成可以唸的句子；rest 是還沒寫完的半句 */
export function splitSentences(buf: string): { sentences: string[]; rest: string } {
  const sentences: string[] = []
  let start = 0
  for (let i = 0; i < buf.length; i++) {
    const ch = buf[i]
    if (SENTENCE_END.test(ch) || (i - start + 1 >= MAX_CHUNK && SOFT_BREAK.test(ch))) {
      sentences.push(buf.slice(start, i + 1))
      start = i + 1
    }
  }
  return { sentences, rest: buf.slice(start) }
}

/** 唸之前拿掉網址、Markdown 符號與表情符號 */
export function cleanForSpeech(text: string): string {
  return text
    .replace(/https?:\/\/\S+/g, '')
    .replace(/[*#_`>|~[\]]/g, '')
    .replace(/\p{Extended_Pictographic}/gu, '')
    .replace(/\s+/g, ' ')
    .trim()
}

// 優先挑臺灣華語，其次其他華語語音（不唸粵語）；同一種語言裡先挑自然的女聲：
// Edge 的 HsiaoChen（Natural）、Chrome 的 Google 國語（臺灣）、Safari 的美佳
const FEMALE = /HsiaoChen|HsiaoYu|Hanhan|Yating|Mei-?Jia|Tingting|Xiaoxiao|Xiaoyi|female|女/i
export function pickVoice<V extends { name: string; lang: string }>(voices: V[]): V | undefined {
  const tw = voices.filter((v) => /^(zh|cmn)[-_](TW|Hant)/i.test(v.lang))
  const pool = tw.length ? tw : voices.filter((v) => /^(zh|cmn)/i.test(v.lang) && !/HK|yue/i.test(v.lang))
  return (
    pool.find((v) => FEMALE.test(v.name) && /Natural|Online/i.test(v.name)) ??
    pool.find((v) => /Google/i.test(v.name)) ??
    pool.find((v) => FEMALE.test(v.name)) ??
    pool.find((v) => /Natural|Online/i.test(v.name)) ??
    pool[0]
  )
}

export type VoiceHooks = {
  /** 開始唸一句 */
  onStart: () => void
  /** 唸完一句 */
  onEnd: () => void
  /** end() 之後全部唸完 */
  onDone: () => void
  /** 這一句唸不出來（裝置沒有中文語音、瀏覽器拒絕等），改用文字對嘴 */
  onFail: (text: string) => void
}

export function createVoice(hooks: VoiceHooks) {
  const synth = typeof window !== 'undefined' && 'speechSynthesis' in window ? window.speechSynthesis : null
  let buf = ''
  let ending = false
  let generation = 0
  let unlocked = false
  // 保留參照：Chrome 會把沒有人引用的 utterance 回收掉，之後就收不到 onend
  const pending = new Set<SpeechSynthesisUtterance>()
  synth?.getVoices() // Chrome 第一次呼叫時才開始載入語音清單

  const settle = () => {
    if (ending && !pending.size) {
      ending = false
      hooks.onDone()
    }
  }

  const say = (raw: string) => {
    const text = cleanForSpeech(raw)
    if (!text) return
    const voices = synth?.getVoices() ?? []
    const v = pickVoice(voices)
    // 語音清單已載入卻沒有華語語音，用英文語音唸中文只會是怪聲，直接改用文字對嘴
    if (!synth || (voices.length && !v)) return hooks.onFail(text)
    const u = new SpeechSynthesisUtterance(text)
    if (v) u.voice = v
    u.lang = v?.lang ?? 'zh-TW'
    u.pitch = 1.05
    const gen = generation
    let started = false
    u.onstart = () => {
      if (gen !== generation) return
      started = true
      hooks.onStart()
    }
    const done = (failed: boolean) => {
      pending.delete(u)
      if (gen !== generation) return
      if (started) hooks.onEnd()
      else if (failed) hooks.onFail(text)
      settle()
    }
    u.onend = () => done(false)
    u.onerror = (e) => done(e.error !== 'interrupted' && e.error !== 'canceled')
    pending.add(u)
    synth.speak(u)
  }

  return {
    supported: !!synth,
    /** 在點擊當下呼叫一次：iOS Safari 只允許由使用者操作開始發聲 */
    unlock() {
      if (!synth || unlocked) return
      unlocked = true
      const u = new SpeechSynthesisUtterance(' ')
      u.volume = 0
      synth.speak(u)
    },
    /** 收到一段串流文字，湊滿的句子先唸 */
    feed(text: string) {
      const { sentences, rest } = splitSentences(buf + text)
      buf = rest
      sentences.forEach(say)
    },
    /** 回答寫完了：唸完剩下的半句，全部唸完後呼叫 onDone */
    end() {
      say(buf)
      buf = ''
      ending = true
      settle()
    },
    /** 立刻停止，不再呼叫任何 hook */
    stop() {
      generation++
      buf = ''
      ending = false
      pending.clear()
      synth?.cancel()
    },
  }
}

export type Voice = ReturnType<typeof createVoice>
