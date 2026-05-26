<template>
  <n-config-provider :theme="naiveTheme" :theme-overrides="themeOverrides">
    <n-message-provider>
      <n-dialog-provider>
        <n-notification-provider>
          <n-loading-bar-provider>
            <Navbar v-if="showNavbar" />
            <main class="min-h-screen bg-gray-50">
              <router-view v-slot="{ Component }">
                <transition name="fade" mode="out-in">
                  <component :is="Component" />
                </transition>
              </router-view>
            </main>
          </n-loading-bar-provider>
        </n-notification-provider>
      </n-dialog-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { darkTheme } from 'naive-ui'
import type { GlobalThemeOverrides } from 'naive-ui'
import Navbar from '@/components/Navbar.vue'

const route = useRoute()
const showNavbar = computed(() => route.name !== 'Login' && route.name !== 'Register')

const themeOverrides: GlobalThemeOverrides = {
  common: {
    primaryColor: '#6c5ce7',
    primaryColorHover: '#7d73ef',
    primaryColorPressed: '#5a4bd1',
    primaryColorSuppl: '#6c5ce7',
    borderRadius: '8px'
  },
  Card: {
    borderRadius: '12px',
    borderColor: '#f0f0f0'
  },
  Button: {
    borderRadius: '8px'
  },
  Input: {
    borderRadius: '8px'
  },
  Tag: {
    borderRadius: '6px'
  }
}

const naiveTheme = computed(() => {
  return null
})
</script>
