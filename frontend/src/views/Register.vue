<template>
  <v-container fluid class="fill-height">
    <v-row justify="center" align="center">
      <v-col cols="12" sm="8" md="6" lg="4">
        <v-card class="elevation-12 pa-6">
          <v-card-title class="text-h5 mb-2 text-center">
            Create Account
          </v-card-title>
          
          <v-card-subtitle v-if="googleInfo" class="text-center mb-4">
            Continue registration for {{ googleInfo.email }}
          </v-card-subtitle>
          
          <v-alert
            v-if="error"
            type="error"
            class="mb-4"
            closable
            @click:close="error = ''"
          >
            {{ error }}
          </v-alert>
          
          <v-form ref="form" v-model="valid" @submit.prevent="handleRegister">
            <v-text-field
              v-model="firstName"
              :readonly="googleInfo && googleInfo.givenName"
              label="First Name"
              prepend-inner-icon="mdi-account"
              variant="outlined"
              :rules="nameRules"
              required
            ></v-text-field>
            
            <v-text-field
              v-model="lastName"
              :readonly="googleInfo && googleInfo.familyName"
              label="Last Name"
              prepend-inner-icon="mdi-account"
              variant="outlined"
              :rules="nameRules"
              required
            ></v-text-field>
            
            <v-text-field
              v-model="email"
              :readonly="googleInfo"
              label="Email"
              prepend-inner-icon="mdi-email"
              variant="outlined"
              :rules="emailRules"
              required
            ></v-text-field>
            
            <v-text-field
              v-model="password"
              label="Password"
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
              label="Confirm Password"
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
                Create Account
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
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import axios from 'axios'

// Router
const router = useRouter()
const route = useRoute()

// Auth store
const authStore = useAuthStore()

// Form data
const valid = ref(false)
const form = ref(null)
const loading = ref(false)
const error = ref('')
const googleInfo = ref(null)
const showPassword = ref(false)
const showConfirmPassword = ref(false)

// Form fields
const firstName = ref('')
const lastName = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')

// Form validation
const nameRules = [
  v => !!v || 'Name is required',
  v => v.length >= 2 || 'Name must be at least 2 characters'
]

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

// Check for Google token in URL
onMounted(async () => {
  const token = route.query.token
  const isMobile = route.query.mobile === 'true'
  
  if (token) {
    try {
      loading.value = true
      
      // Verify the token with your backend
      const response = await axios.post('http://localhost:8080/api/v1/auth/google/verify', {
        id_token: token
      })
      
      if (response.data && response.data.user) {
        googleInfo.value = response.data.user
        firstName.value = googleInfo.value.givenName || ''
        lastName.value = googleInfo.value.familyName || ''
        email.value = googleInfo.value.email || ''
      }
    } catch (err) {
      error.value = 'Invalid or expired token. Please try again.'
      console.error('Token verification error:', err)
    } finally {
      loading.value = false
    }
  }
})

// Handle registration
const handleRegister = async () => {
  if (!form.value.validate()) return
  
  loading.value = true
  error.value = ''
  
  try {
    let registrationData = {
      firstName: firstName.value,
      lastName: lastName.value,
      email: email.value,
      password: password.value,
    }
    
    // If we have Google info, include the Google ID token
    if (googleInfo.value && route.query.token) {
      registrationData.googleToken = route.query.token
    }
    
    // Register the user
    const result = await authStore.register(registrationData)
    
    if (result) {
      // Check if we're in mobile flow
      if (route.query.mobile === 'true') {
        // Use a special URL scheme to return to the Android app
        window.location.href = `uess://auth-callback?status=success&email=${encodeURIComponent(email.value)}`
      } else {
        // Regular web flow - redirect to dashboard
        router.push('/dashboard')
      }
    }
  } catch (err) {
    error.value = err.response?.data?.error || 'Registration failed. Please try again.'
    console.error('Registration error:', err)
  } finally {
    loading.value = false
  }
}
</script>
