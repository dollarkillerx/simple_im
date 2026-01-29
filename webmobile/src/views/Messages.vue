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
    <van-nav-bar title="消息" :border="false">
      <template #right>
        <van-icon name="add-o" size="22" />
      </template>
    </van-nav-bar>

    <div class="page-content">
      <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
        <!-- Search Bar -->
        <van-search
          v-model="searchValue"
          placeholder="搜索"
          shape="round"
          :clearable="false"
        />

        <!-- Empty State -->
        <van-empty
          v-if="chatStore.sortedConversations.length === 0"
          image="search"
          description="暂无消息"
        />

        <!-- Conversation List -->
        <van-cell-group v-else :border="false" class="conv-list">
          <van-swipe-cell
            v-for="conv in chatStore.sortedConversations"
            :key="conv.id"
          >
            <van-cell
              :title="conv.name"
              :label="conv.lastMessage || '暂无消息'"
              center
              clickable
              @click="goToChat(conv)"
            >
              <template #icon>
                <div class="conv-avatar-wrap">
                  <van-image
                    :src="getAvatar(conv)"
                    width="48"
                    height="48"
                    round
                    fit="cover"
                  />
                  <div v-if="conv.type === 'group'" class="group-badge">
                    <van-icon name="friends-o" size="10" />
                  </div>
                </div>
              </template>
              <template #value>
                <div class="conv-right">
                  <span class="conv-time">{{ formatTime(conv.lastTime) }}</span>
                  <van-badge
                    v-if="conv.unread > 0"
                    :content="conv.unread > 99 ? '99+' : conv.unread"
                  />
                </div>
              </template>
            </van-cell>

            <template #right>
              <van-button square type="warning" text="置顶" />
              <van-button square type="danger" text="删除" />
            </template>
          </van-swipe-cell>
        </van-cell-group>
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

.page :deep(.van-nav-bar) {
  background: var(--im-bg-white);
}

.page :deep(.van-nav-bar__title) {
  font-size: 18px;
  font-weight: 600;
}

.page-content {
  flex: 1;
  overflow-y: auto;
  padding-bottom: calc(56px + var(--im-safe-area-bottom));
}

/* Search */
.page :deep(.van-search) {
  padding: 8px 12px;
  background: var(--im-bg-white);
}

.page :deep(.van-search__content) {
  background: var(--im-bg);
}

/* Conversation List */
.conv-list {
  margin-top: 8px;
}

.conv-list :deep(.van-cell) {
  padding: 12px 16px;
}

.conv-list :deep(.van-cell__title) {
  font-weight: 500;
}

.conv-list :deep(.van-cell__label) {
  margin-top: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.conv-avatar-wrap {
  position: relative;
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
  background: var(--im-primary);
  border-radius: 50%;
  border: 2px solid var(--im-bg-white);
  color: #fff;
}

.conv-right {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 6px;
}

.conv-time {
  font-size: 12px;
  color: var(--im-text-muted);
}

/* Swipe buttons */
.conv-list :deep(.van-swipe-cell__right) {
  display: flex;
}

.conv-list :deep(.van-swipe-cell__right .van-button) {
  height: 100%;
  border-radius: 0;
}
</style>
