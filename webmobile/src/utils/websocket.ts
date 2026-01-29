import type { WSMessage } from '@/types'

type MessageHandler = (message: WSMessage) => void

class WebSocketClient {
  private ws: WebSocket | null = null
  private url: string = ''
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null
  private heartbeatTimer: ReturnType<typeof setInterval> | null = null
  private handlers: Set<MessageHandler> = new Set()
  private isManualClose = false
  private reconnectAttempts = 0
  private maxReconnectAttempts = 10
  private reconnectDelay = 1000

  connect(token: string): void {
    if (this.ws?.readyState === WebSocket.OPEN) {
      return
    }

    this.isManualClose = false
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const host = window.location.host
    this.url = `${protocol}//${host}/ws?token=${token}`

    this.createConnection()
  }

  private createConnection(): void {
    try {
      this.ws = new WebSocket(this.url)

      this.ws.onopen = () => {
        console.log('[WS] Connected')
        this.reconnectAttempts = 0
        this.startHeartbeat()
      }

      this.ws.onmessage = (event: MessageEvent) => {
        try {
          const message: WSMessage = JSON.parse(event.data)
          this.handlers.forEach((handler) => handler(message))
        } catch (e) {
          console.error('[WS] Parse message error:', e)
        }
      }

      this.ws.onclose = (event: CloseEvent) => {
        console.log('[WS] Disconnected:', event.code, event.reason)
        this.stopHeartbeat()

        if (!this.isManualClose && this.reconnectAttempts < this.maxReconnectAttempts) {
          this.scheduleReconnect()
        }
      }

      this.ws.onerror = (error: Event) => {
        console.error('[WS] Error:', error)
      }
    } catch (e) {
      console.error('[WS] Create connection error:', e)
      this.scheduleReconnect()
    }
  }

  private scheduleReconnect(): void {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer)
    }

    const delay = Math.min(this.reconnectDelay * Math.pow(2, this.reconnectAttempts), 30000)
    this.reconnectAttempts++

    console.log(`[WS] Reconnecting in ${delay}ms (attempt ${this.reconnectAttempts})`)

    this.reconnectTimer = setTimeout(() => {
      this.createConnection()
    }, delay)
  }

  private startHeartbeat(): void {
    this.stopHeartbeat()
    this.heartbeatTimer = setInterval(() => {
      if (this.ws?.readyState === WebSocket.OPEN) {
        this.ws.send(JSON.stringify({ type: 'ping' }))
      }
    }, 30000)
  }

  private stopHeartbeat(): void {
    if (this.heartbeatTimer) {
      clearInterval(this.heartbeatTimer)
      this.heartbeatTimer = null
    }
  }

  disconnect(): void {
    this.isManualClose = true

    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer)
      this.reconnectTimer = null
    }

    this.stopHeartbeat()

    if (this.ws) {
      this.ws.close()
      this.ws = null
    }

    this.handlers.clear()
  }

  send(data: Record<string, unknown>): void {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(data))
    } else {
      console.warn('[WS] Not connected, cannot send message')
    }
  }

  onMessage(handler: MessageHandler): () => void {
    this.handlers.add(handler)
    return () => {
      this.handlers.delete(handler)
    }
  }

  get isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN
  }
}

export const wsClient = new WebSocketClient()
