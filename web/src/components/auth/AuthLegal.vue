<template>
  <i18n-t keypath="auth2.agree" tag="p" scope="global" class="max-w-[360px] text-[13px] leading-5 text-y-faint">
    <template #policy>
      <router-link to="/acceptable-use" class="text-y-mint-text underline decoration-y-mint-text/40 underline-offset-[3px] hover:decoration-y-mint-text">{{ $t('auth2.policy') }}</router-link>
    </template>
    <template #terms>
      <router-link to="/pages/terms-of-service" class="text-y-mint-text underline decoration-y-mint-text/40 underline-offset-[3px] hover:decoration-y-mint-text">{{ $t('auth2.terms') }}</router-link>
    </template>
  </i18n-t>
</template>
