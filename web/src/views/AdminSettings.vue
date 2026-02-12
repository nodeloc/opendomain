<template>
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold mb-6">{{ $t('adminSettings.title') }}</h1>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Site Settings -->
      <div class="card bg-base-100 shadow-xl lg:col-span-2">
        <div class="card-body">
          <h2 class="card-title">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9" />
            </svg>
            {{ $t('adminSettings.siteSettings') }}
          </h2>
          <div class="space-y-3 mt-2">
            <div class="flex justify-between">
              <span class="opacity-70">{{ $t('adminSettings.siteName') }}</span>
              <span class="font-bold">{{ sysInfo.site?.name || '-' }}</span>
            </div>
            <div class="flex justify-between">
              <span class="opacity-70">{{ $t('adminSettings.siteDescription') }}</span>
              <span>{{ sysInfo.site?.description || '-' }}</span>
            </div>
            <div class="text-xs opacity-50">{{ $t('adminSettings.siteConfigEnvHint') }}</div>
          </div>
          <div class="divider my-2"></div>
          <div class="form-control">
            <label class="label cursor-pointer justify-start gap-3">
              <input
                type="checkbox"
                class="toggle toggle-primary"
                :checked="siteSettings.allow_password_register !== 'false'"
                @change="togglePasswordRegister"
              />
              <div>
                <span class="label-text font-semibold">{{ $t('adminSettings.allowPasswordRegister') }}</span>
                <div class="text-xs opacity-60">{{ $t('adminSettings.allowPasswordRegisterDesc') }}</div>
              </div>
            </label>
          </div>
        </div>
      </div>

      <!-- User Level Quota Settings -->
      <div class="card bg-base-100 shadow-xl lg:col-span-2">
        <div class="card-body">
          <h2 class="card-title">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
            {{ $t('adminSettings.quotaSettings') }}
          </h2>
          <div class="text-sm opacity-70 mb-4">{{ $t('adminSettings.quotaSettingsDesc') }}</div>

          <div v-if="quotaLoading" class="flex justify-center py-8">
            <span class="loading loading-spinner loading-lg"></span>
          </div>

          <div v-else class="overflow-x-auto">
            <table class="table">
              <thead>
                <tr>
                  <th>{{ $t('adminSettings.level') }}</th>
                  <th>{{ $t('adminSettings.trustLevel') }}</th>
                  <th>{{ $t('adminSettings.domainQuota') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="level in quotaLevels" :key="level.key">
                  <td>
                    <span class="px-2.5 py-0.5 rounded-md text-xs font-medium" :class="level.badgeClass">{{ level.label }}</span>
                  </td>
                  <td class="font-mono">{{ level.trustLevel }}</td>
                  <td>
                    <div class="inline-flex items-center gap-2">
                      <input
                        v-model.number="level.value"
                        type="number"
                        min="0"
                        class="input input-bordered input-sm w-24"
                        @focus="onQuotaFocus(level)"
                        @blur="onQuotaBlur(level)"
                      />
                      <span v-if="level.saving" class="loading loading-spinner loading-xs opacity-60"></span>
                      <svg v-else-if="level.saved" xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-green-600" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" /></svg>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

        </div>
      </div>

      <!-- Currency Settings -->
      <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <h2 class="card-title">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            {{ $t('adminSettings.currencySettings') }}
          </h2>
          <div class="text-sm opacity-70 mb-4">{{ $t('adminSettings.currencySymbolDesc') }}</div>

          <div class="form-control">
            <label class="label">
              <span class="label-text font-semibold">{{ $t('adminSettings.currencySymbol') }}</span>
            </label>
            <input
              v-model="currencySymbol"
              type="text"
              maxlength="10"
              class="input input-bordered w-full max-w-xs"
              placeholder="NL"
              @blur="updateCurrencySetting"
            />
            <label class="label">
              <span class="label-text-alt">{{ $t('adminSettings.currencyExample', { symbol: currencySymbol || 'NL' }) }}</span>
            </label>
            <div v-if="currencySaving" class="text-sm text-info mt-2">
              <span class="loading loading-spinner loading-xs mr-1"></span>
              {{ $t('adminSettings.saving') }}
            </div>
            <div v-if="currencySaved" class="text-sm text-success mt-2">
              ✓ {{ $t('adminSettings.saved') }}
            </div>
          </div>

        </div>
      </div>

      <!-- System Info -->
      <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <h2 class="card-title">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            {{ $t('adminSettings.systemInfo') }}
          </h2>
          <div v-if="sysLoading" class="flex justify-center py-4">
            <span class="loading loading-spinner"></span>
          </div>
          <div v-else class="space-y-3">
            <div class="flex justify-between">
              <span class="opacity-70">{{ $t('adminSettings.systemName') }}</span>
              <span class="font-bold">{{ sysInfo.site?.name || '-' }}</span>
            </div>
            <div class="flex justify-between">
              <span class="opacity-70">{{ $t('adminSettings.version') }}</span>
              <span class="font-mono">{{ sysInfo.system?.version || '-' }}</span>
            </div>
            <div class="flex justify-between items-center">
              <span class="opacity-70">{{ $t('adminSettings.environment') }}</span>
              <span class="px-2.5 py-1 rounded-md text-xs font-medium inline-flex items-center gap-1.5" :class="sysInfo.system?.environment === 'production' ? 'bg-green-100 text-green-800' : 'bg-amber-100 text-amber-800'">
                <svg v-if="sysInfo.system?.environment === 'production'" xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M2.166 4.999A11.954 11.954 0 0010 1.944 11.954 11.954 0 0017.834 5c.11.65.166 1.32.166 2.001 0 5.225-3.34 9.67-8 11.317C5.34 16.67 2 12.225 2 7c0-.682.057-1.35.166-2.001zm11.541 3.708a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M11.3 1.046A1 1 0 0112 2v5h4a1 1 0 01.82 1.573l-7 10A1 1 0 018 18v-5H4a1 1 0 01-.82-1.573l7-10a1 1 0 011.12-.38z" clip-rule="evenodd" />
                </svg>
                {{ sysInfo.system?.environment || '-' }}
              </span>
            </div>
            <div class="flex justify-between">
              <span class="opacity-70">{{ $t('adminSettings.goVersion') }}</span>
              <span class="font-mono text-sm">{{ sysInfo.system?.go_version || '-' }}</span>
            </div>
            <div class="flex justify-between items-center">
              <span class="opacity-70">{{ $t('adminSettings.database') }}</span>
              <span class="px-2.5 py-1 rounded-md text-xs font-medium inline-flex items-center gap-1.5" :class="sysInfo.services?.database === 'connected' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'">
                <svg v-if="sysInfo.services?.database === 'connected'" xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
                </svg>
                PostgreSQL
              </span>
            </div>
            <div class="flex justify-between items-center">
              <span class="opacity-70">{{ $t('adminSettings.cache') }}</span>
              <span class="px-2.5 py-1 rounded-md text-xs font-medium inline-flex items-center gap-1.5" :class="sysInfo.services?.redis === 'connected' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'">
                <svg v-if="sysInfo.services?.redis === 'connected'" xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
                </svg>
                Redis
              </span>
            </div>
            <div class="flex justify-between">
              <span class="opacity-70">{{ $t('adminSettings.platform') }}</span>
              <span class="font-mono text-sm">{{ sysInfo.system?.os }}/{{ sysInfo.system?.arch }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- System Statistics -->
      <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <h2 class="card-title">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
            {{ $t('adminSettings.systemStats') }}
          </h2>
          <div v-if="sysLoading" class="flex justify-center py-4">
            <span class="loading loading-spinner"></span>
          </div>
          <div v-else class="space-y-3">
            <div class="flex justify-between">
              <span class="opacity-70">{{ $t('adminSettings.uptime') }}</span>
              <span class="font-bold">{{ sysInfo.uptime?.human || '-' }}</span>
            </div>
            <div class="flex justify-between">
              <span class="opacity-70">{{ $t('adminSettings.totalUsers') }}</span>
              <span class="font-bold">{{ sysInfo.stats?.users ?? '-' }}</span>
            </div>
            <div class="flex justify-between">
              <span class="opacity-70">{{ $t('adminSettings.totalDomains') }}</span>
              <span class="font-bold">{{ sysInfo.stats?.domains ?? '-' }}</span>
            </div>
            <div class="flex justify-between">
              <span class="opacity-70">{{ $t('adminSettings.totalOrders') }}</span>
              <span class="font-bold">{{ sysInfo.stats?.orders ?? '-' }}</span>
            </div>
            <div class="flex justify-between">
              <span class="opacity-70">{{ $t('adminSettings.memoryUsage') }}</span>
              <span class="font-bold">{{ sysInfo.memory?.alloc_mb || '-' }} MB</span>
            </div>
            <div class="flex justify-between">
              <span class="opacity-70">{{ $t('adminSettings.systemMemory') }}</span>
              <span class="font-bold">{{ sysInfo.memory?.sys_mb || '-' }} MB</span>
            </div>
            <div class="flex justify-between">
              <span class="opacity-70">Goroutines</span>
              <span class="font-mono">{{ sysInfo.system?.goroutines ?? '-' }}</span>
            </div>
            <div class="flex justify-between">
              <span class="opacity-70">{{ $t('adminSettings.gcCycles') }}</span>
              <span class="font-mono">{{ sysInfo.memory?.gc_cycles ?? '-' }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Default Nameservers -->
      <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <h2 class="card-title">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01" />
            </svg>
            {{ $t('adminSettings.defaultDNS') }}
          </h2>
          <div class="space-y-2">
            <div class="text-sm opacity-70">{{ $t('adminSettings.defaultDNSDesc') }}</div>
            <div v-if="sysLoading" class="flex justify-center py-4">
              <span class="loading loading-spinner"></span>
            </div>
            <div v-else class="bg-base-200 p-4 rounded font-mono text-sm space-y-2">
              <div class="flex items-center gap-2">
                <span class="px-2 py-0.5 rounded text-xs font-medium bg-blue-100 text-blue-800">NS1</span>
                <span>{{ sysInfo.dns?.ns1 || '-' }}</span>
              </div>
              <div class="flex items-center gap-2">
                <span class="px-2 py-0.5 rounded text-xs font-medium bg-purple-100 text-purple-800">NS2</span>
                <span>{{ sysInfo.dns?.ns2 || '-' }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- OAuth Settings -->
      <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <h2 class="card-title">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
            </svg>
            {{ $t('adminSettings.oauthSettings') }}
          </h2>
          <div v-if="sysLoading" class="flex justify-center py-4">
            <span class="loading loading-spinner"></span>
          </div>
          <div v-else class="space-y-3">
            <div class="flex justify-between items-center">
              <div class="flex items-center gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="opacity-70">
                  <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
                </svg>
                <span class="opacity-70">GitHub</span>
              </div>
              <span class="px-2.5 py-1 rounded-md text-xs font-medium inline-flex items-center gap-1.5" :class="sysInfo.oauth?.github ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-500'">
                <svg v-if="sysInfo.oauth?.github" xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M3 10a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z" clip-rule="evenodd" />
                </svg>
                {{ sysInfo.oauth?.github ? $t('adminSettings.configured') : $t('adminSettings.notConfigured') }}
              </span>
            </div>
            <div class="flex justify-between items-center">
              <div class="flex items-center gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" class="opacity-70">
                  <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92a5.06 5.06 0 0 1-2.2 3.32v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.1z"/>
                  <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
                  <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
                  <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
                </svg>
                <span class="opacity-70">Google</span>
              </div>
              <span class="px-2.5 py-1 rounded-md text-xs font-medium inline-flex items-center gap-1.5" :class="sysInfo.oauth?.google ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-500'">
                <svg v-if="sysInfo.oauth?.google" xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M3 10a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z" clip-rule="evenodd" />
                </svg>
                {{ sysInfo.oauth?.google ? $t('adminSettings.configured') : $t('adminSettings.notConfigured') }}
              </span>
            </div>
            <div class="flex justify-between items-center">
              <div class="flex items-center gap-2">
                <svg viewBox="0 0 102 64" xmlns="http://www.w3.org/2000/svg" width="16" height="10" class="opacity-70">
                  <path d="M51.504 19.32C51.504 23.224 51.312 28.888 50.928 36.312C50.544 43.608 50.352 49.144 50.352 52.92C50.352 54.008 50.416 55.64 50.544 57.816C50.672 59.992 50.736 61.688 50.736 62.904C50.608 63.288 49.36 63.64 46.992 63.96C44.624 64.28 42.992 64.44 42.096 64.44C41.648 64.44 40.944 63.832 39.984 62.616C39.024 61.4 37.808 59.704 36.336 57.528C33.904 54.072 32.368 52.088 31.728 51.576C31.728 51.576 31.216 50.808 30.192 49.272C29.104 47.544 27.76 45.496 26.16 43.128C24.624 40.76 23.856 39.544 23.856 39.48L15.504 29.592C15.44 29.464 15.28 29.112 15.024 28.536C14.768 27.96 14.416 27.672 13.968 27.672C13.776 27.672 13.584 27.864 13.392 28.248C13.264 28.568 13.2 28.92 13.2 29.304C13.2 36.024 13.264 41.304 13.392 45.144C13.456 46.04 13.488 47.16 13.488 48.504C13.488 48.696 13.456 48.92 13.392 49.176C13.328 49.368 13.264 49.528 13.2 49.656V64.44L12.432 65.208L10.32 65.304L0.528 64.536C0.528 62.744 0.592 60.632 0.72 58.2C0.912 55.768 1.008 54.328 1.008 53.88C1.008 52.728 0.944 51.128 0.816 49.08C0.688 46.776 0.624001 44.792 0.624001 43.128L0.816 39.864C1.008 36.216 1.104 34.2 1.104 33.816C1.104 32.28 1.04 29.976 0.912 26.904C0.784 23.832 0.72 21.56 0.72 20.088C0.72 19.192 0.848 17.912 1.104 16.248C1.36 14.712 1.488 13.592 1.488 12.888C1.488 12.696 1.36 11.736 1.104 10.008C0.848 8.216 0.72 6.68 0.72 5.4L0.336 4.72799C0.464 4.088 0.688 3.672 1.008 3.48C1.456 3.416 2.672 3.224 4.656 2.904C6.64 2.52 8.24 2.328 9.456 2.328C10.544 2.328 11.28 2.456 11.664 2.712C12.112 3.096 12.816 3.768 13.776 4.72799C14.736 5.688 15.184 6.264 15.12 6.456L20.016 14.808C20.016 14.872 20.912 16.152 22.704 18.648L24.24 20.856C27.44 25.592 30.256 29.656 32.688 33.048C35.12 36.376 36.688 38.04 37.392 38.04C37.84 38.04 38.096 37.88 38.16 37.56C38.16 36.536 38.096 35.128 37.968 33.336L37.872 29.592C37.872 28.44 37.936 26.648 38.064 24.216L38.16 18.84C38.16 17.368 38.096 15.16 37.968 12.216C37.84 9.4 37.776 7.224 37.776 5.688C37.776 4.152 37.872 3.064 38.064 2.424C38.32 1.784 38.8 1.336 39.504 1.08C40.208 0.759996 41.36 0.599997 42.96 0.599997C46.16 0.599997 48.144 0.759996 48.912 1.08C49.68 1.336 50.096 1.816 50.16 2.51999C50.224 3.16 50.416 3.704 50.736 4.152V10.2L51.024 11.064L50.928 14.328C50.928 16.76 51.12 18.424 51.504 19.32ZM100.475 49.368C100.795 49.368 101.083 49.464 101.339 49.656C101.659 49.784 101.819 50.008 101.819 50.328V60.6C101.819 61.24 101.595 61.784 101.147 62.232C100.699 62.616 100.251 62.808 99.8025 62.808L94.7145 62.616C92.0905 62.616 90.2025 62.712 89.0505 62.904L86.7465 63C86.7465 62.872 85.8185 62.744 83.9625 62.616L78.9705 62.232C77.6265 62.232 76.2825 62.424 74.9385 62.808C73.5945 63.192 72.9225 63.416 72.9225 63.48C71.8345 63.48 70.5865 63.32 69.1785 63C67.5145 62.744 66.2985 62.616 65.5305 62.616C65.5305 61.4 65.4025 59.512 65.1465 56.952C64.8905 54.264 64.7625 52.216 64.7625 50.808C64.7625 50.808 64.7625 50.68 64.7625 50.424C64.8265 50.168 64.8585 49.528 64.8585 48.504C64.8585 46.968 64.8265 45.528 64.7625 44.184L64.2825 43.704V27.96C63.9625 26.04 63.8025 24.344 63.8025 22.872C63.8025 21.464 63.9625 20.568 64.2825 20.184V11.256C64.0905 11.128 63.9945 10.584 63.9945 9.62399V8.184C63.9945 6.2 63.8345 4.536 63.5145 3.192L63.2265 2.808C63.2265 2.168 64.2185 1.656 66.2025 1.272C68.1865 0.823997 70.9065 0.599997 74.3625 0.599997C74.8745 0.599997 75.2905 0.759996 75.6105 1.08C75.9305 1.336 76.0905 1.656 76.0905 2.04C76.0905 2.744 75.9945 3.576 75.8025 4.536C75.3545 6.904 75.1305 9.144 75.1305 11.256C75.1305 11.32 75.0665 12.44 74.9385 14.616C74.8105 16.728 74.6185 18.584 74.3625 20.184L74.9385 20.664V32.184C75.3225 32.888 75.6425 34.488 75.8985 36.984C76.2185 39.416 76.4105 40.92 76.4745 41.496C76.7945 44.568 77.0825 46.936 77.3385 48.6C77.5945 50.2 77.9145 51.192 78.2985 51.576C78.5545 51.704 78.8745 51.832 79.2585 51.96C79.6425 52.088 79.9625 52.184 80.2185 52.248L83.6745 52.152C86.3625 51.96 90.0745 51.864 94.8105 51.864C95.0665 51.8 95.3865 51.736 95.7705 51.672C96.2185 51.608 96.5705 51.576 96.8265 51.576C98.0425 51.512 98.9385 51.352 99.5145 51.096C100.091 50.84 100.411 50.264 100.475 49.368Z" fill="url(#paint0_settings)"/>
                  <defs>
                    <linearGradient id="paint0_settings" x1="-6" y1="35" x2="104" y2="35" gradientUnits="userSpaceOnUse">
                      <stop stop-color="#009966"/>
                      <stop offset="0.790909" stop-color="#FF9933"/>
                    </linearGradient>
                  </defs>
                </svg>
                <span class="opacity-70">NodeLoc</span>
              </div>
              <span class="px-2.5 py-1 rounded-md text-xs font-medium inline-flex items-center gap-1.5" :class="sysInfo.oauth?.nodeloc ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-500'">
                <svg v-if="sysInfo.oauth?.nodeloc" xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M3 10a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z" clip-rule="evenodd" />
                </svg>
                {{ sysInfo.oauth?.nodeloc ? $t('adminSettings.configured') : $t('adminSettings.notConfigured') }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- Payment Settings -->
      <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <h2 class="card-title">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z" />
            </svg>
            {{ $t('adminSettings.paymentSettings') }}
          </h2>
          <div v-if="sysLoading" class="flex justify-center py-4">
            <span class="loading loading-spinner"></span>
          </div>
          <div v-else class="space-y-3">
            <div class="flex justify-between items-center">
              <span class="opacity-70">NodeLoc Pay</span>
              <span class="px-2.5 py-1 rounded-md text-xs font-medium inline-flex items-center gap-1.5" :class="sysInfo.payment?.configured ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-500'">
                <svg v-if="sysInfo.payment?.configured" xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M3 10a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z" clip-rule="evenodd" />
                </svg>
                {{ sysInfo.payment?.configured ? $t('adminSettings.configured') : $t('adminSettings.notConfigured') }}
              </span>
            </div>
            <div v-if="sysInfo.payment?.configured" class="flex justify-between items-center">
              <span class="opacity-70">{{ $t('adminSettings.testMode') }}</span>
              <span class="px-2.5 py-1 rounded-md text-xs font-medium inline-flex items-center gap-1.5" :class="sysInfo.payment?.test_mode ? 'bg-amber-100 text-amber-800' : 'bg-blue-100 text-blue-800'">
                <svg v-if="sysInfo.payment?.test_mode" xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M2.166 4.999A11.954 11.954 0 0010 1.944 11.954 11.954 0 0017.834 5c.11.65.166 1.32.166 2.001 0 5.225-3.34 9.67-8 11.317C5.34 16.67 2 12.225 2 7c0-.682.057-1.35.166-2.001zm11.541 3.708a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                </svg>
                {{ sysInfo.payment?.test_mode ? $t('common.yes') : $t('common.no') }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- Email Settings -->
      <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <h2 class="card-title">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
            </svg>
            {{ $t('adminSettings.emailSettings') }}
          </h2>
          <div v-if="sysLoading" class="flex justify-center py-4">
            <span class="loading loading-spinner"></span>
          </div>
          <div v-else class="space-y-3">
            <div class="flex justify-between items-center">
              <span class="opacity-70">{{ $t('adminSettings.status') }}</span>
              <span class="px-2.5 py-1 rounded-md text-xs font-medium inline-flex items-center gap-1.5" :class="sysInfo.email?.configured ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-500'">
                <svg v-if="sysInfo.email?.configured" xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M3 10a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z" clip-rule="evenodd" />
                </svg>
                {{ sysInfo.email?.configured ? $t('adminSettings.configured') : $t('adminSettings.notConfigured') }}
              </span>
            </div>
            <div v-if="sysInfo.email?.configured" class="flex justify-between">
              <span class="opacity-70">{{ $t('adminSettings.smtpServer') }}</span>
              <span class="font-mono text-sm">{{ sysInfo.email?.host }}</span>
            </div>
            <div v-if="sysInfo.email?.configured" class="flex justify-between">
              <span class="opacity-70">{{ $t('adminSettings.port') }}</span>
              <span>{{ sysInfo.email?.port }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Maintenance -->
      <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <h2 class="card-title">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
            {{ $t('adminSettings.maintenance') }}
          </h2>
          <div class="space-y-3">
            <button class="btn btn-outline btn-block" :disabled="clearingCache" @click="clearCache">
              <span v-if="clearingCache" class="loading loading-spinner loading-sm"></span>
              <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              {{ $t('adminSettings.clearCache') }}
            </button>
            <button class="btn btn-outline btn-block" @click="refreshSystemInfo">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              {{ $t('adminSettings.refreshSystemInfo') }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import axios from '../utils/axios'
import { useToast } from '../composables/useToast'

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
  { key: 'quota_normal', label: 'Normal', trustLevel: '0', value: 2, originalValue: 2, saving: false, saved: false, badgeClass: 'bg-gray-100 text-gray-700' },
  { key: 'quota_basic', label: 'Basic', trustLevel: '1', value: 3, originalValue: 3, saving: false, saved: false, badgeClass: 'bg-blue-100 text-blue-800' },
  { key: 'quota_member', label: 'Member', trustLevel: '2', value: 5, originalValue: 5, saving: false, saved: false, badgeClass: 'bg-green-100 text-green-800' },
  { key: 'quota_regular', label: 'Regular', trustLevel: '3', value: 10, originalValue: 10, saving: false, saved: false, badgeClass: 'bg-amber-100 text-amber-800' },
  { key: 'quota_leader', label: 'Leader', trustLevel: '4', value: 20, originalValue: 20, saving: false, saved: false, badgeClass: 'bg-red-100 text-red-800' },
])

const currencySymbol = ref('NL')
const currencySaving = ref(false)
const currencySaved = ref(false)

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
