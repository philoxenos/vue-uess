<script setup>
import { computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from './stores/auth'
import AppLayout from './components/AppLayout.vue'

const route = useRoute()
const authStore = useAuthStore()

const isLoginPage = computed(() => route.path === '/')
const isAuthenticated = computed(() => authStore.isAuthenticated)

// Check if there's a stored session on app mount
onMounted(async () => {
  if (localStorage.getItem('accessToken') && !isAuthenticated.value) {
    await authStore.refreshSession()
  }
})
</script>

<template>
  <!-- Use layout component for authenticated pages -->
  <AppLayout v-if="isAuthenticated && !isLoginPage">
    <router-view />
  </AppLayout>
  
  <!-- Simple view for login page -->
  <router-view v-else />
</template>

<style>
html, body {
  overflow-y: auto;
}
</style>
