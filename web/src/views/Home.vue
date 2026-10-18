<template>
  <div ref="homePage" class="home-page" :class="{ 'motion-enabled': motionEnabled, 'motion-running': motionActive }">
    <!-- Hero -->
    <section ref="hero" data-theme="yunamo" class="home-hero relative isolate flex min-h-[720px] items-center overflow-hidden bg-y-bg text-y-text sm:min-h-[840px]">
      <img
        src="/brand/planet-horizon-2880x1800.jpg"
        alt=""
        aria-hidden="true"
        class="hero-planet pointer-events-none absolute inset-0 h-full w-full object-cover object-[50%_80%]"
        fetchpriority="high"
      />
      <div
        aria-hidden="true"
        class="pointer-events-none absolute inset-0"
        style="background: linear-gradient(90deg, rgb(10 11 13 / 0.35) 0%, rgb(10 11 13 / 0) 50%, rgb(10 11 13 / 0.2) 100%);"
      ></div>
      <HomeOrbit :active="motionActive" />

      <button
        type="button"
        class="motion-toggle"
        :aria-pressed="motionPaused"
        :disabled="reducedMotion"
        @click="toggleMotion"
      >
        <Pause v-if="!motionPaused && !reducedMotion" :size="12" aria-hidden="true" />
        <Play v-else :size="12" aria-hidden="true" />
        {{ $t(reducedMotion ? 'home2.reducedMotion' : motionPaused ? 'home2.resumeMotion' : 'home2.pauseMotion') }}
      </button>

      <div class="y-container relative flex flex-wrap items-center gap-x-16 gap-y-12 py-20 sm:py-24">
        <div class="hero-copy flex min-w-0 flex-1 flex-col gap-6" style="flex-basis: 420px;">
          <p class="hero-enter y-eyebrow flex items-center gap-2.5" style="--enter-delay: 0ms; text-shadow: 0 1px 12px rgb(0 0 0 / 0.5);"><span class="brand-signal" aria-hidden="true"></span>{{ $t('home2.brandEyebrow') }}</p>
          <h1 class="hero-enter y-headline -mt-2 text-[clamp(2.5rem,4.6vw,4rem)] [&:lang(zh)]:text-[clamp(2.25rem,4.6vw,3.75rem)]" style="--enter-delay: 100ms; text-shadow: 0 2px 30px rgb(0 0 0 / 0.35);">
            {{ $t('home2.titleLead') }}<br />
            {{ $t('home2.titleMid') }}<em>{{ $t('home2.titleAccent') }}</em>
          </h1>
          <p class="hero-enter max-w-[500px] text-lead text-y-text/85" style="--enter-delay: 200ms;">
            {{ $t('home2.lead', { example: `yourname.${rootLabel}` }) }}
          </p>
        </div>

        <!-- Availability form -->
        <div class="hero-enter search-panel flex min-w-0 flex-1 flex-col gap-4" style="--enter-delay: 300ms; flex-basis: 420px; max-width: 560px;">
          <div class="domain-preview">
            <p class="preview-eyebrow">{{ $t('home2.previewLabel') }}</p>
            <button type="button" class="preview-domain" :aria-label="$t('home2.useExample', { name: `${exampleName}.${rootLabel}` })" @click="useExample(exampleName)">
              <Globe2 :size="19" :stroke-width="1.5" class="shrink-0 text-y-mint" aria-hidden="true" />
              <span class="preview-address">
                <Transition name="domain-swap" mode="out-in">
                  <span :key="exampleName" class="preview-name">{{ exampleName }}</span>
                </Transition><span class="text-y-muted">.{{ rootLabel }}</span>
              </span>
              <ArrowUpRight :size="17" class="ml-auto shrink-0 text-y-muted" aria-hidden="true" />
            </button>
            <p class="preview-hint">{{ $t('home2.previewHint') }}</p>
          </div>
          <form class="flex flex-col gap-3.5" @submit.prevent="searchDomain">
            <label for="home-sub" class="font-mono text-[12px] font-semibold uppercase leading-4 tracking-[0.1em] text-y-muted">
              {{ $t('home2.checkLabel') }}
            </label>
            <div class="domain-input flex items-stretch overflow-hidden rounded-md border border-y-line-strong bg-y-bg/80 backdrop-blur transition-colors focus-within:border-y-mint">
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
            <button type="submit" class="search-submit y-btn-ink-block gap-2.5 font-mono uppercase tracking-[0.12em]" :disabled="loading || !rootDomains.length">
              <span v-if="loading" class="loading loading-spinner loading-sm"></span>
              {{ loading ? $t('home2.checking') : $t('home2.checkButton') }}
              <ArrowRight v-if="!loading" :size="16" :stroke-width="2" aria-hidden="true" />
            </button>
            <div class="example-chips">
              <span>{{ $t('home2.tryLabel') }}</span>
              <button v-for="name in exampleNames" :key="name" type="button" :aria-label="$t('home2.useExample', { name })" @click="useExample(name)">{{ name }}<ArrowUpRight :size="11" aria-hidden="true" /></button>
            </div>

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

      <a href="#home-features" class="explore-link" @click.prevent="exploreFeatures">
        <span class="scroll-track" aria-hidden="true"><span></span></span>
        {{ $t('home2.explore') }}
        <ArrowDown :size="13" aria-hidden="true" />
      </a>

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

    <!-- Live numbers and built-in capabilities -->
    <section class="border-b border-y-line">
      <div class="y-container flex flex-wrap items-center justify-center gap-x-10 gap-y-3 py-7 text-[15px] font-medium text-y-faint">
        <span class="y-label inline-flex items-center gap-2"><span class="y-dot-live" aria-hidden="true"></span>{{ $t('home2.stripLabel') }}</span>
        <span v-for="item in powerStrip" :key="item.key">
          <span v-if="item.value" class="font-mono text-y-text">{{ item.value }}</span>
          {{ item.label }}
        </span>
      </div>
    </section>

    <!-- Bento features -->
    <section id="home-features" class="y-container flex flex-col gap-12 py-24 sm:py-28">
      <div data-reveal class="flex max-w-[640px] flex-col gap-3">
        <p class="y-eyebrow">{{ $t('home2.featuresEyebrow') }}</p>
        <h2 class="y-title">{{ $t('home2.featuresTitle') }}</h2>
      </div>
      <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <!-- Full DNS control (wide) -->
        <article data-reveal class="feature-card dns-card y-card flex flex-wrap items-center gap-8 overflow-hidden p-8 sm:col-span-2 lg:col-span-2" @pointermove="moveSpotlight">
          <div class="flex min-w-0 flex-col gap-3" style="flex: 1 1 240px;">
            <h3 class="text-h3">{{ $t('home2.fDnsTitle') }}</h3>
            <p class="text-y-muted">{{ $t('home2.fDnsDesc') }}</p>
          </div>
          <div class="y-inset min-w-0 overflow-hidden font-mono text-[13px] leading-5" style="flex: 1 1 320px;">
            <div class="flex gap-1.5 border-b border-y-line px-4 py-2.5">
              <span class="h-2.5 w-2.5 rounded-full bg-y-line-strong"></span>
              <span class="h-2.5 w-2.5 rounded-full bg-y-line-strong"></span>
              <span class="h-2.5 w-2.5 rounded-full bg-y-line-strong"></span>
              <span class="ml-auto text-[10px] text-y-faint">{{ $t('home2.dnsExample') }}</span>
            </div>
            <div class="dns-records flex flex-col gap-2 overflow-x-auto whitespace-pre p-4">
              <div><span class="text-y-mint-text">CNAME</span>  <span class="text-y-text">blog</span>  <span class="text-y-faint">→</span> <span class="text-y-muted">web.example.net</span></div>
              <div><span class="text-y-mint-text">A    </span>  <span class="text-y-text">@   </span>  <span class="text-y-faint">→</span> <span class="text-y-muted">203.0.113.10</span></div>
              <div><span class="text-y-mint-text">TXT  </span>  <span class="text-y-text">@   </span>  <span class="text-y-faint">→</span> <span class="text-y-muted">"verify=yunamo"</span></div>
              <div class="text-y-mint-text">✓ {{ $t('home2.propagated') }}</div>
            </div>
          </div>
        </article>

        <article data-reveal class="feature-card y-card flex flex-col gap-3 p-8" style="--reveal-delay: 70ms;" @pointermove="moveSpotlight">
          <div class="text-[56px] font-extrabold leading-none tracking-[-0.05em] text-y-mint-text">$0</div>
          <h3 class="text-xl font-semibold tracking-[-0.01em] text-y-text">{{ $t('home2.fFreeTitle') }}</h3>
          <p class="text-y-muted">{{ $t('home2.fFreeDesc') }}</p>
        </article>

        <article data-reveal class="feature-card y-card flex flex-col gap-3 p-8" @pointermove="moveSpotlight">
          <Zap :size="28" :stroke-width="1.5" class="text-y-mint-text" aria-hidden="true" />
          <h3 class="mt-2 text-xl font-semibold tracking-[-0.01em] text-y-text">{{ $t('home2.fLiveTitle') }}</h3>
          <p class="text-y-muted">{{ $t('home2.fLiveDesc') }}</p>
        </article>

        <article data-reveal class="feature-card y-card flex flex-col gap-3 p-8" style="--reveal-delay: 70ms;" @pointermove="moveSpotlight">
          <ShieldCheck :size="28" :stroke-width="1.5" class="text-y-mint-text" aria-hidden="true" />
          <h3 class="mt-2 text-xl font-semibold tracking-[-0.01em] text-y-text">{{ $t('home2.fTrustTitle') }}</h3>
          <p class="text-y-muted">{{ $t('home2.fTrustDesc') }}</p>
        </article>

        <article data-reveal class="feature-card y-card is-hoverable flex flex-col gap-3 p-8" style="--reveal-delay: 140ms;" @pointermove="moveSpotlight">
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
        <div data-reveal class="flex max-w-[640px] flex-col gap-3">
          <p class="y-eyebrow">{{ $t('home2.howEyebrow') }}</p>
          <h2 class="y-title">{{ $t('home2.howTitle') }}</h2>
        </div>
        <ol data-reveal class="launch-steps m-0 grid list-none gap-x-6 gap-y-10 p-0 sm:grid-cols-3">
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
        <div data-reveal class="flex min-w-0 flex-col gap-3" style="flex: 1 1 240px;">
          <p class="y-eyebrow">{{ $t('home2.faqEyebrow') }}</p>
          <h2 class="y-title">{{ $t('home2.faqTitle') }}</h2>
        </div>
        <div data-reveal class="y-faq min-w-0 flex flex-col border-t border-y-line" style="flex: 2 1 480px; --reveal-delay: 100ms;">
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
      <div data-reveal class="cta-card y-card relative isolate overflow-hidden px-6 py-16 text-center sm:py-20" @pointermove="moveSpotlight">
        <div class="y-glow pointer-events-none absolute -bottom-80 left-1/2 -z-10 h-[520px] w-[800px] -translate-x-1/2" aria-hidden="true"></div>
        <h2 class="y-headline mx-auto max-w-[720px] text-[clamp(2rem,5vw,3.5rem)] [&:lang(zh)]:text-[clamp(1.75rem,4vw,2.75rem)]">{{ $t('home2.ctaTitle') }}</h2>
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
import { ChevronDown, ArrowRight, ArrowUpRight, ArrowDown, Info, Zap, ShieldCheck, Code2, Send, CircleAlert, Globe2, Pause, Play } from 'lucide-vue-next'
import HomeOrbit from '../components/HomeOrbit.vue'
import { useHomeMotion } from '../composables/useHomeMotion'
import { useAuthStore } from '../stores/auth'
import { useSiteConfigStore } from '../stores/siteConfig'
import axios from '../utils/axios'
import { useToast } from '../composables/useToast'
import { useCurrency } from '../composables/useCurrency'
import { GITHUB_URL, TELEGRAM_GROUP_URL, DOCS_PATH } from '../config/links'

