import { defineStore } from 'pinia'
import axios from 'axios'
import { ref, computed } from 'vue'
import { jwtDecode } from 'jwt-decode'

const BASE_URL = 'http://localhost:8080/api/v1'

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref(JSON.parse(localStorage.getItem('user') || 'null'))
  const accessToken = ref(localStorage.getItem('accessToken') || null)
  const refreshToken = ref(localStorage.getItem('refreshToken') || null)
  const loading = ref(false)
  
  // Getters
  const isAuthenticated = computed(() => {
    if (!accessToken.value) return false
    
    try {
      const decoded = jwtDecode(accessToken.value)
      // Check if token is expired
      return decoded.exp * 1000 > Date.now()
    } catch (error) {
      return false
    }
  })
  
  const userRoles = computed(() => {
    if (!user.value) return []
    return user.value.roles || []
  })
  
  const hasRole = (role) => {
    return userRoles.value.includes(role)
  }
  
  // Actions
  const setAuthHeader = () => {
    if (accessToken.value) {
      axios.defaults.headers.common['Authorization'] = `Bearer ${accessToken.value}`
    } else {
      delete axios.defaults.headers.common['Authorization']
    }
  }
  
  const storeTokens = (response) => {
    accessToken.value = response.access_token
    refreshToken.value = response.refresh_token
    user.value = response.user
    
    localStorage.setItem('accessToken', accessToken.value)
    localStorage.setItem('refreshToken', refreshToken.value)
    localStorage.setItem('user', JSON.stringify(user.value))
    
    setAuthHeader()
  }
  
  const login = async (credentials) => {
    loading.value = true
    try {
      const response = await axios.post(`${BASE_URL}/auth/login`, credentials)
      
      if (response.data.access_token) {
        storeTokens(response.data)
        return true
      }
      return false
    } catch (error) {
      console.error('Login error:', error)
      throw error
    } finally {
      loading.value = false
    }
  }
  
  const loginWithGoogle = async (idToken) => {
    loading.value = true
    try {
      const response = await axios.post(`${BASE_URL}/auth/google`, { id_token: idToken })
      
      if (response.data.access_token) {
        storeTokens(response.data)
        return true
      }
      return false
    } catch (error) {
      console.error('Google login error:', error)
      throw error
    } finally {
      loading.value = false
    }
  }
  
  const register = async (userData) => {
    loading.value = true
    try {
      const response = await axios.post(`${BASE_URL}/auth/register`, userData)
      
      if (response.data.access_token) {
        storeTokens(response.data)
        return true
      }
      return false
    } catch (error) {
      console.error('Registration error:', error)
      throw error
    } finally {
      loading.value = false
    }
  }
  
  const refreshSession = async () => {
    if (!refreshToken.value) return false
    
    try {
      const response = await axios.post(`${BASE_URL}/auth/refresh`, {
        refresh_token: refreshToken.value
      })
      
      if (response.data.access_token) {
        storeTokens(response.data)
        return true
      }
      return false
    } catch (error) {
      console.error('Token refresh error:', error)
      logout()
      return false
    }
  }
  
  const logout = async () => {
    try {
      if (refreshToken.value) {
        // Try to revoke the token on server
        await axios.post(`${BASE_URL}/auth/logout`, {
          refresh_token: refreshToken.value
        }).catch(() => {}) // Ignore errors on logout
      }
    } finally {
      // Clear local state regardless of server response
      accessToken.value = null
      refreshToken.value = null
      user.value = null
      
      localStorage.removeItem('accessToken')
      localStorage.removeItem('refreshToken')
      localStorage.removeItem('user')
      
      // Remove auth header
      delete axios.defaults.headers.common['Authorization']
    }
  }
  
  const forgotPassword = async (email) => {
    loading.value = true
    try {
      await axios.post(`${BASE_URL}/auth/forgot-password`, { email })
      return true
    } catch (error) {
      console.error('Forgot password error:', error)
      throw error
    } finally {
      loading.value = false
    }
  }
  
  const resetPassword = async (token, newPassword) => {
    loading.value = true
    try {
      await axios.post(`${BASE_URL}/auth/reset-password`, { 
        token,
        new_password: newPassword
      })
      return true
    } catch (error) {
      console.error('Reset password error:', error)
      throw error
    } finally {
      loading.value = false
    }
  }
  
  // Initialize auth header if token exists
  setAuthHeader()
  
  // If token exists but is expired, try to refresh
  if (accessToken.value && !isAuthenticated.value && refreshToken.value) {
    refreshSession()
  }
  
  return {
    user,
    accessToken,
    loading,
    isAuthenticated,
    userRoles,
    login,
    loginWithGoogle,
    register,
    refreshSession,
    logout,
    forgotPassword,
    resetPassword,
    hasRole
  }
})
