<template>
  <div v-if="show" class="fixed inset-0 z-50 overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
    <div class="flex items-center justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:p-0">
      <!-- Background overlay -->
      <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" aria-hidden="true" @click="closeModal"></div>

      <!-- Center spacer -->
      <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>

      <!-- Modal panel -->
      <div class="relative inline-block align-middle bg-white rounded-lg text-left overflow-visible shadow-xl transform transition-all sm:my-8 sm:max-w-2xl sm:w-full">
        <!-- Loading state -->
        <div v-if="loading" class="px-4 py-12 sm:p-12">
          <div class="text-center">
            <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
            <p class="mt-2 text-gray-600">Loading secret...</p>
          </div>
        </div>

        <!-- Error state -->
        <div v-else-if="error" class="px-4 py-5 sm:p-6">
          <div class="rounded-md bg-red-50 p-4">
            <p class="text-sm font-medium text-red-800">{{ error }}</p>
          </div>
          <div class="mt-4 flex justify-end">
            <button @click="closeModal" class="px-4 py-2 text-sm font-medium text-gray-700 hover:text-gray-900">
              Close
            </button>
          </div>
        </div>

        <!-- Secret details -->
        <div v-else-if="secret" class="bg-white">
          <!-- Header -->
          <div class="px-4 py-5 sm:px-6 border-b border-gray-200">
            <div class="flex items-start justify-between">
              <div class="flex-1">
                <h3 class="text-lg leading-6 font-medium text-gray-900" id="modal-title">
                  {{ secret.name }}
                </h3>
                <div class="mt-2 flex items-center space-x-2">
                  <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                    {{ secret.type }}
                  </span>
                  <span v-if="secret.category" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
                    {{ secret.category }}
                  </span>
                </div>
              </div>
              <button @click="closeModal" class="ml-3 text-gray-400 hover:text-gray-500">
                <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
          </div>

          <!-- Body -->
          <div class="px-4 py-5 sm:p-6 space-y-4">
            <!-- Value -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Secret Value
              </label>
              <div class="relative">
                <input
                  :type="showValue ? 'text' : 'password'"
                  :value="secret.value"
                  readonly
                  class="block w-full px-3 py-2 pr-24 border border-gray-300 rounded-md bg-gray-50 text-gray-900 sm:text-sm"
                />
                <div class="absolute inset-y-0 right-0 flex items-center pr-2 space-x-1">
                  <button
                    @click="showValue = !showValue"
                    class="px-2 py-1 text-xs font-medium text-blue-600 hover:text-blue-900"
                    :title="showValue ? 'Hide' : 'Show'"
                  >
                    {{ showValue ? 'Hide' : 'Show' }}
                  </button>
                  <button
                    @click="copyToClipboard"
                    class="px-2 py-1 text-xs font-medium text-blue-600 hover:text-blue-900"
                    title="Copy to clipboard"
                  >
                    {{ copied ? 'Copied!' : 'Copy' }}
                  </button>
                </div>
              </div>
            </div>

            <!-- Tags -->
            <div v-if="secret.tags && secret.tags.length > 0">
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Tags
              </label>
              <div class="flex flex-wrap gap-2">
                <span v-for="tag in secret.tags" :key="tag" class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-green-100 text-green-800">
                  {{ tag }}
                </span>
              </div>
            </div>

            <!-- Metadata -->
            <div v-if="secret.metadata && Object.keys(secret.metadata).length > 0">
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Additional Information
              </label>
              <div class="bg-gray-50 rounded-md p-3 space-y-2">
                <div v-for="(value, key) in secret.metadata" :key="key" class="flex justify-between text-sm">
                  <span class="font-medium text-gray-500 capitalize">{{ key }}:</span>
                  <a v-if="key === 'url'" :href="value" target="_blank" rel="noopener noreferrer" class="text-blue-600 hover:text-blue-900 truncate ml-2">
                    {{ value }}
                  </a>
                  <span v-else class="text-gray-900 truncate ml-2">{{ value }}</span>
                </div>
              </div>
            </div>

            <!-- Timestamps -->
            <div class="grid grid-cols-2 gap-4 pt-4 border-t border-gray-200">
              <div>
                <label class="block text-xs font-medium text-gray-500">Created</label>
                <p class="mt-1 text-sm text-gray-900">{{ formatDate(secret.created_at) }}</p>
              </div>
              <div>
                <label class="block text-xs font-medium text-gray-500">Updated</label>
                <p class="mt-1 text-sm text-gray-900">{{ formatDate(secret.updated_at) }}</p>
              </div>
            </div>
          </div>

          <!-- Footer -->
          <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
            <button
              @click="handleDelete"
              :disabled="deleting"
              class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-red-600 text-base font-medium text-white hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50"
            >
              {{ deleting ? 'Deleting...' : 'Delete' }}
            </button>
            <button
              @click="closeModal"
              :disabled="deleting"
              class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50"
            >
              Close
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useSecretsStore } from '../stores/secrets'

const props = defineProps({
  show: {
    type: Boolean,
    required: true
  },
  secretId: {
    type: String,
    default: null
  }
})

const emit = defineEmits(['close', 'deleted'])

const secretsStore = useSecretsStore()

const secret = ref(null)
const loading = ref(false)
const error = ref('')
const showValue = ref(false)
const copied = ref(false)
const deleting = ref(false)

// Watch for modal open and secretId changes
watch(() => [props.show, props.secretId], async ([newShow, newSecretId]) => {
  if (newShow && newSecretId) {
    await loadSecret(newSecretId)
  } else if (!newShow) {
    resetModal()
  }
}, { immediate: true })

const loadSecret = async (id) => {
  loading.value = true
  error.value = ''
  secret.value = null
  showValue.value = false
  
  try {
    secret.value = await secretsStore.fetchSecret(id)
  } catch (err) {
    error.value = err || 'Failed to load secret'
  } finally {
    loading.value = false
  }
}

const resetModal = () => {
  secret.value = null
  error.value = ''
  showValue.value = false
  copied.value = false
  deleting.value = false
}

const closeModal = () => {
  if (!deleting.value) {
    emit('close')
  }
}

const copyToClipboard = async () => {
  if (secret.value?.value) {
    try {
      await navigator.clipboard.writeText(secret.value.value)
      copied.value = true
      setTimeout(() => {
        copied.value = false
      }, 2000)
    } catch (err) {
      console.error('Failed to copy:', err)
    }
  }
}

const handleDelete = async () => {
  if (!confirm('Are you sure you want to delete this secret? This action cannot be undone.')) {
    return
  }

  deleting.value = true
  try {
    await secretsStore.deleteSecret(props.secretId)
    emit('deleted', props.secretId)
    emit('close')
  } catch (err) {
    error.value = err || 'Failed to delete secret'
  } finally {
    deleting.value = false
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
</script>
