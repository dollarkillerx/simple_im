import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Conversation, Message, WSMessage, Friend, Group, MessageTypeValue } from '@/types'
import { MessageType } from '@/types'
import { friendApi, groupApi, messageApi } from '@/api/rpc'
import { wsClient } from '@/utils/websocket'
import { useUserStore } from './user'

export const useChatStore = defineStore('chat', () => {
  const friends = ref<Friend[]>([])
  const groups = ref<Group[]>([])
  const conversations = ref<Map<string, Conversation>>(new Map())
  const messages = ref<Map<string, Message[]>>(new Map())
  const currentChatId = ref<string | null>(null)
  const pendingRequests = ref<Array<{ id: number; user_id: number; username: string; nickname: string; avatar: string; created_at: string }>>([])

  const sortedConversations = computed(() => {
    return Array.from(conversations.value.values()).sort((a, b) => {
      return new Date(b.lastTime).getTime() - new Date(a.lastTime).getTime()
    })
  })

  const currentMessages = computed(() => {
    if (!currentChatId.value) return []
    return messages.value.get(currentChatId.value) || []
  })

  const totalUnread = computed(() => {
    let count = 0
    conversations.value.forEach((conv) => {
      count += conv.unread
    })
    return count
  })

  function getChatId(type: 'private' | 'group', targetId: number): string {
    return `${type}_${targetId}`
  }

  async function loadFriends() {
    try {
      const list = await friendApi.getList()
      friends.value = list

      // Update or create conversations for friends
      list.forEach((friend) => {
        const chatId = getChatId('private', friend.user_id)
        if (!conversations.value.has(chatId)) {
          conversations.value.set(chatId, {
            id: chatId,
            type: 'private',
            targetId: friend.user_id,
            name: friend.nickname || friend.username,
            avatar: friend.avatar,
            lastMessage: '',
            lastTime: friend.created_at,
            unread: 0,
          })
        }
      })
    } catch (e) {
      console.error('Load friends failed:', e)
    }
  }

  async function loadGroups() {
    try {
      const list = await groupApi.getList()
      groups.value = list

      // Update or create conversations for groups
      list.forEach((group) => {
        const chatId = getChatId('group', group.id)
        if (!conversations.value.has(chatId)) {
          conversations.value.set(chatId, {
            id: chatId,
            type: 'group',
            targetId: group.id,
            name: group.name,
            avatar: group.avatar,
            lastMessage: '',
            lastTime: group.joined_at,
            unread: 0,
          })
        }
      })
    } catch (e) {
      console.error('Load groups failed:', e)
    }
  }

  async function loadPendingRequests() {
    try {
      pendingRequests.value = await friendApi.getPending()
    } catch (e) {
      console.error('Load pending requests failed:', e)
    }
  }

  async function loadMessages(chatId: string, beforeId?: number) {
    const conv = conversations.value.get(chatId)
    if (!conv) return []

    try {
      const params: { receiver_id?: number; group_id?: number; before_id?: number; limit?: number } = {
        limit: 50,
      }

      if (conv.type === 'private') {
        params.receiver_id = conv.targetId
      } else {
        params.group_id = conv.targetId
      }

      if (beforeId) {
        params.before_id = beforeId
      }

      const list = await messageApi.getHistory(params)

      // Map list to Message type
      const mappedList: Message[] = list.map((item) => ({
        ...item,
        msg_type: item.msg_type as MessageTypeValue,
        Sender: item.Sender
          ? {
              id: item.Sender.id,
              username: item.Sender.username,
              nickname: item.Sender.nickname,
              avatar: item.Sender.avatar,
            }
          : undefined,
      }))

      if (!beforeId) {
        messages.value.set(chatId, mappedList)
      } else {
        const existing = messages.value.get(chatId) || []
        messages.value.set(chatId, [...mappedList, ...existing])
      }

      return list
    } catch (e) {
      console.error('Load messages failed:', e)
      return []
    }
  }

  async function sendMessage(params: {
    receiverId?: number
    groupId?: number
    msgType?: MessageTypeValue
    content?: string
    fileUrl?: string
    fileName?: string
    fileSize?: number
  }) {
    const userStore = useUserStore()
    if (!userStore.userInfo) return

    const chatId = params.groupId
      ? getChatId('group', params.groupId)
      : getChatId('private', params.receiverId!)

    const result = await messageApi.send({
      receiver_id: params.receiverId,
      group_id: params.groupId,
      msg_type: params.msgType || MessageType.Text,
      content: params.content,
      file_url: params.fileUrl,
      file_name: params.fileName,
      file_size: params.fileSize,
    })

    // Add message locally
    const newMessage: Message = {
      id: result.id,
      sender_id: userStore.userInfo.id,
      receiver_id: params.receiverId || 0,
      group_id: params.groupId || 0,
      msg_type: params.msgType || MessageType.Text,
      content: params.content || '',
      file_url: params.fileUrl || '',
      file_name: params.fileName || '',
      file_size: params.fileSize || 0,
      created_at: result.created_at,
      Sender: userStore.userInfo,
    }

    const existing = messages.value.get(chatId) || []
    messages.value.set(chatId, [...existing, newMessage])

    // Update conversation
    const conv = conversations.value.get(chatId)
    if (conv) {
      conv.lastMessage = params.content || '[文件]'
      conv.lastTime = result.created_at
    }

    return result
  }

  function handleWSMessage(msg: WSMessage) {
    const userStore = useUserStore()
    if (!userStore.userInfo) return

    if (msg.type === 'message') {
      const isSelf = msg.sender_id === userStore.userInfo.id
      const chatId = msg.group_id
        ? getChatId('group', msg.group_id)
        : isSelf
          ? getChatId('private', msg.receiver_id!)
          : getChatId('private', msg.sender_id)

      // Don't add if it's our own message (already added when sending)
      if (isSelf) return

      const message: Message = {
        id: msg.id,
        sender_id: msg.sender_id,
        receiver_id: msg.receiver_id || 0,
        group_id: msg.group_id || 0,
        msg_type: msg.msg_type,
        content: msg.content || '',
        file_url: msg.file_url || '',
        file_name: msg.file_name || '',
        file_size: msg.file_size || 0,
        created_at: msg.created_at,
        Sender: {
          id: msg.sender_id,
          username: msg.sender_name,
          nickname: msg.sender_name,
          avatar: '',
        },
      }

      const existing = messages.value.get(chatId) || []
      messages.value.set(chatId, [...existing, message])

      // Update conversation
      let conv = conversations.value.get(chatId)
      if (!conv) {
        conv = {
          id: chatId,
          type: msg.group_id ? 'group' : 'private',
          targetId: msg.group_id || msg.sender_id,
          name: msg.group_name || msg.sender_name,
          avatar: '',
          lastMessage: '',
          lastTime: msg.created_at,
          unread: 0,
        }
        conversations.value.set(chatId, conv)
      }

      conv.lastMessage = msg.content || '[文件]'
      conv.lastTime = msg.created_at

      if (chatId !== currentChatId.value) {
        conv.unread++
      }
    } else if (msg.type === 'friend_request') {
      loadPendingRequests()
    } else if (msg.type === 'friend_accepted') {
      loadFriends()
    }
  }

  function setCurrentChat(chatId: string | null) {
    currentChatId.value = chatId
    if (chatId) {
      const conv = conversations.value.get(chatId)
      if (conv) {
        conv.unread = 0
      }
    }
  }

  function initWSListener() {
    wsClient.onMessage(handleWSMessage)
  }

  return {
    friends,
    groups,
    conversations,
    messages,
    currentChatId,
    pendingRequests,
    sortedConversations,
    currentMessages,
    totalUnread,
    getChatId,
    loadFriends,
    loadGroups,
    loadPendingRequests,
    loadMessages,
    sendMessage,
    setCurrentChat,
    initWSListener,
  }
})
