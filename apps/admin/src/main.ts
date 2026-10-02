import { createApp } from 'vue'
import { Quasar, Notify, Dialog } from 'quasar'
import quasarLang from 'quasar/lang/zh-TW'
import '@quasar/extras/material-icons/material-icons.css'
import 'quasar/dist/quasar.css'
import App from './App.vue'
import { router } from './router'

createApp(App)
  .use(Quasar, {
    plugins: { Notify, Dialog },
    lang: quasarLang,
    config: { brand: { primary: '#3B2A20', secondary: '#426631', accent: '#8DAE78', negative: '#8E3546' } },
  })
  .use(router)
  .mount('#app')
