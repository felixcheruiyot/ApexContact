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
          <!-- Admin panel link -->
          <RouterLink v-if="auth.isAdmin" to="/admin"
            class="text-accent-red hover:opacity-80 text-sm font-semibold transition-opacity flex items-center gap-1.5">
            <span class="w-1.5 h-1.5 rounded-full bg-accent-red animate-pulse" />
            Admin Panel
          </RouterLink>
          <!-- Promoter/broadcaster dashboard link -->
          <RouterLink v-else-if="auth.isPromoter" to="/dashboard"
            class="text-text-muted hover:text-white text-sm font-medium transition-colors">
            Dashboard
          </RouterLink>

          <!-- User avatar dropdown -->
          <div class="relative group">
            <button class="flex items-center gap-2 text-sm font-medium">
              <div class="w-8 h-8 rounded-full bg-accent-red flex items-center justify-center text-white font-bold text-xs">
                {{ initials }}
              </div>
            </button>
            <!-- Dropdown -->
            <div class="absolute right-0 top-full mt-2 w-48 bg-bg-elevated border border-white/10 rounded-xl
                        shadow-xl opacity-0 invisible group-hover:opacity-100 group-hover:visible
                        transition-all duration-200">
              <div class="px-4 py-3 border-b border-white/5">
                <p class="text-white text-sm font-semibold truncate">{{ auth.user?.full_name }}</p>
                <p class="text-text-muted text-xs truncate">{{ auth.user?.email }}</p>
                <span class="inline-block mt-1 text-xs font-medium px-1.5 py-0.5 rounded capitalize"
                  :class="roleBadgeClass">{{ auth.user?.role }}</span>
              </div>
              <RouterLink to="/profile"
                class="flex items-center gap-2 px-4 py-2.5 text-sm text-text-muted hover:text-white
                       hover:bg-white/5 transition-colors">
                <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                </svg>
                My Profile
              </RouterLink>
              <button @click="handleLogout"
                class="w-full text-left flex items-center gap-2 px-4 py-2.5 text-sm text-text-muted hover:text-white
                       hover:bg-white/5 transition-colors rounded-b-xl">
                <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                </svg>
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

const roleBadgeClass = computed(() => {
  switch (auth.user?.role) {
    case 'admin': return 'bg-accent-red/20 text-accent-red'
    case 'promoter': return 'bg-accent-orange/20 text-accent-orange'
    case 'broadcaster': return 'bg-purple-500/20 text-purple-400'
    default: return 'bg-white/10 text-text-muted'
  }
})

function handleLogout() {
  auth.logout()
  router.push('/')
}
</script>
