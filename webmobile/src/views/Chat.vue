<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useRouter } from 'vue-router'
import { showToast } from 'vant'
import { useChatStore } from '@/stores/chat'
import { useUserStore } from '@/stores/user'
import { MessageType } from '@/types'
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
            <!-- Avatar (other) -->
            <van-image
              v-if="!isSelf(msg.sender_id)"
              :src="getAvatar(msg.Sender?.nickname || msg.Sender?.username || 'U', msg.Sender?.avatar)"
              width="36"
              height="36"
              round
              fit="cover"
              class="msg-avatar"
            />

            <div class="msg-main">
              <!-- Sender Name (group only) -->
              <div v-if="!isSelf(msg.sender_id) && type === 'group'" class="msg-sender">
                {{ msg.Sender?.nickname || msg.Sender?.username }}
              </div>

              <!-- Bubble -->
              <div class="msg-bubble" :class="{ 'bubble-self': isSelf(msg.sender_id) }">
                <!-- Text -->
                <template v-if="msg.msg_type === MessageType.Text">
                  <span class="msg-text">{{ msg.content }}</span>
                </template>

                <!-- Image -->
                <template v-else-if="msg.msg_type === MessageType.Image">
                  <van-image :src="msg.file_url" width="180" fit="cover" radius="8" />
                </template>

                <!-- File -->
                <template v-else-if="msg.msg_type === MessageType.File">
                  <div class="file-card">
                    <van-icon name="description" size="28" />
                    <div class="file-info">
                      <span class="file-name">{{ msg.file_name }}</span>
                      <span class="file-size">{{ formatFileSize(msg.file_size) }}</span>
                    </div>
                  </div>
                </template>
              </div>
            </div>

            <!-- Avatar (self) -->
            <van-image
              v-if="isSelf(msg.sender_id)"
              :src="getAvatar(userStore.userInfo?.nickname || 'U', userStore.userInfo?.avatar)"
              width="36"
              height="36"
              round
              fit="cover"
              class="msg-avatar"
            />
          </div>
        </template>
      </div>
    </div>

    <!-- Input Area -->
    <div class="input-bar">
      <van-field
        v-model="messageInput"
        type="textarea"
        :autosize="{ minHeight: 36, maxHeight: 100 }"
        placeholder="输入消息..."
        :border="false"
        @keydown.enter.exact.prevent="sendMessage"
      />
      <van-button
        type="primary"
        size="small"
        round
        :disabled="!messageInput.trim()"
        :loading="sending"
        class="send-btn"
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
  background: var(--im-bg);
}

.chat-page :deep(.van-nav-bar) {
  background: var(--im-bg-white);
  box-shadow: var(--im-shadow-xs);
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
  gap: 8px;
}

/* Time Divider */
.time-divider {
  text-align: center;
  padding: 8px 0;
}

.time-divider::before {
  content: attr(data-time);
}

.time-divider {
  font-size: 12px;
  color: var(--im-text-muted);
  background: var(--im-bg);
  padding: 4px 12px;
  border-radius: 12px;
  display: inline-block;
  margin: 8px auto;
  align-self: center;
}

/* Message Row */
.message-row {
  display: flex;
  align-items: flex-start;
  gap: 8px;
}

.message-row.is-self {
  flex-direction: row-reverse;
}

.msg-avatar {
  flex-shrink: 0;
}

.msg-main {
  max-width: 70%;
  display: flex;
  flex-direction: column;
}

.is-self .msg-main {
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
  border-radius: 16px;
  border-top-left-radius: 4px;
  box-shadow: var(--im-shadow-xs);
  word-break: break-word;
}

.msg-bubble.bubble-self {
  background: var(--im-primary);
  color: #fff;
  border-top-left-radius: 16px;
  border-top-right-radius: 4px;
}

.msg-text {
  font-size: 15px;
  line-height: 1.5;
  white-space: pre-wrap;
}

/* File Card */
.file-card {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 180px;
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
  display: flex;
  align-items: flex-end;
  gap: 8px;
  padding: 8px 12px;
  padding-bottom: calc(8px + var(--im-safe-area-bottom));
  background: var(--im-bg-white);
  border-top: 1px solid var(--im-border-light);
}

.input-bar :deep(.van-field) {
  flex: 1;
  background: var(--im-bg);
  border-radius: 18px;
  padding: 6px 14px;
}

.input-bar :deep(.van-field__control) {
  font-size: 15px;
  line-height: 1.4;
}

.send-btn {
  flex-shrink: 0;
  height: 36px;
  padding: 0 16px;
}
</style>
