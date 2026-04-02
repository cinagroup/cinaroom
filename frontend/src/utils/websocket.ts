/**
 * WebSocket 连接管理器
 * 支持心跳、自动重连、消息队列
 */

export interface WSOptions {
  /** WebSocket URL */
  url: string
  /** 子协议 */
  protocols?: string | string[]
  /** 心跳间隔（毫秒），默认 30000 */
  heartbeatInterval?: number
  /** 重连间隔（毫秒），默认 3000 */
  reconnectInterval?: number
  /** 最大重连次数，默认 10 */
  maxReconnectAttempts?: number
  /** 连接超时（毫秒），默认 10000 */
  connectTimeout?: number
  /** 收到消息回调 */
  onMessage?: (data: unknown) => void
  /** 连接打开回调 */
  onOpen?: (event: Event) => void
  /** 连接关闭回调 */
  onClose?: (event: CloseEvent) => void
  /** 连接错误回调 */
  onError?: (event: Event) => void
  /** 重连中回调 */
  onReconnecting?: (attempt: number) => void
  /** 重连失败回调 */
  onReconnectFailed?: () => void
}

export class WebSocketManager {
  private ws: WebSocket | null = null
  private options: Required<Pick<WSOptions, 'heartbeatInterval' | 'reconnectInterval' | 'maxReconnectAttempts' | 'connectTimeout'>> & Omit<WSOptions, 'heartbeatInterval' | 'reconnectInterval' | 'maxReconnectAttempts' | 'connectTimeout'>
  private reconnectAttempts = 0
  private heartbeatTimer: ReturnType<typeof setInterval> | null = null
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null
  private connectTimer: ReturnType<typeof setTimeout> | null = null
  private messageQueue: unknown[] = []
  private _isConnected = false
  private _isClosed = false

  constructor(options: WSOptions) {
    this.options = {
      heartbeatInterval: 30000,
      reconnectInterval: 3000,
      maxReconnectAttempts: 10,
      connectTimeout: 10000,
      ...options
    }
  }

  get isConnected(): boolean {
    return this._isConnected
  }

  /** 建立 WebSocket 连接 */
  connect(): void {
    if (this.ws?.readyState === WebSocket.OPEN) return
    this._isClosed = false

    try {
      this.ws = new WebSocket(this.options.url, this.options.protocols)
    } catch (err) {
      console.error('[WS] 创建连接失败:', err)
      this.scheduleReconnect()
      return
    }

    // 连接超时
    this.connectTimer = setTimeout(() => {
      if (this.ws?.readyState !== WebSocket.OPEN) {
        console.warn('[WS] 连接超时')
        this.ws?.close()
      }
    }, this.options.connectTimeout)

    this.ws.onopen = (event) => {
      this.clearTimers()
      this._isConnected = true
      this.reconnectAttempts = 0
      console.info('[WS] 连接已建立')

      // 发送队列中的消息
      while (this.messageQueue.length > 0) {
        const msg = this.messageQueue.shift()
        this.doSend(msg)
      }

      this.startHeartbeat()
      this.options.onOpen?.(event)
    }

    this.ws.onmessage = (event) => {
      let data: unknown
      try {
        data = JSON.parse(event.data)
      } catch {
        data = event.data
      }
      this.options.onMessage?.(data)
    }

    this.ws.onclose = (event) => {
      this._isConnected = false
      this.stopHeartbeat()
      console.info('[WS] 连接已关闭', event.code, event.reason)
      this.options.onClose?.(event)

      if (!this._isClosed) {
        this.scheduleReconnect()
      }
    }

    this.ws.onerror = (event) => {
      console.error('[WS] 连接错误')
      this.options.onError?.(event)
    }
  }

  /** 发送消息 */
  send(data: unknown): void {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.doSend(data)
    } else {
      this.messageQueue.push(data)
    }
  }

  /** 关闭连接 */
  close(): void {
    this._isClosed = true
    this.clearTimers()
    this.stopHeartbeat()
    
    if (this.ws) {
      this.ws.onclose = null
      this.ws.onerror = null
      this.ws.onmessage = null
      this.ws.onopen = null
      this.ws.close(1000, 'Client closed')
      this.ws = null
    }
    this._isConnected = false
  }

  /** 重新连接 */
  reconnect(): void {
    this.close()
    this._isClosed = false
    this.reconnectAttempts = 0
    this.connect()
  }

  private doSend(data: unknown): void {
    if (!this.ws) return
    const payload = typeof data === 'string' ? data : JSON.stringify(data)
    this.ws.send(payload)
  }

  private startHeartbeat(): void {
    this.stopHeartbeat()
    this.heartbeatTimer = setInterval(() => {
      if (this.ws?.readyState === WebSocket.OPEN) {
        this.ws.send(JSON.stringify({ type: 'ping', timestamp: Date.now() }))
      }
    }, this.options.heartbeatInterval)
  }

  private stopHeartbeat(): void {
    if (this.heartbeatTimer) {
      clearInterval(this.heartbeatTimer)
      this.heartbeatTimer = null
    }
  }

  private scheduleReconnect(): void {
    if (this._isClosed) return
    
    this.reconnectAttempts++
    if (this.reconnectAttempts > this.options.maxReconnectAttempts) {
      console.error('[WS] 已达最大重连次数')
      this.options.onReconnectFailed?.()
      return
    }

    const delay = Math.min(
      this.options.reconnectInterval * Math.pow(1.5, this.reconnectAttempts - 1),
      30000 // 最大 30 秒
    )

    console.info(`[WS] ${delay}ms 后进行第 ${this.reconnectAttempts} 次重连...`)
    this.options.onReconnecting?.(this.reconnectAttempts)

    this.reconnectTimer = setTimeout(() => {
      this.connect()
    }, delay)
  }

  private clearTimers(): void {
    if (this.connectTimer) {
      clearTimeout(this.connectTimer)
      this.connectTimer = null
    }
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer)
      this.reconnectTimer = null
    }
  }
}

/**
 * 创建 WebSocket 终端连接 URL
 * 根据当前页面协议决定 ws/wss
 */
export function buildWsUrl(path: string, params?: Record<string, string>): string {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const host = window.location.host
  const base = `${protocol}//${host}${path}`
  
  if (params) {
    const query = new URLSearchParams(params).toString()
    return `${base}?${query}`
  }
  return base
}