const router = useRouter()
const { t, locale } = useI18n()
const authStore = useAuthStore()
const siteConfigStore = useSiteConfigStore()
const toast = useToast()
const { currencySymbol } = useCurrency()
const { homePage, hero, motionPaused, reducedMotion, motionEnabled, motionActive,
  exampleNames, exampleName, toggleMotion, moveSpotlight } = useHomeMotion()

const githubUrl = GITHUB_URL
const telegramUrl = TELEGRAM_GROUP_URL
const docsPath = DOCS_PATH

const searchQuery = ref('')
const selectedRootDomain = ref('')
const rootDomains = ref([])
const searchResult = ref(null)
const loading = ref(false)

const isAuthenticated = computed(() => authStore.isAuthenticated)

const publicStats = ref(null)
const DNS_RECORD_TYPE_COUNT = 9

const powerStrip = computed(() => {
  const fmt = (n) => Number(n).toLocaleString(locale.value)
  const items = []
  if (publicStats.value) {
    items.push({ key: 'names', value: fmt(publicStats.value.total_domains || 0), label: t('home2.statNames') })
    items.push({ key: 'builders', value: fmt(publicStats.value.total_users || 0), label: t('home2.statBuilders') })
  }
  if (rootDomains.value.length) {
    items.push({ key: 'suffixes', value: fmt(rootDomains.value.length), label: t('home2.statSuffixes') })
  }
  items.push({ key: 'records', value: fmt(DNS_RECORD_TYPE_COUNT), label: t('home2.statRecordTypes') })
  items.push({ key: 'dnssec', value: '', label: 'DNSSEC' })
  items.push({ key: 'ns', value: '', label: t('home2.stripCustomNs') })
  return items
})

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

