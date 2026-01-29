import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/Login.vue'),
      meta: { guest: true },
    },
    {
      path: '/',
      name: 'Main',
      component: () => import('@/views/Main.vue'),
      meta: { auth: true },
      children: [
        {
          path: '',
          name: 'Messages',
          component: () => import('@/views/Messages.vue'),
        },
        {
          path: 'contacts',
          name: 'Contacts',
          component: () => import('@/views/Contacts.vue'),
        },
        {
          path: 'profile',
          name: 'Profile',
          component: () => import('@/views/Profile.vue'),
        },
      ],
    },
    {
      path: '/chat/:type/:id',
      name: 'Chat',
      component: () => import('@/views/Chat.vue'),
      meta: { auth: true },
      props: true,
    },
    {
      path: '/add-friend',
      name: 'AddFriend',
      component: () => import('@/views/AddFriend.vue'),
      meta: { auth: true },
    },
    {
      path: '/friend-requests',
      name: 'FriendRequests',
      component: () => import('@/views/FriendRequests.vue'),
      meta: { auth: true },
    },
    {
      path: '/create-group',
      name: 'CreateGroup',
      component: () => import('@/views/CreateGroup.vue'),
      meta: { auth: true },
    },
    {
      path: '/group/:id',
      name: 'GroupInfo',
      component: () => import('@/views/GroupInfo.vue'),
      meta: { auth: true },
      props: true,
    },
  ],
})

router.beforeEach(async (to, _from, next) => {
  const userStore = useUserStore()

  // Try to restore user info if we have a token
  if (userStore.token && !userStore.userInfo) {
    await userStore.fetchUserInfo()
  }

  if (to.meta.auth && !userStore.isLoggedIn) {
    next('/login')
  } else if (to.meta.guest && userStore.isLoggedIn) {
    next('/')
  } else {
    next()
  }
})

export default router
