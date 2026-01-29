<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { showToast } from 'vant'
import { groupApi } from '@/api/rpc'
import { getDefaultAvatar, formatTime } from '@/utils/format'

const props = defineProps<{
  id: string
}>()

const router = useRouter()

const loading = ref(true)
const groupInfo = ref<{
  id: number
  name: string
  avatar: string
  owner_id: number
  owner_name: string
  created_at: string
  members: Array<{
    user_id: number
    username: string
    nickname: string
    avatar: string
    role: number
    joined_at: string
  }>
} | null>(null)

function getAvatar(item: { avatar?: string; name?: string; nickname?: string; username?: string }) {
  return item.avatar || getDefaultAvatar(item.name || item.nickname || item.username || 'G')
}

function getRoleName(role: number) {
  switch (role) {
    case 1:
      return '群主'
    case 2:
      return '管理员'
    default:
      return ''
  }
}

function getRoleColor(role: number) {
  switch (role) {
    case 1:
      return 'linear-gradient(135deg, #f59e0b 0%, #fbbf24 100%)'
    case 2:
      return 'linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%)'
    default:
      return ''
  }
}

async function loadGroupInfo() {
  loading.value = true
  try {
    groupInfo.value = await groupApi.getInfo(parseInt(props.id))
  } catch (e) {
    showToast((e as Error).message || '加载失败')
  } finally {
    loading.value = false
  }
}

function goBack() {
  router.back()
}

function goToChat() {
  router.push(`/chat/group/${props.id}`)
}

onMounted(() => {
  loadGroupInfo()
})
</script>

<template>
  <div class="page">
    <div class="page-header">
      <div class="header-left" @click="goBack">
        <van-icon name="arrow-left" size="22" />
      </div>
      <h1 class="header-title">群组详情</h1>
      <div class="header-right">
        <van-icon name="ellipsis" size="22" />
      </div>
    </div>

    <div class="page-content">
      <van-skeleton :loading="loading" :row="5" class="skeleton">
        <template v-if="groupInfo">
          <!-- Group Header Card -->
          <div class="group-card">
            <div class="group-bg"></div>
            <div class="group-content">
              <van-image
                :src="getAvatar({ ...groupInfo, name: groupInfo.name })"
                width="80"
                height="80"
                round
                fit="cover"
                class="group-avatar"
              />
              <h2 class="group-name">{{ groupInfo.name }}</h2>
              <div class="group-stats">
                <div class="stat-item">
                  <span class="stat-value">{{ groupInfo.members?.length || 0 }}</span>
                  <span class="stat-label">成员</span>
                </div>
                <div class="stat-divider"></div>
                <div class="stat-item">
                  <span class="stat-value">{{ formatTime(groupInfo.created_at) }}</span>
                  <span class="stat-label">创建于</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Members Section -->
          <div class="section">
            <div class="section-header">
              <span class="section-title">群成员</span>
              <span class="section-count">{{ groupInfo.members?.length || 0 }}</span>
            </div>

            <div class="member-grid">
              <div
                v-for="member in groupInfo.members"
                :key="member.user_id"
                class="member-card"
              >
                <div class="member-avatar-wrap">
                  <van-image
                    :src="getAvatar(member)"
                    width="48"
                    height="48"
                    round
                    fit="cover"
                  />
                  <div
                    v-if="member.role <= 2"
                    class="role-badge"
                    :style="{ background: getRoleColor(member.role) }"
                  >
                    <van-icon :name="member.role === 1 ? 'star' : 'shield-o'" size="10" />
                  </div>
                </div>
                <span class="member-name">{{ member.nickname || member.username }}</span>
                <span v-if="member.role <= 2" class="member-role">{{ getRoleName(member.role) }}</span>
              </div>
            </div>
          </div>

          <!-- Action Button -->
          <div class="action-section">
            <van-button
              type="primary"
              block
              round
              size="large"
              @click="goToChat"
            >
              <van-icon name="chat-o" style="margin-right: 8px" />
              发送消息
            </van-button>
          </div>
        </template>
      </van-skeleton>
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
  z-index: 10;
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

.header-left:active,
.header-right:active {
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

.skeleton {
  padding: 20px;
}

.group-card {
  position: relative;
  background: var(--im-bg-white);
  border-radius: var(--im-radius-xl);
  overflow: hidden;
  box-shadow: var(--im-shadow-md);
  margin-bottom: 24px;
}

.group-bg {
  height: 80px;
  background: var(--im-primary-gradient);
}

.group-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 0 20px 24px;
  margin-top: -40px;
}

.group-avatar {
  border: 4px solid var(--im-bg-white);
  box-shadow: var(--im-shadow-md);
}

.group-name {
  font-size: 22px;
  font-weight: 700;
  color: var(--im-text-primary);
  margin-top: 12px;
}

.group-stats {
  display: flex;
  align-items: center;
  gap: 24px;
  margin-top: 16px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.stat-value {
  font-size: 16px;
  font-weight: 600;
  color: var(--im-text-primary);
}

.stat-label {
  font-size: 12px;
  color: var(--im-text-muted);
  margin-top: 2px;
}

.stat-divider {
  width: 1px;
  height: 24px;
  background: var(--im-border);
}

.section {
  margin-bottom: 24px;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
  padding: 0 4px;
}

.section-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--im-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.section-count {
  font-size: 13px;
  color: var(--im-primary);
  font-weight: 600;
}

.member-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(80px, 1fr));
  gap: 12px;
}

.member-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  padding: 12px 8px;
  background: var(--im-bg-white);
  border-radius: var(--im-radius-lg);
  box-shadow: var(--im-shadow-xs);
}

.member-avatar-wrap {
  position: relative;
}

.role-badge {
  position: absolute;
  bottom: -2px;
  right: -2px;
  width: 18px;
  height: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  border: 2px solid var(--im-bg-white);
  color: #fff;
}

.member-name {
  font-size: 12px;
  color: var(--im-text-primary);
  text-align: center;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 100%;
}

.member-role {
  font-size: 10px;
  color: var(--im-primary);
  font-weight: 500;
}

.action-section {
  padding: 8px 0;
}

.action-section :deep(.van-button) {
  height: 52px;
  font-size: 16px;
  font-weight: 600;
}
</style>
