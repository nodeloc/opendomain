<template>
  <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-8">
    <h1 class="text-4xl font-bold">{{ $t('user.profileSettings') }}</h1>

    <div class="card bg-base-200 shadow-xl">
      <div class="card-body">
        <h2 class="card-title">{{ $t('user.personalInfo') }}</h2>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="form-control">
            <label class="label">
              <span class="label-text">{{ $t('user.username') }}</span>
            </label>
            <input type="text" :value="user?.username" class="input input-bordered" disabled />
          </div>

          <div class="form-control">
            <label class="label">
              <span class="label-text">{{ $t('user.email') }}</span>
            </label>
            <input type="email" :value="user?.email" class="input input-bordered" disabled />
          </div>

          <div class="form-control">
            <label class="label">
              <span class="label-text">{{ $t('user.accountLevel') }}</span>
            </label>
            <input type="text" :value="user?.user_level" class="input input-bordered" disabled />
          </div>

          <div class="form-control">
            <label class="label">
              <span class="label-text">{{ $t('user.domainQuota') }}</span>
            </label>
            <input type="text" :value="user?.domain_quota" class="input input-bordered" disabled />
          </div>
        </div>

        <div class="divider"></div>

        <h3 class="text-xl font-semibold">{{ $t('user.inviteCode') }}</h3>
        <div class="flex gap-2">
          <input
            type="text"
            :value="user?.invite_code"
            class="input input-bordered flex-1"
            readonly
          />
          <button class="btn btn-primary" @click="copyInviteCode">{{ $t('user.copyInviteCode') }}</button>
        </div>
        <p class="text-sm opacity-70">
          {{ $t('user.inviteCodeDesc') }}
        </p>
      </div>
    </div>

    <!-- Change Password Section -->
    <div class="card bg-base-200 shadow-xl">
      <div class="card-body">
        <h2 class="card-title">{{ $t('user.changePassword') }}</h2>
        <p class="text-sm opacity-70 mb-4">{{ $t('user.changePasswordDesc') }}</p>

        <form @submit.prevent="changePassword" class="space-y-4">
          <div class="form-control">
            <label class="label">
              <span class="label-text">{{ $t('user.currentPassword') }}</span>
            </label>
            <input
              v-model="passwordForm.currentPassword"
              type="password"
              class="input input-bordered"
              :placeholder="$t('user.enterCurrentPassword')"
              required
            />
          </div>

          <div class="form-control">
            <label class="label">
              <span class="label-text">{{ $t('user.newPassword') }}</span>
            </label>
            <input
              v-model="passwordForm.newPassword"
              type="password"
              class="input input-bordered"
              :placeholder="$t('user.enterNewPassword')"
              minlength="6"
              required
            />
          </div>

          <div class="form-control">
            <label class="label">
              <span class="label-text">{{ $t('user.confirmNewPassword') }}</span>
            </label>
            <input
              v-model="passwordForm.confirmPassword"
              type="password"
              class="input input-bordered"
              :placeholder="$t('user.enterConfirmPassword')"
              minlength="6"
              required
            />
          </div>

          <div class="flex gap-2">
            <button type="submit" class="btn btn-primary" :disabled="changingPassword">
              <span v-if="changingPassword" class="loading loading-spinner"></span>
              {{ changingPassword ? $t('user.changing') : $t('user.changePasswordBtn') }}
            </button>
            <button type="button" class="btn btn-ghost" @click="resetPasswordForm">
              {{ $t('common.cancel') }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '../stores/auth'
import { useToast } from '../composables/useToast'
import axios from '../utils/axios'

const { t } = useI18n()
const authStore = useAuthStore()
const user = computed(() => authStore.user)
const toast = useToast()

const changingPassword = ref(false)
const passwordForm = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

onMounted(async () => {
  if (!authStore.user) {
    await authStore.fetchProfile()
  }
})

const copyInviteCode = () => {
  if (user.value?.invite_code) {
    navigator.clipboard.writeText(user.value.invite_code)
    toast.success(t('user.inviteCodeCopied'))
  }
}

const changePassword = async () => {
  if (passwordForm.newPassword !== passwordForm.confirmPassword) {
    toast.error(t('user.passwordMismatch'))
    return
  }

  if (passwordForm.newPassword.length < 6) {
    toast.error(t('user.passwordTooShort'))
    return
  }

  changingPassword.value = true
  try {
    await axios.put('/api/user/change-password', {
      current_password: passwordForm.currentPassword,
      new_password: passwordForm.newPassword
    })
    toast.success(t('user.passwordChangeSuccess'))
    resetPasswordForm()
  } catch (error) {
    toast.error(t('user.passwordChangeFailed') + ': ' + (error.response?.data?.error || error.message))
  } finally {
    changingPassword.value = false
  }
}

const resetPasswordForm = () => {
  passwordForm.currentPassword = ''
  passwordForm.newPassword = ''
  passwordForm.confirmPassword = ''
}
</script>
