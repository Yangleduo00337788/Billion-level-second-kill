/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,ts,tsx}'],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Inter', 'system-ui', '-apple-system', 'sans-serif']
      },
      colors: {
        primary: {
          DEFAULT: '#6c5ce7',
          50: '#f0effd',
          100: '#d5d0fa',
          200: '#b8b1f7',
          300: '#9b92f3',
          400: '#7d73ef',
          500: '#6c5ce7',
          600: '#5a4bd1',
          700: '#4a3dbb',
          800: '#3b30a5',
          900: '#2b248f'
        },
        dark: '#1a1a2e'
      }
    }
  },
  plugins: []
}
