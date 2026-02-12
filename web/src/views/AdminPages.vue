<template>
  <div class="container mx-auto px-4 py-8">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-3xl font-bold">{{ $t('admin.pageManagement.title') }}</h1>
      <button @click="openCreateModal" class="btn btn-primary">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M10 3a1 1 0 011 1v5h5a1 1 0 110 2h-5v5a1 1 0 11-2 0v-5H4a1 1 0 110-2h5V4a1 1 0 011-1z" clip-rule="evenodd" />
        </svg>
        {{ $t('admin.pageManagement.createPage') }}
      </button>
    </div>

    <!-- Filter -->
    <div class="mb-4 flex gap-2">
      <button 
        @click="filterCategory = ''"
        :class="['px-4 py-2 rounded', filterCategory === '' ? 'bg-primary text-white' : 'bg-base-200']"
      >
        {{ $t('admin.pageManagement.filters.all') }}
      </button>
      <button 
        @click="filterCategory = 'company'"
        :class="['px-4 py-2 rounded', filterCategory === 'company' ? 'bg-primary text-white' : 'bg-base-200']"
      >
        {{ $t('admin.pageManagement.filters.company') }}
      </button>
      <button 
        @click="filterCategory = 'resources'"
        :class="['px-4 py-2 rounded', filterCategory === 'resources' ? 'bg-primary text-white' : 'bg-base-200']"
      >
        {{ $t('admin.pageManagement.filters.resources') }}
      </button>
    </div>

    <!-- Pages List -->
    <div class="overflow-x-auto bg-base-100 rounded-lg shadow">
      <table class="table w-full">
        <thead>
          <tr>
            <th>{{ $t('admin.pageManagement.table.id') }}</th>
            <th>{{ $t('admin.pageManagement.table.title') }}</th>
            <th>{{ $t('admin.pageManagement.table.slug') }}</th>
            <th>{{ $t('admin.pageManagement.table.category') }}</th>
            <th>{{ $t('admin.pageManagement.table.status') }}</th>
            <th>{{ $t('admin.pageManagement.table.displayOrder') }}</th>
            <th>{{ $t('admin.pageManagement.table.updatedAt') }}</th>
            <th>{{ $t('admin.pageManagement.table.actions') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="page in filteredPages" :key="page.id">
            <td>{{ page.id }}</td>
            <td>{{ page.title }}</td>
            <td><code class="text-sm">{{ page.slug }}</code></td>
            <td>
              <span :class="['px-2 py-0.5 rounded text-xs font-medium', 
                page.category === 'company' ? 'bg-blue-100 text-blue-800' : 'bg-green-100 text-green-800']">
                {{ $t(`admin.pageManagement.category.${page.category}`) }}
              </span>
            </td>
            <td>
              <span :class="['px-2 py-0.5 rounded text-xs font-medium',
                page.is_published ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800']">
                {{ page.is_published ? $t('admin.pageManagement.status.published') : $t('admin.pageManagement.status.draft') }}
              </span>
            </td>
            <td>{{ page.display_order }}</td>
            <td>{{ formatDate(page.updated_at) }}</td>
            <td>
              <div class="flex gap-2">
                <button @click="openEditModal(page)" class="btn btn-sm btn-ghost">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                    <path d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z" />
                  </svg>
                </button>
                <button @click="confirmDelete(page)" class="btn btn-sm btn-ghost text-error">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clip-rule="evenodd" />
                  </svg>
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
      <div v-if="filteredPages.length === 0" class="text-center py-8 text-gray-500">
        {{ $t('admin.pageManagement.noData') }}
      </div>
    </div>

    <!-- Create/Edit Modal -->
    <div v-if="showModal" class="modal modal-open">
      <div class="modal-box max-w-2xl">
        <h3 class="font-bold text-lg mb-4">{{ isEditing ? $t('admin.pageManagement.modal.editTitle') : $t('admin.pageManagement.modal.createTitle') }}</h3>
        <form @submit.prevent="handleSubmit">
          <div class="form-control mb-4">
            <label class="label">
              <span class="label-text">{{ $t('admin.pageManagement.modal.fields.title') }}</span>
            </label>
            <input 
              v-model="formData.title" 
              type="text" 
              class="input input-bordered" 
              required 
            />
          </div>

          <div class="form-control mb-4">
            <label class="label">
              <span class="label-text">{{ $t('admin.pageManagement.modal.fields.slug') }}</span>
            </label>
            <input 
              v-model="formData.slug" 
              type="text" 
              class="input input-bordered" 
              :placeholder="$t('admin.pageManagement.modal.fields.slugPlaceholder')" 
              required 
            />
            <label class="label">
              <span class="label-text-alt">{{ $t('admin.pageManagement.modal.fields.slugHint') }}</span>
            </label>
          </div>

          <div class="form-control mb-4">
            <label class="label">
              <span class="label-text">{{ $t('admin.pageManagement.modal.fields.category') }}</span>
            </label>
            <select v-model="formData.category" class="select select-bordered" required>
              <option value="company">{{ $t('admin.pageManagement.category.company') }}</option>
              <option value="resources">{{ $t('admin.pageManagement.category.resources') }}</option>
            </select>
          </div>

          <div class="form-control mb-4">
            <label class="label">
              <span class="label-text">{{ $t('admin.pageManagement.modal.fields.content') }}</span>
            </label>
            <textarea 
              v-model="formData.content" 
              class="textarea textarea-bordered h-64" 
              :placeholder="$t('admin.pageManagement.modal.fields.contentPlaceholder')"
              required
            ></textarea>
          </div>

          <div class="form-control mb-4">
            <label class="label">
              <span class="label-text">{{ $t('admin.pageManagement.modal.fields.displayOrder') }}</span>
            </label>
            <input 
              v-model.number="formData.display_order" 
              type="number" 
              class="input input-bordered" 
              min="0"
            />
            <label class="label">
              <span class="label-text-alt">{{ $t('admin.pageManagement.modal.fields.displayOrderHint') }}</span>
            </label>
          </div>

          <div class="form-control mb-4">
            <label class="label cursor-pointer justify-start gap-2">
              <input 
                v-model="formData.is_published" 
                type="checkbox" 
                class="checkbox checkbox-primary" 
              />
              <span class="label-text">{{ $t('admin.pageManagement.modal.fields.isPublished') }}</span>
            </label>
          </div>

          <div class="modal-action">
            <button type="button" @click="closeModal" class="btn btn-ghost">{{ $t('admin.pageManagement.modal.actions.cancel') }}</button>
            <button type="submit" class="btn btn-primary">{{ isEditing ? $t('admin.pageManagement.modal.actions.save') : $t('admin.pageManagement.modal.actions.create') }}</button>
          </div>
        </form>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <div v-if="showDeleteModal" class="modal modal-open">
      <div class="modal-box">
        <h3 class="font-bold text-lg">{{ $t('admin.pageManagement.delete.confirmTitle') }}</h3>
        <p class="py-4">{{ $t('admin.pageManagement.delete.confirmMessage', { title: pageToDelete?.title }) }}</p>
        <div class="modal-action">
          <button @click="showDeleteModal = false" class="btn btn-ghost">{{ $t('admin.pageManagement.delete.cancel') }}</button>
          <button @click="handleDelete" class="btn btn-error">{{ $t('admin.pageManagement.delete.confirm') }}</button>
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

const { t } = useI18n()
const toast = useToast()

const pages = ref([])
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
