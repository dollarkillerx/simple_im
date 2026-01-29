<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useChatStore } from '@/stores/chat'
import { formatTime, getDefaultAvatar } from '@/utils/format'

const router = useRouter()
const chatStore = useChatStore()

const searchValue = ref('')
const refreshing = ref(false)

async function onRefresh() {
  await Promise.all([
    chatStore.loadFriends(),
    chatStore.loadGroups(),
    chatStore.loadPendingRequests(),
  ])
  refreshing.value = false
}

function goToChat(conv: { type: string; targetId: number }) {
  router.push(`/chat/${conv.type}/${conv.targetId}`)
}

function getAvatar(conv: { avatar: string; name: string }) {
  return conv.avatar || getDefaultAvatar(conv.name)
}
</script>

<template>
  <div class="page">
    <van-nav-bar title="消息" fixed placeholder />

    <div class="page-content">
      <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
        <van-search
          v-model="searchValue"
          placeholder="搜索"
          shape="round"
          background="transparent"
        />

        <div v-if="chatStore.sortedConversations.length === 0" class="empty-state">
          <van-empty description="暂无消息" />
        </div>

        <div v-else class="conversation-list">
          <van-swipe-cell
            v-for="conv in chatStore.sortedConversations"
            :key="conv.id"
          >
            <div class="conversation-item" @click="goToChat(conv)">
              <div class="conv-avatar">
                <van-image
                  :src="getAvatar(conv)"
                  :alt="conv.name"
                  width="48"
                  height="48"
                  round
                  fit="cover"
                />
                <span v-if="conv.type === 'group'" class="group-badge">
                  <van-icon name="friends-o" size="12" />
                </span>
              </div>
              <div class="conv-content">
                <div class="conv-header">
                  <span class="conv-name">{{ conv.name }}</span>
                  <span class="conv-time">{{ formatTime(conv.lastTime) }}</span>
                </div>
                <div class="conv-message">
                  <span class="conv-text">{{ conv.lastMessage || '暂无消息' }}</span>
                  <van-badge
                    v-if="conv.unread > 0"
                    :content="conv.unread > 99 ? '99+' : conv.unread"
                  />
                </div>
              </div>
            </div>
            <template #right>
              <van-button square type="danger" text="删除" class="swipe-btn" />
            </template>
          </van-swipe-cell>
        </div>
      </van-pull-refresh>
    </div>
  </div>
</template>

<style scoped>
.page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: var(--im-bg);
}

.page-content {
  flex: 1;
  overflow-y: auto;
  padding-bottom: calc(var(--im-tabbar-height) + var(--im-safe-area-bottom));
}

.conversation-list {
  background-color: var(--im-bg-white);
}

.conversation-item {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  cursor: pointer;
  transition: background-color var(--im-transition-fast);
}

.conversation-item:active {
  background-color: var(--im-bg-gray);
}

.conv-avatar {
  position: relative;
  flex-shrink: 0;
  margin-right: 12px;
}

.group-badge {
  position: absolute;
  bottom: -2px;
  right: -2px;
  width: 18px;
  height: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--im-primary);
  border-radius: 50%;
  border: 2px solid var(--im-bg-white);
  color: var(--im-text-white);
}

.conv-content {
  flex: 1;
  min-width: 0;
  padding: 2px 0;
  border-bottom: 1px solid var(--im-border-light);
}

.conversation-item:last-child .conv-content {
  border-bottom: none;
}

.conv-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 6px;
}

.conv-name {
  font-size: 16px;
  font-weight: 500;
  color: var(--im-text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.conv-time {
  font-size: 12px;
  color: var(--im-text-placeholder);
  flex-shrink: 0;
  margin-left: 8px;
}

.conv-message {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.conv-text {
  flex: 1;
  font-size: 13px;
  color: var(--im-text-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.swipe-btn {
  height: 100%;
}
</style>
