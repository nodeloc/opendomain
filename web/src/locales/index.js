import { createI18n } from 'vue-i18n'
import zhCN from './zh-CN'
import enUS from './en-US'

const messages = {
  'zh-CN': zhCN,
  'en-US': enUS
}

// 获取浏览器语言或本地存储的语言
const getLocale = () => {
  const savedLocale = localStorage.getItem('locale')
  if (savedLocale && messages[savedLocale]) {
    return savedLocale
  }
  
  const browserLocale = navigator.language
  if (messages[browserLocale]) {
    return browserLocale
  }
  
  // 默认中文
  return 'zh-CN'
}

const i18n = createI18n({
  legacy: false,
  locale: getLocale(),
  fallbackLocale: 'zh-CN',
  messages,
  globalInjection: true
})

export default i18n
