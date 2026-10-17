import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import axios from '../utils/axios'

// Module-level cache so the docs layout (sidebar, search) and the page view share one request.
const pages = ref([])
const loaded = ref(false)
const loading = ref(false)
const failed = ref(false)
let inflight = null

// Known categories first, in reading order; anything else follows alphabetically.
const CATEGORY_ORDER = ['resources', 'docs', 'guides', 'help', 'company', 'legal']

export const loadDocsPages = (force = false) => {
  if (loaded.value && !force) return Promise.resolve(pages.value)
  if (!inflight) {
    loading.value = true
    inflight = axios
      .get('/api/public/pages')
      .then((res) => {
        pages.value = Array.isArray(res.data?.pages) ? res.data.pages : []
        loaded.value = true
        failed.value = false
      })
      .catch(() => {
        failed.value = true
      })
      .finally(() => {
        loading.value = false
        inflight = null
      })
  }
  return inflight.then(() => pages.value)
}

const categoryRank = (c) => {
  const i = CATEGORY_ORDER.indexOf(c)
  return i === -1 ? CATEGORY_ORDER.length : i
}

export function useDocsPages() {
  const { t, te } = useI18n()

  const groupLabel = (category) => (te(`docs.groups.${category}`) ? t(`docs.groups.${category}`) : category)

  const groups = computed(() => {
    const byCategory = new Map()
    for (const p of pages.value) {
      const key = p.category || 'docs'
      if (!byCategory.has(key)) byCategory.set(key, [])
      byCategory.get(key).push(p)
    }
    return [...byCategory.entries()]
      .sort(([a], [b]) => categoryRank(a) - categoryRank(b) || a.localeCompare(b))
      .map(([category, items]) => ({
        key: category,
        label: groupLabel(category),
        items: items
          .slice()
          .sort((a, b) => (a.display_order ?? 0) - (b.display_order ?? 0) || a.id - b.id)
          .map((p) => ({ slug: p.slug, title: p.title, to: `/pages/${p.slug}`, group: groupLabel(category) })),
      }))
  })

  const policies = computed(() => ({
    key: 'policies',
    label: t('docs.groups.policies'),
    items: [
      { title: t('shell.acceptableUse'), to: '/acceptable-use', group: t('docs.groups.policies') },
      { title: t('shell.reportAbuse'), to: '/report', group: t('docs.groups.policies') },
    ],
  }))

  const sidebarGroups = computed(() => [...groups.value, policies.value])

  // Prev/next walk the database pages only; the policy links leave the docs layout.
  const orderedPages = computed(() => groups.value.flatMap((g) => g.items))

  const searchable = computed(() => sidebarGroups.value.flatMap((g) => g.items))

  const findBySlug = (slug) => orderedPages.value.find((p) => p.slug === slug) || null

  const neighbours = (slug) => {
    const list = orderedPages.value
    const i = list.findIndex((p) => p.slug === slug)
    if (i === -1) return { prev: null, next: null }
    return { prev: list[i - 1] || null, next: list[i + 1] || null }
  }

  return {
    pages,
    loaded,
    loading,
    failed,
    load: loadDocsPages,
    groups,
    sidebarGroups,
    orderedPages,
    searchable,
    groupLabel,
    findBySlug,
    neighbours,
  }
}
