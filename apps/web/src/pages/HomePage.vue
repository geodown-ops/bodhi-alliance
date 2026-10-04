<script setup lang="ts">
// 首頁：最上面是線上問答的 Sunny（畫境＋單一對話框），往下捲是原本的招募頁（GitHub Pages 版），配色換成官網的深咖啡 × 淺綠。
import { nextTick, reactive, ref } from 'vue'
import { api, ApiError } from '../api'
import GuideStage from '../components/GuideStage.vue'

const roles = [
  { label: '禪修中心／道場', kind: 'center' },
  { label: '贊助商家（飯店・水療・餐飲・商店）', kind: 'sponsor' },
  { label: '志工', kind: 'other' },
  { label: '禪修教練', kind: 'other' },
  { label: '其他（請在留言說明）', kind: 'other' },
]
const blank = { role: '', name: '', org: '', email: '', phone: '', area: '', scale: '', msg: '', website: '' }
const form = reactive({ ...blank })
const err = ref('')
const sending = ref(false)
const sentRole = ref('')
const receipt = ref<HTMLElement>()

function go(id: string) {
  document.getElementById(id)?.scrollIntoView({ behavior: 'smooth', block: 'start' })
}

async function submit() {
  const missing = []
  if (!form.role) missing.push('身份')
  if (!form.name.trim()) missing.push('姓名')
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email.trim())) missing.push('可用的電子郵件')
  if (missing.length) {
    err.value = '還缺：' + missing.join('、')
    return
  }
  err.value = ''
  const role = roles.find((r) => r.label === form.role)!
  sending.value = true
  try {
    await api.registerPartner({
      kind: role.kind,
      org_name: form.org.trim() || form.name.trim(),
      contact_name: form.name.trim(),
      email: form.email.trim(),
      phone: form.phone.trim(),
      region: form.area.trim(),
      monthly_scale: form.scale.trim(),
      message: role.kind === 'other' ? `身份：${role.label}\n${form.msg.trim()}`.trim() : form.msg.trim(),
      website: form.website,
    })
    sentRole.value = form.role
    await nextTick()
    receipt.value?.scrollIntoView({ behavior: 'smooth', block: 'start' })
  } catch (e) {
    err.value = e instanceof ApiError ? e.message : '送出失敗，請稍後再試'
  } finally {
    sending.value = false
  }
}

function again() {
  Object.assign(form, blank)
  sentRole.value = ''
}
</script>

