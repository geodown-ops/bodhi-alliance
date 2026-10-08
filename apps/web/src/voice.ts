// Sunny 的聲音：把回答一句一句唸出來。回答是串流進來的，湊滿一句就先唸，不必等整段寫完。
// 伺服器有設定雲端語音（臺灣華語神經語音）時用雲端語音，比較自然；沒有設定或某一句取不到時，
// 改用瀏覽器內建的語音合成（Web Speech API），音色取決於使用者裝置上的中文語音。

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

// 語音引擎常把多音字唸錯，朗讀前換成同音字（畫面上的字不變）。
// 「覺行」唸 jué xíng。「覺」在覺行、覺知、覺悟、感覺都唸 jué，常被唸成睡覺的 jiào；只有睡覺、午覺這類才唸 jiào。
/** 把容易唸錯的字換成唸起來正確的同音字，只用在朗讀 */
export function fixPronunciation(text: string): string {
  return text
    .replace(/覺行/g, '絕形') // 「行」唸 xíng，不是銀行的 háng
    .replace(/([睡午晚一]?)覺/g, (m, sleep: string) => (sleep ? m : '絕'))
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

/** 向伺服器要一句雲端語音（MP3） */
export type CloudSpeech = (text: string) => Promise<ArrayBuffer>

export function createVoice(hooks: VoiceHooks) {
  const synth = typeof window !== 'undefined' && 'speechSynthesis' in window ? window.speechSynthesis : null
  const AudioCtx: typeof AudioContext | undefined =
    typeof window !== 'undefined'
      ? (window.AudioContext ?? (window as unknown as { webkitAudioContext?: typeof AudioContext }).webkitAudioContext)
      : undefined
  let buf = ''
  let ending = false
  let generation = 0
  let unlocked = false
  // 保留參照：Chrome 會把沒有人引用的 utterance 回收掉，之後就收不到 onend
  const pending = new Set<SpeechSynthesisUtterance>()
  synth?.getVoices() // Chrome 第一次呼叫時才開始載入語音清單

  // 雲端語音：每句一送出就先去取聲音，播放則照順序一句接一句
  let cloud: CloudSpeech | null = null
  let ctx: AudioContext | null = null
  let chain: Promise<void> = Promise.resolve()
  let queued = 0 // 還沒唸完的雲端語音句數
  let playing: AudioBufferSourceNode | null = null

  const settle = () => {
    if (ending && !pending.size && !queued) {
      ending = false
      hooks.onDone()
    }
  }

  const play = (audio: AudioBuffer, gen: number) =>
    new Promise<void>((resolve) => {
      const src = ctx!.createBufferSource()
      src.buffer = audio
      src.connect(ctx!.destination)
      src.onended = () => {
        if (playing === src) playing = null
        if (gen === generation) hooks.onEnd()
        resolve()
      }
      playing = src
      hooks.onStart()
      src.start()
    })

  const sayCloud = (text: string) => {
    const gen = generation
    const audio = cloud!(fixPronunciation(text))
      .then((b) => ctx!.decodeAudioData(b))
      .catch(() => null)
    queued++
    chain = chain
      .then(async () => {
        const a = await audio
        if (gen !== generation) return
        // 這一句取不到雲端語音，改用瀏覽器語音唸，唸完再接下一句
        if (a) await play(a, gen)
        else await new Promise<void>((resolve) => say(text, resolve))
      })
      .finally(() => {
        if (gen !== generation) return
        queued--
        settle()
      })
  }

  const speak = (raw: string) => {
    const text = cleanForSpeech(raw)
    if (!text) return
    if (cloud && ctx) sayCloud(text)
    else say(text)
  }

  /** 用瀏覽器語音唸；after 在唸完或唸不出來時呼叫 */
  const say = (text: string, after?: () => void) => {
    const voices = synth?.getVoices() ?? []
    const v = pickVoice(voices)
    // 語音清單已載入卻沒有華語語音，用英文語音唸中文只會是怪聲，直接改用文字對嘴
    if (!synth || (voices.length && !v)) {
      hooks.onFail(text)
      after?.()
      return
    }
    const u = new SpeechSynthesisUtterance(fixPronunciation(text))
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
      after?.()
      settle()
    }
    u.onend = () => done(false)
    u.onerror = (e) => done(e.error !== 'interrupted' && e.error !== 'canceled')
    pending.add(u)
    synth.speak(u)
  }

  return {
    supported: !!synth || !!AudioCtx,
    /** 伺服器有雲端語音時呼叫一次；傳 null 改回瀏覽器語音 */
    useCloud(fn: CloudSpeech | null) {
      cloud = AudioCtx ? fn : null
    },
    /** 在點擊當下呼叫：iOS Safari 只允許由使用者操作開始發聲 */
    unlock() {
      if (cloud && AudioCtx) {
        try {
          ctx ??= new AudioCtx()
          if (ctx.state === 'suspended') void ctx.resume()
          // 先放一段無聲的聲音，之後非點擊當下開始播放的聲音才不會被擋
          const src = ctx.createBufferSource()
          src.buffer = ctx.createBuffer(1, 1, 22050)
          src.connect(ctx.destination)
          src.start()
        } catch {
          ctx = null
        }
      }
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
      sentences.forEach(speak)
    },
    /** 回答寫完了：唸完剩下的半句，全部唸完後呼叫 onDone */
    end() {
      speak(buf)
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
      queued = 0
      chain = Promise.resolve()
      try {
        playing?.stop()
      } catch {
        /* 已經停了 */
      }
      playing = null
    },
  }
}

export type Voice = ReturnType<typeof createVoice>
