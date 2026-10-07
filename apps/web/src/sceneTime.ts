// Sunny 的場景依觀看者裝置上的時間切換（頁面打開時決定一次）：
//   首頁：06:00–22:00 晨霧草原（菩提樹版），其餘時間黃昏湖景（蘆葦）。
//   線上問答頁：08:00–16:00 松林雪山，其餘時間晨霧草原（菩提樹版）。
// 網址加 ?scene=meadow、dusk、pines 或 day（原本的藍天湖景）可以固定其中一種。

export type SceneTime = 'day' | 'dusk' | 'meadow' | 'pines'
export type ScenePage = 'home' | 'guide'

// 各頁的白天時段 [起, 迄)（小時）與時段內、時段外的場景
export const SCHEDULE: Record<ScenePage, { from: number; until: number; inside: SceneTime; outside: SceneTime }> = {
  home: { from: 6, until: 22, inside: 'meadow', outside: 'dusk' },
  guide: { from: 8, until: 16, inside: 'pines', outside: 'meadow' },
}

export function sceneTime(
  page: ScenePage = 'guide',
  now: Date = new Date(),
  search: string = typeof location !== 'undefined' ? location.search : '',
): SceneTime {
  const forced = new URLSearchParams(search).get('scene')
  if (forced === 'day' || forced === 'dusk' || forced === 'meadow' || forced === 'pines') return forced
  const s = SCHEDULE[page]
  const h = now.getHours()
  return h >= s.from && h < s.until ? s.inside : s.outside
}
