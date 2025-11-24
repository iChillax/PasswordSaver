<template>
  <div class="min-h-screen bg-gray-50">
    <nav class="bg-white shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center">
            <h1 class="text-2xl font-bold text-gray-900">Secrets Manager</h1>
          </div>
          <div class="flex items-center space-x-4">
            <span class="text-gray-700">{{ authStore.user?.email }}</span>
            <router-link to="/settings" class="text-gray-600 hover:text-gray-900">Settings</router-link>
            <button @click="handleLogout" class="text-gray-600 hover:text-gray-900">Logout</button>
          </div>
        </div>
      </div>
    </nav>

    <div class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
        <div class="bg-white overflow-hidden shadow rounded-lg">
          <div class="px-4 py-5 sm:p-6">
            <dt class="text-sm font-medium text-gray-500 truncate">Total Secrets</dt>
            <dd class="mt-1 text-3xl font-extrabold text-gray-900">{{ totalSecrets }}</dd>
          </div>
        </div>
      </div>

      <div class="mt-8">
        <div class="flex justify-between items-center mb-6">
          <h2 class="text-lg font-medium text-gray-900">Your Secrets</h2>
          <router-link to="/secrets" class="bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700">
            View All
          </router-link>
        </div>
        <p class="text-gray-600">Welcome to your Secrets Manager. Click "View All" to manage your secrets.</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()
const totalSecrets = ref(0)

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}
</script>
