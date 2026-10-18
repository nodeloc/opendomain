<template>
  <AuthCard :title="sent ? $t('authEmail.forgot.sentTitle') : $t('authEmail.forgot.title')">
    <template #description>
      <span v-if="!sent">{{ $t('authEmail.forgot.description') }}</span>
      <i18n-t v-else keypath="authEmail.forgot.sent" tag="span" scope="global">
        <template #email><span class="break-all font-medium text-y-text">{{ sentTo }}</span></template>
      </i18n-t>
    </template>

    <div v-if="sent" class="flex flex-col gap-3.5">
      <div class="alert alert-success py-3 text-sm" role="status">
        <MailCheck :size="18" :stroke-width="1.5" class="shrink-0 text-y-mint-text" />
        <span>{{ $t('authEmail.forgot.spamHint') }}</span>
      </div>

      <div v-if="error" class="alert alert-error py-3 text-sm" role="alert">
        <CircleX :size="18" :stroke-width="1.5" class="shrink-0" />
        <span>{{ error }}</span>
      </div>

      <router-link to="/login" class="btn btn-primary mt-1.5 w-full">{{ $t('authEmail.forgot.backToSignIn') }}</router-link>
      <button type="button" class="btn btn-outline w-full" :disabled="loading" @click="send(true)">
        <span v-if="loading" class="loading loading-spinner loading-sm"></span>
        <span v-else>{{ $t('authEmail.forgot.sendAgain') }}</span>
      </button>
    </div>

    <form v-else class="flex flex-col gap-3.5" @submit.prevent="send(false)">
      <label class="flex flex-col gap-1.5">
        <span class="text-[13px] text-y-muted">{{ $t('login.email') }}</span>
        <input
          v-model="email"
          type="email"
          :placeholder="$t('auth2.emailPlaceholder')"
          class="input input-bordered w-full"
          autocomplete="email"
          required
        />
      </label>

      <div v-if="error" class="alert alert-error py-3 text-sm" role="alert">
        <CircleX :size="18" :stroke-width="1.5" class="shrink-0" />
        <span>{{ error }}</span>
      </div>

      <button type="submit" class="btn btn-primary mt-1.5 w-full" :disabled="loading">
        <span v-if="loading" class="loading loading-spinner loading-sm"></span>
        <span v-else>{{ $t('authEmail.forgot.submit') }}</span>
      </button>
    </form>

    <template #footer>
      <p v-if="!sent" class="text-sm text-y-muted">
        {{ $t('authEmail.forgot.remembered') }}
        <router-link to="/login" class="font-medium text-y-mint-text hover:underline">{{ $t('auth2.signInLink') }}</router-link>
      </p>
      <AuthLegal />
    </template>
  </AuthCard>
</template>

<script setup>
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { CircleX, MailCheck } from 'lucide-vue-next'
import axios from '../utils/axios'
import { useToast } from '../composables/useToast'
import AuthCard from '../components/ui/AuthCard.vue'
import AuthLegal from '../components/auth/AuthLegal.vue'

const route = useRoute()
const { t } = useI18n()
const toast = useToast()

const email = ref(typeof route.query.email === 'string' ? route.query.email : '')
const sentTo = ref('')
const sent = ref(false)
const loading = ref(false)
const error = ref('')

const send = async (again) => {
  const target = again ? sentTo.value : email.value.trim()
  if (!target || loading.value) return
  loading.value = true
  error.value = ''
  try {
    await axios.post('/api/auth/forgot-password', { email: target })
    sentTo.value = target
    sent.value = true
    if (again) toast.success(t('authEmail.forgot.resent'))
  } catch (err) {
    error.value = err.response?.data?.error || t('authEmail.forgot.failed')
  } finally {
    loading.value = false
  }
}
</script>
