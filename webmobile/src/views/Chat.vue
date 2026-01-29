<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { showToast, showActionSheet } from 'vant'
import { useChatStore } from '@/stores/chat'
import { useUserStore } from '@/stores/user'
import { MessageType } from '@/types'
import { formatChatTime, getDefaultAvatar, formatFileSize } from '@/utils/format'

const props = defineProps<{
  type: 'private' | 'group'
  id: string
}>()

const router = useRouter()
const route = useRoute()
const chatStore = useChatStore()
const userStore = useUserStore()

const messageInput = ref('')
const messagesContainer = ref<HTMLElement | null>(null)
const loading = ref(false)
const sending = ref(false)

const chatId = computed(() => chatStore.getChatId(props.type, parseInt(props.id)))
const conversation = computed(() => chatStore.conversations.get(chatId.value))
const messages = computed(() => chatStore.messages.get(chatId.value) || [])

const chatTitle = computed(() => {
  return conversation.value?.name || (props.type === 'group' ? '群聊' : '私聊')
})

function getAvatar(name: string, avatar?: string) {
  return avatar || getDefaultAvatar(name)
}

function isSelf(senderId: number) {
  return senderId === userStore.userInfo?.id
}

function scrollToBottom() {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

async function loadMessages() {
  loading.value = true
  await chatStore.loadMessages(chatId.value)
  loading.value = false
  scrollToBottom()
}

async function sendMessage() {
  if (!messageInput.value.trim() || sending.value) return

  const content = messageInput.value.trim()
  messageInput.value = ''
  sending.value = true

  try {
    const params: { receiverId?: number; groupId?: number; content: string } = {
      content,
    }

    if (props.type === 'group') {
      params.groupId = parseInt(props.id)
    } else {
      params.receiverId = parseInt(props.id)
    }

    await chatStore.sendMessage(params)
    scrollToBottom()
  } catch (e) {
    showToast((e as Error).message || '发送失败')
  } finally {
    sending.value = false
  }
}

async function handleAttachment() {
  const actions = [
    { name: '图片', value: 'image' },
    { name: '文件', value: 'file' },
  ]

  try {
    const action = await showActionSheet({ actions })
    if (action.value === 'image') {
      // TODO: Implement image upload
      showToast('图片上传功能开发中')
    } else if (action.value === 'file') {
      // TODO: Implement file upload
      showToast('文件上传功能开发中')
    }
  } catch {
    // User cancelled
  }
}

function goBack() {
  router.back()
}

function goToInfo() {
  if (props.type === 'group') {
    router.push(`/group/${props.id}`)
  }
}

onMounted(() => {
  chatStore.setCurrentChat(chatId.value)
  loadMessages()
})

onUnmounted(() => {
  chatStore.setCurrentChat(null)
})

watch(
  () => chatStore.messages.get(chatId.value)?.length,
  () => {
    scrollToBottom()
  }
)
</script>

<template>
  <div class="chat-page">
    <van-nav-bar
      :title="chatTitle"
      left-arrow
      fixed
      placeholder
      @click-left="goBack"
    >
      <template #right>
        <van-icon
          v-if="type === 'group'"
          name="ellipsis"
          size="22"
          @click="goToInfo"
        />
      </template>
    </van-nav-bar>

    <div ref="messagesContainer" class="messages-container">
      <van-loading v-if="loading" class="loading-indicator" />

      <div class="messages-list">
        <div
          v-for="msg in messages"
          :key="msg.id"
          class="message-item"
          :class="{ 'message-self': isSelf(msg.sender_id) }"
        >
          <van-image
            v-if="!isSelf(msg.sender_id)"
            :src="getAvatar(msg.Sender?.nickname || msg.Sender?.username || 'U', msg.Sender?.avatar)"
            width="36"
            height="36"
            round
            fit="cover"
            class="message-avatar"
          />

          <div class="message-content">
            <div v-if="!isSelf(msg.sender_id) && type === 'group'" class="message-sender">
              {{ msg.Sender?.nickname || msg.Sender?.username }}
            </div>

            <div class="message-bubble">
              <!-- Text message -->
              <template v-if="msg.msg_type === MessageType.Text">
                {{ msg.content }}
              </template>

              <!-- Image message -->
              <template v-else-if="msg.msg_type === MessageType.Image">
                <van-image
                  :src="msg.file_url"
                  width="200"
                  fit="contain"
                  radius="8"
                  @click="() => {}"
                />
              </template>

              <!-- File message -->
              <template v-else-if="msg.msg_type === MessageType.File">
                <div class="file-message">
                  <van-icon name="description" size="32" />
                  <div class="file-info">
                    <span class="file-name">{{ msg.file_name }}</span>
                    <span class="file-size">{{ formatFileSize(msg.file_size) }}</span>
                  </div>
                </div>
              </template>
            </div>

            <div class="message-time">
              {{ formatChatTime(msg.created_at) }}
            </div>
          </div>

          <van-image
            v-if="isSelf(msg.sender_id)"
            :src="getAvatar(userStore.userInfo?.nickname || 'U', userStore.userInfo?.avatar)"
            width="36"
            height="36"
            round
            fit="cover"
            class="message-avatar"
          />
        </div>
      </div>

      <van-empty v-if="!loading && messages.length === 0" description="暂无消息" />
    </div>

    <div class="input-area">
      <van-icon
        name="add-o"
        size="24"
        class="input-action"
        @click="handleAttachment"
      />
      <van-field
        v-model="messageInput"
        class="message-input"
        placeholder="输入消息..."
        type="text"
        :border="false"
        @keyup.enter="sendMessage"
      />
      <van-button
        type="primary"
        size="small"
        round
        :disabled="!messageInput.trim() || sending"
        :loading="sending"
        @click="sendMessage"
      >
        发送
      </van-button>
    </div>
  </div>
</template>

<style scoped>
.chat-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: var(--im-bg-chat);
}

