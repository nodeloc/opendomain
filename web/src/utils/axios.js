import axios from 'axios'

const instance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '',
  timeout: 10000,
})

// 请求拦截器
instance.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
instance.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    if (error.response?.status === 401) {
      // Token 过期或无效
      localStorage.removeItem('token')
      window.location.href = '/login'
    } else if (error.response?.status === 403) {
      // 如果是管理员端点返回403，说明token没有admin权限
      // 检查是否是管理员端点
      if (error.config?.url?.includes('/admin/')) {
        console.error('Admin access denied. Please login with an admin account.')
        // 可选：显示友好提示
        if (window.confirm('You need admin privileges to access this page. Would you like to logout and login with an admin account?')) {
          localStorage.removeItem('token')
          window.location.href = '/login'
        } else {
          window.location.href = '/dashboard'
        }
        return Promise.reject(error)
      }
    }
    return Promise.reject(error)
  }
)

export default instance
