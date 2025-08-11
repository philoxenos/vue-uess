<template>
  <v-app>
    <app-layout>
      <v-container>
        <v-row>
          <v-col cols="12">
            <v-card class="mb-5">
              <v-card-title class="text-h4">
                Dashboard
                <v-spacer></v-spacer>
                <span class="text-subtitle-1">Welcome, {{ currentUser?.firstName || 'Admin' }}</span>
              </v-card-title>
              <v-card-text>
                <v-row>
                  <v-col cols="12" md="4">
                    <v-card class="mx-auto" color="primary">
                      <v-card-text>
                        <div class="text-h4 text-white">{{ userCount }}</div>
                        <div class="text-subtitle-1 text-white">Total Users</div>
                      </v-card-text>
                      <v-card-actions>
                        <v-btn variant="text" to="/users" class="text-white">
                          View All
                          <v-icon end icon="mdi-arrow-right"></v-icon>
                        </v-btn>
                      </v-card-actions>
                    </v-card>
                  </v-col>
                  <v-col cols="12" md="4">
                    <v-card class="mx-auto" color="success">
                      <v-card-text>
                        <div class="text-h4 text-white">{{ activeUserCount }}</div>
                        <div class="text-subtitle-1 text-white">Active Users</div>
                      </v-card-text>
                    </v-card>
                  </v-col>
                  <v-col cols="12" md="4">
                    <v-card class="mx-auto" color="info">
                      <v-card-text>
                        <div class="text-h4 text-white">{{ newUsersCount }}</div>
                        <div class="text-subtitle-1 text-white">New Users (Last 7 Days)</div>
                      </v-card-text>
                    </v-card>
                  </v-col>
                </v-row>
              </v-card-text>
            </v-card>
            
            <v-card>
              <v-card-title class="text-h5">
                Recent Users
                <v-spacer></v-spacer>
                <v-btn color="primary" to="/users/new" prepend-icon="mdi-plus">
                  Add User
                </v-btn>
              </v-card-title>
              <v-card-text>
                <v-table>
                  <thead>
                    <tr>
                      <th>Name</th>
                      <th>Email</th>
                      <th>Status</th>
                      <th>Created</th>
                      <th>Actions</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="user in recentUsers" :key="user.id">
                      <td>{{ `${user.firstName} ${user.lastName}` }}</td>
                      <td>{{ user.email }}</td>
                      <td>
                        <v-chip :color="user.isActive ? 'success' : 'error'">
                          {{ user.isActive ? 'Active' : 'Inactive' }}
                        </v-chip>
                      </td>
                      <td>{{ formatDate(user.createdAt) }}</td>
                      <td>
                        <v-btn icon variant="text" :to="`/users/edit/${user.id}`">
                          <v-icon>mdi-pencil</v-icon>
                        </v-btn>
                      </td>
                    </tr>
                    <tr v-if="recentUsers.length === 0">
                      <td colspan="5" class="text-center">No users found</td>
                    </tr>
                  </tbody>
                </v-table>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </v-container>
    </app-layout>
  </v-app>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import axios from 'axios'
import AppLayout from '../components/AppLayout.vue'

const recentUsers = ref([])
const userCount = ref(0)
const activeUserCount = ref(0)
const newUsersCount = ref(0)
const loading = ref(false)

const currentUser = computed(() => {
  const userStr = localStorage.getItem('user')
  return userStr ? JSON.parse(userStr) : null
})

const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', { 
    year: 'numeric', 
    month: 'short', 
    day: 'numeric' 
  })
}

const fetchDashboardData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    
    // Get recent users
    const usersResponse = await axios.get('http://localhost:8080/api/v1/users', {
      headers: {
        Authorization: `Bearer ${token}`
      }
    })
    
    recentUsers.value = usersResponse.data.data.slice(0, 5)
    userCount.value = usersResponse.data.data.length
    activeUserCount.value = usersResponse.data.data.filter(user => user.isActive).length
    
    // Calculate new users in the last 7 days
    const sevenDaysAgo = new Date()
    sevenDaysAgo.setDate(sevenDaysAgo.getDate() - 7)
    
    newUsersCount.value = usersResponse.data.data.filter(user => {
      const userDate = new Date(user.createdAt)
      return userDate >= sevenDaysAgo
    }).length
    
  } catch (error) {
    console.error('Error fetching dashboard data:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchDashboardData()
})
</script>
