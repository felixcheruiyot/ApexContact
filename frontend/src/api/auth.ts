import client from './client'
import type { ApiResponse, User } from '@/types'

export const authApi = {
  register(data: { email: string; password: string; full_name: string; phone: string }) {
    return client.post<ApiResponse<{ token: string; user: User }>>('/auth/register', data)
  },

  login(data: { email: string; password: string }) {
    return client.post<ApiResponse<{ token: string; user: User }>>('/auth/login', data)
  },

  logout() {
    return client.post('/auth/logout')
  },
}
