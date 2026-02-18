<template>
  <footer class="bg-bg-surface border-t border-white/5 mt-16">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div class="grid grid-cols-1 md:grid-cols-3 gap-8 mb-8">
        <div>
          <div class="flex items-center gap-0.5 mb-4">
            <span class="font-display text-xl text-accent-red tracking-widest">APEX</span>
            <span class="font-display text-xl text-white tracking-widest">CONTACT</span>
          </div>
          <p class="text-text-muted text-sm leading-relaxed max-w-xs">
            Premium live sports streaming. Boxing and car racing events — live, secure, unstoppable.
          </p>
        </div>

        <div>
          <h4 class="text-white font-semibold text-sm mb-4">Sports</h4>
          <ul class="space-y-2">
            <li><a href="#" class="text-text-muted text-sm hover:text-white transition-colors">Boxing</a></li>
            <li><a href="#" class="text-text-muted text-sm hover:text-white transition-colors">Car Racing</a></li>
            <li><a href="#" class="text-text-muted text-sm hover:text-white transition-colors">Upcoming Events</a></li>
          </ul>
        </div>

        <div>
          <h4 class="text-white font-semibold text-sm mb-4">Platform</h4>
          <ul class="space-y-2">
            <li><a href="#" class="text-text-muted text-sm hover:text-white transition-colors">For Promoters</a></li>
            <li><a href="#" class="text-text-muted text-sm hover:text-white transition-colors">Help Center</a></li>
            <li><a href="#" class="text-text-muted text-sm hover:text-white transition-colors">Privacy Policy</a></li>
          </ul>
        </div>
      </div>

      <div class="border-t border-white/5 pt-8 flex flex-col md:flex-row items-center justify-between gap-4">
        <p class="text-text-muted text-xs">
          &copy; {{ new Date().getFullYear() }} ApexContact. All rights reserved.
        </p>
        <p class="text-text-muted text-xs">
          Powered by M-Pesa · Secure streaming · Anti-piracy protected
        </p>
      </div>
    </div>
  </footer>
</template>
