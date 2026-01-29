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
    <div class="page-header">
      <div class="header-left" @click="goBack">
        <van-icon name="arrow-left" size="22" />
      </div>
      <h1 class="header-title">好友请求</h1>
      <div class="header-right"></div>
    </div>

    <div class="page-content">
      <div v-if="chatStore.pendingRequests.length === 0" class="empty-state">
        <div class="empty-icon">
          <van-icon name="friends-o" size="48" />
        </div>
        <p class="empty-text">暂无好友请求</p>
      </div>

      <div v-else class="request-list">
        <div
          v-for="request in chatStore.pendingRequests"
          :key="request.id"
          class="request-card"
        >
          <van-image
            :src="getAvatar(request)"
            width="56"
            height="56"
            round
            fit="cover"
            class="request-avatar"
          />

          <div class="request-info">
            <div class="request-name">{{ request.nickname || request.username }}</div>
            <div class="request-time">{{ formatTime(request.created_at) }}</div>
          </div>

          <div class="request-actions">
            <button
              class="action-btn reject"
              :disabled="processing.has(request.id)"
              @click="handleRequest(request.id, false)"
            >
              <van-icon name="cross" size="16" />
            </button>
            <button
              class="action-btn accept"
              :disabled="processing.has(request.id)"
              @click="handleRequest(request.id, true)"
            >
              <van-icon name="success" size="16" />
            </button>
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

.page-header {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  background-color: var(--im-bg-white);
  box-shadow: var(--im-shadow-xs);
}

.header-left,
.header-right {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  border-radius: 50%;
}

.header-left:active {
  background-color: var(--im-bg-card);
}

.header-title {
  flex: 1;
  text-align: center;
  font-size: 17px;
  font-weight: 600;
  color: var(--im-text-primary);
}

.page-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 60px 24px;
}

.empty-icon {
  width: 80px;
  height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(99, 102, 241, 0.1);
  border-radius: 50%;
  color: var(--im-primary);
  margin-bottom: 16px;
}

.empty-text {
  font-size: 16px;
  color: var(--im-text-muted);
}

.request-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.request-card {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 16px;
  background: var(--im-bg-white);
  border-radius: var(--im-radius-lg);
  box-shadow: var(--im-shadow-sm);
}

.request-avatar {
  flex-shrink: 0;
  box-shadow: var(--im-shadow-xs);
}

.request-info {
  flex: 1;
  min-width: 0;
}

.request-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--im-text-primary);
  margin-bottom: 4px;
}

.request-time {
  font-size: 13px;
  color: var(--im-text-muted);
}

.request-actions {
  display: flex;
  gap: 10px;
}

.action-btn {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
}

.action-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.action-btn.reject {
  background: rgba(239, 68, 68, 0.1);
  color: var(--im-danger);
}

.action-btn.reject:active:not(:disabled) {
  background: rgba(239, 68, 68, 0.2);
}

.action-btn.accept {
  background: var(--im-primary-gradient);
  color: #fff;
  box-shadow: var(--im-shadow-primary);
}

.action-btn.accept:active:not(:disabled) {
  transform: scale(0.95);
}
</style>
