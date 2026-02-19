import client from './client'
import type { ApiResponse, User, FraudFlag, Event } from '@/types'

export interface PlatformStats {
  total_users: number
  total_events: number
  live_events: number
  total_revenue: number
  active_viewers: number
  fraud_flags_open: number
}

export const adminApi = {
  analytics() {
    return client.get<ApiResponse<PlatformStats>>('/admin/analytics')
  },

  listUsers() {
    return client.get<ApiResponse<User[]>>('/admin/users')
  },

  listFraudFlags() {
    return client.get<ApiResponse<FraudFlag[]>>('/admin/fraud')
  },

  lockUser(userId: string) {
    return client.post<ApiResponse<string>>(`/admin/users/${userId}/lock`)
  },

  unlockUser(userId: string) {
    return client.post<ApiResponse<string>>(`/admin/users/${userId}/unlock`)
  },

  liveEvents() {
    return client.get<ApiResponse<Event[]>>('/events', { params: { status: 'live' } })
  },

  listAllEvents() {
    return client.get<ApiResponse<Event[]>>('/admin/events')
  },

  updateEvent(eventId: string, data: Partial<Event> & { status?: string }) {
    return client.put<ApiResponse<string>>(`/admin/events/${eventId}`, data)
  },
}
