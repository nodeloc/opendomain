<template>
  <section id="power-source" data-reveal class="y-container pb-24 pt-16 sm:pt-20" aria-labelledby="power-source-title">
    <div data-theme="yunamo" class="power-source-panel">
      <div class="power-visual">
        <div class="power-art" role="img" :aria-label="$t('home2.powerCrystalAlt')">
          <EnergyCrystal :active="motionEnabled" />
        </div>
        <button type="button" class="power-motion-toggle" :aria-pressed="motionPaused" :disabled="reducedMotion" @click="$emit('toggle-motion')">
          <Pause v-if="motionEnabled" :size="11" aria-hidden="true" /><Play v-else :size="11" aria-hidden="true" />
          {{ $t(reducedMotion ? 'home2.reducedMotion' : motionPaused ? 'home2.resumeMotion' : 'home2.pauseMotion') }}
        </button>
        <p class="power-caption"><span aria-hidden="true"></span>{{ $t('home2.powerCaption') }}</p>
      </div>
      <div class="power-copy">
        <p class="y-eyebrow">{{ $t('home2.powerEyebrow') }}</p>
        <h2 id="power-source-title" class="power-title">{{ $t('home2.powerTitleLead') }}<br /><em>{{ $t('home2.powerTitleAccent') }}</em></h2>
        <p class="power-description">{{ $t('home2.powerDescription') }}</p>
        <p class="power-detail">{{ $t('home2.powerDetail') }}</p>
        <div class="power-actions">
          <a href="#home-sub" class="power-claim btn btn-primary gap-2.5" @click.prevent="$emit('claim')">
            <Zap :size="16" :stroke-width="1.5" aria-hidden="true" />{{ $t('home2.powerCta') }}<ArrowUpRight :size="16" aria-hidden="true" />
          </a>
          <router-link :to="docsPath" class="power-docs btn btn-outline">{{ $t('home2.ctaDocs') }}</router-link>
        </div>
        <p class="power-signature">Yunamo. Your dynamo.</p>
      </div>
    </div>
  </section>
</template>

<script setup>
import { Pause, Play, Zap, ArrowUpRight } from 'lucide-vue-next'
import EnergyCrystal from './EnergyCrystal.vue'

defineProps({
  motionEnabled: { type: Boolean, default: true },
  motionPaused: { type: Boolean, default: false },
  reducedMotion: { type: Boolean, default: false },
  docsPath: { type: String, required: true },
})
defineEmits(['claim', 'toggle-motion'])
</script>

<style scoped>
.power-source-panel { position: relative; display: grid; grid-template-columns: 1fr 1fr; align-items: center; overflow: hidden; border: 1px solid #263b35; border-radius: 22px; background: radial-gradient(ellipse at 23% 44%, #10221c 0, #0d1512 38%, #0d1012 75%); color: #ededed; }
.power-visual { min-width: 0; position: relative; padding: 30px 0 48px; }
.power-art { position: relative; isolation: isolate; }
.power-motion-toggle { position: absolute; top: 16px; left: 20px; display: inline-flex; align-items: center; gap: 7px; min-height: 44px; padding: 0 12px; color: #9ab6ac; font-family: 'Geist Mono', monospace; font-size: 10px; border-radius: 999px; border: 1px solid rgb(136 197 177 / .15); background: rgb(10 17 15 / .3); transition: color .2s, border-color .2s; }
.power-motion-toggle:hover:not(:disabled) { color: #bfffe9; border-color: rgb(62 230 180 / .4); }
.power-motion-toggle:disabled { cursor: default; }
.power-caption { position: absolute; bottom: 39px; left: 16px; right: 16px; display: flex; align-items: center; justify-content: center; gap: 9px; font-family: 'Geist Mono', monospace; font-size: 9px; letter-spacing: .08em; text-transform: uppercase; color: #91b2a5; }
.power-caption > span { width: 4px; height: 4px; background: #75edc5; border-radius: 50%; }
.power-copy { display: flex; flex-direction: column; align-items: flex-start; gap: 22px; padding: 64px 52px 64px 20px; min-width: 0; }
.power-title { font-weight: 800; text-transform: uppercase; font-size: clamp(2rem, 3.3vw, 3.125rem); letter-spacing: -.045em; line-height: 1.04; }
.power-title em { font-style: normal; color: rgb(var(--y-mint)); }
.power-description { font-size: 16px; line-height: 1.8; color: #b4c4be; max-width: 420px; }
.power-detail { color: #8fa89d; font-size: 14px; line-height: 1.8; max-width: 395px; }
.power-actions { display: flex; flex-wrap: wrap; gap: 12px; margin-top: 5px; }
.power-claim { min-height: 46px; padding-inline: 20px; }
.power-docs { min-height: 46px; padding-inline: 18px; border-color: #385046; color: #c0d8ca; }
.power-docs:hover { border-color: rgb(var(--y-mint)); background: rgb(var(--y-mint) / .09); color: #e9fff2; }
.power-claim > svg:last-child { transition: transform .2s; }
.power-claim:hover > svg:last-child { transform: translate(2px, -2px); }
.power-signature { font-family: 'Geist Mono', monospace; color: #789589; font-size: 11px; }
@media (max-width: 767px) {
  .power-source-panel { grid-template-columns: 1fr; border-radius: 17px; }
  .power-visual { padding: 24px 12px 30px; max-width: 470px; width: 100%; margin-inline: auto; }
  .power-motion-toggle { left: 16px; top: 16px; }
  .power-caption { bottom: 28px; font-size: 8px; letter-spacing: .035em; }
  .power-copy { padding: 12px 26px 40px; gap: 20px; }
  .power-title { font-size: clamp(2rem, 7.8vw, 3.1rem); }
  .power-description { font-size: 15px; }
}
@media (prefers-reduced-motion: reduce) { .power-source-panel * { transition: none !important; } }
</style>
