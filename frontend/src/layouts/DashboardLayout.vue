<template>
  <div class="min-h-screen bg-bg flex">
    <!-- Sidebar -->
    <aside class="w-64 bg-bg-surface border-r border-white/5 flex flex-col shrink-0">
      <RouterLink to="/" class="flex items-center gap-2 px-6 py-5 border-b border-white/5">
        <span class="font-display text-xl text-accent-red tracking-wider">APEX</span>
        <span class="font-display text-xl text-white tracking-wider">CONTACT</span>
      </RouterLink>

      <nav class="flex-1 px-4 py-6 space-y-1">
        <template v-if="auth.isAdmin">
          <NavItem to="/admin" icon="📊" label="Overview" />
          <NavItem to="/admin/users" icon="👥" label="Users" />
          <NavItem to="/admin/fraud" icon="🛡️" label="Fraud Monitor" />
        </template>
        <template v-else>
          <NavItem to="/dashboard" icon="📊" label="Overview" />
          <NavItem to="/dashboard/create" icon="➕" label="Create Event" />
        </template>
      </nav>

      <div class="px-4 py-4 border-t border-white/5">
        <button @click="auth.logout(); router.push('/')"
          class="w-full text-left px-3 py-2 text-text-muted hover:text-white hover:bg-bg-elevated
                 rounded-lg transition-colors text-sm">
          Sign out
        </button>
      </div>
    </aside>

    <!-- Main -->
    <div class="flex-1 flex flex-col min-w-0">
      <header class="bg-bg-surface border-b border-white/5 px-8 py-4 flex items-center justify-between">
        <h1 class="text-text-muted text-sm">
          Welcome back, <span class="text-white font-semibold">{{ auth.user?.full_name }}</span>
        </h1>
        <span class="badge-upcoming capitalize">{{ auth.user?.role }}</span>
      </header>
      <main class="flex-1 p-8 overflow-auto">
        <RouterView />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { RouterView, RouterLink, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import NavItem from '@/components/layout/NavItem.vue'

const auth = useAuthStore()
const router = useRouter()
</script>
