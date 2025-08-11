<template>
  <v-app>
    <app-layout>
      <v-container>
        <v-card>
          <v-card-title class="text-h4">
            {{ isEditing ? 'Edit User' : 'Create New User' }}
          </v-card-title>
          <v-card-text>
            <v-form ref="form" v-model="valid" @submit.prevent="saveUser">
              <v-row>
                <v-col cols="12" md="6">
                  <v-text-field
                    v-model="user.firstName"
                    label="First Name"
                    :rules="nameRules"
                  ></v-text-field>
                </v-col>
                <v-col cols="12" md="6">
                  <v-text-field
                    v-model="user.lastName"
                    label="Last Name"
                    :rules="nameRules"
                  ></v-text-field>
                </v-col>
                <v-col cols="12" md="6">
                  <v-text-field
                    v-model="user.email"
                    label="Email"
                    :rules="emailRules"
                    required
                  ></v-text-field>
                </v-col>
                <v-col cols="12" md="6">
                  <v-text-field
                    v-model="user.googleId"
                    label="Google ID"
                    hint="Optional - Used for mobile app login"
                  ></v-text-field>
                </v-col>
                <v-col cols="12" md="6" v-if="!isEditing">
                  <v-text-field
                    v-model="user.password"
                    label="Password"
                    :rules="passwordRules"
                    type="password"
                    required
                  ></v-text-field>
                </v-col>
                <v-col cols="12" md="6" v-if="!isEditing">
                  <v-text-field
                    v-model="passwordConfirm"
                    label="Confirm Password"
                    :rules="confirmPasswordRules"
                    type="password"
                    required
                  ></v-text-field>
                </v-col>
                <v-col cols="12" md="6">
                  <v-switch
                    v-model="user.isActive"
                    label="Active"
                    color="success"
                  ></v-switch>
                </v-col>
                <v-col cols="12" md="6">
                  <v-switch
                    v-model="user.isAdmin"
                    label="Admin Access"
                    color="primary"
                  ></v-switch>
                </v-col>
              </v-row>
              
              <v-alert v-if="errorMessage" type="error" class="mt-3">
                {{ errorMessage }}
              </v-alert>
              
              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn
                  color="grey"
                  variant="outlined"
                  to="/users"
                >
                  Cancel
                </v-btn>
                <v-btn
                  color="primary"
                  type="submit"
                  :loading="loading"
                  :disabled="!valid"
                >
                  {{ isEditing ? 'Update' : 'Create' }}
                </v-btn>
              </v-card-actions>
            </v-form>
          </v-card-text>
        </v-card>
      </v-container>
    </app-layout>
  </v-app>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import AppLayout from '../components/AppLayout.vue'

const route = useRoute()
const router = useRouter()
const valid = ref(false)
const loading = ref(false)
const errorMessage = ref('')
const passwordConfirm = ref('')

const user = ref({
  firstName: '',
  lastName: '',
  email: '',
  password: '',
  googleId: '',
  isActive: true,
  isAdmin: false
})

const isEditing = computed(() => {
  return route.params.id !== undefined
})

const nameRules = [
  v => (v && v.length > 0) || 'Name is required'
]

const emailRules = [
  v => !!v || 'Email is required',
  v => /^\w+([.-]?\w+)*@\w+([.-]?\w+)*(\.\w{2,3})+$/.test(v) || 'Email must be valid'
]

const passwordRules = [
  v => isEditing.value || !!v || 'Password is required',
  v => isEditing.value || v.length >= 6 || 'Password must be at least 6 characters'
]

const confirmPasswordRules = [
  v => isEditing.value || !!v || 'Password confirmation is required',
  v => isEditing.value || v === user.value.password || 'Passwords must match'
]

const fetchUser = async (id) => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const response = await axios.get(`http://localhost:8080/api/v1/users/${id}`, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    })
    
    // Remove password as it's not returned from API
    const userData = response.data.data
    delete userData.password
    
    user.value = userData
  } catch (error) {
    console.error('Error fetching user:', error)
    errorMessage.value = error.response?.data?.error || 'An error occurred while fetching user data'
  } finally {
    loading.value = false
  }
}

const saveUser = async () => {
  loading.value = true
  errorMessage.value = ''
  
  try {
    const token = localStorage.getItem('token')
    
    if (isEditing.value) {
      // Update existing user
      const userData = { ...user.value }
      delete userData.password // Don't send password on update
      
      await axios.put(
        `http://localhost:8080/api/v1/users/${route.params.id}`,
        userData,
        {
          headers: {
            Authorization: `Bearer ${token}`
          }
        }
      )
    } else {
      // Create new user
      await axios.post(
        'http://localhost:8080/api/v1/auth/register',
        user.value,
        {
          headers: {
            Authorization: `Bearer ${token}`
          }
        }
      )
    }
    
    router.push('/users')
  } catch (error) {
    console.error('Error saving user:', error)
    errorMessage.value = error.response?.data?.error || 'An error occurred while saving the user'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  if (isEditing.value) {
    fetchUser(route.params.id)
  }
})
</script>
