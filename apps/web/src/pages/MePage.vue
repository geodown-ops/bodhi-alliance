<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Dialog, Notify } from 'quasar'
import { api, ApiError, type Group } from '../api'
import {
  eventTime,
  logout,
  me,
  publicEvents,
  publicVenues,
  type AssociationFeed,
  type PracticeEvent,
  type Venue,
  type Volunteer,
  type Wallet,
} from '../account'

const route = useRoute()
const router = useRouter()
const venues = ref<Venue[]>([])
const profile = ref<Volunteer | null>(null)
const wallet = ref<Wallet>({ balance: 0, entries: [] })
const groups = ref<Group[]>([])
const myEvents = ref<PracticeEvent[]>([])
const openEvents = ref<PracticeEvent[]>([])
const kind = ref<'all' | 'online' | 'offline'>('all')
const feed = ref<AssociationFeed>({ issues: [], events: [], notices: [] })
const notVolunteer = ref(false)
const loading = ref(true)
const editing = ref(false)
const form = reactive({ display_name: '', legal_name: '', phone: '', line_id: '', wants_coach: false })

const toast = (e: unknown) => Notify.create({ type: 'negative', message: e instanceof ApiError ? e.message : '發生錯誤' })
const coins = (n: number) => n.toLocaleString('zh-TW')

async function load() {
  try {
    profile.value = await me.profile()
  } catch (e) {
    if (e instanceof ApiError && e.message.includes('還不是志工')) notVolunteer.value = true
    else toast(e)
    loading.value = false
    return
  }
  // 覺行小組會員才有錢包與活動；協會會員才有會刊、行事曆與通知
  const p = profile.value
  const none = Promise.reject()
  none.catch(() => {})
  const [w, mine, open, g, v, f] = await Promise.allSettled([
    p.in_groups ? me.wallet() : none,
    p.in_groups ? me.events() : none,
    p.in_groups ? publicEvents() : none,
    p.in_groups ? api.groups() : none,
    p.in_groups ? publicVenues() : none,
    p.in_association ? me.association() : none,
  ])
  if (w.status === 'fulfilled') wallet.value = w.value
  if (mine.status === 'fulfilled') myEvents.value = mine.value
  if (open.status === 'fulfilled') openEvents.value = open.value
  if (g.status === 'fulfilled') groups.value = g.value
  if (v.status === 'fulfilled') venues.value = v.value
  if (f.status === 'fulfilled') feed.value = f.value
  loading.value = false
}

// 切換會員身分：至少保留一個
async function setMembership(key: 'in_groups' | 'in_association', on: boolean) {
  const p = profile.value!
  const next = { in_groups: p.in_groups, in_association: p.in_association, [key]: on }
  if (!next.in_groups && !next.in_association) {
    Notify.create({ type: 'warning', message: '至少要保留一個會員身分' })
    return
  }
  const name = key === 'in_groups' ? '覺行小組' : '世界佛教教育協會'
  const apply = () => run(() => me.memberships(next), on ? `已加入${name}` : `已退出${name}`)
  if (on) return apply()
  Dialog.create({
    title: `退出${name}`,
    message: key === 'in_groups' ? '退出後看不到菩提幣錢包與共修活動；餘額會保留，重新加入就看得到。' : '退出後就看不到協會會刊、會員行事曆與通知。',
    cancel: true,
  }).onOk(apply)
}

const dayFmt = new Intl.DateTimeFormat('zh-TW', { year: 'numeric', month: 'numeric', day: 'numeric' })
const day = (s: string) => dayFmt.format(new Date(s))

onMounted(async () => {
  await load()
  // 從協會頁的「加入會員」過來：已經是會員就直接加入協會
  if (route.query.join === 'association' && profile.value) {
    await router.replace({ query: {} })
    if (!profile.value.in_association) await setMembership('in_association', true)
    return
  }
  // 從覺行小組頁的「在這裡發起活動」過來時，直接打開發起表單並選好場域
  const v = venues.value.find((x) => x.id === route.query.venue)
  if (v && profile.value) {
    // 先清掉網址參數再開表單；換網址會讓已開的對話框自動關掉
    await router.replace({ query: {} })
    startCreate()
    pickVenue(v.id)
  }
})

