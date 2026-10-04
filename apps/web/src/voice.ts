// Sunny 的聲音：用瀏覽器內建的語音合成（Web Speech API）唸出回答，不需要另外的伺服器或金鑰。
// 串流進來的文字逐句切開、依序唸；每一句開始唸的時候通知呼叫端（用來對嘴）。
const synth: SpeechSynthesis | undefined = typeof window !== 'undefined' ? window.speechSynthesis : undefined

export const voiceSupported = !!synth

// 優先用台灣華語的自然女聲（Edge 的曉臻、曉雨），其次台灣華語女聲，再其次任何中文語音
let voice: SpeechSynthesisVoice | null = null
function pickVoice() {
  const all = synth?.getVoices() ?? []
  const tw = all.filter((v) => /^zh[-_](TW|Hant)/i.test(v.lang))
  const pool = tw.length ? tw : all.filter((v) => /^(zh|cmn)/i.test(v.lang))
  voice =
    pool.find((v) => /Natural/i.test(v.name) && /曉臻|曉雨|HsiaoChen|HsiaoYu/i.test(v.name)) ??
    pool.find((v) => /female|女|Hsiao|Yating|Mei-?Jia|Hanhan/i.test(v.name)) ??
    pool[0] ??
    null
}
pickVoice()
synth?.addEventListener?.('voiceschanged', pickVoice)

/** 這台裝置有沒有中文語音；沒有就不出聲（用英文語音唸中文會很怪），只照字幕對嘴 */
export function hasChineseVoice() {
  if (!voice) pickVoice()
  return !!voice
}

/** 在使用者按下送出的當下呼叫：有些瀏覽器（iOS Safari）要在點擊事件裡先發過一次聲才允許之後播放 */
export function unlockVoice() {
  if (!synth) return
  synth.cancel()
  synth.speak(new SpeechSynthesisUtterance(''))
}

export function stopVoice() {
  synth?.cancel()
}

const SENTENCE = /^[\s\S]*?[。！？!?；;\n]+/
const MAX_CHUNK = 60   // 太長的句子在逗號處先切，避免部分瀏覽器唸到一半中斷

/** 逐句唸出串流文字。push() 餵入片段；end() 表示文字到齊，全部唸完後呼叫 done */
export function createSpeaker(onSentence: (text: string) => void) {
  let buf = ''
  let pending = 0
  let ended = false
  let stopped = false
  let done: (() => void) | null = null

  const finishIfIdle = () => {
    if (ended && pending === 0 && !stopped) {
      stopped = true
      done?.()
    }
  }

  const say = (raw: string) => {
    const text = raw.replace(/[*#>`_~|]/g, '').trim()   // 去掉 Markdown 符號，免得被唸出來
    if (!text || !synth) return
    const u = new SpeechSynthesisUtterance(text)
    u.lang = voice?.lang ?? 'zh-TW'
    if (voice) u.voice = voice
    u.rate = 1
    u.pitch = 1.1
    let started = false
    u.onstart = () => {
      started = true
      if (!stopped) onSentence(text)
    }
    u.onerror = () => {
      if (!started && !stopped) onSentence(text)   // 這句沒唸出來，至少照字幕對嘴
      pending--
      finishIfIdle()
    }
    u.onend = () => {
      pending--
      finishIfIdle()
    }
    pending++
    synth.speak(u)
  }

  return {
    push(chunk: string) {
      buf += chunk
      for (let m = SENTENCE.exec(buf); m; m = SENTENCE.exec(buf)) {
        say(m[0])
        buf = buf.slice(m[0].length)
      }
      const comma = buf.lastIndexOf('，')
      if (buf.length > MAX_CHUNK && comma > 0) {
        say(buf.slice(0, comma + 1))
        buf = buf.slice(comma + 1)
      }
    },
    end(onDone: () => void) {
      say(buf)
      buf = ''
      ended = true
      done = onDone
      finishIfIdle()
    },
    cancel() {
      stopped = true
      synth?.cancel()
    },
  }
}
