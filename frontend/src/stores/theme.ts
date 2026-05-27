import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export type ThemeMode = 'light' | 'dark' | 'system'

export const useThemeStore = defineStore('theme', () => {
  const mode = ref<ThemeMode>((localStorage.getItem('theme') as ThemeMode) || 'system')
  const isDark = ref(false)

  function applyTheme() {
    let dark = false
    if (mode.value === 'dark') {
      dark = true
    } else if (mode.value === 'system') {
      dark = window.matchMedia('(prefers-color-scheme: dark)').matches
    }
    isDark.value = dark

    if (dark) {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }

  function setMode(newMode: ThemeMode) {
    mode.value = newMode
    localStorage.setItem('theme', newMode)
    applyTheme()
  }

  // Listen for system theme changes
  const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
  mediaQuery.addEventListener('change', () => {
    if (mode.value === 'system') {
      applyTheme()
    }
  })

  // Initialize
  applyTheme()

  return { mode, isDark, setMode }
})
