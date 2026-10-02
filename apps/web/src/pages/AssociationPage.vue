<script setup lang="ts">
const turns = [
  { name: '轉念', pinyin: 'zhuǎn niàn', text: '放下壓力，回歸中道平衡。不以人情與虧欠維持一個團體。' },
  { name: '轉識', pinyin: 'zhuǎn shí', text: '透過呼吸覺察，安定身心。看清資源與付出真正流向哪裡。' },
  { name: '轉依', pinyin: 'zhuǎn yī', text: '觀想與身心覺知統合，進入與法界合一。一間中心的量能，成為全聯盟的量能。' },
]
</script>

<template>
  <q-page class="page">
    <h1>世界佛教教育協會</h1>
    <p class="lead">修行的核心，是把覺察落在日常裡：轉念、轉識、轉依。</p>

    <h2>宗旨：實修唯識三轉</h2>
    <div class="grid">
      <div v-for="t in turns" :key="t.name" class="card">
        <h3 class="q-my-none">{{ t.name }}</h3>
        <div class="status-chip q-my-sm">{{ t.pinyin }}</div>
        <p class="q-mb-none">{{ t.text }}</p>
      </div>
    </div>
    <p class="q-mt-md">菩提幣聯盟是這三轉在「共同生活」這一層的實作，讓服務、供養與受用在同一個身體裡循環。</p>

    <h2>緣起</h2>
    <p>
      太虛大師於 1923 年創立，旨在推動佛教國際化；後由弘化大和尚於 2011 年復辦並改名，現推動全球佛教交流、文化活動與和平倡議。
      世界禪修中心、彌勒心流靜坐、1BN.AI 世界靜坐日，以及菩提幣與時間銀行，都在同一個宗旨底下。
    </p>
    <q-btn outline color="secondary" no-caps href="https://www.sunnylife.world" target="_blank" label="前往 sunnylife.world" icon-right="open_in_new" />
  </q-page>
</template>
