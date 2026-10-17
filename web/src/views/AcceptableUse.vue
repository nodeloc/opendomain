<template>
  <div>
    <!-- Title hero -->
    <section class="relative overflow-hidden border-b border-y-line">
      <div class="y-dot-grid pointer-events-none absolute inset-0" aria-hidden="true"></div>
      <div class="y-container relative flex flex-col gap-4 pb-14 pt-16 sm:pt-20">
        <p class="y-eyebrow">{{ $t('aup.eyebrow') }}</p>
        <h1 class="y-headline text-[clamp(2.5rem,6vw,4.5rem)]">{{ $t('aup.title') }}</h1>
        <p class="m-0 max-w-[640px] text-lead text-y-muted">{{ $t('aup.lead', { site: siteName }) }}</p>
        <div class="flex flex-wrap gap-x-6 gap-y-2 text-sm text-y-faint">
          <span>{{ $t('aup.effective', { date: $t('aup.effectiveDate') }) }}</span>
          <span aria-hidden="true">·</span>
          <span>{{ $t('aup.version') }}</span>
          <span aria-hidden="true">·</span>
          <i18n-t keypath="aup.contact" tag="span">
            <template #email>
              <a :href="`mailto:${abuseEmail}`" class="text-y-mint-text hover:text-y-mint-hover">{{ abuseEmail }}</a>
            </template>
          </i18n-t>
        </div>
      </div>
    </section>

    <div class="y-container flex flex-wrap items-start gap-x-16">
      <!-- TOC -->
      <nav class="w-full py-10 text-sm leading-5 lg:sticky lg:top-24 lg:w-56 lg:shrink-0" :aria-label="$t('aup.tocLabel')">
        <div class="y-label pb-2">{{ $t('aup.tocLabel') }}</div>
        <a
          v-for="item in toc"
          :key="item.id"
          :href="`#${item.id}`"
          class="block py-1.5 transition-colors"
          :class="activeSection === item.id ? 'font-medium text-y-mint-text' : 'text-y-muted hover:text-y-text'"
        >
          {{ item.n }}. {{ item.label }}
        </a>
      </nav>

      <main class="min-w-0 max-w-[760px] flex-1 pb-24 pt-10" style="flex-basis: 560px;">
        <!-- TL;DR cards -->
        <div class="grid gap-4 pb-8 sm:grid-cols-2">
          <div class="y-mint-soft flex flex-col gap-2.5 rounded-lg border p-5">
            <div class="flex items-center gap-2 font-bold text-y-mint-text">
              <Check :size="18" :stroke-width="2" aria-hidden="true" /> {{ $t('aup.goTitle') }}
            </div>
            <ul class="m-0 flex list-none flex-col gap-1.5 p-0 text-[15px] text-y-text">
              <li v-for="i in 4" :key="i" class="flex items-baseline gap-2.5">
                <span class="y-dot relative top-[-1px] bg-y-mint" aria-hidden="true"></span>{{ $t(`aup.go${i}`) }}
              </li>
            </ul>
          </div>
          <div class="y-danger-soft flex flex-col gap-2.5 rounded-lg border p-5">
            <div class="flex items-center gap-2 font-bold text-y-danger-text">
              <X :size="18" :stroke-width="2" aria-hidden="true" /> {{ $t('aup.neverTitle') }}
            </div>
            <ul class="m-0 flex list-none flex-col gap-1.5 p-0 text-[15px] text-y-text">
              <li v-for="i in 4" :key="i" class="flex items-baseline gap-2.5">
                <span class="y-dot relative top-[-1px] bg-y-danger" aria-hidden="true"></span>{{ $t(`aup.never${i}`) }}
              </li>
            </ul>
          </div>
        </div>

        <!-- 01 Who can use -->
        <section id="s1" class="doc-section scroll-mt-24">
          <h2 class="doc-h2"><span>01</span>{{ $t('aup.s1Title', { site: siteName }) }}</h2>
          <p class="doc-p">{{ $t('aup.s1Body') }}</p>
        </section>

        <!-- 02 Welcome -->
        <section id="s2" class="doc-section scroll-mt-24">
          <h2 class="doc-h2"><span>02</span>{{ $t('aup.s2Title') }}</h2>
          <ul class="doc-list">
            <li v-for="key in ['s2a', 's2b', 's2c', 's2d']" :key="key">{{ $t(`aup.${key}`) }}</li>
          </ul>
        </section>

        <!-- 03 Not allowed -->
        <section id="s3" class="doc-section scroll-mt-24">
          <h2 class="doc-h2"><span>03</span>{{ $t('aup.s3Title') }}</h2>
          <ul class="doc-list">
            <li v-for="item in notAllowed" :key="item.t">
              <strong class="font-semibold text-y-text">{{ $t(`aup.${item.t}`, { site: siteName }) }}</strong>
              {{ $t(`aup.${item.d}`) }}
            </li>
          </ul>
        </section>

        <!-- 04 Fair use -->
        <section id="s4" class="doc-section scroll-mt-24">
          <h2 class="doc-h2"><span>04</span>{{ $t('aup.s4Title') }}</h2>
          <ul class="doc-list">
            <li>{{ $t('aup.s4Quota') }}</li>
            <li>{{ $t('aup.s4Inactive') }}</li>
            <li>
              <i18n-t keypath="aup.s4Reserved" tag="span">
                <template #www><span class="font-mono text-y-text">www</span></template>
                <template #mail><span class="font-mono text-y-text">mail</span></template>
                <template #admin><span class="font-mono text-y-text">admin</span></template>
              </i18n-t>
            </li>
          </ul>
        </section>

        <!-- 05 Enforcement -->
        <section id="s5" class="doc-section scroll-mt-24">
          <h2 class="doc-h2"><span>05</span>{{ $t('aup.s5Title') }}</h2>
          <p class="doc-p">{{ $t('aup.s5Body') }}</p>
          <div class="grid gap-3 sm:grid-cols-3">
            <div class="y-panel p-4">
              <div class="font-mono text-xs text-y-warn">{{ $t('aup.step', { n: 1 }) }}</div>
              <div class="mt-1 font-semibold text-y-text">{{ $t('aup.s5WarnTitle') }}</div>
              <div class="mt-1 text-sm leading-[22px] text-y-muted">{{ $t('aup.s5WarnDesc') }}</div>
            </div>
            <div class="y-panel p-4">
              <div class="font-mono text-xs text-y-danger-text">{{ $t('aup.step', { n: 2 }) }}</div>
              <div class="mt-1 font-semibold text-y-text">{{ $t('aup.s5SuspendTitle') }}</div>
              <div class="mt-1 text-sm leading-[22px] text-y-muted">{{ $t('aup.s5SuspendDesc') }}</div>
            </div>
            <div class="y-panel p-4">
              <div class="font-mono text-xs text-y-danger">{{ $t('aup.step', { n: 3 }) }}</div>
              <div class="mt-1 font-semibold text-y-text">{{ $t('aup.s5RemoveTitle') }}</div>
              <div class="mt-1 text-sm leading-[22px] text-y-muted">{{ $t('aup.s5RemoveDesc') }}</div>
            </div>
          </div>
        </section>

        <!-- 06 Reporting -->
        <section id="s6" class="doc-section scroll-mt-24">
          <h2 class="doc-h2"><span>06</span>{{ $t('aup.s6Title') }}</h2>
          <i18n-t keypath="aup.s6Body" tag="p" class="doc-p">
            <template #form>
              <router-link to="/report" class="doc-link">{{ $t('aup.s6Form') }}</router-link>
            </template>
            <template #email>
              <a :href="`mailto:${abuseEmail}`" class="doc-link">{{ abuseEmail }}</a>
            </template>
          </i18n-t>
        </section>

        <!-- 07 Appeals -->
        <section id="s7" class="doc-section scroll-mt-24">
          <h2 class="doc-h2"><span>07</span>{{ $t('aup.s7Title') }}</h2>
          <i18n-t keypath="aup.s7Body" tag="p" class="doc-p">
            <template #email>
              <a :href="`mailto:${supportEmail}`" class="doc-link">{{ supportEmail }}</a>
            </template>
          </i18n-t>
        </section>

        <!-- 08 Changes -->
        <section id="s8" class="doc-section scroll-mt-24">
          <h2 class="doc-h2"><span>08</span>{{ $t('aup.s8Title') }}</h2>
          <i18n-t keypath="aup.s8Body" tag="p" class="doc-p">
            <template #announcements>
              <router-link to="/announcements" class="doc-link">{{ $t('aup.s8Announcements') }}</router-link>
            </template>
          </i18n-t>
        </section>

        <!-- Bottom CTA -->
        <div class="y-card mt-4 flex flex-wrap items-center justify-between gap-4 p-6">
          <div>
            <div class="text-lg font-bold text-y-text">{{ $t('aup.ctaTitle') }}</div>
            <div class="mt-1 text-sm text-y-muted">{{ $t('aup.ctaDesc') }}</div>
          </div>
          <router-link
            to="/report"
            class="inline-flex min-h-[44px] items-center rounded-md border border-y-danger-text/70 px-5 text-[15px] font-semibold text-y-danger-text transition-colors hover:bg-y-danger/10"
          >
            {{ $t('aup.ctaButton') }}
          </router-link>
        </div>
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useI18n } from 'vue-i18n'
import { Check, X } from 'lucide-vue-next'
import { useSiteConfigStore } from '../stores/siteConfig'
import { ABUSE_EMAIL, SUPPORT_EMAIL } from '../config/links'

