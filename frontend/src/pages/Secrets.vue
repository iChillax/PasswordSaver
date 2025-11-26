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
      <!-- Loading state -->
      <div v-if="loading" class="text-center py-12">
        <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        <p class="mt-2 text-gray-600">Loading secrets...</p>
      </div>

      <!-- Empty state -->
      <div v-else-if="secrets.length === 0" class="text-center py-12 bg-white rounded-lg shadow">
        <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
        </svg>
        <h3 class="mt-2 text-sm font-medium text-gray-900">No secrets</h3>
        <p class="mt-1 text-sm text-gray-500">Get started by creating a new secret.</p>
        <div class="mt-6">
          <button @click="showCreateModal = true" class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700">
            + New Secret
          </button>
        </div>
      </div>

      <!-- Secrets grid -->
      <div v-else class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
        <div v-for="secret in secrets" :key="secret.id" class="bg-white overflow-hidden shadow rounded-lg hover:shadow-md transition-shadow">
          <div class="p-6">
            <div class="flex items-start justify-between">
              <div class="flex-1 min-w-0">
                <h3 class="text-lg font-medium text-gray-900 truncate">{{ secret.name }}</h3>
                <div class="mt-2 flex items-center space-x-2">
                  <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                    {{ secret.type }}
                  </span>
                  <span v-if="secret.category" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
                    {{ secret.category }}
                  </span>
                </div>
                <div v-if="secret.tags && secret.tags.length > 0" class="mt-2 flex flex-wrap gap-1">
                  <span v-for="tag in secret.tags.slice(0, 3)" :key="tag" class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-green-100 text-green-800">
                    {{ tag }}
                  </span>
                  <span v-if="secret.tags.length > 3" class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-gray-100 text-gray-600">
                    +{{ secret.tags.length - 3 }}
                  </span>
                </div>
              </div>
            </div>
            <div class="mt-4 flex space-x-3">
              <router-link :to="`/secrets/${secret.id}`" class="text-sm font-medium text-blue-600 hover:text-blue-900">
                View Details →
              </router-link>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Create Secret Modal -->
    <CreateSecretModal 
      :show="showCreateModal" 
      @close="showCreateModal = false"
      @created="handleSecretCreated"
    />

    <!-- Toast Notification -->
    <Toast 
      :show="showToast"
      :type="toastType"
      :message="toastMessage"
      @close="showToast = false"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useSecretsStore } from '../stores/secrets'
import CreateSecretModal from '../components/CreateSecretModal.vue'
import Toast from '../components/Toast.vue'

const secretsStore = useSecretsStore()
const secrets = ref([])
const loading = ref(false)
const showCreateModal = ref(false)
const showToast = ref(false)
const toastType = ref('success')
const toastMessage = ref('')

const loadSecrets = async () => {
  loading.value = true
  try {
    const data = await secretsStore.fetchSecrets()
    secrets.value = data.data || []
  } catch (err) {
    console.error('Failed to load secrets:', err)
    showToastMessage('error', 'Failed to load secrets')
  } finally {
    loading.value = false
  }
}

const handleSecretCreated = (newSecret) => {
  // Add the new secret to the list
  secrets.value.push(newSecret)
  
  // Show success toast
  showToastMessage('success', 'Secret created successfully!')
}

const showToastMessage = (type, message) => {
  toastType.value = type
  toastMessage.value = message
  showToast.value = true
  
  // Auto-hide after 3 seconds
  setTimeout(() => {
    showToast.value = false
  }, 3000)
}

onMounted(() => {
  loadSecrets()
})
</script>
