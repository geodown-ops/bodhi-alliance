// Sunny 的場景依觀看者裝置上的時間切換：早上六點到下午五點是晨霧草原，其餘時間是黃昏湖景。
// 頁面打開時決定一次；網址加 ?scene=meadow、?scene=dusk 或 ?scene=day（原本的藍天湖景）可以固定其中一種。

export type SceneTime = 'day' | 'dusk' | 'meadow'

export const DAY_FROM = 6 // 06:00 起是白天（晨霧草原）
export const DAY_UNTIL = 17 // 17:00 起換成黃昏

export function sceneTime(
  now: Date = new Date(),
  search: string = typeof location !== 'undefined' ? location.search : '',
): SceneTime {
  const forced = new URLSearchParams(search).get('scene')
  if (forced === 'day' || forced === 'dusk' || forced === 'meadow') return forced
  const h = now.getHours()
  return h >= DAY_FROM && h < DAY_UNTIL ? 'meadow' : 'dusk'
}
