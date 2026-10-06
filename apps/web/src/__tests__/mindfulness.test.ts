import { describe, expect, it } from 'vitest'
import { attitudes, parts, references, weeks } from '../mindfulness'

describe('mindfulness course', () => {
  it('has eight weeks in order, two per foundation', () => {
    expect(weeks.map((w) => w.n)).toEqual([1, 2, 3, 4, 5, 6, 7, 8])
    for (const p of parts) for (const n of p.weeks) expect(weeks[n - 1].part).toBe(p.name)
  })

  it('lists the nine attitudes', () => {
    expect(attitudes).toHaveLength(9)
  })

  it('cites only known references', () => {
    for (const w of weeks) for (const s of w.sections) for (const k of s.refs ?? []) expect(references[k], k).toBeDefined()
  })
})
