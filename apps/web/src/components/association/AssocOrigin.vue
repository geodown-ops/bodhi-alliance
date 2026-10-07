<script setup lang="ts">
// 第二章「協會起源」：世代變遷、全人健康、社會責任（簡報第 6–20 頁），圖表依簡報數據重畫。
import BarChart from './BarChart.vue'
import GenerationChart from './GenerationChart.vue'
import HexWholeHealth from './HexWholeHealth.vue'
import PeopleIcons from './PeopleIcons.vue'
import RingStat from './RingStat.vue'
import SdgWheel from './SdgWheel.vue'
import { aims } from './aims'

const levels = [
  { name: '嚴重症狀', note: '積極醫療', range: '11–15%' },
  { name: '重感症狀', note: '需接受建議', range: '32–38%' },
  { name: '中度症狀', note: '積極找尋治療', range: '19–26%' },
  { name: '輕度症狀', note: '以前有症狀', range: '18–22%' },
  { name: '輕微症狀', note: '以前無症狀', range: '6–8%' },
]
const certMajor = ['健康管理師（中國、國際、進階…）', '保健營養規畫師', '長期照顧健康管理師', '健康餐飲調配師', '體重管理師']
const certMinor = ['抗衰老健康管理師', '醫學美容諮詢師', '美容醫學美容師']

// 每位目標學員：[前文, 重點, 後文]
const audiences = [
  ['希望改善', '緊張、焦慮、失眠、長期疲勞或長期疼痛', '等困擾者'],
  ['處於', '高壓競爭環境', '的上班族、專業經理人、主管或經營者'],
  ['', '醫護、心理、社工、教育、人資', '等相關領域工作者及系所學生'],
  ['需要在工作挑戰與家庭責任中找到平衡點者', '', ''],
  ['希望增進', '情緒管理能力', '與改善人際關係者'],
  ['希望成長自我並提昇生活品質與幸福感者', '', ''],
]

const csrDims = ['公司治理', '員工、供應鏈、顧客', '環境績效', '社會關懷']
const csrAll = [66.1, 85.5, 67.4, 76.5]
const industries = [
  { name: '電信業', score: 85.8, dims: [72.6, 89.8, 89.7, 91.0], color: '#a02a27' },
  { name: '金融保險業', score: 73.5, dims: [67.3, 86.9, 63.1, 81.3], color: '#a02a27' },
  { name: '電子科技業', score: 73.1, dims: [65.7, 91.7, 67.6, 71.3], color: '#a02a27' },
  { name: '傳統產業', score: 70.5, dims: [65.2, 79.8, 67.7, 71.3], color: '#a02a27' },
  { name: '服務業', score: 70.4, dims: [62.4, 84.8, 63.2, 75.2], color: '#a02a27' },
]

const issues = [
  { label: '法令遵循', value: 64.9 },
  { label: '經營績效', value: 62.2 },
  { label: '內稽內控與風險品質管理（涵蓋產品品質追蹤、危機處理）', value: 59.5 },
  { label: '能源管理', value: 56.8 },
  { label: '公司治理與誠信經營', value: 51.4 },
  { label: '職場安全衛生', value: 51.4 },
  { label: '員工權益與人權（含友善職場、溝通管道）', value: 51.4 },
  { label: '永續發展策略', value: 48.6 },
  { label: '人才管理與培育（教育訓練）', value: 45.9 },
  { label: '廢物、污染管理', value: 45.9 },
]
const sdgs = [
  { label: '教育品質', value: 86.5 },
  { label: '就業與經濟成長', value: 83.8 },
  { label: '責任消費與生產', value: 83.8 },
  { label: '氣候行動', value: 83.8 },
  { label: '可負擔能源', value: 78.4 },
  { label: '健康與福祉', value: 75.7 },
  { label: '工作、創新和基礎設施', value: 73.0 },
  { label: '性別平等', value: 70.3 },
  { label: '永續城市', value: 70.3 },
  { label: '和平與正義制度', value: 67.6 },
]

