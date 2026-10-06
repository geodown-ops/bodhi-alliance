export type Avatar = {
  /** idle | thinking | talking */
  setState(s: 'idle' | 'thinking' | 'talking'): void
  /** 沒有聲音時：把剛串流到的字幕排進對嘴佇列 */
  speak(text: string): void
  /** 語音開始／唸完一句：唸的時候嘴巴跟著開合 */
  voice(on: boolean): void
  /** 回答結束：嘴型停下後回到冥想 */
  finish(): void
  dispose(): void
}

export function createAvatar(
  canvas: HTMLCanvasElement,
  url: string,
  /** time：day 是藍天湖景（預設），dusk 是黃昏湖景，meadow 是晨霧草原，pines 是松林雪山 */
  opts?: { onProgress?: (p: number) => void; onIdle?: () => void; time?: 'day' | 'dusk' | 'meadow' | 'pines' },
): Promise<Avatar>
