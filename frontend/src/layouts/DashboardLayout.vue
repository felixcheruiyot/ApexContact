<template>
  <div class="min-h-screen bg-bg flex">

    <!-- Mobile backdrop -->
    <Transition
      enter-active-class="transition-opacity duration-200"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition-opacity duration-200"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div v-if="sidebarOpen"
        class="fixed inset-0 z-30 bg-black/60 lg:hidden"
        @click="sidebarOpen = false"
      />
    </Transition>

    <!-- Sidebar -->
    <aside
      :class="[
        'fixed inset-y-0 left-0 z-40 w-64 bg-bg-surface border-r border-white/5 flex flex-col shrink-0',
        'transition-transform duration-200 ease-in-out',
        'lg:relative lg:translate-x-0',
        sidebarOpen ? 'translate-x-0' : '-translate-x-full',
      ]"
    >
      <!-- Logo -->
      <RouterLink to="/" class="flex items-center gap-2 px-5 py-5 border-b border-white/5 shrink-0">
        <svg width="22" height="22" viewBox="0 0 28 28" fill="none" xmlns="http://www.w3.org/2000/svg">
          <circle cx="14" cy="14" r="14" fill="#E8002D"/>
          <polygon points="11,8 22,14 11,20" fill="white"/>
          <rect x="5" y="11" width="2.5" height="6" rx="1" fill="white"/>
          <rect x="8.5" y="9" width="2.5" height="10" rx="1" fill="white"/>
        </svg>
        <span class="font-display text-lg text-accent-red tracking-wider">LIVE</span>
        <span class="font-display text-lg text-white tracking-wider">STREAMIFY</span>
      </RouterLink>

      <!-- Role section badge -->
      <div class="px-4 pt-4 pb-2 shrink-0">
        <div :class="sectionBadgeClass"
          class="flex items-center gap-2 px-3 py-2 rounded-lg text-xs font-bold uppercase tracking-widest">
          <span class="w-2 h-2 rounded-full animate-pulse" :class="sectionDotClass" />
          {{ sectionLabel }}
        </div>
      </div>

      <!-- Nav -->
      <nav class="flex-1 px-3 py-3 space-y-0.5 overflow-y-auto">
        <template v-if="auth.isAdmin">
          <NavItem to="/admin" :exact="true" icon="grid" label="Platform Overview" />
          <NavItem to="/admin/users" icon="users" label="Users" />
          <NavItem to="/admin/events" icon="film" label="Events" />
          <NavItem to="/admin/fraud" icon="shield" label="Fraud Monitor" />
        </template>
        <template v-else-if="auth.user?.role === 'promoter'">
          <NavItem to="/dashboard" :exact="true" icon="grid" label="My Events" />
          <NavItem to="/dashboard/create" icon="plus" label="Create Event" />
          <NavItem to="/dashboard/revenue" icon="bar-chart" label="Revenue" />
          <div class="pt-2 pb-1 px-2">
            <p class="text-[10px] font-bold uppercase tracking-widest text-text-muted/60">Commentary</p>
          </div>
          <NavItem to="/commentary/create" icon="mic" label="Host a Lobby" />
        </template>
        <template v-else-if="auth.user?.role === 'broadcaster'">
          <NavItem to="/dashboard" :exact="true" icon="grid" label="Overview" />
        </template>
      </nav>

      <!-- Bottom: profile + sign out -->
      <div class="px-3 pb-4 border-t border-white/5 shrink-0 pt-3 space-y-0.5">
        <NavItem to="/profile" icon="user" label="My Profile" />

        <div class="flex items-center gap-3 px-3 py-2 mt-1">
          <div class="w-8 h-8 rounded-full flex items-center justify-center shrink-0"
            :class="avatarBgClass">
            <span class="text-xs font-bold" :class="avatarTextClass">{{ initials }}</span>
          </div>
          <div class="min-w-0 flex-1">
            <p class="text-white text-sm font-medium truncate">{{ auth.user?.full_name }}</p>
            <p class="text-xs capitalize truncate" :class="avatarTextClass">{{ auth.user?.role }}</p>
          </div>
        </div>

        <button
          @click="handleSignOut"
          class="w-full text-left px-3 py-2 text-text-muted hover:text-white hover:bg-bg-elevated
                 rounded-lg transition-colors text-sm flex items-center gap-2.5"
        >
          <svg class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
          </svg>
          Sign out
        </button>
      </div>
    </aside>

    <!-- Main content -->
    <div class="flex-1 flex flex-col min-w-0">
      <!-- Top bar -->
      <header class="h-14 border-b border-white/5 px-4 sm:px-8 flex items-center justify-between shrink-0 bg-bg-surface/50">
        <div class="flex items-center gap-3">
          <!-- Hamburger — mobile only -->
          <button
            @click="sidebarOpen = !sidebarOpen"
            class="lg:hidden p-1.5 rounded-lg text-text-muted hover:text-white hover:bg-white/5 transition-colors"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M4 6h16M4 12h16M4 18h16" />
            </svg>
          </button>
          <p class="text-white font-semibold text-sm">{{ pageTitle }}</p>
        </div>
        <div class="flex items-center gap-3">
          <span class="text-text-muted text-xs hidden sm:block">{{ auth.user?.email }}</span>
          <RouterLink to="/profile"
            class="w-8 h-8 rounded-full flex items-center justify-center text-xs font-bold transition-opacity hover:opacity-80"
            :class="[avatarBgClass, avatarTextClass]"
          >
            {{ initials }}
          </RouterLink>
        </div>
      </header>

      <main class="flex-1 p-4 sm:p-6 lg:p-8 overflow-auto">
        <RouterView />
      </main>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { RouterView, RouterLink, useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import NavItem from '@/components/layout/NavItem.vue'

const auth = useAuthStore()
const router = useRouter()
const route = useRoute()

const sidebarOpen = ref(false)

// Close sidebar on navigation (mobile)
watch(() => route.path, () => { sidebarOpen.value = false })

const initials = computed(() =>
  (auth.user?.full_name ?? '')
    .split(' ')
    .slice(0, 2)
    .map((n: string) => n[0])
    .join('')
    .toUpperCase() || '?',
)

const sectionLabel = computed(() => {
  if (auth.isAdmin) return 'Admin Panel'
  if (auth.user?.role === 'broadcaster') return 'Broadcaster'
  return 'Dashboard'
})

const sectionBadgeClass = computed(() => {
  if (auth.isAdmin) return 'bg-accent-red/10 border border-accent-red/20 text-accent-red'
  if (auth.user?.role === 'broadcaster') return 'bg-purple-500/10 border border-purple-500/20 text-purple-400'
  return 'bg-accent-orange/10 border border-accent-orange/20 text-accent-orange'
})

const sectionDotClass = computed(() => {
  if (auth.isAdmin) return 'bg-accent-red'
  if (auth.user?.role === 'broadcaster') return 'bg-purple-400'
  return 'bg-accent-orange'
})

const avatarBgClass = computed(() => {
  if (auth.isAdmin) return 'bg-accent-red/20'
  if (auth.user?.role === 'broadcaster') return 'bg-purple-500/20'
  return 'bg-accent-orange/20'
})

const avatarTextClass = computed(() => {
  if (auth.isAdmin) return 'text-accent-red'
  if (auth.user?.role === 'broadcaster') return 'text-purple-400'
  return 'text-accent-orange'
})

const pageTitle = computed(() => {
  const p = route.path
  if (p === '/admin') return 'Platform Overview'
  if (p === '/dashboard') return 'My Events'
  if (p.startsWith('/admin/users')) return 'User Management'
  if (p.startsWith('/admin/events')) return 'Event Management'
  if (p.startsWith('/admin/fraud')) return 'Fraud Monitor'
  if (p.startsWith('/dashboard/create')) return 'Create Event'
  if (p.startsWith('/dashboard/edit')) return 'Edit Event'
  if (p.startsWith('/dashboard/analytics')) return 'Event Analytics'
  if (p.startsWith('/dashboard/revenue')) return 'Revenue & Analytics'
  if (p.startsWith('/dashboard/commentary')) return 'Host a Commentary'
  if (p.startsWith('/profile')) return 'My Profile'
  return sectionLabel.value
})

function handleSignOut() {
  auth.logout()
  router.push('/')
}
</script>
