<template>
  <div v-if="show" class="fixed bottom-4 right-4 z-50 animate-slide-up">
    <div :class="toastClasses" class="rounded-lg shadow-lg p-4 max-w-sm">
      <div class="flex items-start">
        <div class="flex-shrink-0">
          <!-- Success icon -->
          <svg v-if="type === 'success'" class="h-6 w-6 text-green-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <!-- Error icon -->
          <svg v-else-if="type === 'error'" class="h-6 w-6 text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <!-- Info icon -->
          <svg v-else class="h-6 w-6 text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <div class="ml-3 flex-1">
          <p :class="textClasses" class="text-sm font-medium">
            {{ message }}
          </p>
        </div>
        <div class="ml-4 flex-shrink-0 flex">
          <button @click="close" :class="buttonClasses" class="inline-flex rounded-md focus:outline-none focus:ring-2 focus:ring-offset-2">
            <span class="sr-only">Close</span>
            <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
            </svg>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  show: {
    type: Boolean,
    default: false
  },
  type: {
    type: String,
    default: 'info',
    validator: (value) => ['success', 'error', 'info'].includes(value)
  },
  message: {
    type: String,
    required: true
  }
})

const emit = defineEmits(['close'])

const toastClasses = computed(() => {
  switch (props.type) {
    case 'success':
      return 'bg-green-50 border border-green-200'
    case 'error':
      return 'bg-red-50 border border-red-200'
    default:
      return 'bg-blue-50 border border-blue-200'
  }
})

const textClasses = computed(() => {
  switch (props.type) {
    case 'success':
      return 'text-green-800'
    case 'error':
      return 'text-red-800'
    default:
      return 'text-blue-800'
  }
})

const buttonClasses = computed(() => {
  switch (props.type) {
    case 'success':
      return 'text-green-400 hover:text-green-500 focus:ring-green-500'
    case 'error':
      return 'text-red-400 hover:text-red-500 focus:ring-red-500'
    default:
      return 'text-blue-400 hover:text-blue-500 focus:ring-blue-500'
  }
})

const close = () => {
  emit('close')
}
</script>

<style scoped>
@keyframes slide-up {
  from {
    transform: translateY(100%);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

.animate-slide-up {
  animation: slide-up 0.3s ease-out;
}
</style>
