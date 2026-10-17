<template>
  <AuthCard
    :eyebrow="$t('authPage.registerEyebrow')"
    :title="$t('authPage.registerTitle')"
    :description="$t('authPage.registerDesc')"
  >
    <YLoader v-if="configLoading" />

    <template v-else>
      <form v-if="siteConfig.allow_password_register" class="space-y-4" @submit.prevent="handleRegister">
        <label class="y-field">
          <span class="y-field-label">{{ $t('register.username') }}</span>
          <input
            v-model="form.username"
            type="text"
            :placeholder="$t('register.usernamePlaceholder')"
            class="input input-bordered w-full"
            autocomplete="username"
            required
            minlength="3"
          />
        </label>

        <label class="y-field">
          <span class="y-field-label">{{ $t('register.email') }}</span>
          <input
            v-model="form.email"
            type="email"
            :placeholder="$t('register.emailPlaceholder')"
            class="input input-bordered w-full"
            autocomplete="email"
            required
          />
        </label>

        <label class="y-field">
          <span class="y-field-label">{{ $t('register.password') }}</span>
          <input
            v-model="form.password"
            type="password"
            :placeholder="$t('register.passwordPlaceholder')"
            class="input input-bordered w-full"
            autocomplete="new-password"
            required
            minlength="6"
          />
        </label>

        <label class="y-field">
          <span class="y-field-label">{{ $t('register.inviteCode') }}</span>
          <input
            v-model="form.invite_code"
            type="text"
            :placeholder="$t('register.inviteCodePlaceholder')"
            class="input input-bordered w-full font-mono"
          />
        </label>

        <div v-if="error" class="alert alert-error py-3 text-sm">
          <CircleX :size="18" :stroke-width="1.5" class="shrink-0" />
          <span>{{ error }}</span>
        </div>

        <div v-if="success" class="alert alert-success py-3 text-sm">
          <CircleCheck :size="18" :stroke-width="1.5" class="shrink-0 text-y-mint-text" />
          <span>{{ success }}</span>
        </div>

        <button type="submit" class="btn btn-primary w-full" :disabled="loading">
          <span v-if="loading" class="loading loading-spinner loading-sm"></span>
          <span v-else>{{ $t('register.signUpButton') }}</span>
        </button>
      </form>

      <div v-else class="alert alert-warning py-3 text-sm">
        <TriangleAlert :size="18" :stroke-width="1.5" class="shrink-0" />
        <span>{{ $t('register.registrationDisabled') }}</span>
      </div>

      <template v-if="hasOAuth">
        <div v-if="siteConfig.allow_password_register" class="my-6 flex items-center gap-3 text-xs text-y-faint">
          <span class="h-px flex-1 bg-y-line"></span>
          <span class="font-mono uppercase tracking-[0.08em]">{{ $t('authPage.continueWith') }}</span>
          <span class="h-px flex-1 bg-y-line"></span>
        </div>
        <OAuthButtons :class="{ 'mt-4': !siteConfig.allow_password_register }" :oauth="siteConfig.oauth" :api-base="apiBase" />
      </template>
    </template>

    <template #footer>
      {{ $t('register.haveAccount') }}
      <router-link to="/login" class="font-medium text-y-mint-text hover:underline">{{ $t('register.login') }}</router-link>
    </template>
  </AuthCard>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { CircleX, CircleCheck, TriangleAlert } from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'
import { useSiteConfigStore } from '../stores/siteConfig'
import AuthCard from '../components/ui/AuthCard.vue'
import OAuthButtons from '../components/ui/OAuthButtons.vue'
import YLoader from '../components/brand/YLoader.vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const siteConfigStore = useSiteConfigStore()
const apiBase = import.meta.env.VITE_API_BASE_URL || ''

const configLoading = computed(() => !siteConfigStore.loaded)
const siteConfig = computed(() => ({
  allow_password_register: siteConfigStore.allowPasswordRegister,
  oauth: siteConfigStore.oauth,
}))

const hasOAuth = computed(() => {
  const o = siteConfigStore.oauth
  return o?.github || o?.google || o?.nodeloc
})

const form = ref({
  username: '',
  email: '',
  password: '',
  invite_code: route.query.ref || '',
})

const loading = ref(false)
const error = ref('')
const success = ref('')

const handleRegister = async () => {
  loading.value = true
  error.value = ''
  success.value = ''

  const payload = {
    username: form.value.username,
    email: form.value.email,
    password: form.value.password,
  }

  if (form.value.invite_code) {
    payload.invite_code = form.value.invite_code
  }

  const result = await authStore.register(payload)

  if (result.success) {
    success.value = result.message
    setTimeout(() => {
      router.push('/login')
    }, 2000)
  } else {
    error.value = result.message
  }

  loading.value = false
}
</script>
