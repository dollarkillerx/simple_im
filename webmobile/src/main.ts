import { createApp } from 'vue'
import { createPinia } from 'pinia'
import {
  Button,
  Cell,
  CellGroup,
  Icon,
  Image as VanImage,
  List,
  NavBar,
  Popup,
  Field,
  Form,
  Toast,
  Dialog,
  ActionSheet,
  Tabbar,
  TabbarItem,
  Badge,
  Loading,
  Empty,
  Search,
  PullRefresh,
  Uploader,
  Tab,
  Tabs,
  Skeleton,
  SwipeCell,
  Notify,
  ConfigProvider,
} from 'vant'
import 'vant/lib/index.css'
import './styles/index.css'
import App from './App.vue'
import router from './router'

const app = createApp(App)

// Vant components
app.use(Button)
app.use(Cell)
app.use(CellGroup)
app.use(Icon)
app.use(VanImage)
app.use(List)
app.use(NavBar)
app.use(Popup)
app.use(Field)
app.use(Form)
app.use(Toast)
app.use(Dialog)
app.use(ActionSheet)
app.use(Tabbar)
app.use(TabbarItem)
app.use(Badge)
app.use(Loading)
app.use(Empty)
app.use(Search)
app.use(PullRefresh)
app.use(Uploader)
app.use(Tab)
app.use(Tabs)
app.use(Skeleton)
app.use(SwipeCell)
app.use(Notify)
app.use(ConfigProvider)

app.use(createPinia())
app.use(router)

app.mount('#app')
