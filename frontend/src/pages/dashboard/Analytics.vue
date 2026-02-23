<template>
  <div class="space-y-8">
    <RouterLink to="/dashboard" class="text-text-muted hover:text-white text-sm flex items-center gap-2">
      <ArrowLeft class="w-4 h-4" /> Back to Dashboard
    </RouterLink>

    <div v-if="analytics">
      <h1 class="text-white font-bold text-2xl mb-8">Event Analytics</h1>

      <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div class="card p-5">
          <p class="text-text-muted text-xs uppercase tracking-wider mb-2">Tickets Sold</p>
          <p class="text-white font-bold text-3xl">{{ analytics.total_tickets }}</p>
        </div>
        <div class="card p-5">
          <p class="text-text-muted text-xs uppercase tracking-wider mb-2">Total Revenue</p>
          <p class="text-status-success font-bold text-3xl">
            KES {{ analytics.total_revenue.toLocaleString() }}
          </p>
        </div>
        <div class="card p-5">
          <p class="text-text-muted text-xs uppercase tracking-wider mb-2">Peak Viewers</p>
          <p class="text-white font-bold text-3xl">{{ analytics.peak_viewers }}</p>
        </div>
        <div class="card p-5">
          <p class="text-text-muted text-xs uppercase tracking-wider mb-2">Your Share (70%)</p>
          <p class="text-accent-orange font-bold text-3xl">
            KES {{ (analytics.total_revenue * 0.7).toLocaleString() }}
          </p>
        </div>
      </div>
    </div>

    <div v-else class="flex items-center justify-center py-24">
      <div class="w-8 h-8 border-4 border-accent-red border-t-transparent rounded-full animate-spin" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import { ArrowLeft } from 'lucide-vue-next'
import { eventsApi } from '@/api/events'

const route = useRoute()
const analytics = ref<any>(null)

onMounted(async () => {
  const res = await eventsApi.analytics(route.params.eventId as string)
  analytics.value = res.data.data
})
</script>
