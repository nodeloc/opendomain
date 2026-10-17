<template>
  <div>
    <!-- Hero -->
    <section data-theme="yunamo" class="relative isolate flex min-h-[720px] items-center overflow-hidden bg-y-bg text-y-text sm:min-h-[840px]">
      <img
        src="/brand/planet-horizon-2880x1800.jpg"
        alt=""
        aria-hidden="true"
        class="pointer-events-none absolute inset-0 h-full w-full object-cover object-[50%_80%]"
      />
      <div
        aria-hidden="true"
        class="pointer-events-none absolute inset-0"
        style="background: linear-gradient(90deg, rgb(10 11 13 / 0.35) 0%, rgb(10 11 13 / 0) 50%, rgb(10 11 13 / 0.2) 100%);"
      ></div>

      <div class="y-container relative flex flex-wrap items-center gap-x-16 gap-y-12 py-20 sm:py-24">
        <div class="flex min-w-0 flex-1 flex-col gap-6" style="flex-basis: 420px;">
          <h1 class="y-headline text-[clamp(2.75rem,6.4vw,5.25rem)]" style="text-shadow: 0 2px 30px rgb(0 0 0 / 0.35);">
            {{ $t('home2.titleLead') }}<br />
            {{ $t('home2.titleMid') }} <em>{{ $t('home2.titleAccent') }}</em>
          </h1>
          <p class="max-w-[500px] text-lead text-y-text/85">
            {{ $t('home2.lead', { example: `yourname.${rootLabel}` }) }}
          </p>
        </div>

        <!-- Availability form -->
        <div class="flex min-w-0 flex-1 flex-col gap-4" style="flex-basis: 420px; max-width: 560px;">
          <form class="flex flex-col gap-3.5" @submit.prevent="searchDomain">
            <label for="home-sub" class="font-mono text-[12px] font-semibold uppercase leading-4 tracking-[0.1em] text-y-muted">
              {{ $t('home2.checkLabel') }}
            </label>
            <div class="flex items-stretch overflow-hidden rounded-md border border-y-line-strong bg-y-bg/80 backdrop-blur transition-colors focus-within:border-y-mint">
              <input
                id="home-sub"
                v-model="searchQuery"
                type="text"
                :placeholder="$t('home2.placeholder')"
                class="h-[52px] min-w-0 flex-1 bg-transparent px-4 font-mono text-base text-y-text outline-none placeholder:text-y-faint"
                autocomplete="off"
                autocapitalize="off"
                spellcheck="false"
                @input="onQueryInput"
              />
              <span class="relative flex shrink-0 items-center border-l border-y-line-strong">
                <select
                  v-model="selectedRootDomain"
                  class="h-full max-w-[11rem] cursor-pointer appearance-none bg-transparent pl-3.5 pr-8 font-mono text-[15px] font-medium text-y-text outline-none sm:max-w-none"
                  :aria-label="$t('home2.suffixLabel')"
                >
                  <option value="" disabled>{{ rootDomains.length ? $t('home2.suffixLabel') : $t('home2.noSuffixes') }}</option>
                  <option v-for="domain in rootDomains" :key="domain.id" :value="domain.id" class="bg-y-surface">.{{ domain.domain }}</option>
                </select>
                <ChevronDown :size="16" :stroke-width="1.5" class="pointer-events-none absolute right-3 text-y-faint" />
              </span>
            </div>
            <button type="submit" class="y-btn-ink-block gap-2.5 font-mono uppercase tracking-[0.12em]" :disabled="loading || !rootDomains.length">
              <span v-if="loading" class="loading loading-spinner loading-sm"></span>
              {{ loading ? $t('home2.checking') : $t('home2.checkButton') }}
              <ArrowRight v-if="!loading" :size="16" :stroke-width="2" aria-hidden="true" />
            </button>

            <transition name="result">
              <div v-if="searchResult" class="flex flex-col gap-3">
                <!-- Available -->
                <div
                  v-if="searchResult.available && !searchResult.levelTooLow"
                  class="flex flex-wrap items-center gap-x-4 gap-y-3 rounded-md border bg-y-bg/85 py-3 pl-4 pr-3"
                  style="border-color: rgb(62 230 180 / 0.45);"
                >
                  <span class="y-dot-live" aria-hidden="true"></span>
                  <span class="min-w-0 flex-1 font-mono font-medium text-y-text">
                    {{ searchResult.fullDomain }} <span class="font-sans text-y-mint-text">{{ $t('home2.isAvailable') }}</span>
                  </span>
                  <p v-if="!searchResult.isFree" class="w-full pl-5 text-sm text-y-muted">
                    {{ $t('landing.priceNote', { price: `${currencySymbol}${searchResult.pricePerYear || searchResult.lifetimePrice}` }) }}
                  </p>
                  <button v-if="!isAuthenticated" type="button" class="btn btn-primary btn-sm gap-1.5" @click="$router.push('/register')">
                    {{ $t('landing.signUpToClaim') }} <ArrowRight :size="16" :stroke-width="1.5" />
                  </button>
                  <button v-else type="button" class="btn btn-primary btn-sm gap-1.5" @click="registerDomain">
                    {{ searchResult.isFree ? $t('home2.claimIt') : $t('landing.checkout') }} <ArrowRight :size="16" :stroke-width="1.5" />
                  </button>
                </div>

                <!-- Level too low -->
                <div v-else-if="searchResult.levelTooLow" class="flex items-start gap-3 rounded-md border border-y-line-strong bg-y-bg/85 px-4 py-3.5">
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

                <!-- Taken -->
                <div v-else-if="!searchResult.invalid" class="flex flex-col gap-2.5 rounded-md border border-y-line-strong bg-y-bg/85 px-4 py-3.5">
                  <div class="flex flex-wrap items-center gap-3">
                    <span class="y-dot bg-y-danger" aria-hidden="true"></span>
                    <span class="font-mono font-medium text-y-text">
                      {{ searchResult.fullDomain }} <span class="font-sans text-y-danger-text">{{ $t('home2.isTaken') }}</span>
                    </span>
                    <button
                      v-if="searchResult.domainInfo"
                      type="button"
                      class="ml-auto text-sm font-medium text-y-mint-text hover:text-y-mint-hover"
                      @click="$router.push(`/whois/${searchResult.fullDomain}`)"
                    >
                      {{ $t('home2.viewDetails') }}
                    </button>
                  </div>
                  <div class="flex flex-wrap items-center gap-2">
                    <span class="mr-1 text-sm text-y-muted">{{ $t('home2.tryLabel') }}</span>
                    <button
                      v-for="s in suggestions"
                      :key="s"
                      type="button"
                      class="min-h-[36px] rounded-sm border border-y-line-strong bg-y-surface-2 px-3 font-mono text-[13px] font-medium text-y-muted transition-colors hover:border-y-mint hover:text-y-text"
                      @click="pickSuggestion(s)"
                    >
                      {{ s }}.{{ searchResult.rootDomain?.domain }}
                    </button>
                  </div>
                </div>

                <!-- Invalid / error -->
                <div v-else role="alert" class="flex items-start gap-2.5 rounded-md border border-y-line-strong bg-y-bg/85 px-4 py-3.5 text-sm text-y-danger-text">
                  <CircleAlert :size="16" :stroke-width="1.5" class="mt-0.5 shrink-0" />
                  <span>{{ searchResult.message }}</span>
                </div>
              </div>
            </transition>
          </form>
        </div>
      </div>

      <a
        :href="telegramUrl"
        target="_blank"
        rel="noopener"
        class="group absolute bottom-6 right-6 z-10 inline-flex items-center gap-2.5 rounded-full bg-y-text py-2.5 pl-3 pr-4 text-sm font-bold text-y-bg transition-shadow hover:shadow-[0_0_0_4px_rgb(62_230_180_/_0.25)]"
      >
        <span class="flex h-[26px] w-[26px] items-center justify-center rounded-full bg-y-bg text-y-mint-text">
          <Send :size="14" :stroke-width="2" aria-hidden="true" />
        </span>
        {{ $t('home2.telegramHelp') }}
      </a>
    </section>

    <!-- Works with -->
    <section class="border-b border-y-line">
      <div class="y-container flex flex-wrap items-center justify-center gap-x-10 gap-y-3 py-7 text-[15px] font-medium text-y-faint">
        <span class="y-label">{{ $t('home2.worksWith') }}</span>
        <span>Vercel</span><span>Netlify</span><span>GitHub Pages</span><span>Cloudflare Pages</span><span>Render</span><span>{{ $t('home2.anyVps') }}</span>
      </div>
    </section>

    <!-- Bento features -->
    <section class="y-container flex flex-col gap-12 py-24 sm:py-28">
      <div class="flex max-w-[640px] flex-col gap-3">
        <p class="y-eyebrow">{{ $t('home2.featuresEyebrow') }}</p>
        <h2 class="y-title">{{ $t('home2.featuresTitle') }}</h2>
      </div>
      <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <!-- Full DNS control (wide) -->
        <article class="y-card flex flex-wrap items-center gap-8 overflow-hidden p-8 sm:col-span-2 lg:col-span-2">
          <div class="flex min-w-0 flex-col gap-3" style="flex: 1 1 240px;">
            <h3 class="text-h3">{{ $t('home2.fDnsTitle') }}</h3>
            <p class="text-y-muted">{{ $t('home2.fDnsDesc') }}</p>
          </div>
          <div class="y-inset min-w-0 overflow-hidden font-mono text-[13px] leading-5" style="flex: 1 1 320px;">
            <div class="flex gap-1.5 border-b border-y-line px-4 py-2.5">
              <span class="h-2.5 w-2.5 rounded-full bg-y-line-strong"></span>
              <span class="h-2.5 w-2.5 rounded-full bg-y-line-strong"></span>
              <span class="h-2.5 w-2.5 rounded-full bg-y-line-strong"></span>
            </div>
            <div class="flex flex-col gap-2 overflow-x-auto whitespace-pre p-4">
              <div><span class="text-y-mint-text">CNAME</span>  <span class="text-y-text">blog</span>  <span class="text-y-faint">→</span> <span class="text-y-muted">my-site.pages.dev</span></div>
              <div><span class="text-y-mint-text">A    </span>  <span class="text-y-text">@   </span>  <span class="text-y-faint">→</span> <span class="text-y-muted">203.0.113.10</span></div>
              <div><span class="text-y-mint-text">TXT  </span>  <span class="text-y-text">@   </span>  <span class="text-y-faint">→</span> <span class="text-y-muted">"verify=yunamo"</span></div>
              <div class="text-y-mint-text">✓ {{ $t('home2.propagated') }}</div>
            </div>
          </div>
        </article>

        <article class="y-card flex flex-col gap-3 p-8">
          <div class="text-[56px] font-extrabold leading-none tracking-[-0.05em] text-y-mint-text">$0</div>
          <h3 class="text-xl font-semibold tracking-[-0.01em] text-y-text">{{ $t('home2.fFreeTitle') }}</h3>
          <p class="text-y-muted">{{ $t('home2.fFreeDesc') }}</p>
        </article>

        <article class="y-card flex flex-col gap-3 p-8">
          <Zap :size="28" :stroke-width="1.5" class="text-y-mint-text" aria-hidden="true" />
          <h3 class="mt-2 text-xl font-semibold tracking-[-0.01em] text-y-text">{{ $t('home2.fLiveTitle') }}</h3>
          <p class="text-y-muted">{{ $t('home2.fLiveDesc') }}</p>
        </article>

        <article class="y-card flex flex-col gap-3 p-8">
          <ShieldCheck :size="28" :stroke-width="1.5" class="text-y-mint-text" aria-hidden="true" />
          <h3 class="mt-2 text-xl font-semibold tracking-[-0.01em] text-y-text">{{ $t('home2.fTrustTitle') }}</h3>
          <p class="text-y-muted">{{ $t('home2.fTrustDesc') }}</p>
        </article>

        <article class="y-card is-hoverable flex flex-col gap-3 p-8">
          <Code2 :size="28" :stroke-width="1.5" class="text-y-mint-text" aria-hidden="true" />
          <h3 class="mt-2 text-xl font-semibold tracking-[-0.01em] text-y-text">{{ $t('home2.fOpenTitle') }}</h3>
          <p class="text-y-muted">{{ $t('home2.fOpenDesc') }}</p>
          <a :href="githubUrl" target="_blank" rel="noopener" class="mt-auto inline-flex items-center gap-1.5 pt-2 text-sm font-semibold text-y-mint-text hover:text-y-mint-hover">
            {{ $t('home2.fOpenLink') }} <ArrowUpRight :size="15" :stroke-width="1.5" aria-hidden="true" />
          </a>
        </article>
      </div>
    </section>

    <!-- How it works -->
    <section class="y-band">
      <div class="y-band-inner flex flex-col gap-12">
        <div class="flex max-w-[640px] flex-col gap-3">
          <p class="y-eyebrow">{{ $t('home2.howEyebrow') }}</p>
          <h2 class="y-title">{{ $t('home2.howTitle') }}</h2>
        </div>
        <ol class="m-0 grid list-none gap-x-6 gap-y-10 p-0 sm:grid-cols-3">
          <li class="flex flex-col gap-2 border-t border-y-mint pt-6">
            <span class="font-mono text-[13px] text-y-mint-text">01</span>
            <h3 class="m-0 text-xl font-semibold leading-7 text-y-text">{{ $t('home2.step1Title') }}</h3>
            <p class="m-0 text-y-muted">{{ $t('home2.step1Desc', { suffix: '.' + rootLabel }) }}</p>
          </li>
          <li class="flex flex-col gap-2 border-t border-y-line-strong pt-6">
            <span class="font-mono text-[13px] text-y-faint">02</span>
            <h3 class="m-0 text-xl font-semibold leading-7 text-y-text">{{ $t('home2.step2Title') }}</h3>
            <p class="m-0 text-y-muted">{{ step2Desc }}</p>
          </li>
          <li class="flex flex-col gap-2 border-t border-y-line-strong pt-6">
            <span class="font-mono text-[13px] text-y-faint">03</span>
            <h3 class="m-0 text-xl font-semibold leading-7 text-y-text">{{ $t('home2.step3Title') }}</h3>
            <p class="m-0 text-y-muted">{{ $t('home2.step3Desc') }}</p>
          </li>
        </ol>
      </div>
    </section>

    <!-- FAQ -->
    <section class="y-band">
      <div class="y-band-inner flex flex-wrap gap-x-16 gap-y-10">
        <div class="flex min-w-0 flex-col gap-3" style="flex: 1 1 240px;">
          <p class="y-eyebrow">{{ $t('home2.faqEyebrow') }}</p>
          <h2 class="y-title">{{ $t('home2.faqTitle') }}</h2>
        </div>
        <div class="y-faq min-w-0 flex flex-col border-t border-y-line" style="flex: 2 1 480px;">
          <details>
            <summary>{{ $t('home2.faqFreeQ') }} <span class="y-faq-plus" aria-hidden="true">+</span></summary>
            <p class="m-0 mt-3 text-y-muted">{{ $t('home2.faqFreeA') }}</p>
          </details>
          <details>
            <summary>{{ $t('home2.faqQuotaQ') }} <span class="y-faq-plus" aria-hidden="true">+</span></summary>
            <p class="m-0 mt-3 text-y-muted">{{ $t('home2.faqQuotaA') }}</p>
          </details>
          <details>
            <summary>{{ $t('home2.faqNsQ') }} <span class="y-faq-plus" aria-hidden="true">+</span></summary>
            <p class="m-0 mt-3 text-y-muted">{{ $t('home2.faqNsA') }}</p>
          </details>
          <details>
            <summary>{{ $t('home2.faqSuffixQ') }} <span class="y-faq-plus" aria-hidden="true">+</span></summary>
            <p class="m-0 mt-3 text-y-muted">{{ $t('home2.faqSuffixA') }}</p>
          </details>
          <details>
            <summary>{{ $t('home2.faqRulesQ') }} <span class="y-faq-plus" aria-hidden="true">+</span></summary>
            <i18n-t keypath="home2.faqRulesA" tag="p" class="m-0 mt-3 text-y-muted">
              <template #policy>
                <router-link to="/acceptable-use" class="font-medium text-y-mint-text underline decoration-y-mint-text/40 underline-offset-[3px] hover:text-y-mint-hover">{{ $t('home2.faqPolicyLink') }}</router-link>
              </template>
              <template #report>
                <router-link to="/report" class="font-medium text-y-mint-text underline decoration-y-mint-text/40 underline-offset-[3px] hover:text-y-mint-hover">{{ $t('home2.faqReportLink') }}</router-link>
              </template>
            </i18n-t>
          </details>
        </div>
      </div>
    </section>

    <!-- CTA -->
    <section class="y-container pb-24">
      <div class="y-card relative isolate overflow-hidden px-6 py-16 text-center sm:py-20">
        <div class="y-glow pointer-events-none absolute -bottom-80 left-1/2 -z-10 h-[520px] w-[800px] -translate-x-1/2" aria-hidden="true"></div>
        <h2 class="y-headline mx-auto max-w-[720px] text-[clamp(2rem,5vw,3.5rem)]">{{ $t('home2.ctaTitle') }}</h2>
        <div class="mt-6 flex flex-wrap items-center justify-center gap-3">
          <a href="#home-sub" class="btn btn-primary px-6" @click.prevent="focusSearch">{{ $t('home2.ctaPrimary') }}</a>
          <router-link :to="docsPath" class="btn btn-outline px-6">{{ $t('home2.ctaDocs') }}</router-link>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ChevronDown, ArrowRight, ArrowUpRight, Info, Zap, ShieldCheck, Code2, Send, CircleAlert } from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'
