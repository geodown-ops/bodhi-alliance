<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { Notify } from 'quasar'
import { api, type ChainInfo } from '../api'

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
  { title: '服務', text: '協助覺行小組的正念減壓活動，或在任一中心的活動中服務，依梯級表計算。' },
  { title: '核發', text: '活動結束後發起人在官網送審，核准後入帳到錢包。' },
  { title: '兌換', text: '在官網把菩提幣換成共好企業提供的券。' },
  { title: '核銷', text: '到店出示身份 QR，店員掃碼核銷。' },
]

// 鏈上合約資訊由 api 提供，換到正式鏈時頁面自動跟著更新
const chain = ref<ChainInfo>({ enabled: false })
onMounted(async () => {
  try {
    chain.value = await api.chainInfo()
  } catch {
    // 連不到 api 時只顯示說明文字
  }
})
async function copyContract() {
  try {
    await navigator.clipboard.writeText(chain.value.contract ?? '')
    Notify.create({ type: 'positive', message: '已複製合約地址' })
  } catch {
    Notify.create({ type: 'negative', message: '無法複製，請手動選取地址' })
  }
}
const lookups = [
  {
    title: '看整個菩提幣',
    text: '打開合約頁，可以看到總發行量、有多少地址持有菩提幣（Holders），以及每一筆轉帳（Transfers）。',
  },
  {
    title: '看自己的菩提幣',
    text: '登入後到「我的個人頁」，在「鏈上菩提幣」點你的地址，就會在 PolygonScan 看到這個地址持有多少菩提幣。也可以複製地址，貼到 PolygonScan 最上方的搜尋框。',
  },
  {
    title: '看某一筆交易',
    text: '在個人頁點「查看這筆交易」。交易頁的 Status 顯示 Success 就代表已經寫進區塊鏈；From 是金庫地址，To 是你的地址，數量是 1 枚菩提幣。',
  },
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

    <h2>在區塊鏈上查詢菩提幣</h2>
    <p>
      每位會員入會都會得到 1 枚菩提幣，這筆紀錄寫在 Polygon 區塊鏈上。區塊鏈像一本公開的帳本，寫進去就不能被任何人改，也不需要相信官網，自己就能查證。
    </p>

    <div class="card q-mb-md">
      <h3 class="q-mt-none">菩提幣合約</h3>
      <table class="facts">
        <tbody>
          <tr><th>名稱</th><td>菩提幣（代號 BODHI）</td></tr>
          <tr><th>網路</th><td>{{ chain.network ?? 'Polygon Amoy 測試鏈' }}<span v-if="chain.chain_id" class="text-caption">（Chain ID {{ chain.chain_id }}）</span></td></tr>
          <tr>
            <th>合約地址</th>
            <td>
              <template v-if="chain.contract">
                <a :href="chain.contract_url" target="_blank" rel="noopener" class="mono addr">{{ chain.contract }}</a>
                <q-btn flat dense round size="sm" icon="content_copy" aria-label="複製合約地址" @click="copyContract" />
              </template>
              <span v-else>整理中</span>
            </td>
          </tr>
          <tr><th>總量</th><td>固定 5 億枚，部署時全部存入金庫，之後不能再增發</td></tr>
          <tr><th>小數</th><td>2 位，最小單位 0.01 枚</td></tr>
          <tr><th>轉帳規則</th><td>只能在金庫、中心、共好企業等機構地址與會員之間流動；會員之間不能互轉，這條規則寫在合約裡</td></tr>
        </tbody>
      </table>
    </div>

    <h3>什麼是 PolygonScan</h3>
    <p>
      PolygonScan 是 Polygon 區塊鏈的公開查詢網站（區塊瀏覽器），由經營 Etherscan 的團隊維護，可以把它想成鏈上帳本的「查帳網站」。不用註冊、不用下載錢包，打開網頁就能查任何地址的餘額、任何一筆交易的時間與金額，以及合約本身的資訊。菩提幣現在在 Polygon 的測試鏈 Amoy 上運作，對應的網址是<a :href="chain.explorer ?? 'https://amoy.polygonscan.com'" target="_blank" rel="noopener">{{ (chain.explorer ?? 'https://amoy.polygonscan.com').replace('https://', '') }}</a>；正式上線後會改到 Polygon 主網，網址是 polygonscan.com。網站介面是英文。
    </p>

    <h3>怎麼查</h3>
    <div class="grid three">
      <div v-for="(s, i) in lookups" :key="s.title" class="card">
        <div class="status-chip">{{ i + 1 }}</div>
        <h3 class="q-my-sm">{{ s.title }}</h3>
        <p class="q-mb-none">{{ s.text }}</p>
      </div>
    </div>

    <h3>看懂幾個名詞</h3>
    <ul>
      <li><b>地址（Address）</b>：0x 開頭、共 42 個字元的一串英數字，就像帳號。每位會員都有一個，由平台替你保管，不用自己記私鑰。</li>
      <li><b>交易雜湊（Transaction Hash）</b>：每筆交易的編號，也是 0x 開頭，用它可以找到那一筆交易。</li>
      <li><b>區塊（Block）</b>：鏈上每隔幾秒把新交易打包成一個區塊，交易所在的區塊編號越早，代表越早寫入。</li>
      <li><b>Token</b>：PolygonScan 上把菩提幣這類合約發行的幣叫做 Token，合約頁會標示 BODHI。</li>
    </ul>
    <p class="note">目前在測試鏈上的菩提幣沒有任何市場價值，也無法在交易所買賣；正式上線前，合約地址會在本頁更新公告。</p>

    <h2>現況</h2>
    <p>
      菩提幣目前是 v1.0 草案，正處於「法務諮詢與治理籌組」階段。覺行小組活動的登錄、送審與錢包餘額已經開放；志工定性與券的定性兩題取得律師結論後，才會啟動兌換功能。
      現在可以先<router-link to="/groups">參加覺行小組</router-link>（隨興或定期相約一起做正念減壓，任何人都可發起，三人以上即可）或<router-link to="/partners">登記成為共好企業</router-link>。
    </p>
  </q-page>
</template>

<style scoped>
.facts {
  border-collapse: collapse;
  width: 100%;
}
.facts th,
.facts td {
  padding: 6px 0;
  text-align: left;
  vertical-align: top;
  border-top: 1px solid var(--rule, rgba(0, 0, 0, 0.08));
}
.facts tr:first-child th,
.facts tr:first-child td {
  border-top: none;
}
.facts th {
  width: 6.5em;
  font-weight: 600;
  white-space: nowrap;
  padding-right: 12px;
}
.grid.three {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}
@media (max-width: 800px) {
  .grid.three {
    grid-template-columns: 1fr;
  }
}
.mono {
  font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
}
.addr {
  word-break: break-all;
}
</style>
