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
    <van-nav-bar
      title="创建群组"
      left-arrow
      fixed
      placeholder
      @click-left="goBack"
    >
      <template #right>
        <van-button
          type="primary"
          size="small"
          :disabled="!canCreate"
          :loading="loading"
          @click="handleCreate"
        >
          创建
        </van-button>
      </template>
    </van-nav-bar>

    <div class="page-content">
      <van-cell-group inset class="group-info-section">
        <van-field
          v-model="groupName"
          label="群名称"
          placeholder="请输入群组名称"
          :maxlength="50"
          clearable
        />
      </van-cell-group>

      <div class="section-title">
        选择群成员（可选）
        <span v-if="selectedMembers.length > 0" class="selected-count">
          已选 {{ selectedMembers.length }} 人
        </span>
      </div>

      <van-cell-group inset>
        <van-cell
          v-for="friend in chatStore.friends"
          :key="friend.id"
          :title="friend.nickname || friend.username"
          clickable
          center
          @click="toggleMember(friend.user_id)"
        >
          <template #icon>
            <van-image
              :src="getAvatar(friend)"
              width="40"
              height="40"
              round
              fit="cover"
              class="friend-avatar"
            />
          </template>
          <template #right-icon>
            <van-icon
              :name="isMemberSelected(friend.user_id) ? 'success' : 'circle'"
              :color="isMemberSelected(friend.user_id) ? 'var(--im-primary)' : 'var(--im-text-placeholder)'"
              size="20"
            />
          </template>
        </van-cell>
      </van-cell-group>

      <van-empty
        v-if="chatStore.friends.length === 0"
        description="暂无好友可添加"
      />
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
  padding: 12px 0;
}

.group-info-section {
  margin-bottom: 12px;
}

.section-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 24px;
  font-size: 14px;
  color: var(--im-text-secondary);
}

.selected-count {
  color: var(--im-primary);
  font-weight: 500;
}

.friend-avatar {
  margin-right: 12px;
}
</style>
