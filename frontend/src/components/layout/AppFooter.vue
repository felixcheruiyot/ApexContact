<template>
  <footer class="border-t border-accent-red/30 bg-bg-surface mt-24">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-16 pb-10">

      <!-- Top grid: brand + 3 columns -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-10 mb-14">

        <!-- Brand -->
        <div class="lg:col-span-1">
          <RouterLink to="/" class="flex items-center gap-1.5 mb-4 group">
            <span class="font-display text-xl tracking-widest text-white">LIVE</span>
            <span class="w-2 h-2 rounded-full bg-accent-red shrink-0 mx-0.5" />
            <span class="font-display text-xl tracking-widest text-white">STREAMIFY</span>
          </RouterLink>
          <p class="text-text-muted text-sm leading-relaxed max-w-[220px] mb-5">
            Live streaming and monetisation for African knowledge creators.
          </p>
          <a
            href="mailto:hello@livestreamify.com"
            class="text-text-muted text-xs hover:text-white transition-colors"
          >
            hello@livestreamify.com
          </a>
        </div>

        <!-- Explore -->
        <div>
          <h4 class="text-white font-semibold text-sm tracking-wide mb-5 uppercase">Platform</h4>
          <ul class="space-y-3">
            <li v-for="link in exploreLinks" :key="link.label">
              <RouterLink
                :to="link.to"
                class="text-text-muted text-sm hover:text-white transition-colors duration-200"
              >
                {{ link.label }}
              </RouterLink>
            </li>
          </ul>
        </div>

        <!-- Creators -->
        <div>
          <h4 class="text-white font-semibold text-sm tracking-wide mb-5 uppercase">Creators</h4>
          <ul class="space-y-3">
            <li v-for="link in creatorLinks" :key="link.label">
              <RouterLink
                :to="link.to"
                class="text-text-muted text-sm hover:text-white transition-colors duration-200"
              >
                {{ link.label }}
              </RouterLink>
            </li>
          </ul>
        </div>

        <!-- Support -->
        <div>
          <h4 class="text-white font-semibold text-sm tracking-wide mb-5 uppercase">Support</h4>
          <ul class="space-y-3">
            <li v-for="link in supportLinks" :key="link.label">
              <a
                :href="link.to"
                class="text-text-muted text-sm hover:text-white transition-colors duration-200"
              >
                {{ link.label }}
              </a>
            </li>
          </ul>
        </div>
      </div>

      <!-- Bottom bar -->
      <div class="border-t border-white/5 pt-8 flex flex-col sm:flex-row items-center justify-between gap-4">
        <!-- Copyright -->
        <p class="text-text-muted text-xs">
          &copy; {{ new Date().getFullYear() }} Live Streamify. All rights reserved.
          Powered by M-Pesa &middot; Anti-piracy protected.
        </p>

        <a
          href="mailto:hello@livestreamify.com"
          class="text-text-muted text-xs hover:text-white transition-colors"
        >
          hello@livestreamify.com
        </a>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'

const exploreLinks = [
  { label: 'Home', to: '/' },
  { label: 'Use Cases', to: '/use-cases' },
  { label: 'For Creators', to: '/promoters' },
  { label: 'Try Free Stream', to: '/try' },
  { label: 'Privacy Policy', to: '/privacy' },
]

const creatorLinks = [
  { label: 'Start Hosting', to: '/promoters' },
  { label: 'Pricing & Payouts', to: '/promoters' },
  { label: 'Try a Test Stream', to: '/try' },
]

const supportLinks = [
  { label: 'Help Center', to: 'mailto:support@livestreamify.com' },
  { label: 'Contact', to: 'mailto:hello@livestreamify.com' },
]
</script>
