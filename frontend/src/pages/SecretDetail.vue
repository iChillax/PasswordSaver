<template>
  <div class="min-h-screen bg-gray-50">
    <nav class="bg-white shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center">
            <router-link to="/dashboard" class="text-gray-600 hover:text-gray-900 mr-4">← Back</router-link>
            <h1 class="text-2xl font-bold text-gray-900">Secret Details</h1>
          </div>
        </div>
      </div>
    </nav>

    <div class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <div v-if="loading" class="text-center py-12">
        <p class="text-gray-600">Loading secret...</p>
      </div>
      <div v-else-if="secret" class="bg-white shadow rounded-lg p-6">
        <div class="grid grid-cols-1 gap-6">
          <div>
            <label class="block text-sm font-medium text-gray-700">Name</label>
            <p class="mt-1 text-gray-900">{{ secret.name }}</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700">Type</label>
            <p class="mt-1 text-gray-900">{{ secret.type }}</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700">Value</label>
            <div class="mt-1 flex items-center space-x-2">
              <input
                :type="showValue ? 'text' : 'password'"
                :value="secret.value"
                readonly
                class="flex-1 px-3 py-2 border border-gray-300 rounded-md bg-gray-50"
              />
              <button @click="showValue = !showValue" class="text-blue-600 hover:text-blue-900">
                {{ showValue ? 'Hide' : 'Show' }}
              </button>
              <button @click="copyToClipboard" class="text-blue-600 hover:text-blue-900">
                Copy
              </button>
            </div>
          </div>
          <div v-if="secret.category">
            <label class="block text-sm font-medium text-gray-700">Category</label>
            <p class="mt-1 text-gray-900">{{ secret.category }}</p>
          </div>
          <div v-if="secret.tags && secret.tags.length" >
            <label class="block text-sm font-medium text-gray-700">Tags</label>
            <div class="mt-1 flex flex-wrap gap-2">
              <span v-for="tag in secret.tags" :key="tag" class="bg-blue-100 text-blue-800 px-3 py-1 rounded-full text-sm">
                {{ tag }}
              </span>
            </div>
          </div>
        </div>
        <div class="mt-6 flex space-x-4">
          <button @click="handleDelete" class="bg-red-600 text-white px-4 py-2 rounded-md hover:bg-red-700">
            Delete
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useSecretsStore } from '../stores/secrets'

const route = useRoute()
const router = useRouter()
const secretsStore = useSecretsStore()

const secret = ref(null)
const loading = ref(false)
const showValue = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    secret.value = await secretsStore.fetchSecret(route.params.id)
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
})

const copyToClipboard = () => {
  navigator.clipboard.writeText(secret.value.value)
  alert('Copied to clipboard!')
}

const handleDelete = async () => {
  if (confirm('Are you sure you want to delete this secret?')) {
    try {
      await secretsStore.deleteSecret(route.params.id)
      router.push('/dashboard')
    } catch (err) {
      console.error(err)
    }
  }
}
</script>
