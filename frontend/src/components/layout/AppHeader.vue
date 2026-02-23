<template>
  <header class="sticky top-0 z-50 bg-bg/80 backdrop-blur-md border-b border-white/5">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
      <!-- Logo -->
      <RouterLink to="/" class="flex items-center gap-2 shrink-0">
        <svg width="28" height="28" viewBox="0 0 28 28" fill="none" xmlns="http://www.w3.org/2000/svg">
          <circle cx="14" cy="14" r="14" fill="#E8002D"/>
          <polygon points="11,8 22,14 11,20" fill="white"/>
          <rect x="5" y="11" width="2.5" height="6" rx="1" fill="white"/>
          <rect x="8.5" y="9" width="2.5" height="10" rx="1" fill="white"/>
        </svg>
        <span class="font-display text-2xl text-accent-red tracking-widest">LIVE</span>
        <span class="font-display text-2xl text-white tracking-widest">STREAMIFY</span>
      </RouterLink>

      <!-- Desktop nav links -->
      <nav class="hidden md:flex items-center gap-6">
        <RouterLink to="/" class="text-text-muted hover:text-white transition-colors text-sm font-medium"
          active-class="text-white">
          Home
        </RouterLink>
        <RouterLink to="/use-cases" class="text-text-muted hover:text-white transition-colors text-sm font-medium"
          active-class="text-white">
          Use Cases
        </RouterLink>
        <RouterLink to="/" class="flex items-center gap-1.5 text-text-muted hover:text-accent-orange transition-colors text-sm font-medium"
          :class="{ 'text-accent-orange': isCommentaryRoute }">
          <Mic class="w-4 h-4" /> Live Rooms
        </RouterLink>
        <RouterLink to="/promoters" class="text-text-muted hover:text-white transition-colors text-sm font-medium"
          active-class="text-white">
          For Hosts
        </RouterLink>
      </nav>

      <!-- Right side -->
      <div class="flex items-center gap-2 sm:gap-3">
        <template v-if="auth.isAuthenticated">
          <RouterLink v-if="auth.isAdmin" to="/admin"
            class="text-accent-red hover:opacity-80 text-sm font-semibold transition-opacity hidden sm:flex items-center gap-1.5">
            <span class="w-1.5 h-1.5 rounded-full bg-accent-red animate-pulse" />
            Admin Panel
          </RouterLink>
          <RouterLink v-else to="/dashboard"
            class="text-text-muted hover:text-white text-sm font-medium transition-colors hidden sm:block">
            Dashboard
          </RouterLink>

          <!-- User avatar dropdown (desktop) -->
          <div class="relative group hidden md:block">
            <button class="flex items-center gap-2 text-sm font-medium">
              <div class="w-8 h-8 rounded-full bg-accent-red flex items-center justify-center text-white font-bold text-xs">
                {{ initials }}
              </div>
            </button>
            <div class="absolute right-0 top-full mt-2 w-48 bg-bg-elevated border border-white/10 rounded-xl
                        shadow-xl opacity-0 invisible group-hover:opacity-100 group-hover:visible
                        transition-all duration-200">
              <div class="px-4 py-3 border-b border-white/5">
                <p class="text-white text-sm font-semibold truncate">{{ auth.user?.full_name }}</p>
                <p class="text-text-muted text-xs truncate">{{ auth.user?.email }}</p>
                <span class="inline-block mt-1 text-xs font-medium px-1.5 py-0.5 rounded capitalize"
                  :class="roleBadgeClass">{{ auth.isAdmin ? 'Admin' : 'Member' }}</span>
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
              <RouterLink to="/commentary/create"
                class="flex items-center gap-2 px-4 py-2.5 text-sm text-accent-orange hover:text-orange-400
                       hover:bg-accent-orange/5 transition-colors">
                <Mic class="w-3.5 h-3.5" />
                Host a Live Room
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
          <RouterLink to="/login" class="btn-ghost text-sm py-2 px-3 sm:px-4 hidden sm:inline-flex">Sign in</RouterLink>
          <RouterLink to="/register" class="btn-primary text-sm py-2 px-3 sm:px-4">Get Started</RouterLink>
        </template>

        <!-- Hamburger (mobile only) -->
        <button @click="mobileMenuOpen = !mobileMenuOpen"
          class="md:hidden p-2 rounded-lg text-text-muted hover:text-white hover:bg-white/5 transition-colors"
          :aria-expanded="mobileMenuOpen">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path v-if="!mobileMenuOpen" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M4 6h16M4 12h16M4 18h16" />
            <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </div>

    <!-- Mobile menu drawer -->
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0 -translate-y-2"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 -translate-y-2"
    >
      <div v-if="mobileMenuOpen"
        class="md:hidden border-t border-white/5 bg-bg/95 backdrop-blur-md">
        <nav class="px-4 py-4 space-y-1">
          <RouterLink to="/" @click="mobileMenuOpen = false"
            class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium text-text-muted hover:text-white hover:bg-white/5 transition-colors"
            active-class="text-white bg-white/5">
            Home
          </RouterLink>
          <RouterLink to="/use-cases" @click="mobileMenuOpen = false"
            class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium text-text-muted hover:text-white hover:bg-white/5 transition-colors"
            active-class="text-white bg-white/5">
            <Lightbulb class="w-4 h-4" /> Use Cases
          </RouterLink>

          <!-- Live Rooms section -->
          <div class="border-t border-white/5 pt-3 mt-3">
            <p class="px-3 pb-1 text-[10px] font-bold uppercase tracking-widest text-text-muted">Live Rooms</p>
            <RouterLink to="/" @click="mobileMenuOpen = false"
              class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium text-text-muted hover:text-accent-orange hover:bg-accent-orange/5 transition-colors">
              <Mic class="w-4 h-4" /> Browse Rooms
            </RouterLink>
            <RouterLink v-if="auth.isAuthenticated" to="/commentary/create" @click="mobileMenuOpen = false"
              class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium text-accent-orange hover:bg-accent-orange/10 transition-colors">
              + Host a Live Room
            </RouterLink>
            <RouterLink v-else to="/login" @click="mobileMenuOpen = false"
              class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium text-text-muted hover:text-accent-orange hover:bg-accent-orange/5 transition-colors">
              + Host a Live Room
            </RouterLink>
          </div>

          <div class="border-t border-white/5 pt-3 mt-3 space-y-1">
            <RouterLink to="/promoters" @click="mobileMenuOpen = false"
              class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium text-text-muted hover:text-white hover:bg-white/5 transition-colors">
              For Hosts
            </RouterLink>

            <template v-if="auth.isAuthenticated">
              <RouterLink v-if="auth.isAdmin" to="/admin" @click="mobileMenuOpen = false"
                class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-semibold text-accent-red hover:bg-accent-red/10 transition-colors">
                <span class="w-1.5 h-1.5 rounded-full bg-accent-red animate-pulse" />
                Admin Panel
              </RouterLink>
              <RouterLink v-else to="/dashboard" @click="mobileMenuOpen = false"
                class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium text-text-muted hover:text-white hover:bg-white/5 transition-colors">
                Dashboard
              </RouterLink>
              <RouterLink to="/profile" @click="mobileMenuOpen = false"
                class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium text-text-muted hover:text-white hover:bg-white/5 transition-colors">
                My Profile
              </RouterLink>
              <button @click="handleLogout"
                class="w-full text-left flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium text-text-muted hover:text-white hover:bg-white/5 transition-colors">
                Sign out
              </button>
            </template>
            <template v-else>
              <RouterLink to="/login" @click="mobileMenuOpen = false"
                class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium text-text-muted hover:text-white hover:bg-white/5 transition-colors">
                Sign in
              </RouterLink>
              <RouterLink to="/register" @click="mobileMenuOpen = false"
                class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium text-accent-red hover:bg-accent-red/10 transition-colors">
                Get Started
              </RouterLink>
            </template>
          </div>
        </nav>
      </div>
    </Transition>
  </header>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { RouterLink, useRouter, useRoute } from 'vue-router'
import { Mic, Lightbulb } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const router = useRouter()
const route = useRoute()

const mobileMenuOpen = ref(false)

// Close mobile menu on navigation
watch(() => route.path, () => { mobileMenuOpen.value = false })

const isCommentaryRoute = computed(() => route.path.startsWith('/commentary'))

const initials = computed(() => {
  const name = auth.user?.full_name ?? ''
  return name.split(' ').map((n) => n[0]).join('').slice(0, 2).toUpperCase()
})

const roleBadgeClass = computed(() => {
  switch (auth.user?.role) {
    case 'admin': return 'bg-accent-red/20 text-accent-red'
    case 'member': return 'bg-accent-orange/20 text-accent-orange'
    default: return 'bg-white/10 text-text-muted'
  }
})

function handleLogout() {
  mobileMenuOpen.value = false
  auth.logout()
  router.push('/')
}
</script>
