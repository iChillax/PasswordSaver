<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Navigation -->
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
      <!-- Stats Card -->
      <div class="mb-6">
        <div class="bg-white overflow-hidden shadow rounded-lg">
          <div class="px-4 py-5 sm:p-6">
            <div class="flex items-center justify-between">
              <div>
                <dt class="text-sm font-medium text-gray-500 truncate">Total Secrets</dt>
                <dd class="mt-1 text-3xl font-extrabold text-gray-900">{{ secrets.length }}</dd>
              </div>
              <div class="flex-shrink-0">
                <button 
                  @click="showCreateModal = true" 
                  class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                >
                  <svg class="-ml-1 mr-2 h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
                  </svg>
                  New Secret
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Secrets Section -->
      <div class="bg-white shadow rounded-lg">
        <div class="px-4 py-5 sm:p-6">
          <h2 class="text-lg font-medium text-gray-900 mb-4">Your Secrets</h2>

          <!-- Loading state -->
          <div v-if="loading" class="text-center py-12">
            <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
            <p class="mt-2 text-gray-600">Loading secrets...</p>
          </div>

          <!-- Empty state -->
          <div v-else-if="secrets.length === 0" class="text-center py-12">
            <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
            </svg>
            <h3 class="mt-2 text-sm font-medium text-gray-900">No secrets</h3>
            <p class="mt-1 text-sm text-gray-500">Get started by creating a new secret.</p>
            <div class="mt-6">
              <button 
                @click="showCreateModal = true" 
                class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700"
              >
                <svg class="-ml-1 mr-2 h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
                </svg>
                New Secret
              </button>
            </div>
          </div>

          <!-- Secrets grid -->
          <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            <div 
              v-for="secret in secrets" 
              :key="secret.id" 
              class="bg-white border border-gray-200 rounded-lg overflow-hidden transition-all"
              :class="{ 'ring-2 ring-blue-500': expandedSecretId === secret.id }"
            >
              <!-- Card Header (Always Visible) -->
              <div 
                class="p-4 cursor-pointer hover:bg-gray-50"
                @click.stop="toggleSecret(secret.id)"
              >
                <div class="flex items-start justify-between">
                  <div class="flex-1 min-w-0">
                    <h3 class="text-base font-medium text-gray-900 truncate">{{ secret.name }}</h3>
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
                  <div class="flex items-center space-x-2 ml-3">
                    <!-- Trash Icon -->
                    <button
                      @click.stop="handleQuickDelete(secret.id)"
                      :disabled="deleting && deletingSecretId === secret.id"
                      class="p-1.5 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-md transition-colors disabled:opacity-50"
                      title="Delete secret"
                    >
                      <svg 
                        xmlns="http://www.w3.org/2000/svg" 
                        width="18" 
                        height="18" 
                        viewBox="0 0 24 24" 
                        fill="none" 
                        stroke="currentColor" 
                        stroke-width="2" 
                        stroke-linecap="round" 
                        stroke-linejoin="round"
                        class="lucide lucide-trash-2"
                      >
                        <path d="M3 6h18"/>
                        <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/>
                        <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/>
                        <line x1="10" x2="10" y1="11" y2="17"/>
                        <line x1="14" x2="14" y1="11" y2="17"/>
                      </svg>
                    </button>
                    <!-- Expand Arrow -->
                    <svg 
                      class="h-5 w-5 text-gray-400 transition-transform"
                      :class="{ 'transform rotate-90': expandedSecretId === secret.id }"
                      fill="none" 
                      viewBox="0 0 24 24" 
                      stroke="currentColor"
                    >
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                    </svg>
                  </div>
                </div>
              </div>

              <!-- Expanded Details -->
              <div 
                v-if="expandedSecretId === secret.id"
                class="border-t border-gray-200 bg-gray-50 p-4 space-y-4"
                @click.stop
              >
                <!-- Loading state -->
                <div v-if="loadingDetail" class="text-center py-4">
                  <div class="inline-block animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600"></div>
                  <p class="mt-2 text-sm text-gray-600">Loading details...</p>
                </div>

                <!-- Error state -->
                <div v-else-if="detailError" class="rounded-md bg-red-50 p-3">
                  <p class="text-sm text-red-800">{{ detailError }}</p>
                </div>

                <!-- Secret details -->
                <div v-else-if="secretDetail" class="space-y-4">
                  <!-- Account Type Display -->
                  <template v-if="secretDetail.type === 'account'">
                    <div>
                      <label class="block text-sm font-medium text-gray-700 mb-2">
                        Username
                      </label>
                      <div class="relative">
                        <input
                          :value="getAccountUsername()"
                          readonly
                          type="text"
                          class="block w-full px-3 py-2 pr-20 border border-gray-300 rounded-md bg-white text-gray-900 text-sm"
                        />
                        <div class="absolute inset-y-0 right-0 flex items-center pr-2">
                          <button
                            @click.stop="copyToClipboard(getAccountUsername())"
                            class="px-2 py-1 text-xs font-medium text-blue-600 hover:text-blue-900"
                          >
                            {{ copied ? 'Copied!' : 'Copy' }}
                          </button>
                        </div>
                      </div>
                    </div>

                    <div>
                      <label class="block text-sm font-medium text-gray-700 mb-2">
                        Password
                      </label>
                      <div class="relative">
                        <input
                          :type="showValue ? 'text' : 'password'"
                          :value="getAccountPassword()"
                          readonly
                          class="block w-full px-3 py-2 pr-24 border border-gray-300 rounded-md bg-white text-gray-900 text-sm"
                        />
                        <div class="absolute inset-y-0 right-0 flex items-center pr-2 space-x-1">
                          <button
                            @click.stop="showValue = !showValue"
                            class="px-2 py-1 text-xs font-medium text-blue-600 hover:text-blue-900"
                          >
                            {{ showValue ? 'Hide' : 'Show' }}
                          </button>
                          <button
                            @click.stop="copyToClipboard(getAccountPassword())"
                            class="px-2 py-1 text-xs font-medium text-blue-600 hover:text-blue-900"
                          >
                            {{ copied ? 'Copied!' : 'Copy' }}
                          </button>
                        </div>
                      </div>
                    </div>
                  </template>

                  <!-- Value (for non-account types) -->
                  <div v-else>
                    <label class="block text-sm font-medium text-gray-700 mb-2">
                      Secret Value
                    </label>
                    <div class="relative">
                      <input
                        :type="showValue ? 'text' : 'password'"
                        :value="secretDetail.value"
                        readonly
                        class="block w-full px-3 py-2 pr-24 border border-gray-300 rounded-md bg-white text-gray-900 text-sm"
                      />
                      <div class="absolute inset-y-0 right-0 flex items-center pr-2 space-x-1">
                        <button
                          @click.stop="showValue = !showValue"
                          class="px-2 py-1 text-xs font-medium text-blue-600 hover:text-blue-900"
                        >
                          {{ showValue ? 'Hide' : 'Show' }}
                        </button>
                        <button
                          @click.stop="copyToClipboard(secretDetail.value)"
                          class="px-2 py-1 text-xs font-medium text-blue-600 hover:text-blue-900"
                        >
                          {{ copied ? 'Copied!' : 'Copy' }}
                        </button>
                      </div>
                    </div>
                  </div>

                  <!-- Notes -->
                  <div v-if="secretDetail.notes">
                    <label class="block text-sm font-medium text-gray-700 mb-2">
                      Notes
                    </label>
                    <div class="bg-blue-50 border border-blue-200 rounded-md p-3">
                      <p class="text-sm text-gray-700 whitespace-pre-wrap">{{ secretDetail.notes }}</p>
                    </div>
                  </div>

                  <!-- All Tags -->
                  <div v-if="secretDetail.tags && secretDetail.tags.length > 0">
                    <label class="block text-sm font-medium text-gray-700 mb-2">
                      Tags
                    </label>
                    <div class="flex flex-wrap gap-2">
                      <span v-for="tag in secretDetail.tags" :key="tag" class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-green-100 text-green-800">
                        {{ tag }}
                      </span>
                    </div>
                  </div>

                  <!-- Metadata -->
                  <div v-if="secretDetail.metadata && Object.keys(secretDetail.metadata).length > 0">
                    <label class="block text-sm font-medium text-gray-700 mb-2">
                      Additional Information
                    </label>
                    <div class="bg-white rounded-md p-3 space-y-2 border border-gray-200">
                      <div v-for="(value, key) in secretDetail.metadata" :key="key" class="flex justify-between text-sm">
                        <span class="font-medium text-gray-500 capitalize">{{ key }}:</span>
                        <a v-if="key === 'url'" :href="value" target="_blank" rel="noopener noreferrer" class="text-blue-600 hover:text-blue-900 truncate ml-2">
                          {{ value }}
                        </a>
                        <span v-else class="text-gray-900 truncate ml-2">{{ value }}</span>
                      </div>
                    </div>
                  </div>

                  <!-- Timestamps -->
                  <div class="grid grid-cols-2 gap-4 pt-3 border-t border-gray-200">
                    <div>
                      <label class="block text-xs font-medium text-gray-500">Created</label>
                      <p class="mt-1 text-sm text-gray-900">{{ formatDate(secretDetail.created_at) }}</p>
                    </div>
                    <div>
                      <label class="block text-xs font-medium text-gray-500">Updated</label>
                      <p class="mt-1 text-sm text-gray-900">{{ formatDate(secretDetail.updated_at) }}</p>
                    </div>
                  </div>

                  <!-- Actions -->
                  <div class="flex justify-end space-x-3 pt-3 border-t border-gray-200">
                    <button
                      @click.stop="handleDelete(secret.id)"
                      :disabled="deleting"
                      class="px-4 py-2 text-sm font-medium text-red-600 hover:text-red-900 disabled:opacity-50"
                    >
                      {{ deleting ? 'Deleting...' : 'Delete' }}
                    </button>
                  </div>
                </div>
              </div>
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
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useSecretsStore } from '../stores/secrets'
import CreateSecretModal from '../components/CreateSecretModal.vue'
import Toast from '../components/Toast.vue'