.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 12px 12px 0;
  -webkit-overflow-scrolling: touch;
}

.loading-indicator {
  display: flex;
  justify-content: center;
  padding: 16px;
}

.messages-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding-bottom: 12px;
}

.message-item {
  display: flex;
  align-items: flex-start;
  gap: 8px;
}

.message-item.message-self {
  flex-direction: row-reverse;
}

.message-avatar {
  flex-shrink: 0;
}

.message-content {
  max-width: 70%;
  display: flex;
  flex-direction: column;
}

.message-self .message-content {
  align-items: flex-end;
}

.message-sender {
  font-size: 12px;
  color: var(--im-text-secondary);
  margin-bottom: 4px;
  padding-left: 8px;
}

.message-bubble {
  padding: 10px 12px;
  border-radius: 12px;
  background-color: var(--im-bubble-other);
  word-break: break-word;
  font-size: 15px;
  line-height: 1.5;
  box-shadow: var(--im-shadow-sm);
}

.message-self .message-bubble {
  background-color: var(--im-bubble-self);
  border-top-right-radius: 4px;
}

.message-item:not(.message-self) .message-bubble {
  border-top-left-radius: 4px;
}

.message-time {
  font-size: 11px;
  color: var(--im-text-placeholder);
  margin-top: 4px;
  padding: 0 8px;
}

.file-message {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 180px;
}

.file-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.file-name {
  font-size: 14px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-size {
  font-size: 12px;
  color: var(--im-text-secondary);
}

.input-area {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  padding-bottom: calc(8px + var(--im-safe-area-bottom));
  background-color: var(--im-bg-white);
  border-top: 1px solid var(--im-border-light);
}

.input-action {
  color: var(--im-text-secondary);
  cursor: pointer;
  flex-shrink: 0;
}

.message-input {
  flex: 1;
  background-color: var(--im-bg-gray);
  border-radius: 20px;
  padding: 0 16px;
}

.message-input :deep(.van-field__body) {
  padding: 8px 0;
}
</style>
