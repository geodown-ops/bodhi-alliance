import { describe, expect, it } from 'vitest'
import { sceneTime } from '../sceneTime'

describe('sceneTime', () => {
  it('shows the dawn meadow by default, all day', () => {
    expect(sceneTime('')).toBe('meadow')
    expect(sceneTime('?scene=noon')).toBe('meadow')
  })

  it('can still open the old lake scenes with ?scene=', () => {
    expect(sceneTime('?scene=dusk')).toBe('dusk')
    expect(sceneTime('?scene=day')).toBe('day')
    expect(sceneTime('?scene=meadow')).toBe('meadow')
  })
})
