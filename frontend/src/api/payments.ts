import client from './client'
import type { ApiResponse } from '@/types'

export const paymentsApi = {
  initiate(data: { event_id: string; phone_number: string }) {
    return client.post<ApiResponse<{ payment_id: string; amount: number; currency: string; message: string }>>(
      '/payments/initiate',
      data
    )
  },

  status(paymentId: string) {
    return client.get<ApiResponse<{ status: string }>>(`/payments/status/${paymentId}`)
  },

  subscribe(eventId: string, data: { payment_id: string; device_fingerprint: string }) {
    return client.post<ApiResponse<{ stream_token: string }>>(`/stream/${eventId}/subscribe`, data)
  },

  // token is optional — admins can omit it (backend detects admin role from JWT)
  getStreamToken(eventId: string, token?: string) {
    return client.get<ApiResponse<{ hls_url: string }>>(`/stream/${eventId}/token`, {
      params: token ? { token } : undefined,
    })
  },
}
