<template>
  <div class="flex w-full items-start gap-12 py-8 lg:py-10">
    <YLoader v-if="loading" class="min-h-[50vh] flex-1" :label="$t('page.loading')" />

    <div v-else-if="error" class="mx-auto w-full max-w-xl py-10">
      <div class="y-panel flex flex-col items-center px-6 py-14 text-center">
        <span class="inline-flex h-10 w-10 items-center justify-center rounded-md border border-y-line bg-y-surface-2 text-y-danger-text">
          <CircleX :size="20" :stroke-width="1.5" />
        </span>
        <h1 class="y-h3 mt-5">{{ notFound ? $t('docs.notFoundTitle') : $t('uiPublic.page.errorTitle') }}</h1>
        <p class="mt-2 max-w-md text-sm leading-6 text-y-muted">{{ notFound ? $t('docs.notFoundDesc') : error }}</p>
        <router-link :to="notFound ? DOCS_PATH : '/'" class="btn btn-primary mt-6">
          {{ notFound ? $t('docs.backToDocs') : $t('uiPublic.backHome') }}
        </router-link>
      </div>
    </div>

    <template v-else-if="page">
      <article class="min-w-0 max-w-[760px] flex-1 pb-16">
        <nav class="flex flex-wrap items-center gap-2 text-[13px] text-y-faint" :aria-label="$t('docs.crumbRoot')">
          <router-link :to="DOCS_PATH" class="transition-colors hover:text-y-text">{{ $t('docs.crumbRoot') }}</router-link>
          <template v-if="groupName">
            <span aria-hidden="true">/</span>
            <span>{{ groupName }}</span>
          </template>
          <span aria-hidden="true">/</span>
          <span class="text-y-muted">{{ page.title }}</span>
        </nav>

        <header class="mt-7 flex flex-col gap-3.5">
          <h1 class="break-words text-[clamp(2rem,5vw,2.75rem)] font-extrabold leading-[1.05] tracking-[-0.035em] text-y-text">
            {{ page.title }}
          </h1>
          <p class="flex flex-wrap items-center gap-x-4 gap-y-1 text-[13px] text-y-faint">
            <span>{{ $t('docs.updated', { date: formatDate(page.updated_at) }) }}</span>
            <span aria-hidden="true">·</span>
            <span>{{ $t('docs.readTime', { n: readMinutes }) }}</span>
          </p>
        </header>

        <div
          ref="bodyRef"
          class="docs-body prose mt-8 max-w-none"
          v-html="rendered"
          @click="onBodyClick"
        ></div>

        <div v-if="prevNext.prev || prevNext.next" class="mt-12 flex flex-wrap gap-4 border-t border-y-line pt-7">
          <router-link v-if="prevNext.prev" :to="prevNext.prev.to" class="docs-pn">
            <span class="flex items-center gap-1.5 text-[13px] text-y-faint">
              <ArrowLeft :size="14" :stroke-width="1.5" />
              {{ $t('docs.previous') }}
            </span>
            <span class="mt-1 block font-semibold text-y-text">{{ prevNext.prev.title }}</span>
          </router-link>
          <span v-else class="hidden flex-1 sm:block" aria-hidden="true"></span>
          <router-link v-if="prevNext.next" :to="prevNext.next.to" class="docs-pn text-right">
            <span class="flex items-center justify-end gap-1.5 text-[13px] text-y-faint">
              {{ $t('docs.next') }}
              <ArrowRight :size="14" :stroke-width="1.5" />
            </span>
            <span class="mt-1 block font-semibold text-y-text">{{ prevNext.next.title }}</span>
          </router-link>
        </div>

        <div class="mt-8 flex flex-wrap items-center justify-between gap-x-6 gap-y-3 text-sm text-y-faint">
          <p v-if="helpful" class="text-y-muted">{{ $t('docs.thanks') }}</p>
          <p v-else>
            {{ $t('docs.helpful') }}
            <button type="button" class="ml-1 font-medium text-y-mint-text hover:underline" @click="helpful = 'yes'">{{ $t('docs.yes') }}</button>
            <span aria-hidden="true"> · </span>
            <button type="button" class="font-medium text-y-mint-text hover:underline" @click="helpful = 'no'">{{ $t('docs.no') }}</button>
          </p>
          <p>
            {{ $t('docs.needHelp') }}
            <a :href="TELEGRAM_GROUP_URL" target="_blank" rel="noopener" class="font-medium text-y-mint-text hover:underline">{{ $t('docs.askTelegram') }}</a>
          </p>
        </div>
      </article>

      <nav v-if="toc.length" class="sticky top-24 hidden w-[200px] shrink-0 flex-col gap-2.5 py-2 xl:flex" :aria-label="$t('docs.onThisPage')">
        <div class="mb-1 font-mono text-[11px] uppercase tracking-[0.1em] text-y-faint">{{ $t('docs.onThisPage') }}</div>
        <a
          v-for="h in toc"
          :key="h.id"
          :href="`#${h.id}`"
          class="block text-[13px] leading-5 transition-colors"
          :class="[h.level === 3 ? 'pl-4' : '', activeId === h.id ? 'text-y-mint-text' : 'text-y-faint hover:text-y-text']"
          @click.prevent="scrollToHeading(h.id)"
        >{{ h.text }}</a>
      </nav>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { CircleX, ArrowLeft, ArrowRight } from 'lucide-vue-next'
