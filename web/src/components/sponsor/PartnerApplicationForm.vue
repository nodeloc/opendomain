<template>
  <div class="y-card p-5 sm:p-8">
    <!-- Success state -->
    <div v-if="sent" ref="successRef" role="status" tabindex="-1" class="flex scroll-mt-28 flex-col gap-3.5 outline-none">
      <span class="flex h-12 w-12 items-center justify-center rounded-full bg-y-mint" style="box-shadow: 0 0 24px rgb(var(--y-mint) / 0.5);" aria-hidden="true">
        <Check :size="24" :stroke-width="2.5" class="text-y-on-mint" />
      </span>
      <h3 class="m-0 text-[28px] font-extrabold leading-[34px] tracking-[-0.03em] text-y-text">{{ $t('partner.successTitle') }}</h3>
      <i18n-t keypath="partner.successBody" tag="p" class="m-0 text-y-muted">
        <template #name><span class="font-medium text-y-text">{{ submitted.name }}</span></template>
        <template #email><span class="break-all font-mono text-y-text">{{ submitted.email }}</span></template>
      </i18n-t>
      <p v-if="submitted.id" class="m-0 font-mono text-xs uppercase tracking-[0.08em] text-y-faint">
        {{ $t('partner.successRef', { id: submitted.id }) }}
      </p>
      <div class="pt-1">
        <button type="button" class="btn btn-outline" @click="resetForm">{{ $t('partner.another') }}</button>
      </div>
    </div>

    <!-- Form -->
    <form v-else class="flex flex-col gap-6" novalidate @submit.prevent="submit">
      <div class="flex flex-wrap items-baseline justify-between gap-x-4 gap-y-1">
        <h3 class="m-0 text-xl font-bold tracking-[-0.02em] text-y-text">{{ $t('partner.formTitle') }}</h3>
        <span class="text-xs text-y-faint"><span class="text-y-danger-text" aria-hidden="true">*</span> {{ $t('partner.required') }}</span>
      </div>

      <div class="grid gap-6 sm:grid-cols-2">
        <div class="y-field">
          <label for="partner-name" class="y-field-label">
            {{ $t('partner.nameLabel') }} <span class="text-y-danger-text" aria-hidden="true">*</span>
          </label>
          <input
            id="partner-name"
            ref="nameRef"
            v-model="form.name"
            type="text"
            name="name"
            :maxlength="LIMITS.name"
            :placeholder="$t('partner.namePlaceholder')"
            autocomplete="name"
            required
            class="input input-bordered w-full"
            :class="{ 'input-error': errors.name }"
            :aria-invalid="errors.name ? 'true' : 'false'"
            :aria-describedby="errors.name ? 'partner-name-error' : undefined"
          />
          <p v-if="errors.name" id="partner-name-error" class="y-field-error">{{ errors.name }}</p>
        </div>

        <div class="y-field">
          <label for="partner-email" class="y-field-label">
            {{ $t('partner.emailLabel') }} <span class="text-y-danger-text" aria-hidden="true">*</span>
          </label>
          <input
            id="partner-email"
            v-model="form.email"
            type="email"
            name="email"
            inputmode="email"
            :maxlength="LIMITS.email"
            placeholder="you@example.com"
            autocomplete="email"
            autocapitalize="off"
            spellcheck="false"
            required
            class="input input-bordered w-full"
            :class="{ 'input-error': errors.email }"
            :aria-invalid="errors.email ? 'true' : 'false'"
            :aria-describedby="errors.email ? 'partner-email-error partner-email-hint' : 'partner-email-hint'"
          />
          <p v-if="errors.email" id="partner-email-error" class="y-field-error">{{ errors.email }}</p>
          <p id="partner-email-hint" class="y-field-hint">{{ $t('partner.emailHint') }}</p>
        </div>

        <div class="y-field">
          <label for="partner-org" class="y-field-label">
            {{ $t('partner.organizationLabel') }} <span class="font-normal text-y-faint">{{ $t('partner.optional') }}</span>
          </label>
          <input
            id="partner-org"
            v-model="form.organization"
            type="text"
            name="organization"
            :maxlength="LIMITS.organization"
            :placeholder="$t('partner.organizationPlaceholder')"
            autocomplete="organization"
            class="input input-bordered w-full"
            :class="{ 'input-error': errors.organization }"
            :aria-invalid="errors.organization ? 'true' : 'false'"
            :aria-describedby="errors.organization ? 'partner-org-error' : undefined"
          />
          <p v-if="errors.organization" id="partner-org-error" class="y-field-error">{{ errors.organization }}</p>
        </div>

        <div class="y-field">
          <label for="partner-website" class="y-field-label">
            {{ $t('partner.websiteLabel') }} <span class="font-normal text-y-faint">{{ $t('partner.optional') }}</span>
          </label>
          <input
            id="partner-website"
            v-model="form.website"
            type="url"
            name="website"
            inputmode="url"
            :maxlength="LIMITS.website"
            placeholder="https://example.com"
            autocomplete="url"
            autocapitalize="off"
            spellcheck="false"
            class="input input-bordered w-full font-mono text-[15px]"
            :class="{ 'input-error': errors.website }"
            :aria-invalid="errors.website ? 'true' : 'false'"
            :aria-describedby="errors.website ? 'partner-website-error' : undefined"
          />
          <p v-if="errors.website" id="partner-website-error" class="y-field-error">{{ errors.website }}</p>
        </div>
      </div>

      <fieldset
        class="m-0 flex min-w-0 flex-col gap-2.5 border-0 p-0"
        :aria-invalid="errors.partnership_type ? 'true' : undefined"
        :aria-describedby="errors.partnership_type ? 'partner-type-error' : undefined"
      >
        <legend class="float-left w-full p-0 pb-2.5 text-sm font-medium text-y-text">
          {{ $t('partner.typeLegend') }} <span class="text-y-danger-text" aria-hidden="true">*</span>
        </legend>
        <div class="grid gap-2.5 sm:grid-cols-2">
          <label
            v-for="type in typeOptions"
            :key="type.id"
            class="y-choice partner-choice cursor-pointer"
            :class="{ 'is-active': form.partnership_type === type.id, 'partner-choice-error': errors.partnership_type }"
          >
            <input
              :ref="(el) => (typeRefs[type.id] = el)"
              v-model="form.partnership_type"
              type="radio"
              name="partnership_type"
              :value="type.id"
              class="sr-only"
              required
            />
            <span class="font-semibold text-y-text">{{ type.label }}</span>
            <span class="text-[13px] leading-[18px] text-y-muted">{{ type.hint }}</span>
          </label>
        </div>
        <p v-if="errors.partnership_type" id="partner-type-error" class="y-field-error">{{ errors.partnership_type }}</p>
      </fieldset>

      <div class="y-field">
        <label for="partner-message" class="y-field-label">
          {{ $t('partner.messageLabel') }} <span class="text-y-danger-text" aria-hidden="true">*</span>
        </label>
        <textarea
          id="partner-message"
          v-model="form.message"
          name="message"
          rows="6"
          :maxlength="LIMITS.message"
          :placeholder="$t('partner.messagePlaceholder')"
          required
          class="textarea textarea-bordered w-full text-[15px] leading-6"
          :class="{ 'textarea-error': errors.message }"
          :aria-invalid="errors.message ? 'true' : 'false'"
          :aria-describedby="errors.message ? 'partner-message-error partner-message-count' : 'partner-message-count'"
        ></textarea>
        <div class="flex items-start justify-between gap-4">
          <p v-if="errors.message" id="partner-message-error" class="y-field-error">{{ errors.message }}</p>
          <p id="partner-message-count" class="y-field-hint ml-auto shrink-0 font-mono">
            {{ $t('partner.messageCount', { n: form.message.length, max: LIMITS.message }) }}
          </p>
        </div>
      </div>

      <div class="y-field sm:max-w-[calc(50%-0.75rem)]">
        <label for="partner-contact" class="y-field-label">
          {{ $t('partner.contactLabel') }} <span class="font-normal text-y-faint">{{ $t('partner.optional') }}</span>
        </label>
        <input
          id="partner-contact"
          v-model="form.contact"
          type="text"
          name="contact"
          :maxlength="LIMITS.contact"
          :placeholder="$t('partner.contactPlaceholder')"
          autocomplete="off"
          class="input input-bordered w-full"
          :class="{ 'input-error': errors.contact }"
          :aria-invalid="errors.contact ? 'true' : 'false'"
          :aria-describedby="errors.contact ? 'partner-contact-error' : undefined"
        />
        <p v-if="errors.contact" id="partner-contact-error" class="y-field-error">{{ errors.contact }}</p>
      </div>

      <!-- Honeypot: bots fill every field, people never see this one. -->
      <div class="partner-hp" aria-hidden="true">
        <label for="partner-fax">{{ $t('partner.honeypotLabel') }}</label>
        <input id="partner-fax" v-model="form.fax" type="text" name="fax" tabindex="-1" autocomplete="off" />
      </div>

      <div v-if="formError" role="alert" class="flex items-start gap-2.5 text-sm text-y-danger-text">
        <CircleAlert :size="16" :stroke-width="1.5" class="mt-0.5 shrink-0" aria-hidden="true" />
        <span>{{ formError }}</span>
      </div>

      <div class="flex flex-wrap items-center gap-x-4 gap-y-3">
        <button type="submit" class="y-btn-ink min-h-[48px] w-full rounded-md px-6 sm:w-auto" :disabled="sending" :aria-busy="sending ? 'true' : 'false'">
          <span v-if="sending" class="loading loading-spinner loading-sm" aria-hidden="true"></span>
          <Send v-else :size="16" :stroke-width="1.75" aria-hidden="true" />
          {{ sending ? $t('partner.sending') : $t('partner.submit') }}
        </button>
        <span class="max-w-[320px] text-[13px] leading-5 text-y-faint">{{ $t('partner.privacy') }}</span>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, reactive, computed, watch, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { Check, CircleAlert, Send } from 'lucide-vue-next'
