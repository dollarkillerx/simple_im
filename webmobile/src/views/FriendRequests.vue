<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { showToast, showSuccessToast } from 'vant'
import { useChatStore } from '@/stores/chat'
import { friendApi } from '@/api/rpc'
import { formatTime, getDefaultAvatar } from '@/utils/format'

const router = useRouter()
const chatStore = useChatStore()

const processing = ref<Set<number>>(new Set())

function getAvatar(item: { avatar: string; nickname?: string; username?: string }) {
  return item.avatar || getDefaultAvatar(item.nickname || item.username || 'U')
}

async function handleRequest(requestId: number, accept: boolean) {
  if (processing.value.has(requestId)) return

  processing.value.add(requestId)

  try {
    await friendApi.accept(requestId, accept)
    showSuccessToast(accept ? '已同意' : '已拒绝')
    await chatStore.loadPendingRequests()
    if (accept) {
      await chatStore.loadFriends()
    }
  } catch (e) {
    showToast((e as Error).message || '操作失败')
  } finally {
    processing.value.delete(requestId)
  }
}

function goBack() {
  router.back()
}
</script>

<template>
  <div class="page">
    <van-nav-bar
      title="好友请求"
      left-arrow
      fixed
      placeholder
      @click-left="goBack"
    />

    <div class="page-content">
      <van-empty
        v-if="chatStore.pendingRequests.length === 0"
        description="暂无好友请求"
      />

      <div v-else class="request-list">
        <div
          v-for="request in chatStore.pendingRequests"
          :key="request.id"
          class="request-item"
        >
          <van-image
            :src="getAvatar(request)"
            width="48"
            height="48"
            round
            fit="cover"
            class="request-avatar"
          />

          <div class="request-info">
            <div class="request-name">{{ request.nickname || request.username }}</div>
            <div class="request-time">{{ formatTime(request.created_at) }}</div>
          </div>

          <div class="request-actions">
            <van-button
              size="small"
              plain
              :loading="processing.has(request.id)"
              @click="handleRequest(request.id, false)"
            >
              拒绝
            </van-button>
            <van-button
              type="primary"
              size="small"
              :loading="processing.has(request.id)"
              @click="handleRequest(request.id, true)"
            >
              同意
            </van-button>
          </div>
        </div>
      </div>
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
}

.request-list {
  background-color: var(--im-bg-white);
}

.request-item {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid var(--im-border-light);
}

.request-item:last-child {
  border-bottom: none;
}

.request-avatar {
  flex-shrink: 0;
  margin-right: 12px;
}

.request-info {
  flex: 1;
  min-width: 0;
}

.request-name {
  font-size: 16px;
  font-weight: 500;
  color: var(--im-text-primary);
  margin-bottom: 4px;
}

.request-time {
  font-size: 12px;
  color: var(--im-text-placeholder);
}

.request-actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}
</style>
