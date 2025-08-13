<template>
  <v-app>
    <v-navigation-drawer v-model="drawer" app>
      <v-list>
        <v-list-item prepend-avatar="https://randomuser.me/api/portraits/men/78.jpg">
          <template v-slot:append>
            <v-btn
              variant="text"
              icon="mdi-menu-down"
              @click.stop="userMenu = !userMenu"
            ></v-btn>
          </template>
          <v-list-item-title>{{ authStore.user?.firstName || 'User' }}</v-list-item-title>
          <v-list-item-subtitle>{{ authStore.user?.email }}</v-list-item-subtitle>
        </v-list-item>
        
        <v-divider></v-divider>
        
        <v-list v-if="userMenu" density="compact">
          <v-list-item prepend-icon="mdi-account" title="Profile" @click="openProfile"></v-list-item>
          <v-list-item prepend-icon="mdi-logout" title="Logout" @click="logout"></v-list-item>
        </v-list>
      </v-list>

      <v-divider></v-divider>

      <v-list density="compact" nav>
        <v-list-item to="/dashboard" prepend-icon="mdi-view-dashboard" title="Dashboard"></v-list-item>
        <v-list-item v-if="isAdmin" to="/users" prepend-icon="mdi-account-multiple" title="Users"></v-list-item>
      </v-list>
    </v-navigation-drawer>

    <v-app-bar app color="primary" dark>
      <v-app-bar-nav-icon @click="drawer = !drawer"></v-app-bar-nav-icon>
      <v-toolbar-title>UESS System</v-toolbar-title>
      <v-spacer></v-spacer>
      <v-btn icon @click="logout">
        <v-icon>mdi-logout</v-icon>
      </v-btn>
    </v-app-bar>

    <v-main>
      <v-container fluid>
        <slot></slot>
      </v-container>
    </v-main>

    <v-footer app color="primary" class="text-center d-flex justify-center" dark>
      <div>UESS System &copy; {{ new Date().getFullYear() }}</div>
    </v-footer>
    
    <v-dialog v-model="profileDialog" max-width="500">
      <v-card>
        <v-card-title>User Profile</v-card-title>
        <v-card-text>
          <v-list lines="two">
            <v-list-item>
              <template v-slot:prepend>
                <v-icon icon="mdi-email"></v-icon>
              </template>
              <v-list-item-title>Email</v-list-item-title>
              <v-list-item-subtitle>{{ authStore.user?.email }}</v-list-item-subtitle>
            </v-list-item>
            
            <v-list-item>
              <template v-slot:prepend>
                <v-icon icon="mdi-account"></v-icon>
              </template>
              <v-list-item-title>Name</v-list-item-title>
              <v-list-item-subtitle>{{ authStore.user?.firstName }} {{ authStore.user?.lastName }}</v-list-item-subtitle>
            </v-list-item>
            
            <v-list-item>
              <template v-slot:prepend>
                <v-icon icon="mdi-shield-account"></v-icon>
              </template>
              <v-list-item-title>Roles</v-list-item-title>
              <v-list-item-subtitle>{{ authStore.userRoles.join(', ') }}</v-list-item-subtitle>
            </v-list-item>
            
            <v-list-item>
              <template v-slot:prepend>
                <v-icon icon="mdi-google"></v-icon>
              </template>
              <v-list-item-title>Login Method</v-list-item-title>
              <v-list-item-subtitle>
                {{ authStore.user?.googleSub ? 'Google Account' : 'Email & Password' }}
              </v-list-item-subtitle>
            </v-list-item>
          </v-list>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="primary" text @click="profileDialog = false">
            Close
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-app>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()
const drawer = ref(null)
const userMenu = ref(false)
const profileDialog = ref(false)

const isAdmin = computed(() => {
  return authStore.userRoles.includes('admin')
})

const openProfile = () => {
  userMenu.value = false
  profileDialog.value = true
}

const logout = async () => {
  await authStore.logout()
  router.push('/')
}
</script>
