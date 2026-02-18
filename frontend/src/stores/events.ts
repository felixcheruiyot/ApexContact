import { defineStore } from 'pinia'
import { ref } from 'vue'
import { eventsApi } from '@/api/events'
import type { Event } from '@/types'

export const useEventsStore = defineStore('events', () => {
  const events = ref<Event[]>([])
  const currentEvent = ref<Event | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchEvents(params?: { sport?: string; status?: string }) {
    loading.value = true
    error.value = null
    try {
      const res = await eventsApi.list(params)
      events.value = res.data.data ?? []
    } catch (e: any) {
      error.value = e.response?.data?.error ?? 'Failed to load events'
    } finally {
      loading.value = false
    }
  }

  async function fetchEvent(id: string) {
    loading.value = true
    error.value = null
    try {
      const res = await eventsApi.get(id)
      currentEvent.value = res.data.data ?? null
    } catch (e: any) {
      error.value = e.response?.data?.error ?? 'Event not found'
    } finally {
      loading.value = false
    }
  }

  const liveEvents = () => events.value.filter((e) => e.status === 'live')
  const upcomingEvents = () => events.value.filter((e) => e.status === 'scheduled')

  return { events, currentEvent, loading, error, fetchEvents, fetchEvent, liveEvents, upcomingEvents }
})
