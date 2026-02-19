import { defineStore } from 'pinia'
import { ref } from 'vue'
import { commentaryApi } from '@/api/commentary'
import type { Event, LobbyParticipant, LobbyMessage, LobbyRole } from '@/types'

export const useCommentaryStore = defineStore('commentary', () => {
  const lobbies = ref<Event[]>([])
  const current = ref<Event | null>(null)
  const participantCount = ref(0)
  const participants = ref<LobbyParticipant[]>([])
  const messages = ref<LobbyMessage[]>([])
  const myNickname = ref<string | null>(null)
  const myRole = ref<LobbyRole>('listener')
  const livekitToken = ref<string | null>(null)
  const livekitUrl = ref<string | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchLobbies(params?: { sport?: string; status?: string }) {
    loading.value = true
    error.value = null
    try {
      const res = await commentaryApi.list(params)
      lobbies.value = res.data.data ?? []
    } catch (e: any) {
      error.value = e.response?.data?.error ?? 'Failed to load lobbies'
    } finally {
      loading.value = false
    }
  }

  async function fetchDetail(id: string) {
    loading.value = true
    error.value = null
    try {
      const res = await commentaryApi.get(id)
      const detail = res.data.data
      if (detail) {
        current.value = detail.event
        participantCount.value = detail.participant_count
      }
    } catch (e: any) {
      error.value = e.response?.data?.error ?? 'Lobby not found'
    } finally {
      loading.value = false
    }
  }

  async function joinLobby(id: string, nickname: string) {
    const res = await commentaryApi.join(id, nickname)
    const data = res.data.data
    if (data) {
      myNickname.value = data.nickname
      myRole.value = data.role as LobbyRole
    }
    return data
  }

  async function startRoom(id: string) {
    const res = await commentaryApi.start(id)
    if (current.value) {
      current.value = { ...current.value, status: 'live' }
    }
    return res.data.data
  }

  async function endRoom(id: string) {
    await commentaryApi.end(id)
    if (current.value) {
      current.value = { ...current.value, status: 'completed' }
    }
  }

  async function fetchToken(id: string) {
    const res = await commentaryApi.getToken(id)
    const data = res.data.data
    if (data) {
      livekitToken.value = data.token
      livekitUrl.value = data.livekit_url
    }
    return data
  }

  async function fetchMessages(id: string, page = 1) {
    const res = await commentaryApi.messages(id, page)
    if (page === 1) {
      messages.value = res.data.data ?? []
    } else {
      messages.value = [...messages.value, ...(res.data.data ?? [])]
    }
    return res.data
  }

  function addMessage(msg: LobbyMessage) {
    messages.value.push(msg)
  }

  function reset() {
    current.value = null
    participants.value = []
    messages.value = []
    myNickname.value = null
    myRole.value = 'listener'
    livekitToken.value = null
    livekitUrl.value = null
  }

  const liveLobbies = () => lobbies.value.filter((l) => l.status === 'live')
  const upcomingLobbies = () => lobbies.value.filter((l) => l.status === 'scheduled')

  return {
    lobbies, current, participantCount, participants, messages,
    myNickname, myRole, livekitToken, livekitUrl, loading, error,
    fetchLobbies, fetchDetail, joinLobby, startRoom, endRoom,
    fetchToken, fetchMessages, addMessage, reset,
    liveLobbies, upcomingLobbies,
  }
})
