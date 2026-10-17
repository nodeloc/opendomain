<template>
  <div class="y-container flex flex-wrap items-start gap-x-12 gap-y-12 py-14 sm:py-16">
    <!-- Main column -->
    <div class="flex min-w-0 flex-1 flex-col gap-7" style="flex-basis: 560px;">
      <div class="flex flex-col gap-3.5">
        <p class="font-mono text-eyebrow uppercase tracking-[0.1em] text-y-danger-text">{{ $t('report.eyebrow') }}</p>
        <h1 class="y-headline text-[clamp(2.5rem,5.5vw,4rem)]">{{ $t('report.title') }}</h1>
        <i18n-t keypath="report.lead" tag="p" class="m-0 max-w-[600px] text-lead text-y-muted">
          <template #policy>
            <router-link to="/acceptable-use" class="text-y-mint-text underline decoration-y-mint-text/40 underline-offset-[3px] hover:text-y-mint-hover">{{ $t('report.policyLink') }}</router-link>
          </template>
        </i18n-t>
      </div>

      <!-- Success state -->
      <div v-if="sent" role="status" class="y-mint-soft flex flex-col gap-3.5 rounded-xl border p-8">
        <span class="flex h-12 w-12 items-center justify-center rounded-full bg-y-mint" style="box-shadow: 0 0 24px rgb(var(--y-mint) / 0.5);" aria-hidden="true">
          <Check :size="24" :stroke-width="2.5" class="text-y-on-mint" />
        </span>
        <h2 class="m-0 text-[28px] font-extrabold leading-[34px] tracking-[-0.03em] text-y-text">{{ $t('report.successTitle') }}</h2>
        <i18n-t keypath="report.successBody" tag="p" class="m-0 text-y-muted">
          <template #id><span class="font-mono text-y-text">{{ reportId }}</span></template>
          <template #name><span class="font-mono text-y-text">{{ submittedName }}</span></template>
        </i18n-t>
        <div>
          <button type="button" class="btn btn-outline" @click="resetForm">{{ $t('report.another') }}</button>
        </div>
      </div>

      <!-- Form -->
      <form v-else class="flex flex-col gap-6" novalidate @submit.prevent="submit">
        <div class="y-field">
          <label for="report-sub" class="y-field-label">
            {{ $t('report.subdomainLabel') }} <span class="text-y-danger-text" aria-hidden="true">*</span>
          </label>
          <div class="flex max-w-[520px] items-stretch overflow-hidden rounded-md border border-y-line-strong bg-y-bg font-mono text-base transition-colors focus-within:border-y-mint/70">
            <input
              id="report-sub"
              v-model="subdomain"
              type="text"
              :placeholder="$t('report.subdomainPlaceholder')"
              class="min-h-[48px] min-w-0 flex-1 bg-transparent px-4 text-right text-y-text outline-none placeholder:text-y-faint"
              autocomplete="off"
              autocapitalize="off"
              spellcheck="false"
            />
            <span v-if="rootDomains.length" class="relative flex shrink-0 items-center border-l border-y-line-strong">
              <select
                v-model="selectedRootDomain"
                class="h-full max-w-[11rem] cursor-pointer appearance-none bg-transparent pl-3 pr-8 text-[15px] font-medium text-y-muted outline-none sm:max-w-none"
                :aria-label="$t('report.suffixLabel')"
              >
                <option v-for="domain in rootDomains" :key="domain.id" :value="domain.id" class="bg-y-surface">.{{ domain.domain }}</option>
              </select>
              <ChevronDown :size="16" :stroke-width="1.5" class="pointer-events-none absolute right-3 text-y-faint" />
            </span>
            <span v-else class="flex items-center pr-4 text-y-faint">.…</span>
          </div>
        </div>

        <fieldset class="m-0 flex flex-col gap-2.5 border-0 p-0">
          <legend class="float-left w-full p-0 pb-2.5 text-sm font-medium text-y-text">
            {{ $t('report.categoryLegend') }} <span class="text-y-danger-text" aria-hidden="true">*</span>
          </legend>
          <div class="grid gap-2.5 sm:grid-cols-2 lg:grid-cols-3">
            <button
              v-for="cat in categories"
              :key="cat.id"
              type="button"
              class="y-choice min-h-[76px]"
              :aria-pressed="category === cat.id"
              @click="category = cat.id; error = ''"
            >
              <span class="font-semibold text-y-text">{{ cat.label }}</span>
              <span class="text-[13px] leading-[18px] text-y-muted">{{ cat.hint }}</span>
            </button>
          </div>
        </fieldset>

        <div class="y-field">
          <label for="report-url" class="y-field-label">
            {{ $t('report.urlLabel') }} <span class="font-normal text-y-faint">{{ $t('report.optional') }}</span>
          </label>
          <input
            id="report-url"
            v-model="url"
            type="url"
            placeholder="https://suspicious-site.example.com/login"
            class="input input-bordered w-full font-mono text-[15px]"
            autocomplete="off"
            spellcheck="false"
          />
          <p class="y-field-hint">{{ $t('report.urlHint') }}</p>
        </div>

        <div class="y-field">
          <label for="report-details" class="y-field-label">
            {{ $t('report.detailsLabel') }} <span class="text-y-danger-text" aria-hidden="true">*</span>
          </label>
          <textarea
            id="report-details"
            v-model="details"
            rows="5"
            :maxlength="5000"
            :placeholder="$t('report.detailsPlaceholder')"
            class="textarea textarea-bordered w-full text-[15px]"
          ></textarea>
          <p class="y-field-hint self-end font-mono">{{ $t('report.detailsCount', { n: details.length, max: 5000 }) }}</p>
        </div>

        <div class="y-field max-w-[520px]">
          <label for="report-email" class="y-field-label">
            {{ $t('report.emailLabel') }} <span class="font-normal text-y-faint">{{ $t('report.emailOptional') }}</span>
          </label>
          <input
            id="report-email"
            v-model="email"
            type="email"
            placeholder="you@example.com"
            autocomplete="email"
            class="input input-bordered w-full"
          />
        </div>

        <div v-if="error" role="alert" class="flex items-start gap-2.5 text-sm text-y-danger-text">
          <CircleAlert :size="16" :stroke-width="1.5" class="mt-0.5 shrink-0" />
          <span>{{ error }}</span>
        </div>

        <div class="flex flex-wrap items-center gap-4">
          <button type="submit" class="y-btn-ink min-h-[48px] rounded-md px-6" :disabled="sending">
            <span v-if="sending" class="loading loading-spinner loading-sm"></span>
            {{ sending ? $t('report.sending') : $t('report.submit') }}
          </button>
          <span class="max-w-[320px] text-[13px] leading-5 text-y-faint">{{ $t('report.confidential') }}</span>
        </div>
      </form>
    </div>

    <!-- Side panel -->
    <aside class="flex w-full flex-col gap-4 lg:w-[360px] lg:shrink-0">
      <div class="y-card flex flex-col gap-[18px] p-6">
        <div class="y-eyebrow">{{ $t('report.nextEyebrow') }}</div>
        <ol class="m-0 flex list-none flex-col gap-[18px] p-0">
          <li v-for="(step, i) in nextSteps" :key="i" class="flex gap-3.5">
            <span :class="i === 0 ? 'y-step-num' : 'y-step-num-muted'">{{ i + 1 }}</span>
            <div>
              <div class="font-semibold text-y-text">{{ step.title }}</div>
              <div class="mt-0.5 text-sm leading-[22px] text-y-muted">{{ step.desc }}</div>
            </div>
          </li>
        </ol>
      </div>
      <div class="y-danger-soft rounded-xl border px-6 py-5 text-sm leading-[22px] text-y-text">
        <strong class="font-semibold text-y-danger-text">{{ $t('report.urgentTitle') }}</strong>
        {{ ' ' }}
        <i18n-t keypath="report.urgentBody" tag="span">
          <template #email>
            <a :href="`mailto:${abuseEmail}`" class="font-medium text-y-mint-text hover:text-y-mint-hover">{{ abuseEmail }}</a>
          </template>
        </i18n-t>
      </div>
    </aside>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { Check, ChevronDown, CircleAlert } from 'lucide-vue-next'
