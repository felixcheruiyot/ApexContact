<template>
  <header class="sticky top-0 z-50 bg-bg/90 backdrop-blur-md border-b border-white/5">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">

      <!-- Logo -->
      <RouterLink to="/" class="flex items-center gap-1.5 shrink-0 group">
        <span class="font-display text-2xl tracking-widest text-white group-hover:text-white transition-colors">LIVE</span>
        <span class="w-2 h-2 rounded-full bg-accent-red animate-pulse shrink-0 mx-0.5" />
        <span class="font-display text-2xl tracking-widest text-white group-hover:text-white transition-colors">STREAMIFY</span>
      </RouterLink>

      <!-- Desktop nav -->
      <nav class="hidden md:flex items-center gap-8">
        <RouterLink
          v-for="link in navLinks"
          :key="link.to"
          :to="link.to"
          class="relative text-sm font-medium text-text-muted hover:text-white transition-colors duration-200 py-1 nav-link"
          active-class="text-white"
        >
          {{ link.label }}
          <span class="nav-underline absolute bottom-0 left-0 w-0 h-0.5 bg-white rounded-full transition-all duration-300" />
        </RouterLink>
      </nav>

      <!-- Right side -->
      <div class="flex items-center gap-2 sm:gap-3">
        <template v-if="auth.isAuthenticated">
          <!-- Admin indicator (desktop) -->
          <RouterLink
            v-if="auth.isAdmin"
            to="/admin"
            class="hidden sm:flex items-center gap-1.5 text-accent-red hover:opacity-80 text-sm font-semibold transition-opacity"
          >
            <span class="w-1.5 h-1.5 rounded-full bg-accent-red animate-pulse" />
            Admin
          </RouterLink>
          <!-- Dashboard link (desktop) -->
          <RouterLink
            v-else
            to="/dashboard"
            class="hidden sm:block text-text-muted hover:text-white text-sm font-medium transition-colors"
          >
            Dashboard
          </RouterLink>

          <!-- Avatar dropdown -->
          <div class="relative hidden md:block" ref="dropdownRef">
            <button
              @click="dropdownOpen = !dropdownOpen"
              class="flex items-center gap-2 focus:outline-none"
              aria-haspopup="true"
              :aria-expanded="dropdownOpen"
            >
              <div class="w-9 h-9 rounded-full bg-gradient-to-br from-accent-red to-accent-orange
                          flex items-center justify-center text-white font-bold text-xs ring-2 ring-white/10
                          hover:ring-accent-red/50 transition-all duration-200">
                {{ initials }}
              </div>
            </button>

            <!-- Dropdown panel -->
            <Transition
              enter-active-class="transition duration-150 ease-out"
              enter-from-class="opacity-0 scale-95 translate-y-1"
              enter-to-class="opacity-100 scale-100 translate-y-0"
              leave-active-class="transition duration-100 ease-in"
              leave-from-class="opacity-100 scale-100 translate-y-0"
              leave-to-class="opacity-0 scale-95 translate-y-1"
            >
              <div
                v-if="dropdownOpen"
                class="absolute right-0 top-full mt-2 w-52 bg-bg-elevated border border-white/10
                       rounded-2xl shadow-2xl shadow-black/60 overflow-hidden"
              >
                <!-- User info -->
                <div class="px-4 py-3.5 border-b border-white/5">
                  <p class="text-white text-sm font-semibold truncate">{{ auth.user?.full_name }}</p>
                  <p class="text-text-muted text-xs truncate mt-0.5">{{ auth.user?.email }}</p>
                  <span
                    class="inline-block mt-1.5 text-xs font-semibold px-2 py-0.5 rounded-full capitalize"
                    :class="roleBadgeClass"
                  >
                    {{ auth.isAdmin ? 'Admin' : 'Member' }}
                  </span>
                </div>

                <!-- Links -->
                <div class="py-1">
                  <RouterLink
                    to="/profile"
                    @click="dropdownOpen = false"
                    class="flex items-center gap-2.5 px-4 py-2.5 text-sm text-text-muted
                           hover:text-white hover:bg-white/5 transition-colors"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                    </svg>
                    My Profile
                  </RouterLink>

                  <RouterLink
                    v-if="auth.isAdmin"
                    to="/admin"
                    @click="dropdownOpen = false"
                    class="flex items-center gap-2.5 px-4 py-2.5 text-sm text-accent-red
                           hover:bg-accent-red/10 transition-colors"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
                    </svg>
                    Admin Panel
                  </RouterLink>

                  <RouterLink
                    v-else
                    to="/dashboard"
                    @click="dropdownOpen = false"
                    class="flex items-center gap-2.5 px-4 py-2.5 text-sm text-text-muted
                           hover:text-white hover:bg-white/5 transition-colors"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M3 7h18M3 12h18M3 17h18" />
                    </svg>
                    Dashboard
                  </RouterLink>

                  <RouterLink
                    to="/commentary/create"
                    @click="dropdownOpen = false"
                    class="flex items-center gap-2.5 px-4 py-2.5 text-sm text-accent-orange
                           hover:bg-accent-orange/10 transition-colors"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z" />
                    </svg>
                    Host a Live Room
                  </RouterLink>
                </div>

                <!-- Sign out -->
                <div class="border-t border-white/5 py-1">
                  <button
                    @click="handleLogout"
                    class="w-full text-left flex items-center gap-2.5 px-4 py-2.5 text-sm
                           text-text-muted hover:text-white hover:bg-white/5 transition-colors"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                    </svg>
                    Sign out
                  </button>
                </div>
              </div>
            </Transition>
          </div>
        </template>

        <template v-else>
          <!-- Sign in -->
          <RouterLink
            to="/login"
            class="hidden sm:inline-flex items-center px-5 py-2 rounded-lg border border-white/20
                   hover:border-white/40 text-white text-sm font-medium transition-all duration-200
                   hover:bg-white/5"
          >
            Sign in
          </RouterLink>
          <!-- Try Free — primary action for new visitors -->
          <RouterLink
            to="/try"
            class="inline-flex items-center gap-1.5 px-5 py-2 rounded-lg bg-accent-red hover:bg-accent-red-hover
                   text-white text-sm font-semibold transition-all duration-200 active:scale-95"
          >
            <span class="w-1.5 h-1.5 rounded-full bg-white/80 animate-pulse" />
            Try Free
          </RouterLink>
        </template>

        <!-- Hamburger (mobile) -->
        <button
          @click="mobileMenuOpen = !mobileMenuOpen"
          class="md:hidden p-2 rounded-xl text-text-muted hover:text-white hover:bg-white/5 transition-colors"
          :aria-expanded="mobileMenuOpen"
          aria-label="Toggle menu"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path v-if="!mobileMenuOpen" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M4 6h16M4 12h16M4 18h16" />
            <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </div>

    <!-- Mobile drawer -->
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0 -translate-y-2"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 -translate-y-2"
    >
      <div
        v-if="mobileMenuOpen"
        class="md:hidden border-t border-white/5 bg-bg/98 backdrop-blur-md"
      >
        <div class="px-4 py-6 space-y-1">
          <!-- Nav links -->
          <RouterLink
            v-for="link in navLinks"
            :key="link.to"
            :to="link.to"
            @click="mobileMenuOpen = false"
            class="flex items-center px-4 py-3 rounded-xl text-sm font-medium text-text-muted
                   hover:text-white hover:bg-white/5 transition-colors"
            active-class="text-white bg-white/5"
          >
            {{ link.label }}
          </RouterLink>

          <!-- Divider -->
          <div class="border-t border-white/5 my-3" />

          <template v-if="auth.isAuthenticated">
            <!-- Authenticated mobile links -->
            <RouterLink
              v-if="auth.isAdmin"
              to="/admin"
              @click="mobileMenuOpen = false"
              class="flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-semibold
                     text-accent-red hover:bg-accent-red/10 transition-colors"
            >
              <span class="w-1.5 h-1.5 rounded-full bg-accent-red animate-pulse" />
              Admin Panel
            </RouterLink>
            <RouterLink
              v-else
              to="/dashboard"
              @click="mobileMenuOpen = false"
              class="flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-medium
                     text-text-muted hover:text-white hover:bg-white/5 transition-colors"
            >
              Dashboard
            </RouterLink>

            <RouterLink
              to="/profile"
              @click="mobileMenuOpen = false"
              class="flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-medium
                     text-text-muted hover:text-white hover:bg-white/5 transition-colors"
            >
              My Profile
            </RouterLink>

            <RouterLink
              to="/commentary/create"
              @click="mobileMenuOpen = false"
              class="flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-medium
                     text-accent-orange hover:bg-accent-orange/10 transition-colors"
            >
              Host a Live Room
            </RouterLink>

            <button
              @click="handleLogout"
              class="w-full text-left flex items-center gap-3 px-4 py-3 rounded-xl text-sm
                     font-medium text-text-muted hover:text-white hover:bg-white/5 transition-colors"
            >
              Sign out
            </button>
          </template>

          <template v-else>
            <!-- Unauthenticated mobile buttons -->
            <div class="pt-2 flex flex-col gap-2.5">
              <RouterLink
                to="/try"
                @click="mobileMenuOpen = false"
                class="w-full text-center px-5 py-3 rounded-lg bg-accent-red hover:bg-accent-red-hover
                       text-white text-sm font-semibold transition-all active:scale-95"
              >
                Try Free — No Account Needed
              </RouterLink>
              <RouterLink
                to="/login"
                @click="mobileMenuOpen = false"
                class="w-full text-center px-5 py-3 rounded-lg border border-white/20
                       hover:border-white/40 text-white text-sm font-medium transition-all
                       hover:bg-white/5"
              >
                Sign in
              </RouterLink>
            </div>
          </template>
        </div>
      </div>
    </Transition>
  </header>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { RouterLink, useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const router = useRouter()
const route = useRoute()

const mobileMenuOpen = ref(false)
const dropdownOpen = ref(false)
const dropdownRef = ref<HTMLElement | null>(null)

const navLinks = [
  { to: '/', label: 'Home' },
  { to: '/promoters', label: 'For Creators' },
]

// Close mobile menu on route change
watch(() => route.path, () => {
  mobileMenuOpen.value = false
  dropdownOpen.value = false
})

// Close dropdown on outside click
function handleOutsideClick(e: MouseEvent) {
  if (dropdownRef.value && !dropdownRef.value.contains(e.target as Node)) {
    dropdownOpen.value = false
  }
}

onMounted(() => document.addEventListener('click', handleOutsideClick))
onUnmounted(() => document.removeEventListener('click', handleOutsideClick))

const initials = computed(() => {
  const name = auth.user?.full_name ?? ''
  return name
    .split(' ')
    .map((n) => n[0])
    .join('')
    .slice(0, 2)
    .toUpperCase()
})

const roleBadgeClass = computed(() => {
  switch (auth.user?.role) {
    case 'admin':
      return 'bg-accent-red/20 text-accent-red'
    case 'member':
      return 'bg-accent-orange/20 text-accent-orange'
    default:
      return 'bg-white/10 text-text-muted'
  }
})

function handleLogout() {
  mobileMenuOpen.value = false
  dropdownOpen.value = false
  auth.logout()
  router.push('/')
}
</script>

<style scoped>
.nav-link:hover .nav-underline {
  width: 100%;
}
.nav-link.router-link-active .nav-underline {
  width: 100%;
}
</style>
