import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export type ThemeMode = 'light' | 'dark' | 'auto'

export const useAppStore = defineStore('app', () => {
  // State
  const sidebarCollapsed = ref(false)
  const theme = ref<ThemeMode>((localStorage.getItem('theme-mode') as ThemeMode) || 'auto')
  const breadcrumbs = ref<Array<{ title: string; path?: string }>>([])

  // Getters
  const isDark = computed(() => {
    if (theme.value === 'auto') {
      return window.matchMedia('(prefers-color-scheme: dark)').matches
    }
    return theme.value === 'dark'
  })

  // Actions
  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  function setTheme(mode: ThemeMode) {
    theme.value = mode
    localStorage.setItem('theme-mode', mode)
    applyTheme()
  }

  function toggleTheme() {
    const cycle: ThemeMode[] = ['light', 'dark', 'auto']
    const idx = cycle.indexOf(theme.value)
    setTheme(cycle[(idx + 1) % cycle.length])
  }

  function setBreadcrumbs(items: Array<{ title: string; path?: string }>) {
    breadcrumbs.value = items
  }

  function applyTheme() {
    const effectiveTheme = isDark.value ? 'dark' : 'light'
    document.documentElement.setAttribute('data-theme', effectiveTheme)
    if (effectiveTheme === 'dark') {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }

  // 初始化主题
  function initTheme() {
    applyTheme()
    // 监听系统主题变化
    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
      if (theme.value === 'auto') {
        applyTheme()
      }
    })
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
