<template>
  <header class="sticky top-0 z-50 bg-bg/80 backdrop-blur-md border-b border-white/5">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
      <!-- Logo -->
      <RouterLink to="/" class="flex items-center gap-2">
        <svg width="28" height="28" viewBox="0 0 28 28" fill="none" xmlns="http://www.w3.org/2000/svg">
          <circle cx="14" cy="14" r="14" fill="#E8002D"/>
          <polygon points="11,8 22,14 11,20" fill="white"/>
          <rect x="5" y="11" width="2.5" height="6" rx="1" fill="white"/>
          <rect x="8.5" y="9" width="2.5" height="10" rx="1" fill="white"/>
        </svg>
        <span class="font-display text-2xl text-accent-red tracking-widest">LIVE</span>
        <span class="font-display text-2xl text-white tracking-widest">STREAMIFY</span>
      </RouterLink>

      <!-- Nav links -->
      <nav class="hidden md:flex items-center gap-6">
        <RouterLink to="/" class="text-text-muted hover:text-white transition-colors text-sm font-medium"
          active-class="text-white">
          Home
        </RouterLink>
        <RouterLink to="/promoters" class="text-text-muted hover:text-white transition-colors text-sm font-medium"
          active-class="text-white">
          For Promoters
        </RouterLink>
        <RouterLink :to="{ name: 'home', query: { sport: 'boxing' } }"
          class="text-text-muted hover:text-white transition-colors text-sm font-medium"
          :class="{ 'text-white': route.query.sport === 'boxing' }">
          Boxing
        </RouterLink>
        <RouterLink :to="{ name: 'home', query: { sport: 'racing' } }"
          class="text-text-muted hover:text-white transition-colors text-sm font-medium"
          :class="{ 'text-white': route.query.sport === 'racing' }">
          Racing
        </RouterLink>
      </nav>

      <!-- Right: auth actions -->
      <div class="flex items-center gap-3">
        <template v-if="auth.isAuthenticated">
          <RouterLink v-if="auth.isPromoter" to="/dashboard"
            class="text-text-muted hover:text-white text-sm font-medium transition-colors">
            Dashboard
          </RouterLink>
          <div class="relative group">
            <button class="flex items-center gap-2 text-sm font-medium">
              <div class="w-8 h-8 rounded-full bg-accent-red flex items-center justify-center text-white font-bold text-xs">
                {{ initials }}
              </div>
            </button>
            <!-- Dropdown -->
            <div class="absolute right-0 top-full mt-2 w-44 bg-bg-elevated border border-white/10 rounded-xl
                        shadow-xl opacity-0 invisible group-hover:opacity-100 group-hover:visible
                        transition-all duration-200">
              <div class="px-4 py-3 border-b border-white/5">
                <p class="text-white text-sm font-semibold truncate">{{ auth.user?.full_name }}</p>
                <p class="text-text-muted text-xs truncate">{{ auth.user?.email }}</p>
              </div>
              <button @click="handleLogout"
                class="w-full text-left px-4 py-2.5 text-sm text-text-muted hover:text-white
                       hover:bg-white/5 transition-colors rounded-b-xl">
                Sign out
              </button>
            </div>
          </div>
        </template>
        <template v-else>
          <RouterLink to="/login" class="btn-ghost text-sm py-2 px-4">Sign in</RouterLink>
          <RouterLink to="/register" class="btn-primary text-sm py-2 px-4">Get Started</RouterLink>
        </template>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink, useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const router = useRouter()
const route = useRoute()

const initials = computed(() => {
  const name = auth.user?.full_name ?? ''
  return name.split(' ').map((n) => n[0]).join('').slice(0, 2).toUpperCase()
})

function handleLogout() {
  auth.logout()
  router.push('/')
}
</script>
