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
    <div class="page-header">
      <div class="header-left" @click="goBack">
        <van-icon name="arrow-left" size="22" />
      </div>
      <h1 class="header-title">添加好友</h1>
      <div class="header-right"></div>
    </div>

    <div class="page-content">
      <div class="search-card">
        <div class="search-icon">
          <van-icon name="search" size="20" />
        </div>
        <input
          v-model="username"
          class="search-input"
          placeholder="输入用户名搜索"
          :disabled="loading"
          @keyup.enter="handleAdd"
        />
        <van-button
          type="primary"
          size="small"
          round
          :loading="loading"
          :disabled="!username.trim()"
          @click="handleAdd"
        >
          添加
        </van-button>
      </div>

      <div class="tips-card">
        <div class="tip-item">
          <div class="tip-icon">
            <van-icon name="user-o" size="20" />
          </div>
          <div class="tip-content">
            <h4>搜索用户名</h4>
            <p>输入对方的用户名来查找好友</p>
          </div>
        </div>
        <div class="tip-item">
          <div class="tip-icon">
            <van-icon name="passed" size="20" />
          </div>
          <div class="tip-content">
            <h4>等待确认</h4>
            <p>对方同意后即可成为好友</p>
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
  transition: background-color 0.2s;
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
  padding: 20px 16px;
}

.search-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: var(--im-bg-white);
  border-radius: var(--im-radius-lg);
  box-shadow: var(--im-shadow-sm);
}

.search-icon {
  color: var(--im-text-muted);
}

.search-input {
  flex: 1;
  border: none;
  background: transparent;
  font-size: 15px;
  color: var(--im-text-primary);
  outline: none;
}

.search-input::placeholder {
  color: var(--im-text-muted);
}

.search-input:disabled {
  opacity: 0.6;
}

.tips-card {
  margin-top: 24px;
  background: var(--im-bg-white);
  border-radius: var(--im-radius-lg);
  padding: 8px 0;
  box-shadow: var(--im-shadow-sm);
}

.tip-item {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 16px 20px;
}

.tip-item + .tip-item {
  border-top: 1px solid var(--im-border-light);
}

.tip-icon {
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(99, 102, 241, 0.1);
  border-radius: 12px;
  color: var(--im-primary);
}

.tip-content h4 {
  font-size: 15px;
  font-weight: 600;
  color: var(--im-text-primary);
  margin-bottom: 2px;
}

.tip-content p {
  font-size: 13px;
  color: var(--im-text-muted);
}
</style>
