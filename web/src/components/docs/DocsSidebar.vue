<template>
  <nav :aria-label="$t('docs.nav')" class="flex flex-col gap-6 text-sm leading-5">
    <div v-if="loading && !groups.length" class="flex flex-col gap-2 px-3" aria-hidden="true">
      <span class="h-3 w-20 rounded-sm bg-y-surface-2"></span>
      <span v-for="n in 5" :key="n" class="h-4 rounded-sm bg-y-surface"></span>
    </div>

    <div v-if="failed && !groups.length" class="flex flex-col items-start gap-2 px-3 text-y-faint">
      <span>{{ $t('docs.pagesUnavailable') }}</span>
      <button type="button" class="text-y-mint-text hover:underline" @click="$emit('retry')">{{ $t('docs.retry') }}</button>
    </div>

    <div v-for="group in allGroups" :key="group.key" class="flex flex-col gap-0.5">
      <div class="px-3 pb-1.5 font-mono text-[11px] uppercase tracking-[0.1em] text-y-faint">{{ group.label }}</div>
      <router-link
        v-for="item in group.items"
        :key="item.to"
        :to="item.to"
        class="docs-side-link"
        :class="{ 'is-active': isActive(item) }"
        :aria-current="isActive(item) ? 'page' : undefined"
      >{{ item.title }}</router-link>
    </div>
  </nav>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const props = defineProps({
  groups: { type: Array, default: () => [] },
  policies: { type: Object, default: null },
  loading: { type: Boolean, default: false },
  failed: { type: Boolean, default: false },
})
defineEmits(['retry'])

const route = useRoute()
const allGroups = computed(() => (props.policies ? [...props.groups, props.policies] : props.groups))
const isActive = (item) => route.path === item.to
</script>

<style scoped>
.docs-side-link {
  position: relative;
  display: block;
  border-radius: 8px;
  padding: 6px 12px;
  color: rgb(var(--y-muted));
  transition: color 150ms, background-color 150ms;
}
.docs-side-link:hover {
  color: rgb(var(--y-text));
  background: rgb(var(--y-surface));
}
.docs-side-link.is-active {
  color: rgb(var(--y-text));
  background: rgb(var(--y-surface));
}
.docs-side-link.is-active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 6px;
  bottom: 6px;
  width: 2px;
  border-radius: 2px;
  background: rgb(var(--y-mint));
}
</style>
