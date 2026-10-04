import { describe, expect, it } from 'vitest'
import { sceneTime } from '../sceneTime'

const at = (h: number, m = 0) => new Date(2026, 9, 4, h, m)

describe('sceneTime', () => {
  it('shows the blue-sky lake from 6am until 5pm', () => {
    expect(sceneTime(at(6, 0), '')).toBe('day')
    expect(sceneTime(at(12, 30), '')).toBe('day')
    expect(sceneTime(at(16, 59), '')).toBe('day')
  })

  it('shows the dusk lake at other times', () => {
    expect(sceneTime(at(5, 59), '')).toBe('dusk')
    expect(sceneTime(at(17, 0), '')).toBe('dusk')
    expect(sceneTime(at(23, 0), '')).toBe('dusk')
    expect(sceneTime(at(0, 0), '')).toBe('dusk')
  })

  it('can be pinned with ?scene= for testing', () => {
    expect(sceneTime(at(12), '?scene=dusk')).toBe('dusk')
    expect(sceneTime(at(22), '?scene=day')).toBe('day')
    expect(sceneTime(at(22), '?scene=noon')).toBe('dusk')
  })
})