import axios from '../../utils/axios'

const { t } = useI18n()

const TYPES = ['sponsorship', 'resources', 'business', 'community', 'other']
const LIMITS = { name: 100, email: 254, organization: 200, website: 500, message: 5000, contact: 200 }
const MESSAGE_MIN = 10
const FIELD_ORDER = ['name', 'email', 'organization', 'website', 'partnership_type', 'message', 'contact']
const FIELD_IDS = {
  name: 'partner-name',
  email: 'partner-email',
  organization: 'partner-org',
  website: 'partner-website',
  message: 'partner-message',
  contact: 'partner-contact',
}

const emptyForm = () => ({
  name: '',
  email: '',
  organization: '',
  website: '',
  partnership_type: '',
  message: '',
  contact: '',
  fax: '',
})

const form = reactive(emptyForm())
const errors = reactive({})
const formError = ref('')
const sending = ref(false)
const sent = ref(false)
const attempted = ref(false)
const submitted = reactive({ id: null, name: '', email: '' })
const nameRef = ref(null)
const successRef = ref(null)
const typeRefs = {}

const typeOptions = computed(() => TYPES.map((id) => ({
  id,
  label: t(`partner.types.${id}`),
  hint: t(`partner.typeHints.${id}`),
})))

const isHttpUrl = (value) => {
  try {
    const u = new URL(value)
    return (u.protocol === 'http:' || u.protocol === 'https:') && !!u.hostname
  } catch {
    return false
  }
}

