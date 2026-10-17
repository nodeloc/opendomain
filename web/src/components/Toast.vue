<template>
  <div class="pointer-events-none fixed right-4 top-4 z-[9999] flex w-[min(24rem,calc(100vw-2rem))] flex-col gap-2">
    <TransitionGroup name="toast">
      <div
        v-for="notification in notifications"
        :key="notification.id"
        class="pointer-events-auto flex items-start gap-3 rounded-md border bg-y-surface px-4 py-3 text-sm"
        :class="toneClass(notification.type)"
        role="status"
      >
        <component :is="iconFor(notification.type)" :size="18" :stroke-width="1.5" class="mt-0.5 shrink-0" :class="iconClass(notification.type)" />
        <span class="flex-1 whitespace-pre-line leading-5 text-y-text">{{ notification.message }}</span>
        <button type="button" class="-mr-1 shrink-0 rounded-sm p-0.5 text-y-faint hover:text-y-text" :aria-label="$t('common.close')" @click="remove(notification.id)">
          <X :size="16" :stroke-width="1.5" />
        </button>
      </div>
    </TransitionGroup>
  </div>
</template>

<script setup>
import { CircleCheck, CircleX, TriangleAlert, Info, X } from 'lucide-vue-next'
import { useToast } from '../composables/useToast'

const { notifications, remove } = useToast()

const iconFor = (type) => ({ success: CircleCheck, error: CircleX, warning: TriangleAlert }[type] || Info)

const toneClass = (type) => ({
  success: 'border-y-mint/40',
  error: 'border-y-danger/40',
  warning: 'border-y-warn/40',
}[type] || 'border-y-line-strong')

const iconClass = (type) => ({
  success: 'text-y-mint-text',
  error: 'text-y-danger-text',
  warning: 'text-y-warn',
}[type] || 'text-y-info')
</script>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.25s cubic-bezier(0.2, 0.8, 0.2, 1);
}
.toast-enter-from {
  opacity: 0;
  transform: translateY(-8px);
}
.toast-leave-to {
  opacity: 0;
  transform: translateX(16px);
}
.toast-move {
  transition: transform 0.25s ease;
}
</style>
