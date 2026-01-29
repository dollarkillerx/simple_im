<script setup lang="ts">
import { onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { useChatStore } from '@/stores/chat'

const userStore = useUserStore()
const chatStore = useChatStore()

onMounted(async () => {
  if (userStore.token) {
    await userStore.fetchUserInfo()
    if (userStore.isLoggedIn) {
      userStore.initWebSocket()
      chatStore.initWSListener()
      chatStore.loadFriends()
      chatStore.loadGroups()
      chatStore.loadPendingRequests()
    }
  }
})
</script>

<template>
  <van-config-provider>
    <router-view v-slot="{ Component }">
      <transition name="slide-right" mode="out-in">
        <component :is="Component" />
      </transition>
    </router-view>
  </van-config-provider>
</template>

<style>
/* App-level styles are in styles/index.css */
</style>
