# Simple IM - WebMobile

基于 Vue 3 + Vant 4 的即时通讯 Web 应用，支持 Mobile 和 PC 端自适应布局。

## 技术栈

- **Vue 3** - 渐进式 JavaScript 框架
- **Vant 4** - 移动端 UI 组件库
- **Pinia** - Vue 状态管理
- **Vue Router** - 路由管理
- **TypeScript** - 类型安全
- **Vite** - 构建工具

## 功能特性

- 用户登录/注册
- 好友列表与管理
- 好友请求发送与处理
- 群组创建与管理
- 实时消息收发 (WebSocket)
- 消息历史记录
- 响应式布局 (Mobile & PC 自适应)

## 开发

```bash
# 安装依赖
pnpm install

# 启动开发服务器
pnpm dev

# 构建生产版本
pnpm build

# 预览生产版本
pnpm preview
```

## 项目结构

```
src/
├── api/          # API 接口封装
├── assets/       # 静态资源
├── components/   # 公共组件
├── router/       # 路由配置
├── stores/       # Pinia 状态管理
├── styles/       # 全局样式
├── types/        # TypeScript 类型定义
├── utils/        # 工具函数
├── views/        # 页面组件
├── App.vue       # 根组件
└── main.ts       # 入口文件
```

## 页面说明

| 页面 | 路径 | 说明 |
|------|------|------|
| 登录 | `/login` | 用户登录/注册 |
| 消息 | `/` | 会话列表 |
| 通讯录 | `/contacts` | 好友和群组列表 |
| 我的 | `/profile` | 个人信息 |
| 聊天 | `/chat/:type/:id` | 聊天页面 |
| 添加好友 | `/add-friend` | 搜索添加好友 |
| 好友请求 | `/friend-requests` | 待处理的好友请求 |
| 创建群组 | `/create-group` | 创建新群组 |
| 群组详情 | `/group/:id` | 群组成员列表 |

## 后端接口

应用通过 JSON-RPC 2.0 协议与后端通信：

- `user.login` - 用户登录
- `user.register` - 用户注册
- `user.info` - 获取用户信息
- `friend.list` - 获取好友列表
- `friend.add` - 添加好友
- `friend.accept` - 处理好友请求
- `friend.pending` - 获取待处理请求
- `group.list` - 获取群组列表
- `group.create` - 创建群组
- `group.info` - 获取群组详情
- `message.send` - 发送消息
- `message.history` - 获取历史消息

## 响应式设计

- **Mobile (< 768px)**: 全屏布局
- **PC (>= 768px)**: 居中卡片式布局，最大宽度 480px
