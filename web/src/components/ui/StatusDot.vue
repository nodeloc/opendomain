<template>
  <span class="inline-flex items-center gap-2 text-sm" :class="textClass">
    <span class="y-dot" :class="dotClass" aria-hidden="true"></span>
    <slot>{{ label }}</slot>
  </span>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  tone: {
    type: String,
    default: 'muted',
    validator: (v) => ['ok', 'warn', 'danger', 'info', 'muted'].includes(v),
  },
  label: { type: String, default: '' },
})

const dotClass = computed(() => ({
  ok: 'bg-y-mint',
  warn: 'bg-y-warn',
  danger: 'bg-y-danger',
  info: 'bg-y-info',
  muted: 'bg-y-faint',
}[props.tone]))

const textClass = computed(() => ({
  ok: 'text-y-text',
  warn: 'text-y-warn',
  danger: 'text-y-danger-text',
  info: 'text-y-info',
  muted: 'text-y-muted',
}[props.tone]))
</script>
