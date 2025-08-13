import axios from 'axios'
import { useAuthStore } from '@/stores/auth'

// Create axios instance
const api = axios.create({
  baseURL: 'http://localhost:8080/api/v1'
})

// Add a request interceptor
api.interceptors.request.use(
  config => {
    // Add token to request if available
    const token = localStorage.getItem('accessToken')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => Promise.reject(error)
)

// Add a response interceptor
api.interceptors.response.use(
  response => response,
  async error => {
    const originalRequest = error.config
    
    // If the error is not 401 or the request has already been retried, reject
    if (error.response?.status !== 401 || originalRequest._retry) {
      return Promise.reject(error)
    }
    
    // Mark as retried
    originalRequest._retry = true
    
    try {
      // Try to refresh token
      const authStore = useAuthStore()
      const refreshed = await authStore.refreshSession()
      
      if (refreshed) {
        // Retry the original request with new token
        originalRequest.headers.Authorization = `Bearer ${localStorage.getItem('accessToken')}`
        return api(originalRequest)
      }
      
      // If refresh failed, logout and reject
      await authStore.logout()
      return Promise.reject(error)
    } catch (refreshError) {
      // If refresh fails, logout and reject
      const authStore = useAuthStore()
      await authStore.logout()
      return Promise.reject(refreshError)
    }
  }
)

export default api
