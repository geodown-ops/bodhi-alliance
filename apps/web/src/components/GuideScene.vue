<script setup lang="ts">
// 組長所在的畫境：解說員的 3D 黃昏湖景與 VRM 人物（three.js）。
// 下面的 CSS／SVG 插畫是底圖：模型載入前、或瀏覽器不支援 WebGL 時看到的就是它。
import { onBeforeUnmount, onMounted, ref } from 'vue'
import type { Avatar } from '../guide3d/avatar.js'

defineProps<{ state: 'idle' | 'listening' | 'speaking' }>()

const canvas = ref<HTMLCanvasElement>()
const progress = ref(0)
const phase = ref<'loading' | 'ready' | 'fallback'>('loading')
let avatar: Avatar | null = null
let gone = false

onMounted(async () => {
  try {
    const { createAvatar } = await import('../guide3d/avatar.js')
    const a = await createAvatar(canvas.value!, '/models/bodhi.vrm', { onProgress: (p) => (progress.value = p) })
    if (gone) return a.dispose()
    avatar = a
    phase.value = 'ready'
  } catch (e) {
    console.error('guide avatar', e)
    phase.value = 'fallback'
  }
})
onBeforeUnmount(() => {
  gone = true
  avatar?.dispose()
})

// 給對話框呼叫：思考時半閉眼、回答時鏡頭推近對嘴、結束後回到冥想
const leaf =
  'M93.6 0C91.2 16.8 85.2 31.2 75.6 40.8C45.6 69.6 0 100.8 0 153.6C0 204 48 230.4 84 231.6C88.8 231.6 92.4 228 93.6 224.4C94.8 228 98.4 231.6 103.2 231.6C139.2 230.4 187.2 204 187.2 153.6C187.2 100.8 141.6 69.6 111.6 40.8C102 31.2 96 16.8 93.6 0Z'

defineExpose({
  think: () => avatar?.setState('thinking'),
  speak: (text: string) => avatar?.speak(text),
  voice: (on: boolean) => avatar?.voice(on),
  finish: () => avatar?.finish(),
  rest: () => avatar?.setState('idle'),
})
</script>

<template>
  <div :class="['scene', state, phase]" aria-hidden="true">
    <div class="sun" />
    <svg class="hills" viewBox="0 0 1200 200" preserveAspectRatio="none">
      <path d="M0 140 C150 90 260 120 380 100 S620 60 760 95 S1020 70 1200 110 V200 H0Z" fill="#b99a7a" opacity=".55" />
      <path d="M0 165 C200 130 330 150 520 135 S860 120 1000 140 S1130 135 1200 145 V200 H0Z" fill="#8f7357" opacity=".7" />
    </svg>
    <div class="lake">
      <div class="glint" />
    </div>
    <svg class="guide" viewBox="0 0 120 120">
      <ellipse class="halo" cx="60" cy="58" rx="54" ry="54" />
      <circle cx="60" cy="30" r="12" fill="#3b2a20" />
      <path d="M60 44 C44 44 38 58 36 74 L22 92 C40 100 80 100 98 92 L84 74 C82 58 76 44 60 44Z" fill="#3b2a20" />
      <ellipse cx="60" cy="98" rx="44" ry="5" fill="#3b2a20" opacity=".25" />
    </svg>
    <svg class="reeds left" viewBox="0 0 200 260" preserveAspectRatio="xMinYMax meet">
      <g stroke="#f6f2e8" stroke-width="2" fill="none" opacity=".85">
        <path d="M20 260 C24 180 30 120 46 40" />
        <path d="M44 260 C46 190 56 140 70 80" />
        <path d="M70 260 C70 200 74 170 92 120" />
        <path d="M8 260 C8 210 4 170 0 130" />
      </g>
      <g fill="#f6f2e8" opacity=".9">
        <ellipse cx="46" cy="38" rx="4" ry="16" transform="rotate(14 46 38)" />
        <ellipse cx="70" cy="78" rx="4" ry="14" transform="rotate(18 70 78)" />
        <ellipse cx="92" cy="118" rx="3.5" ry="12" transform="rotate(24 92 118)" />
      </g>
    </svg>
    <svg class="reeds right" viewBox="0 0 200 260" preserveAspectRatio="xMaxYMax meet">
      <g stroke="#f6f2e8" stroke-width="2" fill="none" opacity=".85">
        <path d="M180 260 C176 190 168 130 150 60" />
        <path d="M156 260 C154 200 142 160 126 110" />
        <path d="M196 260 C198 220 200 190 200 160" />
      </g>
      <g fill="#f6f2e8" opacity=".9">
        <ellipse cx="150" cy="58" rx="4" ry="16" transform="rotate(-14 150 58)" />
        <ellipse cx="126" cy="108" rx="3.5" ry="13" transform="rotate(-20 126 108)" />
      </g>
    </svg>
    <canvas ref="canvas" class="stage3d" />
    <div class="loading">
      <svg viewBox="0 0 188 232"><path :d="leaf" /></svg>
      <span>覺行小組長Sunny準備中…</span>
      <span class="bar"><span :style="{ width: `${Math.round(progress * 100)}%` }" /></span>
    </div>
  </div>
