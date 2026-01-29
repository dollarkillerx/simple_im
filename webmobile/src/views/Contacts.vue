<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useChatStore } from '@/stores/chat'
import { getDefaultAvatar } from '@/utils/format'

const router = useRouter()
const chatStore = useChatStore()

const activeTab = ref('friends')
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
    <van-nav-bar title="通讯录" fixed placeholder>
      <template #right>
        <van-icon
          name="add-o"
          size="22"
          @click="activeTab === 'friends' ? goToAddFriend() : goToCreateGroup()"
        />
      </template>
    </van-nav-bar>

    <div class="page-content">
      <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
        <van-tabs v-model:active="activeTab" sticky offset-top="46">
          <van-tab title="好友" name="friends">
            <div class="contact-section">
              <!-- Friend Requests Entry -->
              <div
                v-if="chatStore.pendingRequests.length > 0"
                class="request-entry"
                @click="goToFriendRequests"
              >
                <div class="request-icon">
                  <van-icon name="friends-o" size="24" />
                </div>
                <div class="request-content">
                  <span class="request-title">新的好友</span>
                  <van-badge :content="chatStore.pendingRequests.length" />
                </div>
                <van-icon name="arrow" color="var(--im-text-placeholder)" />
              </div>

              <van-cell-group v-if="friends.length > 0" inset>
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
                      width="40"
                      height="40"
                      round
                      fit="cover"
                      class="friend-avatar"
                    />
                  </template>
                </van-cell>
              </van-cell-group>

              <van-empty
                v-else-if="chatStore.pendingRequests.length === 0"
                description="暂无好友"
              />
            </div>
          </van-tab>

          <van-tab title="群组" name="groups">
            <div class="contact-section">
              <van-cell-group v-if="groups.length > 0" inset>
                <van-cell
                  v-for="group in groups"
                  :key="group.id"
                  :title="group.name"
                  :label="`${group.owner_name} 创建`"
                  is-link
                  center
                >
                  <template #icon>
                    <van-image
                      :src="getAvatar({ ...group, name: group.name })"
                      width="40"
                      height="40"
                      round
                      fit="cover"
                      class="friend-avatar"
                    />
                  </template>
                  <template #right-icon>
                    <div class="group-actions">
                      <van-icon
                        name="info-o"
                        size="20"
                        @click.stop="goToGroupInfo(group.id)"
                      />
                      <van-icon
                        name="chat-o"
                        size="20"
                        @click.stop="goToChat('group', group.id)"
                      />
                    </div>
                  </template>
                </van-cell>
              </van-cell-group>

              <van-empty v-else description="暂无群组" />
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

.page-content {
  flex: 1;
  overflow-y: auto;
  padding-bottom: calc(var(--im-tabbar-height) + var(--im-safe-area-bottom));
}

.contact-section {
  padding: 12px 0;
}

.request-entry {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  margin: 0 12px 12px;
  background-color: var(--im-bg-white);
  border-radius: var(--im-radius-lg);
  cursor: pointer;
}

.request-entry:active {
  background-color: var(--im-bg-gray);
}

.request-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #ff976a 0%, #ff6034 100%);
  border-radius: 50%;
  color: var(--im-text-white);
  margin-right: 12px;
}

.request-content {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 8px;
}

.request-title {
  font-size: 15px;
  font-weight: 500;
}

.friend-avatar {
  margin-right: 12px;
}

.group-actions {
  display: flex;
  align-items: center;
  gap: 16px;
  color: var(--im-text-secondary);
}

.group-actions .van-icon {
  cursor: pointer;
}
</style>