<template>
  <q-page class="recruit">
    <!-- ============ 線上問答 Sunny ============ -->
    <section class="home-guide" aria-label="線上問答">
      <h1 class="home-title">一即一切<br />一切即一</h1>
      <GuideStage compact />
    </section>

    <!-- ============ 開頭 ============ -->
    <header class="hero">
      <div class="shell">
        <div class="hero-copy">
          <svg class="mark" viewBox="0 0 187.2 231.6" role="img" aria-label="菩提幣 logo">
            <path d="M93.6 0C91.2 16.8 85.2 31.2 75.6 40.8C45.6 69.6 0 100.8 0 153.6C0 204 48 230.4 84 231.6C88.8 231.6 92.4 228 93.6 224.4C94.8 228 98.4 231.6 103.2 231.6C139.2 230.4 187.2 204 187.2 153.6C187.2 100.8 141.6 69.6 111.6 40.8C102 31.2 96 16.8 93.6 0Z" />
          </svg>
          <p class="eyebrow">Sunny life · 成員招募</p>
          <div class="hairline"></div>
          <p class="hero-lede">
            一個人的服務，回到所有人身上；所有人的供養，回到每一個人身上。<b>菩提幣</b>把禪修中心的服務時數，變成共好企業共通的記帳單位——志工與禪修教練以服務換幣，在共好企業兌換住宿、餐飲與療程。
          </p>
          <div class="cta-row">
            <a class="btn" href="#join" @click.prevent="go('join')">加入共好企業</a>
            <a class="btn btn-ghost" href="#how" @click.prevent="go('how')">先看看怎麼運作</a>
          </div>
          <p class="hero-note">幣不販售 · 不提領 · 不可兌現 · 只在共好企業間循環</p>
        </div>

        <figure class="ring">
          <svg viewBox="0 0 560 560" role="img" aria-label="菩提幣的四段循環：服務、核發、菩提幣入帳、兌換成券、核銷、每月歸集，幣總量固定只在共好企業間循環">
            <defs>
              <marker id="ar" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="7" markerHeight="7" orient="auto-start-reverse">
                <path d="M 0 0 L 10 5 L 0 10 z" style="fill: var(--band-rule)" />
              </marker>
            </defs>

            <circle cx="280" cy="280" r="150" fill="none" stroke-width="1" stroke-dasharray="2 6" opacity=".55" style="stroke: var(--band-rule)" />

            <path class="r-arc" d="M 311.2 133.3 A 150 150 0 0 1 426.7 248.8" marker-end="url(#ar)" />
            <path class="r-arc" d="M 426.7 311.2 A 150 150 0 0 1 311.2 426.7" marker-end="url(#ar)" />
            <path class="r-arc" d="M 248.8 426.7 A 150 150 0 0 1 133.3 311.2" marker-end="url(#ar)" />
            <path class="r-arc" d="M 133.3 248.8 A 150 150 0 0 1 248.8 133.3" marker-end="url(#ar)" />

            <text class="r-act" x="418" y="141" text-anchor="middle">核發</text>
            <text class="r-actsub" x="418" y="159" text-anchor="middle">決策小組審核名單</text>

            <text class="r-act" x="418" y="412" text-anchor="middle">兌換</text>
            <text class="r-actsub" x="418" y="430" text-anchor="middle">在官網換券・幾秒完成</text>

            <text class="r-act" x="142" y="412" text-anchor="middle">核銷</text>
            <text class="r-actsub" x="142" y="430" text-anchor="middle">現場掃身份 QR</text>

            <text class="r-act" x="142" y="141" text-anchor="middle">歸集</text>
            <text class="r-actsub" x="142" y="159" text-anchor="middle">每月・純贊助</text>

            <circle class="r-node" cx="280" cy="130" r="48" />
            <text class="r-name" x="280" y="126" text-anchor="middle">服務</text>
            <text class="r-sub" x="280" y="146" text-anchor="middle">時數・梯級</text>

            <circle class="r-node" cx="430" cy="280" r="48" />
            <text class="r-name" x="430" y="276" text-anchor="middle">菩提幣</text>
            <text class="r-sub" x="430" y="296" text-anchor="middle">入帳到錢包</text>

            <circle class="r-node" cx="280" cy="430" r="48" />
            <text class="r-name" x="280" y="426" text-anchor="middle">券</text>
            <text class="r-sub" x="280" y="446" text-anchor="middle">住宿・餐飲・課程</text>

            <circle class="r-node" cx="130" cy="280" r="48" />
            <text class="r-name" x="130" y="270" text-anchor="middle">供應方</text>
            <text class="r-sub" x="130" y="290" text-anchor="middle">共好企業</text>
            <text class="r-sub" x="130" y="304" text-anchor="middle">中心・贊助商家</text>

            <text class="r-mid" x="280" y="272" text-anchor="middle">幣總量固定</text>
            <text class="r-mid" x="280" y="296" text-anchor="middle">全程無新台幣流動</text>
          </svg>
        </figure>

        <div class="ring-list" aria-hidden="true">
          <div><span class="rs">服務</span><span><span class="rn">投入時數</span><br /><span class="rd">依梯級表計算</span></span></div>
          <div><span class="rs">核發</span><span><span class="rn">菩提幣入帳</span><br /><span class="rd">決策小組審核名單</span></span></div>
          <div><span class="rs">兌換</span><span><span class="rn">在官網換成券</span><br /><span class="rd">住宿・餐飲・課程</span></span></div>
          <div><span class="rs">核銷</span><span><span class="rn">現場掃身份 QR</span><br /><span class="rd">供應方提供服務</span></span></div>
          <div><span class="rs">歸集</span><span><span class="rn">每月回到共好企業金庫</span><br /><span class="rd">純贊助・不換現金</span></span></div>
        </div>
      </div>
    </header>

    <!-- ============ 宗旨 ============ -->
    <section class="band" id="purpose">
      <div class="shell">
        <div class="sec-head">
          <p class="eyebrow">宗旨</p>
          <h2>實修唯識三轉</h2>
          <p>世界佛教教育協會的修行核心，是把覺察落在日常裡：轉念、轉識、轉依。Sunny life 是這三轉在「共同生活」這一層的實作——讓服務、供養與受用在同一個身體裡循環。</p>
        </div>

        <div class="triad">
          <div>
            <div class="glyph">轉念</div>
            <div class="romaji">zhuǎn niàn</div>
            <p>放下壓力，回歸中道平衡。不以人情與虧欠維持一個團體。</p>
          </div>
          <div>
            <div class="glyph">轉識</div>
            <div class="romaji">zhuǎn shí</div>
            <p>透過呼吸覺察，安定身心。看清資源與付出真正流向哪裡。</p>
          </div>
          <div>
            <div class="glyph">轉依</div>
            <div class="romaji">zhuǎn yī</div>
            <p>觀想與身心覺知統合，進入與法界合一。一間中心的量能，成為所有共好企業的量能。</p>
          </div>
        </div>

        <div class="origin">
          <div class="k">緣起</div>
          <div>
            <p>太虛大師於 1923 年創立，旨在推動佛教國際化；後由弘化大和尚於 2011 年復辦並改名，現推動全球佛教交流、文化活動與和平倡議。世界禪修中心、彌勒心流靜坐、1BN.AI 世界靜坐日，以及菩提幣與時間銀行，都在同一個宗旨底下。</p>
          </div>
        </div>
      </div>
    </section>

    <!-- ============ 為什麼 ============ -->
    <section class="band" id="why">
      <div class="shell">
        <div class="sec-head">
          <p class="eyebrow">為什麼需要共好企業</p>
          <h2>感謝說得再多，<br />都補不上三個缺口</h2>
        </div>

        <div class="gaps">
          <div>
            <h3>留任靠人情，不靠制度</h3>
            <p>核心志工一旦離開，服務量能立刻出現缺口。中心沒有任何可衡量、可累積的方式，承認他們付出的時間。</p>
          </div>
          <div>
            <h3>各中心各自為政</h3>
            <p>甲中心的志工到乙中心支援，時數不被承認、也換不到任何東西。共好企業的整體量能因此無法互相調度。</p>
          </div>
          <div>
            <h3>資源明明就閒置著</h3>
            <p>淡季的空房、餐飲量能、課程名額，邊際成本極低，卻沒有管道回饋給最該得到它們的人。</p>
          </div>
        </div>

        <p class="pull">用一個共好企業共通的記帳單位，把「服務時間」與「閒置量能」對接起來。</p>
      </div>
    </section>

    <!-- ============ 怎麼運作 ============ -->
    <section class="band" id="how">
      <div class="shell">
        <div class="sec-head">
          <p class="eyebrow">怎麼運作</p>
          <h2>以梯級計，不以時薪計</h2>
          <p>菩提幣決策小組公告一張梯級表，<strong>表上的數字就是實拿的幣數</strong>，不需要再乘任何係數。以活動為單位，符合禪修中心的實務。1 菩提幣 = 1 新台幣。</p>
        </div>

        <div class="tw">
          <table>
            <caption>服務梯級表（初始建議值，由菩提幣決策小組決議）</caption>
            <thead>
              <tr><th>梯級</th><th>投入</th><th class="num">菩提幣</th></tr>
            </thead>
            <tbody>
              <tr><td>半日服務</td><td>4 小時</td><td class="num">1,500</td></tr>
              <tr><td>全日服務</td><td>8 小時</td><td class="num">3,000</td></tr>
              <tr><td>全日・帶領禪修</td><td>8 小時，需教練資格</td><td class="num">6,000</td></tr>
              <tr><td>禪修營三日全程</td><td>24 小時</td><td class="num">10,000</td></tr>
              <tr><td>禪修營七日全程</td><td>56 小時</td><td class="num">25,000</td></tr>
              <tr><td>零星支援</td><td>每小時，未達半日者</td><td class="num">300</td></tr>
            </tbody>
          </table>
        </div>
        <p class="after-table">梯級<strong>刻意不成比例</strong>：三日全程 10,000 幣高於三個全日的 9,000 幣，用來鼓勵長天期的完整投入——這正是禪修營最需要、也最難找到人的部分。</p>

        <div class="vouchers">
          <div>
            <h4>品項券</h4>
            <p>對應一項具體的商品或服務，兌換時即鎖定一份供應量。</p>
            <ul>
              <li>禪修中心單人房一晚</li>
              <li>水療中心 60 分鐘療程</li>
              <li>三日禪修營名額</li>
              <li>贊助飯店標準房一晚</li>
            </ul>
          </div>
          <div>
            <h4>面額券</h4>
            <p>可抵用等額消費，供零星消費使用。可部分核銷，餘額留在券上。</p>
            <ul>
              <li>餐飲 200 / 500 / 1,000 元券</li>
              <li>商店 500 元券</li>
              <li>結緣品抵用券</li>
            </ul>
          </div>
        </div>
        <p class="after-table">券的面額等於該品項的對外公開台幣售價，商家不必自訂。<strong>有效期到了自動退回等額菩提幣，不沒收</strong>；未核銷的券不結算給商家——服務尚未提供，不應先拿到幣。</p>
      </div>
    </section>

    <!-- ============ 誰可以加入 ============ -->
    <section class="band" id="who">
      <div class="shell">
        <div class="sec-head">
          <p class="eyebrow">誰可以加入</p>
          <h2>三種身份，一個循環</h2>
          <p>首波規模為 10 間中心、約 2,000 位活躍服務者、30 家外部贊助商家。試辦期先從 3 間中心、約 300 位服務者開始。</p>
        </div>

        <div class="roles">
          <div>
            <div class="rk">Role 一</div>
            <h3>禪修中心與道場</h3>
            <p>核發菩提幣給自己的志工與教練，同時把住宿、餐飲與課程名額開放給所有共好企業使用。推派委員進入菩提幣決策小組，共同持有金庫的多簽鑰匙。</p>
            <div class="gives"><span>核發額度</span><span>決策小組席次</span><span>自主查帳工具</span><span>開放供應</span></div>
          </div>
          <div>
            <div class="rk">Role 二</div>
            <h3>贊助商家</h3>
            <p>飯店、水療、餐飲、實體商店。在後台自行建立券種與每月贊助額度，額度用罄即暫停接受新兌換，已在志工手上的券照常核銷。不計入任何中心的配額，屬所有共好企業共享的紅利。</p>
            <div class="gives"><span>自建券種</span><span>每月額度自訂</span><span>掃碼核銷</span><span>對帳單</span></div>
          </div>
          <div>
            <div class="rk">Role 三</div>
            <h3>志工與禪修教練</h3>
            <p>在任一中心服務，時數在所有共好企業通認。核准當下入帳、手機收到通知；在官網把幣換成券，到店出示身份 QR 即可核銷。</p>
            <div class="gives"><span>跨中心通認</span><span>手機兌換</span><span>券可轉贈</span><span>一鍵凍結</span></div>
          </div>
        </div>
      </div>
    </section>

    <!-- ============ 流程與現況 ============ -->
    <section class="band" id="stage">
      <div class="shell">
        <div class="sec-head">
          <p class="eyebrow">接下來會發生什麼</p>
          <h2>登記之後</h2>
        </div>

        <ol class="steps">
          <li><span><b>收到你的登記</b><span>籌備小組在七個工作日內回覆，確認身份與聯絡方式。</span></span></li>
          <li><span><b>參加共好企業說明會</b><span>說明企劃書全文、菩提幣決策小組席次分配與表決規則，回答你的疑問。</span></span></li>
          <li><span><b>填回可承受的規模</b><span>中心與商家回填每月可承受的贊助規模——經濟模型的所有數字都等這組資料校準。</span></span></li>
          <li><span><b>法務結論到齊後啟動</b><span>志工定性與券的定性兩題取得律師結論、首波成員名單確認，才決定是否啟動。</span></span></li>
        </ol>

        <div class="disclose">
          <span class="tag">現況揭露</span>
          <p>菩提幣目前是 <strong>v1.0 草案</strong>，尚未經菩提幣決策小組決議；官網已經上線，菩提幣的核發、錢包與兌換要等法務結論到齊後才會開放。現在正處於企劃書中的 <strong>Phase 0</strong>：法務諮詢與治理籌組。</p>
          <p>其中兩題沒有回退路徑，必須先有答案：<strong>以服務換取住宿與給付，在勞動法上如何定性</strong>；以及<strong>可轉贈的券是否被認定為商品禮券</strong>。登記加入不代表任何承諾與義務，你隨時可以退出。</p>
          <p>菩提幣不販售、不提領、不可兌現，志工甚至沒有鏈上地址；全程沒有任何新台幣移轉。這不是投資，也不是可交易的資產。</p>
        </div>
      </div>
    </section>

    <!-- ============ 登記 ============ -->
    <section class="signup" id="join">
      <div class="shell">
        <div class="sec-head">
          <p class="eyebrow">登記加入</p>
          <h2>把你的位置留下來</h2>
          <p>填好送出後，籌備小組會在七個工作日內回覆。</p>
        </div>

        <form v-if="!sentRole" novalidate @submit.prevent="submit">
          <div class="f f-wide">
            <label for="role">我的身份 <span class="req">*</span></label>
            <select id="role" v-model="form.role" required>
              <option value="">請選擇</option>
              <option v-for="r in roles" :key="r.label">{{ r.label }}</option>
            </select>
          </div>
          <div class="f">
            <label for="name">姓名 <span class="req">*</span></label>
            <input id="name" v-model="form.name" type="text" required autocomplete="name" placeholder="王小明" />
          </div>
          <div class="f">
            <label for="org">單位名稱</label>
            <input id="org" v-model="form.org" type="text" autocomplete="organization" placeholder="中心、道場或商號（個人可留空）" />
          </div>
          <div class="f">
            <label for="email">電子郵件 <span class="req">*</span></label>
            <input id="email" v-model="form.email" type="email" required autocomplete="email" placeholder="name@example.com" />
          </div>
          <div class="f">
            <label for="phone">聯絡電話</label>
            <input id="phone" v-model="form.phone" type="tel" autocomplete="tel" placeholder="09xx-xxx-xxx" />
          </div>
          <div class="f">
            <label for="area">所在地區</label>
            <input id="area" v-model="form.area" type="text" placeholder="例：台北・南投・馬來西亞" />
          </div>
          <div class="f">
            <label for="scale">每月可承受的規模</label>
            <input id="scale" v-model="form.scale" type="text" placeholder="中心／商家填：房數、餐飲或課程名額" />
          </div>
          <div class="f f-wide">
            <label for="msg">想說的話</label>
            <textarea id="msg" v-model="form.msg" placeholder="你目前在哪裡服務、能提供什麼、或對這個企劃的疑慮"></textarea>
          </div>
          <input v-model="form.website" class="hp" tabindex="-1" autocomplete="off" aria-hidden="true" />
          <div class="form-foot">
            <button class="btn" type="submit" :disabled="sending">{{ sending ? '送出中…' : '送出登記' }}</button>
            <span class="fine">登記不代表任何承諾。資料僅供籌備小組聯繫使用。</span>
          </div>
          <p v-if="err" class="err">{{ err }}</p>
        </form>

        <div v-else ref="receipt" class="receipt">
          <h3>已收到你的登記</h3>
          <p>籌備小組會在七個工作日內回覆。</p>
          <p v-if="sentRole === '志工' || sentRole === '禪修教練'">志工與禪修教練也可以先在官網建立志工帳號，之後的服務時數與菩提幣都會記在這個帳號。</p>
          <div class="r-row">
            <router-link v-if="sentRole === '志工' || sentRole === '禪修教練'" class="btn" to="/join">建立志工帳號</router-link>
            <button class="btn btn-ghost" type="button" @click="again">再登記一筆</button>
          </div>
        </div>
      </div>
    </section>
  </q-page>
