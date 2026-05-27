/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,ts,tsx}'],
  theme: {
    extend: {
      fontFamily: {
        sans: ['-apple-system', 'BlinkMacSystemFont', 'SF Pro Display', 'SF Pro Text', 'Helvetica Neue', 'PingFang SC', 'Microsoft YaHei', 'sans-serif'],
        mono: ['SF Mono', 'Fira Code', 'Consolas', 'monospace']
      },
      colors: {
        primary: {
          DEFAULT: '#EB9463',
          50: '#fef5ef',
          100: '#fce6d5',
          200: '#f9cfaa',
          300: '#f4b179',
          400: '#f0a87e',
          500: '#EB9463',
          600: '#D67E4E',
          700: '#b86a3e',
          800: '#9a5734',
          900: '#7d462b'
        },
        dark: '#1d1d1f',
        paper: '#fafaf8'
      },
      backdropBlur: {
        xs: '2px',
        sm: '8px',
        md: '12px',
        lg: '16px',
        xl: '24px',
        '2xl': '40px',
        '3xl': '64px'
      },
      borderRadius: {
        'glass': '12px',
        'card': '16px'
      },
      boxShadow: {
        'glass': '0 8px 32px rgba(0, 0, 0, 0.08)',
        'glass-lg': '0 12px 40px rgba(0, 0, 0, 0.12)',
        'paper': '0 1px 3px rgba(0, 0, 0, 0.05), 0 2px 8px rgba(0, 0, 0, 0.04)',
        'paper-hover': '0 2px 6px rgba(0, 0, 0, 0.08), 0 4px 16px rgba(0, 0, 0, 0.06)'
      },
      animation: {
        'float': 'float 6s ease-in-out infinite',
        'pulse-slow': 'pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite'
      },
      keyframes: {
        float: {
          '0%, 100%': { transform: 'translateY(0px)' },
          '50%': { transform: 'translateY(-10px)' }
        }
      }
    }
  },
  plugins: []
}
