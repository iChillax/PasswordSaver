<template>
  <div class="min-h-screen bg-gray-50">
    <nav class="bg-white shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center">
            <router-link to="/dashboard" class="text-gray-600 hover:text-gray-900 mr-4">← Back</router-link>
            <h1 class="text-2xl font-bold text-gray-900">Settings</h1>
          </div>
        </div>
      </div>
    </nav>

    <div class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <div class="bg-white shadow rounded-lg p-6">
        <div class="mb-6">
          <h2 class="text-lg font-medium text-gray-900 mb-4">Profile</h2>
          <div>
            <label class="block text-sm font-medium text-gray-700">Email</label>
            <p class="mt-1 text-gray-900">{{ authStore.user?.email }}</p>
          </div>
        </div>

        <div class="border-t pt-6">
          <h2 class="text-lg font-medium text-gray-900 mb-4">Security</h2>
          <button @click="handleLogout" class="bg-red-600 text-white px-4 py-2 rounded-md hover:bg-red-700">
            Logout
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}
</script>
