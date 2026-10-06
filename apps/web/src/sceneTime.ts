// Sunny 的場景：全天都是晨霧草原（菩提樹版）。
// 網址加 ?scene=dusk（黃昏湖景）或 ?scene=day（藍天湖景）仍可打開舊場景。

export type SceneTime = 'day' | 'dusk' | 'meadow'

export function sceneTime(search: string = typeof location !== 'undefined' ? location.search : ''): SceneTime {
  const forced = new URLSearchParams(search).get('scene')
  if (forced === 'day' || forced === 'dusk' || forced === 'meadow') return forced
  return 'meadow'
}
