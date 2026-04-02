import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export type ThemeMode = 'light' | 'dark'

export const useAppStore = defineStore('app', () => {
  // State
  const sidebarCollapsed = ref(false)
  const theme = ref<ThemeMode>((localStorage.getItem('theme') as ThemeMode) || 'light')
  const breadcrumbs = ref<Array<{ title: string; path?: string }>>([])

  // Getters
  const isDark = computed(() => theme.value === 'dark')

  // Actions
  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  function setTheme(mode: ThemeMode) {
    theme.value = mode
    localStorage.setItem('theme', mode)
    document.documentElement.setAttribute('data-theme', mode)
    if (mode === 'dark') {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }

  function toggleTheme() {
    setTheme(theme.value === 'light' ? 'dark' : 'light')
  }

  function setBreadcrumbs(items: Array<{ title: string; path?: string }>) {
    breadcrumbs.value = items
  }

  // 初始化主题
  function initTheme() {
    const saved = localStorage.getItem('theme') as ThemeMode | null
    if (saved) {
      setTheme(saved)
    } else if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
      setTheme('dark')
    }
  }

  return {
    sidebarCollapsed,
    theme,
    breadcrumbs,
    isDark,
    toggleSidebar,
    setTheme,
    toggleTheme,
    setBreadcrumbs,
    initTheme
  }
})