</template>

<style scoped>
.home-guide {
  position: relative;
  height: calc(100svh - 50px);
  min-height: 560px;   /* 和畫境的最小高度一致，矮螢幕上才不會蓋到下面的招募內容 */
}
/* 「一即一切／一切即一」放在 Sunny 畫境左側的天空，深咖啡字配上淡淡的光暈，壓在亮處也讀得清楚 */
.home-title {
  position: absolute;
  z-index: 1;
  left: clamp(22px, 6vw, 104px);
  top: 14%;
  margin: 0;
  font-family: var(--wenkai);
  font-weight: 700;
  font-size: clamp(40px, 5vw, 72px);
  line-height: 1.3;
  letter-spacing: 0.06em;
  color: var(--ink);
  text-shadow: 0 0 18px rgba(251, 249, 243, 0.85), 0 0 4px rgba(251, 249, 243, 0.6);
  pointer-events: none;
}
/* 深咖啡 × 淺綠（brand/coffee-green）。變數名沿用招募頁：--gold 是綠色強調色。 */
.recruit {
  --ground: #f6f2e8;
  --ground-raised: #fbf9f3;
  --ground-sunk: #ece5d6;
  --ink: #3b2a20;
  --ink-soft: #5e4c40;
  --ink-faint: #6b5c4f;
  --rule: #ddd3c2;
  --rule-strong: #bfb29c;
  --gold: #426631;
  --crimson: #8e3546;
  --crimson-wash: #f4e3e2;
  --band: #3b2a20;
  --band-2: #4d392b;
  --band-ink: #f6f2e8;
  --band-soft: #cdbfaf;
  --band-gold: #b8d8a0;
  --band-rule: #5a4536;
  --mono: 'Spline Sans Mono', ui-monospace, SFMono-Regular, Menlo, monospace;

  background: var(--ground);
  color: var(--ink);
  font-family: var(--sans);
  font-size: 16.5px;
  line-height: 1.9;
}
.shell {
  max-width: 1040px;
  margin: 0 auto;
  padding: 0 30px;
}
.eyebrow {
  font-family: var(--mono);
  font-size: 11.5px;
  font-weight: 500;
  letter-spacing: 0.2em;
  text-transform: uppercase;
  color: var(--gold);
  margin: 0;
}
h2 {
  font-family: var(--wenkai);
  font-weight: 700;
  font-size: clamp(27px, 3.6vw, 38px);
  line-height: 1.4;
  letter-spacing: 0.02em;
  margin: 12px 0 0;
  text-wrap: balance;
}
h3 {
  font-family: var(--wenkai);
  font-weight: 700;
  font-size: 20px;
  line-height: 1.5;
  margin: 0;
  letter-spacing: normal;
}
p {
  margin: 0;
  max-width: 64ch;
}
strong {
  font-weight: 500;
  color: var(--ink);
}
section,
header {
  scroll-margin-top: 56px;
}

