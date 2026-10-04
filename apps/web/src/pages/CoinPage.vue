<script setup lang="ts">
// 數字取自企劃書 v1.0 草案的初始建議值，正式數字以菩提幣決策小組決議為準
const tiers = [
  ['半日服務', '4 小時', '1,500'],
  ['全日服務', '8 小時', '3,000'],
  ['全日・帶領禪修', '8 小時，需教練資格', '6,000'],
  ['禪修營三日全程', '24 小時', '10,000'],
  ['禪修營七日全程', '56 小時', '25,000'],
  ['零星支援', '每小時，未達半日者', '300'],
]
const steps = [
  { title: '服務', text: '在任一中心的活動中服務，依梯級表計算。' },
  { title: '核發', text: '活動結束後造冊送審，菩提幣決策小組核准後入帳。' },
  { title: '兌換', text: '在官網把菩提幣換成共好企業提供的券。' },
  { title: '核銷', text: '到店出示身份 QR，店員掃碼核銷。' },
]
</script>

<template>
  <q-page class="page">
    <h1>菩提幣介紹</h1>
    <p class="lead">
      把禪修中心的服務時數，變成共好企業共通的記帳單位。志工與禪修教練以服務換幣，在共好企業兌換住宿、餐飲與課程。
    </p>
    <p class="note">菩提幣不販售、不提領、不可兌現，只在共好企業間循環。這不是投資，也不是可交易的資產。</p>

    <h2>怎麼運作</h2>
    <div class="grid">
      <div v-for="(s, i) in steps" :key="s.title" class="card">
        <div class="status-chip">第 {{ i + 1 }} 步</div>
        <h3 class="q-my-sm">{{ s.title }}</h3>
        <p class="q-mb-none">{{ s.text }}</p>
      </div>
    </div>

    <h2>以梯級計，不以時薪計</h2>
    <p>表上的數字就是實拿的幣數，不需要再乘任何係數。以活動為單位，符合禪修中心的實務。1 菩提幣 = 1 新台幣。</p>
    <table class="table">
      <thead>
        <tr><th>梯級</th><th>投入</th><th style="text-align: right">菩提幣</th></tr>
      </thead>
      <tbody>
        <tr v-for="t in tiers" :key="t[0]">
          <td>{{ t[0] }}</td><td>{{ t[1] }}</td><td class="num">{{ t[2] }}</td>
        </tr>
      </tbody>
    </table>
    <p class="text-caption q-mt-sm">初始建議值，正式梯級表由菩提幣決策小組決議公告。</p>

    <h2>兩種券</h2>
    <div class="grid">
      <div class="card">
        <h3 class="q-mt-none">品項券</h3>
        <p>對應一項具體的商品或服務，例如禪修中心單人房一晚、三日禪修營名額、贊助飯店標準房一晚。</p>
      </div>
      <div class="card">
        <h3 class="q-mt-none">面額券</h3>
        <p>可抵用等額消費，供零星消費使用，可部分核銷，餘額留在券上。</p>
      </div>
    </div>
    <p>券的面額等於該品項的對外公開台幣售價。券到期會自動退回等額菩提幣，不沒收。</p>

    <h2>現況</h2>
    <p>
      菩提幣目前是 v1.0 草案，正處於「法務諮詢與治理籌組」階段。志工定性與券的定性兩題取得律師結論後，才會啟動錢包與兌換功能。
      現在可以先<router-link to="/groups">參加覺行小組</router-link>或<router-link to="/partners">登記成為共好企業</router-link>。
    </p>
  </q-page>
</template>
