// Sunny 的湖景依觀看者裝置上的時間切換：早上六點到下午五點是藍天白雲，其餘時間是黃昏。
// 頁面打開時決定一次；網址加 ?scene=day 或 ?scene=dusk 可以固定其中一種（測試用）。

export type SceneTime = 'day' | 'dusk' | 'meadow'

export const DAY_FROM = 6 // 06:00 起是白天
export const DAY_UNTIL = 17 // 17:00 起換成黃昏

export function sceneTime(
  now: Date = new Date(),
  search: string = typeof location !== 'undefined' ? location.search : '',
): SceneTime {
  const forced = new URLSearchParams(search).get('scene')
  if (forced === 'day' || forced === 'dusk' || forced === 'meadow') return forced   // meadow：晨霧草原（預覽中）
  const h = now.getHours()
  return h >= DAY_FROM && h < DAY_UNTIL ? 'day' : 'dusk'
}