/* ================= 開頭：一即一切 ================= */
.hero {
  background: radial-gradient(120% 90% at 50% 0%, var(--band-2) 0%, var(--band) 62%);
  color: var(--band-ink);
  border-bottom: 1px solid var(--band-rule);
  overflow: hidden;
}
.hero .shell {
  padding-top: 82px;
  padding-bottom: 78px;
  display: grid;
  grid-template-columns: minmax(0, 1.05fr) minmax(0, 0.95fr);
  gap: 0 56px;
  align-items: center;
}
.hero-copy {
  min-width: 0;
}
.hero .eyebrow {
  color: var(--band-gold);
}
.mark {
  display: block;
  height: 52px;
  width: auto;
  fill: var(--band-gold);
  margin: 0 0 26px;
}
.hairline {
  height: 1px;
  background: var(--band-rule);
  margin: 28px 0;
  max-width: 420px;
}
.hero-lede {
  color: var(--band-soft);
  font-size: 17px;
  line-height: 1.95;
  max-width: 46ch;
  font-weight: 300;
}
.hero-lede b {
  color: var(--band-ink);
  font-weight: 500;
}
.cta-row {
  display: flex;
  flex-wrap: wrap;
  gap: 14px;
  margin-top: 34px;
  align-items: center;
}
.btn {
  font-family: var(--sans);
  font-size: 15px;
  font-weight: 500;
  padding: 13px 30px;
  border: 1px solid var(--band-gold);
  background: var(--band-gold);
  color: #3b2a20;
  text-decoration: none;
  cursor: pointer;
  border-radius: 1px;
  letter-spacing: 0.06em;
  line-height: 1.6;
  transition: background 0.18s, color 0.18s;
}
.btn:hover {
  background: #cbe5b6;
}
.btn:disabled {
  opacity: 0.6;
  cursor: default;
}
.btn-ghost {
  background: transparent;
  color: var(--band-gold);
}
.btn-ghost:hover {
  background: rgba(184, 216, 160, 0.14);
}
.hero-note {
  font-family: var(--mono);
  font-size: 11.5px;
  letter-spacing: 0.06em;
  color: var(--band-soft);
  margin-top: 18px;
}

