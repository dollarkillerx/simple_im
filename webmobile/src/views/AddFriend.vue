<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { showToast, showLoadingToast, closeToast, showSuccessToast } from 'vant'
import { friendApi } from '@/api/rpc'

const router = useRouter()

const username = ref('')
const loading = ref(false)

async function handleAdd() {
  if (!username.value.trim()) {
    showToast('请输入用户名')
    return
  }

  loading.value = true
  showLoadingToast({ message: '请稍候...', forbidClick: true })

  try {
    await friendApi.add({ username: username.value.trim() })
    closeToast()
    showSuccessToast('好友请求已发送')
    username.value = ''
    setTimeout(() => {
      router.back()
    }, 1000)
  } catch (e) {
    closeToast()
    showToast((e as Error).message || '添加失败')
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
      title="添加好友"
      left-arrow
      fixed
      placeholder
      @click-left="goBack"
    />

    <div class="page-content">
      <div class="search-section">
        <van-search
          v-model="username"
          placeholder="输入用户名搜索"
          show-action
          :disabled="loading"
          @search="handleAdd"
        >
          <template #action>
            <van-button
              type="primary"
              size="small"
              :loading="loading"
              :disabled="!username.trim()"
              @click="handleAdd"
            >
              添加
            </van-button>
          </template>
        </van-search>
      </div>

      <div class="tips">
        <p>输入对方的用户名来添加好友</p>
        <p>对方同意后即可成为好友</p>
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

.search-section {
  padding: 16px 0;
  background-color: var(--im-bg-white);
}

.search-section :deep(.van-search) {
  padding-right: 12px;
}

.search-section :deep(.van-search__content) {
  background-color: var(--im-bg-gray);
}

.search-section :deep(.van-search__action) {
  padding-left: 12px;
}

.tips {
  padding: 24px;
  text-align: center;
  color: var(--im-text-secondary);
  font-size: 14px;
  line-height: 1.8;
}
</style>
