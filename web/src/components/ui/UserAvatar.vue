<template>
  <span
    class="inline-flex shrink-0 items-center justify-center overflow-hidden rounded-full border border-y-line bg-y-surface-2 font-mono font-medium uppercase text-y-mint-text"
    :style="{ width: size + 'px', height: size + 'px', fontSize: Math.round(size * 0.4) + 'px' }"
  >
    <img v-if="src && !failed" :src="src" :alt="name" class="h-full w-full object-cover" @error="failed = true" />
    <span v-else>{{ initials }}</span>
  </span>
</template>

<script setup>
import { computed, ref } from 'vue'

const props = defineProps({
  src: { type: String, default: '' },
  name: { type: String, default: '' },
  size: { type: Number, default: 32 },
})

const failed = ref(false)
const initials = computed(() => (props.name || '?').trim().slice(0, 2))
</script>
