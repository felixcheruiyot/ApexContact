/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        // ApexContact design system
        bg: {
          DEFAULT: '#0A0A0F',
          surface: '#141418',
          elevated: '#1E1E26',
          card: '#17171C',
        },
        accent: {
          red: '#E8002D',
          'red-hover': '#C0001F',
          orange: '#FF6B00',
          gold: '#F59E0B',
          'gold-hover': '#D97706',
        },
        text: {
          primary: '#FFFFFF',
          muted: '#A0A0B0',
        },
        status: {
          success: '#00C853',
          warning: '#FFB300',
          error: '#FF1744',
        },
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
        display: ['Bebas Neue', 'Impact', 'sans-serif'],
      },
      animation: {
        'pulse-red': 'pulse-red 2s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        'fade-in': 'fadeIn 0.3s ease-in-out',
        'slide-up': 'slideUp 0.4s ease-out',
        'float': 'float 6s ease-in-out infinite',
        'shimmer': 'shimmer 2s linear infinite',
      },
      keyframes: {
        'pulse-red': {
          '0%, 100%': { opacity: 1 },
          '50%': { opacity: 0.5 },
        },
        fadeIn: {
          from: { opacity: 0 },
          to: { opacity: 1 },
        },
        slideUp: {
          from: { transform: 'translateY(20px)', opacity: 0 },
          to: { transform: 'translateY(0)', opacity: 1 },
        },
        float: {
          '0%, 100%': { transform: 'translateY(0px)' },
          '50%': { transform: 'translateY(-12px)' },
        },
        shimmer: {
          from: { backgroundPosition: '-200% 0' },
          to: { backgroundPosition: '200% 0' },
        },
      },
      aspectRatio: {
        video: '16 / 9',
      },
    },
  },
  plugins: [
    // Hide scrollbar while keeping scroll functionality
    function ({ addUtilities }) {
      addUtilities({
        '.scrollbar-none': {
          '-ms-overflow-style': 'none',
          'scrollbar-width': 'none',
          '&::-webkit-scrollbar': { display: 'none' },
        },
        '.gradient-card': {
          background: 'linear-gradient(135deg, #141418 0%, #1E1E26 100%)',
        },
        '.text-gradient': {
          background: 'linear-gradient(90deg, #E8002D, #FF6B00)',
          '-webkit-background-clip': 'text',
          '-webkit-text-fill-color': 'transparent',
          'background-clip': 'text',
        },
      })
    },
  ],
}
