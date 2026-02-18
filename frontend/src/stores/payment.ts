import { defineStore } from 'pinia'
import { ref } from 'vue'
import { paymentsApi } from '@/api/payments'

const POLL_INTERVAL_MS = 3000
const POLL_TIMEOUT_MS = 5 * 60 * 1000 // 5 minutes — matches M-Pesa STK prompt expiry

function tokenKey(eventId: string) {
  return `stream_token:${eventId}`
}

export const usePaymentStore = defineStore('payment', () => {
  const paymentId = ref<string | null>(null)
  const status = ref<'idle' | 'pending' | 'polling' | 'success' | 'failed' | 'cancelled' | 'timeout'>('idle')
  const streamToken = ref<string | null>(null)
  const error = ref<string | null>(null)

  let pollInterval: ReturnType<typeof setInterval> | null = null
  let pollTimeout: ReturnType<typeof setTimeout> | null = null

  // ── Token persistence ────────────────────────────────────────────────────

  function getStoredToken(eventId: string): string | null {
    return localStorage.getItem(tokenKey(eventId))
  }

  function clearStoredToken(eventId: string) {
    localStorage.removeItem(tokenKey(eventId))
  }

  // ── Payment initiation ───────────────────────────────────────────────────

  async function initiate(eventId: string, phoneNumber: string) {
    error.value = null
    status.value = 'pending'
    try {
      const res = await paymentsApi.initiate({ event_id: eventId, phone_number: phoneNumber })
      paymentId.value = res.data.data!.payment_id
      startPolling()
    } catch (e: any) {
      error.value = e.response?.data?.error ?? 'Payment initiation failed. Please try again.'
      status.value = 'failed'
    }
  }

  // ── Polling ──────────────────────────────────────────────────────────────

  function startPolling() {
    status.value = 'polling'

    pollInterval = setInterval(async () => {
      if (!paymentId.value) return
      try {
        const res = await paymentsApi.status(paymentId.value)
        const s = res.data.data?.status
        if (s === 'success') {
          status.value = 'success'
          stopPolling()
        } else if (s === 'failed') {
          error.value = 'Payment was declined. Please try again.'
          status.value = 'failed'
          stopPolling()
        } else if (s === 'cancelled') {
          error.value = 'Payment was cancelled.'
          status.value = 'cancelled'
          stopPolling()
        }
      } catch {
        // transient network error — keep polling
      }
    }, POLL_INTERVAL_MS)

    // Hard timeout: if the webhook never fires within 5 min, give up
    pollTimeout = setTimeout(() => {
      if (status.value === 'polling') {
        error.value = 'Payment timed out. If you were charged, please contact support.'
        status.value = 'timeout'
        stopPolling()
      }
    }, POLL_TIMEOUT_MS)
  }

  function stopPolling() {
    if (pollInterval) { clearInterval(pollInterval); pollInterval = null }
    if (pollTimeout)  { clearTimeout(pollTimeout);   pollTimeout  = null }
  }

  // ── Subscription ─────────────────────────────────────────────────────────

  async function createSubscription(eventId: string, deviceFingerprint: string) {
    if (!paymentId.value) return
    try {
      const res = await paymentsApi.subscribe(eventId, {
        payment_id: paymentId.value,
        device_fingerprint: deviceFingerprint,
      })
      const token = res.data.data?.stream_token ?? null
      streamToken.value = token
      if (token) {
        // Persist so EventDetail can show "Watch Now" without re-purchasing
        localStorage.setItem(tokenKey(eventId), token)
      }
    } catch (e: any) {
      error.value = e.response?.data?.error ?? 'Could not activate your subscription. Please contact support.'
      throw e // re-throw so MpesaModal can show a specific error state
    }
  }

  // ── Reset ────────────────────────────────────────────────────────────────

  function reset() {
    stopPolling()
    paymentId.value = null
    status.value = 'idle'
    streamToken.value = null
    error.value = null
  }

  return {
    paymentId,
    status,
    streamToken,
    error,
    initiate,
    createSubscription,
    getStoredToken,
    clearStoredToken,
    reset,
  }
})
