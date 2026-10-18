<template>
  <div>
    <!-- Hero -->
    <section class="relative overflow-hidden border-b border-y-line">
      <div class="y-dot-grid pointer-events-none absolute inset-0" aria-hidden="true"></div>
      <div class="y-glow pointer-events-none absolute -top-48 right-[-160px] h-[560px] w-[560px]" aria-hidden="true"></div>
      <div class="y-container relative flex flex-col gap-5 pb-16 pt-16 sm:pb-20 sm:pt-24">
        <p class="y-eyebrow m-0 flex items-center gap-2">
          <Heart :size="14" :stroke-width="2" aria-hidden="true" />
          {{ $t('sponsor.eyebrow') }}
        </p>
        <i18n-t keypath="sponsor.heroTitle" tag="h1" class="y-headline m-0 max-w-[920px] text-[clamp(2.75rem,7vw,5.25rem)]">
          <template #accent><em>{{ $t('sponsor.heroAccent') }}</em></template>
        </i18n-t>
        <p class="y-lead m-0 max-w-[640px]">{{ $t('sponsor.heroLead') }}</p>
        <div class="mt-2 flex flex-wrap items-center gap-3">
          <a href="#support" class="y-btn-ink min-h-[48px] px-6" @click.prevent="scrollToSupport">
            {{ $t('sponsor.ctaSupport') }}
            <ArrowDown :size="16" :stroke-width="2" aria-hidden="true" />
          </a>
          <a href="#become-partner" class="btn btn-outline min-h-[48px] rounded-full px-6" @click.prevent="scrollToPartner()">
            <Handshake :size="16" :stroke-width="1.5" aria-hidden="true" />
            {{ $t('sponsor.becomePartner') }}
          </a>
        </div>
        <p class="m-0 mt-3 font-mono text-xs uppercase tracking-[0.08em] text-y-faint">{{ $t('sponsor.brandLine') }}</p>
      </div>
    </section>

    <!-- Partners and sponsors -->
    <section class="y-container flex flex-col gap-10 py-16 sm:py-20" aria-labelledby="sponsor-partners-title">
      <div class="flex flex-wrap items-end justify-between gap-x-8 gap-y-5">
        <div class="flex min-w-0 max-w-[680px] flex-col gap-3">
          <p class="y-eyebrow m-0">{{ $t('sponsor.partnersEyebrow') }}</p>
          <i18n-t id="sponsor-partners-title" keypath="sponsor.partnersTitle" tag="h2" class="y-title m-0">
            <template #accent><em>{{ $t('sponsor.partnersAccent') }}</em></template>
          </i18n-t>
          <p class="m-0 text-y-muted">{{ $t('sponsor.partnersLead') }}</p>
        </div>
        <a href="#become-partner" class="btn btn-outline btn-sm gap-1.5" @click.prevent="scrollToPartner()">
          <Plus :size="14" :stroke-width="1.75" aria-hidden="true" />
          {{ $t('sponsor.becomePartner') }}
        </a>
      </div>

      <!-- Loading skeleton -->
      <div v-if="loading" class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3" role="status" :aria-label="$t('sponsor.loading')">
        <div v-for="i in 3" :key="i" class="y-card flex flex-col gap-4 p-6">
          <div class="h-12 w-28 animate-pulse rounded-md bg-y-surface-2"></div>
          <div class="h-4 w-1/2 animate-pulse rounded bg-y-surface-2"></div>
          <div class="h-3 w-full animate-pulse rounded bg-y-surface-2"></div>
          <div class="h-3 w-4/5 animate-pulse rounded bg-y-surface-2"></div>
        </div>
      </div>

      <!-- Error -->
      <div v-else-if="loadError" role="alert" class="y-warn-soft flex flex-wrap items-center gap-x-4 gap-y-3 rounded-lg border px-5 py-4">
        <TriangleAlert :size="18" :stroke-width="1.5" class="shrink-0 text-y-warn" aria-hidden="true" />
        <span class="min-w-0 flex-1 text-sm text-y-text">{{ $t('sponsor.loadFailed') }}</span>
        <button type="button" class="btn btn-outline btn-sm gap-1.5" @click="load">
          <RefreshCw :size="14" :stroke-width="1.5" aria-hidden="true" />
          {{ $t('sponsor.retry') }}
        </button>
      </div>

      <!-- Empty -->
      <div
        v-else-if="groups.length === 0"
        class="relative flex flex-col items-start gap-4 overflow-hidden rounded-xl border border-dashed border-y-line-strong p-8 sm:flex-row sm:items-center sm:gap-6"
      >
        <span class="y-icon-tile h-14 w-14" aria-hidden="true"><HeartHandshake :size="26" :stroke-width="1.5" /></span>
        <div class="min-w-0 flex-1">
          <h3 class="m-0 text-xl font-bold text-y-text">{{ $t('sponsor.emptyTitle') }}</h3>
          <p class="m-0 mt-1.5 max-w-[560px] text-sm leading-6 text-y-muted">{{ $t('sponsor.emptyDesc') }}</p>
        </div>
        <a href="#become-partner" class="y-btn-ink shrink-0" @click.prevent="scrollToPartner()">{{ $t('sponsor.becomePartner') }}</a>
      </div>

      <!-- Tier groups -->
      <div v-else class="flex flex-col gap-12">
        <div v-for="group in groups" :key="group.tier" class="flex flex-col gap-5">
          <div class="flex items-center gap-3">
            <span class="y-label">{{ $t(`sponsor.tiers.${group.tier}`) }}</span>
            <span class="font-mono text-xs text-y-faint">{{ group.items.length }}</span>
            <span class="h-px flex-1 bg-y-line" aria-hidden="true"></span>
          </div>
          <div class="grid gap-4" :class="group.tier === 'partner' ? 'sm:grid-cols-2 lg:grid-cols-3' : 'grid-cols-2 lg:grid-cols-4'">
            <component
              :is="s.website_url ? 'a' : 'div'"
              v-for="s in group.items"
              :key="s.id"
              :href="s.website_url || undefined"
              :target="s.website_url ? '_blank' : undefined"
              :rel="s.website_url ? 'noopener sponsored' : undefined"
              :aria-label="s.website_url ? $t('sponsor.visitSite', { name: s.name }) : undefined"
              class="y-card group flex min-w-0 flex-col"
              :class="group.tier === 'partner' ? 'gap-4 p-6' : 'gap-3 p-4 sm:p-5'"
            >
              <div class="flex items-center" :class="group.tier === 'partner' ? 'h-12' : 'h-9'">
                <SponsorLogo :src="logoFor(s)" :name="s.name" :small="group.tier !== 'partner'" />
              </div>
              <div class="flex min-w-0 items-center gap-1.5">
                <span class="truncate font-semibold text-y-text" :class="group.tier === 'partner' ? 'text-[17px]' : 'text-[15px]'">{{ s.name }}</span>
                <ArrowUpRight
                  v-if="s.website_url"
                  :size="15"
                  :stroke-width="1.75"
                  class="shrink-0 text-y-faint transition-colors group-hover:text-y-mint-text"
                  aria-hidden="true"
                />
              </div>
              <p
                v-if="descriptionFor(s)"
                class="m-0 text-y-muted"
                :class="group.tier === 'partner' ? 'line-clamp-3 text-sm leading-6' : 'line-clamp-2 text-[13px] leading-5'"
              >
                {{ descriptionFor(s) }}
              </p>
            </component>
          </div>
        </div>
      </div>
    </section>

    <!-- Ways to support -->
    <section id="support" class="y-band scroll-mt-20" aria-labelledby="sponsor-ways-title">
      <div class="y-container flex flex-col gap-10 py-16 sm:py-20">
        <div class="flex max-w-[680px] flex-col gap-3">
          <p class="y-eyebrow m-0">{{ $t('sponsor.waysEyebrow') }}</p>
          <i18n-t id="sponsor-ways-title" keypath="sponsor.waysTitle" tag="h2" class="y-title m-0">
            <template #accent><em>{{ $t('sponsor.waysAccent') }}</em></template>
          </i18n-t>
          <p class="m-0 text-y-muted">{{ $t('sponsor.waysLead') }}</p>
        </div>

        <div class="grid items-stretch gap-4" :class="waysGridClass">
          <!-- Payment skeleton while loading -->
          <div v-if="loading" class="y-card flex flex-col gap-4 p-6" aria-hidden="true">
            <div class="h-10 w-10 animate-pulse rounded-md bg-y-surface-2"></div>
            <div class="h-4 w-1/3 animate-pulse rounded bg-y-surface-2"></div>
            <div class="h-12 w-full animate-pulse rounded-md bg-y-surface-2"></div>
            <div class="h-3 w-3/4 animate-pulse rounded bg-y-surface-2"></div>
          </div>

          <!-- USDT TRC20 -->
          <article v-if="usdtAddress" class="y-card flex min-w-0 flex-col gap-4 p-6">
            <div class="flex flex-wrap items-center gap-3">
              <span class="y-icon-tile" aria-hidden="true"><Coins :size="20" :stroke-width="1.5" /></span>
              <h3 class="m-0 text-lg font-bold text-y-text">{{ $t('sponsor.usdtTitle') }}</h3>
              <span class="y-chip-warn ml-auto font-mono">{{ $t('sponsor.usdtNetwork') }}</span>
            </div>
            <p class="m-0 text-sm leading-6 text-y-muted">{{ $t('sponsor.usdtDesc') }}</p>
            <div class="flex flex-col gap-2">
              <span id="sponsor-usdt-label" class="y-label">{{ $t('sponsor.usdtAddressLabel') }}</span>
              <div
                class="y-inset select-all break-all px-4 py-3 font-mono text-[14px] leading-6 text-y-text"
                aria-labelledby="sponsor-usdt-label"
              >{{ usdtAddress }}</div>
            </div>
            <button type="button" class="btn btn-outline btn-sm gap-1.5 self-start" @click="copyAddress">
              <Check v-if="copiedAddress" :size="14" :stroke-width="2" class="text-y-mint-text" aria-hidden="true" />
              <Copy v-else :size="14" :stroke-width="1.5" aria-hidden="true" />
              {{ copiedAddress ? $t('sponsor.copied') : $t('sponsor.copy') }}
            </button>
            <div class="y-warn-soft mt-auto flex items-start gap-2.5 rounded-md border px-3.5 py-3 text-[13px] leading-5 text-y-text" role="note">
              <TriangleAlert :size="16" :stroke-width="1.75" class="mt-0.5 shrink-0 text-y-warn" aria-hidden="true" />
              <span>{{ $t('sponsor.usdtWarning') }}</span>
            </div>
          </article>

          <!-- GitHub Sponsors -->
          <article v-if="githubSponsorsUrl" class="y-card flex min-w-0 flex-col gap-4 p-6">
            <div class="flex items-center gap-3">
              <span class="y-icon-tile" aria-hidden="true"><Github :size="20" :stroke-width="1.5" /></span>
              <h3 class="m-0 text-lg font-bold text-y-text">{{ $t('sponsor.githubTitle') }}</h3>
            </div>
            <p class="m-0 text-sm leading-6 text-y-muted">{{ $t('sponsor.githubDesc') }}</p>
            <div class="mt-auto pt-2">
              <a :href="githubSponsorsUrl" target="_blank" rel="noopener" class="y-btn-ink">
                <Heart :size="15" :stroke-width="2" aria-hidden="true" />
                {{ $t('sponsor.githubButton') }}
              </a>
            </div>
          </article>

          <!-- Business and resource partnerships -->
          <article class="y-card flex min-w-0 flex-col gap-4 p-6">
            <div class="flex items-center gap-3">
              <span class="y-icon-tile" aria-hidden="true"><Handshake :size="20" :stroke-width="1.5" /></span>
              <h3 class="m-0 text-lg font-bold leading-6 text-y-text">{{ $t('sponsor.businessTitle') }}</h3>
            </div>
            <p class="m-0 text-sm leading-6 text-y-muted">{{ $t('sponsor.businessDesc') }}</p>
            <ul class="m-0 flex list-none flex-wrap gap-2 p-0">
              <li v-for="key in ['servers', 'bandwidth', 'domains', 'longTerm']" :key="key" class="y-chip">
                {{ $t(`sponsor.businessItems.${key}`) }}
              </li>
            </ul>
            <div class="mt-auto pt-2">
              <a
                href="#become-partner"
                class="btn btn-outline btn-sm h-auto min-h-[2.125rem] max-w-full gap-1.5 py-1.5"
                @click.prevent="scrollToPartner('resources')"
              >
                <ArrowDown :size="14" :stroke-width="1.75" class="shrink-0" aria-hidden="true" />
                <span class="text-left font-medium">{{ $t('partner.resourcesCta') }}</span>
              </a>
            </div>
          </article>
        </div>

        <p v-if="!loading && !usdtAddress && !githubSponsorsUrl" class="m-0 flex items-start gap-2 text-sm text-y-faint">
          <Sparkles :size="16" :stroke-width="1.5" class="mt-0.5 shrink-0 text-y-mint-text" aria-hidden="true" />
          <span>{{ $t('sponsor.moreComing') }}</span>
        </p>
      </div>
    </section>

    <!-- Become a partner -->
    <section id="become-partner" class="y-band scroll-mt-20" aria-labelledby="partner-title">
      <!-- On small screens the form follows the intro directly; the details move below it. -->
      <div class="y-container grid gap-10 py-16 sm:py-20 lg:grid-cols-[minmax(0,380px)_minmax(0,1fr)] lg:grid-rows-[auto_1fr] lg:gap-x-12 lg:gap-y-8">
        <div class="flex min-w-0 flex-col gap-3 lg:col-start-1 lg:row-start-1">
          <p class="y-eyebrow m-0 flex items-center gap-2">
            <Handshake :size="14" :stroke-width="2" aria-hidden="true" />
            {{ $t('partner.eyebrow') }}
          </p>
          <i18n-t id="partner-title" keypath="partner.title" tag="h2" class="y-title m-0 outline-none" tabindex="-1">
            <template #accent><em>{{ $t('partner.titleAccent') }}</em></template>
          </i18n-t>
          <p class="m-0 text-y-muted">{{ $t('partner.lead') }}</p>
        </div>

        <div class="min-w-0 lg:col-start-2 lg:row-span-2 lg:row-start-1">
          <PartnerApplicationForm ref="partnerFormRef" />
        </div>

        <div class="flex min-w-0 flex-col gap-8 lg:col-start-1 lg:row-start-2">
          <div class="flex flex-col gap-3.5">
            <span class="y-label">{{ $t('partner.perksTitle') }}</span>
            <ul class="m-0 flex list-none flex-col gap-3 p-0">
              <li v-for="key in ['logo', 'reach', 'team']" :key="key" class="flex items-start gap-3 text-sm leading-6 text-y-text">
                <Check :size="16" :stroke-width="2" class="mt-1 shrink-0 text-y-mint-text" aria-hidden="true" />
                <span>{{ $t(`partner.perks.${key}`) }}</span>
              </li>
            </ul>
          </div>

          <div class="flex flex-col gap-3.5">
            <span class="y-label">{{ $t('partner.nextTitle') }}</span>
            <ol class="m-0 flex list-none flex-col gap-4 p-0">
              <li v-for="(key, i) in ['read', 'reply', 'plan']" :key="key" class="flex gap-3.5">
                <span :class="i === 0 ? 'y-step-num' : 'y-step-num-muted'" aria-hidden="true">{{ i + 1 }}</span>
                <div>
                  <div class="font-semibold text-y-text">{{ $t(`partner.next.${key}Title`) }}</div>
                  <div class="mt-0.5 text-sm leading-[22px] text-y-muted">{{ $t(`partner.next.${key}Desc`) }}</div>
                </div>
              </li>
            </ol>
          </div>

          <i18n-t keypath="partner.fallback" tag="p" class="m-0 text-sm leading-6 text-y-muted">
            <template #email>
              <a :href="partnerMailto" class="break-all font-medium text-y-mint-text hover:text-y-mint-hover">{{ supportEmail }}</a>
            </template>
          </i18n-t>
        </div>
      </div>
    </section>

    <!-- What your support sustains -->
    <section class="y-band" aria-labelledby="sponsor-sustains-title">
      <div class="y-container flex flex-col gap-10 py-16 sm:py-20">
        <div class="flex max-w-[680px] flex-col gap-3">
          <p class="y-eyebrow m-0">{{ $t('sponsor.sustainsEyebrow') }}</p>
          <i18n-t id="sponsor-sustains-title" keypath="sponsor.sustainsTitle" tag="h2" class="y-title m-0">
            <template #accent><em>{{ $t('sponsor.sustainsAccent') }}</em></template>
          </i18n-t>
        </div>
        <ol class="m-0 grid list-none gap-x-8 gap-y-10 p-0 md:grid-cols-3">
          <li v-for="(item, i) in sustains" :key="item.key" class="flex flex-col gap-3 border-t border-y-line-strong pt-5">
            <div class="flex items-center justify-between">
              <span class="font-mono text-sm font-medium text-y-mint-text">{{ String(i + 1).padStart(2, '0') }}</span>
              <component :is="item.icon" :size="20" :stroke-width="1.5" class="text-y-faint" aria-hidden="true" />
            </div>
            <h3 class="m-0 text-xl font-bold tracking-[-0.02em] text-y-text">{{ item.title }}</h3>
            <p class="m-0 text-sm leading-6 text-y-muted">{{ item.desc }}</p>
          </li>
        </ol>
      </div>
    </section>

    <!-- Other ways to help -->
    <section class="y-band" aria-labelledby="sponsor-help-title">
      <div class="y-container flex flex-col gap-10 py-16 sm:py-20">
        <div class="flex max-w-[680px] flex-col gap-3">
          <p class="y-eyebrow m-0">{{ $t('sponsor.helpEyebrow') }}</p>
          <h2 id="sponsor-help-title" class="y-title m-0">{{ $t('sponsor.helpTitle') }}</h2>
        </div>
        <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          <component
            :is="item.href ? 'a' : 'button'"
            v-for="item in helpItems"
            :key="item.key"
            :href="item.href"
            :target="item.href ? '_blank' : undefined"
            :rel="item.href ? 'noopener' : undefined"
            :type="item.href ? undefined : 'button'"
            class="y-card group flex flex-col gap-3 p-5 text-left"
            @click="item.action ? item.action() : undefined"
          >
            <span class="y-icon-tile" aria-hidden="true"><component :is="item.icon" :size="18" :stroke-width="1.5" /></span>
            <span class="flex items-center gap-1.5 font-semibold text-y-text">
              {{ item.title }}
              <ArrowUpRight
                v-if="item.href"
                :size="15"
                :stroke-width="1.75"
                class="shrink-0 text-y-faint transition-colors group-hover:text-y-mint-text"
                aria-hidden="true"
              />
            </span>
            <span class="text-sm leading-6 text-y-muted">{{ item.desc }}</span>
          </component>
        </div>
      </div>
    </section>

    <!-- Closing -->
    <section class="y-container pb-20">
      <div class="y-mint-soft relative flex flex-wrap items-center gap-x-8 gap-y-5 overflow-hidden rounded-xl border p-7 sm:p-9">
        <div class="y-glow pointer-events-none absolute -bottom-40 -left-24 h-[360px] w-[360px]" aria-hidden="true"></div>
        <div class="relative min-w-0 flex-1" style="flex-basis: 420px;">
          <h2 class="m-0 text-2xl font-extrabold tracking-[-0.03em] text-y-text sm:text-[28px] sm:leading-[34px]">{{ $t('sponsor.thanksTitle') }}</h2>
          <i18n-t keypath="sponsor.thanksDesc" tag="p" class="m-0 mt-2 text-y-muted">
            <template #email>
              <a :href="partnerMailto" class="break-all font-medium text-y-mint-text hover:text-y-mint-hover">{{ supportEmail }}</a>
            </template>
          </i18n-t>
        </div>
        <a href="#become-partner" class="y-btn-ink relative shrink-0" @click.prevent="scrollToPartner()">{{ $t('sponsor.becomePartner') }}</a>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  Heart, HeartHandshake, Handshake, ArrowDown, ArrowUpRight, Plus, TriangleAlert, RefreshCw,
  Coins, Copy, Check, Github, Sparkles, Server, Network, CodeXml, Star, Share2, Bug, GitPullRequest,
} from 'lucide-vue-next'
import axios from '../utils/axios'
import { copyText } from '../utils/clipboard'
import { useToast } from '../composables/useToast'
import { useTheme } from '../composables/useTheme'
import { useSiteConfigStore } from '../stores/siteConfig'
import { GITHUB_URL, SUPPORT_EMAIL } from '../config/links'
import SponsorLogo from '../components/sponsor/SponsorLogo.vue'
import PartnerApplicationForm from '../components/sponsor/PartnerApplicationForm.vue'