const validateAll = () => {
  const e = {}
  const name = form.name.trim()
  const email = form.email.trim()
  const organization = form.organization.trim()
  const website = form.website.trim()
  const message = form.message.trim()
  const contact = form.contact.trim()

  if (!name) e.name = t('partner.errNameRequired')
  else if ([...name].length > LIMITS.name) e.name = t('partner.errNameLong', { max: LIMITS.name })

  if (!email) e.email = t('partner.errEmailRequired')
  else if (email.length > LIMITS.email) e.email = t('partner.errEmailLong')
  else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) e.email = t('partner.errEmailInvalid')

  if ([...organization].length > LIMITS.organization) e.organization = t('partner.errOrganizationLong', { max: LIMITS.organization })

  if (website) {
    if (website.length > LIMITS.website) e.website = t('partner.errWebsiteLong', { max: LIMITS.website })
    else if (!/^https?:\/\//i.test(website) || !isHttpUrl(website)) e.website = t('partner.errWebsiteInvalid')
  }

  if (!TYPES.includes(form.partnership_type)) e.partnership_type = t('partner.errTypeRequired')

  const messageLen = [...message].length
  if (messageLen < MESSAGE_MIN) e.message = t('partner.errMessageShort', { min: MESSAGE_MIN })
  else if (messageLen > LIMITS.message) e.message = t('partner.errMessageLong', { max: LIMITS.message })

  if ([...contact].length > LIMITS.contact) e.contact = t('partner.errContactLong', { max: LIMITS.contact })

  return e
}

