<script setup lang="ts">
import { useRouter } from 'vue-router'
import { showConfirmDialog } from 'vant'
import { useUserStore } from '@/stores/user'
import { getDefaultAvatar } from '@/utils/format'

const router = useRouter()
const userStore = useUserStore()

function getAvatar() {
  return userStore.userInfo?.avatar || getDefaultAvatar(userStore.userInfo?.nickname || 'U')
}

async function handleLogout() {
  try {
    await showConfirmDialog({
      title: '退出登录',
      message: '确定要退出登录吗？',
      confirmButtonColor: 'var(--im-danger)',
    })
    userStore.logout()
    router.replace('/login')
  } catch {
    // User cancelled
  }
}
</script>

<template>
  <div class="page">
    <van-nav-bar title="我的" fixed placeholder />

    <div class="page-content">
      <div class="profile-header">
        <div class="profile-avatar">
          <van-image
            :src="getAvatar()"
            width="72"
            height="72"
            round
            fit="cover"
          />
        </div>
        <div class="profile-info">
          <h2 class="profile-name">{{ userStore.userInfo?.nickname }}</h2>
          <p class="profile-username">@{{ userStore.userInfo?.username }}</p>
        </div>
      </div>

      <van-cell-group inset class="settings-group">
        <van-cell title="账号设置" icon="setting-o" is-link />
        <van-cell title="消息通知" icon="bell" is-link />
        <van-cell title="隐私设置" icon="shield-o" is-link />
      </van-cell-group>

      <van-cell-group inset class="settings-group">
        <van-cell title="关于我们" icon="info-o" is-link />
        <van-cell title="帮助与反馈" icon="question-o" is-link />
      </van-cell-group>

      <div class="logout-section">
        <van-button
          type="danger"
          plain
          block
          round
          @click="handleLogout"
        >
          退出登录
        </van-button>
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
  padding-bottom: calc(var(--im-tabbar-height) + var(--im-safe-area-bottom));
}

.profile-header {
  display: flex;
  align-items: center;
  padding: 24px 20px;
  margin: 12px;
  background-color: var(--im-bg-white);
  border-radius: var(--im-radius-lg);
}

.profile-avatar {
  flex-shrink: 0;
  margin-right: 16px;
}

.profile-info {
  flex: 1;
  min-width: 0;
}

.profile-name {
  font-size: 20px;
  font-weight: 600;
  color: var(--im-text-primary);
  margin-bottom: 4px;
}

.profile-username {
  font-size: 14px;
  color: var(--im-text-secondary);
}

.settings-group {
  margin: 12px;
}

.logout-section {
  padding: 24px 24px;
}
</style>
