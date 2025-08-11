<template>
  <v-app>
    <app-layout>
      <v-container>
        <v-card>
          <v-card-title class="text-h4 d-flex">
            Users
            <v-spacer></v-spacer>
            <v-btn color="primary" prepend-icon="mdi-plus" to="/users/new">
              Add User
            </v-btn>
          </v-card-title>
          <v-card-text>
            <v-text-field
              v-model="search"
              label="Search Users"
              prepend-icon="mdi-magnify"
              single-line
              hide-details
              class="mb-4"
            ></v-text-field>
            
            <v-data-table
              :headers="headers"
              :items="filteredUsers"
              :loading="loading"
              :items-per-page="10"
              class="elevation-1"
            >
              <template v-slot:item.name="{ item }">
                {{ `${item.firstName} ${item.lastName}` }}
              </template>
              
              <template v-slot:item.isActive="{ item }">
                <v-chip :color="item.isActive ? 'success' : 'error'">
                  {{ item.isActive ? 'Active' : 'Inactive' }}
                </v-chip>
              </template>
              
              <template v-slot:item.isAdmin="{ item }">
                <v-chip v-if="item.isAdmin" color="primary">
                  Admin
                </v-chip>
                <span v-else>User</span>
              </template>
              
              <template v-slot:item.createdAt="{ item }">
                {{ formatDate(item.createdAt) }}
              </template>
              
              <template v-slot:item.actions="{ item }">
                <v-btn icon variant="text" :to="`/users/edit/${item.id}`">
                  <v-icon>mdi-pencil</v-icon>
                </v-btn>
                <v-btn icon variant="text" color="error" @click="confirmDelete(item)">
                  <v-icon>mdi-delete</v-icon>
                </v-btn>
              </template>
            </v-data-table>
          </v-card-text>
        </v-card>
      </v-container>
      
      <!-- Delete Confirmation Dialog -->
      <v-dialog v-model="deleteDialog" max-width="500">
        <v-card>
          <v-card-title class="text-h5">Confirm Delete</v-card-title>
          <v-card-text>
            Are you sure you want to delete this user? This action cannot be undone.
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="grey" text @click="deleteDialog = false">Cancel</v-btn>
            <v-btn color="error" @click="deleteUser" :loading="deleteLoading">Delete</v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
    </app-layout>
  </v-app>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import axios from 'axios'
import AppLayout from '../components/AppLayout.vue'

const users = ref([])
const loading = ref(false)
const search = ref('')
const deleteDialog = ref(false)
const deleteLoading = ref(false)
const userToDelete = ref(null)

const headers = [
  { title: 'Name', key: 'name', sortable: true },
  { title: 'Email', key: 'email', sortable: true },
  { title: 'Status', key: 'isActive', sortable: true },
  { title: 'Role', key: 'isAdmin', sortable: true },
  { title: 'Created', key: 'createdAt', sortable: true },
  { title: 'Actions', key: 'actions', sortable: false }
]

const filteredUsers = computed(() => {
  if (!search.value) return users.value
  
  const searchLower = search.value.toLowerCase()
  return users.value.filter(user => 
    user.email.toLowerCase().includes(searchLower) || 
    (user.firstName && user.firstName.toLowerCase().includes(searchLower)) ||
    (user.lastName && user.lastName.toLowerCase().includes(searchLower))
  )
})

const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', { 
    year: 'numeric', 
    month: 'short', 
    day: 'numeric' 
  })
}

const fetchUsers = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const response = await axios.get('http://localhost:8080/api/v1/users', {
      headers: {
        Authorization: `Bearer ${token}`
      }
    })
    users.value = response.data.data
  } catch (error) {
    console.error('Error fetching users:', error)
  } finally {
    loading.value = false
  }
}

const confirmDelete = (user) => {
  userToDelete.value = user
  deleteDialog.value = true
}

const deleteUser = async () => {
  if (!userToDelete.value) return
  
  deleteLoading.value = true
  try {
    const token = localStorage.getItem('token')
    await axios.delete(`http://localhost:8080/api/v1/users/${userToDelete.value.id}`, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    })
    
    // Remove user from the list
    users.value = users.value.filter(user => user.id !== userToDelete.value.id)
    deleteDialog.value = false
  } catch (error) {
    console.error('Error deleting user:', error)
  } finally {
    deleteLoading.value = false
  }
}

onMounted(() => {
  fetchUsers()
})
</script>
