<template>
  <svg
    v-if="variant === 'mark'"
    viewBox="0 0 64 64"
    xmlns="http://www.w3.org/2000/svg"
    role="img"
    aria-label="Yunamo"
    class="y-logo"
  >
    <path
      d="M32 54V35M32 35C32 26.5 24.5 22 16.5 14.5M32 35C32 26.5 39.5 22 47.5 14.5"
      stroke="currentColor"
      stroke-width="6.5"
      stroke-linecap="round"
      stroke-linejoin="round"
      fill="none"
    />
    <circle cx="15" cy="13" r="7" class="y-logo-node" />
    <circle cx="49" cy="13" r="7" class="y-logo-node" />
  </svg>
  <svg
    v-else
    viewBox="0 0 313 80"
    xmlns="http://www.w3.org/2000/svg"
    role="img"
    aria-label="Yunamo"
    class="y-logo"
  >
    <path
      d="M32 54V35M32 35C32 26.5 24.5 22 16.5 14.5M32 35C32 26.5 39.5 22 47.5 14.5"
      stroke="currentColor"
      stroke-width="6.5"
      stroke-linecap="round"
      stroke-linejoin="round"
      fill="none"
    />
    <circle cx="15" cy="13" r="7" class="y-logo-node" />
    <circle cx="49" cy="13" r="7" class="y-logo-node" />
    <g
      transform="translate(78 11) scale(1.15)"
      fill="none"
      stroke="currentColor"
      stroke-width="6"
      stroke-linecap="round"
      stroke-linejoin="round"
    >
      <path transform="translate(3 0)" d="M0 26A11 11 0 0 0 22 26M22 15V44A11 11 0 0 1 11 55H5M0 15V26" />
      <path transform="translate(36 0)" d="M0 15V26A11 11 0 0 0 22 26V15M22 26V37" />
      <path transform="translate(69 0)" d="M0 37V26A11 11 0 0 1 22 26V37M0 26V15" />
      <path transform="translate(102 0)" d="M22 15V37M11 15A11 11 0 1 0 11.01 15" />
      <path transform="translate(135 0)" d="M0 37V23A8 8 0 0 1 16 23V37M16 23A8 8 0 0 1 32 23V37M0 23V15" />
      <path transform="translate(178 0)" d="M11 15A11 11 0 1 0 11.01 15" />
    </g>
  </svg>
</template>

<script setup>
defineProps({
  variant: {
    type: String,
    default: 'lockup',
    validator: (v) => ['lockup', 'mark'].includes(v),
  },
})
</script>

<style scoped>
.y-logo {
  display: block;
  color: rgb(var(--y-text));
}
.y-logo-node {
  fill: rgb(var(--y-node));
}
</style>
