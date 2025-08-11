<template>
  <v-app>
    <v-container>
      <v-row justify="center" align="center" class="fill-height">
        <v-col cols="12" sm="8" md="4">
          <v-card class="elevation-12 pa-5">
            <v-card-title class="text-h5 text-center">
              Admin Login
            </v-card-title>
            <v-card-text>
              <v-form ref="form" v-model="valid" @submit.prevent="login">
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
              </v-form>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </v-app>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

const router = useRouter()
const valid = ref(false)
const email = ref('')
const password = ref('')
const errorMessage = ref('')
const loading = ref(false)

const emailRules = [
  v => !!v || 'Email is required',
  v => /^\w+([.-]?\w+)*@\w+([.-]?\w+)*(\.\w{2,3})+$/.test(v) || 'Email must be valid'
]

const passwordRules = [
  v => !!v || 'Password is required',
  v => v.length >= 6 || 'Password must be at least 6 characters'
]

const login = async () => {
  loading.value = true
  errorMessage.value = ''
  
  try {
    const response = await axios.post('http://localhost:8080/api/v1/auth/login', {
      email: email.value,
      password: password.value
    })
    
    if (response.data.token) {
      localStorage.setItem('token', response.data.token)
      localStorage.setItem('user', JSON.stringify(response.data.user))
      router.push('/dashboard')
    } else {
      errorMessage.value = 'Invalid login response'
    }
  } catch (error) {
    console.error('Login error:', error)
    errorMessage.value = error.response?.data?.error || 'An error occurred during login'
  } finally {
    loading.value = false
  }
}
</script>
