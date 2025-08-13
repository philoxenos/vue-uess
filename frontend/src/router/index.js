import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

// Views
import Dashboard from '../views/Dashboard.vue'
import Login from '../views/Login.vue'
import Register from '../views/Register.vue'
import ResetPassword from '../views/ResetPassword.vue'
import Users from '../views/Users.vue'
import UserForm from '../views/UserForm.vue'

const routes = [
  { 
    path: '/', 
    component: Login,
    meta: { requiresGuest: true }
  },
  { 
    path: '/register', 
    component: Register,
    meta: { requiresGuest: true }
  },
  { 
    path: '/reset-password', 
    component: ResetPassword,
    meta: { requiresGuest: true }
  },
  { 
    path: '/dashboard', 
    component: Dashboard, 
    meta: { 
      requiresAuth: true 
    } 
  },
  { 
    path: '/users', 
    component: Users, 
    meta: { 
      requiresAuth: true,
      roles: ['admin']
    } 
  },
  { 
    path: '/users/new', 
    component: UserForm, 
    meta: { 
      requiresAuth: true,
      roles: ['admin'] 
    } 
  },
  { 
    path: '/users/edit/:id', 
    component: UserForm, 
    meta: { 
      requiresAuth: true,
      roles: ['admin'] 
    } 
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

// Navigation guard
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  
  // Check if token needs refreshing
  if (localStorage.getItem('accessToken') && localStorage.getItem('refreshToken') && !authStore.isAuthenticated) {
    await authStore.refreshSession()
  }
  
  // Handle route access
  if (to.matched.some(record => record.meta.requiresAuth)) {
    // Check if user is authenticated
    if (!authStore.isAuthenticated) {
      next({ path: '/' })
      return
    }
    
    // Check for required roles
    if (to.meta.roles) {
      const hasRequiredRole = to.meta.roles.some(role => 
        authStore.userRoles.includes(role)
      )
      
      if (!hasRequiredRole) {
        next({ path: '/dashboard' })
        return
      }
    }
    
    next()
  } 
  // Handle guest-only routes
  else if (to.matched.some(record => record.meta.requiresGuest)) {
    if (authStore.isAuthenticated) {
      next({ path: '/dashboard' })
      return
    }
    next()
  }
  // Public routes
  else {
    next()
  }
})

export default router