const { t } = useI18n()
const siteConfigStore = useSiteConfigStore()

const siteName = computed(() => siteConfigStore.siteName || 'Yunamo')
const abuseEmail = ABUSE_EMAIL
const supportEmail = SUPPORT_EMAIL

const notAllowed = [
  { t: 's3PhishingT', d: 's3PhishingD' },
  { t: 's3MalwareT', d: 's3MalwareD' },
  { t: 's3FraudT', d: 's3FraudD' },
  { t: 's3SpamT', d: 's3SpamD' },
  { t: 's3ImpersonationT', d: 's3ImpersonationD' },
  { t: 's3IllegalT', d: 's3IllegalD' },
  { t: 's3SquattingT', d: 's3SquattingD' },
  { t: 's3AttacksT', d: 's3AttacksD' },
]

const toc = computed(() => [
  { id: 's1', n: 1, label: t('aup.s1Toc', { site: siteName.value }) },
  { id: 's2', n: 2, label: t('aup.s2Title') },
  { id: 's3', n: 3, label: t('aup.s3Title') },
  { id: 's4', n: 4, label: t('aup.s4Title') },
  { id: 's5', n: 5, label: t('aup.s5Title') },
  { id: 's6', n: 6, label: t('aup.s6Toc') },
  { id: 's7', n: 7, label: t('aup.s7Title') },
  { id: 's8', n: 8, label: t('aup.s8Toc') },
])