const applyErrors = (next) => {
  for (const key of Object.keys(errors)) delete errors[key]
  Object.assign(errors, next)
}

// Re-validate as the user types, but only after the first submit so empty fields are not flagged on load.
watch(form, () => {
  if (!attempted.value) return
  applyErrors(validateAll())
  if (formError.value && Object.keys(errors).length === 0) formError.value = ''
})

const focusField = (key) => {
  if (key === 'partnership_type') {
    const el = typeRefs[form.partnership_type] || typeRefs[TYPES[0]]
    el?.focus()
    return
  }
  document.getElementById(FIELD_IDS[key])?.focus()
}

const submit = async () => {
  if (sending.value) return
  attempted.value = true
  formError.value = ''
  const next = validateAll()
  applyErrors(next)
  const firstInvalid = FIELD_ORDER.find((key) => next[key])
  if (firstInvalid) {
    formError.value = t('partner.errSummary')
    focusField(firstInvalid)
    return
  }

  sending.value = true
  try {
    const body = {
      name: form.name.trim(),
      email: form.email.trim(),
      organization: form.organization.trim(),
      website: form.website.trim(),
      partnership_type: form.partnership_type,
      message: form.message.trim(),
      contact: form.contact.trim(),
      fax: form.fax,
    }
    const { data } = await axios.post('/api/public/partner-applications', body)
    submitted.id = data?.id ?? null
    submitted.name = body.name
    submitted.email = body.email
    sent.value = true
    await nextTick()
    // The confirmation is much shorter than the form, so bring it back into view.
    successRef.value?.scrollIntoView({ behavior: 'smooth', block: 'nearest' })
    successRef.value?.focus({ preventScroll: true })
  } catch (err) {
    const status = err.response?.status
    const serverError = err.response?.data?.error
    if (status === 429) formError.value = t('partner.errRateLimited')
    else if (status === 400 && serverError) formError.value = serverError
    else formError.value = t('partner.errGeneric')
  } finally {
    sending.value = false
  }
}

const resetForm = async () => {
  Object.assign(form, emptyForm())
  applyErrors({})
  formError.value = ''
  attempted.value = false
  sent.value = false
  submitted.id = null
  submitted.name = ''
  submitted.email = ''
  await nextTick()
  nameRef.value?.focus()
}

const setType = (type) => {
  if (!TYPES.includes(type)) return
  if (sent.value) resetForm()
  form.partnership_type = type
}

defineExpose({ setType })
</script>

<style scoped>
.y-field-error {
  margin: 0;
  font-size: 0.8125rem;
  line-height: 1.25rem;
  color: rgb(var(--y-danger-text));
}
/* The global focus style paints inputs mint, which would hide a field that failed validation. */
.input-error,
.input-error:focus,
.textarea-error,
.textarea-error:focus {
  border-color: rgb(var(--y-danger) / 0.8);
  outline-color: rgb(var(--y-danger) / 0.4);
}
.partner-choice:focus-within {
  outline: 2px solid rgb(var(--y-mint) / 0.6);
  outline-offset: 2px;
}
.partner-choice-error:not(.is-active) {
  border-color: rgb(var(--y-danger) / 0.5);
}
.partner-hp {
  position: absolute;
  left: -10000px;
  top: auto;
  width: 1px;
  height: 1px;
  overflow: hidden;
}
</style>
