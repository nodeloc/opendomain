<template>
  <div ref="rootRef" class="relative w-full">
    <label class="relative flex w-full">
      <span class="sr-only">{{ $t('docs.searchLabel') }}</span>
      <Search :size="16" :stroke-width="1.75" class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-y-faint" />
      <input
        ref="inputRef"
        v-model="query"
        type="search"
        role="combobox"
        :aria-expanded="open"
        aria-controls="docs-search-results"
        aria-autocomplete="list"
        :aria-activedescendant="open && results[active] ? `docs-search-opt-${active}` : undefined"
        :placeholder="$t('docs.search')"
        class="docs-search-input h-10 w-full rounded-md border border-y-line bg-y-surface pl-9 pr-16 text-sm text-y-text placeholder:text-y-faint"
        autocomplete="off"
        spellcheck="false"
        @focus="open = true"
        @blur="open = false"
        @keydown="onKeydown"
      />
      <span class="y-kbd pointer-events-none absolute right-2.5 top-1/2 -translate-y-1/2" aria-hidden="true">{{ shortcut }}</span>
    </label>

    <div
      v-if="open"
      id="docs-search-results"
      role="listbox"
      class="absolute left-0 right-0 top-full z-50 mt-2 max-h-[min(24rem,70vh)] overflow-y-auto rounded-md border border-y-line bg-y-surface p-1.5"
    >
      <template v-if="results.length">
        <button
          v-for="(item, i) in results"
          :id="`docs-search-opt-${i}`"
          :key="item.to"
          type="button"
          role="option"
          :aria-selected="i === active"
          class="flex w-full items-center justify-between gap-3 rounded-sm px-3 py-2 text-left text-sm transition-colors"
          :class="i === active ? 'bg-y-surface-2 text-y-text' : 'text-y-muted'"
          @mousedown.prevent
          @mouseenter="active = i"
          @click="go(item)"
        >
          <span class="truncate">{{ item.title }}</span>
          <span class="shrink-0 font-mono text-[11px] uppercase tracking-[0.08em] text-y-faint">{{ item.group }}</span>
        </button>
      </template>
      <p v-else class="px-3 py-2 text-sm text-y-faint">{{ $t('docs.noResults', { q: query }) }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { Search } from 'lucide-vue-next'

const props = defineProps({
  items: { type: Array, default: () => [] },
})

const router = useRouter()
const rootRef = ref(null)
const inputRef = ref(null)
const query = ref('')
const open = ref(false)
const active = ref(0)

const isMac = typeof navigator !== 'undefined' && /Mac|iPhone|iPad/.test(navigator.platform || navigator.userAgent)
const shortcut = isMac ? '⌘K' : 'Ctrl K'

const results = computed(() => {
  const q = query.value.trim().toLowerCase()
  if (!q) return props.items
  return props.items
    .map((item) => {
      const title = item.title.toLowerCase()
      const at = title.indexOf(q)
      const slugHit = item.slug && item.slug.includes(q)
      return { item, score: at === 0 ? 0 : at > 0 ? 1 : slugHit ? 2 : -1 }
    })
    .filter((r) => r.score >= 0)
    .sort((a, b) => a.score - b.score)
    .map((r) => r.item)
})

watch(query, () => {
  active.value = 0
  open.value = true
})

const go = (item) => {
  if (!item) return
  query.value = ''
  open.value = false
  inputRef.value?.blur()
  router.push(item.to)
}

const onKeydown = (e) => {
  if (e.key === 'ArrowDown') {
    e.preventDefault()
    open.value = true
    active.value = Math.min(active.value + 1, results.value.length - 1)
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()
    active.value = Math.max(active.value - 1, 0)
  } else if (e.key === 'Enter') {
    e.preventDefault()
    go(results.value[active.value])
  } else if (e.key === 'Escape') {
    open.value = false
    inputRef.value?.blur()
  }
}

const onGlobalKeydown = (e) => {
  if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 'k') {
    e.preventDefault()
    inputRef.value?.focus()
    inputRef.value?.select()
  }
}

onMounted(() => window.addEventListener('keydown', onGlobalKeydown))
onBeforeUnmount(() => window.removeEventListener('keydown', onGlobalKeydown))
</script>

<style scoped>
.docs-search-input {
  outline: none;
  transition: border-color 150ms;
}
.docs-search-input:focus {
  border-color: rgb(var(--y-mint) / 0.7);
}
.docs-search-input::-webkit-search-cancel-button {
  display: none;
}
</style>
