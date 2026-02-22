<template>
  <footer class="bg-bg-surface border-t border-white/5 mt-16">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8 mb-8">
        <!-- Brand -->
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
            The open platform for live knowledge and event monetisation. Host paid or free events — sports, mentoring, education, and more.
          </p>
          <p class="text-text-muted text-xs mt-4">
            <a href="mailto:hello@livestreamify.com" class="hover:text-white transition-colors">hello@livestreamify.com</a>
          </p>
        </div>

        <!-- Explore -->
        <div>
          <h4 class="text-white font-semibold text-sm mb-4">Explore</h4>
          <ul class="space-y-2">
            <li><RouterLink to="/" class="text-text-muted text-sm hover:text-white transition-colors">🔴 Live Now</RouterLink></li>
            <li><RouterLink to="/" class="text-text-muted text-sm hover:text-white transition-colors">📅 Upcoming Events</RouterLink></li>
            <li><RouterLink to="/" class="text-text-muted text-sm hover:text-accent-orange transition-colors">🎙 Live Rooms</RouterLink></li>
            <li><RouterLink to="/use-cases" class="text-text-muted text-sm hover:text-white transition-colors">💡 Use Cases</RouterLink></li>
          </ul>
        </div>

        <!-- Use Cases -->
        <div>
          <h4 class="text-white font-semibold text-sm mb-4">Use Cases</h4>
          <ul class="space-y-2">
            <li><RouterLink to="/use-cases" class="text-text-muted text-sm hover:text-white transition-colors">🥊 Sports & Combat</RouterLink></li>
            <li><RouterLink to="/use-cases" class="text-text-muted text-sm hover:text-white transition-colors">🏫 Mentoring & Coaching</RouterLink></li>
            <li><RouterLink to="/use-cases" class="text-text-muted text-sm hover:text-white transition-colors">🌍 Visa & Immigration</RouterLink></li>
            <li><RouterLink to="/use-cases" class="text-text-muted text-sm hover:text-white transition-colors">📚 Education</RouterLink></li>
            <li><RouterLink to="/use-cases" class="text-text-muted text-sm hover:text-white transition-colors">🎵 Music & Performances</RouterLink></li>
            <li><RouterLink to="/use-cases" class="text-text-muted text-sm hover:text-white transition-colors">⚖️ Legal Consultations</RouterLink></li>
          </ul>
        </div>

        <!-- Platform -->
        <div>
          <h4 class="text-white font-semibold text-sm mb-4">Platform</h4>
          <ul class="space-y-2">
            <li><RouterLink to="/promoters" class="text-text-muted text-sm hover:text-white transition-colors">Host an Event</RouterLink></li>
            <li><RouterLink to="/commentary/create" class="text-text-muted text-sm hover:text-accent-orange transition-colors">Start a Live Room</RouterLink></li>
            <li><RouterLink to="/promoters" class="text-text-muted text-sm hover:text-white transition-colors">Pricing</RouterLink></li>
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
