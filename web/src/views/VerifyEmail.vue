<template>
  <AuthCard :title="title">
    <template v-if="description" #description>{{ description }}</template>

    <YLoader v-if="state === 'verifying'" compact :label="$t('authEmail.verify.verifying')" />

    <div v-else-if="state === 'success'" class="flex flex-col gap-3.5">
      <div class="alert alert-success py-3 text-sm" role="status">
        <CircleCheck :size="18" :stroke-width="1.5" class="shrink-0 text-y-mint-text" />
        <span>
          <i18n-t v-if="verifiedEmail" keypath="authEmail.verify.success" tag="span" scope="global">
            <template #email><span class="break-all font-medium">{{ verifiedEmail }}</span></template>
          </i18n-t>
          <span v-else>{{ $t('authEmail.verify.successNoEmail') }}</span>
          <span class="block">{{ $t('authEmail.verify.successHint') }}</span>
        </span>
      </div>
      <router-link v-if="isAuthenticated" to="/domains" class="btn btn-primary mt-1.5 w-full">{{ $t('authEmail.verify.goDomains') }}</router-link>
      <router-link v-else to="/login" class="btn btn-primary mt-1.5 w-full">{{ $t('authEmail.verify.signIn') }}</router-link>
    </div>

    <div v-else class="flex flex-col gap-3.5">
      <div class="alert py-3 text-sm" :class="state === 'error' ? 'alert-error' : 'alert-warning'" role="alert">
        <CircleX v-if="state === 'error'" :size="18" :stroke-width="1.5" class="shrink-0" />
        <TriangleAlert v-else :size="18" :stroke-width="1.5" class="shrink-0" />
        <span>{{ failureMessage }}</span>
      </div>

      <button v-if="state === 'error'" type="button" class="btn btn-outline mt-1.5 w-full" @click="verify(true)">
        {{ $t('authEmail.verify.retry') }}
      </button>

      <button
        v-if="isAuthenticated"
        type="button"
        class="btn btn-primary w-full"
        :class="{ 'mt-1.5': state !== 'error' }"
        :disabled="resendDisabled"
        @click="handleResend"
      >
        <span v-if="resending" class="loading loading-spinner loading-sm"></span>
        <span v-else-if="secondsLeft > 0">{{ $t('authEmail.banner.resendIn', { s: secondsLeft }) }}</span>
        <span v-else>{{ $t('authEmail.verify.sendNew') }}</span>
      </button>
      <router-link v-else to="/login" class="btn btn-primary w-full" :class="{ 'mt-1.5': state !== 'error' }">
        {{ $t('authEmail.verify.signInToResend') }}
      </router-link>
    </div>

    <template #footer>
      <AuthLegal />
    </template>
  </AuthCard>
</template>

<script>
// Keyed by token at module scope so a remount (HMR, quick re-navigation) reuses the
// in-flight request instead of spending a single-use token twice.
const inflight = new Map()
</script>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { CircleX, CircleCheck, TriangleAlert } from 'lucide-vue-next'
import axios from '../utils/axios'
import { useAuthStore } from '../stores/auth'
import { useResendVerification } from '../composables/useResendVerification'
import AuthCard from '../components/ui/AuthCard.vue'
import AuthLegal from '../components/auth/AuthLegal.vue'
import YLoader from '../components/brand/YLoader.vue'

const route = useRoute()
const { t } = useI18n()
const authStore = useAuthStore()
const { resend, sending: resending, secondsLeft, disabled: resendDisabled } = useResendVerification()

const token = typeof route.query.token === 'string' ? route.query.token.trim() : ''

// verifying | success | invalid | missing | error
const state = ref(token ? 'verifying' : 'missing')
const verifiedEmail = ref('')
const isAuthenticated = computed(() => authStore.isAuthenticated)

const title = computed(() => {
  if (state.value === 'success') return t('authEmail.verify.successTitle')
  if (state.value === 'invalid') return t('authEmail.verify.invalidTitle')
  return t('authEmail.verify.title')
})

const description = computed(() => (state.value === 'verifying' ? t('authEmail.verify.verifyingDesc') : ''))

const failureMessage = computed(() => {
  if (state.value === 'missing') return t('authEmail.verify.missing')
  if (state.value === 'invalid') return t('authEmail.verify.invalid')
  return t('authEmail.verify.failed')
})

const markVerified = async (email, refresh = true) => {
  verifiedEmail.value = email || authStore.user?.email || ''
  state.value = 'success'
  if (refresh && authStore.isAuthenticated) await authStore.fetchProfile()
}

const verify = async (retry = false) => {
  if (!token) return
  if (retry) inflight.delete(token)
  state.value = 'verifying'

  let request = inflight.get(token)
  if (!request) {
    request = axios.post('/api/auth/verify-email', { token })
    inflight.set(token, request)
  }

  try {
    const { data } = await request
    await markVerified(data?.email)
  } catch (err) {
    inflight.delete(token)
    state.value = err.response?.data?.code === 'invalid_token' ? 'invalid' : 'error'
  }
}

const handleResend = async () => {
  const result = await resend()
  // resend() has already refreshed the profile in this case.
  if (result === 'already_verified') await markVerified('', false)
}

onMounted(() => verify())
</script>
