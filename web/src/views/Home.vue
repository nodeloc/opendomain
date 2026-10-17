<template>
  <div>
    <!-- Hero -->
    <section class="relative isolate overflow-hidden border-b border-y-line">
      <div class="y-dot-grid absolute inset-0 -z-20" aria-hidden="true"></div>
      <img
        v-if="isDark"
        src="/brand/planet-horizon-2880x1800.jpg"
        alt=""
        aria-hidden="true"
        class="pointer-events-none absolute inset-x-0 bottom-0 -z-10 h-[62%] w-full object-cover object-bottom opacity-90"
      />
      <div v-else class="planet-light pointer-events-none absolute left-1/2 top-[76%] -z-10 aspect-square w-[200%] -translate-x-1/2 rounded-full sm:w-[160%]" aria-hidden="true"></div>

      <div class="y-container pb-44 pt-16 sm:pb-60 sm:pt-24">
        <p class="y-eyebrow">{{ $t('landing.eyebrow') }}</p>
        <h1 class="y-headline mt-5 text-[clamp(3rem,9vw,7rem)]">
          {{ $t('landing.headlineLead') }}<br />
          {{ $t('landing.headlineMid') }} <em>{{ $t('landing.headlineAccent') }}</em>
        </h1>
        <p class="y-lead mt-6 max-w-xl">{{ $t('landing.lead') }}</p>

        <!-- Search -->
        <div class="mt-10 max-w-3xl">
          <form class="y-panel flex flex-col gap-2 p-2 sm:flex-row" @submit.prevent="searchDomain">
            <label class="flex h-12 min-w-0 flex-1 items-center rounded-md border border-y-line bg-y-bg transition-colors focus-within:border-y-mint/70">
              <Search :size="18" :stroke-width="1.5" class="ml-4 shrink-0 text-y-faint" />
              <input
                v-model="searchQuery"
                type="text"
                :placeholder="$t('landing.placeholder')"
                class="h-full min-w-0 flex-1 bg-transparent px-3 font-mono text-[15px] text-y-text outline-none placeholder:text-y-faint"
                autocomplete="off"
                autocapitalize="off"
                spellcheck="false"
                required
              />
              <span class="relative flex h-full shrink-0 items-center border-l border-y-line">
                <select
                  v-model="selectedRootDomain"
                  class="h-full max-w-[11rem] appearance-none bg-transparent pl-3 pr-8 font-mono text-[15px] text-y-muted outline-none sm:max-w-none"
                  :aria-label="$t('home.chooseExtension')"
                >
                  <option value="" disabled>{{ $t('home.chooseExtension') }}</option>
                  <option v-for="domain in rootDomains" :key="domain.id" :value="domain.id">.{{ domain.domain }}</option>
                </select>
                <ChevronDown :size="16" :stroke-width="1.5" class="pointer-events-none absolute right-3 text-y-faint" />
              </span>
            </label>
            <button type="submit" class="btn btn-primary h-12 min-h-12 px-6" :disabled="loading">
              <span v-if="loading" class="loading loading-spinner loading-sm"></span>
              {{ loading ? $t('landing.checking') : $t('landing.check') }}
            </button>
          </form>

          <transition name="result">
            <div v-if="searchResult" class="mt-3">
              <!-- Available -->
              <div
                v-if="searchResult.available && !searchResult.levelTooLow"
                class="y-mint-soft flex flex-col gap-3 rounded-lg border px-5 py-4 sm:flex-row sm:items-center"
              >
                <div class="min-w-0 flex-1">
                  <p class="flex flex-wrap items-center gap-x-2 gap-y-1">
                    <span class="y-dot bg-y-mint"></span>
                    <span class="y-domain">{{ searchResult.fullDomain }}</span>
                    <span class="text-sm text-y-muted">{{ $t('landing.isAvailable') }}</span>
                  </p>
                  <p class="mt-1 pl-4 text-sm text-y-muted">
                    {{ searchResult.isFree ? $t('landing.freeNote') : $t('landing.priceNote', { price: `${currencySymbol}${searchResult.pricePerYear || searchResult.lifetimePrice}` }) }}
                  </p>
                </div>
                <button v-if="!isAuthenticated" type="button" class="btn btn-primary btn-sm gap-1.5" @click="$router.push('/register')">
                  {{ $t('landing.signUpToClaim') }} <ArrowRight :size="16" :stroke-width="1.5" />
                </button>
                <button v-else type="button" class="btn btn-primary btn-sm gap-1.5" @click="registerDomain">
                  {{ searchResult.isFree ? $t('landing.claimIt') : $t('landing.checkout') }} <ArrowRight :size="16" :stroke-width="1.5" />
                </button>
              </div>

              <!-- Level too low -->
              <div v-else-if="searchResult.levelTooLow" class="y-panel flex items-start gap-3 px-5 py-4">
                <Info :size="18" :stroke-width="1.5" class="mt-0.5 shrink-0 text-y-info" />
                <div class="text-sm">
                  <p class="font-medium text-y-text">{{ $t('home.levelTooLow') }}</p>
                  <p class="mt-1 text-y-muted">
                    {{ $t('home.requiredLevel') }}: <span class="text-y-text">{{ $t(`home.levelNames.${searchResult.requiredLevel}`) }}</span>
                    <span class="mx-2 text-y-faint">·</span>
                    {{ $t('home.yourLevel') }}: <span class="text-y-text">{{ $t(`home.levelNames.${searchResult.yourLevel}`) }}</span>
                  </p>
                </div>
              </div>

              <!-- Taken or error -->
              <div v-else class="y-panel flex flex-col gap-3 px-5 py-4 sm:flex-row sm:items-center">
                <div class="min-w-0 flex-1">
                  <p class="flex flex-wrap items-center gap-x-2 gap-y-1">
                    <span class="y-dot bg-y-danger"></span>
                    <template v-if="searchResult.fullDomain">
                      <span class="y-domain">{{ searchResult.fullDomain }}</span>
                      <span class="text-sm text-y-danger-text">{{ $t('landing.isTaken') }}</span>
                    </template>
                    <span v-else class="text-sm text-y-danger-text">{{ searchResult.message }}</span>
                  </p>
                  <p class="mt-1 pl-4 text-sm text-y-muted">{{ $t('landing.tryAnother') }}</p>
                </div>
                <button
                  v-if="searchResult.domainInfo"
                  type="button"
                  class="btn btn-outline btn-sm"
                  @click="$router.push(`/whois/${searchResult.fullDomain}`)"
                >
                  {{ $t('landing.viewWhois') }}
                </button>
              </div>
            </div>
          </transition>

          <div class="mt-5 flex flex-wrap gap-2">
            <span class="y-chip"><Check :size="12" :stroke-width="2" class="text-y-mint-text" />{{ $t('landing.chipNoCard') }}</span>
            <span class="y-chip"><Check :size="12" :stroke-width="2" class="text-y-mint-text" />{{ $t('landing.chipMinutes') }}</span>
            <span class="y-chip"><Check :size="12" :stroke-width="2" class="text-y-mint-text" />{{ $t('landing.chipNoLock') }}</span>
          </div>
        </div>
      </div>
    </section>

    <!-- How it works -->
    <section class="y-container py-24">
      <div class="grid gap-12 lg:grid-cols-[1fr_1.1fr] lg:items-center">
        <div>
          <p class="y-eyebrow">{{ $t('landing.howEyebrow') }}</p>
          <h2 class="y-h1 mt-4 max-w-md">{{ $t('landing.howTitle') }}</h2>
          <ol class="mt-10 space-y-6">
            <li v-for="(step, i) in steps" :key="step.title" class="flex gap-5">
              <span class="font-mono text-sm text-y-mint-text">0{{ i + 1 }}</span>
              <div>
                <h3 class="y-h3">{{ step.title }}</h3>
                <p class="mt-1 text-[15px] leading-6 text-y-muted">{{ step.desc }}</p>
              </div>
            </li>
          </ol>
        </div>

        <div class="y-tile overflow-hidden">
          <div class="relative border-b border-y-line px-6 py-8">
            <svg viewBox="0 0 420 150" class="h-auto w-full" aria-hidden="true">
              <path
                d="M70 75 H150 C190 75 190 25 230 25 H300 M150 75 H300 M150 75 C190 75 190 125 230 125 H300"
                fill="none"
                stroke="rgb(var(--y-text))"
                stroke-width="2"
                stroke-linecap="round"
              />
              <circle cx="70" cy="75" r="6" fill="rgb(var(--y-text))" />
              <circle cx="302" cy="25" r="6" fill="rgb(var(--y-node))" />
              <circle cx="302" cy="75" r="6" fill="rgb(var(--y-node))" />
              <circle cx="302" cy="125" r="6" fill="rgb(var(--y-node))" />
              <text x="70" y="104" text-anchor="middle" class="node-label">{{ rootLabel }}</text>
              <text x="318" y="29" class="node-label node-label-strong">blog.</text>
              <text x="318" y="79" class="node-label node-label-strong">docs.</text>
              <text x="318" y="129" class="node-label node-label-strong">api.</text>
            </svg>
          </div>
          <table class="y-table font-mono text-[13px]">
            <thead>
              <tr>
                <th>Name</th>
                <th>Type</th>
                <th>Content</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="row in sampleRecords" :key="row.name">
                <td class="text-y-text">{{ row.name }}</td>
                <td><span class="y-tag">{{ row.type }}</span></td>
                <td class="text-y-muted">{{ row.content }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </section>

    <!-- Features -->
    <section class="border-y border-y-line bg-y-surface/40">
      <div class="y-container py-24">
        <p class="y-eyebrow">{{ $t('landing.featuresEyebrow') }}</p>
        <h2 class="y-h1 mt-4 max-w-xl">{{ $t('landing.featuresTitle') }}</h2>
        <div class="mt-12 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <div v-for="feature in features" :key="feature.title" class="y-tile p-6">
            <span class="y-icon-tile"><component :is="feature.icon" :size="20" :stroke-width="1.5" /></span>
            <h3 class="mt-6 text-[17px] font-semibold tracking-[-0.01em] text-y-text">{{ feature.title }}</h3>
            <p class="mt-2 text-sm leading-6 text-y-muted">{{ feature.desc }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Stats -->
    <section class="y-container py-16">
      <div class="y-panel grid divide-y divide-y-line sm:grid-cols-3 sm:divide-x sm:divide-y-0">
        <div v-for="stat in stats" :key="stat.label" class="px-6 py-7">
          <p class="y-label">{{ stat.label }}</p>
          <p class="mt-3 font-mono text-4xl font-medium tracking-tight text-y-text">{{ stat.value }}</p>
        </div>
      </div>
    </section>

    <!-- CTA -->
    <section class="y-container pb-8 pt-8">
      <div class="y-tile relative isolate overflow-hidden px-6 py-20 text-center">
        <div class="y-glow absolute left-1/2 top-1/2 -z-10 h-[420px] w-[720px] -translate-x-1/2 -translate-y-1/2" aria-hidden="true"></div>
        <div class="y-dot-grid absolute inset-0 -z-20" aria-hidden="true"></div>
        <h2 class="y-headline text-[clamp(2.25rem,6vw,4.5rem)]">{{ $t('landing.ctaTitle') }}</h2>
        <p class="y-lead mx-auto mt-5 max-w-md">{{ $t('landing.ctaDesc') }}</p>
        <div class="mt-8 flex flex-wrap justify-center gap-3">
          <router-link :to="isAuthenticated ? '/domains' : '/register'" class="btn btn-primary gap-1.5 px-6">
            {{ $t('landing.ctaButton') }} <ArrowRight :size="16" :stroke-width="1.5" />
          </router-link>
          <router-link to="/announcements" class="btn btn-outline px-6">{{ $t('landing.ctaSecondary') }}</router-link>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { Search, ChevronDown, ArrowRight, Check, Info, Zap, ListTree, Network, Activity, KeyRound, Sparkles } from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'
import { useSiteConfigStore } from '../stores/siteConfig'
import axios from '../utils/axios'
import { useToast } from '../composables/useToast'
import { useCurrency } from '../composables/useCurrency'
import { useTheme } from '../composables/useTheme'

const router = useRouter()
const { t } = useI18n()
const authStore = useAuthStore()
const siteConfigStore = useSiteConfigStore()
const toast = useToast()
const { currencySymbol } = useCurrency()
const { isDark } = useTheme()

const searchQuery = ref('')
const selectedRootDomain = ref('')
const rootDomains = ref([])
const searchResult = ref(null)
const loading = ref(false)

const totalDomains = ref(null)
const totalUsers = ref(null)

const isAuthenticated = computed(() => authStore.isAuthenticated)

const rootLabel = computed(() => rootDomains.value[0]?.domain || 'yunamo.com')

const sampleRecords = computed(() => {
  const root = rootLabel.value
  return [
    { name: `blog.${root}`, type: 'CNAME', content: 'my-site.pages.dev' },
    { name: `docs.${root}`, type: 'CNAME', content: 'docs.vercel.app' },
    { name: `api.${root}`, type: 'A', content: '203.0.113.10' },
  ]
})

const steps = computed(() => [
  { title: t('landing.step1Title'), desc: t('landing.step1Desc') },
  { title: t('landing.step2Title'), desc: t('landing.step2Desc') },
  { title: t('landing.step3Title'), desc: t('landing.step3Desc') },
])

const features = computed(() => [
  { icon: Zap, title: t('landing.fFastTitle'), desc: t('landing.fFastDesc') },
  { icon: ListTree, title: t('landing.fRecordsTitle'), desc: t('landing.fRecordsDesc') },
  { icon: Network, title: t('landing.fRootsTitle'), desc: t('landing.fRootsDesc') },
  { icon: Activity, title: t('landing.fHealthTitle'), desc: t('landing.fHealthDesc') },
  { icon: KeyRound, title: t('landing.fSecureTitle'), desc: t('landing.fSecureDesc') },
  { icon: Sparkles, title: t('landing.fFreeTitle'), desc: t('landing.fFreeDesc') },
])

const formatCount = (value) => (value === null || value === undefined ? '—' : Number(value).toLocaleString())

const stats = computed(() => [
  { label: t('landing.statDomains'), value: formatCount(totalDomains.value) },
  { label: t('landing.statUsers'), value: formatCount(totalUsers.value) },
  { label: t('landing.statRoots'), value: rootDomains.value.length ? rootDomains.value.length : '—' },
])

onMounted(async () => {
  await fetchRootDomains()
  await fetchPublicStats()
})

const fetchPublicStats = async () => {
  try {
    const response = await axios.get('/api/public/stats')
    totalDomains.value = response.data.total_domains
    totalUsers.value = response.data.total_users
  } catch (error) {
    console.error('Failed to fetch stats:', error)
  }
}

const fetchRootDomains = async () => {
  try {
    const response = await axios.get('/api/public/root-domains')
    rootDomains.value = response.data.root_domains || []
    if (rootDomains.value.length > 0) {
      selectedRootDomain.value = rootDomains.value[0].id
    }
  } catch (error) {
    console.error('Failed to fetch root domains:', error)
  }
}

const levelOrder = (level) => {
  const map = { normal: 0, basic: 1, member: 2, regular: 3, leader: 4 }
  return map[level] ?? 0
}

const searchDomain = async () => {
  if (!searchQuery.value || !selectedRootDomain.value) return

  loading.value = true
  searchResult.value = null

  try {
    const response = await axios.get('/api/domains/search', {
      params: {
        subdomain: searchQuery.value,
        root_domain_id: selectedRootDomain.value,
      },
    })

    const rootDomain = rootDomains.value.find((d) => d.id === selectedRootDomain.value)

    const minLevel = rootDomain?.min_user_level || 'normal'
    const userLevel = authStore.user?.user_level || 'normal'
    const tooLow = isAuthenticated.value && levelOrder(userLevel) < levelOrder(minLevel)

    searchResult.value = {
      available: response.data.available,
      levelTooLow: tooLow,
      requiredLevel: minLevel,
      yourLevel: userLevel,
      message: response.data.available
        ? `${response.data.full_domain} is available!`
        : `${response.data.full_domain} is already taken.`,
      fullDomain: response.data.full_domain,
      rootDomain: rootDomain,
      isFree: rootDomain?.is_free ?? true,
      pricePerYear: rootDomain?.price_per_year,
      lifetimePrice: rootDomain?.lifetime_price,
      domainInfo: response.data.domain_info || null,
    }
  } catch (error) {
    searchResult.value = {
      available: false,
      message: error.response?.data?.error || 'Search failed',
    }
  } finally {
    loading.value = false
  }
}

const registerDomain = async () => {
  if (!isAuthenticated.value) {
    router.push('/register')
    return
  }

  if (searchResult.value && !searchResult.value.isFree) {
    router.push({
      path: '/checkout',
      query: {
        subdomain: searchQuery.value,
        root_domain_id: selectedRootDomain.value,
      },
    })
    return
  }

  try {
    await axios.post('/api/domains', {
      subdomain: searchQuery.value,
      root_domain_id: selectedRootDomain.value,
    })

    toast.success(t('landing.registered'))
    router.push('/domains')
  } catch (error) {
    const errorData = error.response?.data
    if (errorData?.requires_payment) {
      router.push({
        path: '/checkout',
        query: {
          subdomain: searchQuery.value,
          root_domain_id: selectedRootDomain.value,
        },
      })
    } else if (error.response?.status === 403 && errorData?.required_level) {
      searchResult.value = {
        ...searchResult.value,
        levelTooLow: true,
        requiredLevel: errorData.required_level,
        yourLevel: errorData.your_level || authStore.user?.user_level || 'normal',
      }
    } else {
      toast.error(errorData?.error || 'Registration failed')
    }
  }
}
</script>

<style scoped>
.planet-light {
  background: radial-gradient(closest-side, rgb(var(--y-surface)) 97%, transparent 100%);
  box-shadow: 0 -1px 0 0 rgb(var(--y-mint-deep) / 0.8), 0 -30px 90px -50px rgb(var(--y-mint-deep) / 0.35);
}
.node-label {
  font-family: 'Geist Mono', ui-monospace, monospace;
  font-size: 12px;
  fill: rgb(var(--y-faint));
}
.node-label-strong {
  fill: rgb(var(--y-text));
}
.result-enter-active {
  transition: all 0.25s cubic-bezier(0.2, 0.8, 0.2, 1);
}
.result-leave-active {
  transition: all 0.15s ease;
}
.result-enter-from,
.result-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}
</style>
