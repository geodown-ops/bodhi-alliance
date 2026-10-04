export type Avatar = {
  /** idle | thinking | talking */
  setState(s: 'idle' | 'thinking' | 'talking'): void
  /** 把剛串流到的字幕排進對嘴佇列 */
  speak(text: string): void
  /** 語音開始唸某一句時呼叫：從這句重新對嘴 */
  say(text: string): void
  /** 回答結束：唸完佇列裡的字後回到冥想 */
  finish(): void
  dispose(): void
}

export function createAvatar(
  canvas: HTMLCanvasElement,
  url: string,
  opts?: { onProgress?: (p: number) => void; onIdle?: () => void },
): Promise<Avatar>
