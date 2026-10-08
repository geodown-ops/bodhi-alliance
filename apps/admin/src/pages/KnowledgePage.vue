<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { Notify } from 'quasar'
import { api } from '../session'

type Doc = { id: string; title: string; category: string; status: string; chars: number; uploaded_by: string | null; updated_at: string }

const docs = ref<Doc[]>([])
const categories = ref<Record<string, string>>({})
const router = useRouter()
const statusLabel: Record<string, string> = { draft: '草稿', published: '已上架' }

const columns = [
  { name: 'category', label: '分類', field: 'category', format: (v: string) => categories.value[v] ?? v, sortable: true },
  { name: 'title', label: '標題', field: 'title', sortable: true },
  { name: 'status', label: '狀態', field: 'status', format: (v: string) => statusLabel[v] ?? v, sortable: true },
  { name: 'chars', label: '字數', field: 'chars', sortable: true },
  { name: 'uploaded_by', label: '上傳者', field: 'uploaded_by' },
  { name: 'updated_at', label: '更新', field: 'updated_at', format: (v: string) => new Date(v).toLocaleString('zh-TW'), sortable: true },
]

async function load() {
  try {
    const res = await api.guideGet<{ documents: Doc[]; categories: Record<string, string> }>('/guide/admin/documents')
    docs.value = res.documents
    categories.value = res.categories
  } catch (e) {
    Notify.create({ type: 'negative', message: (e as Error).message })
  }
}
onMounted(load)

const showUpload = ref(false)
const form = reactive({ title: '', category: 'practice', file: null as File | null, body: '', source: 'file' as 'file' | 'text' })
const uploading = ref(false)

async function upload() {
  const fd = new FormData()
  fd.append('title', form.title)
  fd.append('category', form.category)
  if (form.source === 'file' && form.file) fd.append('file', form.file)
  else fd.append('body', form.body)
  uploading.value = true
  try {
    const res = await api.guideUpload<{ id: string }>('/guide/admin/documents', fd)
    showUpload.value = false
    Object.assign(form, { title: '', file: null, body: '' })
    router.push(`/knowledge/${res.id}`)
  } catch (e) {
    Notify.create({ type: 'negative', message: (e as Error).message })
  } finally {
    uploading.value = false
  }
}
</script>

<template>
  <q-page class="admin-page">
    <div class="row items-center q-mb-sm">
      <h1 class="q-mb-none">知識庫</h1>
      <q-space />
      <q-btn color="secondary" unelevated no-caps icon="upload" label="上傳知識" @click="showUpload = true" />
    </div>
    <p class="text-grey-8">
      上傳後是草稿，可以先在文件頁「試問」看 Sunny 會怎麼回答，確認沒問題再上架。只有「已上架」的文件會被 Sunny 使用。
      「共修腳本」是帶領共修用的流程，會出現在官網「帶我共修」的選單裡。
    </p>
    <q-table
      :rows="docs"
      :columns="columns"
      row-key="id"
      flat
      bordered
      no-data-label="還沒有任何知識文件"
      :rows-per-page-options="[25, 50, 0]"
      @row-click="(_e: Event, row: Doc) => router.push(`/knowledge/${row.id}`)"
    />

    <q-dialog v-model="showUpload">
      <q-card style="width: 600px; max-width: 95vw">
        <q-card-section class="text-h6">上傳知識</q-card-section>
        <q-card-section class="q-gutter-md">
          <q-select
            v-model="form.category"
            :options="Object.entries(categories).map(([value, label]) => ({ value, label }))"
            emit-value
            map-options
            label="分類"
            outlined
            dense
          />
          <q-btn-toggle
            v-model="form.source"
            :options="[{ label: '上傳檔案', value: 'file' }, { label: '貼上文字', value: 'text' }]"
            no-caps
            unelevated
            toggle-color="secondary"
          />
          <q-file v-if="form.source === 'file'" v-model="form.file" label="選擇檔案（.txt、.md、.docx、.pdf，20 MB 以內）" outlined accept=".txt,.md,.markdown,.docx,.pdf" />
          <q-input v-else v-model="form.body" type="textarea" label="內容" outlined autogrow input-style="min-height: 200px" />
          <q-input v-model="form.title" :label="form.source === 'file' ? '標題（空白則用檔名）' : '標題 *'" outlined dense />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat no-caps label="取消" v-close-popup />
          <q-btn
            color="secondary"
            unelevated
            no-caps
            label="上傳為草稿"
            :loading="uploading"
            :disable="form.source === 'file' ? !form.file : !form.body.trim() || !form.title.trim()"
            @click="upload"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>
