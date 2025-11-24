<template>
  <div class="min-h-screen bg-gray-50">
    <nav class="bg-white shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center">
            <router-link to="/dashboard" class="text-gray-600 hover:text-gray-900 mr-4">← Back</router-link>
            <h1 class="text-2xl font-bold text-gray-900">Secrets</h1>
          </div>
          <div class="flex items-center space-x-4">
            <button @click="showCreateModal = true" class="bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700">
              + New Secret
            </button>
          </div>
        </div>
      </div>
    </nav>

    <div class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <div v-if="loading" class="text-center py-12">
        <p class="text-gray-600">Loading secrets...</p>
      </div>
      <div v-else-if="secrets.length === 0" class="text-center py-12">
        <p class="text-gray-600">No secrets yet. Create one to get started!</p>
      </div>
      <div v-else class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
        <div v-for="secret in secrets" :key="secret.id" class="bg-white overflow-hidden shadow rounded-lg p-6">
          <h3 class="text-lg font-medium text-gray-900">{{ secret.name }}</h3>
          <p class="text-sm text-gray-500">{{ secret.type }}</p>
          <p v-if="secret.category" class="text-sm text-gray-500">{{ secret.category }}</p>
          <div class="mt-4 flex space-x-2">
            <router-link :to="`/secrets/${secret.id}`" class="text-blue-600 hover:text-blue-900">View</router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useSecretsStore } from '../stores/secrets'

const secretsStore = useSecretsStore()
const secrets = ref([])
const loading = ref(false)
const showCreateModal = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    const data = await secretsStore.fetchSecrets()
    secrets.value = data.data || []
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
})
</script>
