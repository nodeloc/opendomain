<template>
  <div class="y-page">
    <PageHeader
      :eyebrow="$t('uiAdmin2.eyebrow')"
      :title="$t('uiAdmin2.pages.title')"
      :description="$t('uiAdmin2.pages.description')"
    >
      <template #actions>
        <button type="button" class="btn btn-primary gap-1.5" @click="openCreateModal">
          <Plus :size="16" :stroke-width="1.5" />
          {{ $t('uiAdmin2.pages.create') }}
        </button>
      </template>
    </PageHeader>

    <section class="y-panel">
      <div class="y-panel-header">
        <div class="max-w-full overflow-x-auto">
          <div class="tabs tabs-boxed tabs-sm w-max" role="tablist" :aria-label="$t('uiAdmin2.pages.filterLabel')">
            <button
              v-for="option in categoryFilters"
              :key="option.value || 'all'"
              type="button"
              role="tab"
              class="tab"
              :class="{ 'tab-active': filterCategory === option.value }"
              :aria-selected="filterCategory === option.value"
              @click="filterCategory = option.value"
            >
              {{ option.label }}
            </button>
          </div>
        </div>
        <span class="font-mono text-xs text-y-faint">{{ filteredPages.length }} / {{ pages.length }}</span>
      </div>

      <YLoader v-if="loading" :label="$t('uiAdmin2.loading')" />

      <EmptyState
        v-else-if="filteredPages.length === 0"
        :title="$t('uiAdmin2.pages.emptyTitle')"
        :description="filterCategory ? $t('uiAdmin2.pages.emptyFilteredDescription') : $t('uiAdmin2.pages.emptyDescription')"
      >
        <button type="button" class="btn btn-outline btn-sm gap-1.5" @click="openCreateModal">
          <Plus :size="14" :stroke-width="1.5" />
          {{ $t('uiAdmin2.pages.create') }}
        </button>
      </EmptyState>

      <div v-else class="overflow-x-auto">
        <table class="y-table">
          <thead>
            <tr>
              <th>{{ $t('admin.pageManagement.table.id') }}</th>
              <th>{{ $t('admin.pageManagement.table.title') }}</th>
              <th>{{ $t('admin.pageManagement.table.slug') }}</th>
              <th>{{ $t('admin.pageManagement.table.category') }}</th>
              <th>{{ $t('admin.pageManagement.table.status') }}</th>
              <th class="text-right">{{ $t('admin.pageManagement.table.displayOrder') }}</th>
              <th>{{ $t('admin.pageManagement.table.updatedAt') }}</th>
              <th class="text-right">{{ $t('admin.pageManagement.table.actions') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="page in filteredPages" :key="page.id">
              <td class="font-mono text-xs text-y-faint">{{ page.id }}</td>
              <td class="font-medium text-y-text">{{ page.title }}</td>
              <td class="whitespace-nowrap font-mono text-sm text-y-muted">/pages/{{ page.slug }}</td>
              <td class="whitespace-nowrap">
                <span class="y-chip">{{ $t(`admin.pageManagement.category.${page.category}`) }}</span>
              </td>
              <td class="whitespace-nowrap">
                <StatusDot
                  :tone="page.is_published ? 'ok' : 'muted'"
                  :label="page.is_published ? $t('admin.pageManagement.status.published') : $t('admin.pageManagement.status.draft')"
                />
              </td>
              <td class="text-right font-mono">{{ page.display_order }}</td>
              <td class="whitespace-nowrap text-y-muted">{{ formatDate(page.updated_at) }}</td>
              <td>
                <div class="flex justify-end gap-1">
                  <a
                    v-if="page.is_published"
                    :href="`/pages/${page.slug}`"
                    target="_blank"
                    rel="noopener"
                    class="btn btn-ghost btn-sm btn-square"
                    :title="`/pages/${page.slug}`"
                    :aria-label="`/pages/${page.slug}`"
                  >
                    <ExternalLink :size="16" :stroke-width="1.5" />
                  </a>
                  <button
                    type="button"
                    class="btn btn-ghost btn-sm btn-square"
                    :title="$t('uiAdmin2.edit')"
                    :aria-label="$t('uiAdmin2.edit')"
                    @click="openEditModal(page)"
                  >
                    <Pencil :size="16" :stroke-width="1.5" />
                  </button>
                  <button
                    type="button"
                    class="btn btn-ghost btn-sm btn-square text-y-danger-text"
                    :title="$t('uiAdmin2.delete')"
                    :aria-label="$t('uiAdmin2.delete')"
                    @click="confirmDelete(page)"
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
    <div v-if="showModal" class="modal modal-open">
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
          {{ isEditing ? $t('uiAdmin2.pages.editTitle') : $t('uiAdmin2.pages.createTitle') }}
        </h3>

        <form class="mt-6 space-y-4" @submit.prevent="handleSubmit">
          <label class="y-field">
            <span class="y-field-label">{{ $t('uiAdmin2.pages.fields.title') }}</span>
            <input v-model="formData.title" type="text" class="input input-bordered w-full" required />
          </label>

          <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <label class="y-field">
              <span class="y-field-label">{{ $t('uiAdmin2.pages.fields.slug') }}</span>
              <input
                v-model="formData.slug"
                type="text"
                class="input input-bordered w-full font-mono"
                :placeholder="$t('admin.pageManagement.modal.fields.slugPlaceholder')"
                autocomplete="off"
                autocapitalize="off"
                spellcheck="false"
                required
              />
              <span class="y-field-hint font-mono">{{ $t('uiAdmin2.pages.fields.slugHint', { slug: formData.slug || 'slug' }) }}</span>
            </label>

            <label class="y-field">
              <span class="y-field-label">{{ $t('uiAdmin2.pages.fields.category') }}</span>
              <select v-model="formData.category" class="select select-bordered w-full" required>
                <option value="company">{{ $t('admin.pageManagement.category.company') }}</option>
                <option value="resources">{{ $t('admin.pageManagement.category.resources') }}</option>
              </select>
            </label>
          </div>

          <label class="y-field">
            <span class="y-field-label">{{ $t('uiAdmin2.pages.fields.content') }}</span>
            <textarea
              v-model="formData.content"
              class="textarea textarea-bordered h-72 w-full font-mono text-sm leading-6"
              :placeholder="$t('admin.pageManagement.modal.fields.contentPlaceholder')"
              spellcheck="false"
              required
            ></textarea>
            <span class="y-field-hint">{{ $t('uiAdmin2.pages.fields.contentHint') }}</span>
          </label>

          <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <label class="y-field">
              <span class="y-field-label">{{ $t('uiAdmin2.pages.fields.displayOrder') }}</span>
              <input
                v-model.number="formData.display_order"
                type="number"
                class="input input-bordered w-full font-mono"
                min="0"
              />
              <span class="y-field-hint">{{ $t('uiAdmin2.pages.fields.displayOrderHint') }}</span>
            </label>

            <label class="flex cursor-pointer items-start gap-3 self-start rounded-md border border-y-line px-4 py-3 sm:mt-7">
              <input v-model="formData.is_published" type="checkbox" class="checkbox checkbox-sm mt-0.5" />
              <span class="flex flex-col gap-0.5">
                <span class="y-field-label">{{ $t('uiAdmin2.pages.fields.isPublished') }}</span>
                <span class="y-field-hint">{{ $t('uiAdmin2.pages.fields.isPublishedHint') }}</span>
              </span>
            </label>
          </div>

          <div class="modal-action">
            <button type="button" class="btn btn-ghost" @click="closeModal">{{ $t('uiAdmin2.cancel') }}</button>
            <button type="submit" class="btn btn-primary">
              {{ isEditing ? $t('uiAdmin2.saveChanges') : $t('uiAdmin2.pages.create') }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <div v-if="showDeleteModal" class="modal modal-open">
      <div class="modal-box">
        <h3 class="text-lg font-semibold text-y-text">{{ $t('uiAdmin2.pages.deleteTitle') }}</h3>
        <p class="mt-2 text-sm leading-6 text-y-muted">
          {{ $t('admin.pageManagement.delete.confirmMessage', { title: pageToDelete?.title }) }}
        </p>
        <div class="modal-action">
          <button type="button" class="btn btn-ghost" @click="showDeleteModal = false">{{ $t('uiAdmin2.cancel') }}</button>
          <button type="button" class="btn btn-error" @click="handleDelete">{{ $t('uiAdmin2.delete') }}</button>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import axios from '../utils/axios'
import { useToast } from '../composables/useToast'
import { Plus, Pencil, Trash2, X, ExternalLink } from 'lucide-vue-next'
import PageHeader from '../components/ui/PageHeader.vue'
import StatusDot from '../components/ui/StatusDot.vue'
import EmptyState from '../components/ui/EmptyState.vue'
import YLoader from '../components/brand/YLoader.vue'

const { t } = useI18n()
const toast = useToast()

const pages = ref([])
const loading = ref(true)
const filterCategory = ref('')
const showModal = ref(false)
const showDeleteModal = ref(false)
const isEditing = ref(false)
const pageToDelete = ref(null)

const formData = ref({
  title: '',
  slug: '',
  content: '',
  category: 'company',
  is_published: true,
  display_order: 0
})

const categoryFilters = computed(() => [
  { value: '', label: t('admin.pageManagement.filters.all') },
  { value: 'company', label: t('admin.pageManagement.filters.company') },
  { value: 'resources', label: t('admin.pageManagement.filters.resources') }
])

const filteredPages = computed(() => {
  if (!filterCategory.value) return pages.value
  return pages.value.filter(page => page.category === filterCategory.value)
})

const fetchPages = async () => {
  try {
    const response = await axios.get('/api/admin/pages')
    pages.value = response.data.pages || []
  } catch (error) {
    toast.error(t('admin.pageManagement.fetchError'))
  } finally {
    loading.value = false
  }
}

const openCreateModal = () => {
  isEditing.value = false
  formData.value = {
    title: '',
    slug: '',
    content: '',
    category: 'company',
    is_published: true,
    display_order: 0
  }
  showModal.value = true
}

const openEditModal = (page) => {
  isEditing.value = true
  formData.value = {
    id: page.id,
    title: page.title,
    slug: page.slug,
    content: page.content,
    category: page.category,
    is_published: page.is_published,
    display_order: page.display_order
  }
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
}

const handleSubmit = async () => {
  try {
    if (isEditing.value) {
      await axios.put(`/api/admin/pages/${formData.value.id}`, formData.value)
      toast.success(t('admin.pageManagement.updateSuccess'))
    } else {
      await axios.post('/api/admin/pages', formData.value)
      toast.success(t('admin.pageManagement.createSuccess'))
    }
    closeModal()
    fetchPages()
  } catch (error) {
    toast.error(error.response?.data?.error || t('admin.pageManagement.operationFailed'))
  }
}

const confirmDelete = (page) => {
  pageToDelete.value = page
  showDeleteModal.value = true
}

const handleDelete = async () => {
  try {
    await axios.delete(`/api/admin/pages/${pageToDelete.value.id}`)
    toast.success(t('admin.pageManagement.deleteSuccess'))
    showDeleteModal.value = false
    pageToDelete.value = null
    fetchPages()
  } catch (error) {
    toast.error(t('admin.pageManagement.deleteFailed'))
  }
}

const formatDate = (date) => {
  return new Date(date).toLocaleString('zh-CN')
}

onMounted(() => {
  fetchPages()
})
</script>
