<template>
  <div class="y-page">
    <PageHeader
      :eyebrow="$t('shell.adminArea')"
      :title="$t('uiAdmin1.rootDomains.title')"
      :description="$t('uiAdmin1.rootDomains.description')"
    >
      <template #actions>
        <button type="button" class="btn btn-primary gap-2" @click="showCreateModal = true">
          <Plus :size="16" :stroke-width="1.5" />
          {{ $t('admin.addRootDomain') }}
        </button>
      </template>
    </PageHeader>

    <section class="y-panel overflow-hidden">
      <YLoader v-if="loading" :label="$t('common.loading')" />

      <EmptyState
        v-else-if="rootDomains.length === 0"
        :title="$t('admin.noRootDomains')"
        :description="$t('admin.noRootDomainsDesc')"
      >
        <button type="button" class="btn btn-outline gap-2" @click="showCreateModal = true">
          <Plus :size="16" :stroke-width="1.5" />
          {{ $t('admin.addRootDomain') }}
        </button>
      </EmptyState>

      <div v-else class="overflow-x-auto">
        <table class="y-table">
          <thead>
            <tr>
              <th>{{ $t('uiAdmin1.rootDomains.columns.domain') }}</th>
              <th>{{ $t('uiAdmin1.rootDomains.columns.status') }}</th>
              <th>{{ $t('uiAdmin1.rootDomains.columns.priority') }}</th>
              <th>{{ $t('uiAdmin1.rootDomains.columns.registrations') }}</th>
              <th>{{ $t('uiAdmin1.rootDomains.columns.dnssec') }}</th>
              <th>{{ $t('uiAdmin1.rootDomains.columns.pricing') }}</th>
              <th class="text-right">{{ $t('uiAdmin1.rootDomains.columns.actions') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="domain in rootDomains" :key="domain.id">
              <td class="min-w-[14rem]">
                <div class="flex flex-wrap items-center gap-2">
                  <span class="y-domain text-[15px]">.{{ domain.domain }}</span>
                  <span v-if="domain.is_hot" class="y-chip-warn">
                    <Flame :size="12" :stroke-width="1.5" />
                    {{ $t('admin.hot') }}
                  </span>
                  <span v-if="domain.is_new" class="y-chip-info">
                    <Sparkles :size="12" :stroke-width="1.5" />
                    {{ $t('admin.new') }}
                  </span>
                </div>
                <p v-if="domain.description" class="mt-1 max-w-xs text-xs leading-5 text-y-muted">{{ domain.description }}</p>
              </td>
              <td class="whitespace-nowrap">
                <StatusDot :tone="domain.is_active ? 'ok' : 'muted'" :label="domain.is_active ? $t('admin.active') : $t('admin.inactive')" />
              </td>
              <td class="font-mono text-sm text-y-muted">{{ domain.priority }}</td>
              <td>
                <router-link
                  :to="`/admin/root-domains/${domain.id}/domains`"
                  class="inline-flex items-center gap-1.5 rounded-sm font-mono text-sm text-y-text transition-colors hover:text-y-mint-text"
                  :title="$t('uiAdmin1.rootDomains.viewRegistrations')"
                >
                  <Users :size="14" :stroke-width="1.5" class="text-y-faint" />
                  {{ domain.registration_count || 0 }}
                  <ArrowUpRight :size="14" :stroke-width="1.5" class="text-y-faint" />
                </router-link>
              </td>
              <td class="whitespace-nowrap">
                <button
                  v-if="domain.use_default_nameservers !== false"
                  type="button"
                  class="cursor-pointer underline-offset-2 hover:underline"
                  :class="dnssecBadgeClass(domain.id)"
                  :title="$t('uiAdmin1.rootDomains.manageDnssec')"
                  @click="openDNSSEC(domain)"
                >
                  <ShieldCheck :size="12" :stroke-width="1.5" />
                  {{ dnssecBadgeText(domain.id) }}
                </button>
                <span v-else class="text-xs text-y-faint">{{ $t('uiAdmin1.rootDomains.externalNs') }}</span>
              </td>
              <td class="whitespace-nowrap">
                <span v-if="domain.is_free" class="y-chip">{{ $t('admin.free') }}</span>
                <div v-else class="space-y-0.5">
                  <p v-if="domain.price_per_year" class="font-mono text-sm text-y-text">
                    {{ $t('uiAdmin1.rootDomains.perYear', { price: formatPrice(domain.price_per_year) }) }}
                  </p>
                  <p v-if="domain.lifetime_price" class="font-mono text-xs text-y-muted">
                    {{ $t('uiAdmin1.rootDomains.lifetime', { price: formatPrice(domain.lifetime_price) }) }}
                  </p>
                </div>
              </td>
              <td>
                <div class="flex justify-end gap-1">
                  <button
                    type="button"
                    class="btn btn-ghost btn-sm btn-square"
                    :title="$t('uiAdmin1.rootDomains.edit')"
                    :aria-label="$t('uiAdmin1.rootDomains.edit')"
                    @click="editDomain(domain)"
                  >
                    <Pencil :size="16" :stroke-width="1.5" />
                  </button>
                  <button
                    type="button"
                    class="btn btn-ghost btn-sm btn-square text-y-danger-text"
                    :title="$t('uiAdmin1.rootDomains.delete')"
                    :aria-label="$t('uiAdmin1.rootDomains.delete')"
                    @click="deleteDomain(domain)"
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
    <dialog :class="{ 'modal': true, 'modal-open': showCreateModal || showEditModal }">
      <div class="modal-box w-11/12 max-w-2xl">
        <h3 class="text-lg font-semibold text-y-text">
          {{ showEditModal ? $t('uiAdmin1.rootDomains.form.editTitle') : $t('uiAdmin1.rootDomains.form.createTitle') }}
        </h3>

        <form class="mt-5 space-y-6" @submit.prevent="showEditModal ? updateDomain() : createDomain()">
          <div class="space-y-4">
            <div class="y-field">
              <label class="y-field-label" for="root-domain-name">{{ $t('uiAdmin1.rootDomains.form.domain') }}</label>
              <div class="relative">
                <span class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 font-mono text-y-faint">.</span>
                <input
                  id="root-domain-name"
                  v-model="formData.domain"
                  type="text"
                  placeholder="example.com"
                  class="input input-bordered w-full pl-6 font-mono"
                  :disabled="showEditModal"
                  required
                />
              </div>
            </div>

            <div class="y-field">
              <label class="y-field-label" for="root-domain-description">
                {{ $t('uiAdmin1.rootDomains.form.description') }}
                <span class="font-normal text-y-faint">· {{ $t('uiAdmin1.rootDomains.form.optional') }}</span>
              </label>
              <textarea
                id="root-domain-description"
                v-model="formData.description"
                class="textarea textarea-bordered h-20 w-full"
                :placeholder="$t('uiAdmin1.rootDomains.form.descriptionPlaceholder')"
              ></textarea>
            </div>

            <div class="y-field sm:max-w-[12rem]">
              <label class="y-field-label" for="root-domain-priority">{{ $t('uiAdmin1.rootDomains.form.priority') }}</label>
              <input
                id="root-domain-priority"
                v-model.number="formData.priority"
                type="number"
                class="input input-bordered w-full font-mono"
                min="0"
                max="100"
              />
              <span class="y-field-hint">{{ $t('uiAdmin1.rootDomains.form.priorityHint') }}</span>
            </div>
          </div>

          <div class="y-divider space-y-4 pt-5">
            <p class="y-label">{{ $t('uiAdmin1.rootDomains.form.pricing') }}</p>
            <div class="y-field">
              <label class="flex cursor-pointer items-center gap-3">
                <input v-model="formData.is_free" type="checkbox" class="toggle toggle-primary toggle-sm" />
                <span class="text-sm font-medium text-y-text">
                  {{ formData.is_free ? $t('uiAdmin1.rootDomains.form.freeDomain') : $t('uiAdmin1.rootDomains.form.paidDomain') }}
                </span>
              </label>
              <span class="y-field-hint">{{ $t('uiAdmin1.rootDomains.form.pricingHint') }}</span>
            </div>

            <div v-if="!formData.is_free" class="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div class="y-field">
                <label class="y-field-label" for="root-domain-price-year">
                  {{ $t('uiAdmin1.rootDomains.form.pricePerYear', { symbol: currencySymbol }) }}
                </label>
                <input
                  id="root-domain-price-year"
                  v-model.number="formData.price_per_year"
                  type="number"
                  step="0.01"
                  min="0"
                  class="input input-bordered w-full font-mono"
                  placeholder="0.00"
                />
                <span class="y-field-hint">{{ $t('uiAdmin1.rootDomains.form.pricePerYearHint') }}</span>
              </div>

              <div class="y-field">
                <label class="y-field-label" for="root-domain-price-lifetime">
                  {{ $t('uiAdmin1.rootDomains.form.lifetimePrice', { symbol: currencySymbol }) }}
                </label>
                <input
                  id="root-domain-price-lifetime"
                  v-model.number="formData.lifetime_price"
                  type="number"
                  step="0.01"
                  min="0"
                  class="input input-bordered w-full font-mono"
                  placeholder="0.00"
                />
                <span class="y-field-hint">{{ $t('uiAdmin1.rootDomains.form.lifetimePriceHint') }}</span>
              </div>
            </div>
          </div>

          <div class="y-divider space-y-4 pt-5">
            <p class="y-label">{{ $t('uiAdmin1.rootDomains.form.nameservers') }}</p>
            <div class="y-field">
              <label class="flex cursor-pointer items-center gap-3">
                <input v-model="formData.use_default_nameservers" type="checkbox" class="toggle toggle-primary toggle-sm" />
                <span class="text-sm font-medium text-y-text">
                  {{ formData.use_default_nameservers ? $t('uiAdmin1.rootDomains.form.usingDefaultNs') : $t('uiAdmin1.rootDomains.form.usingCustomNs') }}
                </span>
              </label>
              <span class="y-field-hint">{{ $t('uiAdmin1.rootDomains.form.nsHint') }}</span>
            </div>

            <div v-if="formData.use_default_nameservers" class="y-inset flex items-start gap-3 px-4 py-3">
              <Info :size="16" :stroke-width="1.5" class="mt-0.5 shrink-0 text-y-info" />
              <div class="min-w-0">
                <p class="text-sm font-medium text-y-text">{{ $t('uiAdmin1.rootDomains.form.defaultNsTitle') }}</p>
                <p class="mt-0.5 break-all font-mono text-xs text-y-muted">ns1.nodelook.com, ns2.nodelook.com</p>
              </div>
            </div>

            <div v-else class="space-y-2">
              <div v-for="(ns, index) in formData.nameservers" :key="index" class="flex gap-2">
                <input
                  v-model="formData.nameservers[index]"
                  type="text"
                  class="input input-bordered flex-1 font-mono"
                  placeholder="ns1.example.com"
                />
                <button
                  v-if="formData.nameservers.length > 1"
                  type="button"
                  class="btn btn-ghost btn-square"
                  :aria-label="$t('uiAdmin1.rootDomains.form.removeNameserver')"
                  @click="formData.nameservers.splice(index, 1)"
                >
                  <X :size="16" :stroke-width="1.5" />
                </button>
              </div>
              <button type="button" class="btn btn-ghost btn-sm gap-2" @click="formData.nameservers.push('')">
                <Plus :size="16" :stroke-width="1.5" />
                {{ $t('uiAdmin1.rootDomains.form.addNameserver') }}
              </button>
            </div>
          </div>

          <div class="y-divider space-y-4 pt-5">
            <p class="y-label">{{ $t('uiAdmin1.rootDomains.form.display') }}</p>
            <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div class="y-field">
                <span class="y-field-label">{{ $t('uiAdmin1.rootDomains.form.status') }}</span>
                <label class="flex h-10 cursor-pointer items-center gap-3">
                  <input v-model="formData.is_active" type="checkbox" class="toggle toggle-primary toggle-sm" />
                  <span class="text-sm text-y-muted">
                    {{ formData.is_active ? $t('uiAdmin1.rootDomains.form.active') : $t('uiAdmin1.rootDomains.form.inactive') }}
                  </span>
                </label>
              </div>

              <div class="y-field">
                <label class="y-field-label" for="root-domain-min-level">{{ $t('uiAdmin1.rootDomains.form.minLevel') }}</label>
                <select id="root-domain-min-level" v-model="formData.min_user_level" class="select select-bordered w-full">
                  <option value="normal">{{ $t('uiAdmin1.rootDomains.form.levels.normal') }}</option>
                  <option value="basic">{{ $t('uiAdmin1.rootDomains.form.levels.basic') }}</option>
                  <option value="member">{{ $t('uiAdmin1.rootDomains.form.levels.member') }}</option>
                  <option value="regular">{{ $t('uiAdmin1.rootDomains.form.levels.regular') }}</option>
                  <option value="leader">{{ $t('uiAdmin1.rootDomains.form.levels.leader') }}</option>
                </select>
                <span class="y-field-hint">{{ $t('uiAdmin1.rootDomains.form.minLevelHint') }}</span>
              </div>

              <div class="y-field">
                <span class="y-field-label">{{ $t('uiAdmin1.rootDomains.form.markHot') }}</span>
                <label class="flex h-10 cursor-pointer items-center gap-3">
                  <input v-model="formData.is_hot" type="checkbox" class="toggle toggle-primary toggle-sm" />
                  <span class="text-sm text-y-muted">
                    {{ formData.is_hot ? $t('uiAdmin1.rootDomains.form.hot') : $t('uiAdmin1.rootDomains.form.normal') }}
                  </span>
                </label>
              </div>

              <div class="y-field">
                <span class="y-field-label">{{ $t('uiAdmin1.rootDomains.form.markNew') }}</span>
                <label class="flex h-10 cursor-pointer items-center gap-3">
                  <input v-model="formData.is_new" type="checkbox" class="toggle toggle-primary toggle-sm" />
                  <span class="text-sm text-y-muted">
                    {{ formData.is_new ? $t('uiAdmin1.rootDomains.form.isNew') : $t('uiAdmin1.rootDomains.form.regular') }}
                  </span>
                </label>
              </div>
            </div>
          </div>

          <div class="modal-action">
            <button type="button" class="btn btn-ghost" @click="closeModals">{{ $t('common.cancel') }}</button>
            <button type="submit" class="btn btn-primary" :disabled="submitting">
              <span v-if="submitting" class="loading loading-spinner loading-sm"></span>
              <span v-else>{{ showEditModal ? $t('uiAdmin1.rootDomains.form.update') : $t('uiAdmin1.rootDomains.form.create') }}</span>
            </button>
          </div>
        </form>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button type="button" @click="closeModals">{{ $t('common.close') }}</button>
      </form>
    </dialog>

    <!-- DNSSEC Modal -->
    <dialog :class="{ 'modal': true, 'modal-open': showDNSSECModal }">
      <div class="modal-box w-11/12 max-w-2xl">
        <h3 class="flex flex-wrap items-baseline gap-x-2 text-lg font-semibold text-y-text">
          DNSSEC
          <span class="y-domain text-base">.{{ dnssecDomain?.domain }}</span>
        </h3>
        <p class="mt-1 text-sm leading-6 text-y-muted">{{ $t('uiAdmin1.rootDomains.dnssec.intro') }}</p>

        <YLoader v-if="dnssecState.loading" :size="28" />

        <div v-else class="mt-5 space-y-6">
          <div class="space-y-4">
            <div class="flex flex-wrap items-center justify-between gap-3">
              <StatusDot
                :tone="dnssecState.enabled ? 'ok' : 'muted'"
                :label="dnssecState.enabled ? $t('uiAdmin1.rootDomains.dnssecSigned') : $t('uiAdmin1.rootDomains.dnssecUnsigned')"
              />
              <button
                type="button"
                class="btn btn-sm"
                :class="dnssecState.enabled ? 'btn-outline text-y-danger-text' : 'btn-primary'"
                :disabled="dnssecState.submitting"
                @click="toggleRootDNSSEC"
              >
                <span v-if="dnssecState.submitting" class="loading loading-spinner loading-xs"></span>
                <span v-else>{{ dnssecState.enabled ? $t('uiAdmin1.rootDomains.dnssec.disable') : $t('uiAdmin1.rootDomains.dnssec.enable') }}</span>
              </button>
            </div>

            <div v-if="!dnssecState.zoneExists" class="alert alert-warning items-start text-sm">
              <TriangleAlert :size="16" :stroke-width="1.5" class="mt-0.5 shrink-0" />
              <span>{{ $t('uiAdmin1.rootDomains.dnssec.zoneMissing') }}</span>
            </div>

            <div v-if="dnssecState.enabled" class="space-y-3">
              <div class="alert alert-info items-start text-sm">
                <Info :size="16" :stroke-width="1.5" class="mt-0.5 shrink-0" />
                <i18n-t keypath="uiAdmin1.rootDomains.dnssec.dsInstruction" tag="span">
                  <template #domain>
                    <strong class="font-mono font-medium">.{{ dnssecDomain?.domain }}</strong>
                  </template>
                </i18n-t>
              </div>

              <p class="y-label">{{ $t('uiAdmin1.rootDomains.dnssec.dsRecords') }}</p>
              <p v-if="dnssecState.dsRecords.length === 0" class="text-sm text-y-faint">{{ $t('uiAdmin1.rootDomains.dnssec.noDs') }}</p>
              <div v-for="(ds, idx) in dnssecState.dsRecords" :key="idx" class="y-inset flex items-start gap-2 py-2 pl-3 pr-2">
                <code class="flex-1 break-all py-0.5 font-mono text-xs leading-5 text-y-text">{{ ds }}</code>
                <button
                  type="button"
                  class="btn btn-ghost btn-xs btn-square shrink-0"
                  :title="$t('uiAdmin1.rootDomains.dnssec.copy')"
                  :aria-label="$t('uiAdmin1.rootDomains.dnssec.copy')"
                  @click="copyText(ds)"
                >
                  <Copy :size="14" :stroke-width="1.5" />
                </button>
              </div>

              <div v-for="key in dnssecState.keys" :key="key.id" class="space-y-1.5 pt-1">
                <p class="y-label">DNSKEY · {{ key.keytype.toUpperCase() }} · {{ key.algorithm }}</p>
                <div class="y-inset flex items-start gap-2 py-2 pl-3 pr-2">
                  <code class="flex-1 break-all py-0.5 font-mono text-xs leading-5 text-y-text">{{ key.dnskey }}</code>
                  <button
                    type="button"
                    class="btn btn-ghost btn-xs btn-square shrink-0"
                    :title="$t('uiAdmin1.rootDomains.dnssec.copy')"
                    :aria-label="$t('uiAdmin1.rootDomains.dnssec.copy')"
                    @click="copyText(key.dnskey)"
                  >
                    <Copy :size="14" :stroke-width="1.5" />
                  </button>
                </div>
              </div>
            </div>
          </div>

          <div class="y-divider space-y-3 pt-5">
            <h4 class="text-sm font-semibold text-y-text">{{ $t('uiAdmin1.rootDomains.dnssec.delegationsTitle') }}</h4>
            <p class="text-sm leading-6 text-y-muted">{{ $t('uiAdmin1.rootDomains.dnssec.delegationsDesc') }}</p>
            <div class="flex flex-wrap items-center gap-3">
              <button
                type="button"
                class="btn btn-outline btn-sm gap-2"
                :disabled="dnssecState.repairing || dnssecState.repair?.running"
                @click="repairDelegations"
              >
                <span v-if="dnssecState.repairing || dnssecState.repair?.running" class="loading loading-spinner loading-xs"></span>
                <Wrench v-else :size="14" :stroke-width="1.5" />
                {{ $t('uiAdmin1.rootDomains.dnssec.repair') }}
              </button>
              <span v-if="dnssecState.repair" class="font-mono text-xs text-y-faint">
                <template v-if="dnssecState.repair.running">
                  {{ $t('uiAdmin1.rootDomains.dnssec.repairRunning', { total: dnssecState.repair.total }) }}
                </template>
                <template v-else>
                  {{ $t('uiAdmin1.rootDomains.dnssec.repairDone', {
                    delegated: dnssecState.repair.delegated,
                    withDs: dnssecState.repair.with_ds,
                    noZone: dnssecState.repair.no_zone,
                    errors: dnssecState.repair.errors?.length || 0,
                  }) }}
                </template>
              </span>
            </div>
            <div
              v-if="dnssecState.repair && !dnssecState.repair.running && dnssecState.repair.errors?.length"
              class="y-inset max-h-40 space-y-1 overflow-auto p-3"
            >
              <p v-for="(e, i) in dnssecState.repair.errors" :key="i" class="break-all font-mono text-xs text-y-danger-text">{{ e }}</p>
            </div>
          </div>

          <div class="y-divider space-y-3 pt-5">
            <h4 class="text-sm font-semibold text-y-text">{{ $t('uiAdmin1.rootDomains.dnssec.orphansTitle') }}</h4>
            <p class="text-sm leading-6 text-y-muted">{{ $t('uiAdmin1.rootDomains.dnssec.orphansDesc') }}</p>
            <div class="flex flex-wrap items-center gap-3">
              <button type="button" class="btn btn-outline btn-sm gap-2" :disabled="orphans.scanning" @click="scanOrphans">
                <span v-if="orphans.scanning" class="loading loading-spinner loading-xs"></span>
                <ScanSearch v-else :size="14" :stroke-width="1.5" />
                {{ $t('uiAdmin1.rootDomains.dnssec.scanOrphans') }}
              </button>
              <button
                v-if="orphans.zones.length > 0"
                type="button"
                class="btn btn-error btn-sm gap-2"
                :disabled="orphans.cleaning || dnssecState.orphanCleanup?.running"
                @click="cleanupOrphans"
              >
                <span v-if="orphans.cleaning || dnssecState.orphanCleanup?.running" class="loading loading-spinner loading-xs"></span>
                <Trash2 v-else :size="14" :stroke-width="1.5" />
                {{ $t('uiAdmin1.rootDomains.dnssec.deleteOrphans', { count: orphans.zones.length }) }}
              </button>
              <span v-if="orphans.scanned" class="font-mono text-xs text-y-faint">
                {{ $t('uiAdmin1.rootDomains.dnssec.orphansFound', { count: orphans.zones.length }) }}
              </span>
              <span v-if="dnssecState.orphanCleanup" class="font-mono text-xs text-y-faint">
                <template v-if="dnssecState.orphanCleanup.running">
                  {{ $t('uiAdmin1.rootDomains.dnssec.cleaning', { deleted: dnssecState.orphanCleanup.deleted, total: dnssecState.orphanCleanup.total }) }}
                </template>
                <template v-else>
                  {{ $t('uiAdmin1.rootDomains.dnssec.cleanupDone', {
                    deleted: dnssecState.orphanCleanup.deleted,
                    total: dnssecState.orphanCleanup.total,
                    errors: dnssecState.orphanCleanup.errors?.length || 0,
                  }) }}
                </template>
              </span>
            </div>
            <div v-if="orphans.zones.length > 0" class="y-inset max-h-32 overflow-auto p-3 font-mono text-xs leading-5 text-y-muted">
              {{ orphans.zones.slice(0, 200).join(', ') }}<span v-if="orphans.zones.length > 200" class="text-y-faint"> {{ $t('uiAdmin1.rootDomains.dnssec.more', { count: orphans.zones.length - 200 }) }}</span>
            </div>
            <div
              v-if="dnssecState.orphanCleanup && !dnssecState.orphanCleanup.running && dnssecState.orphanCleanup.errors?.length"
              class="y-inset max-h-32 space-y-1 overflow-auto p-3"
            >
              <p v-for="(e, i) in dnssecState.orphanCleanup.errors" :key="i" class="break-all font-mono text-xs text-y-danger-text">{{ e }}</p>
            </div>
          </div>
        </div>

        <div class="modal-action">
          <button type="button" class="btn btn-ghost" @click="closeDNSSEC">{{ $t('common.close') }}</button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button type="button" @click="closeDNSSEC">{{ $t('common.close') }}</button>
      </form>
    </dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  Plus, Flame, Sparkles, Users, ArrowUpRight, ShieldCheck, Pencil, Trash2, Info, X, TriangleAlert, Copy,
  Wrench, ScanSearch,
} from 'lucide-vue-next'
import axios from '../utils/axios'
import { useToast } from '../composables/useToast'
import { useCurrency } from '../composables/useCurrency'
import PageHeader from '../components/ui/PageHeader.vue'
import StatusDot from '../components/ui/StatusDot.vue'
import EmptyState from '../components/ui/EmptyState.vue'
import YLoader from '../components/brand/YLoader.vue'

const { t } = useI18n()
const toast = useToast()
const { formatPrice, currencySymbol } = useCurrency()
const loading = ref(true)
const submitting = ref(false)
const rootDomains = ref([])
const showCreateModal = ref(false)
const showEditModal = ref(false)
const editingDomain = ref(null)

const formData = ref({
  domain: '',
  description: '',
  priority: 0,
  is_active: true,
  is_hot: false,
  is_new: false,
  is_free: true,
  min_user_level: 'normal',
  price_per_year: null,
  lifetime_price: null,
  use_default_nameservers: true,
  nameservers: [''],
})

// DNSSEC state per root domain (id -> status) plus the modal state
const dnssecStatuses = ref({})
const showDNSSECModal = ref(false)
const dnssecDomain = ref(null)
const dnssecState = ref({
  loading: false,
  submitting: false,
  repairing: false,
  repair: null,
  orphanCleanup: null,
  zoneExists: false,
  enabled: false,
  keys: [],
  dsRecords: [],
})
let repairPollTimer = null

const pollRepair = (domainId) => {
  clearTimeout(repairPollTimer)
  repairPollTimer = setTimeout(async () => {
    if (!showDNSSECModal.value || dnssecDomain.value?.id !== domainId) return
    try {
      const response = await axios.get(`/api/admin/root-domains/${domainId}/dnssec`)
      applyRootDNSSEC(response.data)
      if (response.data.repair?.running || response.data.orphan_cleanup?.running) pollRepair(domainId)
    } catch (error) {
      console.error('Failed to poll repair status:', error)
    }
  }, 3000)
}

const orphans = ref({ scanning: false, cleaning: false, scanned: false, zones: [] })

const scanOrphans = async () => {
  const domain = dnssecDomain.value
  if (!domain) return
  orphans.value.scanning = true
  try {
    const response = await axios.get(`/api/admin/root-domains/${domain.id}/orphan-zones`)
    orphans.value.zones = response.data.zones || []
    orphans.value.scanned = true
  } catch (error) {
    toast.error(error.response?.data?.error || 'Failed to scan orphan zones')
  } finally {
    orphans.value.scanning = false
  }
}

const cleanupOrphans = async () => {
  const domain = dnssecDomain.value
  if (!domain) return
  if (!confirm(`Delete ${orphans.value.zones.length} orphan zone(s) under .${domain.domain} from PowerDNS?\n\nThese names have no domain in the database. This cannot be undone.`)) return
  orphans.value.cleaning = true
  try {
    const response = await axios.post(`/api/admin/root-domains/${domain.id}/orphan-zones/cleanup`)
    dnssecState.value.orphanCleanup = response.data.orphan_cleanup
    toast.success(response.data.message || 'Cleanup started')
    orphans.value.zones = []
    orphans.value.scanned = false
    pollRepair(domain.id)
  } catch (error) {
    if (error.response?.data?.orphan_cleanup) dnssecState.value.orphanCleanup = error.response.data.orphan_cleanup
    toast.error(error.response?.data?.error || 'Failed to start cleanup')
  } finally {
    orphans.value.cleaning = false
  }
}

const repairDelegations = async () => {
  const domain = dnssecDomain.value
  if (!domain) return
  if (!confirm(`Publish NS/DS delegations for all subdomains of .${domain.domain} into its zone?`)) return
  dnssecState.value.repairing = true
  try {
    const response = await axios.post(`/api/admin/root-domains/${domain.id}/dnssec/repair`)
    dnssecState.value.repair = response.data.repair
    toast.success(response.data.message || 'Repair started')
    pollRepair(domain.id)
  } catch (error) {
    if (error.response?.data?.repair) dnssecState.value.repair = error.response.data.repair
    toast.error(error.response?.data?.error || 'Failed to start repair')
  } finally {
    dnssecState.value.repairing = false
  }
}

onMounted(async () => {
  await fetchRootDomains()
})

const fetchRootDomains = async () => {
  loading.value = true
  try {
    const response = await axios.get('/api/admin/root-domains')
    rootDomains.value = response.data.root_domains || []
    fetchDNSSECStatuses()
  } catch (error) {
    console.error('Failed to fetch root domains:', error)
    toast.error(error.response?.data?.error || 'Failed to fetch root domains')
  } finally {
    loading.value = false
  }
}

const createDomain = async () => {
  submitting.value = true
  try {
    // Filter out empty nameservers
    const nameservers = formData.value.nameservers.filter(ns => ns.trim() !== '')

    const payload = {
      ...formData.value,
      nameservers: formData.value.use_default_nameservers ? [] : nameservers,
    }

    await axios.post('/api/admin/root-domains', payload)
    toast.success('Root domain created successfully!')
    closeModals()
    await fetchRootDomains()
  } catch (error) {
    toast.error(error.response?.data?.error || 'Failed to create root domain')
  } finally {
    submitting.value = false
  }
}

const editDomain = (domain) => {
  editingDomain.value = domain

  // Parse nameservers from JSON string
  let nameservers = ['']
  try {
    if (domain.nameservers) {
      nameservers = JSON.parse(domain.nameservers)
      if (!Array.isArray(nameservers) || nameservers.length === 0) {
        nameservers = ['']
      }
    }
  } catch (e) {
    console.error('Failed to parse nameservers:', e)
    nameservers = ['']
  }

  formData.value = {
    domain: domain.domain,
    description: domain.description || '',
    priority: domain.priority,
    is_active: domain.is_active,
    is_hot: domain.is_hot,
    is_new: domain.is_new,
    is_free: domain.is_free ?? true,
    min_user_level: domain.min_user_level || 'normal',
    price_per_year: domain.price_per_year || null,
    lifetime_price: domain.lifetime_price || null,
    use_default_nameservers: domain.use_default_nameservers ?? true,
    nameservers: nameservers,
  }
  showEditModal.value = true
}

const updateDomain = async () => {
  submitting.value = true
  try {
    // Filter out empty nameservers
    const nameservers = formData.value.nameservers.filter(ns => ns.trim() !== '')

    await axios.put(`/api/admin/root-domains/${editingDomain.value.id}`, {
      description: formData.value.description || null,
      priority: formData.value.priority,
      is_active: formData.value.is_active,
      is_hot: formData.value.is_hot,
      is_new: formData.value.is_new,
      is_free: formData.value.is_free,
      min_user_level: formData.value.min_user_level || 'normal',
      price_per_year: formData.value.is_free ? null : (formData.value.price_per_year || null),
      lifetime_price: formData.value.is_free ? null : (formData.value.lifetime_price || null),
      use_default_nameservers: formData.value.use_default_nameservers,
      nameservers: formData.value.use_default_nameservers ? [] : nameservers,
    })
    toast.success('Root domain updated successfully!')
    closeModals()
    await fetchRootDomains()
  } catch (error) {
    toast.error(error.response?.data?.error || 'Failed to update root domain')
  } finally {
    submitting.value = false
  }
}

const deleteDomain = async (domain) => {
  if (!confirm(`Are you sure you want to delete .${domain.domain}?\n\nThis cannot be undone.`)) {
    return
  }

  try {
    await axios.delete(`/api/admin/root-domains/${domain.id}`)
    toast.success('Root domain deleted successfully!')
    await fetchRootDomains()
  } catch (error) {
    toast.error(error.response?.data?.error || 'Failed to delete root domain')
  }
}

const fetchDNSSECStatuses = async () => {
  await Promise.all(
    rootDomains.value
      .filter(d => d.use_default_nameservers !== false)
      .map(async (d) => {
        try {
          const response = await axios.get(`/api/admin/root-domains/${d.id}/dnssec`)
          dnssecStatuses.value[d.id] = response.data
        } catch (error) {
          console.error(`Failed to fetch DNSSEC status for ${d.domain}:`, error)
        }
      })
  )
}

const dnssecBadgeClass = (id) => {
  const status = dnssecStatuses.value[id]
  if (!status) return 'y-chip'
  return status.enabled ? 'y-chip-mint' : 'y-chip-warn'
}

const dnssecBadgeText = (id) => {
  const status = dnssecStatuses.value[id]
  if (!status) return t('uiAdmin1.rootDomains.dnssecUnknown')
  return status.enabled ? t('uiAdmin1.rootDomains.dnssecSigned') : t('uiAdmin1.rootDomains.dnssecUnsigned')
}

const applyRootDNSSEC = (data) => {
  dnssecState.value.zoneExists = !!data.zone_exists
  dnssecState.value.enabled = !!data.enabled
  dnssecState.value.keys = data.keys || []
  dnssecState.value.dsRecords = data.ds_records || []
  if (data.repair !== undefined) dnssecState.value.repair = data.repair
  if (data.orphan_cleanup !== undefined) dnssecState.value.orphanCleanup = data.orphan_cleanup
  if (dnssecDomain.value) {
    dnssecStatuses.value[dnssecDomain.value.id] = data
  }
}

const openDNSSEC = async (domain) => {
  dnssecDomain.value = domain
  showDNSSECModal.value = true
  dnssecState.value.loading = true
  dnssecState.value.repair = null
  dnssecState.value.orphanCleanup = null
  orphans.value = { scanning: false, cleaning: false, scanned: false, zones: [] }
  try {
    const response = await axios.get(`/api/admin/root-domains/${domain.id}/dnssec`)
    applyRootDNSSEC(response.data)
    if (response.data.repair?.running || response.data.orphan_cleanup?.running) pollRepair(domain.id)
  } catch (error) {
    toast.error(error.response?.data?.error || 'Failed to fetch DNSSEC status')
  } finally {
    dnssecState.value.loading = false
  }
}

const closeDNSSEC = () => {
  clearTimeout(repairPollTimer)
  showDNSSECModal.value = false
  dnssecDomain.value = null
}

const toggleRootDNSSEC = async () => {
  const domain = dnssecDomain.value
  if (!domain) return
  if (dnssecState.value.enabled) {
    if (!confirm(`Disable DNSSEC for .${domain.domain}?\n\nRemove the DS record at the registrar FIRST and wait for its TTL to expire. Otherwise .${domain.domain} and every subdomain under it will fail validation.`)) {
      return
    }
  }
  dnssecState.value.submitting = true
  try {
    const action = dnssecState.value.enabled ? 'disable' : 'enable'
    const response = await axios.post(`/api/admin/root-domains/${domain.id}/dnssec/${action}`)
    applyRootDNSSEC(response.data)
    toast.success(response.data.message || `DNSSEC ${action}d`)
  } catch (error) {
    toast.error(error.response?.data?.error || 'DNSSEC operation failed')
  } finally {
    dnssecState.value.submitting = false
  }
}

const copyText = async (text) => {
  try {
    await navigator.clipboard.writeText(text)
    toast.success('Copied!')
  } catch {
    toast.error('Failed to copy')
  }
}

const closeModals = () => {
  showCreateModal.value = false
  showEditModal.value = false
  editingDomain.value = null
  formData.value = {
    domain: '',
    description: '',
    priority: 0,
    is_active: true,
    is_hot: false,
    is_new: false,
    is_free: true,
    min_user_level: 'normal',
    price_per_year: null,
    lifetime_price: null,
    use_default_nameservers: true,
    nameservers: [''],
  }
}
</script>