const useExample = (name) => {
  searchQuery.value = name
  searchResult.value = null
  focusSearch()
}

const exploreFeatures = () => {
  document.getElementById('home-features')?.scrollIntoView({ behavior: motionEnabled.value ? 'smooth' : 'instant' })
}

onMounted(async () => {
  await Promise.all([fetchRootDomains(), fetchPublicStats()])
})

const fetchPublicStats = async () => {
  try {
    const response = await axios.get('/api/public/stats')
    publicStats.value = response.data
  } catch (error) {
    publicStats.value = null
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
    } else if (errorData?.code === 'email_not_verified') {
      toast.error(t('authEmail.claim.notVerified'), 6000)
    } else {
      toast.error(errorData?.error || 'Registration failed')
    }
  }
}
</script>

<style scoped>
.home-hero { --hero-edge: 24px; }
.hero-planet {
  transform: scale(1.035);
  animation: planet-drift 36s ease-in-out infinite alternate;
  animation-play-state: paused;
}
.motion-running .hero-planet { animation-play-state: running; }
.motion-toggle {
  position: absolute;
  z-index: 2;
  top: var(--hero-edge);
  right: var(--hero-edge);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 7px;
  min-height: 44px;
  padding: 0 13px;
  border: 1px solid rgb(255 255 255 / .14);
  border-radius: 999px;
  background: rgb(10 11 13 / .5);
  color: #c2cdc9;
  font-family: 'Geist Mono', monospace;
  font-size: 10px;
  transition: color .2s, border-color .2s, background .2s;
}
.motion-toggle:hover:not(:disabled) { color: #ededed; border-color: rgb(62 230 180 / .5); background: rgb(10 11 13 / .8); }
.motion-toggle:disabled { cursor: default; }
.brand-signal { width: 6px; height: 6px; flex-shrink: 0; border-radius: 50%; background: #3ee6b4; box-shadow: 0 0 14px rgb(62 230 180 / .45); }
.motion-enabled .hero-enter { animation: hero-arrive .85s cubic-bezier(.16, 1, .3, 1) both; animation-delay: var(--enter-delay, 0ms); }
.search-panel {
  border: 1px solid rgb(255 255 255 / .14);
  border-radius: 18px;
  padding: 26px;
  background: linear-gradient(135deg, rgb(17 25 24 / .86), rgb(10 11 13 / .74));
  box-shadow: 0 16px 80px rgb(0 0 0 / .22), inset 0 1px 0 rgb(255 255 255 / .04);
  backdrop-filter: blur(16px);
}
.domain-preview { border-bottom: 1px solid rgb(255 255 255 / .1); padding-bottom: 19px; margin-bottom: 3px; }
.preview-eyebrow { font-size: 12px; color: #b1bdb8; margin-bottom: 15px; }
.preview-domain { display: flex; align-items: center; gap: 11px; width: 100%; min-height: 44px; text-align: left; }
.preview-address { min-width: 0; display: flex; align-items: baseline; overflow: hidden; font-family: 'Geist Mono', monospace; font-size: clamp(14px, 1.45vw, 20px); white-space: nowrap; letter-spacing: -.045em; }
.preview-name { display: inline-block; color: #ededed; }
.preview-domain > svg:last-child { transition: transform .2s, color .2s; }
.preview-domain:hover > svg:last-child { transform: translate(2px, -2px); color: #3ee6b4; }
.preview-hint { padding-left: 30px; color: #8a9e95; font-family: 'Geist Mono', monospace; font-size: 10px; margin-top: 5px; }
.domain-input { transition: border-color .2s, box-shadow .2s; }
.domain-input:focus-within { box-shadow: 0 0 0 3px rgb(62 230 180 / .1); }
.search-submit { position: relative; overflow: hidden; transition: background .2s, box-shadow .2s; }
.search-submit::after { content: ''; position: absolute; inset: 0; background: linear-gradient(110deg, transparent 25%, rgb(62 230 180 / .22) 50%, transparent 75%); transform: translateX(-110%); }
.search-submit:hover:not(:disabled)::after { transform: translateX(110%); transition: transform .7s ease; }
.search-submit svg { transition: transform .2s; }
.search-submit:hover:not(:disabled) svg { transform: translateX(3px); }
.example-chips { display: flex; flex-wrap: wrap; align-items: center; gap: 2px 6px; color: #9aa8a2; font-size: 11px; }
.example-chips button { display: inline-flex; align-items: center; gap: 4px; min-height: 32px; padding: 3px 7px; border: 1px solid transparent; border-radius: 5px; font-family: 'Geist Mono', monospace; color: #b7c7bf; transition: border-color .2s, background .2s, color .2s; }
.example-chips button:hover { color: #3ee6b4; border-color: rgb(62 230 180 / .2); background: rgb(62 230 180 / .06); }
.explore-link { position: absolute; bottom: 25px; left: max(24px, calc((100% - 1160px) / 2 + 24px)); display: flex; align-items: center; gap: 10px; min-height: 44px; font-family: 'Geist Mono', monospace; font-size: 10px; color: #b2c0ba; transition: color .2s; }
.explore-link:hover { color: #3ee6b4; }
.scroll-track { display: inline-block; position: relative; width: 13px; height: 21px; border: 1px solid rgb(186 210 199 / .45); border-radius: 8px; }
.scroll-track span { position: absolute; top: 4px; left: 5px; height: 4px; width: 1px; background: #3ee6b4; animation: scroll-hint 2.5s ease-in-out infinite; animation-play-state: paused; }
.motion-running .scroll-track span { animation-play-state: running; }
#home-features { scroll-margin-top: 80px; }
[data-reveal] { transition: opacity .7s ease, transform .7s cubic-bezier(.16, 1, .3, 1), border-color .25s; transition-delay: var(--reveal-delay, 0ms); }
.motion-observed.motion-enabled [data-reveal]:not(.is-revealed) { opacity: 0; transform: translateY(24px); }
.feature-card { position: relative; }
@media (hover: hover) and (pointer: fine) {
  .feature-card:hover, .cta-card:hover { border-color: rgb(var(--y-mint) / .35); background-image: radial-gradient(380px circle at var(--spot-x, 50%) var(--spot-y, 50%), rgb(var(--y-mint) / .065), transparent 75%); }
  .motion-enabled .feature-card.is-revealed:hover { transform: translateY(-5px); transition-delay: 0ms; }
}
.motion-enabled .dns-card.is-revealed .dns-records > div { animation: record-arrive .6s ease both; }
.motion-enabled .dns-card.is-revealed .dns-records > div:nth-child(2) { animation-delay: .2s; }
.motion-enabled .dns-card.is-revealed .dns-records > div:nth-child(3) { animation-delay: .4s; }
.motion-enabled .dns-card.is-revealed .dns-records > div:nth-child(4) { animation-delay: .7s; }
.launch-steps > li { position: relative; }
.launch-steps > li::before { content: ''; position: absolute; top: -1px; left: 0; width: 100%; height: 1px; background: linear-gradient(90deg, rgb(var(--y-mint)), rgb(var(--y-mint) / .15)); transform-origin: left; }
.motion-enabled .launch-steps.is-revealed > li::before { animation: connect-step .8s ease both; }
.motion-enabled .launch-steps.is-revealed > li:nth-child(2)::before { animation-delay: .3s; }
.motion-enabled .launch-steps.is-revealed > li:nth-child(3)::before { animation-delay: .6s; }
.motion-enabled .y-faq details[open] > p { animation: record-arrive .25s ease; }
.domain-swap-enter-active, .domain-swap-leave-active { transition: opacity .25s ease, transform .25s ease; }
.domain-swap-enter-from { opacity: 0; transform: translateY(9px); }
.domain-swap-leave-to { opacity: 0; transform: translateY(-9px); }
@keyframes hero-arrive { from { opacity: 0; transform: translateY(20px); } to { opacity: 1; transform: translateY(0); } }
@keyframes planet-drift { from { transform: scale(1.035) translateY(3px); } to { transform: scale(1.06) translateY(-7px); } }
@keyframes scroll-hint { 0%, 100% { transform: translateY(0); opacity: .35; } 50% { transform: translateY(6px); opacity: 1; } }
@keyframes record-arrive { from { opacity: 0; transform: translateY(5px); } to { opacity: 1; transform: translateY(0); } }
@keyframes connect-step { from { transform: scaleX(0); } to { transform: scaleX(1); } }
@media (max-width: 639px) {
  .home-hero { min-height: 900px; align-items: flex-start; }
  .home-hero > .y-container { padding-top: 96px; padding-bottom: 130px; row-gap: 32px; }
  .home-hero .hero-copy { gap: 20px; }
  .home-hero h1 { font-size: clamp(2.5rem, 9vw, 3.5rem); }
  .search-panel { padding: 20px; border-radius: 14px; }
  .preview-address { font-size: clamp(12px, 3.7vw, 18px); }
  .example-chips { gap: 0 2px; }
  .example-chips button { font-size: 10px; padding: 3px 5px; }
  .explore-link { bottom: 73px; left: 24px; }
  .home-hero > a:last-child { bottom: 22px; right: 20px; }
}
@media (prefers-reduced-motion: reduce) {
  .home-page *, .home-page *::before, .home-page *::after { animation: none !important; transition: none !important; }
  [data-reveal] { opacity: 1 !important; transform: none !important; }
}
.home-page:not(.motion-enabled) *,
.home-page:not(.motion-enabled) *::before,
.home-page:not(.motion-enabled) *::after { animation-play-state: paused !important; }
.home-page:not(.motion-enabled) .domain-swap-enter-active,
.home-page:not(.motion-enabled) .domain-swap-leave-active { transition: none; }
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
