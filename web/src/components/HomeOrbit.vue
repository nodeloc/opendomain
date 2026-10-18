<template>
  <div class="home-orbit" :class="{ 'orbit-active': active }" aria-hidden="true">
    <svg class="orbit-map" viewBox="0 0 1440 840" preserveAspectRatio="xMidYMid slice" fill="none">
      <defs>
        <linearGradient id="home-orbit-line" x1="80" y1="610" x2="1380" y2="350" gradientUnits="userSpaceOnUse">
          <stop stop-color="#3ee6b4" stop-opacity="0" />
          <stop offset=".4" stop-color="#3ee6b4" stop-opacity=".3" />
          <stop offset=".8" stop-color="#7ab8ff" stop-opacity=".25" />
          <stop offset="1" stop-color="#3ee6b4" stop-opacity="0" />
        </linearGradient>
      </defs>
      <g stroke="url(#home-orbit-line)" stroke-width=".8">
        <path d="M-160 810C160 250 1010 210 1600 570" />
        <path d="M-120 500C310 910 1180 730 1520 230" />
        <path d="M200 920C370 210 1000 130 1330 650" />
      </g>
      <g stroke="#9af7dc" stroke-linecap="round">
        <path class="orbit-signal signal-one" pathLength="1000" d="M-160 810C160 250 1010 210 1600 570" />
        <path class="orbit-signal signal-two" pathLength="1000" d="M-120 500C310 910 1180 730 1520 230" />
        <path class="orbit-signal signal-three" pathLength="1000" d="M200 920C370 210 1000 130 1330 650" />
      </g>
      <g fill="#d7fff3">
        <circle v-for="star in stars" :key="star.x" class="orbit-star" :cx="star.x" :cy="star.y" :r="star.r" :style="{ '--star-delay': star.delay + 's' }" />
      </g>
    </svg>
    <div class="horizon-halo"></div>
  </div>
</template>

<script setup>
defineProps({ active: { type: Boolean, default: false } })

const stars = [
  { x: 95, y: 125, r: 1.2, delay: -1 }, { x: 345, y: 96, r: 1, delay: -3 },
  { x: 605, y: 188, r: 1.4, delay: -2 }, { x: 885, y: 88, r: 1.1, delay: -5 },
  { x: 1190, y: 170, r: 1.4, delay: -4 }, { x: 1330, y: 300, r: 1, delay: 0 },
  { x: 185, y: 355, r: 1.3, delay: -2.5 }, { x: 1065, y: 320, r: 1.1, delay: -1.5 },
]
</script>

<style scoped>
.home-orbit { position: absolute; inset: 0; overflow: hidden; pointer-events: none; }
.orbit-map { width: 100%; height: 100%; opacity: .7; }
.orbit-signal { stroke-width: 2; stroke-dasharray: 2 998; animation: orbit-travel 16s linear infinite; animation-play-state: paused; opacity: .85; }
.signal-two { animation-duration: 22s; animation-delay: -9s; }
.signal-three { animation-duration: 19s; animation-delay: -4s; }
.orbit-star { animation: star-breathe 6s ease-in-out infinite; animation-delay: var(--star-delay); animation-play-state: paused; opacity: .35; }
.horizon-halo { position: absolute; width: 38%; height: 12%; left: 60%; top: 58%; border-radius: 50%; background: radial-gradient(ellipse, rgb(62 230 180 / .12), transparent 70%); animation: horizon-breathe 8s ease-in-out infinite alternate; animation-play-state: paused; }
.orbit-active .orbit-signal,
.orbit-active .orbit-star,
.orbit-active .horizon-halo { animation-play-state: running; }
@keyframes orbit-travel { to { stroke-dashoffset: -1000; } }
@keyframes star-breathe { 50% { opacity: .9; } }
@keyframes horizon-breathe { to { opacity: .4; transform: scale(1.2); } }
@media (max-width: 640px) { .orbit-map { opacity: .4; } .signal-three { display: none; } }
@media (prefers-reduced-motion: reduce) { .orbit-signal, .orbit-star, .horizon-halo { animation: none !important; } }
</style>