/* 環形循環圖 */
.ring {
  margin: 0;
}
.ring svg {
  display: block;
  width: 100%;
  height: auto;
}
.r-node {
  fill: var(--band-2);
  stroke: var(--band-gold);
  stroke-width: 1;
}
.r-name {
  fill: var(--band-ink);
  font-family: var(--sans);
  font-size: 17px;
  font-weight: 500;
}
.r-sub {
  fill: var(--band-soft);
  font-family: var(--mono);
  font-size: 11px;
}
.r-arc {
  stroke: var(--band-rule);
  stroke-width: 1.25;
  fill: none;
}
.r-act {
  fill: var(--band-gold);
  font-family: var(--sans);
  font-size: 14.5px;
  font-weight: 500;
}
.r-actsub {
  fill: var(--band-soft);
  font-family: var(--mono);
  font-size: 10.5px;
}
.r-mid {
  fill: var(--band-soft);
  font-family: var(--sans);
  font-size: 13px;
}
.ring-list {
  display: none;
}

/* ================= 一般段落 ================= */
section.band {
  border-bottom: 1px solid var(--rule);
}
section.band > .shell {
  padding-top: 76px;
  padding-bottom: 76px;
}
.sec-head {
  max-width: 66ch;
}
.sec-head p:not(.eyebrow) {
  margin-top: 16px;
  color: var(--ink-soft);
}

