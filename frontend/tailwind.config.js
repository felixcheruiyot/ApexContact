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
        },
        accent: {
          red: '#E8002D',
          'red-hover': '#C0001F',
          orange: '#FF6B00',
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
      },
      aspectRatio: {
        video: '16 / 9',
      },
    },
  },
  plugins: [],
}
