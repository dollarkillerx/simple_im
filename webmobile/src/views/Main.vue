<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useChatStore } from '@/stores/chat'

const route = useRoute()
const router = useRouter()
const chatStore = useChatStore()

const active = ref(0)

// Sync active tab with route
watch(
  () => route.name,
  (name) => {
    if (name === 'Messages') active.value = 0
    else if (name === 'Contacts') active.value = 1
    else if (name === 'Profile') active.value = 2
  },
  { immediate: true }
)

function onChange(index: number) {
  const routes: string[] = ['/', '/contacts', '/profile']
  if (routes[index]) {
    router.push(routes[index])
  }
}
</script>

<template>
  <div class="main-page">
    <router-view />

    <van-tabbar
      v-model="active"
      :safe-area-inset-bottom="true"
      :border="false"
      active-color="var(--im-primary)"
      inactive-color="var(--im-text-muted)"
      @change="onChange"
    >
      <van-tabbar-item>
        <template #icon="props">
          <div class="tab-icon-wrap">
            <van-icon :name="props.active ? 'chat' : 'chat-o'" size="24" />
            <div v-if="chatStore.totalUnread > 0" class="tab-badge">
              {{ chatStore.totalUnread > 99 ? '99+' : chatStore.totalUnread }}
            </div>
          </div>
        </template>
        消息
      </van-tabbar-item>

      <van-tabbar-item>
        <template #icon="props">
          <div class="tab-icon-wrap">
            <van-icon :name="props.active ? 'friends' : 'friends-o'" size="24" />
            <div v-if="chatStore.pendingRequests.length > 0" class="tab-badge">
              {{ chatStore.pendingRequests.length }}
            </div>
          </div>
        </template>
        通讯录
      </van-tabbar-item>

      <van-tabbar-item>
        <template #icon="props">
          <van-icon :name="props.active ? 'contact' : 'contact-o'" size="24" />
        </template>
        我的
      </van-tabbar-item>
    </van-tabbar>
  </div>
</template>

<style scoped>
.main-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: var(--im-bg);
}

.main-page :deep(.van-tabbar) {
  height: 56px;
  box-shadow: 0 -2px 12px rgba(0, 0, 0, 0.05);
}

.main-page :deep(.van-tabbar-item__text) {
  font-size: 11px;
  margin-top: 2px;
}

.main-page :deep(.van-tabbar-item--active .van-tabbar-item__text) {
  font-weight: 600;
}

.tab-icon-wrap {
  position: relative;
  display: inline-block;
}

.tab-badge {
  position: absolute;
  top: -4px;
  right: -10px;
  min-width: 16px;
  height: 16px;
  padding: 0 4px;
  font-size: 10px;
  font-weight: 600;
  color: #fff;
  background: linear-gradient(135deg, #ef4444 0%, #f87171 100%);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
