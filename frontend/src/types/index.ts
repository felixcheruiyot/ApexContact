export type UserRole = 'member' | 'admin'
export type EventStatus = 'draft' | 'pending_review' | 'scheduled' | 'live' | 'completed' | 'cancelled' | 'declined'
export type SportType =
  | 'sales'
  | 'mentoring'
  | 'education'
  | 'music'
  | 'fitness'
  | 'legal'
  | 'business'
  | 'gaming'
  | 'cooking'
  | 'visa'
  | 'community'
  | 'other'
export type PaymentStatus = 'pending' | 'success' | 'failed' | 'cancelled'
export type EventType = 'video' | 'audio_video' | 'audio' | 'commentary'
export type LobbyRole = 'host' | 'speaker' | 'listener'

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
  event_type: EventType
  teaser_hook: string
  is_public: boolean
  created_at: string
}

export interface LobbyParticipant {
  id: string
  event_id: string
  user_id: string
  nickname: string
  role: LobbyRole
  joined_at: string
}

export interface LobbyMessage {
  id: string
  event_id: string
  user_id: string
  nickname: string
  content: string
  message_type: 'text' | 'meme'
  created_at: string
}

export interface ChatEvent {
  type: 'message' | 'reaction' | 'joined' | 'left' | 'speaker_granted' | 'speaker_revoked'
  nickname: string
  content?: string
  user_id?: string
  created_at: string
}

export interface LobbyDetail {
  event: Event
  participant_count: number
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
