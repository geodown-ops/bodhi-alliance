<script setup lang="ts">
import { reactive, ref } from 'vue'
import { api, ApiError } from '../api'

const kinds = [
  { label: '禪修中心／道場', value: 'center' },
  { label: '贊助商家（飯店・水療・餐飲・商店）', value: 'sponsor' },
  { label: '其他（請在留言說明）', value: 'other' },
]
const form = reactive({
  kind: 'sponsor',
  org_name: '',
  contact_name: '',
  email: '',
  phone: '',
  region: '',
  offerings: '',
  monthly_scale: '',
  message: '',
  website: '',
})
const sending = ref(false)
const sent = ref(false)
const error = ref('')

async function submit() {
  sending.value = true
  error.value = ''
  try {
    await api.registerPartner({ ...form })
    sent.value = true
  } catch (e) {
    error.value = e instanceof ApiError ? e.message : '送出失敗'
  } finally {
    sending.value = false
  }
}
</script>

<template>
  <q-page class="page">
    <h1>共好企業登記及管理</h1>
    <p class="lead">淡季的空房、餐飲量能、課程名額，邊際成本極低，卻沒有管道回饋給最該得到它們的人。</p>
    <p>共好企業提供券，讓志工與禪修教練用服務換來的菩提幣兌換。核銷是純贊助，全程沒有新台幣移轉。</p>

    <div class="grid">
      <div class="card">
        <h3 class="q-mt-none">禪修中心與道場</h3>
        <p class="q-mb-none">核發菩提幣給自己的志工與教練，同時把住宿、餐飲與課程名額開放給所有共好企業使用，並推派委員進入菩提幣決策小組。</p>
      </div>
      <div class="card">
        <h3 class="q-mt-none">贊助商家</h3>
        <p class="q-mb-none">飯店、水療、餐飲、實體商店。自行建立券種與每月贊助額度，額度用罄即暫停接受新兌換，已在志工手上的券照常核銷。</p>
      </div>
    </div>

    <h2>企業後台</h2>
    <p>
      通過審核的共好企業會拿到後台帳號，可以自建券種、管理門市與店員、用手機掃碼核銷、查看對帳單。
      <span class="status-chip">籌備中，法務結論到齊後開放</span>
    </p>

    <h2>登記之後</h2>
    <ol>
      <li>籌備小組在七個工作日內回覆，確認身份與聯絡方式。</li>
      <li>參加共好企業說明會，了解企劃書全文、菩提幣決策小組席次與表決規則。</li>
      <li>回填每月可承受的贊助規模，經濟模型的數字都等這組資料校準。</li>
      <li>法務結論到齊、首波成員名單確認後，才決定是否啟動。登記不代表任何承諾，你隨時可以退出。</li>
    </ol>

    <h2 id="register">登記</h2>
    <div v-if="sent" class="note"><strong>已收到你的登記。</strong>籌備小組會在七個工作日內回覆。</div>
    <q-form v-else class="card q-gutter-md" @submit.prevent="submit">
      <q-select v-model="form.kind" :options="kinds" emit-value map-options label="身份 *" outlined />
      <q-input v-model="form.org_name" label="單位名稱 *" outlined :rules="[(v) => !!v.trim() || '請填寫單位名稱']" />
      <q-input v-model="form.contact_name" label="聯絡人 *" outlined :rules="[(v) => !!v.trim() || '請填寫聯絡人']" />
      <q-input v-model="form.email" type="email" label="電子郵件 *" outlined :rules="[(v) => /.+@.+\..+/.test(v) || '請填寫正確的電子郵件']" />
      <q-input v-model="form.phone" label="聯絡電話" outlined />
      <q-input v-model="form.region" label="所在地區" outlined />
      <q-input v-model="form.offerings" type="textarea" autogrow label="可提供的品項（例如：雙人房一晚、午餐套餐）" outlined />
      <q-input v-model="form.monthly_scale" label="每月可承受的規模（例如：每月 10 間房晚）" outlined />
      <q-input v-model="form.message" type="textarea" autogrow label="想說的話" outlined />
      <input v-model="form.website" class="hp" tabindex="-1" autocomplete="off" aria-hidden="true" />
      <p v-if="error" class="text-negative q-mb-none">{{ error }}</p>
      <q-btn type="submit" color="secondary" unelevated no-caps size="lg" :loading="sending" label="送出登記" />
    </q-form>
  </q-page>
</template>

<style scoped>
.hp {
  position: absolute;
  left: -9999px;
}
</style>
