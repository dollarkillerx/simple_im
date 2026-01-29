<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useRouter } from 'vue-router'
import { showToast, showLoadingToast, closeToast } from 'vant'
import { useChatStore } from '@/stores/chat'
import { useUserStore } from '@/stores/user'
import { fileApi } from '@/api/rpc'
import { MessageType, type MessageTypeValue } from '@/types'
import { formatChatTime, getDefaultAvatar, formatFileSize } from '@/utils/format'

const props = defineProps<{
  type: 'private' | 'group'
  id: string
}>()

const router = useRouter()
const chatStore = useChatStore()
const userStore = useUserStore()

const messageInput = ref('')
const messagesContainer = ref<HTMLElement | null>(null)
const imageInput = ref<HTMLInputElement | null>(null)
const fileInput = ref<HTMLInputElement | null>(null)
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

function shouldShowTime(index: number) {
  if (index === 0) return true
  const curr = messages.value[index]
  const prev = messages.value[index - 1]
  if (!prev || !curr) return true
  return new Date(curr.created_at).getTime() - new Date(prev.created_at).getTime() > 300000
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

function selectImage() {
  imageInput.value?.click()
}

function selectFile() {
  fileInput.value?.click()
}

async function handleImageSelect(e: Event) {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return
  target.value = '' // Reset input

  await uploadAndSend(file, MessageType.Image)
}

async function handleFileSelect(e: Event) {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return
  target.value = '' // Reset input

  await uploadAndSend(file, MessageType.File)
}

async function uploadAndSend(file: File, msgType: MessageTypeValue) {
  showLoadingToast({
    message: '上传中...',
    forbidClick: true,
    duration: 0,
  })

  try {
    // Upload file
    const result = await fileApi.upload(file)

    // Send message with file info
    const params: {
      receiverId?: number
      groupId?: number
      msgType: MessageTypeValue
      fileUrl: string
      fileName: string
      fileSize: number
    } = {
      msgType,
      fileUrl: result.url,
      fileName: result.filename,
      fileSize: result.size,
    }

    if (props.type === 'group') {
      params.groupId = parseInt(props.id)
    } else {
      params.receiverId = parseInt(props.id)
    }

    await chatStore.sendMessage(params)
    scrollToBottom()
    closeToast()
  } catch (e) {
    closeToast()
    showToast((e as Error).message || '发送失败')
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

function openImage(url: string) {
  window.open(url)
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
    <!-- Header -->
    <van-nav-bar
      :title="chatTitle"
      left-arrow
      :border="false"
      @click-left="goBack"
    >
      <template #right>
        <van-icon
          v-if="type === 'group'"
          name="ellipsis"
          size="20"
          @click="goToInfo"
        />
      </template>
    </van-nav-bar>

    <!-- Messages -->
    <div ref="messagesContainer" class="messages-container">
      <div v-if="loading" class="loading-wrap">
        <van-loading size="24" />
      </div>

      <div v-else-if="messages.length === 0" class="empty-wrap">
        <van-empty image="search" description="暂无消息" />
      </div>

      <div v-else class="messages-list">
        <template v-for="(msg, index) in messages" :key="msg.id">
          <!-- Time Divider -->
          <div v-if="shouldShowTime(index)" class="time-divider">
            {{ formatChatTime(msg.created_at) }}
          </div>

          <!-- Message Item -->
          <div class="message-row" :class="{ 'is-self': isSelf(msg.sender_id) }">
            <!-- Avatar -->
            <img
              class="msg-avatar"
              :src="isSelf(msg.sender_id)
                ? getAvatar(userStore.userInfo?.nickname || 'U', userStore.userInfo?.avatar)
                : getAvatar(msg.Sender?.nickname || msg.Sender?.username || 'U', msg.Sender?.avatar)"
            />

            <div class="msg-content">
              <!-- Sender Name (group only, not self) -->
              <div v-if="!isSelf(msg.sender_id) && type === 'group'" class="msg-sender">
                {{ msg.Sender?.nickname || msg.Sender?.username }}
              </div>

              <!-- Bubble -->
              <div
                class="msg-bubble"
                :class="{
                  'bubble-self': isSelf(msg.sender_id),
                  'bubble-image': Number(msg.msg_type) === 2
                }"
              >
                <!-- Text -->
                <template v-if="Number(msg.msg_type) === 1">
                  {{ msg.content }}
                </template>

                <!-- Image -->
                <template v-else-if="Number(msg.msg_type) === 2">
                  <img
                    v-if="msg.file_url"
                    :src="msg.file_url"
                    class="msg-image"
                    @click="openImage(msg.file_url)"
                  />
                  <span v-else>{{ msg.file_name || '[图片]' }}</span>
                </template>

                <!-- File -->
                <template v-else-if="Number(msg.msg_type) === 3">
                  <a :href="msg.file_url" target="_blank" class="file-card">
                    <van-icon name="description" size="28" />
                    <div class="file-info">
                      <span class="file-name">{{ msg.file_name }}</span>
                      <span class="file-size">{{ formatFileSize(msg.file_size) }}</span>
                    </div>
                  </a>
                </template>

                <!-- Fallback -->
                <template v-else>
                  {{ msg.content || msg.file_name || '[消息]' }}
                </template>
              </div>
            </div>
          </div>
        </template>
      </div>
    </div>

    <!-- Input Area -->
    <div class="input-bar">
      <!-- Hidden file inputs -->
      <input
        ref="imageInput"
        type="file"
        accept="image/*"
        style="display: none"
        @change="handleImageSelect"
      />
      <input
        ref="fileInput"
        type="file"
        style="display: none"
        @change="handleFileSelect"
      />

      <!-- Action buttons -->
      <button class="action-btn" @click="selectImage" title="发送图片">
        <van-icon name="photo-o" size="22" />
      </button>
      <button class="action-btn" @click="selectFile" title="发送文件">
        <van-icon name="description" size="22" />
      </button>

      <!-- Text input -->
      <input
        v-model="messageInput"
        type="text"
        class="msg-input"
        placeholder="输入消息..."
        @keydown.enter="sendMessage"
      />

      <!-- Send button -->
      <button
        class="send-btn"
        :disabled="!messageInput.trim() || sending"
        @click="sendMessage"
      >
        发送
      </button>
    </div>
  </div>
</template>

<style scoped>
.chat-page {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  flex-direction: column;
  background: var(--im-bg);
}

@media (min-width: 768px) {
  .chat-page {
    left: 50%;
    transform: translateX(-50%);
    max-width: 420px;
    box-shadow: 0 0 20px rgba(0, 0, 0, 0.1);
  }
}

.chat-page :deep(.van-nav-bar) {
  flex-shrink: 0;
  background: var(--im-bg-white);
}

/* Messages Container */
.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 12px 16px;
  -webkit-overflow-scrolling: touch;
}

.loading-wrap,
.empty-wrap {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 200px;
}

.messages-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

/* Time Divider */
.time-divider {
  font-size: 12px;
  color: var(--im-text-muted);
  text-align: center;
  padding: 8px 0;
}

/* Message Row */
.message-row {
  display: flex;
  align-items: flex-start;
  gap: 10px;
}

.message-row.is-self {
  flex-direction: row-reverse;
}

.msg-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
  flex-shrink: 0;
}

.msg-content {
  max-width: 70%;
  display: flex;
  flex-direction: column;
}

.is-self .msg-content {
  align-items: flex-end;
}

.msg-sender {
  font-size: 12px;
  color: var(--im-text-muted);
  margin-bottom: 4px;
  margin-left: 4px;
}

/* Bubble */
.msg-bubble {
  padding: 10px 14px;
  background: var(--im-bg-white);
  border-radius: 18px;
  border-top-left-radius: 4px;
  box-shadow: var(--im-shadow-xs);
  word-break: break-word;
  font-size: 15px;
  line-height: 1.5;
}

.msg-bubble.bubble-self {
  background: var(--im-primary);
  color: #fff;
  border-top-left-radius: 18px;
  border-top-right-radius: 4px;
}

.msg-bubble.bubble-image {
  padding: 4px;
  background: transparent;
  box-shadow: none;
}

.msg-bubble.bubble-image.bubble-self {
  background: transparent;
}

.msg-image {
  max-width: 200px;
  max-height: 300px;
  border-radius: 8px;
  cursor: pointer;
}

/* File Card */
.file-card {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 160px;
  text-decoration: none;
  color: inherit;
}

.file-info {
  flex: 1;
  overflow: hidden;
}

.file-name {
  display: block;
  font-size: 14px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-size {
  font-size: 12px;
  opacity: 0.7;
}

/* Input Bar */
.input-bar {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  padding-bottom: calc(10px + env(safe-area-inset-bottom, 0px));
  background: var(--im-bg-white);
  border-top: 1px solid var(--im-border-light);
}

.action-btn {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  border-radius: 50%;
  background: transparent;
  color: var(--im-text-secondary);
  cursor: pointer;
}

.action-btn:active {
  background: var(--im-bg);
}

.msg-input {
  flex: 1;
  height: 40px;
  padding: 0 16px;
  border: none;
  border-radius: 20px;
  background: var(--im-bg);
  font-size: 15px;
  outline: none;
}

.msg-input::placeholder {
  color: var(--im-text-muted);
}

.send-btn {
  flex-shrink: 0;
  height: 40px;
  padding: 0 20px;
  border: none;
  border-radius: 20px;
  background: var(--im-primary);
  color: #fff;
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
}

.send-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.send-btn:not(:disabled):active {
  opacity: 0.8;
}
</style>
