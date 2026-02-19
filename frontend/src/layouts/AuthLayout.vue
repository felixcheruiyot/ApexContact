<template>
  <div class="min-h-screen bg-bg flex items-center justify-center p-4">
    <!-- Background gradient -->
    <div class="fixed inset-0 bg-gradient-to-br from-bg via-bg to-accent-red/10 pointer-events-none" />

    <div class="w-full max-w-md relative z-10">
      <!-- Logo -->
      <RouterLink to="/" class="flex items-center justify-center mb-8">
        <span class="font-display text-4xl text-accent-red tracking-widest">LIVE</span>
        <span class="font-display text-4xl text-white tracking-widest">STREAMIFY</span>
      </RouterLink>

      <RouterView />
    </div>
  </div>
</template>

<script setup lang="ts">
import { RouterView, RouterLink } from 'vue-router'
</script>
