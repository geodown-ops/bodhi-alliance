import { describe, expect, it } from 'vitest'
import { attitudes, lessons, parts, references } from '../mindfulness'

describe('mindfulness course', () => {
  it('has eight lessons in order, two per foundation', () => {
    expect(lessons.map((w) => w.n)).toEqual([1, 2, 3, 4, 5, 6, 7, 8])
    for (const p of parts) for (const n of p.lessons) expect(lessons[n - 1].part).toBe(p.name)
  })

  it('lists the nine attitudes', () => {
    expect(attitudes).toHaveLength(9)
  })

  it('cites only known references', () => {
    for (const w of lessons) for (const s of w.sections) for (const k of s.refs ?? []) expect(references[k], k).toBeDefined()
  })
})
