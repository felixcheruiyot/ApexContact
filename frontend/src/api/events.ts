import client from './client'
import type { ApiResponse, Event } from '@/types'

export const eventsApi = {
  list(params?: { sport?: string; status?: string }) {
    return client.get<ApiResponse<Event[]>>('/events', { params })
  },

  get(id: string) {
    return client.get<ApiResponse<Event>>(`/events/${id}`)
  },

  create(data: Partial<Event>) {
    return client.post<ApiResponse<Event>>('/events', data)
  },

  update(id: string, data: Partial<Event>) {
    return client.put<ApiResponse<Event>>(`/events/${id}`, data)
  },

  // Promoter endpoints
  myEvents() {
    return client.get<ApiResponse<Event[]>>('/promoter/events')
  },

  analytics(eventId: string) {
    return client.get(`/promoter/analytics/${eventId}`)
  },

  revenue() {
    return client.get('/promoter/revenue')
  },
}
