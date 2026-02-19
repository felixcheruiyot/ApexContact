import client from './client'
import type { ApiResponse, User } from '@/types'

export interface ProfileUpdateData {
  full_name: string
  phone: string
  age?: number | null
  gender: string
  country: string
}

export const profileApi = {
  get() {
    return client.get<ApiResponse<User>>('/profile')
  },

  update(data: ProfileUpdateData) {
    return client.put<ApiResponse<string>>('/profile', data)
  },
}
