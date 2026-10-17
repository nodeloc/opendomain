<template>
  <div class="flex flex-col items-center justify-center gap-3" :class="compact ? 'py-4' : 'py-12'" role="status" :aria-label="label || 'Loading'">
    <svg viewBox="0 0 64 64" xmlns="http://www.w3.org/2000/svg" :style="{ width: size + 'px', height: size + 'px' }" class="y-loader">
      <path
        d="M32 54V35M32 35C32 26.5 24.5 22 16.5 14.5M32 35C32 26.5 39.5 22 47.5 14.5"
        stroke="currentColor"
        stroke-width="6.5"
        stroke-linecap="round"
        stroke-linejoin="round"
        fill="none"
      />
      <circle cx="15" cy="13" r="7" class="node node-a" />
      <circle cx="49" cy="13" r="7" class="node node-b" />
    </svg>
    <span v-if="label" class="font-mono text-xs uppercase tracking-[0.08em] text-y-faint">{{ label }}</span>
  </div>
</template>

<script setup>
defineProps({
  size: { type: Number, default: 36 },
  label: { type: String, default: '' },
  compact: { type: Boolean, default: false },
})
</script>

<style scoped>
.y-loader {
  color: rgb(var(--y-line-strong));
}
.node {
  fill: rgb(var(--y-node));
  transform-box: fill-box;
  transform-origin: center;
  animation: node-pulse 1.2s cubic-bezier(0.2, 0.8, 0.2, 1) infinite;
}
.node-b {
  animation-delay: 0.6s;
}
@keyframes node-pulse {
  0%, 100% { opacity: 0.35; transform: scale(0.85); }
  50% { opacity: 1; transform: scale(1.1); }
}
</style>
