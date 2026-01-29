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
    <div class="page-header">
      <h1 class="page-title">我的</h1>
    </div>

    <div class="page-content">
      <!-- Profile Card -->
      <div class="profile-card">
        <div class="profile-bg"></div>
        <div class="profile-content">
          <div class="profile-avatar">
            <van-image
              :src="getAvatar()"
              width="80"
              height="80"
              round
              fit="cover"
            />
            <div class="avatar-edit">
              <van-icon name="photograph" size="14" />
            </div>
          </div>
          <h2 class="profile-name">{{ userStore.userInfo?.nickname }}</h2>
          <p class="profile-id">ID: {{ userStore.userInfo?.id }}</p>
          <p class="profile-username">@{{ userStore.userInfo?.username }}</p>
        </div>
      </div>

      <!-- Settings Groups -->
      <div class="settings-section">
        <div class="section-title">账户设置</div>
        <div class="settings-card">
          <div class="settings-item">
            <div class="item-icon" style="background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%)">
              <van-icon name="manager-o" size="18" />
            </div>
            <span class="item-text">个人资料</span>
            <van-icon name="arrow" size="16" color="var(--im-text-muted)" />
          </div>
          <div class="settings-item">
            <div class="item-icon" style="background: linear-gradient(135deg, #10b981 0%, #34d399 100%)">
              <van-icon name="shield-o" size="18" />
            </div>
            <span class="item-text">账号安全</span>
            <van-icon name="arrow" size="16" color="var(--im-text-muted)" />
          </div>
          <div class="settings-item">
            <div class="item-icon" style="background: linear-gradient(135deg, #f59e0b 0%, #fbbf24 100%)">
              <van-icon name="bell" size="18" />
            </div>
            <span class="item-text">消息通知</span>
            <van-icon name="arrow" size="16" color="var(--im-text-muted)" />
          </div>
        </div>
      </div>

      <div class="settings-section">
        <div class="section-title">其他</div>
        <div class="settings-card">
          <div class="settings-item">
            <div class="item-icon" style="background: linear-gradient(135deg, #06b6d4 0%, #22d3ee 100%)">
              <van-icon name="question-o" size="18" />
            </div>
            <span class="item-text">帮助与反馈</span>
            <van-icon name="arrow" size="16" color="var(--im-text-muted)" />
          </div>
          <div class="settings-item">
            <div class="item-icon" style="background: linear-gradient(135deg, #8b5cf6 0%, #a78bfa 100%)">
              <van-icon name="info-o" size="18" />
            </div>
            <span class="item-text">关于我们</span>
            <span class="item-value">v1.0.0</span>
            <van-icon name="arrow" size="16" color="var(--im-text-muted)" />
          </div>
        </div>
      </div>

      <div class="logout-section">
        <van-button
          type="danger"
          plain
          block
          round
          size="large"
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

.page-header {
  padding: 16px 20px 8px;
  background-color: var(--im-bg-white);
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: var(--im-text-primary);
}

.page-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  padding-bottom: calc(var(--im-tabbar-height) + var(--im-safe-area-bottom) + 16px);
}

.profile-card {
  position: relative;
  background: var(--im-bg-white);
  border-radius: var(--im-radius-xl);
  overflow: hidden;
  box-shadow: var(--im-shadow-md);
  margin-bottom: 24px;
}

.profile-bg {
  height: 80px;
  background: var(--im-primary-gradient);
}

.profile-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 0 20px 24px;
  margin-top: -40px;
}

.profile-avatar {
  position: relative;
}

.profile-avatar :deep(.van-image) {
  border: 4px solid var(--im-bg-white);
  box-shadow: var(--im-shadow-md);
}

.avatar-edit {
  position: absolute;
  bottom: 0;
  right: 0;
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--im-primary-gradient);
  border-radius: 50%;
  border: 2px solid var(--im-bg-white);
  color: #fff;
  cursor: pointer;
}

.profile-name {
  font-size: 22px;
  font-weight: 700;
  color: var(--im-text-primary);
  margin-top: 12px;
}

.profile-id {
  font-size: 12px;
  color: var(--im-text-muted);
  margin-top: 4px;
}

.profile-username {
  font-size: 14px;
  color: var(--im-text-secondary);
  margin-top: 2px;
}

.settings-section {
  margin-bottom: 24px;
}

.section-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--im-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 12px;
  padding-left: 4px;
}

.settings-card {
  background: var(--im-bg-white);
  border-radius: var(--im-radius-lg);
  overflow: hidden;
  box-shadow: var(--im-shadow-sm);
}

.settings-item {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 14px 16px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.settings-item:active {
  background-color: var(--im-bg-card);
}

.settings-item + .settings-item {
  border-top: 1px solid var(--im-border-light);
}

.item-icon {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 10px;
  color: #fff;
}

.item-text {
  flex: 1;
  font-size: 15px;
  font-weight: 500;
  color: var(--im-text-primary);
}

.item-value {
  font-size: 14px;
  color: var(--im-text-muted);
  margin-right: 4px;
}

.logout-section {
  padding: 8px 0;
}

.logout-section :deep(.van-button) {
  height: 48px;
  font-size: 16px;
  font-weight: 600;
}
</style>