import axios from '../utils/axios'
import { ABUSE_EMAIL } from '../config/links'

const route = useRoute()
const { t } = useI18n()

const abuseEmail = ABUSE_EMAIL

const rootDomains = ref([])
const selectedRootDomain = ref('')
const subdomain = ref('')
const category = ref('')
const url = ref('')
const details = ref('')
const email = ref('')
const sending = ref(false)
const sent = ref(false)
const reportId = ref(null)
const submittedName = ref('')
const error = ref('')

const categories = computed(() => [
  { id: 'phishing', label: t('report.catPhishing'), hint: t('report.catPhishingHint') },
  { id: 'malware', label: t('report.catMalware'), hint: t('report.catMalwareHint') },
  { id: 'scam', label: t('report.catScam'), hint: t('report.catScamHint') },
  { id: 'spam', label: t('report.catSpam'), hint: t('report.catSpamHint') },
  { id: 'impersonation', label: t('report.catImpersonation'), hint: t('report.catImpersonationHint') },
  { id: 'other', label: t('report.catOther'), hint: t('report.catOtherHint') },
])

const nextSteps = computed(() => [
  { title: t('report.next1Title'), desc: t('report.next1Desc') },
  { title: t('report.next2Title'), desc: t('report.next2Desc') },
  { title: t('report.next3Title'), desc: t('report.next3Desc') },
])

