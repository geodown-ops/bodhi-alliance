<script setup lang="ts">
// 密碼欄：右邊有眼睛圖示，按住（滑鼠或手指）時顯示密碼，放開就再藏起來
import { ref } from 'vue'

const model = defineModel<string>({ default: '' })
const show = ref(false)
</script>

<template>
  <q-input v-model="model" :type="show ? 'text' : 'password'">
    <template #append>
      <q-icon
        :name="show ? 'visibility' : 'visibility_off'"
        class="eye"
        role="button"
        tabindex="0"
        :aria-label="show ? '放開隱藏密碼' : '按住顯示密碼'"
        title="按住顯示密碼"
        @pointerdown.prevent="show = true"
        @pointerup="show = false"
        @pointerleave="show = false"
        @pointercancel="show = false"
        @contextmenu.prevent
        @keydown.space.prevent="show = true"
        @keyup.space="show = false"
        @blur="show = false"
      />
    </template>
  </q-input>
</template>

<style scoped>
.eye {
  cursor: pointer;
  touch-action: none;
  user-select: none;
  -webkit-user-select: none;
  -webkit-touch-callout: none;
}
</style>
