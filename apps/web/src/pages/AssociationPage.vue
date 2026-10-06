<script setup lang="ts">
// 世界佛教教育協會介紹：內容依協會簡報《世界佛教教育協會介紹》整理，簡介與任務的文字以協會提供的《協會簡介》為準。
// 插圖放在 public/images/association/（依簡報插圖重新繪製，提示詞見 docs/association-illustrations.md）；
// 圖檔還沒放進去時不顯示該張圖。
import { reactive } from 'vue'

const missing = reactive(new Set<string>())
const img = (name: string) => `/images/association/${name}.webp`

const pillars = [
  { name: '慈悲關懷', text: '人文關懷、社會救助、護生及生態環保。' },
  { name: '智慧無礙', text: '雲端教學、佛學研究、生命教育、終身學習。' },
  { name: '善巧方便', text: '佛學資訊、音像紀錄、養生及身心健康、食品安全，以及社會型企業。' },
]

const tasks = [
  '推動慈悲關懷工作，包括人文關懷、社會救助、與護生及生態環保等。',
  '推動智慧無礙工作，包括雲端教學、佛學研究、生命教育、終身學習等。',
  '推動善巧方便工作，包括佛學資訊、音像紀錄，養生及身心健康、食品安全，以及社會型企業等。',
  '結合科技、醫藥、宗教、政府相關單位及民間社團組織等，推動終身學習及社會型企業策略聯盟。',
  '舉辦研討會、座談會、職業訓練、組織學習、關懷訪視等各項相關活動，以促進佛教普及弘化及世界永續發展之目的。',
]

const aims = [
  { key: '交流', text: '推廣並促進國內外產官學研界交流合作，讓佛教能普及弘化及永續發展。' },
  { key: '療癒', text: '透過佛教八正道傳統精神，讓現代人面臨的身心亞健康得以舒緩修復。' },
  { key: '善企業', text: '藉由傳統佛法教育，讓企業走向善管理的善企業，進而善盡社會責任。' },
]

const origins = [
  {
    title: '世代變遷',
    image: 'generations',
    text: 'Z 世代是第一個自小同時生活在電子虛擬與現實世界的原生世代。多數年輕人沒有宗教信仰、也不信任宗教組織，卻有大部分的時間處於壓力與孤獨之中。聖嚴法師指出，佛教衰退的問題在於教育和人才，而未必是思想。',
  },
  {
    title: '全人健康',
    image: 'wellbeing',
    text: '市場上的健康相關證照大多著重身體健康，身心靈的照顧才剛成為趨勢。協會專注在「身心靈亞健康」：緊張、焦慮、失眠、長期疲勞，高壓環境中的上班族與經營者，醫護、心理、社工、教育等助人工作者，以及希望在工作與家庭間找到平衡、提升生活品質的人。',
  },
  {
    title: '社會責任',
    image: 'responsibility',
    text: '「品德 → 善管理 → 善企業 → 善經濟」。企業家從自我的品德要求，走向能利他的善管理、善企業：創造員工與社會的福祉，創造環境的永續。佛法教育提供以人為中心的管理模式、圓融的思想、諸行無常與諸法無我的智慧，而僧團的管理也揭示了企業管理的關鍵。',
  },
]

const quality = [
  { title: '教學系統化', text: '教學教案、課件與習題、備課參考資料、實例影片，教師與學員各有專屬系統。' },
  { title: '數位典藏', text: '課程內容以數位方式保存，可線上觀看與下載。' },
  { title: '知識管理暨線上學習', text: '依主題整理的知識庫與線上課程，隨時可以學習。' },
  { title: '開發模組化', text: '依「分析、設計、發展、實施、評估」五個階段開發課程，每一階段都有檢核、回饋與修正。' },
  { title: '執行標準化', text: '問卷調查 → 課程規劃 → 課程製作 → 課程說明 → 正式開課，課後收集教學回饋。' },
]

const badges = ['學員與老師未來課程的身份認證', '與學員溝通活動訊息的管道', '結業或證書的數位證明', '學員與老師社群互動的應用']

const units = [
  {
    name: '指導委員會',
    items: ['提供市場與產業課程需求趨勢', '審核課程教材的適用性與適法性', '評核課程師資的教學專業性', '稽核課程上架後的發展績效', '決議課程在市場上的發展生命週期'],
  },
  {
    name: '課程發展中心',
    items: ['調研課程的市場需求', '開發課程的教材與教具', '舉辦課程說明會或體驗會', '負責整體課程招生計畫與執行', '評量上課的實際狀況與回饋'],
  },
  {
    name: '認證中心',
    items: ['制定、更新課程內容的標準與實施流程', '依指導委員會的評核結果核發師資證照', '評測後給予學員結業證明或證照'],
  },
]

const members = [
  { role: '策略夥伴', text: '以策略投資的方式，支援協會在初始階段及未來營運過程中的財務需求。' },
  { role: '資源協助者', text: '提供協會多面向的必要產業資源（例如企業培訓需求），以完成協會的任務。' },
  { role: '推廣協同者', text: '協助會員招募、課程招生與相關活動的執行，藉由會員的力量推廣協會。' },
]
</script>

