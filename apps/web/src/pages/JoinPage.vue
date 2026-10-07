<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ApiError } from '../api'
import { me, register } from '../account'

// 加入會員：覺行小組與世界佛教教育協會共用同一份會員資料，報名時選要加入哪一邊（之後在個人頁可以切換）。
// 從小組或活動的「報名」按鈕進來時，順便加入那個小組或活動；從協會頁進來時預設加入協會。
const route = useRoute()
const router = useRouter()
const forAssociation = route.query.for === 'association'
const form = reactive({
  legal_name: '',
  email: '',
  password: '',
  display_name: '',
  line_id: '',
  website: '',
  in_groups: !forAssociation,
  in_association: forAssociation,
})
const sending = ref(false)
const error = ref('')
const picked = computed(() => form.in_groups || form.in_association)

async function submit() {
  if (!picked.value) {
    error.value = '請至少選擇加入覺行小組或世界佛教教育協會其中一個'
    return
  }
  sending.value = true
  error.value = ''
  try {
    await register({ ...form })
  } catch (e) {
    error.value = e instanceof ApiError ? e.message : '報名失敗'
    sending.value = false
    return
  }
  // 帳號已經建好；加入小組或活動失敗時，到個人頁再加一次就好
  try {
    if (form.in_groups && typeof route.query.group === 'string') await me.join(route.query.group)
    if (form.in_groups && typeof route.query.event === 'string')
      await me.joinEvent(route.query.event, route.query.role === 'helper' ? 'helper' : 'participant')
  } catch {
    /* ignore */
  }
  router.push('/me')
}
</script>

<template>
  <q-page class="page narrow">
    <h1>{{ forAssociation ? '加入世界佛教教育協會會員' : '加入會員' }}</h1>
    <p class="lead">
      覺行小組和世界佛教教育協會共用同一個會員帳號。留下資料就完成報名，之後用電子郵件和密碼登入個人頁，隨時可以切換要加入哪一邊。
    </p>
    <p>已經是會員？<router-link :to="{ path: '/login', query: { next: '/me' } }">直接登入</router-link>，在個人頁就能加入另一邊。</p>

    <q-form class="card card-form" @submit.prevent="submit">
      <div class="choices">
        <q-checkbox v-model="form.in_groups" class="choice">
          <div>
            <div class="text-weight-bold">加入覺行小組</div>
            <div class="text-caption">參加或發起正念減壓共修，擁有菩提幣錢包</div>
          </div>
        </q-checkbox>
        <q-checkbox v-model="form.in_association" class="choice">
          <div>
            <div class="text-weight-bold">加入世界佛教教育協會</div>
            <div class="text-caption">收到協會會刊、會員行事曆與協會通知</div>
          </div>
        </q-checkbox>
      </div>
      <q-input v-model="form.legal_name" label="真實姓名 *" hint="核對身分用，不會公開" outlined :rules="[(v) => !!v.trim() || '請填寫真實姓名']" />
      <q-input v-model="form.email" type="email" label="電子郵件 *" hint="登入帳號" autocomplete="username" outlined :rules="[(v) => /.+@.+\..+/.test(v) || '請填寫正確的電子郵件']" />
      <q-input v-model="form.password" type="password" label="密碼 *" autocomplete="new-password" outlined :rules="[(v) => v.length >= 10 || '至少 10 個字元']" />
      <q-input v-model="form.display_name" label="暱稱" hint="活動頁上顯示的名字；不填就用真實姓名" outlined />
      <q-input v-model="form.line_id" label="LINE ID" hint="方便我們用 LINE 聯絡你" outlined />
      <input v-model="form.website" class="hp" tabindex="-1" autocomplete="off" aria-hidden="true" />
      <p v-if="error" class="text-negative q-mb-none">{{ error }}</p>
      <q-btn type="submit" color="secondary" unelevated no-caps size="lg" :loading="sending" :disable="!picked" label="送出報名" />
      <p class="text-caption q-mb-none">真實姓名、電子郵件與 LINE ID 只有管理員看得到。</p>
    </q-form>
  </q-page>
</template>

<style scoped>
.narrow {
  max-width: 640px;
}
.choices {
  display: grid;
  gap: 8px;
}
.choice {
  align-items: flex-start;
}
.hp {
  position: absolute;
  left: -9999px;
}
</style>
