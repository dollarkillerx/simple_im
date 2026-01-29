<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { showToast, showLoadingToast, closeToast, showSuccessToast } from 'vant'
import { useChatStore } from '@/stores/chat'
import { groupApi } from '@/api/rpc'
import { getDefaultAvatar } from '@/utils/format'

const router = useRouter()
const chatStore = useChatStore()

const groupName = ref('')
const selectedMembers = ref<number[]>([])
const loading = ref(false)

const canCreate = computed(() => {
  return groupName.value.trim().length > 0
})

function getAvatar(item: { avatar: string; nickname?: string; username?: string }) {
  return item.avatar || getDefaultAvatar(item.nickname || item.username || 'U')
}

function toggleMember(userId: number) {
  const index = selectedMembers.value.indexOf(userId)
  if (index === -1) {
    selectedMembers.value.push(userId)
  } else {
    selectedMembers.value.splice(index, 1)
  }
}

function isMemberSelected(userId: number) {
  return selectedMembers.value.includes(userId)
}

async function handleCreate() {
  if (!groupName.value.trim()) {
    showToast('请输入群组名称')
    return
  }

  loading.value = true
  showLoadingToast({ message: '创建中...', forbidClick: true })

  try {
    await groupApi.create(groupName.value.trim(), selectedMembers.value)
    closeToast()
    showSuccessToast('群组创建成功')
    await chatStore.loadGroups()
    setTimeout(() => {
      router.back()
    }, 1000)
  } catch (e) {
    closeToast()
    showToast((e as Error).message || '创建失败')
  } finally {
    loading.value = false
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
      <h1 class="header-title">创建群组</h1>
      <van-button
        type="primary"
        size="small"
        round
        :disabled="!canCreate"
        :loading="loading"
        @click="handleCreate"
      >
        创建
      </van-button>
    </div>

    <div class="page-content">
      <div class="form-section">
        <div class="form-label">群组名称</div>
        <div class="input-card">
          <van-icon name="cluster-o" size="20" color="var(--im-primary)" />
          <input
            v-model="groupName"
            class="form-input"
            placeholder="给群组起个名字"
            maxlength="50"
          />
          <span class="char-count">{{ groupName.length }}/50</span>
        </div>
      </div>

      <div class="form-section">
        <div class="form-label">
          选择成员
          <span v-if="selectedMembers.length > 0" class="selected-count">
            已选 {{ selectedMembers.length }} 人
          </span>
        </div>

        <div v-if="chatStore.friends.length > 0" class="member-list">
          <div
            v-for="friend in chatStore.friends"
            :key="friend.id"
            class="member-item"
            :class="{ selected: isMemberSelected(friend.user_id) }"
            @click="toggleMember(friend.user_id)"
          >
            <van-image
              :src="getAvatar(friend)"
              width="44"
              height="44"
              round
              fit="cover"
              class="member-avatar"
            />
            <span class="member-name">{{ friend.nickname || friend.username }}</span>
            <div class="check-icon">
              <van-icon v-if="isMemberSelected(friend.user_id)" name="success" size="16" />
            </div>
          </div>
        </div>

        <div v-else class="empty-state">
          <van-icon name="friends-o" size="32" color="var(--im-text-muted)" />
          <p>暂无好友可添加</p>
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
  gap: 12px;
  padding: 12px 16px;
  background-color: var(--im-bg-white);
  box-shadow: var(--im-shadow-xs);
}

.header-left {
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
  font-size: 17px;
  font-weight: 600;
  color: var(--im-text-primary);
}

.page-content {
  flex: 1;
  overflow-y: auto;
  padding: 20px 16px;
}

.form-section {
  margin-bottom: 24px;
}

.form-label {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 13px;
  font-weight: 600;
  color: var(--im-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 12px;
  padding-left: 4px;
}

.selected-count {
  color: var(--im-primary);
  text-transform: none;
}

.input-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  background: var(--im-bg-white);
  border-radius: var(--im-radius-lg);
  box-shadow: var(--im-shadow-sm);
}

.form-input {
  flex: 1;
  border: none;
  background: transparent;
  font-size: 15px;
  color: var(--im-text-primary);
  outline: none;
}

.form-input::placeholder {
  color: var(--im-text-muted);
}

.char-count {
  font-size: 12px;
  color: var(--im-text-muted);
}

.member-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
  gap: 12px;
}

.member-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 16px 8px;
  background: var(--im-bg-white);
  border-radius: var(--im-radius-lg);
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
}

.member-item:active {
  transform: scale(0.98);
}

.member-item.selected {
  background: rgba(99, 102, 241, 0.1);
  box-shadow: inset 0 0 0 2px var(--im-primary);
}

.member-avatar {
  box-shadow: var(--im-shadow-xs);
}

.member-name {
  font-size: 13px;
  color: var(--im-text-primary);
  text-align: center;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 100%;
}

.check-icon {
  position: absolute;
  top: 8px;
  right: 8px;
  width: 22px;
  height: 22px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--im-primary-gradient);
  border-radius: 50%;
  color: #fff;
  opacity: 0;
  transform: scale(0.8);
  transition: all 0.2s;
}

.member-item.selected .check-icon {
  opacity: 1;
  transform: scale(1);
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 40px;
  color: var(--im-text-muted);
  font-size: 14px;
}
</style>