const router = useRouter()
const authStore = useAuthStore()
const secretsStore = useSecretsStore()

const secrets = ref([])
const loading = ref(false)
const showCreateModal = ref(false)
const expandedSecretId = ref(null)
const secretDetail = ref(null)
const loadingDetail = ref(false)
const detailError = ref('')
const showValue = ref(false)
const copied = ref(false)
const deleting = ref(false)
const deletingSecretId = ref(null)
const showToast = ref(false)
const toastType = ref('success')
const toastMessage = ref('')

const loadSecrets = async () => {
  loading.value = true
  try {
    const data = await secretsStore.fetchSecrets(100) // Load up to 100 secrets
    secrets.value = data.data || []
  } catch (err) {
    console.error('Failed to load secrets:', err)
    showToastMessage('error', 'Failed to load secrets')
  } finally {
    loading.value = false
  }
}

const handleSecretCreated = async (newSecret) => {
  // Reload the secrets list to get fresh data from server
  await loadSecrets()
  
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

const toggleSecret = async (id) => {
  if (expandedSecretId.value === id) {
    // Collapse if already expanded
    expandedSecretId.value = null
    secretDetail.value = null
    showValue.value = false
  } else {
    // Expand and load details
    expandedSecretId.value = id
    showValue.value = false
    await loadSecretDetail(id)
  }
}

const loadSecretDetail = async (id) => {
  loadingDetail.value = true
  detailError.value = ''
  secretDetail.value = null
  
  try {
    secretDetail.value = await secretsStore.fetchSecret(id)
  } catch (err) {
    detailError.value = err || 'Failed to load secret details'
  } finally {
    loadingDetail.value = false
  }
}

const copyToClipboard = async (text) => {
  try {
    await navigator.clipboard.writeText(text)
    copied.value = true
    setTimeout(() => {
      copied.value = false
    }, 2000)
  } catch (err) {
    console.error('Failed to copy:', err)
  }
}

const handleQuickDelete = async (id) => {
  if (!confirm('Are you sure you want to delete this secret? This action cannot be undone.')) {
    return
  }

  deleting.value = true
  deletingSecretId.value = id
  
  try {
    await secretsStore.deleteSecret(id)
    
    // Close expanded view if this secret was expanded
    if (expandedSecretId.value === id) {
      expandedSecretId.value = null
      secretDetail.value = null
    }
    
    await loadSecrets()
    showToastMessage('success', 'Secret deleted successfully!')
  } catch (err) {
    showToastMessage('error', err || 'Failed to delete secret')
  } finally {
    deleting.value = false
    deletingSecretId.value = null
  }
}

const handleDelete = async (id) => {
  if (!confirm('Are you sure you want to delete this secret? This action cannot be undone.')) {
    return
  }

  deleting.value = true
  deletingSecretId.value = id
  
  try {
    await secretsStore.deleteSecret(id)
    expandedSecretId.value = null
    secretDetail.value = null
    await loadSecrets()
    showToastMessage('success', 'Secret deleted successfully!')
  } catch (err) {
    showToastMessage('error', err || 'Failed to delete secret')
  } finally {
    deleting.value = false
    deletingSecretId.value = null
  }
}

const formatDate = (dateString) => {
  if (!dateString) return 'N/A'
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// Helper functions for account type
const getAccountData = () => {
  if (!secretDetail.value || secretDetail.value.type !== 'account') {
    return { username: '', password: '' }
  }
  try {
    return JSON.parse(secretDetail.value.value)
  } catch (err) {
    console.error('Failed to parse account data:', err)
    return { username: '', password: '' }
  }
}

const getAccountUsername = () => {
  return getAccountData().username || ''
}

const getAccountPassword = () => {
  return getAccountData().password || ''
}

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}

onMounted(() => {
  loadSecrets()
})
</script>
