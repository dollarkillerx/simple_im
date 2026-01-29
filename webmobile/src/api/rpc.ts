import type { ApiResponse } from '@/types'

let requestId = 0

function getToken(): string | null {
  return localStorage.getItem('token')
}

export async function rpcCall<T>(method: string, params?: Record<string, unknown>): Promise<T> {
  const token = getToken()
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  }

  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  const response = await fetch('/api/rpc', {
    method: 'POST',
    headers,
    body: JSON.stringify({
      jsonrpc: '2.0',
      id: String(++requestId),
      method,
      params: params || {},
    }),
  })

  if (!response.ok) {
    throw new Error(`HTTP error: ${response.status}`)
  }

  const data: ApiResponse<T> = await response.json()

  if (data.error) {
    throw new Error(data.error.message)
  }

  return data.result as T
}

// Auth API
export const authApi = {
  login: (username: string, password: string) =>
    rpcCall<{ user: { id: number; username: string; nickname: string; avatar: string }; token: string }>('user.login', { username, password }),

  register: (username: string, password: string, nickname?: string) =>
    rpcCall<{ user: { id: number; username: string; nickname: string; avatar: string }; token: string }>('user.register', { username, password, nickname }),

  getUserInfo: (userId?: number) =>
    rpcCall<{ id: number; username: string; nickname: string; avatar: string }>('user.info', userId ? { user_id: userId } : {}),
}

// Friend API
export const friendApi = {
  getList: () =>
    rpcCall<Array<{ id: number; user_id: number; username: string; nickname: string; avatar: string; created_at: string }>>('friend.list'),

  add: (params: { friend_id?: number; username?: string }) =>
    rpcCall<{ id: number; message: string }>('friend.add', params),

  accept: (requestId: number, accept: boolean) =>
    rpcCall<{ message: string }>('friend.accept', { request_id: requestId, accept }),

  getPending: () =>
    rpcCall<Array<{ id: number; user_id: number; username: string; nickname: string; avatar: string; created_at: string }>>('friend.pending'),
}

// Group API
export const groupApi = {
  getList: () =>
    rpcCall<Array<{ id: number; name: string; avatar: string; owner_id: number; owner_name: string; role: number; joined_at: string }>>('group.list'),

  create: (name: string, memberIds?: number[], avatar?: string) =>
    rpcCall<{ id: number; name: string; avatar: string; owner_id: number }>('group.create', { name, member_ids: memberIds, avatar }),

  getInfo: (groupId: number) =>
    rpcCall<{
      id: number
      name: string
      avatar: string
      owner_id: number
      owner_name: string
      created_at: string
      members: Array<{ user_id: number; username: string; nickname: string; avatar: string; role: number; joined_at: string }>
    }>('group.info', { group_id: groupId }),

  join: (groupId: number) =>
    rpcCall<{ message: string }>('group.join', { group_id: groupId }),
}

// Message API
export const messageApi = {
  send: (params: {
    receiver_id?: number
    group_id?: number
    msg_type?: number
    content?: string
    file_url?: string
    file_name?: string
    file_size?: number
  }) => rpcCall<{ id: number; created_at: string }>('message.send', params),

  getHistory: (params: { receiver_id?: number; group_id?: number; before_id?: number; limit?: number }) =>
    rpcCall<Array<{
      id: number
      sender_id: number
      receiver_id: number
      group_id: number
      msg_type: number
      content: string
      file_url: string
      file_name: string
      file_size: number
      created_at: string
      Sender?: { id: number; username: string; nickname: string; avatar: string }
    }>>('message.history', params),
}
