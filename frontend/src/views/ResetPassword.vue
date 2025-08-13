<template>
  <v-container fluid class="fill-height">
    <v-row justify="center" align="center">
      <v-col cols="12" sm="8" md="6" lg="4">
        <v-card class="elevation-12 pa-6">
          <v-card-title class="text-h5 mb-4 text-center">
            Reset Password
          </v-card-title>
          
          <v-alert
            v-if="error"
            type="error"
            class="mb-4"
            closable
            @click:close="error = ''"
          >
            {{ error }}
          </v-alert>
          
          <v-alert
            v-if="success"
            type="success"
            class="mb-4"
            closable
            @click:close="success = ''"
          >
            {{ success }}
          </v-alert>
          
          <!-- Request Password Reset Form -->
          <v-form v-if="!token && !resetSuccess" ref="requestForm" v-model="valid" @submit.prevent="handleRequestReset">
            <v-text-field
              v-model="email"
              label="Email"
              prepend-inner-icon="mdi-email"
              variant="outlined"
              :rules="emailRules"
              required
            ></v-text-field>
            
            <div class="d-flex flex-column gap-4 mt-4">
              <v-btn 
                type="submit"
                color="primary"
                block
                size="large"
                :loading="loading"
              >
                Request Password Reset
              </v-btn>
              
              <v-btn 
                color="secondary"
                variant="outlined"
                block
                :to="{ path: '/' }"
                :disabled="loading"
              >
                Back to Login
              </v-btn>
            </div>
          </v-form>
          
          <!-- Password Reset Form -->
          <v-form v-if="token && !resetSuccess" ref="resetForm" v-model="validReset" @submit.prevent="handleResetPassword">
            <v-text-field
              v-model="password"
              label="New Password"
              prepend-inner-icon="mdi-lock"
              :append-inner-icon="showPassword ? 'mdi-eye' : 'mdi-eye-off'"
              :type="showPassword ? 'text' : 'password'"
              @click:append-inner="showPassword = !showPassword"
              variant="outlined"
              :rules="passwordRules"
              required
            ></v-text-field>
            
            <v-text-field
              v-model="confirmPassword"
              label="Confirm New Password"
              prepend-inner-icon="mdi-lock-check"
              :append-inner-icon="showConfirmPassword ? 'mdi-eye' : 'mdi-eye-off'"
              :type="showConfirmPassword ? 'text' : 'password'"
              @click:append-inner="showConfirmPassword = !showConfirmPassword"
              variant="outlined"
              :rules="[...passwordRules, passwordMatchRule]"
              required
            ></v-text-field>
            
            <div class="d-flex flex-column gap-4 mt-4">
              <v-btn 
                type="submit"
                color="primary"
                block
                size="large"
                :loading="loading"
              >
                Reset Password
              </v-btn>
              
              <v-btn 
                color="secondary"
                variant="outlined"
                block
                :to="{ path: '/' }"
                :disabled="loading"
              >
                Back to Login
              </v-btn>
            </div>
          </v-form>
          
          <!-- Success View -->
          <div v-if="resetSuccess" class="text-center">
            <v-icon color="success" size="64" class="mb-4">mdi-check-circle</v-icon>
            <h3 class="text-h5 mb-4">Password Reset Successfully!</h3>
            <p class="mb-6">Your password has been updated. You can now log in with your new password.</p>
            
            <div class="d-flex justify-center">
              <v-btn 
                color="primary"
                :to="{ path: '/' }"
                size="large"
              >
                Go to Login
              </v-btn>
            </div>
          </div>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

// Router
const router = useRouter()
const route = useRoute()

// Auth store
const authStore = useAuthStore()

// Form data
const valid = ref(false)
const validReset = ref(false)
const requestForm = ref(null)
const resetForm = ref(null)
const loading = ref(false)
const error = ref('')
const success = ref('')
const resetSuccess = ref(false)
const showPassword = ref(false)
const showConfirmPassword = ref(false)

// Form fields
const email = ref('')
const token = ref('')
const password = ref('')
const confirmPassword = ref('')

// Form validation
const emailRules = [
  v => !!v || 'Email is required',
  v => /.+@.+\..+/.test(v) || 'Email must be valid'
]

const passwordRules = [
  v => !!v || 'Password is required',
  v => v.length >= 6 || 'Password must be at least 6 characters'
]

const passwordMatchRule = () => 
  password.value === confirmPassword.value || 'Passwords must match'

// Check for token and email in URL
onMounted(() => {
  token.value = route.query.token || ''
  email.value = route.query.email || ''
  
  // If mobile parameter is present, store it for later redirect
  if (route.query.mobile === 'true') {
    localStorage.setItem('isMobileFlow', 'true')
  }
})

// Handle password reset request
const handleRequestReset = async () => {
  if (!requestForm.value.validate()) return
  
  loading.value = true
  error.value = ''
  success.value = ''
  
  try {
    await authStore.forgotPassword(email.value)
    success.value = 'Password reset instructions have been sent to your email.'
    
    // If coming from mobile app, show success message then redirect back to app
    if (route.query.mobile === 'true') {
      setTimeout(() => {
        window.location.href = `uess://auth-callback?status=reset-requested&email=${encodeURIComponent(email.value)}`
      }, 3000)
    }
  } catch (err) {
    error.value = err.response?.data?.error || 'Failed to send reset instructions. Please try again.'
    console.error('Reset request error:', err)
  } finally {
    loading.value = false
  }
}

// Handle password reset
const handleResetPassword = async () => {
  if (!resetForm.value.validate()) return
  
  loading.value = true
  error.value = ''
  
  try {
    await authStore.resetPassword(token.value, password.value)
    resetSuccess.value = true
    
    // If coming from mobile app, redirect back to app after successful reset
    const isMobileFlow = localStorage.getItem('isMobileFlow') === 'true'
    if (isMobileFlow) {
      localStorage.removeItem('isMobileFlow')
      setTimeout(() => {
        window.location.href = `uess://auth-callback?status=reset-success&email=${encodeURIComponent(email.value)}`
      }, 3000)
    }
  } catch (err) {
    error.value = err.response?.data?.error || 'Password reset failed. The token may be invalid or expired.'
    console.error('Reset error:', err)
  } finally {
    loading.value = false
  }
}
</script>
