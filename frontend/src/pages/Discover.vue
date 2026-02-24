<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12 space-y-10">

    <!-- Header -->
    <div>
      <p class="text-accent-red font-semibold text-xs uppercase tracking-widest mb-2">Public streams</p>
      <h1 class="font-display text-5xl md:text-6xl uppercase tracking-wide text-white">Discover Streams</h1>
      <p class="text-text-muted mt-3 text-base max-w-xl">
        Browse upcoming and live streams from creators. Private events are only accessible via a direct link.
      </p>
    </div>

    <!-- Search + filter -->
    <div class="flex flex-col sm:flex-row gap-3">
      <div class="relative flex-1">
        <svg class="absolute left-3.5 top-1/2 -translate-y-1/2 w-4 h-4 text-text-muted pointer-events-none"
          fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
        <input
          v-model="searchQuery"
          @input="debouncedSearch"
          type="search"
          placeholder="Search streams…"
          class="input pl-10"
        />
      </div>
      <select v-model="selectedCategory" @change="load" class="input sm:w-48">
        <option value="">All categories</option>
        <option value="sales">Sales</option>
        <option value="mentoring">Mentoring</option>
        <option value="education">Education</option>
        <option value="business">Business</option>
        <option value="fitness">Fitness</option>
        <option value="music">Music</option>
        <option value="gaming">Gaming</option>
        <option value="cooking">Cooking</option>
        <option value="community">Community</option>
        <option value="other">Other</option>
      </select>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-20">
      <div class="w-8 h-8 border-2 border-accent-red border-t-transparent rounded-full animate-spin" />
    </div>

    <!-- Events grid -->
    <div v-else-if="events.length" class="grid sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
      <EventCard v-for="event in events" :key="event.id" :event="event" />
    </div>

    <!-- Empty state -->
    <div v-else class="text-center py-20">
      <svg class="w-12 h-12 text-text-muted mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
          d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <p class="text-text-muted text-sm">No public streams found.</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import EventCard from '@/components/events/EventCard.vue'
import client from '@/api/client'
import type { Event } from '@/types'

const events = ref<Event[]>([])
const loading = ref(true)
const searchQuery = ref('')
const selectedCategory = ref('')

let debounceTimer: ReturnType<typeof setTimeout> | null = null

function debouncedSearch() {
  if (debounceTimer) clearTimeout(debounceTimer)
  debounceTimer = setTimeout(load, 350)
}

async function load() {
  loading.value = true
  try {
    const params: Record<string, string> = {}
    if (searchQuery.value.trim()) params.q = searchQuery.value.trim()
    if (selectedCategory.value) params.sport = selectedCategory.value
    const res = await client.get('/discover', { params })
    events.value = res.data.data ?? []
  } catch {
    events.value = []
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>
