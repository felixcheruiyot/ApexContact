<template>
  <div class="space-y-8">
    <div class="flex items-center justify-between">
      <h1 class="text-white font-bold text-2xl">My Events</h1>
      <RouterLink to="/dashboard/create" class="btn-primary text-sm py-2 px-4">+ New Event</RouterLink>
    </div>

    <!-- Stats row -->
    <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
      <div class="card p-5">
        <p class="text-text-muted text-xs uppercase tracking-wider mb-2">Total Events</p>
        <p class="text-white font-bold text-3xl">{{ events.length }}</p>
      </div>
      <div class="card p-5">
        <p class="text-text-muted text-xs uppercase tracking-wider mb-2">Live Now</p>
        <p class="text-accent-red font-bold text-3xl">
          {{ events.filter(e => e.status === 'live').length }}
        </p>
      </div>
      <div class="card p-5">
        <p class="text-text-muted text-xs uppercase tracking-wider mb-2">Upcoming</p>
        <p class="text-white font-bold text-3xl">
          {{ events.filter(e => e.status === 'scheduled').length }}
        </p>
      </div>
    </div>

    <!-- Events table -->
    <div class="card overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead>
            <tr class="border-b border-white/5">
              <th class="text-left text-text-muted text-xs uppercase tracking-wider px-6 py-4">Event</th>
              <th class="text-left text-text-muted text-xs uppercase tracking-wider px-6 py-4">Sport</th>
              <th class="text-left text-text-muted text-xs uppercase tracking-wider px-6 py-4">Date</th>
              <th class="text-left text-text-muted text-xs uppercase tracking-wider px-6 py-4">Status</th>
              <th class="text-left text-text-muted text-xs uppercase tracking-wider px-6 py-4">Price</th>
              <th class="px-6 py-4" />
            </tr>
          </thead>
          <tbody class="divide-y divide-white/5">
            <tr v-for="event in events" :key="event.id" class="hover:bg-bg-elevated transition-colors">
              <td class="px-6 py-4">
                <p class="text-white font-medium text-sm line-clamp-1">{{ event.title }}</p>
              </td>
              <td class="px-6 py-4">
                <span class="text-text-muted text-sm capitalize">{{ event.sport_type }}</span>
              </td>
              <td class="px-6 py-4">
                <span class="text-text-muted text-sm">{{ format(new Date(event.scheduled_at), 'MMM d, yyyy') }}</span>
              </td>
              <td class="px-6 py-4">
                <span v-if="event.status === 'live'" class="badge-live text-xs">Live</span>
                <span v-else class="badge-upcoming text-xs capitalize">{{ event.status }}</span>
              </td>
              <td class="px-6 py-4">
                <span class="text-white font-semibold text-sm">
                  {{ event.currency }} {{ event.price.toLocaleString() }}
                </span>
              </td>
              <td class="px-6 py-4">
                <RouterLink :to="`/dashboard/analytics/${event.id}`"
                  class="text-accent-red text-sm hover:underline">
                  Analytics
                </RouterLink>
              </td>
            </tr>
          </tbody>
        </table>
        <div v-if="!events.length" class="text-center py-16">
          <p class="text-text-muted text-sm">No events yet.</p>
          <RouterLink to="/dashboard/create" class="btn-primary mt-4 inline-block text-sm py-2 px-4">
            Create your first event
          </RouterLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { format } from 'date-fns'
import { eventsApi } from '@/api/events'
import type { Event } from '@/types'

const events = ref<Event[]>([])

onMounted(async () => {
  const res = await eventsApi.myEvents()
  events.value = res.data.data ?? []
})
</script>
