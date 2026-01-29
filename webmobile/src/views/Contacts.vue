<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useChatStore } from '@/stores/chat'
import { getDefaultAvatar } from '@/utils/format'

const router = useRouter()
const chatStore = useChatStore()

const activeTab = ref(0)
const refreshing = ref(false)

const friends = computed(() => chatStore.friends)
const groups = computed(() => chatStore.groups)

async function onRefresh() {
  await Promise.all([
    chatStore.loadFriends(),
    chatStore.loadGroups(),
    chatStore.loadPendingRequests(),
  ])
  refreshing.value = false
}

function goToChat(type: 'private' | 'group', id: number) {
  router.push(`/chat/${type}/${id}`)
}

function goToAddFriend() {
  router.push('/add-friend')
}

function goToFriendRequests() {
  router.push('/friend-requests')
}

function goToCreateGroup() {
  router.push('/create-group')
}

function goToGroupInfo(id: number) {
  router.push(`/group/${id}`)
}

function getAvatar(item: { avatar: string; name?: string; nickname?: string; username?: string }) {
  return item.avatar || getDefaultAvatar(item.name || item.nickname || item.username || 'U')
}
</script>

<template>
  <div class="page">
    <van-nav-bar title="通讯录" :border="false">
      <template #right>
        <van-icon
          :name="activeTab === 0 ? 'add-o' : 'add-o'"
          size="22"
          @click="activeTab === 0 ? goToAddFriend() : goToCreateGroup()"
        />
      </template>
    </van-nav-bar>

    <div class="page-content">
      <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
        <van-tabs v-model:active="activeTab" animated swipeable>
          <van-tab title="好友">
            <div class="tab-content">
              <!-- Friend Requests Entry -->
              <van-cell
                v-if="chatStore.pendingRequests.length > 0"
                is-link
                center
                class="request-cell"
                @click="goToFriendRequests"
              >
                <template #icon>
                  <div class="request-icon">
                    <van-icon name="friends-o" size="20" color="#fff" />
                  </div>
                </template>
                <template #title>
                  <span class="request-title">新的好友请求</span>
                </template>
                <template #value>
                  <van-badge :content="chatStore.pendingRequests.length" />
                </template>
              </van-cell>

              <van-cell-group v-if="friends.length > 0" :border="false">
                <van-cell
                  v-for="friend in friends"
                  :key="friend.id"
                  :title="friend.nickname || friend.username"
                  is-link
                  center
                  @click="goToChat('private', friend.user_id)"
                >
                  <template #icon>
                    <van-image
                      :src="getAvatar(friend)"
                      width="44"
                      height="44"
                      round
                      fit="cover"
                      class="cell-avatar"
                    />
                  </template>
                </van-cell>
              </van-cell-group>

              <van-empty
                v-if="friends.length === 0 && chatStore.pendingRequests.length === 0"
                image="friends"
                description="暂无好友"
              >
                <van-button type="primary" round size="small" @click="goToAddFriend">
                  添加好友
                </van-button>
              </van-empty>
            </div>
          </van-tab>

          <van-tab title="群组">
            <div class="tab-content">
              <van-cell-group v-if="groups.length > 0" :border="false">
                <van-cell
                  v-for="group in groups"
                  :key="group.id"
                  :title="group.name"
                  :label="`${group.owner_name} 创建`"
                  is-link
                  center
                  @click="goToChat('group', group.id)"
                >
                  <template #icon>
                    <van-image
                      :src="getAvatar({ ...group, name: group.name })"
                      width="44"
                      height="44"
                      round
                      fit="cover"
                      class="cell-avatar"
                    />
                  </template>
                  <template #right-icon>
                    <van-icon
                      name="info-o"
                      size="20"
                      class="info-icon"
                      @click.stop="goToGroupInfo(group.id)"
                    />
                  </template>
                </van-cell>
              </van-cell-group>

              <van-empty v-else image="search" description="暂无群组">
                <van-button type="primary" round size="small" @click="goToCreateGroup">
                  创建群组
                </van-button>
              </van-empty>
            </div>
          </van-tab>
        </van-tabs>
      </van-pull-refresh>
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

.page :deep(.van-nav-bar) {
  background: var(--im-bg-white);
}

.page-content {
  flex: 1;
  overflow-y: auto;
  padding-bottom: calc(56px + var(--im-safe-area-bottom));
}

.page :deep(.van-tabs__nav) {
  background: var(--im-bg-white);
}

.page :deep(.van-tab) {
  font-size: 15px;
}

.page :deep(.van-tab--active) {
  font-weight: 600;
}

.page :deep(.van-tabs__line) {
  background: var(--im-primary);
  width: 24px !important;
  height: 3px;
  border-radius: 2px;
}

.tab-content {
  min-height: 300px;
  padding-top: 8px;
}

.request-cell {
  margin: 8px 12px;
  border-radius: 12px;
  background: linear-gradient(135deg, rgba(249, 115, 22, 0.1) 0%, rgba(251, 146, 60, 0.05) 100%);
}

.request-cell :deep(.van-cell__left-icon) {
  margin-right: 12px;
}

.request-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #f97316 0%, #fb923c 100%);
  border-radius: 50%;
}

.request-title {
  font-weight: 500;
}

.cell-avatar {
  margin-right: 12px;
}

.info-icon {
  color: var(--im-text-muted);
  padding: 8px;
  margin-right: -8px;
}

.info-icon:active {
  color: var(--im-primary);
}

.page :deep(.van-empty) {
  padding: 40px 0;
}

.page :deep(.van-empty__bottom) {
  margin-top: 16px;
}
</style>
