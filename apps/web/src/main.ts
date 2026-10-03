import { createApp } from 'vue'
import { Quasar, Dialog, Notify } from 'quasar'
import quasarLang from 'quasar/lang/zh-TW'
import '@quasar/extras/material-icons/material-icons.css'
import 'quasar/dist/quasar.css'
import './styles.css'
import App from './App.vue'
import { router } from './router'

createApp(App)
  .use(Quasar, {
    plugins: { Notify, Dialog },
    lang: quasarLang,
    config: {
      // 深咖啡 × 淺綠（同招募頁 brand/coffee-green 分支）
      brand: { primary: '#3B2A20', secondary: '#426631', accent: '#8DAE78', positive: '#426631', negative: '#8E3546' },
    },
  })
  .use(router)
  .mount('#app')
