/**
 * 时间格式化工具
 * 统一处理时区问题，所有时间数据从后端接收时都是UTC格式
 * 前端根据用户本地时区显示
 */

/**
 * 格式化日期时间为本地时间
 * @param {string|Date} dateString - UTC时间字符串或Date对象
 * @param {object} options - Intl.DateTimeFormat选项
 * @returns {string} 格式化后的本地时间字符串
 */
export function formatDateTime(dateString, options = {}) {
  if (!dateString) return 'N/A'
  
  const date = typeof dateString === 'string' ? new Date(dateString) : dateString
  
  // 检查日期是否有效
  if (isNaN(date.getTime())) return 'Invalid Date'
  
  // 默认选项：显示完整日期和时间
  const defaultOptions = {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false
  }
  
  const finalOptions = { ...defaultOptions, ...options }
  
  return new Intl.DateTimeFormat(undefined, finalOptions).format(date)
}

/**
 * 格式化日期（不含时间）为本地日期
 * @param {string|Date} dateString - UTC时间字符串或Date对象
 * @returns {string} 格式化后的本地日期字符串
 */
export function formatDate(dateString) {
  if (!dateString) return 'N/A'
  
  const date = typeof dateString === 'string' ? new Date(dateString) : dateString
  
  if (isNaN(date.getTime())) return 'Invalid Date'
  
  return new Intl.DateTimeFormat(undefined, {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  }).format(date)
}

/**
 * 格式化为简短的日期时间
 * @param {string|Date} dateString - UTC时间字符串或Date对象
 * @returns {string} 格式化后的简短时间字符串
 */
export function formatShortDateTime(dateString) {
  return formatDateTime(dateString, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false
  })
}

/**
 * 格式化为相对时间（如：2小时前）
 * @param {string|Date} dateString - UTC时间字符串或Date对象
 * @returns {string} 相对时间字符串
 */
export function formatRelativeTime(dateString) {
  if (!dateString) return 'N/A'
  
  const date = typeof dateString === 'string' ? new Date(dateString) : dateString
  
  if (isNaN(date.getTime())) return 'Invalid Date'
  
  const now = new Date()
  const diffInSeconds = Math.floor((now - date) / 1000)
  
  if (diffInSeconds < 60) {
    return `${diffInSeconds}秒前`
  }
  
  const diffInMinutes = Math.floor(diffInSeconds / 60)
  if (diffInMinutes < 60) {
    return `${diffInMinutes}分钟前`
  }
  
  const diffInHours = Math.floor(diffInMinutes / 60)
  if (diffInHours < 24) {
    return `${diffInHours}小时前`
  }
  
  const diffInDays = Math.floor(diffInHours / 24)
  if (diffInDays < 30) {
    return `${diffInDays}天前`
  }
  
  const diffInMonths = Math.floor(diffInDays / 30)
  if (diffInMonths < 12) {
    return `${diffInMonths}个月前`
  }
  
  const diffInYears = Math.floor(diffInMonths / 12)
  return `${diffInYears}年前`
}

/**
 * 格式化为中文友好的日期时间
 * @param {string|Date} dateString - UTC时间字符串或Date对象
 * @returns {string} 中文格式的日期时间
 */
export function formatChineseDateTime(dateString) {
  return formatDateTime(dateString, {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false
  })
}

/**
 * 计算距离到期的天数
 * @param {string|Date} expiryDate - 到期日期
 * @returns {number} 剩余天数（负数表示已过期）
 */
export function daysUntilExpiry(expiryDate) {
  if (!expiryDate) return 0
  
  const expiry = typeof expiryDate === 'string' ? new Date(expiryDate) : expiryDate
  
  if (isNaN(expiry.getTime())) return 0
  
  const now = new Date()
  const diffInMs = expiry - now
  return Math.ceil(diffInMs / (1000 * 60 * 60 * 24))
}

/**
 * 检查日期是否已过期
 * @param {string|Date} dateString - 日期
 * @returns {boolean} 是否已过期
 */
export function isExpired(dateString) {
  if (!dateString) return false
  
  const date = typeof dateString === 'string' ? new Date(dateString) : dateString
  
  if (isNaN(date.getTime())) return false
  
  return date < new Date()
}

/**
 * 格式化为ISO 8601字符串（用于发送到后端）
 * @param {Date} date - Date对象
 * @returns {string} ISO 8601格式的字符串
 */
export function toISO8601(date) {
  return date.toISOString()
}

/**
 * 从UTC字符串创建Date对象
 * @param {string} utcString - UTC时间字符串
 * @returns {Date} Date对象
 */
export function fromUTC(utcString) {
  return new Date(utcString)
}

/**
 * 添加天数到日期
 * @param {string|Date} dateString - 日期
 * @param {number} days - 要添加的天数
 * @returns {Date} 新的Date对象
 */
export function addDays(dateString, days) {
  const date = typeof dateString === 'string' ? new Date(dateString) : new Date(dateString)
  date.setDate(date.getDate() + days)
  return date
}

/**
 * 添加年份到日期
 * @param {string|Date} dateString - 日期
 * @param {number} years - 要添加的年数
 * @returns {Date} 新的Date对象
 */
export function addYears(dateString, years) {
  const date = typeof dateString === 'string' ? new Date(dateString) : new Date(dateString)
  date.setFullYear(date.getFullYear() + years)
  return date
}

// 导出用于Vue模板的默认格式化函数
export default {
  formatDateTime,
  formatDate,
  formatShortDateTime,
  formatRelativeTime,
  formatChineseDateTime,
  daysUntilExpiry,
  isExpired,
  toISO8601,
  fromUTC,
  addDays,
  addYears
}
