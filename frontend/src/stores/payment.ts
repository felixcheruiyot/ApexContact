import { defineStore } from 'pinia'
import { ref } from 'vue'
import { paymentsApi } from '@/api/payments'

export const usePaymentStore = defineStore('payment', () => {
  const paymentId = ref<string | null>(null)
  const status = ref<string>('idle') // idle | pending | polling | success | failed
  const streamToken = ref<string | null>(null)
  const error = ref<string | null>(null)
  let pollInterval: ReturnType<typeof setInterval> | null = null

  async function initiate(eventId: string, phoneNumber: string) {
    error.value = null
    status.value = 'pending'
    try {
      const res = await paymentsApi.initiate({ event_id: eventId, phone_number: phoneNumber })
      paymentId.value = res.data.data!.payment_id
      startPolling()
    } catch (e: any) {
      error.value = e.response?.data?.error ?? 'Payment initiation failed'
      status.value = 'failed'
    }
  }

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
        } else if (s === 'failed' || s === 'cancelled') {
          status.value = s
          stopPolling()
        }
      } catch {
        // keep polling on transient errors
      }
    }, 3000) // poll every 3 seconds
  }

  function stopPolling() {
    if (pollInterval) {
      clearInterval(pollInterval)
      pollInterval = null
    }
  }

  async function createSubscription(eventId: string, deviceFingerprint: string) {
    if (!paymentId.value) return
    try {
      const res = await paymentsApi.subscribe(eventId, {
        payment_id: paymentId.value,
        device_fingerprint: deviceFingerprint,
      })
      streamToken.value = res.data.data?.stream_token ?? null
    } catch (e: any) {
      error.value = e.response?.data?.error ?? 'Subscription creation failed'
    }
  }

  function reset() {
    stopPolling()
    paymentId.value = null
    status.value = 'idle'
    streamToken.value = null
    error.value = null
  }

  return { paymentId, status, streamToken, error, initiate, createSubscription, reset }
})
