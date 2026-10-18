<template>
  <AuthCard :title="$t('auth2.welcome', { siteName })">
    <template #description>
      <i18n-t keypath="auth2.loginSub" tag="span" scope="global">
        <template #domain><span class="font-mono text-y-text">{{ exampleDomain }}</span></template>
      </i18n-t>
    </template>

    <div class="flex flex-col gap-4">
      <template v-if="hasOAuth">
        <OAuthButtons :oauth="siteConfigStore.oauth" :api-base="apiBase" />
        <div class="flex items-center gap-3 font-mono text-[11px] uppercase tracking-[0.08em] text-y-faint" role="separator">
          <span class="h-px flex-1 bg-y-line"></span>
          {{ $t('auth2.or') }}
          <span class="h-px flex-1 bg-y-line"></span>
        </div>
      </template>

      <form class="flex flex-col gap-3.5" @submit.prevent="handleLogin">
        <label class="flex flex-col gap-1.5">
          <span class="text-[13px] text-y-muted">{{ $t('login.email') }}</span>
          <input
            v-model="form.email"
            type="email"
            :placeholder="$t('auth2.emailPlaceholder')"
            class="input input-bordered w-full"
            autocomplete="email"
            required
          />
        </label>

        <label class="flex flex-col gap-1.5">
          <span class="flex items-center justify-between">
            <span class="text-[13px] text-y-muted">{{ $t('login.password') }}</span>
            <router-link to="/forgot-password" class="text-xs text-y-faint hover:text-y-text">{{ $t('login.forgotPassword') }}</router-link>
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

        <div v-if="error" class="alert alert-error py-3 text-sm" role="alert">
          <CircleX :size="18" :stroke-width="1.5" class="shrink-0" />
          <span>{{ error }}</span>
        </div>

        <button type="submit" class="btn btn-primary mt-1.5 w-full" :disabled="loading">
          <span v-if="loading" class="loading loading-spinner loading-sm"></span>
          <span v-else>{{ $t('auth2.signIn') }}</span>
        </button>
      </form>
    </div>

    <template #footer>
      <p class="text-sm text-y-muted">
        {{ $t('auth2.noAccount') }}
        <router-link to="/register" class="font-medium text-y-mint-text hover:underline">{{ $t('auth2.register') }}</router-link>
      </p>
      <AuthLegal />
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
import AuthLegal from '../components/auth/AuthLegal.vue'

const router = useRouter()
const authStore = useAuthStore()
const siteConfigStore = useSiteConfigStore()
const apiBase = import.meta.env.VITE_API_BASE_URL || ''

const siteName = computed(() => siteConfigStore.siteName || 'Yunamo')
const exampleDomain = 'myproject.yunamo.com'

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
