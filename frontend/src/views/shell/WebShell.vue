<template>
  <div class="webshell-container">
    <div class="terminal-header">
      <div class="terminal-title">
        <el-icon><Terminal /></el-icon>
        <span>WebShell - {{ vmName }}</span>
        <el-tag :type="connected ? 'success' : 'danger'" style="margin-left: 10px">
          {{ connected ? '已连接' : '未连接' }}
        </el-tag>
      </div>
      <div class="terminal-toolbar">
        <el-button size="small" @click="handleDisconnect">断开连接</el-button>
        <el-button size="small" @click="handleClear">清屏</el-button>
        <el-button size="small" @click="handleFullscreen">全屏</el-button>
        <el-dropdown size="small">
          <el-button size="small">
            字体 <el-icon><ArrowDown /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="changeFontSize(1)">增大</el-dropdown-item>
              <el-dropdown-item @click="changeFontSize(-1)">减小</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-button size="small" @click="toggleTheme">
          {{ darkTheme ? '亮色' : '暗色' }}
        </el-button>
      </div>
    </div>
    
    <div ref="terminalContainer" class="terminal-container"></div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import 'xterm/css/xterm.css'

const route = useRoute()
const terminalContainer = ref(null)

const vmName = ref('ubuntu-dev')
const connected = ref(true)
const darkTheme = ref(true)
let fontSize = 14

let term = null
let fitAddon = null

const initTerminal = () => {
  if (!terminalContainer.value) return
  
  term = new Terminal({
    fontSize: fontSize,
    fontFamily: 'Menlo, Monaco, "Courier New", monospace',
    theme: darkTheme.value ? {
      background: '#1e1e1e',
      foreground: '#ffffff'
    } : {
      background: '#ffffff',
      foreground: '#000000'
    },
    cursorBlink: true,
    scrollback: 10000
  })
  
  fitAddon = new FitAddon()
  term.loadAddon(fitAddon)
  term.open(terminalContainer.value)
  fitAddon.fit()
  
  term.writeln('\x1b[1;32m欢迎使用 CinaSeek WebShell\x1b[0m')
  term.writeln('\x1b[1;32m正在连接到虚拟机...\x1b[0m\n')
  
  // TODO: 建立 WebSocket 连接
  setTimeout(() => {
    term.writeln('\x1b[1;32m连接成功！\x1b[0m\n')
    term.writeln('\x1b[1;34mubuntu@ubuntu-dev:~$ \x1b[0m')
  }, 1000)
  
  term.onData((data) => {
    // TODO: 发送数据到后端
    console.log('发送数据:', data)
  })
}

const handleDisconnect = () => {
  // TODO: 断开 WebSocket 连接
  connected.value = false
  ElMessage.warning('已断开连接')
}

const handleClear = () => {
  term?.clear()
}

const handleFullscreen = () => {
  document.documentElement.requestFullscreen?.()
}

const changeFontSize = (delta) => {
  fontSize = Math.max(10, Math.min(24, fontSize + delta))
  term.options.fontSize = fontSize
  fitAddon?.fit()
}

const toggleTheme = () => {
  darkTheme.value = !darkTheme.value
  term.options.theme = darkTheme.value ? {
    background: '#1e1e1e',
    foreground: '#ffffff'
  } : {
    background: '#ffffff',
    foreground: '#000000'
  }
}

const handleResize = () => {
  fitAddon?.fit()
}

onMounted(() => {
  initTerminal()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  term?.dispose()
})
</script>

<style scoped lang="scss">
.webshell-container {
  height: calc(100vh - 120px);
  display: flex;
  flex-direction: column;
  background-color: #1e1e1e;
  border-radius: $border-radius;
  overflow: hidden;
}

.terminal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 15px;
  background-color: #2d2d2d;
  border-bottom: 1px solid #404040;
  
  .terminal-title {
    display: flex;
    align-items: center;
    color: #fff;
    font-weight: bold;
  }
  
  .terminal-toolbar {
    display: flex;
    gap: 8px;
  }
}

.terminal-container {
  flex: 1;
  padding: 10px;
  overflow: hidden;
}
</style>
