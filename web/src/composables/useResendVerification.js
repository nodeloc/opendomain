import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import axios from '../utils/axios'
import { useAuthStore } from '../stores/auth'
import { useToast } from './useToast'

// Module-level so the banner and the verify page share one cooldown.
const sending = ref(false)
const cooldownUntil = ref(0)
const now = ref(Date.now())
let ticker = null

const startTicker = () => {
  if (ticker) return
  ticker = setInterval(() => {
    now.value = Date.now()
    if (now.value >= cooldownUntil.value) {
      clearInterval(ticker)
      ticker = null
    }
  }, 1000)
}

export function useResendVerification() {
  const { t } = useI18n()
  const toast = useToast()
  const authStore = useAuthStore()

  const secondsLeft = computed(() => Math.max(0, Math.ceil((cooldownUntil.value - now.value) / 1000)))
  const disabled = computed(() => sending.value || secondsLeft.value > 0)

  const resend = async () => {
    if (disabled.value) return
    sending.value = true
    const email = authStore.user?.email
    try {
      await axios.post('/api/user/resend-verification')
      toast.success(email ? t('authEmail.resend.sent', { email }) : t('authEmail.resend.sentNoEmail'), 5000)
      return 'sent'
    } catch (error) {
      const data = error.response?.data || {}
      if (data.code === 'already_verified') {
        toast.info(t('authEmail.resend.alreadyVerified'))
        await authStore.fetchProfile()
        return 'already_verified'
      }
      if (error.response?.status === 429) {
        const wait = Number(data.retry_after) > 0 ? Number(data.retry_after) : 60
        now.value = Date.now()
        cooldownUntil.value = now.value + wait * 1000
        startTicker()
        toast.error(data.error || t('authEmail.resend.failed'))
        return 'rate_limited'
      }
      toast.error(data.error || t('authEmail.resend.failed'))
      return 'error'
    } finally {
      sending.value = false
    }
  }

  return { resend, sending, secondsLeft, disabled }
}
