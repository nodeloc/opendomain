<template>
  <div class="y-container py-14 sm:py-16">
    <div class="mx-auto max-w-3xl">
      <router-link to="/announcements" class="btn btn-ghost btn-sm -ml-3 gap-1.5 text-y-muted">
        <ArrowLeft :size="16" :stroke-width="1.5" />
        {{ $t('announcement.back') }}
      </router-link>

      <YLoader v-if="loading" class="mt-6" :label="$t('common.loading')" />

      <div v-else-if="error" class="y-panel mt-6 flex flex-col items-center px-6 py-14 text-center">
        <span class="inline-flex h-10 w-10 items-center justify-center rounded-md border border-y-line bg-y-surface-2 text-y-danger-text">
          <CircleX :size="20" :stroke-width="1.5" />
        </span>
        <h1 class="y-h3 mt-5">{{ $t('uiPublic.announcements.loadFailedTitle') }}</h1>
        <p class="mt-2 max-w-md text-sm leading-6 text-y-muted">{{ error }}</p>
        <router-link to="/announcements" class="btn btn-primary mt-6">
          {{ $t('announcement.back') }}
        </router-link>
      </div>

      <article v-else-if="announcement" class="mt-6">
        <header>
          <div class="flex flex-wrap items-center gap-2">
            <span :class="getTypeBadgeClass(announcement.type)">
              <component :is="getTypeIcon(announcement.type)" :size="12" :stroke-width="1.75" />
              {{ formatType(announcement.type) }}
            </span>
            <span v-if="announcement.priority > 0" class="y-chip-warn">
              <Star :size="12" :stroke-width="1.75" />
              {{ $t('announcement.priority') }}
            </span>
          </div>

          <h1 class="y-h1 mt-5 break-words">{{ announcement.title }}</h1>

          <div class="mt-5 flex flex-wrap items-center gap-x-5 gap-y-2 text-sm text-y-faint">
            <span class="inline-flex items-center gap-1.5">
              <Calendar :size="16" :stroke-width="1.5" />
              {{ formatDate(announcement.published_at || announcement.created_at) }}
            </span>
            <span v-if="announcement.author_name" class="inline-flex items-center gap-1.5">
              <User :size="16" :stroke-width="1.5" />
              {{ announcement.author_name }}
            </span>
            <span class="inline-flex items-center gap-1.5">
              <Eye :size="16" :stroke-width="1.5" />
              {{ $t('uiPublic.announcements.viewsCount', { count: announcement.views || 0 }) }}
            </span>
          </div>
        </header>

        <div class="y-panel mt-8 p-6 sm:p-10">
          <div class="prose max-w-none sm:prose-lg" v-html="formatContent(announcement.content)"></div>
        </div>
      </article>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ArrowLeft, Calendar, CircleX, Eye, Megaphone, Sparkles, Star, TriangleAlert, User, Wrench } from 'lucide-vue-next'
import axios from '../utils/axios'
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import YLoader from '../components/brand/YLoader.vue'

const route = useRoute()
const router = useRouter()
const { t, te } = useI18n()

const announcement = ref(null)
const loading = ref(true)
const error = ref(null)

onMounted(async () => {
  await fetchAnnouncement()
})

const fetchAnnouncement = async () => {
  loading.value = true
  error.value = null
  try {
    const response = await axios.get(`/api/public/announcements/${route.params.id}`)
    announcement.value = response.data
  } catch (err) {
    console.error('Failed to fetch announcement:', err)
    error.value = err.response?.data?.error || 'Failed to load announcement'
  } finally {
    loading.value = false
  }
}

const getTypeBadgeClass = (type) => {
  const classes = {
    general: 'y-chip',
    maintenance: 'y-chip-warn',
    update: 'y-chip-mint',
    important: 'y-chip-danger',
  }
  return classes[type] || 'y-chip'
}

const getTypeIcon = (type) => {
  const icons = {
    general: Megaphone,
    maintenance: Wrench,
    update: Sparkles,
    important: TriangleAlert,
  }
  return icons[type] || Megaphone
}

const formatType = (type) => {
  const key = `announcement.${type}`
  return te(key) ? t(key) : type
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

const formatContent = (content) => {
  if (!content) return ''
  return DOMPurify.sanitize(marked.parse(content))
}
</script>
