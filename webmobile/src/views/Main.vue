<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useChatStore } from '@/stores/chat'

const route = useRoute()
const chatStore = useChatStore()

const activeTab = computed(() => {
  if (route.name === 'Messages') return 'messages'
  if (route.name === 'Contacts') return 'contacts'
  if (route.name === 'Profile') return 'profile'
  return 'messages'
})
</script>

<template>
  <div class="main-page">
    <router-view />

    <van-tabbar v-model="activeTab" route>
      <van-tabbar-item to="/" name="messages" icon="chat-o">
        <template #icon="{ active }">
          <van-icon :name="active ? 'chat' : 'chat-o'" />
          <van-badge
            v-if="chatStore.totalUnread > 0"
            :content="chatStore.totalUnread > 99 ? '99+' : chatStore.totalUnread"
            class="tab-badge"
          />
        </template>
        消息
      </van-tabbar-item>
      <van-tabbar-item to="/contacts" name="contacts" icon="friends-o">
        <template #icon="{ active }">
          <van-icon :name="active ? 'friends' : 'friends-o'" />
          <van-badge
            v-if="chatStore.pendingRequests.length > 0"
            :content="chatStore.pendingRequests.length"
            class="tab-badge"
          />
        </template>
        通讯录
      </van-tabbar-item>
      <van-tabbar-item to="/profile" name="profile" icon="contact-o">
        <template #icon="{ active }">
          <van-icon :name="active ? 'contact' : 'contact-o'" />
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
  padding-bottom: var(--im-safe-area-bottom);
}

.tab-badge {
  position: absolute;
  top: -4px;
  right: -8px;
}

.tab-badge :deep(.van-badge) {
  transform: scale(0.9);
}
</style>
