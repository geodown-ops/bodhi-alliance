<script setup lang="ts">
// 組長所在的畫境。目前是 CSS／SVG 的暮色湖景與靜坐剪影；
// 解說員的 3D 湖景與人物模型到位後，換掉這個元件的內容即可，props 不變。
defineProps<{ state: 'idle' | 'listening' | 'speaking' }>()
</script>

<template>
  <div :class="['scene', state]" aria-hidden="true">
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
@media (prefers-reduced-motion: reduce) {
  .guide, .glint, .speaking .halo {
    animation: none;
  }
}
</style>