/* 唯識三轉 */
.triad {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 0 34px;
  margin-top: 48px;
}
.triad > div {
  padding-top: 20px;
  border-top: 2px solid var(--gold);
}
.triad .glyph {
  font-family: var(--wenkai);
  font-weight: 700;
  font-size: 40px;
  line-height: 1.1;
  letter-spacing: 0.06em;
  color: var(--ink);
}
.triad .romaji {
  font-family: var(--mono);
  font-size: 11px;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: var(--ink-faint);
  margin-top: 8px;
}
.triad p {
  margin-top: 12px;
  font-size: 15.5px;
  color: var(--ink-soft);
  line-height: 1.85;
}
.origin {
  margin-top: 52px;
  padding: 24px 0 0;
  border-top: 1px solid var(--rule);
  display: grid;
  grid-template-columns: 150px minmax(0, 1fr);
  gap: 0 28px;
}
.origin .k {
  font-family: var(--mono);
  font-size: 11.5px;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: var(--ink-faint);
  padding-top: 7px;
}
.origin p {
  font-size: 15.5px;
  color: var(--ink-soft);
  max-width: 62ch;
}

/* 為什麼 —— 三個缺口 */
.gaps {
  display: grid;
  gap: 0;
  margin-top: 44px;
  border-top: 1px solid var(--rule);
}
.gaps > div {
  display: grid;
  grid-template-columns: minmax(0, 240px) minmax(0, 1fr);
  gap: 6px 34px;
  padding: 26px 0;
  border-bottom: 1px solid var(--rule);
}
.gaps h3 {
  font-size: 19px;
}
.gaps p {
  font-size: 15.5px;
  color: var(--ink-soft);
  max-width: 60ch;
}
.pull {
  margin-top: 38px;
  font-family: var(--wenkai);
  font-weight: 700;
  font-size: clamp(20px, 2.6vw, 26px);
  line-height: 1.65;
  color: var(--ink);
  max-width: 34ch;
  border-left: 2px solid var(--gold);
  padding-left: 22px;
}

