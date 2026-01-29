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
    <!-- Decorative blobs -->
    <div class="blob blob-1"></div>
    <div class="blob blob-2"></div>
    <div class="blob blob-3"></div>

    <div class="login-container">
      <div class="login-header">
        <div class="login-logo">
          <svg viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
            <defs>
              <linearGradient id="logoGradient" x1="0%" y1="0%" x2="100%" y2="100%">
                <stop offset="0%" style="stop-color:#6366f1"/>
                <stop offset="100%" style="stop-color:#8b5cf6"/>
              </linearGradient>
            </defs>
            <path d="M44 24C44 12.954 35.046 4 24 4S4 12.954 4 24c0 4.87 1.74 9.334 4.632 12.81L6 44l8.19-2.632A19.904 19.904 0 0024 44c11.046 0 20-8.954 20-20z" fill="url(#logoGradient)"/>
            <path d="M16 20h16M16 28h10" stroke="#fff" stroke-width="3" stroke-linecap="round"/>
          </svg>
        </div>
        <h1 class="login-title">Simple IM</h1>
        <p class="login-subtitle">{{ isLogin ? '欢迎回来' : '创建您的账户' }}</p>
      </div>

      <div class="login-card">
        <van-form @submit="handleSubmit">
          <div class="form-group">
            <label class="form-label">用户名</label>
            <van-field
              v-model="username"
              name="username"
              placeholder="请输入用户名"
              :rules="[{ required: true, message: '请输入用户名' }]"
              clearable
              class="custom-field"
            />
          </div>

          <div class="form-group">
            <label class="form-label">密码</label>
            <van-field
              v-model="password"
              type="password"
              name="password"
              placeholder="请输入密码"
              :rules="[{ required: true, message: '请输入密码' }]"
              clearable
              class="custom-field"
            />
          </div>

          <div v-if="!isLogin" class="form-group">
            <label class="form-label">昵称 <span class="optional">(可选)</span></label>
            <van-field
              v-model="nickname"
              name="nickname"
              placeholder="请输入昵称"
              clearable
              class="custom-field"
            />
          </div>

          <van-button
            round
            block
            type="primary"
            native-type="submit"
            class="submit-btn"
          >
            {{ isLogin ? '登 录' : '注 册' }}
          </van-button>

          <div class="login-footer">
            <span class="footer-text">{{ isLogin ? '还没有账户？' : '已有账户？' }}</span>
            <a class="footer-link" @click="toggleMode">
              {{ isLogin ? '立即注册' : '立即登录' }}
            </a>
          </div>
        </van-form>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 50%, #6366f1 100%);
  padding: 20px;
  position: relative;
  overflow: hidden;
}

/* Decorative blobs */
.blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(60px);
  opacity: 0.6;
  animation: float 8s ease-in-out infinite;
}

.blob-1 {
  width: 300px;
  height: 300px;
  background: rgba(139, 92, 246, 0.5);
  top: -100px;
  left: -100px;
}

.blob-2 {
  width: 250px;
  height: 250px;
  background: rgba(99, 102, 241, 0.5);
  bottom: -80px;
  right: -80px;
  animation-delay: -4s;
}

.blob-3 {
  width: 200px;
  height: 200px;
  background: rgba(168, 85, 247, 0.4);
  top: 50%;
  left: 50%;
  animation-delay: -2s;
}

@keyframes float {
  0%, 100% { transform: translate(0, 0) scale(1); }
  50% { transform: translate(30px, -30px) scale(1.1); }
}

.login-container {
  width: 100%;
  max-width: 380px;
  position: relative;
  z-index: 1;
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-logo {
  width: 72px;
  height: 72px;
  margin: 0 auto 16px;
  filter: drop-shadow(0 8px 16px rgba(99, 102, 241, 0.4));
}

.login-logo svg {
  width: 100%;
  height: 100%;
}

.login-title {
  font-size: 32px;
  font-weight: 700;
  color: #fff;
  margin-bottom: 8px;
  letter-spacing: -0.5px;
}

.login-subtitle {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.8);
}

.login-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-radius: 24px;
  padding: 32px 24px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
}

.form-group {
  margin-bottom: 20px;
}

.form-label {
  display: block;
  font-size: 14px;
  font-weight: 600;
  color: var(--im-text-primary);
  margin-bottom: 8px;
}

.optional {
  font-weight: 400;
  color: var(--im-text-muted);
}

.custom-field {
  background: var(--im-bg-input);
  border-radius: 12px;
  padding: 4px 0;
}

.custom-field :deep(.van-field__body) {
  padding: 12px 16px;
}

.custom-field :deep(.van-field__control) {
  font-size: 15px;
}

.submit-btn {
  height: 52px;
  font-size: 16px;
  font-weight: 600;
  margin-top: 8px;
}

.login-footer {
  text-align: center;
  margin-top: 24px;
  font-size: 14px;
}

.footer-text {
  color: var(--im-text-secondary);
}

.footer-link {
  color: var(--im-primary);
  font-weight: 600;
  margin-left: 4px;
  cursor: pointer;
  transition: opacity 0.2s;
}

.footer-link:hover {
  opacity: 0.8;
}

@media (min-width: 768px) {
  .login-page {
    background: var(--im-bg);
  }

  .login-card {
    background: var(--im-bg-white);
  }

  .login-title {
    background: var(--im-primary-gradient);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }

  .login-subtitle {
    color: var(--im-text-secondary);
  }

  .blob {
    display: none;
  }
}
</style>