import { useSiteConfigStore } from '../stores/siteConfig'
import axios from '../utils/axios'
import { useToast } from '../composables/useToast'
import { useCurrency } from '../composables/useCurrency'
import { GITHUB_URL, TELEGRAM_GROUP_URL, DOCS_PATH } from '../config/links'

const router = useRouter()
const { t } = useI18n()
const authStore = useAuthStore()
const siteConfigStore = useSiteConfigStore()
const toast = useToast()
const { currencySymbol } = useCurrency()

const githubUrl = GITHUB_URL
const telegramUrl = TELEGRAM_GROUP_URL
const docsPath = DOCS_PATH

const searchQuery = ref('')
const selectedRootDomain = ref('')
const rootDomains = ref([])
const searchResult = ref(null)
const loading = ref(false)

const isAuthenticated = computed(() => authStore.isAuthenticated)

const rootLabel = computed(() => {
  const root = rootDomains.value.find((d) => d.id === selectedRootDomain.value) || rootDomains.value[0]
  return root?.domain || 'yunamo.com'
})

const step2Desc = computed(() => {
  const o = siteConfigStore.oauth || {}
  const providers = []
  if (o.github) providers.push('GitHub')
  if (o.google) providers.push('Google')
  if (o.nodeloc) providers.push('NodeLoc')
  if (!providers.length) return t('home2.step2DescEmail')
  return t('home2.step2DescProviders', { providers: providers.join(' / ') })
})