const route = useRoute()
const { t, locale } = useI18n()
const toast = useToast()
const { isDark } = useTheme()
const siteConfigStore = useSiteConfigStore()

const supportEmail = SUPPORT_EMAIL
const TIERS = ['partner', 'sponsor', 'backer']

const loading = ref(true)
const loadError = ref(false)
const sponsors = ref([])
const payment = ref({ usdt_trc20_address: '', github_sponsors_url: '' })
const copiedAddress = ref(false)

const usdtAddress = computed(() => (payment.value.usdt_trc20_address || '').trim())
const githubSponsorsUrl = computed(() => (payment.value.github_sponsors_url || '').trim())

const mailto = (subject) => `mailto:${SUPPORT_EMAIL}?subject=${encodeURIComponent(subject)}`
const partnerMailto = computed(() => mailto(t('sponsor.mailSubjectPartner')))
const partnerFormRef = ref(null)

const groups = computed(() => {
  const byTier = { partner: [], sponsor: [], backer: [] }
  for (const s of sponsors.value) {
    const tier = TIERS.includes(s.tier) ? s.tier : 'backer'
    byTier[tier].push(s)
  }
  return TIERS.filter((tier) => byTier[tier].length).map((tier) => ({ tier, items: byTier[tier] }))
})

const waysGridClass = computed(() => {
  const count = (loading.value ? 1 : 0) + (usdtAddress.value ? 1 : 0) + (githubSponsorsUrl.value ? 1 : 0) + 1
  if (count >= 3) return 'md:grid-cols-2 lg:grid-cols-3'
  if (count === 2) return 'md:grid-cols-2'
  return 'max-w-[640px]'
})

