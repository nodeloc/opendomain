<template>
  <AuthCard :title="title">
    <template #description>{{ description }}</template>

    <div v-if="state === 'missing' || state === 'invalid'" class="flex flex-col gap-3.5">
      <div class="alert alert-warning py-3 text-sm" role="alert">
        <TriangleAlert :size="18" :stroke-width="1.5" class="shrink-0" />
        <span>{{ state === 'missing' ? $t('authEmail.reset.missing') : $t('authEmail.reset.invalid') }}</span>
      </div>
      <router-link to="/forgot-password" class="btn btn-primary mt-1.5 w-full">{{ $t('authEmail.reset.requestNew') }}</router-link>
    </div>

    <div v-else-if="state === 'done'" class="flex flex-col gap-3.5">
      <div class="alert alert-success py-3 text-sm" role="status">
        <CircleCheck :size="18" :stroke-width="1.5" class="shrink-0 text-y-mint-text" />
        <span>{{ $t('authEmail.reset.done') }}</span>
      </div>
      <router-link to="/login" class="btn btn-primary mt-1.5 w-full">{{ $t('authEmail.reset.signIn') }}</router-link>
    </div>

    <form v-else class="flex flex-col gap-3.5" @submit.prevent="submit">
      <label class="flex flex-col gap-1.5">
        <span class="text-[13px] text-y-muted">{{ $t('authEmail.reset.password') }}</span>
        <input
          v-model="password"
          type="password"
          :placeholder="$t('auth2.passwordNewPlaceholder')"
          class="input input-bordered w-full"
          autocomplete="new-password"
          required
          minlength="6"
        />
      </label>

      <label class="flex flex-col gap-1.5">
        <span class="text-[13px] text-y-muted">{{ $t('authEmail.reset.confirm') }}</span>
        <input
          v-model="confirm"
          type="password"
          :placeholder="$t('authEmail.reset.confirmPlaceholder')"
          class="input input-bordered w-full"
          autocomplete="new-password"
          required
          minlength="6"
        />
      </label>

      <div v-if="error" class="alert alert-error py-3 text-sm" role="alert">
        <CircleX :size="18" :stroke-width="1.5" class="shrink-0" />
        <span>{{ error }}</span>
      </div>

      <button type="submit" class="btn btn-primary mt-1.5 w-full" :disabled="loading">
        <span v-if="loading" class="loading loading-spinner loading-sm"></span>
        <span v-else>{{ $t('authEmail.reset.submit') }}</span>
      </button>
    </form>

    <template #footer>
      <p v-if="state === 'form'" class="text-sm text-y-muted">
        {{ $t('authEmail.forgot.remembered') }}
        <router-link to="/login" class="font-medium text-y-mint-text hover:underline">{{ $t('auth2.signInLink') }}</router-link>
      </p>
      <router-link v-else-if="state !== 'done'" to="/login" class="text-sm font-medium text-y-mint-text hover:underline">
        {{ $t('authEmail.forgot.backToSignIn') }}
      </router-link>
      <AuthLegal />
    </template>
  </AuthCard>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { CircleX, CircleCheck, TriangleAlert } from 'lucide-vue-next'
import axios from '../utils/axios'
import AuthCard from '../components/ui/AuthCard.vue'
import AuthLegal from '../components/auth/AuthLegal.vue'

const route = useRoute()
const { t } = useI18n()

const token = typeof route.query.token === 'string' ? route.query.token.trim() : ''

const state = ref(token ? 'form' : 'missing')
const password = ref('')
const confirm = ref('')
const loading = ref(false)
const error = ref('')

const title = computed(() => {
  if (state.value === 'missing') return t('authEmail.reset.missingTitle')
  if (state.value === 'invalid') return t('authEmail.reset.invalidTitle')
  if (state.value === 'done') return t('authEmail.reset.doneTitle')
  return t('authEmail.reset.title')
})

const description = computed(() => (state.value === 'form' ? t('authEmail.reset.description') : ''))

const submit = async () => {
  if (loading.value) return
  error.value = ''
  if (password.value.length < 6) {
    error.value = t('authEmail.reset.tooShort')
    return
  }
  if (password.value !== confirm.value) {
    error.value = t('authEmail.reset.mismatch')
    return
  }

  loading.value = true
  try {
    await axios.post('/api/auth/reset-password', { token, password: password.value })
    state.value = 'done'
  } catch (err) {
    const data = err.response?.data
    if (data?.code === 'invalid_token') {
      state.value = 'invalid'
    } else {
      error.value = data?.error || t('authEmail.reset.failed')
    }
  } finally {
    loading.value = false
  }
}
</script>
