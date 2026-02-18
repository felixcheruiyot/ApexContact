<template>
  <div class="space-y-8">
    <h1 class="text-white font-bold text-2xl">Platform Overview</h1>

    <div v-if="stats" class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4">
      <div class="card p-5 col-span-1">
        <p class="text-text-muted text-xs uppercase tracking-wider mb-2">Total Users</p>
        <p class="text-white font-bold text-2xl">{{ stats.total_users }}</p>
      </div>
      <div class="card p-5">
        <p class="text-text-muted text-xs uppercase tracking-wider mb-2">Total Events</p>
        <p class="text-white font-bold text-2xl">{{ stats.total_events }}</p>
      </div>
      <div class="card p-5">
        <p class="text-text-muted text-xs uppercase tracking-wider mb-2">Live Events</p>
        <p class="text-accent-red font-bold text-2xl">{{ stats.live_events }}</p>
      </div>
      <div class="card p-5">
        <p class="text-text-muted text-xs uppercase tracking-wider mb-2">Total Revenue</p>
        <p class="text-status-success font-bold text-2xl">
          KES {{ stats.total_revenue.toLocaleString() }}
        </p>
      </div>
      <div class="card p-5">
        <p class="text-text-muted text-xs uppercase tracking-wider mb-2">Active Viewers</p>
        <p class="text-white font-bold text-2xl">{{ stats.active_viewers }}</p>
      </div>
      <div class="card p-5">
        <p class="text-text-muted text-xs uppercase tracking-wider mb-2">Fraud Flags</p>
        <p class="text-status-error font-bold text-2xl">{{ stats.fraud_flags_open }}</p>
      </div>
    </div>

    <div class="flex gap-4 mt-6">
      <RouterLink to="/admin/users" class="btn-ghost text-sm py-2 px-4">Manage Users</RouterLink>
      <RouterLink to="/admin/fraud" class="btn-primary text-sm py-2 px-4">
        View Fraud Flags
        <span v-if="stats?.fraud_flags_open" class="ml-2 bg-white/20 text-white text-xs rounded-full px-1.5 py-0.5">
          {{ stats.fraud_flags_open }}
        </span>
      </RouterLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import client from '@/api/client'

const stats = ref<any>(null)

onMounted(async () => {
  const res = await client.get('/admin/analytics')
  stats.value = res.data.data
})
</script>
