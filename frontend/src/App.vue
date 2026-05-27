<template>
  <n-config-provider :theme="naiveTheme" :theme-overrides="themeOverrides">
    <n-message-provider>
      <n-dialog-provider>
        <n-notification-provider>
          <n-loading-bar-provider>
            <div class="min-h-screen paper-texture">
              <Navbar v-if="showNavbar" />
              <main class="min-h-screen">
                <router-view v-slot="{ Component }">
                  <transition name="fade" mode="out-in">
                    <component :is="Component" />
                  </transition>
                </router-view>
              </main>
            </div>
          </n-loading-bar-provider>
        </n-notification-provider>
      </n-dialog-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { darkTheme, type GlobalThemeOverrides } from 'naive-ui'
import Navbar from '@/components/Navbar.vue'
import { useThemeStore } from '@/stores/theme'

const route = useRoute()
const themeStore = useThemeStore()
const showNavbar = computed(() => route.name !== 'Login' && route.name !== 'Register' && !route.path.startsWith('/admin'))

const themeOverrides: GlobalThemeOverrides = {
  common: {
    primaryColor: '#EB9463',
    primaryColorHover: '#F0A87E',
    primaryColorPressed: '#D67E4E',
    primaryColorSuppl: '#EB9463',
    borderRadius: '12px',
    fontFamily: '-apple-system, BlinkMacSystemFont, "SF Pro Display", "SF Pro Text", "Helvetica Neue", "PingFang SC", "Microsoft YaHei", sans-serif'
  },
  Card: {
    borderRadius: '16px',
    borderColor: 'rgba(0, 0, 0, 0.06)',
    color: 'rgba(255, 255, 255, 0.65)',
    boxShadow: '0 2px 8px rgba(0, 0, 0, 0.04), 0 8px 32px rgba(0, 0, 0, 0.06)'
  },
  Button: {
    borderRadius: '12px',
    fontWeight: '500'
  },
  Input: {
    borderRadius: '12px',
    color: 'rgba(255, 255, 255, 0.6)',
    border: '1px solid rgba(255, 255, 255, 0.5)',
    boxShadow: '0 2px 8px rgba(0, 0, 0, 0.02)'
  },
  Tag: {
    borderRadius: '8px'
  },
  Avatar: {
    boxShadow: '0 1px 3px rgba(0, 0, 0, 0.1)'
  }
}

const naiveTheme = computed(() => {
  return themeStore.isDark ? darkTheme : null
})
</script>

<style>
/* 全局过渡动画 */
.page-enter-active,
.page-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.page-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.page-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