import axios from '../utils/axios'
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import YLoader from '../components/brand/YLoader.vue'
import { useDocsPages } from '../composables/useDocsPages'
import { TELEGRAM_GROUP_URL, DOCS_PATH } from '../config/links'

const route = useRoute()
const { t, locale } = useI18n()
const { load, findBySlug, neighbours } = useDocsPages()

const page = ref(null)
const loading = ref(true)
const error = ref(null)
const notFound = ref(false)
const helpful = ref(null)

const bodyRef = ref(null)
const toc = ref([])
const activeId = ref('')

const groupName = computed(() => findBySlug(route.params.slug)?.group || '')
const prevNext = computed(() => neighbours(route.params.slug))

const rendered = computed(() => {
  if (!page.value?.content) return ''
  return DOMPurify.sanitize(marked.parse(page.value.content))
})

const readMinutes = computed(() => {
  const text = (page.value?.content || '').replace(/[#>*`|\-\[\]()]/g, ' ').replace(/\s+/g, ' ').trim()
  return Math.max(1, Math.round(text.length / 1200))
})

const formatDate = (date) => {
  if (!date) return ''
  return new Date(date).toLocaleDateString(locale.value === 'zh-CN' ? 'zh-CN' : 'en-GB', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
  })
}

const slugify = (text) =>
  text
    .toLowerCase()
    .trim()
    .replace(/[^\p{L}\p{N}]+/gu, '-')
    .replace(/^-+|-+$/g, '')

const makeCopyButton = () => {
  const btn = document.createElement('button')
  btn.type = 'button'
  btn.className = 'docs-copy'
  btn.textContent = t('docs.copy')
  return btn
}

// Post-process the sanitized HTML: heading ids for the TOC, copy buttons on code
// blocks and tables, and callout styling for blockquotes.
const enhance = async () => {
  await nextTick()
  const root = bodyRef.value
  if (!root) return

  // A markdown body that starts with an H1 repeating the page title would
  // duplicate the header we already render above the content.
  const first = root.firstElementChild
  if (
    first?.tagName === 'H1' &&
    first.textContent.trim().toLowerCase() === (page.value?.title || '').trim().toLowerCase()
  ) {
    first.remove()
  }

  const used = new Set()
  toc.value = [...root.querySelectorAll('h2, h3')].map((h, i) => {
    const base = slugify(h.textContent) || `section-${i + 1}`
    let id = base
    let n = 2
    while (used.has(id)) id = `${base}-${n++}`
    used.add(id)
    h.id = id
    return { id, text: h.textContent, level: h.tagName === 'H3' ? 3 : 2 }
  })

  root.querySelectorAll('pre').forEach((pre) => {
    if (pre.parentElement?.classList.contains('docs-pre')) return
    const wrap = document.createElement('div')
    wrap.className = 'docs-pre'
    pre.parentNode.insertBefore(wrap, pre)
    wrap.appendChild(pre)
    wrap.appendChild(makeCopyButton())
  })

  root.querySelectorAll('table').forEach((table) => {
    if (table.closest('.docs-table-card')) return
    const card = document.createElement('div')
    card.className = 'docs-table-card'
    table.parentNode.insertBefore(card, table)
    const head = document.createElement('div')
    head.className = 'docs-table-head'
    head.appendChild(makeCopyButton())
    const scroller = document.createElement('div')
    scroller.className = 'docs-table-scroll'
    card.appendChild(head)
    card.appendChild(scroller)
    scroller.appendChild(table)
  })

  root.querySelectorAll('blockquote').forEach((q) => {
    const warn = /^\s*\[?!?\s*(warn|attention|注意|警告|caution|danger)/i.test(q.textContent)
    q.classList.add(warn ? 'docs-quote-warn' : 'docs-quote-info')
  })

  updateScrollspy()
}

const onBodyClick = (e) => {
  const btn = e.target.closest('.docs-copy')
  if (!btn) return
  const pre = btn.closest('.docs-pre')
  const card = btn.closest('.docs-table-card')
  let text = ''
  if (pre) {
    text = pre.querySelector('pre')?.innerText || ''
  } else if (card) {
    text = [...card.querySelectorAll('tr')]
      .map((tr) => [...tr.cells].map((c) => c.innerText).join('\t'))
      .join('\n')
  }
  if (text) navigator.clipboard?.writeText(text).catch(() => {})
  btn.textContent = t('docs.copied')
  btn.classList.add('is-copied')
  setTimeout(() => {
    btn.textContent = t('docs.copy')
    btn.classList.remove('is-copied')
  }, 1500)
}

const scrollToHeading = (id) => {
  document.getElementById(id)?.scrollIntoView({ behavior: 'smooth', block: 'start' })
  history.replaceState(null, '', `#${id}`)
  activeId.value = id
}

const updateScrollspy = () => {
  const root = bodyRef.value
  if (!root) return
  const heads = [...root.querySelectorAll('h2[id], h3[id]')]
  let current = heads[0]?.id || ''
  for (const h of heads) {
    if (h.getBoundingClientRect().top <= 120) current = h.id
  }
  activeId.value = current
}

let ticking = false
const onScroll = () => {
  if (ticking) return
  ticking = true
  requestAnimationFrame(() => {
    updateScrollspy()
    ticking = false
  })
}

const fetchPage = async () => {
  try {
    loading.value = true
    error.value = null
    notFound.value = false
    page.value = null
    helpful.value = null
    toc.value = []
    const response = await axios.get(`/api/public/pages/${route.params.slug}`)
    page.value = response.data
    loading.value = false // body must render before enhance() can walk it
    await enhance()
  } catch (err) {
    notFound.value = err.response?.status === 404
    error.value = err.response?.data?.error || t('page.loadFailed')
  } finally {
    loading.value = false
  }
}

watch(() => route.params.slug, (slug, prev) => {
  if (slug && slug !== prev) fetchPage()
})

onMounted(() => {
  load()
  fetchPage()
  window.addEventListener('scroll', onScroll, { passive: true })
})

onBeforeUnmount(() => {
  window.removeEventListener('scroll', onScroll)
})
</script>

<style scoped>
.docs-body {
  counter-reset: docs-h2;
  font-size: 15px;
  line-height: 26px;
}
.docs-body :deep(h1) {
  margin: 32px 0 0;
  font-size: 22px;
  font-weight: 700;
  line-height: 30px;
  letter-spacing: -0.015em;
  color: rgb(var(--y-text));
  scroll-margin-top: 96px;
}
.docs-body :deep(h2) {
  counter-increment: docs-h2;
  display: flex;
  align-items: center;
  gap: 12px;
  margin: 40px 0 0;
  font-size: 26px;
  font-weight: 700;
  line-height: 32px;
  letter-spacing: -0.02em;
  color: rgb(var(--y-text));
  scroll-margin-top: 96px;
}
.docs-body :deep(h2)::before {
  content: counter(docs-h2);
  display: inline-flex;
  width: 28px;
  height: 28px;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  border-radius: 999px;
  background: rgb(var(--y-mint));
  color: rgb(var(--y-on-mint));
  font-family: 'Geist Mono', ui-monospace, monospace;
  font-size: 14px;
  font-weight: 500;
}
.docs-body :deep(h3) {
  margin: 28px 0 0;
  font-size: 19px;
  font-weight: 600;
  line-height: 26px;
  letter-spacing: -0.01em;
  color: rgb(var(--y-text));
  scroll-margin-top: 96px;
}
.docs-body :deep(p),
.docs-body :deep(ul),
.docs-body :deep(ol) {
  margin: 16px 0 0;
}
.docs-body :deep(li) {
  margin-top: 6px;
}
.docs-body :deep(ul > li)::marker {
  color: rgb(var(--y-mint));
}
.docs-body :deep(ol > li)::marker {
  color: rgb(var(--y-mint-text));
  font-family: 'Geist Mono', ui-monospace, monospace;
}
.docs-body :deep(a) {
  color: rgb(var(--y-mint-text));
}
.docs-body :deep(hr) {
  margin: 32px 0;
  border-color: rgb(var(--y-line));
}
.docs-body :deep(code) {
  font-size: 0.875em;
}
.docs-body :deep(code)::before,
.docs-body :deep(code)::after {
  content: none;
}

/* Code blocks: recessed card with a copy button. */
.docs-body :deep(.docs-pre) {
  position: relative;
  margin: 20px 0 0;
}
.docs-body :deep(.docs-pre pre) {
  margin: 0;
  border: 1px solid rgb(var(--y-line));
  border-radius: 12px;
  background: rgb(var(--y-bg));
  padding: 14px 16px;
  padding-right: 72px;
  font-size: 14px;
  line-height: 22px;
  overflow-x: auto;
}
.docs-body :deep(.docs-copy) {
  position: absolute;
  top: 10px;
  right: 10px;
  min-height: 30px;
  cursor: pointer;
  border: 1px solid rgb(var(--y-line-strong));
  border-radius: 6px;
  background: transparent;
  padding: 4px 10px;
  font-family: 'Geist Mono', ui-monospace, monospace;
  font-size: 12px;
  font-weight: 500;
  line-height: 16px;
  color: rgb(var(--y-muted));
  transition: background-color 150ms, color 150ms;
}
.docs-body :deep(.docs-copy:hover) {
  background: rgb(var(--y-surface-2));
  color: rgb(var(--y-text));
}
.docs-body :deep(.docs-copy.is-copied) {
  color: rgb(var(--y-mint-text));
  border-color: rgb(var(--y-mint) / 0.5);
}

/* Tables: card with a copy header row. */
.docs-body :deep(.docs-table-card) {
  margin: 20px 0 0;
  overflow: hidden;
  border: 1px solid rgb(var(--y-line));
  border-radius: 12px;
  background: rgb(var(--y-surface));
}
.docs-body :deep(.docs-table-head) {
  position: relative;
  display: flex;
  min-height: 46px;
  align-items: center;
  justify-content: flex-end;
  border-bottom: 1px solid rgb(var(--y-line));
  padding: 8px;
}
.docs-body :deep(.docs-table-head .docs-copy) {
  position: static;
}
.docs-body :deep(.docs-table-scroll) {
  overflow-x: auto;
}
.docs-body :deep(.docs-table-card table) {
  width: 100%;
  margin: 0;
  border-collapse: collapse;
  font-family: 'Geist Mono', ui-monospace, monospace;
  font-size: 14px;
  line-height: 20px;
}
.docs-body :deep(.docs-table-card thead th) {
  border: 0;
  padding: 12px 16px;
  text-align: left;
  font-size: 12px;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: rgb(var(--y-faint));
  white-space: nowrap;
}
.docs-body :deep(.docs-table-card tbody td) {
  border: 0;
  border-top: 1px solid rgb(var(--y-line));
  padding: 12px 16px;
  color: rgb(var(--y-text));
}
.docs-body :deep(.docs-table-card tbody td:first-child) {
  color: rgb(var(--y-mint-text));
}
.docs-body :deep(.docs-table-card thead) {
  border: 0;
}

/* Blockquotes become callouts. */
.docs-body :deep(blockquote) {
  margin: 20px 0 0;
  border: 1px solid;
  border-radius: 12px;
  padding: 14px 16px;
  font-style: normal;
  font-weight: 400;
}
.docs-body :deep(blockquote.docs-quote-info) {
  background: rgb(var(--y-mint) / 0.07);
  border-color: rgb(var(--y-mint) / 0.3);
}
.docs-body :deep(blockquote.docs-quote-warn) {
  background: rgb(var(--y-warn) / 0.08);
  border-color: rgb(var(--y-warn) / 0.35);
}
.docs-body :deep(blockquote p) {
  margin: 0;
  color: rgb(var(--y-text));
}
.docs-body :deep(blockquote p)::before,
.docs-body :deep(blockquote p)::after {
  content: none;
}
.docs-body :deep(blockquote p + p) {
  margin-top: 8px;
}

.docs-pn {
  flex: 1 1 240px;
  border: 1px solid rgb(var(--y-line));
  border-radius: 12px;
  padding: 16px 20px;
  transition: border-color 150ms;
}
.docs-pn:hover {
  border-color: rgb(var(--y-line-strong));
}
</style>
