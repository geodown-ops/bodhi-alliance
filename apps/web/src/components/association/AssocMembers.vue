<script setup lang="ts">
// 第四章「協會會員招募中」：會員人設與結語（簡報第 34–36 頁）。
// 每種會員：[前文, 重點, 後文]
const members = [
  { role: '策略夥伴', text: ['以策略投資的方式，支援協會在初始階段及未來營運過程中的', '財務需求', '。'], tone: 'light' },
  { role: '資源協助者', text: ['提供協會', '多面向的相關必要產業資源（例如：企業培訓需求等）', '，以完成協會的設定任務。'], tone: 'mid' },
  { role: '推廣協同者', text: ['提供協會', '在會員招募、課程招生與必要輔助活動的執行', '，並藉由會員力量協助協會推廣。'], tone: 'deep' },
]
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

  <div class="closing">
    <img src="/images/association/logo.webp" alt="世界佛教教育協會標誌" class="closing-logo" />
    <p>
      我們有決心與行動，讓世界佛教教育協會成長茁壯，<b>歡迎您的加入</b>
    </p>
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
.closing {
  position: relative;
  display: flex;
  align-items: center;
  gap: 18px;
  margin: 40px -16px 0;
  padding: 120px 28px;
  border-radius: 18px;
  overflow: hidden;
  background:
    linear-gradient(rgba(40, 28, 20, 0.35), rgba(40, 28, 20, 0.35)),
    url('/images/association/tree.webp') center / cover;
}
.closing::before {
  content: '';
  position: absolute;
  left: 0;
  right: 0;
  top: 50%;
  height: 78px;
  transform: translateY(-50%);
  background: rgba(171, 43, 39, 0.92);
}
.closing-logo {
  position: relative;
  flex: none;
  width: 96px;
  height: 96px;
  border-radius: 50%;
  box-shadow: 0 0 0 4px rgba(255, 255, 255, 0.35);
}
.closing p {
  position: relative;
  margin: 0;
  font-size: 1.25rem !important;
  font-weight: 700;
  line-height: 1.6 !important;
  color: #fff !important;
}
.closing b {
  color: #ffe36e;
}
@media (max-width: 680px) {
  .member-grid {
    grid-template-columns: 1fr;
  }
  .closing {
    padding: 90px 18px;
  }
  .closing::before {
    height: 120px;
  }
  .closing-logo {
    width: 64px;
    height: 64px;
  }
  .closing p {
    font-size: 1.05rem !important;
  }
}
</style>
