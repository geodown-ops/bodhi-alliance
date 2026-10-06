<script setup lang="ts">
import { account } from '../account'
</script>

<template>
  <q-page class="page">
    <h1>我的錢包</h1>
    <p class="lead">服務換來的菩提幣、兌換的券、到店核銷用的身份 QR，都在這裡。</p>
    <p><span class="status-chip">餘額與核發紀錄已開放；兌換券籌備中</span></p>
    <p v-if="account.user">你的餘額與核發紀錄在<router-link to="/me">個人頁</router-link>。</p>
    <p v-else>先<router-link to="/join">報名參加覺行小組</router-link>，協助共修活動、審核通過後，菩提幣就會進到你的錢包。</p>

    <h2>開放後可以做什麼</h2>
    <div class="grid">
      <div class="card"><h3 class="q-mt-none">查看餘額</h3><p class="q-mb-none">活動送審核准當下入帳，在個人頁就看得到。</p></div>
      <div class="card"><h3 class="q-mt-none">兌換券</h3><p class="q-mb-none">瀏覽共好企業提供的品項券與面額券，幾秒完成兌換。</p></div>
      <div class="card"><h3 class="q-mt-none">轉贈</h3><p class="q-mb-none">券可以轉贈給家人或同修；幣本身不能轉讓。</p></div>
      <div class="card"><h3 class="q-mt-none">身份 QR</h3><p class="q-mb-none">到店出示會定時更新的身份 QR，店員掃碼核銷；手機遺失可一鍵凍結。</p></div>
    </div>

    <h2>錢包放在哪裡</h2>
    <p>
      你的餘額與券記在 Sunny life 官網的帳本裡，不需要自己保管私鑰，也不用付任何手續費。區塊鏈上只有共好企業金庫、各中心與各共好企業的地址；
      官網每天把整本帳的摘要寫上鏈，各中心都能用查帳工具獨立核對，確保帳本沒有被改動。
    </p>

    <h2>錢包使用者管理</h2>
    <p>各中心的管理員負責志工與禪修教練的造冊核可，可以查詢流水、協助凍結身份 QR。這部分在管理後台操作。</p>
  </q-page>
</template>
