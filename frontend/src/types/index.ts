export type UserRole = 'viewer' | 'promoter' | 'broadcaster' | 'admin'
export type EventStatus = 'draft' | 'pending_review' | 'scheduled' | 'live' | 'completed' | 'cancelled' | 'declined'
export type SportType = 'boxing' | 'racing'
export type PaymentStatus = 'pending' | 'success' | 'failed' | 'cancelled'

export interface User {
  id: string
  email: string
  full_name: string
  phone: string
  role: UserRole
  is_locked: boolean
  age?: number
  gender?: string
  country?: string
  created_at: string
}

export interface Event {
  id: string
  promoter_id: string
  title: string
  description: string
  sport_type: SportType
  scheduled_at: string
  status: EventStatus
  price: number
  currency: string
  thumbnail_url: string
  review_note?: string
  created_at: string
}

export interface Payment {
  id: string
  user_id: string
  event_id: string
  amount: number
  currency: string
  status: PaymentStatus
  phone_number: string
  created_at: string
}

export interface Subscription {
  id: string
  user_id: string
  event_id: string
  payment_id: string
  stream_token: string
  expires_at: string
  created_at: string
}

export interface StreamAnalytics {
  event_id: string
  viewer_count: number
  peak_viewers: number
  total_revenue: number
}

export interface FraudFlag {
  id: string
  user_id: string
  subscription_id: string
  reason: string
  detected_at: string
  resolved: boolean
}

export interface ApiResponse<T> {
  data?: T
  error?: string
  meta?: {
    page: number
    per_page: number
    total: number
  }
}

export interface AuthState {
  token: string | null
  user: User | null
}
