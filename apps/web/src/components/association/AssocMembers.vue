<script setup lang="ts">
import { computed } from 'vue'
import { account } from '../../account'

// 第四章「協會會員招募中」：會員人設。
// 每種會員：[前文, 重點, 後文]
const members = [
  { role: '策略夥伴', text: ['以策略投資的方式，支援協會在初始階段及未來營運過程中的', '財務需求', '。'], tone: 'light' },
  { role: '資源協助者', text: ['提供協會', '多面向的相關必要產業資源（例如：企業培訓需求等）', '，以完成協會的設定任務。'], tone: 'mid' },
  { role: '推廣協同者', text: ['提供協會', '在會員招募、課程招生與必要輔助活動的執行', '，並藉由會員力量協助協會推廣。'], tone: 'deep' },
]

// 已經是會員（例如覺行小組）就到個人頁加入協會；還不是會員就到報名頁，預設勾選協會
const joinTo = computed(() => (account.user ? { path: '/me', query: { join: 'association' } } : { path: '/join', query: { for: 'association' } }))
</script>

<template>
  <section id="members" class="chapter members">
    <p class="kicker">協會會員招募中</p>
    <h2>產業的專業會員</h2>
  </section>

  <h2>會員人設</h2>
  <div class="member-grid">
    <div v-for="m in members" :key="m.role" class="member" :class="m.tone">
      <p class="role">{{ m.role }}</p>
      <p class="text">{{ m.text[0] }}<span class="key">{{ m.text[1] }}</span>{{ m.text[2] }}</p>
    </div>
  </div>

  <div class="join">
    <q-btn color="secondary" unelevated no-caps size="lg" :to="joinTo" label="加入會員" />
  </div>
</template>

<style scoped>
.member-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
}
.member {
  padding: 22px 20px;
  border-radius: 18px;
  box-shadow: 0 10px 24px rgba(47, 111, 147, 0.16);
}
.member.light {
  background: #d4e4f0;
}
.member.mid {
  background: #e6eff6;
}
.member.deep {
  background: #4f8fbf;
}
.role {
  margin: 0 0 10px;
  font-family: var(--wenkai);
  font-size: 1.5rem !important;
  font-weight: 700;
  color: #2f5e24 !important;
}
.member.deep .role {
  color: #fff !important;
}
.text {
  margin: 0;
  color: #2c4a63 !important;
}
.member.deep .text {
  color: #fff !important;
}
.key {
  color: var(--weba);
  font-weight: 600;
  text-decoration: underline;
  text-underline-offset: 4px;
}
.member.deep .key {
  color: #ffe36e;
}
.join {
  margin-top: 32px;
  text-align: center;
}
@media (max-width: 680px) {
  .member-grid {
    grid-template-columns: 1fr;
  }
}
</style>
