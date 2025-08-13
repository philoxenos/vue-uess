<template>
  <v-app>
    <v-container>
      <v-row justify="center" align="center" class="fill-height">
        <v-col cols="12" sm="8" md="4">
          <v-card class="elevation-12 pa-5">
            <v-card-title class="text-h5 text-center">
              UESS Login
            </v-card-title>
            <v-card-text>
              <!-- Google Sign-In Button -->
              <v-btn
                block
                class="mt-2 mb-4"
                color="error"
                prepend-icon="mdi-google"
                :loading="googleLoading"
                @click="handleGoogleSignIn"
              >
                Continue with Google
              </v-btn>
              
              <v-divider class="my-4">
                <span class="text-overline">OR</span>
              </v-divider>
              
              <!-- Regular Login Form -->
              <v-form ref="form" v-model="valid" @submit.prevent="handleLogin">
                <v-text-field
                  v-model="email"
                  :rules="emailRules"
                  label="Email"
                  required
                  prepend-icon="mdi-email"
                />
                <v-text-field
                  v-model="password"
                  :rules="passwordRules"
                  label="Password"
                  type="password"
                  required
                  prepend-icon="mdi-lock"
                />
                <v-alert v-if="errorMessage" type="error" class="mt-3">
                  {{ errorMessage }}
                </v-alert>
                
                <div class="d-flex align-center justify-space-between my-2">
                  <v-checkbox
                    v-model="rememberMe"
                    label="Remember me"
                    hide-details
                  ></v-checkbox>
                  <v-btn
                    variant="text"
                    color="primary"
                    @click="forgotPassword = true"
                    size="small"
                  >
                    Forgot Password?
                  </v-btn>
                </div>
                
                <v-btn
                  class="mt-4"
                  color="primary"
                  block
                  :disabled="!valid"
                  type="submit"
                  :loading="loading"
                >
                  Login
                </v-btn>
                
                <div class="mt-4 text-center">
                  <span>Don't have an account? </span>
                  <v-btn
                    variant="text"
                    color="primary"
                    @click="registerDialog = true"
                    size="small"
                  >
                    Register
                  </v-btn>
                </div>
              </v-form>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
    
    <!-- Forgot Password Dialog -->
    <v-dialog v-model="forgotPassword" max-width="500px">
      <v-card>
        <v-card-title>Reset Password</v-card-title>
        <v-card-text>
          <p class="mb-4">Enter your email address and we'll send you instructions to reset your password.</p>
          <v-form ref="forgotForm" v-model="forgotValid" @submit.prevent="handleForgotPassword">
            <v-text-field
              v-model="forgotEmail"
              :rules="emailRules"
              label="Email"
              required
              prepend-icon="mdi-email"
            ></v-text-field>
            <v-alert v-if="forgotErrorMessage" type="error" class="mt-3">
              {{ forgotErrorMessage }}
            </v-alert>
            <v-alert v-if="forgotSuccessMessage" type="success" class="mt-3">
              {{ forgotSuccessMessage }}
            </v-alert>
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="grey" text @click="forgotPassword = false">
            Cancel
          </v-btn>
          <v-btn 
            color="primary" 
            @click="handleForgotPassword" 
            :loading="forgotLoading"
            :disabled="!forgotValid || forgotLoading"
          >
            Reset Password
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    
    <!-- Registration Dialog -->
    <v-dialog v-model="registerDialog" max-width="500px">
      <v-card>
        <v-card-title>Create Account</v-card-title>
        <v-card-text>
          <v-form ref="registerForm" v-model="registerValid" @submit.prevent="handleRegister">
            <v-text-field
              v-model="registerEmail"
              :rules="emailRules"
              label="Email"
              required
              prepend-icon="mdi-email"
            ></v-text-field>
            <v-text-field
              v-model="registerFirstName"
              :rules="nameRules"
              label="First Name"
              required
              prepend-icon="mdi-account"
            ></v-text-field>
            <v-text-field
              v-model="registerLastName"
              :rules="nameRules"
              label="Last Name"
              required
              prepend-icon="mdi-account"
            ></v-text-field>
            <v-text-field
              v-model="registerPassword"
              :rules="passwordRules"
              label="Password"
              type="password"
              required
              prepend-icon="mdi-lock"
            ></v-text-field>
            <v-text-field
              v-model="registerConfirmPassword"
              :rules="confirmPasswordRules"
              label="Confirm Password"
              type="password"
              required
              prepend-icon="mdi-lock"
            ></v-text-field>
            <v-alert v-if="registerErrorMessage" type="error" class="mt-3">
              {{ registerErrorMessage }}
            </v-alert>
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="grey" text @click="registerDialog = false">
            Cancel
          </v-btn>
          <v-btn 
            color="primary" 
            @click="handleRegister" 
            :loading="registerLoading"
            :disabled="!registerValid || registerLoading"
          >
            Register
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-app>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