/* 梯級表 */
.tw {
  overflow-x: auto;
  margin-top: 40px;
}
table {
  width: 100%;
  border-collapse: collapse;
  font-size: 15.5px;
}
caption {
  text-align: left;
  font-family: var(--mono);
  font-size: 11.5px;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  color: var(--ink-faint);
  padding-bottom: 12px;
}
th {
  text-align: left;
  font-family: var(--mono);
  font-weight: 500;
  font-size: 11.5px;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  color: var(--ink-faint);
  padding: 0 18px 10px 0;
  border-bottom: 1px solid var(--rule-strong);
  white-space: nowrap;
  vertical-align: bottom;
}
td {
  padding: 13px 18px 13px 0;
  border-bottom: 1px solid var(--rule);
  vertical-align: top;
  color: var(--ink-soft);
}
td:first-child {
  color: var(--ink);
}
th:last-child,
td:last-child {
  padding-right: 0;
}
.num {
  font-family: var(--mono);
  font-variant-numeric: tabular-nums;
  text-align: right;
  white-space: nowrap;
  color: var(--ink);
}
th.num {
  text-align: right;
}
.after-table {
  margin-top: 22px;
  font-size: 15.5px;
  color: var(--ink-soft);
}

/* 券 */
.vouchers {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0;
  margin-top: 44px;
  border-top: 2px solid var(--ink);
  border-bottom: 1px solid var(--rule-strong);
}
.vouchers > div {
  padding: 24px 26px 26px 0;
}
.vouchers > div + div {
  padding-left: 30px;
  border-left: 1px solid var(--rule);
}
.vouchers h4 {
  font-family: var(--mono);
  font-size: 11.5px;
  font-weight: 500;
  line-height: 1.6;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  color: var(--gold);
  margin: 0 0 14px;
}
.vouchers p {
  font-size: 15.5px;
  color: var(--ink-soft);
  margin-bottom: 12px;
}
.vouchers ul {
  margin: 0;
  padding-left: 1.1em;
  font-size: 15px;
  color: var(--ink-soft);
}
.vouchers li {
  margin-bottom: 6px;
}

/* 誰可以加入 */
.roles {
  display: grid;
  gap: 0;
  margin-top: 44px;
}
.roles > div {
  padding: 28px 0 28px 26px;
  border-left: 2px solid var(--rule);
  border-bottom: 1px solid var(--rule);
}
.roles > div:first-child {
  border-top: 1px solid var(--rule);
}
.roles > div:hover {
  border-left-color: var(--gold);
}
.roles .rk {
  font-family: var(--mono);
  font-size: 11px;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: var(--ink-faint);
  margin-bottom: 8px;
}
.roles h3 {
  margin-bottom: 10px;
}
.roles p {
  font-size: 15.5px;
  color: var(--ink-soft);
}
.roles .gives {
  margin-top: 14px;
  font-size: 14.5px;
  color: var(--ink-soft);
  display: flex;
  flex-wrap: wrap;
  gap: 6px 10px;
}
.roles .gives span {
  font-family: var(--mono);
  font-size: 11.5px;
  letter-spacing: 0.05em;
  padding: 3px 10px;
  background: var(--ground-sunk);
  color: var(--ink-soft);
}

/* 流程 */
.steps {
  list-style: none;
  margin: 44px 0 0;
  padding: 0;
  counter-reset: s;
}
.steps li {
  counter-increment: s;
  display: grid;
  grid-template-columns: 64px minmax(0, 1fr);
  gap: 0 26px;
  padding: 20px 0;
  border-bottom: 1px solid var(--rule);
}
.steps li:first-child {
  border-top: 1px solid var(--rule);
}
.steps li::before {
  content: '0' counter(s);
  font-family: var(--mono);
  font-size: 13px;
  font-weight: 500;
  color: var(--gold);
  padding-top: 5px;
  letter-spacing: 0.06em;
}
.steps b {
  display: block;
  font-weight: 500;
  color: var(--ink);
  font-family: var(--wenkai);
  font-size: 18px;
}
.steps span {
  font-size: 15.5px;
  color: var(--ink-soft);
}

/* 現況揭露 */
.disclose {
  margin-top: 40px;
  background: var(--crimson-wash);
  border-left: 3px solid var(--crimson);
  padding: 24px 28px;
}
.disclose .tag {
  font-family: var(--mono);
  font-size: 11px;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: var(--crimson);
  display: block;
  margin-bottom: 12px;
}
.disclose p {
  font-size: 15.5px;
  color: var(--ink-soft);
  max-width: 66ch;
}
.disclose p + p {
  margin-top: 10px;
}

