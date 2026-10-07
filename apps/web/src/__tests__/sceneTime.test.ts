import { describe, expect, it } from 'vitest'
import { sceneTime } from '../sceneTime'

const at = (h: number, m = 0) => new Date(2026, 9, 7, h, m)

describe('sceneTime', () => {
  it('home: dawn meadow from 06:00 to 22:00, dusk lake otherwise', () => {
    expect(sceneTime('home', at(5, 59), '')).toBe('dusk')
    expect(sceneTime('home', at(6), '')).toBe('meadow')
    expect(sceneTime('home', at(21, 59), '')).toBe('meadow')
    expect(sceneTime('home', at(22), '')).toBe('dusk')
    expect(sceneTime('home', at(0), '')).toBe('dusk')
  })

  it('guide page: snowy pines from 08:00 to 16:00, dawn meadow otherwise', () => {
    expect(sceneTime('guide', at(7, 59), '')).toBe('meadow')
    expect(sceneTime('guide', at(8), '')).toBe('pines')
    expect(sceneTime('guide', at(15, 59), '')).toBe('pines')
    expect(sceneTime('guide', at(16), '')).toBe('meadow')
    expect(sceneTime('guide', at(23), '')).toBe('meadow')
  })

  it('?scene= fixes a scene on any page at any time', () => {
    expect(sceneTime('home', at(12), '?scene=dusk')).toBe('dusk')
    expect(sceneTime('home', at(23), '?scene=pines')).toBe('pines')
    expect(sceneTime('guide', at(12), '?scene=meadow')).toBe('meadow')
    expect(sceneTime('guide', at(3), '?scene=day')).toBe('day')
    expect(sceneTime('home', at(12), '?scene=noon')).toBe('meadow')
  })
})