const logoFor = (s) => (isDark.value && s.logo_dark_url ? s.logo_dark_url : s.logo_url) || ''

const descriptionFor = (s) => {
  const zh = (s.description_zh || '').trim()
  const en = (s.description_en || '').trim()
  return String(locale.value).startsWith('zh') ? zh || en : en || zh
}

const sustains = computed(() => [
  { key: 'servers', icon: Server, title: t('sponsor.sustains.serversTitle'), desc: t('sponsor.sustains.serversDesc') },
  { key: 'dns', icon: Network, title: t('sponsor.sustains.dnsTitle'), desc: t('sponsor.sustains.dnsDesc') },
  { key: 'dev', icon: CodeXml, title: t('sponsor.sustains.devTitle'), desc: t('sponsor.sustains.devDesc') },
])

const shareSite = async () => {
  const url = window.location.origin + '/'
  if (navigator.share) {
    try {
      await navigator.share({ title: siteConfigStore.siteName || 'Yunamo', text: t('sponsor.shareText'), url })
      return
    } catch (err) {
      if (err?.name === 'AbortError') return
    }
  }
  if (await copyText(url)) toast.success(t('sponsor.linkCopied'))
  else toast.error(t('sponsor.copyFailed'))
}

const helpItems = computed(() => [
  { key: 'star', icon: Star, href: GITHUB_URL, title: t('sponsor.help.starTitle'), desc: t('sponsor.help.starDesc') },
  { key: 'share', icon: Share2, action: shareSite, title: t('sponsor.help.shareTitle'), desc: t('sponsor.help.shareDesc') },
  { key: 'bugs', icon: Bug, href: `${GITHUB_URL}/issues`, title: t('sponsor.help.bugsTitle'), desc: t('sponsor.help.bugsDesc') },
  { key: 'code', icon: GitPullRequest, href: GITHUB_URL, title: t('sponsor.help.codeTitle'), desc: t('sponsor.help.codeDesc') },
])

