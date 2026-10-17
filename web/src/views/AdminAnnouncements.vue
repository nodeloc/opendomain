<template>
  <div class="y-page">
    <PageHeader
      :eyebrow="$t('uiAdmin2.eyebrow')"
      :title="$t('uiAdmin2.announcements.title')"
      :description="$t('uiAdmin2.announcements.description')"
    >
      <template #actions>
        <button type="button" class="btn btn-primary gap-1.5" @click="showCreateModal = true">
          <Plus :size="16" :stroke-width="1.5" />
          {{ $t('uiAdmin2.announcements.create') }}
        </button>
      </template>
    </PageHeader>

    <section class="y-panel">
      <YLoader v-if="loading" :label="$t('uiAdmin2.loading')" />

      <EmptyState
        v-else-if="announcements.length === 0"
        :title="$t('uiAdmin2.announcements.emptyTitle')"
        :description="$t('uiAdmin2.announcements.emptyDescription')"
      >
        <button type="button" class="btn btn-outline btn-sm gap-1.5" @click="showCreateModal = true">
          <Plus :size="14" :stroke-width="1.5" />
          {{ $t('uiAdmin2.announcements.create') }}
        </button>
      </EmptyState>

      <div v-else class="overflow-x-auto">
        <table class="y-table">
          <thead>
            <tr>
              <th>{{ $t('uiAdmin2.announcements.table.title') }}</th>
              <th>{{ $t('uiAdmin2.announcements.table.type') }}</th>
              <th class="text-right">{{ $t('uiAdmin2.announcements.table.priority') }}</th>
              <th class="text-right">{{ $t('uiAdmin2.announcements.table.views') }}</th>
              <th>{{ $t('uiAdmin2.announcements.table.status') }}</th>
              <th>{{ $t('uiAdmin2.announcements.table.date') }}</th>
              <th class="text-right">{{ $t('uiAdmin2.announcements.table.actions') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="announcement in announcements" :key="announcement.id">
              <td class="max-w-md">
                <p class="truncate font-medium text-y-text">{{ announcement.title }}</p>
              </td>
              <td class="whitespace-nowrap">
                <span :class="getTypeBadgeClass(announcement.type)">{{ formatType(announcement.type) }}</span>
              </td>
              <td class="text-right font-mono">{{ announcement.priority }}</td>
              <td class="text-right font-mono text-y-muted">{{ announcement.views }}</td>
              <td class="whitespace-nowrap">
                <StatusDot
                  :tone="announcement.is_published ? 'ok' : 'muted'"
                  :label="announcement.is_published ? $t('uiAdmin2.announcements.published') : $t('uiAdmin2.announcements.draft')"
                />
              </td>
              <td class="whitespace-nowrap text-y-muted">{{ formatDate(announcement.created_at) }}</td>
              <td>
                <div class="flex justify-end gap-1">
                  <button
                    type="button"
                    class="btn btn-ghost btn-sm btn-square"
                    :title="$t('uiAdmin2.edit')"
                    :aria-label="$t('uiAdmin2.edit')"
                    @click="editAnnouncement(announcement)"
                  >
                    <Pencil :size="16" :stroke-width="1.5" />
                  </button>
                  <button
                    type="button"
                    class="btn btn-ghost btn-sm btn-square text-y-danger-text"
                    :title="$t('uiAdmin2.delete')"
                    :aria-label="$t('uiAdmin2.delete')"
                    @click="confirmDelete(announcement)"
                  >
                    <Trash2 :size="16" :stroke-width="1.5" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>

    <!-- Create/Edit Modal -->
    <dialog class="modal" :class="{ 'modal-open': showCreateModal || showEditModal }">
      <div class="modal-box relative w-11/12 max-w-3xl">
        <button
          type="button"
          class="btn btn-ghost btn-sm btn-square absolute right-3 top-3"
          :aria-label="$t('uiAdmin2.close')"
          @click="closeModal"
        >
          <X :size="16" :stroke-width="1.5" />
        </button>
        <h3 class="text-lg font-semibold text-y-text">
          {{ showEditModal ? $t('uiAdmin2.announcements.editTitle') : $t('uiAdmin2.announcements.createTitle') }}
        </h3>
        <p v-if="!showEditModal" class="mt-1 text-sm text-y-muted">{{ $t('uiAdmin2.announcements.modalDescription') }}</p>

        <form class="mt-6 space-y-4" @submit.prevent="handleSubmit">
          <label class="y-field">
            <span class="y-field-label">{{ $t('uiAdmin2.announcements.fields.title') }}</span>
            <input
              v-model="form.title"
              type="text"
              :placeholder="$t('uiAdmin2.announcements.fields.titlePlaceholder')"
              class="input input-bordered w-full"
              required
            />
          </label>

          <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <label class="y-field">
              <span class="y-field-label">{{ $t('uiAdmin2.announcements.fields.type') }}</span>
              <select v-model="form.type" class="select select-bordered w-full" required>
                <option value="general">{{ $t('announcement.general') }}</option>
                <option value="maintenance">{{ $t('announcement.maintenance') }}</option>
                <option value="update">{{ $t('announcement.update') }}</option>
                <option value="important">{{ $t('announcement.important') }}</option>
              </select>
            </label>

            <label class="y-field">
              <span class="y-field-label">{{ $t('uiAdmin2.announcements.fields.priority') }}</span>
              <input
                v-model.number="form.priority"
                type="number"
                min="0"
                max="100"
                class="input input-bordered w-full font-mono"
              />
              <span class="y-field-hint">{{ $t('uiAdmin2.announcements.fields.priorityHint') }}</span>
            </label>
          </div>

          <label class="y-field">
            <span class="y-field-label">{{ $t('uiAdmin2.announcements.fields.content') }}</span>
            <textarea
              v-model="form.content"
              class="textarea textarea-bordered w-full font-mono text-sm leading-6"
              rows="10"
              :placeholder="$t('uiAdmin2.announcements.fields.contentPlaceholder')"
              required
            ></textarea>
            <span class="y-field-hint">{{ $t('uiAdmin2.announcements.fields.contentHint') }}</span>
          </label>

          <label v-if="showEditModal" class="flex cursor-pointer items-start gap-3 rounded-md border border-y-line px-4 py-3">
            <input v-model="form.is_published" type="checkbox" class="checkbox checkbox-sm mt-0.5" />
            <span class="flex flex-col gap-0.5">
              <span class="y-field-label">{{ $t('uiAdmin2.announcements.fields.published') }}</span>
              <span class="y-field-hint">{{ $t('uiAdmin2.announcements.fields.publishedHint') }}</span>
            </span>
          </label>

          <div class="modal-action">
            <button type="button" class="btn btn-ghost" @click="closeModal">{{ $t('uiAdmin2.cancel') }}</button>
            <button type="submit" class="btn btn-primary" :disabled="submitting">
              <span v-if="submitting" class="loading loading-spinner loading-sm"></span>
              {{ showEditModal ? $t('uiAdmin2.saveChanges') : $t('uiAdmin2.create') }}
            </button>
          </div>
        </form>
      </div>
      <form method="dialog" class="modal-backdrop" @click="closeModal">
        <button>{{ $t('uiAdmin2.close') }}</button>
      </form>
    </dialog>
  </div>
</template>
<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import axios from '../utils/axios'
import { useToast } from '../composables/useToast'
import { Plus, Pencil, Trash2, X } from 'lucide-vue-next'
import PageHeader from '../components/ui/PageHeader.vue'
import StatusDot from '../components/ui/StatusDot.vue'
import EmptyState from '../components/ui/EmptyState.vue'
import YLoader from '../components/brand/YLoader.vue'

const { t } = useI18n()
const toast = useToast()
const announcements = ref([])
const loading = ref(true)
const showCreateModal = ref(false)
const showEditModal = ref(false)
const submitting = ref(false)
const editingAnnouncement = ref(null)

const form = ref({
  title: '',
  content: '',
  type: 'general',
  priority: 0,
  is_published: false,
})

onMounted(async () => {
  await fetchAnnouncements()
})

const fetchAnnouncements = async () => {
  loading.value = true
  try {
    const response = await axios.get('/api/admin/announcements')
    announcements.value = response.data.announcements
  } catch (error) {
    console.error('Failed to fetch announcements:', error)
  } finally {
    loading.value = false
  }
}

const handleSubmit = async () => {
  submitting.value = true
  try {
    if (showEditModal.value) {
      await axios.put(`/api/admin/announcements/${editingAnnouncement.value.id}`, form.value)
      toast.success(t('uiAdmin2.announcements.updateSuccess'))
    } else {
      await axios.post('/api/admin/announcements', form.value)
      toast.success(t('uiAdmin2.announcements.createSuccess'))
    }
    closeModal()
    await fetchAnnouncements()
  } catch (error) {
    toast.error(error.response?.data?.error || t('uiAdmin2.announcements.operationFailed'))
  } finally {
    submitting.value = false
  }
}

const editAnnouncement = (announcement) => {
  editingAnnouncement.value = announcement
  form.value = {
    title: announcement.title,
    content: announcement.content,
    type: announcement.type,
    priority: announcement.priority,
    is_published: announcement.is_published,
  }
  showEditModal.value = true
}

const confirmDelete = async (announcement) => {
  if (confirm(t('uiAdmin2.announcements.confirmDelete', { title: announcement.title }))) {
    try {
      await axios.delete(`/api/admin/announcements/${announcement.id}`)
      toast.success(t('uiAdmin2.announcements.deleteSuccess'))
      await fetchAnnouncements()
    } catch (error) {
      toast.error(error.response?.data?.error || t('uiAdmin2.announcements.deleteFailed'))
    }
  }
}

const closeModal = () => {
  showCreateModal.value = false
  showEditModal.value = false
  editingAnnouncement.value = null
  form.value = {
    title: '',
    content: '',
    type: 'general',
    priority: 0,
    is_published: false,
  }
}

const getTypeBadgeClass = (type) => {
  const classes = {
    general: 'y-chip',
    maintenance: 'y-chip-warn',
    update: 'y-chip-info',
    important: 'y-chip-danger',
  }
  return classes[type] || 'y-chip'
}

const formatType = (type) => {
  const types = {
    general: 'announcement.general',
    maintenance: 'announcement.maintenance',
    update: 'announcement.update',
    important: 'announcement.important',
  }
  return types[type] ? t(types[type]) : type
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}
</script>
