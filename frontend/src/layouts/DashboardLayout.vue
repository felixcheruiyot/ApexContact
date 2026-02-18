<template>
  <div class="min-h-screen bg-bg flex">

    <!-- Sidebar -->
    <aside class="w-60 bg-bg-surface border-r border-white/5 flex flex-col shrink-0">

      <!-- Logo -->
      <RouterLink to="/" class="flex items-center gap-2 px-5 py-5 border-b border-white/5 shrink-0">
        <span class="font-display text-xl text-accent-red tracking-wider">APEX</span>
        <span class="font-display text-xl text-white tracking-wider">CONTACT</span>
      </RouterLink>

      <!-- Nav -->
      <nav class="flex-1 px-3 py-5 space-y-0.5 overflow-y-auto">
        <template v-if="auth.isAdmin">
          <NavItem to="/admin" exact icon="📊" label="Overview" />
          <NavItem to="/admin/users" icon="👥" label="Users" />
          <NavItem to="/admin/fraud" icon="🛡️" label="Fraud Monitor" />
        </template>
        <template v-else>
          <NavItem to="/dashboard" exact icon="📊" label="Overview" />
          <NavItem to="/dashboard/create" icon="➕" label="Create Event" />
        </template>
      </nav>

      <!-- User info + sign out -->
      <div class="px-3 py-4 border-t border-white/5 shrink-0">
        <div class="flex items-center gap-3 px-3 py-2 mb-1">
          <div class="w-8 h-8 rounded-full bg-accent-red/20 flex items-center justify-center shrink-0">
            <span class="text-accent-red text-xs font-bold">{{ initials }}</span>
          </div>
          <div class="min-w-0">
            <p class="text-white text-sm font-medium truncate">{{ auth.user?.full_name }}</p>
            <p class="text-text-muted text-xs capitalize truncate">{{ auth.user?.role }}</p>
          </div>
        </div>
        <button
          @click="handleSignOut"
          class="w-full text-left px-3 py-2 text-text-muted hover:text-white hover:bg-bg-elevated
                 rounded-lg transition-colors text-sm flex items-center gap-2"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
          </svg>
          Sign out
        </button>
      </div>
    </aside>

    <!-- Main content -->
    <div class="flex-1 flex flex-col min-w-0">
      <main class="flex-1 p-8 overflow-auto">
        <RouterView />
      </main>
    </div>

  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { RouterView, RouterLink, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import NavItem from '@/components/layout/NavItem.vue'

const auth = useAuthStore()
const router = useRouter()

const initials = computed(() => {
  return (auth.user?.full_name ?? '')
    .split(' ')
    .slice(0, 2)
    .map((n: string) => n[0])
    .join('')
    .toUpperCase() || '?'
})

function handleSignOut() {
  auth.logout()
  router.push('/')
}
</script>
