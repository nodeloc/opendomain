<template>
  <img
    v-if="src && !failed"
    :src="src"
    :alt="name"
    loading="lazy"
    decoding="async"
    referrerpolicy="no-referrer"
    class="h-full w-auto max-w-full object-contain object-left"
    @error="failed = true"
  />
  <span
    v-else
    class="sponsor-monogram inline-flex h-full max-h-full aspect-square items-center justify-center rounded-md border font-mono font-semibold uppercase text-y-mint-text"
    :class="small ? 'text-xs' : 'text-lg'"
    role="img"
    :aria-label="name"
  >{{ monogram }}</span>
</template>

<script setup>
import { ref, computed, watch } from 'vue'

const props = defineProps({
  src: { type: String, default: '' },
  name: { type: String, default: '' },
  small: { type: Boolean, default: false },
})

const failed = ref(false)
watch(() => props.src, () => {
  failed.value = false
})

const monogram = computed(() => {
  const name = (props.name || '').trim()
  if (!name) return '?'
  const first = Array.from(name)[0]
  // Han and other non-Latin scripts read better as a single glyph than as initials.
  if (!/[a-z0-9]/i.test(first)) return first
  const words = name.split(/[\s._-]+/).filter(Boolean)
  const letters = words.length > 1 ? words[0][0] + words[1][0] : name.slice(0, 2)
  return letters.toUpperCase()
})
</script>

<style scoped>
.sponsor-monogram {
  background: rgb(var(--y-mint) / 0.08);
  border-color: rgb(var(--y-mint) / 0.3);
}
</style>