</template>

<style scoped>
.scene {
  position: absolute;
  inset: 0;
  overflow: hidden;
  background: linear-gradient(180deg, #9fbfd6 0%, #d9cbbb 38%, #f4d3a6 52%, #c79c63 100%);
}
.sun {
  position: absolute;
  left: 50%;
  top: 40%;
  width: 46vmin;
  height: 46vmin;
  transform: translate(-50%, -50%);
  border-radius: 50%;
  background: radial-gradient(circle, rgba(255, 236, 200, 0.9) 0%, rgba(255, 220, 170, 0.35) 40%, transparent 70%);
}
.hills {
  position: absolute;
  left: 0;
  right: 0;
  top: 34%;
  width: 100%;
  height: 18%;
}
.lake {
  position: absolute;
  left: 0;
  right: 0;
  top: 52%;
  bottom: 0;
  background: linear-gradient(180deg, #e7c99c 0%, #b6a48c 35%, #7f8a8c 100%);
}
.glint {
  position: absolute;
  left: 50%;
  top: 0;
  width: 18vmin;
  height: 60%;
  transform: translateX(-50%);
  background: linear-gradient(180deg, rgba(255, 240, 210, 0.75), transparent);
  filter: blur(6px);
  animation: shimmer 6s ease-in-out infinite;
}
.guide {
  position: absolute;
  left: 50%;
  top: 52%;
  width: min(22vmin, 180px);
  transform: translate(-50%, -78%);
  animation: breathe 6s ease-in-out infinite;
  transform-origin: 50% 90%;
}
.halo {
  fill: rgba(184, 216, 160, 0);
  transition: fill 0.8s ease;
}
.listening .halo {
  fill: rgba(184, 216, 160, 0.25);
}
.speaking .halo {
  fill: rgba(184, 216, 160, 0.35);
  animation: glow 1.6s ease-in-out infinite;
}
.reeds {
  position: absolute;
  bottom: 0;
  height: 45%;
  width: auto;
}
.reeds.left {
  left: 0;
}
.reeds.right {
  right: 0;
}
/* 手機上對話框佔掉下半部，把地平線與人物往上移 */
@media (max-width: 599px) {
  .sun { top: 26%; }
  .hills { top: 20%; }
  .lake { top: 38%; }
  .guide { top: 38%; }
}
@keyframes breathe {
  0%, 100% { transform: translate(-50%, -78%) scale(1); }
  50% { transform: translate(-50%, -78%) scale(1.025); }
}
@keyframes glow {
  0%, 100% { opacity: 0.6; }
  50% { opacity: 1; }
}
@keyframes shimmer {
  0%, 100% { opacity: 0.7; }
  50% { opacity: 1; }
}

.stage3d {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  display: block;
  opacity: 0;
  transition: opacity 0.9s;
}
.ready .stage3d {
  opacity: 1;
}
.loading {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  background: #3b2a20;
  color: #cfc3b5;
  font-size: 13px;
  letter-spacing: 0.08em;
  transition: opacity 0.9s;
}
.ready .loading,
.fallback .loading {
  opacity: 0;
  pointer-events: none;
}
.loading svg {
  width: 64px;
  height: 79px;
  fill: #b8d8a0;
  animation: leaf 2.4s ease-in-out infinite;
}
.bar {
  width: 160px;
  height: 2px;
  background: #5c4738;
  overflow: hidden;
}
.bar span {
  display: block;
  height: 100%;
  background: #b8d8a0;
  transition: width 0.2s;
}
@keyframes leaf {
  50% { transform: scale(1.05); opacity: 0.85; }
}
@media (prefers-reduced-motion: reduce) {
  .guide, .glint, .speaking .halo, .loading svg {
    animation: none;
  }
}
</style>
