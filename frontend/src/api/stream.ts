import client from './client'

export interface GuestStream {
  guest_id: string
  stream_key: string
  rtmp_url: string
  viewer_url: string
  expires_at: string
  time_limit_seconds: number
}

export async function createGuestStream(title?: string): Promise<GuestStream> {
  const { data } = await client.post('/stream/guest', { title })
  return data.data
}
