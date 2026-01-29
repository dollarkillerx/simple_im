// User types
export interface User {
  id: number
  username: string
  nickname: string
  avatar: string
  status?: number
  created_at?: string
}

// Friend types
export interface Friend {
  id: number
  user_id: number
  username: string
  nickname: string
  avatar: string
  created_at: string
}

export interface FriendRequest {
  id: number
  user_id: number
  username: string
  nickname: string
  avatar: string
  created_at: string
}

// Group types
export interface Group {
  id: number
  name: string
  avatar: string
  owner_id: number
  owner_name: string
  role: number
  joined_at: string
}

export interface GroupMember {
  user_id: number
  username: string
  nickname: string
  avatar: string
  role: number
  joined_at: string
}

export interface GroupInfo {
  id: number
  name: string
  avatar: string
  owner_id: number
  owner_name: string
  created_at: string
  members: GroupMember[]
}

// Message types
export const MessageType = {
  Text: 1,
  Image: 2,
  File: 3,
} as const

export type MessageTypeValue = (typeof MessageType)[keyof typeof MessageType]

export interface Message {
  id: number
  sender_id: number
  receiver_id: number
  group_id: number
  msg_type: MessageTypeValue
  content: string
  file_url: string
  file_name: string
  file_size: number
  created_at: string
  Sender?: User
}

// WebSocket message types
export interface WSMessage {
  id: number
  type: string
  sender_id: number
  sender_name: string
  receiver_id?: number
  group_id?: number
  group_name?: string
  msg_type: MessageTypeValue
  content?: string
  file_url?: string
  file_name?: string
  file_size?: number
  created_at: string
}

// Conversation types
export interface Conversation {
  id: string // 'user_123' or 'group_456'
  type: 'private' | 'group'
  targetId: number
  name: string
  avatar: string
  lastMessage: string
  lastTime: string
  unread: number
}

// API Response types
export interface ApiResponse<T = unknown> {
  jsonrpc: string
  id: number
  result?: T
  error?: {
    code: number
    message: string
  }
}

export interface LoginResult {
  user: User
  token: string
}
