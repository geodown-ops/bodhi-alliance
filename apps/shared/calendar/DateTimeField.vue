<script setup lang="ts">
// 日期＋時間欄位（24 小時制）：日期點了跳月曆，時間可以打字或點時鐘挑。
// 值是本地時間字串 YYYY-MM-DDTHH:mm；hideTime 時只顯示日期（整天的行程）。
import { computed } from 'vue'

const props = defineProps<{ modelValue: string; label: string; hideTime?: boolean }>()
const emit = defineEmits<{ 'update:modelValue': [v: string] }>()

const date = computed(() => props.modelValue.slice(0, 10).replace(/-/g, '/'))
const time = computed(() => props.modelValue.slice(11, 16))
function set(d: string, t: string) {
  if (!/^\d{4}\/\d{2}\/\d{2}$/.test(d)) return
  emit('update:modelValue', `${d.replace(/\//g, '-')}T${/^([01]\d|2[0-3]):[0-5]\d$/.test(t) ? t : '09:00'}`)
}
const timeRule = (v: string) => !v || /^([01]\d|2[0-3]):[0-5]\d$/.test(v) || '24 小時制，例如 19:00'
</script>

<template>
  <div class="dt-field">
    <q-input :model-value="date" :label="label" stack-label outlined dense readonly class="cursor-pointer dt-date">
      <template #append><q-icon name="event" /></template>
      <q-popup-proxy cover transition-show="scale" transition-hide="scale">
        <q-date :model-value="date" mask="YYYY/MM/DD" minimal first-day-of-week="1" @update:model-value="(v: string | null) => v && set(v, time)">
          <div class="row justify-end"><q-btn v-close-popup flat no-caps label="確定" /></div>
        </q-date>
      </q-popup-proxy>
    </q-input>
    <q-input
      v-if="!hideTime"
      :model-value="time"
      label="時間"
      stack-label
      outlined
      dense
      mask="##:##"
      :rules="[timeRule]"
      hide-bottom-space
      class="dt-time"
      @update:model-value="(v) => /^([01]\d|2[0-3]):[0-5]\d$/.test(String(v)) && set(date, String(v))"
    >
      <template #append>
        <q-icon name="schedule" class="cursor-pointer">
          <q-popup-proxy cover transition-show="scale" transition-hide="scale">
            <q-time :model-value="time" format24h :minute-options="[0, 15, 30, 45]" @update:model-value="(v: string | null) => v && set(date, v)">
              <div class="row justify-end"><q-btn v-close-popup flat no-caps label="確定" /></div>
            </q-time>
          </q-popup-proxy>
        </q-icon>
      </template>
    </q-input>
  </div>
</template>

<style scoped>
.dt-field {
  display: grid;
  grid-template-columns: 1.3fr 1fr;
  gap: 8px;
  align-items: start;
}
.dt-date :deep(.q-field__control:before) {
  border-style: solid;
}
</style>