const suggestions = computed(() => {
  const name = searchResult.value?.queriedName || searchQuery.value
  if (!name) return []
  return [`${name}-dev`, `my-${name}`, `${name}-app`]
})

const pickSuggestion = (s) => {
  searchQuery.value = s
  searchResult.value = null
  searchDomain()
}

const onQueryInput = () => {
  searchResult.value = null
}

const focusSearch = () => {
  document.getElementById('home-sub')?.focus()
}

onMounted(async () => {
  await fetchRootDomains()
})

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
  if (!selectedRootDomain.value) return

  const rootDomain = rootDomains.value.find((d) => d.id === selectedRootDomain.value)
  const minLen = rootDomain?.min_length || 3
  const maxLen = rootDomain?.max_length || 63
  const name = searchQuery.value.trim().toLowerCase()

  const invalid = (message) => {
    searchResult.value = { available: false, invalid: true, queriedName: name, message }
  }

  if (!name) return invalid(t('home2.emptyName'))
  if (!/^[a-z0-9-]+$/.test(name)) return invalid(t('home2.badChars'))
  if (name.length < minLen) return invalid(t('home2.tooShort', { n: minLen }))
  if (name.length > maxLen) return invalid(t('home2.tooLong', { n: maxLen }))
  if (name.startsWith('-') || name.endsWith('-')) return invalid(t('home2.hyphenEdge'))

  searchQuery.value = name
  loading.value = true
  searchResult.value = null

  try {
    const response = await axios.get('/api/domains/search', {
      params: {
        subdomain: name,
        root_domain_id: selectedRootDomain.value,
      },
    })

    const minLevel = rootDomain?.min_user_level || 'normal'
    const userLevel = authStore.user?.user_level || 'normal'
    const tooLow = isAuthenticated.value && levelOrder(userLevel) < levelOrder(minLevel)

    searchResult.value = {
      available: response.data.available,
      levelTooLow: tooLow,
      requiredLevel: minLevel,
      yourLevel: userLevel,
      queriedName: name,
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
      invalid: true,
      queriedName: name,
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
