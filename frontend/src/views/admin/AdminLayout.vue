<template>
  <div class="min-h-screen bg-gray-50 flex">
    <aside class="w-56 bg-white border-r border-gray-100 flex-shrink-0 hidden md:block">
      <div class="p-4 border-b border-gray-100">
        <router-link to="/admin" class="text-lg font-semibold text-dark">管理后台</router-link>
      </div>
      <nav class="p-3">
        <router-link v-for="item in menuItems" :key="item.path" :to="item.path" class="flex items-center gap-3 px-3 py-2 text-sm rounded-lg mb-1 transition-colors" :class="$route.path === item.path ? 'bg-primary/10 text-primary font-medium' : 'text-gray-600 hover:bg-gray-50'">{{ item.label }}</router-link>
      </nav>
      <div class="p-3 border-t border-gray-100 mt-auto">
        <router-link to="/" class="flex items-center gap-2 px-3 py-2 text-sm text-gray-500 hover:text-dark rounded-lg transition-colors">← 返回前台</router-link>
      </div>
    </aside>
    <main class="flex-1 overflow-x-auto">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
const menuItems = [
  { path: '/admin', label: '📊 仪表盘' },
  { path: '/admin/users', label: '👥 用户管理' },
  { path: '/admin/articles', label: '📝 文章管理' },
  { path: '/admin/prompts', label: '🔧 Prompt 管理' }
]
</script>
