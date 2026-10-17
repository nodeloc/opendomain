<template>
  <AuthCard
    :eyebrow="$t('authPage.loginEyebrow')"
    :title="$t('authPage.loginTitle', { siteName })"
    :description="$t('authPage.loginDesc')"
  >
    <form class="space-y-4" @submit.prevent="handleLogin">
      <label class="y-field">
        <span class="y-field-label">{{ $t('login.email') }}</span>
        <input
          v-model="form.email"
          type="email"
          :placeholder="$t('login.emailPlaceholder')"
          class="input input-bordered w-full"
          autocomplete="email"
          required
        />
      </label>

      <label class="y-field">
        <span class="flex items-center justify-between">
          <span class="y-field-label">{{ $t('login.password') }}</span>
          <a href="#" class="text-xs text-y-faint hover:text-y-text">{{ $t('login.forgotPassword') }}</a>
        </span>
        <input
          v-model="form.password"
          type="password"
          :placeholder="$t('login.passwordPlaceholder')"
          class="input input-bordered w-full"
          autocomplete="current-password"
          required
        />
      </label>

      <div v-if="error" class="alert alert-error py-3 text-sm">
        <CircleX :size="18" :stroke-width="1.5" class="shrink-0" />
        <span>{{ error }}</span>
      </div>

      <button type="submit" class="btn btn-primary w-full" :disabled="loading">
        <span v-if="loading" class="loading loading-spinner loading-sm"></span>
        <span v-else>{{ $t('login.loginButton') }}</span>
      </button>
    </form>

    <template v-if="hasOAuth">
      <div class="my-6 flex items-center gap-3 text-xs text-y-faint">
        <span class="h-px flex-1 bg-y-line"></span>
        <span class="font-mono uppercase tracking-[0.08em]">{{ $t('authPage.continueWith') }}</span>
        <span class="h-px flex-1 bg-y-line"></span>
      </div>
      <OAuthButtons :oauth="siteConfigStore.oauth" :api-base="apiBase" />
    </template>

    <template #footer>
      {{ $t('login.noAccount') }}
      <router-link to="/register" class="font-medium text-y-mint-text hover:underline">{{ $t('login.signUp') }}</router-link>
    </template>
  </AuthCard>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { CircleX } from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'
import { useSiteConfigStore } from '../stores/siteConfig'
import AuthCard from '../components/ui/AuthCard.vue'
import OAuthButtons from '../components/ui/OAuthButtons.vue'

const router = useRouter()
const authStore = useAuthStore()
const siteConfigStore = useSiteConfigStore()
const apiBase = import.meta.env.VITE_API_BASE_URL || ''

const siteName = computed(() => siteConfigStore.siteName || 'Yunamo')

const hasOAuth = computed(() => {
  const o = siteConfigStore.oauth
  return o?.github || o?.google || o?.nodeloc
})

const form = ref({
  email: '',
  password: '',
})

const loading = ref(false)
const error = ref('')

const handleLogin = async () => {
  loading.value = true
  error.value = ''

  const result = await authStore.login(form.value)

  if (result.success) {
    await authStore.fetchProfile()
    router.push('/dashboard')
  } else {
    error.value = result.message
  }

  loading.value = false
}
</script>
