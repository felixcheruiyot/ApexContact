import client from './client'
import type { ApiResponse, Event, LobbyMessage, LobbyDetail } from '@/types'

export interface CreateCommentaryData {
  title: string
  teaser_hook: string
  description: string
  sport_type: string
  scheduled_at: string
  price: number
  thumbnail_url: string
  event_type?: 'audio_video' | 'audio'
}

export const commentaryApi = {
  list(params?: { page?: number; status?: string; sport?: string }) {
    return client.get<ApiResponse<Event[]>>('/commentary', { params })
  },

  get(id: string) {
    return client.get<ApiResponse<LobbyDetail>>(`/commentary/${id}`)
  },

  create(data: CreateCommentaryData) {
    return client.post<ApiResponse<Event>>('/commentary', data)
  },

  join(id: string, nickname: string) {
    return client.post<ApiResponse<{ joined: boolean; nickname: string; role: string }>>(
      `/commentary/${id}/join`,
      { nickname }
    )
  },

  start(id: string) {
    return client.post<ApiResponse<{ status: string; room: string }>>(`/commentary/${id}/start`)
  },

  end(id: string) {
    return client.post<ApiResponse<{ status: string }>>(`/commentary/${id}/end`)
  },

  getToken(id: string) {
    return client.get<ApiResponse<{ token: string; room: string; livekit_url: string }>>(
      `/commentary/${id}/token`
    )
  },

  messages(id: string, page = 1) {
    return client.get<ApiResponse<LobbyMessage[]>>(`/commentary/${id}/messages`, {
      params: { page },
    })
  },

  suggestNicknames(id: string) {
    return client.get<ApiResponse<string[]>>(`/commentary/${id}/nicknames/suggest`)
  },

  me(id: string) {
    return client.get<ApiResponse<{ has_joined: boolean; nickname?: string; role?: string }>>(
      `/commentary/${id}/me`
    )
  },

  updateParticipant(id: string, userId: string, role: 'speaker' | 'listener') {
    return client.patch<ApiResponse<{ role: string }>>(
      `/commentary/${id}/participants/${userId}`,
      { role }
    )
  },
}
