<template>
  <div v-if="show" class="fixed inset-0 z-50 overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
    <div class="flex items-center justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:p-0">
      <!-- Background overlay -->
      <div class="fixed inset-0 bg-gray-300 opacity-80" aria-hidden="true" @click="closeModal"></div>

      <!-- Center spacer -->
      <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>

      <!-- Modal panel -->
      <div class="relative inline-block align-middle bg-white rounded-2xl text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:max-w-lg sm:w-full">
        <form @submit.prevent="handleSubmit">
          <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-2">
            <div class="sm:flex sm:items-start">
              <div class="mt-3 text-center sm:mt-0 sm:text-left w-full">
                <h3 class="text-lg leading-6 font-medium text-gray-900 mb-4" id="modal-title">
                  Create New Secret
                </h3>
                
                <div class="space-y-4">
                  <!-- Name -->
                  <div>
                    <label for="secret-name" class="block text-sm font-medium text-gray-700">
                      Name <span class="text-red-500">*</span>
                    </label>
                    <input
                      id="secret-name"
                      v-model="formData.name"
                      type="text"
                      name="secret-name"
                      autocomplete="off"
                      required
                      class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                      placeholder="e.g., GitHub Token"
                    />
                  </div>

                  <!-- Type -->
                  <div>
                    <label for="secret-type" class="block text-sm font-medium text-gray-700">
                      Type <span class="text-red-500">*</span>
                    </label>
                    <select
                      id="secret-type"
                      v-model="formData.type"
                      name="secret-type"
                      autocomplete="off"
                      required
                      class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                    >
                      <option value="">Select type...</option>
                      <option value="account">Account (Username + Password)</option>
                      <option value="password">Password</option>
                      <option value="token">Token</option>
                      <option value="api_key">API Key</option>
                      <option value="url">URL</option>
                      <option value="other">Other</option>
                    </select>
                  </div>

                  <!-- Account Type Fields -->
                  <template v-if="formData.type === 'account'">
                    <!-- Username -->
                    <div>
                      <label for="secret-username" class="block text-sm font-medium text-gray-700">
                        Username <span class="text-red-500">*</span>
                      </label>
                      <input
                        id="secret-username"
                        v-model="formData.username"
                        type="text"
                        name="secret-username"
                        autocomplete="off"
                        class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                        placeholder="Enter username"
                      />
                    </div>

                    <!-- Password -->
                    <div>
                      <label for="secret-password" class="block text-sm font-medium text-gray-700">
                        Password <span class="text-red-500">*</span>
                      </label>
                      <input
                        id="secret-password"
                        v-model="formData.password"
                        type="password"
                        name="secret-password"
                        autocomplete="off"
                        class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                        placeholder="Enter password"
                      />
                    </div>
                  </template>

                  <!-- Value (for non-account types) -->
                  <div v-else>
                    <label for="secret-value" class="block text-sm font-medium text-gray-700">
                      Value <span class="text-red-500">*</span>
                    </label>
                    <textarea
                      id="secret-value"
                      v-model="formData.value"
                      name="secret-value"
                      autocomplete="off"
                      required
                      rows="3"
                      class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                      placeholder="Enter your secret value"
                    ></textarea>
                    <p class="mt-1 text-xs text-gray-500">This will be encrypted before storage</p>
                  </div>

                  <!-- Category -->
                  <div>
                    <label for="secret-category" class="block text-sm font-medium text-gray-700">
                      Category
                    </label>
                    <input
                      id="secret-category"
                      v-model="formData.category"
                      type="text"
                      name="secret-category"
                      autocomplete="off"
                      class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                      placeholder="e.g., development, personal, work"
                    />
                  </div>

                  <!-- Notes -->
                  <div>
                    <label for="secret-notes" class="block text-sm font-medium text-gray-700">
                      Notes (optional)
                    </label>
                    <textarea
                      id="secret-notes"
                      v-model="formData.notes"
                      name="secret-notes"
                      rows="2"
                      class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                      placeholder="Add any notes or description about this secret"
                    ></textarea>
                  </div>

                  <!-- Tags -->
                  <div>
                    <label for="secret-tags" class="block text-sm font-medium text-gray-700">
                      Tags
                    </label>
                    <input
                      id="secret-tags"
                      v-model="tagsInput"
                      type="text"
                      name="secret-tags"
                      autocomplete="off"
                      class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                      placeholder="e.g., github, api, production (comma-separated)"
                    />
                    <p class="mt-1 text-xs text-gray-500">Separate tags with commas</p>
                  </div>

                  <!-- Metadata URL (optional) -->
                  <div>
                    <label for="secret-url" class="block text-sm font-medium text-gray-700">
                      Related URL (optional)
                    </label>
                    <input
                      id="secret-url"
                      v-model="metadataUrl"
                      type="url"
                      name="secret-url"
                      autocomplete="url"
                      class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                      placeholder="https://example.com"
                    />
                  </div>

                  <!-- Error message -->
                  <div v-if="error" class="rounded-md bg-red-50 p-4">
                    <p class="text-sm font-medium text-red-800">{{ error }}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Footer buttons -->
          <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
            <button
              type="submit"
              :disabled="loading"
              class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-blue-600 text-base font-medium text-white hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {{ loading ? 'Creating...' : 'Create Secret' }}
            </button>
            <button
              type="button"
              @click="closeModal"
              :disabled="loading"
              class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50"
            >
              Cancel
            </button>
          </div>
        </form>
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
  }
})

const emit = defineEmits(['close', 'created'])

const secretsStore = useSecretsStore()

const formData = ref({
  name: '',
  type: '',
  value: '',
  username: '',
  password: '',
  category: '',
  notes: '',
  tags: [],
  metadata: {}
})

const tagsInput = ref('')
const metadataUrl = ref('')
const loading = ref(false)
const error = ref('')

// Reset form when modal opens
watch(() => props.show, (newVal) => {
  if (newVal) {
    resetForm()
  }
})

const resetForm = () => {
  formData.value = {
    name: '',
    type: '',
    value: '',
    username: '',
    password: '',
    category: '',
    notes: '',
    tags: [],
    metadata: {}
  }
  tagsInput.value = ''
  metadataUrl.value = ''
  error.value = ''
}

const closeModal = () => {
  if (!loading.value) {
    emit('close')
  }
}

const handleSubmit = async () => {
  loading.value = true
  error.value = ''

  try {
    // Parse tags from comma-separated string
    const tags = tagsInput.value
      .split(',')
      .map(tag => tag.trim())
      .filter(tag => tag.length > 0)

    // Build metadata object
    const metadata = {}
    if (metadataUrl.value) {
      metadata.url = metadataUrl.value
    }

    // Prepare value based on type
    let secretValue = formData.value.value
    if (formData.value.type === 'account') {
      // For account type, store username and password as JSON
      if (!formData.value.username || !formData.value.password) {
        error.value = 'Username and password are required for account type'
        loading.value = false
        return
      }
      secretValue = JSON.stringify({
        username: formData.value.username,
        password: formData.value.password
      })
    } else if (!secretValue) {
      error.value = 'Value is required'
      loading.value = false
      return
    }

    // Prepare payload
    const payload = {
      name: formData.value.name,
      type: formData.value.type,
      value: secretValue,
      category: formData.value.category || '',
      notes: formData.value.notes || '',
      tags: tags,
      metadata: metadata
    }

    // Call API
    const newSecret = await secretsStore.createSecret(payload)
    
    // Success
    emit('created', newSecret)
    emit('close')
    resetForm()
  } catch (err) {
    error.value = err || 'Failed to create secret'
  } finally {
    loading.value = false
  }
}
</script>
