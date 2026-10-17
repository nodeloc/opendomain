<template>
  <AuthCard :title="$t('auth2.registerTitle', { siteName })">
    <template #description>
      <i18n-t keypath="auth2.registerSub" tag="span" scope="global">
        <template #domain><span class="font-mono text-y-text">{{ exampleDomain }}</span></template>
      </i18n-t>
    </template>

    <YLoader v-if="configLoading" compact />

    <div v-else class="flex flex-col gap-4">
      <template v-if="hasOAuth">
        <OAuthButtons :oauth="siteConfig.oauth" :api-base="apiBase" />
        <div
          v-if="siteConfig.allow_password_register"
          class="flex items-center gap-3 font-mono text-[11px] uppercase tracking-[0.08em] text-y-faint"
          role="separator"
        >
          <span class="h-px flex-1 bg-y-line"></span>
          {{ $t('auth2.or') }}
          <span class="h-px flex-1 bg-y-line"></span>
        </div>
      </template>

      <form v-if="siteConfig.allow_password_register" class="flex flex-col gap-3.5" @submit.prevent="handleRegister">
        <label class="flex flex-col gap-1.5">
          <span class="text-[13px] text-y-muted">{{ $t('register.username') }}</span>
          <input
            v-model="form.username"
            type="text"
            :placeholder="$t('auth2.usernamePlaceholder')"
            class="input input-bordered w-full"
            autocomplete="username"
            required
            minlength="3"
          />
        </label>

        <label class="flex flex-col gap-1.5">
          <span class="text-[13px] text-y-muted">{{ $t('register.email') }}</span>
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
          <span class="text-[13px] text-y-muted">{{ $t('register.password') }}</span>
          <input
            v-model="form.password"
            type="password"
            :placeholder="$t('auth2.passwordNewPlaceholder')"
            class="input input-bordered w-full"
            autocomplete="new-password"
            required
            minlength="6"
          />
        </label>

        <label class="flex flex-col gap-1.5">
          <span class="flex items-center justify-between">
            <span class="text-[13px] text-y-muted">{{ $t('auth2.inviteCode') }}</span>
            <span class="text-xs text-y-faint">{{ $t('auth2.optional') }}</span>
          </span>
          <input
            v-model="form.invite_code"
            type="text"
            :placeholder="$t('auth2.invitePlaceholder')"
            class="input input-bordered w-full font-mono"
          />
        </label>

        <div v-if="error" class="alert alert-error py-3 text-sm" role="alert">
          <CircleX :size="18" :stroke-width="1.5" class="shrink-0" />
          <span>{{ error }}</span>
        </div>

        <div v-if="success" class="alert alert-success py-3 text-sm" role="status">
          <CircleCheck :size="18" :stroke-width="1.5" class="shrink-0 text-y-mint-text" />
          <span>{{ success }}</span>
        </div>

        <button type="submit" class="btn btn-primary mt-1.5 w-full" :disabled="loading">
          <span v-if="loading" class="loading loading-spinner loading-sm"></span>
          <span v-else>{{ $t('auth2.createAccount') }}</span>
        </button>
      </form>

      <div v-else class="alert alert-warning py-3 text-sm">
        <TriangleAlert :size="18" :stroke-width="1.5" class="shrink-0" />
        <span>{{ $t('register.registrationDisabled') }}</span>
      </div>
    </div>

    <template #footer>
      <p class="text-sm text-y-muted">
        {{ $t('auth2.haveAccount') }}
        <router-link to="/login" class="font-medium text-y-mint-text hover:underline">{{ $t('auth2.signInLink') }}</router-link>
      </p>
      <AuthLegal />
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
import AuthLegal from '../components/auth/AuthLegal.vue'
import YLoader from '../components/brand/YLoader.vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const siteConfigStore = useSiteConfigStore()
const apiBase = import.meta.env.VITE_API_BASE_URL || ''

const siteName = computed(() => siteConfigStore.siteName || 'Yunamo')
const exampleDomain = 'myproject.yunamo.com'

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
