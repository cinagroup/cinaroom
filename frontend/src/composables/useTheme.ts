/**
 * 主题系统 Composable - 参照 CinaToken ThemeContext
 * 三种模式：light / dark / auto（跟随系统 prefers-color-scheme）
 * localStorage 持久化，html 根元素 data-theme 属性驱动
 */
import { ref, computed, watchEffect, onMounted } from 'vue'

export type ThemeMode = 'light' | 'dark' | 'auto'

const STORAGE_KEY = 'theme-mode'

const mode = ref<ThemeMode>((localStorage.getItem(STORAGE_KEY) as ThemeMode) || 'auto')
const systemDark = ref(false)

// 获取系统主题偏好
function getSystemPrefersDark(): boolean {
  return window.matchMedia('(prefers-color-scheme: dark)').matches
}

export function useTheme() {
  const actualTheme = computed(() => {
    if (mode.value === 'auto') {
      return systemDark.value ? 'dark' : 'light'
    }
    return mode.value
  })

  const isDark = computed(() => actualTheme.value === 'dark')

  function setTheme(newMode: ThemeMode) {
    mode.value = newMode
    localStorage.setItem(STORAGE_KEY, newMode)
  }

  function toggleTheme() {
    // light → dark → auto → light
    const cycle: ThemeMode[] = ['light', 'dark', 'auto']
    const idx = cycle.indexOf(mode.value)
    setTheme(cycle[(idx + 1) % cycle.length])
  }

  // 应用主题到 DOM
  function applyTheme() {
    const theme = actualTheme.value
    document.documentElement.setAttribute('data-theme', theme)
    if (theme === 'dark') {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }

  // 监听系统主题变化
  onMounted(() => {
    systemDark.value = getSystemPrefersDark()
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
    const handler = (e: MediaQueryListEvent) => {
      systemDark.value = e.matches
    }
    mediaQuery.addEventListener('change', handler)
  })

  // 自动应用
  watchEffect(() => {
    applyTheme()
  })

  return {
    mode,
    actualTheme,
    isDark,
    setTheme,
    toggleTheme
  }
}
