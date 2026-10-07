<script setup lang="ts">
// 第三章「協會營運特色分享」：核心共識。
</script>

<template>
  <section id="operation" class="chapter">
    <p class="kicker">協會營運特色分享</p>
    <h2>核心共識</h2>
  </section>

  <div class="quote">
    <p class="quote-title">不辦教育，佛教就沒有明天</p>
    <p>
      <b>聖嚴法師</b>有感於當時社會對佛教的誤解，甚至視之為迷信，雖然僧尼無數，寺院儼然，但是<span class="hl">佛教的形象卻普遍地落於民間信仰的層次</span>。
      <b>聖嚴法師</b>懇切呼籲<span class="hl">推動佛教教育事業為當務之急</span>。
    </p>
    <p class="by">1985 年</p>
  </div>
  <div class="lifelong">
    <p class="flow">教育 <span>→</span> 學習 <span>→</span> 終身</p>
    <div class="flow-notes">
      <span>未來生活之準備</span>
      <span>提升生活水平價值</span>
      <span>活到老學到老</span>
    </div>
  </div>
</template>

<style scoped>
.lifelong {
  text-align: center;
}
.flow {
  margin: 8px 0 4px;
  font-family: var(--sans);
  font-size: 2.2rem !important;
  font-weight: 700;
  color: var(--plum) !important;
}
.flow span {
  margin: 0 8px;
  color: var(--ink-faint);
}
.flow-notes {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 10px;
}
.flow-notes span {
  padding: 4px 14px;
  border-radius: 999px;
  background: var(--plum-soft);
  color: var(--plum);
  font-size: 0.95rem;
}
</style>
