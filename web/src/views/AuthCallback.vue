<template>
  <div class="min-h-[70vh] flex items-center justify-center px-4">
    <div class="text-center">
      <div v-if="error" class="space-y-4">
        <div class="alert alert-error max-w-md">
          <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <span>{{ $t('auth.callbackError') }}</span>
        </div>
        <router-link to="/login" class="btn btn-primary">{{ $t('auth.backToLogin') }}</router-link>
      </div>
      <div v-else>
        <span class="loading loading-spinner loading-lg"></span>
        <p class="mt-4 text-base-content/70">{{ $t('auth.loggingIn') }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'

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