// 佛法教育對企業管理的啟發：[前文, 重點, 中間, 重點, 後文]
const dharma = [
  ['一個以', '人為中心', '的管理模式', '', ''],
  ['一種', '圓融', '的辯證法', '思想', ''],
  ['', '諸行無常', '－永恆發展變化規律', '', ''],
  ['諸法無我的', '大智慧', '和', '大境界', ''],
  ['對', '人性細緻剖析', '能為管理者自我提升', '', ''],
  ['對', '僧團的管理', '揭示了企業管理的關鍵', '', ''],
]
</script>

<template>
  <section id="origin" class="chapter">
    <p class="kicker">協會起源</p>
    <h2>世代變遷・全人健康・社會責任</h2>
  </section>

  <!-- 世代變遷 -->
  <h2>世代變遷</h2>
  <h3>Z 世代</h3>
  <div class="figure dark">
    <GenerationChart />
    <p class="takeaway">
      <b>Z 世代</b>受到科技產物，影響很大；可說是<b>第一個</b>自小同時生活在<b>電子虛擬與現實世界的原生世代</b>。
    </p>
    <p class="source">資料來源：CBNData 發布的《2020 Z世代消費態度洞察報告》</p>
  </div>

  <h3>信仰衰退</h3>
  <div class="figure dark">
    <div class="stats3">
      <div class="stat">
        <PeopleIcons :total="3" :filled="1" color="#a02a27" />
        <p class="big teal">三分之一</p>
        <p>Z 世代感覺大部分的時間都處於<b class="gold">壓力</b>和<b class="gold">孤獨</b>狀況</p>
      </div>
      <div class="stat">
        <PeopleIcons :total="4" :filled="3" color="#a02a27" />
        <p class="big teal">四分之三</p>
        <p>Z 世代屬於<b class="gold">無宗教信仰</b>傾向</p>
        <p class="small">但只有 <b class="coral">16%</b> 面對無助時，會尋求<b class="gold">信仰</b>的<b class="gold">求助</b></p>
      </div>
      <div class="stat">
        <RingStat :value="52" color="#a02a27" />
        <p>Z 世代幾乎<b class="gold">不信任</b>宗教組織</p>
      </div>
    </div>
    <p class="takeaway">
      <b>Z 世代</b>無宗教信仰也不信任宗教組織，但是<b>大部分的時間</b>都處於<b>壓力和孤獨的自處狀況</b>。
    </p>
    <p class="source">資料來源：Springtide Research Institute 發布的《2020年宗教與年輕人現況》</p>
  </div>

  <blockquote class="master">
    <p>
      佛教衰退的問題在於<span class="hl-u">教育</span>和<span class="hl-u">人才</span>上，<br />
      而未必是<span class="hl-u">思想</span>的問題。
    </p>
    <cite>— 聖嚴法師</cite>
  </blockquote>

  <div class="aim-link">
    <span class="tag">對應宗旨一</span>
    <p>
      <template v-for="(part, j) in aims[0].parts" :key="j"><b v-if="j % 2">{{ part }}</b><template v-else>{{ part }}</template></template>
    </p>
  </div>

  <!-- 全人健康 -->
  <h2>全人健康</h2>
  <h3>壓力</h3>
  <div class="figure warm">
    <div class="stats3">
      <div class="stat">
        <RingStat :value="68" color="#a02a27" />
        <p><b class="red">Z 世代</b>認為<b>壓力</b>是他們<b>幸福的障礙</b>，相對於 2021 年 65% 比例<b class="red">成長了 3%</b></p>
      </div>
      <div class="stat">
        <RingStat :value="48" color="#a02a27" />
        <p><b class="red">Z 世代</b>認為他們自己是<b>無法</b>或<b>沒有能力有效管理</b>和<b>解決</b>他們自己的<b>壓力</b></p>
      </div>
      <div class="stat">
        <RingStat :value="53" color="#a02a27" />
        <p><b class="red">Z 世代</b>在<b>疫情</b>大流行期間，<b>無法</b>有效<b>處理壓力</b>感到更加<b>孤獨</b>、<b>寂寞</b>與<b>空虛</b></p>
      </div>
    </div>
    <p class="takeaway">
      研究顯示在所有世代中，<b>Z 世代</b>在<b>心理及社會幸福感</b>方面的<b>指數最低</b>，自我感受<b>壓力最大</b>。
    </p>
    <p class="source">資料來源：lululemon 發布的《2022年全球幸福感報告》</p>
  </div>

  <h3>更大壓力</h3>
  <div class="figure stress">
    <div>
      <p class="stress-lead"><b class="sage">Z 世代</b>比<b>老一代</b>承受著<b class="red">更大</b>的<b class="sky">壓力</b></p>
      <p class="huge">91<span>%</span></p>
      <p>Z 世代過去一個月因<b class="sky">壓力</b>經歷<b>至少一種</b>身體或情緒<b class="sky">症狀</b></p>
      <p class="source">資料來源：SocialBee 發布的《Z世代壓力調查》</p>
    </div>
  </div>

  <h3>身體健康</h3>
  <div class="figure market">
    <div class="services">
      <p class="svc">醫療服務<small>具醫療專業</small></p>
      <p class="vs">VS</p>
      <p class="svc">業務服務<small>醫療服務協助</small></p>
    </div>
    <div class="spectrum">
      <span class="tag-market all">全健康市場</span>
      <div class="levels">
        <div class="level">
          <b>{{ levels[0].name }}</b><small>（{{ levels[0].note }} <em>{{ levels[0].range }}</em>）</small>
        </div>
        <div class="sub">
          <span class="tag-market sub-tag">亞健康市場</span>
          <div v-for="l in levels.slice(1, 4)" :key="l.name" class="level">
            <b>{{ l.name }}</b><small>（{{ l.note }} <em>{{ l.range }}</em>）</small>
          </div>
        </div>
        <div class="level">
          <b>{{ levels[4].name }}</b><small>（{{ levels[4].note }} <em>{{ levels[4].range }}</em>）</small>
        </div>
      </div>
    </div>
    <div class="certs">
      <div class="cert">
        <p class="pct">85%<span class="down">▼</span></p>
        <p class="maps">對應嚴重、重感症狀</p>
        <ul>
          <li v-for="c in certMajor" :key="c">{{ c }}</li>
          <li>…</li>
        </ul>
      </div>
      <div class="cert">
        <p class="pct">15%<span class="up">▲</span></p>
        <p class="maps">對應輕度症狀</p>
        <ul>
          <li v-for="c in certMinor" :key="c">{{ c }}</li>
          <li>…</li>
        </ul>
      </div>
    </div>
  </div>
  <p class="takeaway navy">目前市場上絕大部分的健康相關證照，還是著重在身體健康的部分，身心靈的部分是趨勢，比較以減壓療癒為主。</p>

  <h3>身心靈健康</h3>
  <div class="figure">
    <div class="whole">
      <HexWholeHealth />
      <div>
        <p class="q-mb-sm">協會特別適合這些人：</p>
        <ol class="audience">
          <li v-for="(a, i) in audiences" :key="i">{{ a[0] }}<span v-if="a[1]" class="hl">{{ a[1] }}</span>{{ a[2] }}</li>
        </ol>
      </div>
    </div>
    <div class="focus">
      <p><small>市面上的健康證照（上圖 85%、15% 兩類）</small>專注在<span class="hl-u">身體亞健康</span></p>
      <p class="ours"><small>世界佛教教育協會</small>專注在<span class="hl-u">身心靈亞健康</span></p>
    </div>
  </div>

  <div class="aim-link red">
    <span class="tag">對應宗旨二</span>
    <p>
      <template v-for="(part, j) in aims[1].parts" :key="j"><b v-if="j % 2">{{ part }}</b><template v-else>{{ part }}</template></template>
    </p>
  </div>

  <!-- 社會責任 -->
  <h2>社會責任</h2>
  <h3>企業社會責任（CSR）成績</h3>
  <div class="figure csr-all">
    <div>
      <p class="csr-title">2020 年 CSR<br />全產業成績總平均 73 分</p>
      <p class="csr-score">總平均 <b>73.0</b> 分</p>
    </div>
    <div>
      <div v-for="(dim, i) in csrDims" :key="dim" class="mini">
        <span>{{ dim }}</span>
        <i><em :style="{ width: `${csrAll[i]}%` }" /></i>
        <b>{{ csrAll[i].toFixed(1) }}</b>
      </div>
    </div>
  </div>
  <div class="figure">
    <p class="csr-title">五大產業別成績，電信業居冠、金融保險業亞軍</p>
    <div class="table-wrap">
      <table class="ind-table">
        <thead>
          <tr>
            <th>產業</th>
            <th>CSR 總平均</th>
            <th v-for="dim in csrDims" :key="dim">{{ dim }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="ind in industries" :key="ind.name">
            <td><i class="dot" :style="{ background: ind.color }" />{{ ind.name }}</td>
            <td class="total" :style="{ color: ind.color }">{{ ind.score.toFixed(1) }}</td>
            <td v-for="(v, i) in ind.dims" :key="i" :class="{ best: v === Math.max(...industries.map((x) => x.dims[i])) }">{{ v.toFixed(1) }}</td>
          </tr>
        </tbody>
      </table>
    </div>
    <p class="source">每一欄最高分以粗體標示</p>
  </div>

  <div class="two">
    <div class="figure">
      <BarChart
        title="企業認為法令遵循、經營績效，是最重要兩大議題（可複選，%）"
        question="問：請以治理面、社會面、環境面，列出貴公司「前10大重大性議題」為何？"
        :items="issues"
        note="註：僅列出前10項"
      />
    </div>
    <div class="figure">
      <BarChart
        title="教育品質是國內企業呼應最高的 SDGs（可複選，%）"
        question="問：請根據貴公司永續發展藍圖，勾選出所呼應之聯合國永續發展目標（SDGs）？"
        :items="sdgs"
        note="註：僅列出前10項"
        color="#a02a27"
      />
    </div>
  </div>

  <div class="good">
    <p class="chain">品德 <span>➜</span> 善管理 <span>➜</span> 善企業 <span>➜</span> 善經濟</p>
    <p class="good-lead">企業家從自我的品德要求，進而走向能利他的善管理的善企業</p>
    <p class="good-quote">
      「品善德澤利群生，善德者常住真心；企業家從事企業必須為善，創造員工與社會的福祉，創造環境的永續；企業只要是善企業，就能很容易達成理想的市場法則。」
    </p>
  </div>
  <div class="sdg-row">
    <SdgWheel label="環境" :goals="[11, 6, 12, 13]" />
    <SdgWheel label="社會" :goals="[12, 15, 3, 13, 1, 4, 8, 5]" />
    <SdgWheel label="治理" :goals="[12, 1, 3, 8, 13, 9, 2, 7, 16]" />
  </div>

  <h3>佛法教育</h3>
  <div class="dharma">
    <div class="dharma-core">佛法<br />教育</div>
    <ul>
      <li v-for="(d, i) in dharma" :key="i">
        {{ d[0] }}<span class="hl-u">{{ d[1] }}</span>{{ d[2] }}<span v-if="d[3]" class="hl-u">{{ d[3] }}</span>{{ d[4] }}
      </li>
    </ul>
  </div>

  <div class="aim-link">
    <span class="tag">對應宗旨三</span>
    <p>
      <template v-for="(part, j) in aims[2].parts" :key="j"><b v-if="j % 2">{{ part }}</b><template v-else>{{ part }}</template></template>
    </p>
  </div>
</template>

<style scoped>
.stats3 {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 24px;
  text-align: center;
}
.stat p {
  margin: 0 0 6px;
}
.big {
  font-family: var(--wenkai);
  font-size: 1.8rem !important;
  font-weight: 700;
  line-height: 1.3 !important;
}
.small {
  font-size: 0.95rem !important;
}
.teal {
  color: #a02a27 !important;
}
.gold {
  color: #a02a27;
}
.coral {
  color: #a02a27;
  font-size: 1.2rem;
}
.red {
  color: var(--weba);
}
.sky {
  color: var(--sky);
}
.sage {
  color: var(--sage);
}
.figure.dark .takeaway {
  background: #f3e6c8;
  color: #a02a27 !important;
}
.figure.dark .takeaway b {
  color: #a02a27;
}
.figure.dark .source {
  color: #a02a27 !important;
}
.figure.warm {
  background: var(--gold-soft);
  border-color: #f3e6c8;
}
.figure.warm .takeaway {
  background: #f1d9d4;
}
.figure.warm b {
  color: var(--ink);
}
.figure.warm b.red {
  color: var(--weba);
}

.master {
  margin: 24px 0;
  padding: 28px 24px;
  border-radius: 16px;
  background: var(--ground-sunk);
  text-align: center;
}
.master p {
  margin: 0;
  font-family: var(--wenkai);
  font-size: 1.5rem !important;
  line-height: 1.9 !important;
  color: var(--ink) !important;
}
.master cite {
  display: block;
  margin-top: 10px;
  font-style: normal;
  font-weight: 700;
  color: var(--ink-soft);
}

.figure.stress {
  text-align: center;
  background: #f1d9d4;
  border-color: #f1d9d4;
}
.figure.stress p {
  color: var(--ink);
}
.stress-lead {
  font-size: 1.2rem !important;
}
.huge {
  margin: 0;
  font-size: 5.5rem !important;
  line-height: 1.1 !important;
  font-weight: 800;
  color: #a02a27 !important;
  border-bottom: 5px solid #a02a27;
  display: inline-block;
}
.huge span {
  font-size: 2rem;
}

.market {
  display: grid;
  grid-template-columns: minmax(0, 0.8fr) minmax(0, 1.4fr) minmax(0, 1.2fr);
  gap: 20px;
  align-items: center;
}
.services {
  text-align: center;
}
.svc {
  margin: 0;
  font-size: 1.2rem !important;
  font-weight: 700;
  color: #a02a27 !important;
  line-height: 1.4 !important;
}
.svc small {
  display: block;
  font-size: 0.8rem;
  font-weight: 400;
  color: var(--ink-faint);
}
.vs {
  margin: 8px 0;
  font-weight: 700;
  color: var(--ink-faint) !important;
}
.spectrum {
  position: relative;
  padding: 18px 16px 8px;
  border: 1.5px solid var(--sky);
  border-radius: 12px;
}
.tag-market {
  position: absolute;
  top: -12px;
  padding: 1px 10px;
  border-radius: 999px;
  font-size: 0.8rem;
  font-weight: 700;
  letter-spacing: 0.08em;
  color: #f3e6c8;
}
.tag-market.all {
  left: 14px;
  background: var(--sky);
}
.tag-market.sub-tag {
  right: 12px;
  background: var(--weba);
}
.sub {
  position: relative;
  margin: 12px 0 6px;
  padding-top: 6px;
  border: 1.5px solid var(--weba);
  border-radius: 10px;
  background: var(--weba-soft);
}
.level {
  padding: 8px 6px;
  text-align: center;
  line-height: 1.4;
}
.level b {
  display: block;
  font-size: 1.1rem;
}
.level small {
  color: var(--ink-faint);
}
.level em {
  font-style: normal;
  font-weight: 700;
  color: var(--sky);
}
.cert {
  padding: 8px 0;
}
.cert .pct {
  margin: 0;
  font-size: 2rem !important;
  font-weight: 800;
  color: var(--plum) !important;
  line-height: 1.2 !important;
}
.down {
  margin-left: 6px;
  font-size: 1.1rem;
  color: #a02a27;
}
.up {
  margin-left: 6px;
  font-size: 1.1rem;
  color: #a02a27;
}
.maps {
  margin: 0 0 4px;
  font-size: 0.85rem !important;
  color: var(--ink-faint) !important;
}
.cert ul {
  margin: 0;
  padding-left: 1.1em;
}
.cert li {
  font-size: 0.95rem;
  line-height: 1.6;
}
.takeaway.navy {
  background: #a02a27;
  color: #f3e6c8 !important;
  text-align: center;
}

.whole {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
  gap: 24px;
  align-items: center;
}
.audience {
  margin: 0;
  padding-left: 1.4em;
}
.audience li {
  margin: 4px 0;
}
.focus {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  margin-top: 18px;
}
.focus p {
  margin: 0;
  padding: 12px 16px;
  border-radius: 10px;
  background: var(--ground-sunk);
  font-weight: 600;
  text-align: center;
}
.focus small {
  display: block;
  font-weight: 400;
  color: var(--ink-faint);
}
.focus .ours {
  background: var(--weba-soft);
  color: var(--ink) !important;
}

.csr-all {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(0, 1.4fr);
  gap: 24px;
  align-items: center;
  background: linear-gradient(120deg, #f3e6c8, #dfe9d8) !important;
  border-color: #dfe9d8 !important;
}
.csr-title {
  margin: 0 0 10px;
  font-weight: 700;
  font-size: 1.15rem !important;
  line-height: 1.5 !important;
  color: var(--ink) !important;
}
.csr-score {
  margin: 0;
  font-size: 1.1rem !important;
}
.csr-score b {
  font-size: 2.8rem;
  color: #a02a27;
}
.mini {
  display: grid;
  grid-template-columns: minmax(0, 9em) minmax(0, 1fr) 2.8em;
  gap: 10px;
  align-items: center;
  margin: 8px 0;
}
.mini span {
  color: var(--ink-soft);
}
.mini i {
  height: 12px;
  border-radius: 6px;
  background: #f1d9d4;
  overflow: hidden;
}
.mini em {
  display: block;
  height: 100%;
  border-radius: 6px;
  background: #a02a27;
}
.mini b {
  font-variant-numeric: tabular-nums;
  text-align: right;
}
.table-wrap {
  overflow-x: auto;
}
.ind-table {
  width: 100%;
  min-width: 560px;
  border-collapse: collapse;
  font-variant-numeric: tabular-nums;
}
.ind-table th,
.ind-table td {
  padding: 10px 8px;
  border-bottom: 1px solid var(--rule);
  text-align: right;
  color: var(--ink-soft);
}
.ind-table th {
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--ink-faint);
}
.ind-table th:first-child,
.ind-table td:first-child {
  text-align: left;
  color: var(--ink);
  font-weight: 600;
}
.ind-table .total {
  font-size: 1.15rem;
  font-weight: 800;
}
.ind-table .best {
  font-weight: 800;
  color: var(--ink);
}
.dot {
  display: inline-block;
  width: 10px;
  height: 10px;
  margin-right: 8px;
  border-radius: 50%;
}

.two {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.good {
  margin: 24px 0 0;
  padding: 24px 28px;
  border-radius: 16px 16px 0 0;
  background: var(--gold-soft);
}
.chain {
  margin: 0 0 8px;
  font-size: 1.3rem !important;
  font-weight: 800;
  color: var(--weba) !important;
}
.chain span {
  margin: 0 6px;
  color: #a02a27;
}
.good-lead {
  margin: 0 0 6px;
  font-weight: 700;
  font-size: 1.15rem !important;
  color: var(--ink) !important;
}
.good-quote {
  margin: 0;
  font-family: var(--wenkai);
  font-size: 1.12rem !important;
  color: var(--ink) !important;
}
.sdg-row {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 20px;
  padding: 24px;
  border-radius: 0 0 16px 16px;
  background: var(--ground-raised);
  border: 1px solid var(--rule);
  border-top: none;
}

.dharma {
  display: grid;
  grid-template-columns: auto minmax(0, 1fr);
  gap: 0;
  align-items: center;
  margin: 12px 0;
}
.dharma-core {
  display: grid;
  place-items: center;
  width: 150px;
  height: 130px;
  border-radius: 14px;
  background: var(--ink-soft);
  color: #f3e6c8;
  font-family: var(--wenkai);
  font-size: 2rem;
  line-height: 1.25;
  text-align: center;
}
.dharma ul {
  list-style: none;
  margin: 0;
  padding: 0 0 0 28px;
  border-left: 2px solid var(--rule);
  margin-left: 24px;
}
.dharma li {
  position: relative;
  margin: 6px 0;
  font-size: 1.1rem;
}
.dharma li::before {
  content: '';
  position: absolute;
  left: -28px;
  top: 50%;
  width: 20px;
  border-top: 2px solid var(--rule);
}

@media (max-width: 760px) {
  .stats3,
  .market,
  .whole,
  .csr-all,
  .two,
  .sdg-row {
    grid-template-columns: 1fr;
  }
  .focus {
    grid-template-columns: 1fr;
  }
  .huge {
    font-size: 4.2rem !important;
  }
  .dharma {
    grid-template-columns: 1fr;
  }
  .dharma-core {
    width: 100%;
    height: auto;
    padding: 12px;
    font-size: 1.6rem;
  }
  .dharma ul {
    margin: 12px 0 0 12px;
  }
  .good {
    padding: 18px;
  }
}
</style>