const activeSection = ref('s1')
let observer = null

onMounted(() => {
  document.documentElement.style.scrollBehavior = 'smooth'
  observer = new IntersectionObserver(
    (entries) => {
      for (const entry of entries) {
        if (entry.isIntersecting) activeSection.value = entry.target.id
      }
    },
    { rootMargin: '-15% 0px -70% 0px' }
  )
  for (const item of toc.value) {
    const el = document.getElementById(item.id)
    if (el) observer.observe(el)
  }
})

onBeforeUnmount(() => {
  document.documentElement.style.scrollBehavior = ''
  observer?.disconnect()
})
</script>

<style scoped>
.doc-section {
  display: flex;
  flex-direction: column;
  gap: 14px;
  border-top: 1px solid rgb(var(--y-line));
  padding: 32px 0;
  line-height: 26px;
}
.doc-h2 {
  margin: 0;
  display: flex;
  align-items: baseline;
  gap: 14px;
  font-size: 24px;
  line-height: 32px;
  font-weight: 700;
  letter-spacing: -0.02em;
  color: rgb(var(--y-text));
}
.doc-h2 span {
  font-family: 'Geist Mono', ui-monospace, monospace;
  font-size: 14px;
  font-weight: 500;
  color: rgb(var(--y-mint-text));
}
.doc-p {
  margin: 0;
  color: rgb(var(--y-muted));
}
.doc-list {
  margin: 0;
  padding-left: 20px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  color: rgb(var(--y-muted));
}
.doc-list li::marker {
  color: rgb(var(--y-mint-text));
}
.doc-link {
  color: rgb(var(--y-mint-text));
  text-decoration: underline;
  text-decoration-color: rgb(var(--y-mint-text) / 0.4);
  text-underline-offset: 3px;
}
.doc-link:hover {
  color: rgb(var(--y-mint-hover));
}
</style>
