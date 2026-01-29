<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useChatStore } from '@/stores/chat'
import { formatTime, getDefaultAvatar } from '@/utils/format'

const router = useRouter()
const chatStore = useChatStore()

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
  <div class="messages-page">
    <!-- Header -->
    <van-nav-bar title="消息" :border="false">
      <template #right>
        <van-icon name="add-o" size="22" />
      </template>
    </van-nav-bar>

    <!-- Content -->
    <div class="messages-content">
      <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
        <!-- Empty State -->
        <van-empty
          v-if="chatStore.sortedConversations.length === 0"
          image="search"
          description="暂无消息"
        />

        <!-- Conversation List -->
        <div v-else class="conv-list">
          <van-swipe-cell
            v-for="conv in chatStore.sortedConversations"
            :key="conv.id"
          >
            <div class="conv-item" @click="goToChat(conv)">
              <!-- Avatar -->
              <div class="conv-avatar">
                <van-image
                  :src="getAvatar(conv)"
                  width="50"
                  height="50"
                  round
                  fit="cover"
                />
                <div v-if="conv.type === 'group'" class="group-badge">
                  <van-icon name="friends-o" size="10" />
                </div>
              </div>

              <!-- Info -->
              <div class="conv-info">
                <div class="conv-top">
                  <span class="conv-name">{{ conv.name }}</span>
                  <span class="conv-time">{{ formatTime(conv.lastTime) }}</span>
                </div>
                <div class="conv-bottom">
                  <span class="conv-msg">{{ conv.lastMessage || '暂无消息' }}</span>
                  <van-badge
                    v-if="conv.unread > 0"
                    :content="conv.unread > 99 ? '99+' : conv.unread"
                  />
                </div>
              </div>
            </div>

            <template #right>
              <van-button square type="warning" text="置顶" class="swipe-btn" />
              <van-button square type="danger" text="删除" class="swipe-btn" />
            </template>
          </van-swipe-cell>
        </div>
      </van-pull-refresh>
    </div>
  </div>
</template>

<style scoped>
.messages-page {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--im-bg);
}

.messages-page :deep(.van-nav-bar) {
  flex-shrink: 0;
  background: var(--im-bg-white);
}

.messages-content {
  flex: 1;
  overflow-y: auto;
  padding-bottom: calc(60px + env(safe-area-inset-bottom, 0px));
}

/* Conversation List */
.conv-list {
  background: var(--im-bg-white);
  margin-top: 10px;
}

.conv-item {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  background: var(--im-bg-white);
  cursor: pointer;
}

.conv-item:active {
  background: var(--im-bg);
}

.conv-avatar {
  position: relative;
  flex-shrink: 0;
  margin-right: 12px;
}

.group-badge {
  position: absolute;
  bottom: 0;
  right: 0;
  width: 18px;
  height: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--im-primary);
  border-radius: 50%;
  border: 2px solid var(--im-bg-white);
  color: #fff;
}

.conv-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.conv-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
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
  color: var(--im-text-muted);
  flex-shrink: 0;
  margin-left: 8px;
}

.conv-bottom {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.conv-msg {
  font-size: 14px;
  color: var(--im-text-muted);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
}

/* Swipe buttons */
.swipe-btn {
  height: 100%;
}
</style>
