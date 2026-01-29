<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { showToast, showLoadingToast, closeToast } from 'vant'
import { useUserStore } from '@/stores/user'
import { useChatStore } from '@/stores/chat'

const router = useRouter()
const userStore = useUserStore()
const chatStore = useChatStore()

const isLogin = ref(true)
const username = ref('')
const password = ref('')
const nickname = ref('')

async function handleSubmit() {
  if (!username.value || !password.value) {
    showToast('请填写用户名和密码')
    return
  }

  if (username.value.length < 3) {
    showToast('用户名至少3个字符')
    return
  }

  if (password.value.length < 6) {
    showToast('密码至少6个字符')
    return
  }

  showLoadingToast({ message: '请稍候...', forbidClick: true })

  try {
    if (isLogin.value) {
      await userStore.login(username.value, password.value)
    } else {
      await userStore.register(username.value, password.value, nickname.value || undefined)
    }

    chatStore.initWSListener()
    chatStore.loadFriends()
    chatStore.loadGroups()
    chatStore.loadPendingRequests()

    closeToast()
    router.replace('/')
  } catch (e) {
    closeToast()
    showToast((e as Error).message || '操作失败')
  }
}

function toggleMode() {
  isLogin.value = !isLogin.value
}
</script>

<template>
  <div class="login-page">
    <div class="login-header">
      <div class="login-logo">
        <svg viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M44 24C44 12.954 35.046 4 24 4S4 12.954 4 24c0 4.87 1.74 9.334 4.632 12.81L6 44l8.19-2.632A19.904 19.904 0 0024 44c11.046 0 20-8.954 20-20z" fill="#1989fa"/>
          <path d="M16 20h16M16 28h10" stroke="#fff" stroke-width="3" stroke-linecap="round"/>
        </svg>
      </div>
      <h1 class="login-title">Simple IM</h1>
      <p class="login-subtitle">{{ isLogin ? '登录您的账户' : '创建新账户' }}</p>
    </div>

    <div class="login-form">
      <van-form @submit="handleSubmit">
        <van-cell-group inset>
          <van-field
            v-model="username"
            name="username"
            label="用户名"
            placeholder="请输入用户名"
            :rules="[{ required: true, message: '请输入用户名' }]"
            clearable
          />
          <van-field
            v-model="password"
            type="password"
            name="password"
            label="密码"
            placeholder="请输入密码"
            :rules="[{ required: true, message: '请输入密码' }]"
            clearable
          />
          <van-field
            v-if="!isLogin"
            v-model="nickname"
            name="nickname"
            label="昵称"
            placeholder="请输入昵称（可选）"
            clearable
          />
        </van-cell-group>

        <div class="login-actions">
          <van-button round block type="primary" native-type="submit">
            {{ isLogin ? '登 录' : '注 册' }}
          </van-button>
          <div class="login-switch">
            <span>{{ isLogin ? '还没有账户？' : '已有账户？' }}</span>
            <a @click="toggleMode">{{ isLogin ? '立即注册' : '立即登录' }}</a>
          </div>
        </div>
      </van-form>
    </div>
  </div>
</template>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: linear-gradient(180deg, var(--im-primary) 0%, var(--im-primary-light) 100%);
  padding: var(--im-safe-area-top) 0 var(--im-safe-area-bottom);
}

.login-header {
  flex: 0 0 auto;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 60px 24px 40px;
  color: var(--im-text-white);
}

.login-logo {
  width: 80px;
  height: 80px;
  margin-bottom: 16px;
}

.login-logo svg {
  width: 100%;
  height: 100%;
}

.login-title {
  font-size: 28px;
  font-weight: 600;
  margin-bottom: 8px;
}

.login-subtitle {
  font-size: 14px;
  opacity: 0.9;
}

.login-form {
  flex: 1;
  background-color: var(--im-bg-white);
  border-radius: 24px 24px 0 0;
  padding: 32px 16px;
}

.login-actions {
  padding: 24px 16px;
}

.login-actions .van-button {
  margin-bottom: 16px;
}

.login-switch {
  text-align: center;
  font-size: 14px;
  color: var(--im-text-secondary);
}

.login-switch a {
  color: var(--im-primary);
  margin-left: 4px;
  cursor: pointer;
}

@media (min-width: 768px) {
  .login-page {
    justify-content: center;
  }

  .login-header {
    padding-top: 40px;
  }

  .login-form {
    border-radius: 24px;
    max-width: 400px;
    margin: 0 auto;
    width: 100%;
  }
}
</style>
