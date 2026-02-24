<template>
  <div class="space-y-8 max-w-2xl">
    <div>
      <h1 class="text-white font-bold text-2xl">Withdrawal</h1>
      <p class="text-text-muted text-sm mt-1">Withdraw your earnings to your M-Pesa or bank account.</p>
    </div>

    <!-- Balance card -->
    <div class="card p-6 flex items-center justify-between">
      <div>
        <p class="text-text-muted text-xs uppercase tracking-wider mb-1">Available Balance</p>
        <p class="text-white font-bold text-3xl">
          KES {{ balance !== null ? balance.toLocaleString('en-KE', { minimumFractionDigits: 2 }) : '—' }}
        </p>
        <p class="text-text-muted text-xs mt-1">70% of your completed ticket sales, minus past withdrawals</p>
      </div>
      <svg class="w-10 h-10 text-accent-orange opacity-40" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
          d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z" />
      </svg>
    </div>

    <!-- Payout account setup -->
    <div class="card p-6 space-y-5">
      <div class="flex items-center justify-between">
        <h2 class="text-white font-semibold">Payout Account</h2>
        <span v-if="payoutAccount" class="text-status-success text-xs font-semibold">&#10003; Set up</span>
      </div>

      <div v-if="payoutAccount" class="bg-bg-surface rounded-lg p-4 text-sm space-y-1">
        <p class="text-white font-medium capitalize">{{ payoutAccount.account_type === 'mpesa' ? 'M-Pesa' : 'Bank' }}</p>
        <p class="text-text-muted">{{ payoutAccount.account_name }}</p>
        <p class="text-text-muted font-mono">{{ payoutAccount.account_number }}</p>
        <p v-if="payoutAccount.bank_name" class="text-text-muted">{{ payoutAccount.bank_name }}</p>
        <p class="text-text-muted/50 text-xs mt-2">Contact support to change your payout account.</p>
      </div>

      <form v-else @submit.prevent="savePayoutAccount" class="space-y-4">
        <!-- Account type -->
        <div>
          <label class="block text-text-muted text-sm font-medium mb-2">Account Type</label>
          <div class="flex gap-3">
            <label class="flex items-center gap-2 cursor-pointer">
              <input type="radio" v-model="accountForm.account_type" value="mpesa" class="accent-accent-red" />
              <span class="text-white text-sm">M-Pesa</span>
            </label>
            <label class="flex items-center gap-2 cursor-pointer">
              <input type="radio" v-model="accountForm.account_type" value="bank" class="accent-accent-red" />
              <span class="text-white text-sm">Bank Account</span>
            </label>
          </div>
        </div>

        <div>
          <label class="block text-text-muted text-sm font-medium mb-1.5">Account Name</label>
          <input v-model="accountForm.account_name" type="text" class="input" required
            placeholder="As registered with your provider" />
        </div>

        <div>
          <label class="block text-text-muted text-sm font-medium mb-1.5">
            {{ accountForm.account_type === 'mpesa' ? 'Phone Number (e.g. 0712345678)' : 'Account Number' }}
          </label>
          <input v-model="accountForm.account_number" type="text" class="input" required />
        </div>

        <div v-if="accountForm.account_type === 'bank'">
          <label class="block text-text-muted text-sm font-medium mb-1.5">Bank Name</label>
          <input v-model="accountForm.bank_name" type="text" class="input" />
        </div>

        <div v-if="accountError" class="text-status-error text-sm">{{ accountError }}</div>

        <button type="submit" class="btn-primary text-sm py-2.5 px-6" :disabled="savingAccount">
          <span v-if="savingAccount" class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-2" />
          Save Payout Account
        </button>
      </form>
    </div>

    <!-- Withdrawal form -->
    <div v-if="payoutAccount" class="card p-6 space-y-5">
      <h2 class="text-white font-semibold">Request Withdrawal</h2>
      <p class="text-text-muted text-sm">An OTP will be sent to your email to confirm the withdrawal.</p>

      <form @submit.prevent="initiateWithdrawal" class="space-y-4">
        <div>
          <label class="block text-text-muted text-sm font-medium mb-1.5">Amount (KES)</label>
          <input v-model.number="withdrawAmount" type="number" min="100" step="1" class="input"
            placeholder="Minimum KES 100" required />
        </div>

        <div v-if="withdrawError" class="text-status-error text-sm">{{ withdrawError }}</div>
        <div v-if="withdrawSuccess" class="text-status-success text-sm">{{ withdrawSuccess }}</div>

        <button type="submit" class="btn-primary text-sm py-2.5 px-6" :disabled="initiating">
          <span v-if="initiating" class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-2" />
          Request Withdrawal
        </button>
      </form>
    </div>

    <!-- OTP confirmation modal -->
    <Teleport to="body">
      <div v-if="pendingWithdrawalId" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <div class="absolute inset-0 bg-black/70 backdrop-blur-sm" />
        <div class="relative bg-bg-elevated border border-white/10 rounded-2xl w-full max-w-md p-6 shadow-2xl">
          <h3 class="text-white font-bold text-lg mb-2">Confirm Withdrawal</h3>
          <p class="text-text-muted text-sm mb-5">Enter the 6-digit OTP sent to your email address.</p>

          <form @submit.prevent="confirmWithdrawal" class="space-y-4">
            <input
              v-model="otp"
              type="text"
              maxlength="6"
              pattern="[0-9]{6}"
              placeholder="000000"
              class="input text-center text-2xl font-mono tracking-widest"
              required
              autofocus
            />
            <div v-if="otpError" class="text-status-error text-sm text-center">{{ otpError }}</div>
            <div class="flex gap-3">
              <button type="button" @click="pendingWithdrawalId = null; otp = ''"
                class="btn-ghost flex-1 text-sm py-2.5">Cancel</button>
              <button type="submit" class="btn-primary flex-1 text-sm py-2.5" :disabled="confirming">
                <span v-if="confirming" class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-2" />
                Confirm
              </button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>

    <!-- Withdrawal history -->
    <div v-if="history.length" class="card overflow-hidden">
      <div class="px-6 py-4 border-b border-white/5">
        <h2 class="text-white font-semibold">History</h2>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full min-w-[500px]">
          <thead>
            <tr class="border-b border-white/5">
              <th class="text-left text-text-muted text-xs uppercase tracking-wider px-6 py-3">Date</th>
              <th class="text-left text-text-muted text-xs uppercase tracking-wider px-6 py-3">Amount</th>
              <th class="text-left text-text-muted text-xs uppercase tracking-wider px-6 py-3">Status</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="w in history" :key="w.id" class="border-b border-white/5">
              <td class="px-6 py-3 text-text-muted text-sm">{{ formatDate(w.created_at) }}</td>
              <td class="px-6 py-3 text-white font-semibold text-sm">KES {{ w.amount.toLocaleString() }}</td>
              <td class="px-6 py-3">
                <span :class="statusBadge(w.status)" class="text-xs font-semibold px-2 py-0.5 rounded-full capitalize">
                  {{ w.status.replace('_', ' ') }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { format } from 'date-fns'
import client from '@/api/client'

interface PayoutAccount {
  id: string
  account_type: string
  account_number: string
  account_name: string
  bank_name: string
}

interface WithdrawalRecord {
  id: string
  amount: number
  currency: string
  status: string
  created_at: string
}

const balance = ref<number | null>(null)
const payoutAccount = ref<PayoutAccount | null>(null)
const history = ref<WithdrawalRecord[]>([])

const accountForm = ref({ account_type: 'mpesa', account_number: '', account_name: '', bank_name: '' })
const savingAccount = ref(false)
const accountError = ref('')

const withdrawAmount = ref<number | undefined>(undefined)
const initiating = ref(false)
const withdrawError = ref('')
const withdrawSuccess = ref('')

const pendingWithdrawalId = ref<string | null>(null)
const otp = ref('')
const confirming = ref(false)
const otpError = ref('')

onMounted(async () => {
  await Promise.all([loadBalance(), loadPayoutAccount(), loadHistory()])
})

async function loadBalance() {
  try {
    const res = await client.get('/profile/balance')
    balance.value = res.data.data?.balance ?? 0
  } catch { balance.value = 0 }
}

async function loadPayoutAccount() {
  try {
    const res = await client.get('/profile/payout-account')
    payoutAccount.value = res.data.data ?? null
  } catch { payoutAccount.value = null }
}

async function loadHistory() {
  try {
    const res = await client.get('/profile/withdrawals')
    history.value = res.data.data ?? []
  } catch { history.value = [] }
}

async function savePayoutAccount() {
  accountError.value = ''
  savingAccount.value = true
  try {
    const res = await client.post('/profile/payout-account', accountForm.value)
    payoutAccount.value = res.data.data
  } catch (e: any) {
    accountError.value = e.response?.data?.error ?? 'Failed to save payout account'
  } finally {
    savingAccount.value = false
  }
}

async function initiateWithdrawal() {
  withdrawError.value = ''
  withdrawSuccess.value = ''
  initiating.value = true
  try {
    const res = await client.post('/profile/withdrawals', { amount: withdrawAmount.value })
    pendingWithdrawalId.value = res.data.data?.id
    withdrawSuccess.value = 'OTP sent to your email. Enter it below to confirm.'
  } catch (e: any) {
    withdrawError.value = e.response?.data?.error ?? 'Failed to initiate withdrawal'
  } finally {
    initiating.value = false
  }
}

async function confirmWithdrawal() {
  otpError.value = ''
  confirming.value = true
  try {
    await client.post(`/profile/withdrawals/${pendingWithdrawalId.value}/confirm`, { otp: otp.value })
    pendingWithdrawalId.value = null
    otp.value = ''
    withdrawAmount.value = undefined
    withdrawSuccess.value = 'Withdrawal processed successfully!'
    await Promise.all([loadBalance(), loadHistory()])
  } catch (e: any) {
    otpError.value = e.response?.data?.error ?? 'Invalid or expired OTP'
  } finally {
    confirming.value = false
  }
}

function formatDate(dateStr: string) {
  return format(new Date(dateStr), 'MMM d, yyyy')
}

function statusBadge(status: string) {
  switch (status) {
    case 'completed': return 'bg-status-success/20 text-status-success'
    case 'failed': return 'bg-status-error/20 text-status-error'
    case 'processing': return 'bg-blue-500/20 text-blue-400'
    case 'pending_otp': return 'bg-accent-orange/20 text-accent-orange'
    default: return 'bg-white/10 text-text-muted'
  }
}
</script>