<template>
  <q-page class="page">
    <h1>世界佛教教育協會</h1>
    <p class="lead">
      以慈悲關懷、智慧無礙、善巧方便為主軸，推展世界佛教雲端教學、佛學研究、社會教育、人文關懷、永續發展及社會型企業。
    </p>
    <figure v-if="!missing.has('hero')" class="art wide">
      <img :src="img('hero')" alt="" @error="missing.add('hero')" />
    </figure>

    <h2>協會簡介</h2>
    <p>
      本會為依法設立、非以營利為目的之社會團體，以慈悲關懷、智慧無礙、善巧方便為主軸，推展世界佛教雲端教學、佛學研究、社會教育、人文關懷、永續發展及社會型企業為宗旨，
      於民國 111 年 12 月 30 日送件申請成立。
    </p>
    <div class="grid">
      <div v-for="p in pillars" :key="p.name" class="card">
        <h3 class="q-mt-none">{{ p.name }}</h3>
        <p class="q-mb-none">{{ p.text }}</p>
      </div>
    </div>

    <h2>協會任務</h2>
    <p>本會之任務如下，並依相關法令規定推動及執行：</p>
    <ol>
      <li v-for="t in tasks" :key="t">{{ t }}</li>
    </ol>

    <h2>協會宗旨</h2>
    <div class="aims">
      <div v-for="a in aims" :key="a.key" class="aim">
        <span class="status-chip">{{ a.key }}</span>
        <p class="q-mb-none">{{ a.text }}</p>
      </div>
    </div>

    <h2>協會起源：世代變遷・全人健康・社會責任</h2>
    <div v-for="o in origins" :key="o.title" class="origin">
      <figure v-if="!missing.has(o.image)" class="art">
        <img :src="img(o.image)" alt="" @error="missing.add(o.image)" />
      </figure>
      <div>
        <h3 class="q-mt-none">{{ o.title }}</h3>
        <p class="q-mb-none">{{ o.text }}</p>
      </div>
    </div>

    <h2>營運特色</h2>
    <p>
      惠敏法師提出佛教教育「現代化」的方向：培育兼具佛學及資訊知能的人才，建立佛學數位典藏、知識管理與線上學習環境，
      構成「數位神經系統」，並結合人文與藝術資源發展文化創意產業，達成真、善、美的人生目標。
    </p>

    <h3>品質系統</h3>
    <div class="grid">
      <div v-for="q in quality" :key="q.title" class="card">
        <h3 class="q-mt-none">{{ q.title }}</h3>
        <p class="q-mb-none">{{ q.text }}</p>
      </div>
    </div>

    <h3>學習標章</h3>
    <p>數位學習標章是數位神經系統的一環：</p>
    <ul>
      <li v-for="b in badges" :key="b">{{ b }}</li>
    </ul>

    <h3>標準組織</h3>
    <p>
      協會由會員大會、理事會、監事會與理事長領導，下設秘書處（行政、總務、財務、秘書）與工作委員會（公關、業務拓展、募款、資訊、寺廟規劃委員會及課程發展中心），
      課程經審查與紀律把關，由指導委員會與認證暨品質中心確保品質。
    </p>
    <div class="grid">
      <div v-for="u in units" :key="u.name" class="card">
        <h3 class="q-mt-none">{{ u.name }}</h3>
        <ul class="q-mb-none">
          <li v-for="i in u.items" :key="i">{{ i }}</li>
        </ul>
      </div>
    </div>

    <h3>寺廟教育文創</h3>
    <figure v-if="!missing.has('temple')" class="art wide">
      <img :src="img('temple')" alt="" @error="missing.add('temple')" />
    </figure>
    <p>把寺廟與自然結合成可以靜心學習的空間：森林中的樹屋與禪修小屋、頌缽聲浴、靜坐與光影冥想。</p>

    <h2>不辦教育，佛教就沒有明天</h2>
    <p>
      聖嚴法師有感於當時社會對佛教的誤解，佛教的形象普遍落於民間信仰的層次，懇切呼籲推動佛教教育事業為當務之急。
      協會延續這份心願：教育 → 學習 → 終身，為未來生活做準備、提升生活的價值，活到老、學到老。
    </p>

    <h2>協會會員招募中</h2>
    <p>歡迎各產業的專業夥伴加入，以三種身份與協會一起成長：</p>
    <div class="grid">
      <div v-for="m in members" :key="m.role" class="card">
        <h3 class="q-mt-none">{{ m.role }}</h3>
        <p class="q-mb-none">{{ m.text }}</p>
      </div>
    </div>
    <figure v-if="!missing.has('closing')" class="art wide">
      <img :src="img('closing')" alt="" @error="missing.add('closing')" />
    </figure>
    <p class="closing">我們有決心與行動，讓世界佛教教育協會成長茁壯，歡迎您的加入。</p>
  </q-page>
</template>

<style scoped>
.art {
  margin: 16px 0;
}
.art img {
  display: block;
  width: 100%;
  border-radius: 14px;
  object-fit: cover;
}
.art.wide img {
  aspect-ratio: 16 / 7;
}
.aims {
  display: grid;
  gap: 10px;
}
.aim {
  display: flex;
  align-items: baseline;
  gap: 12px;
}
.aim .status-chip {
  flex: none;
}
.origin {
  display: grid;
  grid-template-columns: minmax(0, 2fr) minmax(0, 3fr);
  gap: 20px;
  align-items: center;
  margin: 18px 0;
}
.origin .art {
  margin: 0;
}
.origin .art img {
  aspect-ratio: 4 / 3;
}
.origin:not(:has(.art)) {
  grid-template-columns: 1fr;
}
.closing {
  font-weight: 600;
}
@media (max-width: 680px) {
  .origin {
    grid-template-columns: 1fr;
  }
}
</style>
