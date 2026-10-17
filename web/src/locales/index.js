import { createI18n } from 'vue-i18n'
import zhCN from './zh-CN'
import enUS from './en-US'

const messages = {
  'zh-CN': zhCN,
  'en-US': enUS
}

const isObject = (v) => v && typeof v === 'object' && !Array.isArray(v)

const deepMerge = (target, source) => {
  for (const key of Object.keys(source)) {
    if (isObject(source[key]) && isObject(target[key])) {
      deepMerge(target[key], source[key])
    } else {
      target[key] = source[key]
    }
  }
  return target
}

// Each file in ./extra exports { 'zh-CN': {...}, 'en-US': {...} } and is merged over the base messages.
const extras = import.meta.glob('./extra/*.js', { eager: true })
for (const path of Object.keys(extras).sort()) {
  const mod = extras[path].default || {}
  for (const lang of Object.keys(messages)) {
    if (mod[lang]) deepMerge(messages[lang], mod[lang])
  }
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
