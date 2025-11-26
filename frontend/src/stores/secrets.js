import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../services/api'

export const useSecretsStore = defineStore('secrets', () => {
  const secrets = ref([])
  const currentSecret = ref(null)
  const loading = ref(false)
  const error = ref(null)

  const fetchSecrets = async (limit = 10, offset = 0) => {
    loading.value = true
    error.value = null
    try {
      const response = await api.get('/secrets', {
        params: { limit, offset }
      })
      secrets.value = response.data.data || []
      return response.data
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to fetch secrets'
      throw error.value
    } finally {
      loading.value = false
    }
  }

  const fetchSecret = async (id) => {
    loading.value = true
    error.value = null
    try {
      const response = await api.get(`/secrets/${id}`)
      currentSecret.value = response.data
      return response.data
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to fetch secret'
      throw error.value
    } finally {
      loading.value = false
    }
  }

  const createSecret = async (secretData) => {
    loading.value = true
    error.value = null
    try {
      const response = await api.post('/secrets', secretData)
      // Don't push to store's secrets array - let the component manage its own list
      return response.data
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to create secret'
      throw error.value
    } finally {
      loading.value = false
    }
  }

  const updateSecret = async (id, secretData) => {
    loading.value = true
    error.value = null
    try {
      const response = await api.put(`/secrets/${id}`, secretData)
      const index = secrets.value.findIndex(s => s.id === id)
      if (index !== -1) {
        secrets.value[index] = response.data
      }
      currentSecret.value = response.data
      return response.data
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to update secret'
      throw error.value
    } finally {
      loading.value = false
    }
  }

  const deleteSecret = async (id) => {
    loading.value = true
    error.value = null
    try {
      await api.delete(`/secrets/${id}`)
      secrets.value = secrets.value.filter(s => s.id !== id)
      if (currentSecret.value?.id === id) {
        currentSecret.value = null
      }
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to delete secret'
      throw error.value
    } finally {
      loading.value = false
    }
  }

  const searchSecrets = async (query, limit = 10, offset = 0) => {
    loading.value = true
    error.value = null
    try {
      const response = await api.get('/secrets/search', {
        params: { q: query, limit, offset }
      })
      return response.data
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to search secrets'
      throw error.value
    } finally {
      loading.value = false
    }
  }

  return {
    secrets,
    currentSecret,
    loading,
    error,
    fetchSecrets,
    fetchSecret,
    createSecret,
    updateSecret,
    deleteSecret,
    searchSecrets
  }
})
