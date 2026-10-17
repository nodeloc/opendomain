<template>
  <AuthCard
    v-if="error"
    :standalone="false"
    :title="$t('uiPublic.auth.errorTitle')"
  >
    <div class="y-danger-soft flex items-start gap-3 rounded-md border px-4 py-3">
      <CircleX :size="18" :stroke-width="1.5" class="mt-0.5 shrink-0 text-y-danger-text" />
      <p class="text-sm leading-6 text-y-danger-text">{{ $t('auth.callbackError') }}</p>
    </div>
    <router-link to="/login" class="btn btn-primary mt-4 w-full">{{ $t('auth.backToLogin') }}</router-link>
  </AuthCard>

  <AuthCard
    v-else
    :standalone="false"
    :title="$t('auth2.callbackTitle')"
    :description="$t('uiPublic.auth.loggingInDesc')"
  >
    <YLoader :size="44" compact />
  </AuthCard>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { CircleX } from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'
import AuthCard from '../components/ui/AuthCard.vue'
import YLoader from '../components/brand/YLoader.vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const error = ref(false)

onMounted(async () => {
  const token = route.query.token
  if (!token) {
    error.value = true
    return
  }

  authStore.token = token
  localStorage.setItem('token', token)

  try {
    await authStore.fetchProfile()
    router.replace('/dashboard')
  } catch (e) {
    error.value = true
  }
})
</script>
