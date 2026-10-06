<script setup lang="ts">
import { computed } from 'vue'
import { references } from '../mindfulness'

const props = defineProps<{ keys: string[] }>()
const list = computed(() => [...new Set(props.keys)].map((k) => references[k]).filter(Boolean))
</script>

<template>
  <ol class="refs">
    <li v-for="r in list" :key="r.text">
      {{ r.text }}
      <a v-if="r.url" :href="r.url" target="_blank" rel="noopener">{{ r.url.replace('https://', '') }}</a>
    </li>
  </ol>
</template>

<style scoped>
.refs {
  padding-left: 1.4rem;
}
.refs li {
  font-size: 0.9rem !important;
  line-height: 1.6 !important;
  margin-bottom: 8px;
  overflow-wrap: anywhere;
}
.refs a {
  color: var(--leaf);
  margin-left: 4px;
}
</style>
