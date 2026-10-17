<template>
  <div class="y-page">
    <PageHeader
      :eyebrow="$t('shell.adminArea')"
      :title="$t('uiAdmin3.settings.title')"
      :description="$t('uiAdmin3.settings.description')"
    >
      <template #actions>
        <button type="button" class="btn btn-outline btn-sm gap-1.5" :disabled="sysLoading" @click="refreshSystemInfo">
          <RefreshCw :size="16" :stroke-width="1.5" :class="{ 'animate-spin': sysLoading }" />
          {{ $t('uiAdmin3.settings.refresh') }}
        </button>
      </template>
    </PageHeader>

    <!-- Site and sign-up -->
    <section class="grid gap-4 lg:grid-cols-[16rem_minmax(0,1fr)] lg:gap-10">
      <div class="space-y-1.5">
        <h2 class="y-h3">{{ $t('uiAdmin3.settings.siteTitle') }}</h2>
        <p class="text-sm leading-6 text-y-muted">{{ $t('uiAdmin3.settings.siteDesc') }}</p>
      </div>
      <div class="y-panel">
        <dl>
          <div class="flex flex-col gap-1 border-b border-y-line px-5 py-3.5 sm:flex-row sm:items-center sm:justify-between sm:gap-6">
            <dt class="text-sm text-y-muted">{{ $t('adminSettings.siteName') }}</dt>
            <dd class="text-sm font-medium text-y-text sm:text-right">{{ sysInfo.site?.name || '-' }}</dd>
          </div>
          <div class="flex flex-col gap-1 border-b border-y-line px-5 py-3.5 sm:flex-row sm:items-center sm:justify-between sm:gap-6">
            <dt class="shrink-0 text-sm text-y-muted">{{ $t('adminSettings.siteDescription') }}</dt>
            <dd class="text-sm text-y-text sm:text-right">{{ sysInfo.site?.description || '-' }}</dd>
          </div>
        </dl>
        <p class="flex items-start gap-2 border-b border-y-line px-5 py-3 text-xs leading-5 text-y-faint">
          <Info :size="14" :stroke-width="1.5" class="mt-0.5 shrink-0" />
          {{ $t('adminSettings.siteConfigEnvHint') }}
        </p>
        <label class="flex cursor-pointer items-start justify-between gap-6 px-5 py-4">
          <span class="min-w-0 space-y-1">
            <span class="block text-sm font-medium text-y-text">{{ $t('adminSettings.allowPasswordRegister') }}</span>
            <span class="block text-sm leading-6 text-y-muted">{{ $t('adminSettings.allowPasswordRegisterDesc') }}</span>
          </span>
          <input
            type="checkbox"
            class="toggle toggle-primary mt-0.5 shrink-0"
            :checked="siteSettings.allow_password_register !== 'false'"
            @change="togglePasswordRegister"
          />
        </label>
      </div>
    </section>

    <!-- Domain quotas per level -->
    <section class="grid gap-4 border-t border-y-line pt-8 lg:grid-cols-[16rem_minmax(0,1fr)] lg:gap-10">
      <div class="space-y-1.5">
        <h2 class="y-h3">{{ $t('uiAdmin3.settings.quotaTitle') }}</h2>
        <p class="text-sm leading-6 text-y-muted">{{ $t('adminSettings.quotaSettingsDesc') }}</p>
      </div>
      <div class="y-panel">
        <YLoader v-if="quotaLoading" :size="32" />
        <template v-else>
          <div class="overflow-x-auto">
            <table class="y-table">
              <thead>
                <tr>
                  <th>{{ $t('adminSettings.level') }}</th>
                  <th>{{ $t('adminSettings.trustLevel') }}</th>
                  <th>{{ $t('adminSettings.domainQuota') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="level in quotaLevels" :key="level.key">
                  <td class="whitespace-nowrap">
                    <span :class="level.badgeClass">{{ $t('home.levelNames.' + level.key.replace('quota_', '')) }}</span>
                  </td>
                  <td class="font-mono text-sm text-y-muted">{{ level.trustLevel }}</td>
                  <td>
                    <div class="inline-flex items-center gap-2">
                      <input
                        v-model.number="level.value"
                        type="number"
                        min="0"
                        class="input input-bordered input-sm w-24 font-mono"
                        :aria-label="$t('adminSettings.domainQuota')"
                        @focus="onQuotaFocus(level)"
                        @blur="onQuotaBlur(level)"
                      />
                      <span v-if="level.saving" class="loading loading-spinner loading-xs"></span>
                      <Check v-else-if="level.saved" :size="16" :stroke-width="2" class="text-y-mint-text" :aria-label="$t('adminSettings.saved')" />
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
          <p class="border-t border-y-line px-5 py-3 text-xs text-y-faint">{{ $t('uiAdmin3.settings.autoSave') }}</p>
        </template>
      </div>
    </section>

    <!-- Reserved subdomains -->
    <section class="grid gap-4 border-t border-y-line pt-8 lg:grid-cols-[16rem_minmax(0,1fr)] lg:gap-10">
      <div class="space-y-1.5">
        <h2 class="y-h3">{{ $t('uiAdmin3.settings.blacklistTitle') }}</h2>
        <p class="text-sm leading-6 text-y-muted">{{ $t('uiAdmin3.settings.blacklistDesc') }}</p>
      </div>
      <div class="y-panel p-5">
        <label class="y-field">
          <span class="y-field-label">{{ $t('uiAdmin3.settings.blacklistLabel') }}</span>
          <textarea
            v-model="subdomainBlacklist"
            class="textarea textarea-bordered w-full font-mono text-sm"
            rows="4"
            placeholder="admin,root,api,www,mail,smtp,ftp,ssh,dns,..."
            spellcheck="false"
            @blur="updateBlacklistSetting"
          />
          <span class="y-field-hint">{{ $t('uiAdmin3.settings.blacklistHint') }}</span>
        </label>
        <p v-if="blacklistSaving" class="mt-3 inline-flex items-center gap-1.5 text-xs text-y-muted">
          <span class="loading loading-spinner loading-xs"></span>
          {{ $t('adminSettings.saving') }}
        </p>
        <p v-else-if="blacklistSaved" class="mt-3 inline-flex items-center gap-1.5 text-xs text-y-mint-text">
          <Check :size="14" :stroke-width="2" />
          {{ $t('adminSettings.saved') }}
        </p>
      </div>
    </section>

    <!-- Currency -->
    <section class="grid gap-4 border-t border-y-line pt-8 lg:grid-cols-[16rem_minmax(0,1fr)] lg:gap-10">
      <div class="space-y-1.5">
        <h2 class="y-h3">{{ $t('uiAdmin3.settings.currencyTitle') }}</h2>
        <p class="text-sm leading-6 text-y-muted">{{ $t('adminSettings.currencySymbolDesc') }}</p>
      </div>
      <div class="y-panel p-5">
        <label class="y-field">
          <span class="y-field-label">{{ $t('adminSettings.currencySymbol') }}</span>
          <input
            v-model="currencySymbol"
            type="text"
            maxlength="10"
            class="input input-bordered w-full max-w-xs font-mono"
            placeholder="NL"
            @blur="updateCurrencySetting"
          />
          <span class="y-field-hint font-mono">{{ $t('adminSettings.currencyExample', { symbol: currencySymbol || 'NL' }) }}</span>
        </label>
        <p v-if="currencySaving" class="mt-3 inline-flex items-center gap-1.5 text-xs text-y-muted">
          <span class="loading loading-spinner loading-xs"></span>
          {{ $t('adminSettings.saving') }}
        </p>
        <p v-else-if="currencySaved" class="mt-3 inline-flex items-center gap-1.5 text-xs text-y-mint-text">
          <Check :size="14" :stroke-width="2" />
          {{ $t('adminSettings.saved') }}
        </p>
      </div>
    </section>

    <!-- System information and statistics -->
    <section class="grid gap-4 border-t border-y-line pt-8 lg:grid-cols-[16rem_minmax(0,1fr)] lg:gap-10">
      <div class="space-y-1.5">
        <h2 class="y-h3">{{ $t('uiAdmin3.settings.systemTitle') }}</h2>
        <p class="text-sm leading-6 text-y-muted">{{ $t('uiAdmin3.settings.systemDesc') }}</p>
      </div>

      <div v-if="sysLoading" class="y-panel">
        <YLoader :size="32" :label="$t('common.loading')" />
      </div>

      <div v-else class="space-y-4">
        <div class="y-panel">
          <div class="y-panel-header">
            <h3 class="text-sm font-semibold text-y-text">{{ $t('adminSettings.systemInfo') }}</h3>
          </div>
          <dl>
            <div class="flex items-center justify-between gap-6 border-b border-y-line px-5 py-3">
              <dt class="text-sm text-y-muted">{{ $t('adminSettings.systemName') }}</dt>
              <dd class="text-right text-sm font-medium text-y-text">{{ sysInfo.site?.name || '-' }}</dd>
            </div>
            <div class="flex items-center justify-between gap-6 border-b border-y-line px-5 py-3">
              <dt class="text-sm text-y-muted">{{ $t('adminSettings.version') }}</dt>
              <dd class="break-all text-right font-mono text-sm text-y-text">{{ sysInfo.system?.version || '-' }}</dd>
            </div>
            <div class="flex items-center justify-between gap-6 border-b border-y-line px-5 py-3">
              <dt class="text-sm text-y-muted">{{ $t('adminSettings.environment') }}</dt>
              <dd>
                <StatusDot
                  class="font-mono"
                  :tone="sysInfo.system?.environment === 'production' ? 'ok' : 'warn'"
                  :label="sysInfo.system?.environment || '-'"
                />
              </dd>
            </div>
            <div class="flex items-center justify-between gap-6 border-b border-y-line px-5 py-3">
              <dt class="text-sm text-y-muted">{{ $t('adminSettings.goVersion') }}</dt>
              <dd class="text-right font-mono text-sm text-y-text">{{ sysInfo.system?.go_version || '-' }}</dd>
            </div>
            <div class="flex items-center justify-between gap-6 border-b border-y-line px-5 py-3">
              <dt class="text-sm text-y-muted">{{ $t('adminSettings.database') }}</dt>
              <dd class="flex flex-wrap items-center justify-end gap-x-3 gap-y-1">
                <span class="font-mono text-sm text-y-muted">PostgreSQL</span>
                <StatusDot
                  :tone="sysInfo.services?.database === 'connected' ? 'ok' : 'danger'"
                  :label="sysInfo.services?.database === 'connected' ? $t('uiAdmin3.settings.connected') : $t('uiAdmin3.settings.disconnected')"
                />
              </dd>
            </div>
            <div class="flex items-center justify-between gap-6 border-b border-y-line px-5 py-3">
              <dt class="text-sm text-y-muted">{{ $t('adminSettings.cache') }}</dt>
              <dd class="flex flex-wrap items-center justify-end gap-x-3 gap-y-1">
                <span class="font-mono text-sm text-y-muted">Redis</span>
                <StatusDot
                  :tone="sysInfo.services?.redis === 'connected' ? 'ok' : 'danger'"
                  :label="sysInfo.services?.redis === 'connected' ? $t('uiAdmin3.settings.connected') : $t('uiAdmin3.settings.disconnected')"
                />
              </dd>
            </div>
            <div class="flex items-center justify-between gap-6 px-5 py-3">
              <dt class="text-sm text-y-muted">{{ $t('adminSettings.platform') }}</dt>
              <dd class="text-right font-mono text-sm text-y-text">{{ sysInfo.system?.os }}/{{ sysInfo.system?.arch }}</dd>
            </div>
          </dl>
        </div>

        <div class="y-panel overflow-hidden">
          <div class="y-panel-header">
            <h3 class="text-sm font-semibold text-y-text">{{ $t('adminSettings.systemStats') }}</h3>
          </div>
          <div class="grid grid-cols-2 gap-px bg-y-line sm:grid-cols-4">
            <div v-for="stat in systemStats" :key="stat.key" class="space-y-2 bg-y-surface px-5 py-4">
              <p class="y-label">{{ stat.label }}</p>
              <p class="break-all font-mono text-xl font-medium tracking-tight text-y-text">
                {{ stat.value }}<span v-if="stat.unit" class="text-sm text-y-faint"> {{ stat.unit }}</span>
              </p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Default nameservers -->
    <section class="grid gap-4 border-t border-y-line pt-8 lg:grid-cols-[16rem_minmax(0,1fr)] lg:gap-10">
      <div class="space-y-1.5">
        <h2 class="y-h3">{{ $t('uiAdmin3.settings.nameserversTitle') }}</h2>
        <p class="text-sm leading-6 text-y-muted">{{ $t('adminSettings.defaultDNSDesc') }}</p>
      </div>
      <div class="y-panel">
        <YLoader v-if="sysLoading" :size="32" />
        <dl v-else>
          <div class="flex items-center gap-3 border-b border-y-line px-5 py-3.5">
            <dt><span class="y-tag">NS1</span></dt>
            <dd class="y-domain text-sm">{{ sysInfo.dns?.ns1 || '-' }}</dd>
          </div>
          <div class="flex items-center gap-3 px-5 py-3.5">
            <dt><span class="y-tag">NS2</span></dt>
            <dd class="y-domain text-sm">{{ sysInfo.dns?.ns2 || '-' }}</dd>
          </div>
        </dl>
      </div>
    </section>

    <!-- Integrations: OAuth, payment, email -->
    <section class="grid gap-4 border-t border-y-line pt-8 lg:grid-cols-[16rem_minmax(0,1fr)] lg:gap-10">
      <div class="space-y-1.5">
        <h2 class="y-h3">{{ $t('uiAdmin3.settings.integrationsTitle') }}</h2>
        <p class="text-sm leading-6 text-y-muted">{{ $t('uiAdmin3.settings.integrationsDesc') }}</p>
      </div>

      <div v-if="sysLoading" class="y-panel">
        <YLoader :size="32" />
      </div>

      <div v-else class="space-y-4">
        <!-- OAuth -->
        <div class="y-panel">
          <div class="y-panel-header">
            <h3 class="text-sm font-semibold text-y-text">{{ $t('adminSettings.oauthSettings') }}</h3>
          </div>
          <div class="flex items-center justify-between gap-4 border-b border-y-line px-5 py-3">
            <div class="flex items-center gap-2.5 text-sm text-y-text">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="text-y-muted" aria-hidden="true">
                <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
              </svg>
              GitHub
            </div>
            <StatusDot
              :tone="sysInfo.oauth?.github ? 'ok' : 'muted'"
              :label="sysInfo.oauth?.github ? $t('adminSettings.configured') : $t('adminSettings.notConfigured')"
            />
          </div>
          <div class="flex items-center justify-between gap-4 border-b border-y-line px-5 py-3">
            <div class="flex items-center gap-2.5 text-sm text-y-text">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" aria-hidden="true">
                <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92a5.06 5.06 0 0 1-2.2 3.32v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.1z"/>
                <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
                <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
                <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
              </svg>
              Google
            </div>
            <StatusDot
              :tone="sysInfo.oauth?.google ? 'ok' : 'muted'"
              :label="sysInfo.oauth?.google ? $t('adminSettings.configured') : $t('adminSettings.notConfigured')"
            />
          </div>
          <div class="flex items-center justify-between gap-4 px-5 py-3">
            <div class="flex items-center gap-2.5 text-sm text-y-text">
              <svg viewBox="0 0 102 64" xmlns="http://www.w3.org/2000/svg" width="16" height="10" aria-hidden="true">
                <path d="M51.504 19.32C51.504 23.224 51.312 28.888 50.928 36.312C50.544 43.608 50.352 49.144 50.352 52.92C50.352 54.008 50.416 55.64 50.544 57.816C50.672 59.992 50.736 61.688 50.736 62.904C50.608 63.288 49.36 63.64 46.992 63.96C44.624 64.28 42.992 64.44 42.096 64.44C41.648 64.44 40.944 63.832 39.984 62.616C39.024 61.4 37.808 59.704 36.336 57.528C33.904 54.072 32.368 52.088 31.728 51.576C31.728 51.576 31.216 50.808 30.192 49.272C29.104 47.544 27.76 45.496 26.16 43.128C24.624 40.76 23.856 39.544 23.856 39.48L15.504 29.592C15.44 29.464 15.28 29.112 15.024 28.536C14.768 27.96 14.416 27.672 13.968 27.672C13.776 27.672 13.584 27.864 13.392 28.248C13.264 28.568 13.2 28.92 13.2 29.304C13.2 36.024 13.264 41.304 13.392 45.144C13.456 46.04 13.488 47.16 13.488 48.504C13.488 48.696 13.456 48.92 13.392 49.176C13.328 49.368 13.264 49.528 13.2 49.656V64.44L12.432 65.208L10.32 65.304L0.528 64.536C0.528 62.744 0.592 60.632 0.72 58.2C0.912 55.768 1.008 54.328 1.008 53.88C1.008 52.728 0.944 51.128 0.816 49.08C0.688 46.776 0.624001 44.792 0.624001 43.128L0.816 39.864C1.008 36.216 1.104 34.2 1.104 33.816C1.104 32.28 1.04 29.976 0.912 26.904C0.784 23.832 0.72 21.56 0.72 20.088C0.72 19.192 0.848 17.912 1.104 16.248C1.36 14.712 1.488 13.592 1.488 12.888C1.488 12.696 1.36 11.736 1.104 10.008C0.848 8.216 0.72 6.68 0.72 5.4L0.336 4.72799C0.464 4.088 0.688 3.672 1.008 3.48C1.456 3.416 2.672 3.224 4.656 2.904C6.64 2.52 8.24 2.328 9.456 2.328C10.544 2.328 11.28 2.456 11.664 2.712C12.112 3.096 12.816 3.768 13.776 4.72799C14.736 5.688 15.184 6.264 15.12 6.456L20.016 14.808C20.016 14.872 20.912 16.152 22.704 18.648L24.24 20.856C27.44 25.592 30.256 29.656 32.688 33.048C35.12 36.376 36.688 38.04 37.392 38.04C37.84 38.04 38.096 37.88 38.16 37.56C38.16 36.536 38.096 35.128 37.968 33.336L37.872 29.592C37.872 28.44 37.936 26.648 38.064 24.216L38.16 18.84C38.16 17.368 38.096 15.16 37.968 12.216C37.84 9.4 37.776 7.224 37.776 5.688C37.776 4.152 37.872 3.064 38.064 2.424C38.32 1.784 38.8 1.336 39.504 1.08C40.208 0.759996 41.36 0.599997 42.96 0.599997C46.16 0.599997 48.144 0.759996 48.912 1.08C49.68 1.336 50.096 1.816 50.16 2.51999C50.224 3.16 50.416 3.704 50.736 4.152V10.2L51.024 11.064L50.928 14.328C50.928 16.76 51.12 18.424 51.504 19.32ZM100.475 49.368C100.795 49.368 101.083 49.464 101.339 49.656C101.659 49.784 101.819 50.008 101.819 50.328V60.6C101.819 61.24 101.595 61.784 101.147 62.232C100.699 62.616 100.251 62.808 99.8025 62.808L94.7145 62.616C92.0905 62.616 90.2025 62.712 89.0505 62.904L86.7465 63C86.7465 62.872 85.8185 62.744 83.9625 62.616L78.9705 62.232C77.6265 62.232 76.2825 62.424 74.9385 62.808C73.5945 63.192 72.9225 63.416 72.9225 63.48C71.8345 63.48 70.5865 63.32 69.1785 63C67.5145 62.744 66.2985 62.616 65.5305 62.616C65.5305 61.4 65.4025 59.512 65.1465 56.952C64.8905 54.264 64.7625 52.216 64.7625 50.808C64.7625 50.808 64.7625 50.68 64.7625 50.424C64.8265 50.168 64.8585 49.528 64.8585 48.504C64.8585 46.968 64.8265 45.528 64.7625 44.184L64.2825 43.704V27.96C63.9625 26.04 63.8025 24.344 63.8025 22.872C63.8025 21.464 63.9625 20.568 64.2825 20.184V11.256C64.0905 11.128 63.9945 10.584 63.9945 9.62399V8.184C63.9945 6.2 63.8345 4.536 63.5145 3.192L63.2265 2.808C63.2265 2.168 64.2185 1.656 66.2025 1.272C68.1865 0.823997 70.9065 0.599997 74.3625 0.599997C74.8745 0.599997 75.2905 0.759996 75.6105 1.08C75.9305 1.336 76.0905 1.656 76.0905 2.04C76.0905 2.744 75.9945 3.576 75.8025 4.536C75.3545 6.904 75.1305 9.144 75.1305 11.256C75.1305 11.32 75.0665 12.44 74.9385 14.616C74.8105 16.728 74.6185 18.584 74.3625 20.184L74.9385 20.664V32.184C75.3225 32.888 75.6425 34.488 75.8985 36.984C76.2185 39.416 76.4105 40.92 76.4745 41.496C76.7945 44.568 77.0825 46.936 77.3385 48.6C77.5945 50.2 77.9145 51.192 78.2985 51.576C78.5545 51.704 78.8745 51.832 79.2585 51.96C79.6425 52.088 79.9625 52.184 80.2185 52.248L83.6745 52.152C86.3625 51.96 90.0745 51.864 94.8105 51.864C95.0665 51.8 95.3865 51.736 95.7705 51.672C96.2185 51.608 96.5705 51.576 96.8265 51.576C98.0425 51.512 98.9385 51.352 99.5145 51.096C100.091 50.84 100.411 50.264 100.475 49.368Z" fill="url(#paint0_settings)"/>
                <defs>
                  <linearGradient id="paint0_settings" x1="-6" y1="35" x2="104" y2="35" gradientUnits="userSpaceOnUse">
                    <stop stop-color="#009966"/>
                    <stop offset="0.790909" stop-color="#FF9933"/>
                  </linearGradient>
                </defs>
              </svg>
              NodeLoc
            </div>
            <StatusDot
              :tone="sysInfo.oauth?.nodeloc ? 'ok' : 'muted'"
              :label="sysInfo.oauth?.nodeloc ? $t('adminSettings.configured') : $t('adminSettings.notConfigured')"
            />
          </div>
        </div>

        <!-- Payment -->
        <div class="y-panel">
          <div class="y-panel-header">
            <h3 class="text-sm font-semibold text-y-text">{{ $t('adminSettings.paymentSettings') }}</h3>
          </div>
          <div class="flex items-center justify-between gap-4 px-5 py-3" :class="{ 'border-b border-y-line': sysInfo.payment?.configured }">
            <span class="text-sm text-y-text">NodeLoc Pay</span>
            <StatusDot
              :tone="sysInfo.payment?.configured ? 'ok' : 'muted'"
              :label="sysInfo.payment?.configured ? $t('adminSettings.configured') : $t('adminSettings.notConfigured')"
            />
          </div>
          <div v-if="sysInfo.payment?.configured" class="flex items-center justify-between gap-4 px-5 py-3">
            <span class="text-sm text-y-muted">{{ $t('adminSettings.testMode') }}</span>
            <StatusDot
              :tone="sysInfo.payment?.test_mode ? 'warn' : 'ok'"
              :label="sysInfo.payment?.test_mode ? $t('common.yes') : $t('common.no')"
            />
          </div>
        </div>

        <!-- Email -->
        <div class="y-panel">
          <div class="y-panel-header">
            <h3 class="text-sm font-semibold text-y-text">{{ $t('adminSettings.emailSettings') }}</h3>
          </div>
          <dl>
            <div class="flex items-center justify-between gap-4 px-5 py-3" :class="{ 'border-b border-y-line': sysInfo.email?.configured }">
              <dt class="text-sm text-y-muted">{{ $t('adminSettings.status') }}</dt>
              <dd>
                <StatusDot
                  :tone="sysInfo.email?.configured ? 'ok' : 'muted'"
                  :label="sysInfo.email?.configured ? $t('adminSettings.configured') : $t('adminSettings.notConfigured')"
                />
              </dd>
            </div>
            <div v-if="sysInfo.email?.configured" class="flex items-center justify-between gap-4 border-b border-y-line px-5 py-3">
              <dt class="text-sm text-y-muted">{{ $t('adminSettings.smtpServer') }}</dt>
              <dd class="break-all text-right font-mono text-sm text-y-text">{{ sysInfo.email?.host }}</dd>
            </div>
            <div v-if="sysInfo.email?.configured" class="flex items-center justify-between gap-4 px-5 py-3">
              <dt class="text-sm text-y-muted">{{ $t('adminSettings.port') }}</dt>
              <dd class="font-mono text-sm text-y-text">{{ sysInfo.email?.port }}</dd>
            </div>
          </dl>
        </div>
      </div>
    </section>

    <!-- Maintenance -->
    <section class="grid gap-4 border-t border-y-line pt-8 lg:grid-cols-[16rem_minmax(0,1fr)] lg:gap-10">
      <div class="space-y-1.5">
        <h2 class="y-h3">{{ $t('uiAdmin3.settings.maintenanceTitle') }}</h2>
        <p class="text-sm leading-6 text-y-muted">{{ $t('uiAdmin3.settings.maintenanceDesc') }}</p>
      </div>
      <div class="y-panel">
        <div class="flex flex-col gap-3 border-b border-y-line px-5 py-4 sm:flex-row sm:items-center sm:justify-between sm:gap-6">
          <div class="min-w-0 space-y-1">
            <p class="text-sm font-medium text-y-text">{{ $t('adminSettings.clearCache') }}</p>
            <p class="text-sm text-y-muted">{{ $t('uiAdmin3.settings.clearCacheHint') }}</p>
          </div>
          <button type="button" class="btn btn-outline btn-sm shrink-0 gap-1.5" :disabled="clearingCache" @click="clearCache">
            <span v-if="clearingCache" class="loading loading-spinner loading-xs"></span>
            <Eraser v-else :size="16" :stroke-width="1.5" />
            {{ $t('adminSettings.clearCache') }}
          </button>
        </div>
        <div class="flex flex-col gap-3 px-5 py-4 sm:flex-row sm:items-center sm:justify-between sm:gap-6">
          <div class="min-w-0 space-y-1">
            <p class="text-sm font-medium text-y-text">{{ $t('adminSettings.refreshSystemInfo') }}</p>
            <p class="text-sm text-y-muted">{{ $t('uiAdmin3.settings.systemDesc') }}</p>
          </div>
          <button type="button" class="btn btn-outline btn-sm shrink-0 gap-1.5" :disabled="sysLoading" @click="refreshSystemInfo">
            <RefreshCw :size="16" :stroke-width="1.5" />
            {{ $t('adminSettings.refreshSystemInfo') }}
          </button>
        </div>
      </div>
    </section>
  </div>
</template>
<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { RefreshCw, Check, Info, Eraser } from 'lucide-vue-next'
import axios from '../utils/axios'
import { useToast } from '../composables/useToast'
import PageHeader from '../components/ui/PageHeader.vue'
import StatusDot from '../components/ui/StatusDot.vue'
import YLoader from '../components/brand/YLoader.vue'

const { t } = useI18n()
const toast = useToast()
const quotaLoading = ref(true)
const sysLoading = ref(true)
const clearingCache = ref(false)
const sysInfo = ref({})

const siteSettings = reactive({
  allow_password_register: 'true',
})

const quotaLevels = reactive([
  { key: 'quota_normal', label: 'Normal', trustLevel: '0', value: 2, originalValue: 2, saving: false, saved: false, badgeClass: 'y-chip' },
  { key: 'quota_basic', label: 'Basic', trustLevel: '1', value: 3, originalValue: 3, saving: false, saved: false, badgeClass: 'y-chip' },
  { key: 'quota_member', label: 'Member', trustLevel: '2', value: 5, originalValue: 5, saving: false, saved: false, badgeClass: 'y-chip-info' },
  { key: 'quota_regular', label: 'Regular', trustLevel: '3', value: 10, originalValue: 10, saving: false, saved: false, badgeClass: 'y-chip-mint' },
  { key: 'quota_leader', label: 'Leader', trustLevel: '4', value: 20, originalValue: 20, saving: false, saved: false, badgeClass: 'y-chip-warn' },
])

const systemStats = computed(() => [
  { key: 'uptime', label: t('adminSettings.uptime'), value: sysInfo.value.uptime?.human || '-' },
  { key: 'users', label: t('adminSettings.totalUsers'), value: sysInfo.value.stats?.users ?? '-' },
  { key: 'domains', label: t('adminSettings.totalDomains'), value: sysInfo.value.stats?.domains ?? '-' },
  { key: 'orders', label: t('adminSettings.totalOrders'), value: sysInfo.value.stats?.orders ?? '-' },
  { key: 'alloc', label: t('adminSettings.memoryUsage'), value: sysInfo.value.memory?.alloc_mb || '-', unit: 'MB' },
  { key: 'sys', label: t('adminSettings.systemMemory'), value: sysInfo.value.memory?.sys_mb || '-', unit: 'MB' },
  { key: 'goroutines', label: 'Goroutines', value: sysInfo.value.system?.goroutines ?? '-' },
  { key: 'gc', label: t('adminSettings.gcCycles'), value: sysInfo.value.memory?.gc_cycles ?? '-' },
])

const currencySymbol = ref('NL')
const currencySaving = ref(false)
const currencySaved = ref(false)

const subdomainBlacklist = ref('')
const blacklistSaving = ref(false)
const blacklistSaved = ref(false)

const fetchSettings = async () => {
  quotaLoading.value = true
  try {
    const res = await axios.get('/api/admin/settings')
    const settings = res.data.settings || []
    for (const setting of settings) {
      const level = quotaLevels.find(l => l.key === setting.setting_key)
      if (level) {
        level.value = parseInt(setting.setting_value) || 0
        level.originalValue = level.value
      }
      if (setting.setting_key === 'allow_password_register') siteSettings.allow_password_register = setting.setting_value
      if (setting.setting_key === 'currency_symbol') currencySymbol.value = setting.setting_value || 'NL'
      if (setting.setting_key === 'subdomain_blacklist') subdomainBlacklist.value = setting.setting_value || ''
    }
  } catch (err) {
    console.error('Failed to fetch settings:', err)
  } finally {
    quotaLoading.value = false
  }
}

const togglePasswordRegister = async (e) => {
  const val = e.target.checked ? 'true' : 'false'
  siteSettings.allow_password_register = val
  try {
    await axios.put('/api/admin/settings/allow_password_register', { value: val })
    toast.success(val === 'true' ? t('adminSettings.passwordRegisterEnabled') : t('adminSettings.passwordRegisterDisabled'))
  } catch (err) {
    toast.error(t('adminSettings.updateFailed'))
    siteSettings.allow_password_register = val === 'true' ? 'false' : 'true'
  }
}

const fetchSystemInfo = async () => {
  sysLoading.value = true
  try {
    const res = await axios.get('/api/admin/system-info')
    sysInfo.value = res.data
  } catch (err) {
    console.error('Failed to fetch system info:', err)
  } finally {
    sysLoading.value = false
  }
}

const onQuotaFocus = (level) => {
  level.saved = false
}

const onQuotaBlur = async (level) => {
  if (level.value === level.originalValue) return
  level.saving = true
  level.saved = false
  try {
    await axios.put(`/api/admin/settings/${level.key}`, { value: String(level.value) })
    level.originalValue = level.value
    level.saving = false
    level.saved = true
    setTimeout(() => { level.saved = false }, 2000)
  } catch (err) {
    level.saving = false
    toast.error(t('adminSettings.updateFailed') + ': ' + (err.response?.data?.error || err.message))
  }
}

const updateCurrencySetting = async () => {
  if (!currencySymbol.value || currencySymbol.value.trim() === '') {
    currencySymbol.value = 'NL'
  }
  currencySaving.value = true
  currencySaved.value = false
  try {
    await axios.put('/api/admin/settings/currency_symbol', { value: currencySymbol.value.trim() })
    currencySaving.value = false
    currencySaved.value = true
    setTimeout(() => { currencySaved.value = false }, 2000)
  } catch (err) {
    currencySaving.value = false
    toast.error(t('adminSettings.updateFailed') + ': ' + (err.response?.data?.error || err.message))
  }
}

const updateBlacklistSetting = async () => {
  blacklistSaving.value = true
  blacklistSaved.value = false
  try {
    await axios.put('/api/admin/settings/subdomain_blacklist', { value: subdomainBlacklist.value.trim() })
    blacklistSaving.value = false
    blacklistSaved.value = true
    setTimeout(() => { blacklistSaved.value = false }, 2000)
  } catch (err) {
    blacklistSaving.value = false
    toast.error(t('adminSettings.saveFailed') + ': ' + (err.response?.data?.error || err.message))
  }
}

const clearCache = async () => {
  clearingCache.value = true
  try {
    await axios.post('/api/admin/clear-cache')
    toast.success(t('adminSettings.cacheCleared'))
  } catch (err) {
    toast.error(t('adminSettings.clearCacheFailed') + ': ' + (err.response?.data?.error || err.message))
  } finally {
    clearingCache.value = false
  }
}

const refreshSystemInfo = async () => {
  await fetchSystemInfo()
  toast.success(t('adminSettings.systemInfoRefreshed'))
}

onMounted(() => {
  fetchSettings()
  fetchSystemInfo()
})
</script>