/* ================= 表單 ================= */
.signup {
  background: var(--band);
  color: var(--band-ink);
  border-bottom: 1px solid var(--band-rule);
}
.signup > .shell {
  padding-top: 78px;
  padding-bottom: 82px;
}
.signup .eyebrow {
  color: var(--band-gold);
}
.signup h2 {
  color: var(--band-ink);
}
.signup .sec-head p:not(.eyebrow) {
  color: var(--band-soft);
}
form {
  margin-top: 44px;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 26px 30px;
}
.f {
  display: flex;
  flex-direction: column;
  gap: 7px;
  min-width: 0;
}
.f-wide {
  grid-column: 1/-1;
}
label {
  font-family: var(--mono);
  font-size: 11px;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: var(--band-soft);
}
label .req {
  color: var(--band-gold);
}
input,
select,
textarea {
  font-family: var(--sans);
  font-size: 16px;
  color: var(--band-ink);
  background: transparent;
  border: none;
  border-bottom: 1px solid var(--band-rule);
  padding: 9px 2px;
  border-radius: 0;
  width: 100%;
}
select {
  appearance: none;
  background-image: linear-gradient(45deg, transparent 50%, var(--band-gold) 50%),
    linear-gradient(135deg, var(--band-gold) 50%, transparent 50%);
  background-position: calc(100% - 13px) 19px, calc(100% - 8px) 19px;
  background-size: 5px 5px, 5px 5px;
  background-repeat: no-repeat;
}
select option {
  color: #3b2a20;
  background: #f6f2e8;
}
textarea {
  resize: vertical;
  min-height: 92px;
  line-height: 1.8;
}
input::placeholder,
textarea::placeholder {
  color: #a89a8c;
}
input:focus,
select:focus,
textarea:focus {
  outline: none;
  border-bottom-color: var(--band-gold);
  box-shadow: 0 1px 0 0 var(--band-gold);
}
.hp {
  position: absolute;
  left: -9999px;
  width: 1px;
  height: 1px;
}
.form-foot {
  grid-column: 1/-1;
  display: flex;
  flex-wrap: wrap;
  gap: 18px;
  align-items: center;
  margin-top: 6px;
}
.form-foot .fine {
  font-family: var(--mono);
  font-size: 11px;
  letter-spacing: 0.05em;
  color: var(--band-soft);
  max-width: 42ch;
  line-height: 1.8;
}
.err {
  font-family: var(--mono);
  font-size: 11.5px;
  color: #efa0a0;
  grid-column: 1/-1;
}
.receipt {
  margin-top: 44px;
  border: 1px solid var(--band-gold);
  padding: 26px 28px;
  background: rgba(184, 216, 160, 0.07);
  scroll-margin-top: 72px;
}
.receipt h3 {
  color: var(--band-gold);
  margin-bottom: 12px;
}
.receipt p {
  color: var(--band-soft);
  font-size: 15px;
  max-width: 58ch;
}
.receipt p + p {
  margin-top: 8px;
}
.receipt .r-row {
  display: flex;
  flex-wrap: wrap;
  gap: 14px;
  align-items: center;
  margin-top: 18px;
}

.recruit :focus-visible {
  outline: 2px solid var(--gold);
  outline-offset: 3px;
}

@media (max-width: 900px) {
  .hero .shell {
    grid-template-columns: 1fr;
    gap: 48px 0;
  }
  .ring {
    order: 2;
  }
  .triad {
    grid-template-columns: 1fr;
    gap: 28px 0;
  }
  .triad > div {
    padding-top: 16px;
  }
}
@media (max-width: 680px) {
  .recruit {
    font-size: 16px;
  }
  .home-title {
    top: 22px;
    font-size: 32px;
  }
  .shell {
    padding: 0 22px;
  }
  section.band > .shell,
  .signup > .shell {
    padding-top: 56px;
    padding-bottom: 56px;
  }
  .hero .shell {
    padding-top: 56px;
    padding-bottom: 58px;
  }
  .ring {
    display: none;
  }
  .ring-list {
    display: grid;
    gap: 0;
    border-top: 1px solid var(--band-rule);
    order: 2;
  }
  .ring-list > div {
    display: grid;
    grid-template-columns: 74px minmax(0, 1fr);
    gap: 0 16px;
    padding: 15px 0;
    border-bottom: 1px solid var(--band-rule);
  }
  .ring-list .rs {
    font-family: var(--mono);
    font-size: 11px;
    letter-spacing: 0.1em;
    color: var(--band-gold);
    padding-top: 5px;
  }
  .ring-list .rn {
    color: var(--band-ink);
    font-size: 15.5px;
  }
  .ring-list .rd {
    color: var(--band-soft);
    font-size: 13.5px;
    font-family: var(--mono);
  }
  .origin {
    grid-template-columns: 1fr;
    gap: 8px 0;
  }
  .gaps > div {
    grid-template-columns: 1fr;
    gap: 8px 0;
  }
  .vouchers {
    grid-template-columns: 1fr;
  }
  .vouchers > div {
    padding: 22px 0;
  }
  .vouchers > div + div {
    padding-left: 0;
    border-left: none;
    border-top: 1px solid var(--rule);
  }
  form {
    grid-template-columns: 1fr;
  }
  .steps li {
    grid-template-columns: 1fr;
    gap: 4px 0;
  }
  .steps li::before {
    padding-top: 0;
  }
}
@media (prefers-reduced-motion: reduce) {
  .recruit * {
    animation: none !important;
    transition: none !important;
  }
}
</style>