// Login form
const form = ref(null)
const valid = ref(false)
const email = ref('')
const password = ref('')
const rememberMe = ref(false)
const errorMessage = ref('')
const loading = ref(false)
const googleLoading = ref(false)

// Forgot password dialog
const forgotPassword = ref(false)
const forgotForm = ref(null)
const forgotValid = ref(false)
const forgotEmail = ref('')
const forgotLoading = ref(false)
const forgotErrorMessage = ref('')
const forgotSuccessMessage = ref('')

// Registration dialog
const registerDialog = ref(false)
const registerForm = ref(null)
const registerValid = ref(false)
const registerEmail = ref('')
const registerFirstName = ref('')
const registerLastName = ref('')
const registerPassword = ref('')
const registerConfirmPassword = ref('')
const registerLoading = ref(false)
const registerErrorMessage = ref('')

// Validation rules
const emailRules = [
  v => !!v || 'Email is required',
  v => /^\w+([.-]?\w+)*@\w+([.-]?\w+)*(\.\w{2,3})+$/.test(v) || 'Email must be valid'
]

const passwordRules = [
  v => !!v || 'Password is required',
  v => v.length >= 6 || 'Password must be at least 6 characters'
]

const nameRules = [
  v => !!v || 'Name is required'
]

const confirmPasswordRules = computed(() => [
  v => !!v || 'Please confirm your password',
  v => v === registerPassword.value || 'Passwords do not match'
])

// Google Sign-In setup
let googleAuth = null

const loadGoogleAPI = () => {
  const script = document.createElement('script')
  script.src = 'https://accounts.google.com/gsi/client'
  script.async = true
  script.defer = true
  document.head.appendChild(script)
  
  script.onload = initGoogleSignIn
}

const initGoogleSignIn = () => {
  if (!window.google) return
  
  // Initialize Google Sign-In
  window.google.accounts.id.initialize({
    client_id: 'YOUR_GOOGLE_CLIENT_ID', // Replace with your actual Google Client ID
    callback: handleGoogleCredentialResponse,
    auto_select: false,
    cancel_on_tap_outside: true
  })
}

const handleGoogleSignIn = () => {
  googleLoading.value = true
  if (window.google) {
    window.google.accounts.id.prompt((notification) => {
      if (notification.isNotDisplayed() || notification.isSkippedMoment()) {
        // Try to render the button and then click it automatically
        window.google.accounts.id.renderButton(
          document.createElement('div'),
          { theme: 'outline', size: 'large' }
        )
        document.querySelector('.google-signin-button')?.click()
      }
      googleLoading.value = false
    })
  } else {
    errorMessage.value = 'Google Sign-In not available'
    googleLoading.value = false
  }
}

const handleGoogleCredentialResponse = async (response) => {
  googleLoading.value = true
  try {
    // Call your backend with the ID token
    const success = await authStore.loginWithGoogle(response.credential)
    if (success) {
      router.push('/dashboard')
    }
  } catch (error) {
    console.error('Google auth error:', error)
    errorMessage.value = error.response?.data?.error || 'Error authenticating with Google'
  } finally {
    googleLoading.value = false
  }
}

// Form submission handlers
const handleLogin = async () => {
  if (!form.value.validate()) return
  
  loading.value = true
  errorMessage.value = ''
  
  try {
    const success = await authStore.login({
      email: email.value,
      password: password.value
    })
    
    if (success) {
      router.push('/dashboard')
    } else {
      errorMessage.value = 'Login failed'
    }
  } catch (error) {
    console.error('Login error:', error)
    errorMessage.value = error.response?.data?.error || 'An error occurred during login'
  } finally {
    loading.value = false
  }
}

const handleForgotPassword = async () => {
  if (!forgotForm.value?.validate()) return
  
  forgotLoading.value = true
  forgotErrorMessage.value = ''
  forgotSuccessMessage.value = ''
  
  try {
    await authStore.forgotPassword(forgotEmail.value)
    forgotSuccessMessage.value = 'If your email is registered, you will receive password reset instructions'
    setTimeout(() => {
      forgotPassword.value = false
    }, 3000)
  } catch (error) {
    forgotErrorMessage.value = error.response?.data?.error || 'An error occurred'
  } finally {
    forgotLoading.value = false
  }
}

const handleRegister = async () => {
  if (!registerForm.value?.validate()) return
  
  registerLoading.value = true
  registerErrorMessage.value = ''
  
  try {
    const success = await authStore.register({
      email: registerEmail.value,
      first_name: registerFirstName.value,
      last_name: registerLastName.value,
      password: registerPassword.value,
      confirm_password: registerConfirmPassword.value
    })
    
    if (success) {
      registerDialog.value = false
      router.push('/dashboard')
    }
  } catch (error) {
    console.error('Registration error:', error)
    registerErrorMessage.value = error.response?.data?.error || 'An error occurred during registration'
  } finally {
    registerLoading.value = false
  }
}

// Load Google API on component mount
onMounted(() => {
  loadGoogleAPI()
})
</script>