const selectedRoot = computed(() => rootDomains.value.find((d) => d.id === selectedRootDomain.value))

const fullName = computed(() => {
  const name = subdomain.value.trim().toLowerCase()
  if (!name) return ''
  const root = selectedRoot.value?.domain
  return root ? `${name}.${root}` : name
})

onMounted(async () => {
  try {
    const response = await axios.get('/api/public/root-domains')
    rootDomains.value = response.data.root_domains || []
    if (rootDomains.value.length > 0) {
      selectedRootDomain.value = rootDomains.value[0].id
    }
  } catch (err) {
    console.error('Failed to fetch root domains:', err)
  }
  prefillFromQuery()
})

const prefillFromQuery = () => {
  const q = route.query.domain
  if (!q || typeof q !== 'string') return
  const clean = q.trim().toLowerCase().replace(/^\.+|\.+$/g, '')
  if (!clean) return
  const match = rootDomains.value.find((d) => clean.endsWith('.' + d.domain))
  if (match) {
    subdomain.value = clean.slice(0, -(match.domain.length + 1))
    selectedRootDomain.value = match.id
  } else {
    subdomain.value = clean
  }
}

const validate = () => {
  const name = subdomain.value.trim().toLowerCase()
  const root = selectedRoot.value
  if (!name) return t('report.errSubdomainRequired')
  if (!/^[a-z0-9-]+$/.test(name) || name.startsWith('-') || name.endsWith('-')) {
    return t('report.errSubdomainInvalid', { root: root?.domain || '…' })
  }
  const minLen = root?.min_length || 1
  const maxLen = root?.max_length || 63
  if (name.length < minLen || name.length > maxLen) {
    return t('report.errSubdomainInvalid', { root: root?.domain || '…' })
  }
  if (!category.value) return t('report.errCategory')
  const trimmed = details.value.trim()
  if (trimmed.length < 10) return t('report.errDetailsShort', { min: 10 })
  if (trimmed.length > 5000) return t('report.errDetailsLong', { max: 5000 })
  if (url.value.trim()) {
    try {
      const u = new URL(url.value.trim())
      if (u.protocol !== 'http:' && u.protocol !== 'https:') throw new Error('bad protocol')
    } catch {
      return t('report.errUrl')
    }
  }
  if (email.value.trim() && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.value.trim())) {
    return t('report.errEmail')
  }
  return ''
}

const submit = async () => {
  error.value = validate()
  if (error.value) return

  sending.value = true
  try {
    const body = {
      subdomain: fullName.value,
      category: category.value,
      details: details.value.trim(),
    }
    if (url.value.trim()) body.url = url.value.trim()
    if (email.value.trim()) body.reporter_email = email.value.trim()

    const response = await axios.post('/api/public/abuse-reports', body)
    reportId.value = response.data?.id ?? null
    submittedName.value = fullName.value
    sent.value = true
  } catch (err) {
    const status = err.response?.status
    if (status === 429) {
      error.value = t('report.errRateLimited')
    } else if (status === 400 && err.response?.data?.error) {
      error.value = err.response.data.error
    } else {
      error.value = t('report.errGeneric')
    }
  } finally {
    sending.value = false
  }
}

const resetForm = () => {
  subdomain.value = ''
  category.value = ''
  url.value = ''
  details.value = ''
  email.value = ''
  error.value = ''
  sent.value = false
  reportId.value = null
  submittedName.value = ''
}
</script>
