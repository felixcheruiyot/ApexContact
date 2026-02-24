import client from './client'

export interface GuestStream {
  guest_id: string
  stream_key: string
  rtmp_url: string
  viewer_url: string
  expires_at: string
  time_limit_seconds: number
}

export interface GuestRoom {
  guest_id: string
  room_name: string
  event_type: 'audio' | 'audio_video'
  token: string
  livekit_url: string
  viewer_url: string
  expires_at: string
  time_limit_seconds: number
}

export async function createGuestStream(title?: string): Promise<GuestStream> {
  const { data } = await client.post('/stream/guest', { title })
  return data.data
}

export async function createGuestRoom(eventType: 'audio' | 'audio_video', title?: string): Promise<GuestRoom> {
  const { data } = await client.post('/stream/guest/room', { event_type: eventType, title })
  return data.data
}

export async function getGuestRoomToken(guestId: string): Promise<{ token: string; room_name: string; event_type: string; livekit_url: string }> {
  const { data } = await client.get(`/stream/guest/room/${guestId}`)
  return data.data
}