const joinedGroups = computed(() => new Set(profile.value?.groups.map((g) => g.group_id)))
const joinedEvents = computed(() => new Set(myEvents.value.map((e) => e.id)))
const available = computed(() =>
  openEvents.value.filter((e) => !joinedEvents.value.has(e.id) && (kind.value === 'all' || e.is_online === (kind.value === 'online'))),
)
const now = Date.now()
const ended = (e: PracticeEvent) => new Date(e.ends_at).getTime() <= now

const roleText = { organizer: '發起人', helper: '協辦', participant: '參加' } as const
function claimText(e: PracticeEvent) {
  if (e.status === 'cancelled') return '已取消'
  switch (e.claim_status) {
    case 'submitted':
      return '菩提幣審核中'
    case 'approved':
      return '菩提幣已核發'
    case 'rejected':
      return '審核退回'
  }
  return ended(e) ? '已結束' : '即將舉行'
}

const statusText = computed(() => {
  const p = profile.value
  if (!p) return ''
  if (p.frozen) return '已凍結。如果不是你自己申請的，請聯絡管理員。'
  return {
    pending: `等待${p.center_name || ''}管理員核對身分；不影響參加或發起活動。`,
    verified: '已核對身分。',
    rejected: '管理員退回了你的資料，請依說明修改後儲存，會重新送審。',
  }[p.status]
})

function startEdit() {
  const p = profile.value!
  Object.assign(form, { display_name: p.display_name, legal_name: p.legal_name, phone: p.phone, line_id: p.line_id, wants_coach: p.wants_coach })
  editing.value = true
}

async function save() {
  try {
    profile.value = await me.update({ ...form })
    editing.value = false
  } catch (e) {
    toast(e)
  }
}

async function run(action: () => Promise<unknown>, done?: string) {
  try {
    await action()
    if (done) Notify.create({ type: 'positive', message: done })
    await load()
  } catch (e) {
    toast(e)
  }
}

const joinEvent = (e: PracticeEvent, role: 'participant' | 'helper') =>
  run(() => me.joinEvent(e.id, role), role === 'helper' ? '已登記為協辦' : '已報名')
const joinGroup = (g: Group) => run(() => me.join(g.id))

function leaveGroup(g: { group_id: string; name: string }) {
  Dialog.create({ title: '退出小組', message: `確定退出「${g.name}」？`, cancel: true }).onOk(() => run(() => me.leave(g.group_id)))
}

function leaveEvent(e: PracticeEvent) {
  Dialog.create({ title: '取消報名', message: `確定不參加「${e.title}」？`, cancel: true }).onOk(() => run(() => me.leaveEvent(e.id)))
}

function cancelEvent(e: PracticeEvent) {
  Dialog.create({ title: '取消活動', message: `確定取消「${e.title}」？已報名的人會在個人頁看到活動取消。`, cancel: true }).onOk(() =>
    run(() => me.cancelEvent(e.id), '活動已取消'),
  )
}

