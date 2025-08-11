import { defineStore } from 'pinia'
import axios from 'axios'
import { ref } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(JSON.parse(localStorage.getItem('user') || 'null'))
  const token = ref(localStorage.getItem('token') || null)

  const isAuthenticated = () => !!token.value
  
  const login = async (credentials) => {
    const response = await axios.post('http://localhost:8080/api/v1/auth/login', credentials)
    
    if (response.data.token) {
      token.value = response.data.token
      user.value = response.data.user
      
      localStorage.setItem('token', response.data.token)
      localStorage.setItem('user', JSON.stringify(response.data.user))
      
      // Set default auth header
      axios.defaults.headers.common['Authorization'] = `Bearer ${token.value}`
      
      return true
    }
    
    return false
  }
  
  const logout = () => {
    token.value = null
    user.value = null
    
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    
    // Remove auth header
    delete axios.defaults.headers.common['Authorization']
  }
  
  // Initialize auth header if token exists
  if (token.value) {
    axios.defaults.headers.common['Authorization'] = `Bearer ${token.value}`
  }
  
  return {
    user,
    token,
    isAuthenticated,
    login,
    logout
  }
})
