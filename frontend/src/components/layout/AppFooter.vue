<template>
  <footer class="bg-bg-surface border-t border-white/5 mt-16">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div class="grid grid-cols-1 md:grid-cols-3 gap-8 mb-8">
        <div>
          <div class="flex items-center gap-2 mb-4">
            <svg width="24" height="24" viewBox="0 0 28 28" fill="none" xmlns="http://www.w3.org/2000/svg">
              <circle cx="14" cy="14" r="14" fill="#E8002D"/>
              <polygon points="11,8 22,14 11,20" fill="white"/>
              <rect x="5" y="11" width="2.5" height="6" rx="1" fill="white"/>
              <rect x="8.5" y="9" width="2.5" height="10" rx="1" fill="white"/>
            </svg>
            <span class="font-display text-xl text-accent-red tracking-widest">LIVE</span>
            <span class="font-display text-xl text-white tracking-widest">STREAMIFY</span>
          </div>
          <p class="text-text-muted text-sm leading-relaxed max-w-xs">
            Premium live sports streaming. Boxing and car racing events — live, secure, unstoppable.
          </p>
        </div>

        <div>
          <h4 class="text-white font-semibold text-sm mb-4">Sports</h4>
          <ul class="space-y-2">
            <li><RouterLink :to="{ name: 'home', query: { sport: 'boxing' } }" class="text-text-muted text-sm hover:text-white transition-colors">🥊 Boxing</RouterLink></li>
            <li><RouterLink :to="{ name: 'home', query: { sport: 'racing' } }" class="text-text-muted text-sm hover:text-white transition-colors">🏎️ Car Racing</RouterLink></li>
            <li><RouterLink to="/" class="text-text-muted text-sm hover:text-accent-orange transition-colors">🎙 Commentary Lobbies</RouterLink></li>
          </ul>
        </div>

        <div>
          <h4 class="text-white font-semibold text-sm mb-4">Platform</h4>
          <ul class="space-y-2">
            <li><RouterLink to="/commentary/create" class="text-text-muted text-sm hover:text-accent-orange transition-colors">Host a Commentary</RouterLink></li>
            <li><RouterLink to="/promoters" class="text-text-muted text-sm hover:text-white transition-colors">For Promoters</RouterLink></li>
            <li><a href="mailto:support@livestreamify.com" class="text-text-muted text-sm hover:text-white transition-colors">Help Center</a></li>
            <li><RouterLink to="/privacy" class="text-text-muted text-sm hover:text-white transition-colors">Privacy Policy</RouterLink></li>
          </ul>
        </div>
      </div>

      <div class="border-t border-white/5 pt-8 flex flex-col md:flex-row items-center justify-between gap-4">
        <p class="text-text-muted text-xs">
          &copy; {{ new Date().getFullYear() }} Live Streamify. All rights reserved.
        </p>
        <p class="text-text-muted text-xs">
          Powered by M-Pesa · Secure streaming · Anti-piracy protected
        </p>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'
</script>