// 發起共修活動
const creating = ref(false)
function pad(n: number) {
  return String(n).padStart(2, '0')
}
function localInput(d: Date) {
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`
}
const draft = reactive({ title: '', is_online: false, venue_id: null as string | null, location: '', starts_at: '', ends_at: '', capacity: 6, description: '' })
const venueOptions = computed(() => venues.value.map((v) => ({ label: `${v.name}（${v.center_name}）`, value: v.id })))
// 選了場域，地點自動帶場域地址（還可以再改）
function pickVenue(id: string | null) {
  const prev = venues.value.find((v) => v.id === draft.venue_id)
  draft.venue_id = id
  const v = venues.value.find((x) => x.id === id)
  if (v && (!draft.location || draft.location === prev?.address)) draft.location = v.address || v.name
}
function startCreate() {
  const s = new Date(Date.now() + 2 * 24 * 3600 * 1000)
  s.setHours(19, 0, 0, 0)
  Object.assign(draft, {
    title: '',
    is_online: false,
    venue_id: null,
    location: '',
    starts_at: localInput(s),
    ends_at: localInput(new Date(s.getTime() + 90 * 60 * 1000)),
    capacity: 6,
    description: '',
  })
  creating.value = true
}
async function create() {
  try {
    await me.createEvent({
      ...draft,
      venue_id: draft.is_online ? '' : (draft.venue_id ?? ''),
      capacity: Number(draft.capacity),
      starts_at: new Date(draft.starts_at).toISOString(),
      ends_at: new Date(draft.ends_at).toISOString(),
    })
    creating.value = false
    Notify.create({ type: 'positive', message: '活動已發起，覺行小組頁上看得到' })
    await load()
  } catch (e) {
    toast(e)
  }
}

// 活動結束後送審菩提幣
const claiming = ref<PracticeEvent | null>(null)
const claim = reactive({ attendance: 3, report: '' })
function startClaim(e: PracticeEvent) {
  Object.assign(claim, { attendance: Math.max(3, e.joined), report: '' })
  claiming.value = e
}
async function submitClaim() {
  const e = claiming.value!
  try {
    await me.claim(e.id, { attendance: Number(claim.attendance), report: claim.report })
    claiming.value = null
    Notify.create({ type: 'positive', message: '已送審，核准後菩提幣會進到發起人與協辦志工的錢包' })
    await load()
  } catch (err) {
    toast(err)
  }
}

async function signOut() {
  await logout()
  router.push('/')
}
</script>

<template>
  <q-page class="page">
    <div class="row items-center q-mb-md">
      <h1 class="q-mb-none">我的個人頁</h1>
      <q-space />
      <q-btn flat no-caps icon="logout" label="登出" @click="signOut" />
    </div>
    <q-spinner v-if="loading" color="secondary" size="32px" />

    <div v-else-if="notVolunteer" class="note">
      這個帳號是管理用帳號，沒有個人頁。請到<router-link to="/join">加入會員</router-link>另外建立帳號。
    </div>

    <template v-else-if="profile">
      <div :class="['top', { solo: !profile.in_groups }]">
        <div v-if="profile.in_groups" class="card wallet">
          <div class="text-caption">錢包餘額</div>
          <div class="balance">{{ coins(wallet.balance) }} <span>菩提幣</span></div>
          <p class="q-mb-none text-caption">協助共修活動、審核通過後入帳。兌換券還在籌備中。</p>
        </div>

        <div class="card">
          <div class="row items-start">
            <div class="col">
              <h2 class="q-mt-none q-mb-xs">{{ profile.display_name }}</h2>
              <div>{{ profile.legal_name }}<span v-if="profile.center_name"> · {{ profile.center_name }}</span><span v-if="profile.is_coach"> · 減壓教練</span></div>
              <div class="text-caption">
                {{ profile.email }}{{ profile.line_id ? ` · LINE ${profile.line_id}` : '' }}{{ profile.phone ? ` · ${profile.phone}` : '' }}
              </div>
            </div>
            <q-btn v-if="!editing" flat no-caps color="secondary" icon="edit" label="修改" @click="startEdit" />
          </div>
          <p :class="['q-mt-md', 'q-mb-none', profile.status === 'rejected' ? 'text-negative' : '']">
            <span class="status-chip q-mr-sm">{{ { pending: '待核對', verified: '已核對', rejected: '已退回' }[profile.status] }}</span>{{ statusText }}
          </p>
          <p v-if="profile.review_note" class="q-mt-sm q-mb-none">管理員的說明：{{ profile.review_note }}</p>

          <q-form v-if="editing" class="q-gutter-md q-mt-md" @submit.prevent="save">
            <q-input v-model="form.display_name" label="暱稱" outlined dense />
            <q-input v-model="form.legal_name" label="真實姓名" outlined dense :disable="profile.status === 'verified'" :hint="profile.status === 'verified' ? '核對後要改姓名，請聯絡管理員' : ''" />
            <q-input v-model="form.line_id" label="LINE ID" outlined dense />
            <q-input v-model="form.phone" label="手機" outlined dense />
            <q-checkbox v-model="form.wants_coach" label="我有帶領正念減壓的經驗，想擔任減壓教練" />
            <div class="row q-gutter-sm">
              <q-btn type="submit" color="secondary" unelevated no-caps label="儲存" />
              <q-btn flat no-caps label="取消" @click="editing = false" />
            </div>
          </q-form>
        </div>
      </div>

      <div class="card q-mt-lg">
        <h2 class="q-mt-none q-mb-xs">我的會員身分</h2>
        <p class="text-caption q-mb-sm">覺行小組和世界佛教教育協會共用這個會員帳號，可以只加入一邊，也可以兩邊都加入。</p>
        <div class="memberships">
          <q-toggle :model-value="profile.in_groups" color="secondary" @update:model-value="(v: boolean) => setMembership('in_groups', v)">
            <div>
              <div class="text-weight-bold">覺行小組會員</div>
              <div class="text-caption">菩提幣錢包、參加或發起共修活動</div>
            </div>
          </q-toggle>
          <q-toggle :model-value="profile.in_association" color="secondary" @update:model-value="(v: boolean) => setMembership('in_association', v)">
            <div>
              <div class="text-weight-bold">世界佛教教育協會會員</div>
              <div class="text-caption">協會會刊、會員行事曆、協會通知</div>
            </div>
          </q-toggle>
        </div>
      </div>

      <template v-if="profile.in_association">
        <h2>世界佛教教育協會</h2>
        <div class="assoc">
          <div class="card">
            <h3 class="q-mt-none">協會通知</h3>
            <p v-if="!feed.notices.length" class="q-mb-none">目前沒有通知。</p>
            <div v-for="n in feed.notices" :key="n.id" class="item">
              <div class="text-caption">{{ day(n.created_at) }}</div>
              <div class="text-weight-bold">{{ n.title }}</div>
              <p v-if="n.body" class="q-mb-none pre">{{ n.body }}</p>
            </div>
          </div>
          <div class="card">
            <h3 class="q-mt-none">會員行事曆</h3>
            <p v-if="!feed.events.length" class="q-mb-none">近期沒有排定的協會活動。</p>
            <div v-for="e in feed.events" :key="e.id" class="item">
              <div class="text-caption">{{ e.ends_at ? eventTime({ starts_at: e.starts_at, ends_at: e.ends_at }) : day(e.starts_at) }}</div>
              <div class="text-weight-bold">{{ e.title }}</div>
              <div v-if="e.location" class="text-caption">{{ e.location }}</div>
              <p v-if="e.description" class="q-mb-none pre">{{ e.description }}</p>
            </div>
          </div>
          <div class="card">
            <h3 class="q-mt-none">協會會刊</h3>
            <p v-if="!feed.issues.length" class="q-mb-none">會刊整理中。</p>
            <div v-for="i in feed.issues" :key="i.id" class="item">
              <div class="text-caption">{{ i.issued_on }}</div>
              <a v-if="i.url" :href="i.url" target="_blank" rel="noopener" class="text-weight-bold">{{ i.title }}</a>
              <div v-else class="text-weight-bold">{{ i.title }}</div>
              <p v-if="i.summary" class="q-mb-none pre">{{ i.summary }}</p>
            </div>
          </div>
        </div>
      </template>

      <template v-if="profile.in_groups">
      <div class="row items-center q-mt-lg">
        <h2 class="q-my-none">我的共修活動</h2>
        <q-space />
        <q-btn color="secondary" unelevated no-caps icon="add" label="發起共修活動" @click="startCreate" />
      </div>
      <p v-if="!myEvents.length" class="q-mt-md">還沒有參加或發起活動。從下面挑一場參加，或找兩位以上的朋友，自己發起一場。</p>
      <div v-else class="grid q-mt-md">
        <div v-for="e in myEvents" :key="e.id" class="card">
          <div class="row q-gutter-xs">
            <span class="status-chip">{{ roleText[e.my_role!] }}</span>
            <span class="status-chip">{{ claimText(e) }}</span>
          </div>
          <h3 class="q-my-sm">{{ e.title }}</h3>
          <p class="q-mb-xs">{{ eventTime(e) }}</p>
          <p class="q-mb-xs">{{ e.is_online ? '線上：' : '' }}{{ e.location }}</p>
          <p class="q-mb-xs text-caption">發起人 {{ e.organizer_name }} · 已報名 {{ e.joined }}／{{ e.capacity }} 人</p>
          <p v-if="e.claim_status === 'rejected' && e.claim_note" class="text-negative q-mb-xs">退回原因：{{ e.claim_note }}</p>
          <div class="row q-gutter-sm q-mt-xs">
            <template v-if="e.my_role === 'organizer' && e.status === 'open'">
              <q-btn
                v-if="ended(e) && (!e.claim_status || e.claim_status === 'rejected')"
                color="secondary"
                unelevated
                no-caps
                :label="e.claim_status === 'rejected' ? '修改後重新送審' : '申請菩提幣審核'"
                @click="startClaim(e)"
              />
              <q-btn v-if="!e.claim_status && !ended(e)" flat dense no-caps color="grey-8" label="取消活動" @click="cancelEvent(e)" />
            </template>
            <q-btn
              v-else-if="e.my_role !== 'organizer' && e.status === 'open' && (!e.claim_status || e.claim_status === 'rejected')"
              flat
              dense
              no-caps
              color="grey-8"
              label="取消報名"
              @click="leaveEvent(e)"
            />
          </div>
          <p v-if="e.my_role === 'organizer' && e.status === 'open' && !ended(e)" class="text-caption q-mt-sm q-mb-none">活動結束後，可以在這裡申請菩提幣審核。</p>
        </div>
      </div>

      <h2>可以參加的活動</h2>
      <q-btn-toggle
        v-model="kind"
        no-caps
        unelevated
        toggle-color="secondary"
        class="q-mb-md"
        :options="[
          { label: '全部', value: 'all' },
          { label: '線下', value: 'offline' },
          { label: '線上', value: 'online' },
        ]"
      />
      <p v-if="!available.length">目前沒有其他排定的活動。</p>
      <div v-else class="grid">
        <div v-for="e in available" :key="e.id" class="card">
          <div class="status-chip">{{ e.is_online ? '線上' : '線下' }}</div>
          <h3 class="q-my-sm">{{ e.title }}</h3>
          <p class="q-mb-xs">{{ eventTime(e) }}</p>
          <p v-if="!e.is_online" class="q-mb-xs">{{ e.venue_name ? `${e.venue_name} · ` : '' }}{{ e.location }}</p>
          <p class="q-mb-xs text-caption">發起人 {{ e.organizer_name }} · 已報名 {{ e.joined }}／{{ e.capacity }} 人</p>
          <p v-if="e.description" class="q-mb-sm">{{ e.description }}</p>
          <div v-if="e.joined < e.capacity" class="row q-gutter-sm">
            <q-btn outline color="secondary" no-caps label="報名參加" @click="joinEvent(e, 'participant')" />
            <q-btn flat color="secondary" no-caps label="我來協辦" @click="joinEvent(e, 'helper')" />
          </div>
          <span v-else class="status-chip">已額滿</span>
        </div>
      </div>

      <h2>我的覺行小組</h2>
      <p v-if="!profile.groups.length">還沒加入定期聚會的小組。</p>
      <div v-else class="grid q-mb-md">
        <div v-for="g in profile.groups" :key="g.group_id" class="card">
          <h3 class="q-my-none">{{ g.name }}</h3>
          <p v-if="g.role === 'leader'" class="q-my-xs"><span class="status-chip">組長</span></p>
          <q-btn flat dense no-caps color="grey-8" label="退出" @click="leaveGroup(g)" />
        </div>
      </div>
      <template v-if="groups.some((g) => !joinedGroups.has(g.id))">
        <h3>可以加入的小組</h3>
        <div class="grid">
          <div v-for="g in groups.filter((g) => !joinedGroups.has(g.id))" :key="g.id" class="card">
            <div class="status-chip">{{ g.is_online ? '線上' : g.region }}</div>
            <h3 class="q-my-sm">{{ g.name }}</h3>
            <p v-if="g.schedule" class="q-mb-sm">{{ g.schedule }}</p>
            <q-btn outline color="secondary" no-caps label="加入" @click="joinGroup(g)" />
          </div>
        </div>
      </template>

      <h2>菩提幣紀錄</h2>
      <p v-if="!wallet.entries.length">還沒有紀錄。協助一場共修活動，發起人送審通過後就會入帳。</p>
      <table v-else class="table">
        <tbody>
          <tr v-for="(x, i) in wallet.entries" :key="i">
            <td>{{ new Date(x.created_at).toLocaleDateString('zh-TW') }}</td>
            <td>{{ x.memo }}</td>
            <td class="num">{{ x.amount > 0 ? '+' : '' }}{{ coins(x.amount) }}</td>
          </tr>
        </tbody>
      </table>
      </template>
    </template>

    <q-dialog v-model="creating">
      <q-card style="width: 560px; max-width: 95vw">
        <q-form @submit.prevent="create">
          <q-card-section>
            <div class="text-h6">發起共修活動</div>
            <div class="text-caption">三人以上（含你自己）就可以進行一次正念減壓實作。</div>
          </q-card-section>
          <q-card-section class="q-gutter-md">
            <q-input v-model="draft.title" label="活動名稱 *" outlined dense :rules="[(v) => !!v.trim() || '請填寫活動名稱']" />
            <q-btn-toggle
              v-model="draft.is_online"
              no-caps
              unelevated
              toggle-color="secondary"
              :options="[
                { label: '線下', value: false },
                { label: '線上', value: true },
              ]"
            />
            <q-select
              v-if="!draft.is_online && venueOptions.length"
              :model-value="draft.venue_id"
              :options="venueOptions"
              emit-value
              map-options
              clearable
              label="活動場域"
              hint="選了場域，活動會列在覺行小組頁的這個場域底下；也可以不選，自己填地點"
              outlined
              dense
              @update:model-value="pickVenue"
            />
            <q-input
              v-model="draft.location"
              :label="draft.is_online ? '會議連結或集合方式 *' : '地點 *'"
              :hint="draft.is_online ? '只有報名的人看得到' : ''"
              outlined
              dense
              :rules="[(v) => !!v.trim() || '請填寫地點']"
            />
            <div class="row q-col-gutter-sm">
              <q-input v-model="draft.starts_at" type="datetime-local" label="開始 *" stack-label outlined dense class="col-12 col-sm-6" />
              <q-input v-model="draft.ends_at" type="datetime-local" label="結束 *" stack-label outlined dense class="col-12 col-sm-6" />
            </div>
            <q-input v-model.number="draft.capacity" type="number" label="開放人數（含你自己）*" outlined dense :rules="[(v) => v >= 3 || '至少三人']" />
            <q-input v-model="draft.description" type="textarea" autogrow label="說明" hint="例如：帶一張瑜伽墊、第一次參加也歡迎" outlined dense />
          </q-card-section>
          <q-card-actions align="right">
            <q-btn flat no-caps label="取消" @click="creating = false" />
            <q-btn type="submit" color="secondary" unelevated no-caps label="發起" />
          </q-card-actions>
        </q-form>
      </q-card>
    </q-dialog>

    <q-dialog :model-value="!!claiming" @update:model-value="claiming = null">
      <q-card v-if="claiming" style="width: 520px; max-width: 95vw">
        <q-form @submit.prevent="submitClaim">
          <q-card-section>
            <div class="text-h6">申請菩提幣審核</div>
            <div class="text-grey-8">{{ claiming.title }} · {{ eventTime(claiming) }}</div>
          </q-card-section>
          <q-card-section class="q-gutter-md">
            <p class="q-mb-none">核准後，你和登記協辦的志工、減壓教練會依梯級表拿到菩提幣。</p>
            <q-input v-model.number="claim.attendance" type="number" label="實際出席人數（含你自己）*" outlined dense :rules="[(v) => v >= 3 || '三人以上才能申請']" />
            <q-input v-model="claim.report" type="textarea" autogrow label="活動紀錄" hint="例如出席名單、帶了哪些練習" outlined dense />
          </q-card-section>
          <q-card-actions align="right">
            <q-btn flat no-caps label="取消" @click="claiming = null" />
            <q-btn type="submit" color="secondary" unelevated no-caps label="送審" />
          </q-card-actions>
        </q-form>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<style scoped>
.top {
  display: grid;
  gap: 16px;
  grid-template-columns: minmax(220px, 1fr) 2fr;
}
.top.solo {
  grid-template-columns: 1fr;
}
.memberships {
  display: flex;
  flex-wrap: wrap;
  gap: 8px 32px;
}
.assoc {
  display: grid;
  gap: 16px;
  grid-template-columns: repeat(3, minmax(0, 1fr));
}
.item + .item {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid var(--rule, rgba(0, 0, 0, 0.08));
}
.pre {
  white-space: pre-line;
}
@media (max-width: 700px) {
  .top,
  .assoc {
    grid-template-columns: 1fr;
  }
}
.balance {
  font-size: 2.2rem;
  font-variant-numeric: tabular-nums;
  margin: 4px 0 8px;
}
.balance span {
  font-size: 1rem;
  color: var(--ink-soft);
}
</style>
