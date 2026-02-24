<template>
  <div class="flex min-h-screen bg-bg">
    <!-- Left panel — hidden on mobile, shown on lg+ -->
    <div class="hidden lg:flex relative w-[45%] flex-shrink-0 overflow-hidden">
      <!-- Background image -->
      <img
        src="https://images.pexels.com/photos/3184291/pexels-photo-3184291.jpeg?auto=compress&cs=tinysrgb&w=1260&q=80"
        alt="Live Streamify — Africa's premier live knowledge platform"
        class="absolute inset-0 w-full h-full object-cover"
      />

      <!-- Dark overlay -->
      <div class="absolute inset-0 bg-gradient-to-br from-black/90 via-black/60 to-accent-red/20" />

      <!-- Panel content -->
      <div class="absolute bottom-16 left-10 right-10 flex flex-col gap-6">
        <!-- Logo -->
        <div>
          <RouterLink to="/" class="inline-block">
            <span class="font-display text-3xl text-white tracking-widest">
              LIVE <span class="text-accent-red">STREAMIFY</span>
            </span>
          </RouterLink>
          <!-- Red underline accent -->
          <div class="mt-1.5 h-0.5 w-16 bg-accent-red rounded-full" />
          <p class="mt-3 text-white/70 text-sm leading-relaxed">
            Go live. Get paid. No setup needed.
          </p>
        </div>

        <!-- Feature bullets -->
        <div class="flex flex-col gap-2.5">
          <span class="inline-flex items-center gap-2 self-start bg-white/10 backdrop-blur-sm
                       rounded-full px-4 py-2 text-white/90 text-sm font-medium">
            <span class="text-status-success">&#10003;</span>
            Stream live — no equipment required
          </span>
          <span class="inline-flex items-center gap-2 self-start bg-white/10 backdrop-blur-sm
                       rounded-full px-4 py-2 text-white/90 text-sm font-medium">
            <span class="text-status-success">&#10003;</span>
            M-Pesa payouts, 70% goes to you
          </span>
          <span class="inline-flex items-center gap-2 self-start bg-white/10 backdrop-blur-sm
                       rounded-full px-4 py-2 text-white/90 text-sm font-medium">
            <span class="text-status-success">&#10003;</span>
            Free to start — pay nothing upfront
          </span>
        </div>
      </div>
    </div>

    <!-- Right panel — full width on mobile, 55% on desktop -->
    <div class="flex flex-col justify-center items-center min-h-screen px-6 py-12 flex-1">
      <!-- Mobile-only logo -->
      <RouterLink
        to="/"
        class="lg:hidden mb-10 inline-block"
      >
        <span class="font-display text-3xl text-white tracking-widest">
          LIVE <span class="text-accent-red">STREAMIFY</span>
        </span>
      </RouterLink>

      <!-- Auth form (Login or Register) -->
      <RouterView />

      <!-- Footer -->
      <p class="mt-10 text-text-muted text-xs text-center">
        &copy; 2026 Live Streamify. All rights reserved.
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { RouterView, RouterLink } from 'vue-router'
</script>