let copiedTimer = null
const copyAddress = async () => {
  if (await copyText(usdtAddress.value)) {
    toast.success(t('sponsor.copied'))
    copiedAddress.value = true
    clearTimeout(copiedTimer)
    copiedTimer = setTimeout(() => {
      copiedAddress.value = false
    }, 2000)
  } else {
    toast.error(t('sponsor.copyFailed'))
  }
}

const scrollToSupport = () => {
  document.getElementById('support')?.scrollIntoView({ behavior: 'smooth', block: 'start' })
}

const scrollToPartner = (type) => {
  if (type) partnerFormRef.value?.setType(type)
  document.getElementById('become-partner')?.scrollIntoView({ behavior: 'smooth', block: 'start' })
  // Move focus for keyboard and screen reader users without fighting the smooth scroll or opening a mobile keyboard.
  document.getElementById('partner-title')?.focus({ preventScroll: true })
}

const load = async () => {
  loading.value = true
  loadError.value = false
  try {
    const { data } = await axios.get('/api/public/sponsors')
    sponsors.value = Array.isArray(data?.sponsors) ? data.sponsors : []
    payment.value = {
      usdt_trc20_address: data?.payment?.usdt_trc20_address || '',
      github_sponsors_url: data?.payment?.github_sponsors_url || '',
    }
  } catch (err) {
    console.error('Failed to fetch sponsors:', err)
    loadError.value = true
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await load()
  // The router resets scroll to the top, so honour a #become-partner deep link once the content above has its final height.
  if (route.hash === '#become-partner') {
    await nextTick()
    document.getElementById('become-partner')?.scrollIntoView({ block: 'start' })
  }
})
</script>
