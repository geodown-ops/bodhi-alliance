import { describe, expect, it } from 'vitest'
import { cleanForSpeech, fixPronunciation, pickVoice, splitSentences } from '../voice'

describe('voice', () => {
  it('speaks finished sentences and keeps the unfinished half', () => {
    expect(splitSentences('我是線上覺行小組組長Sunny。可以問我')).toEqual({
      sentences: ['我是線上覺行小組組長Sunny。'],
      rest: '可以問我',
    })
  })

  it('breaks a long sentence at a comma', () => {
    const long = '一'.repeat(70) + '，' + '二'.repeat(5)
    expect(splitSentences(long)).toEqual({ sentences: ['一'.repeat(70) + '，'], rest: '二'.repeat(5) })
  })

  it('reads 覺 as jué except in 睡覺', () => {
    expect(fixPronunciation('歡迎加入覺行小組，覺察呼吸；睡覺前、午覺後也可以練')).toBe(
      '歡迎加入絕行小組，絕察呼吸；睡覺前、午覺後也可以練',
    )
  })

  it('does not read links, markdown or emoji aloud', () => {
    expect(cleanForSpeech('**第一步**：到 https://example.org 報名 🙏')).toBe('第一步：到 報名')
  })
  it('picks a natural Taiwanese Mandarin voice', () => {
    const v = (name: string, lang: string) => ({ name, lang })
    const chrome = [v('Google US English', 'en-US'), v('Microsoft Hanhan - Chinese (Traditional, Taiwan)', 'zh-TW'), v('Google 國語（臺灣）', 'zh-TW')]
    expect(pickVoice(chrome)?.name).toBe('Google 國語（臺灣）')
    const edge = [v('Microsoft YunJhe Online (Natural) - Chinese (Taiwanese Mandarin)', 'zh-TW'), v('Microsoft HsiaoChen Online (Natural) - Chinese (Taiwanese Mandarin)', 'zh-TW')]
    expect(pickVoice(edge)?.name).toContain('HsiaoChen')
    expect(pickVoice([v('Sinji', 'zh-HK'), v('Tingting', 'zh-CN')])?.name).toBe('Tingting')
    expect(pickVoice([v('Samantha', 'en-US')])).toBeUndefined()
  })
})
