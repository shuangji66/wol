import { ref } from 'vue'

const toasts = ref([])

export function useToast() {
  const showToast = (message, type = 'info') => {
    const id = Date.now() + Math.random()
    toasts.value.push({ id, message, type })
    setTimeout(() => {
      const idx = toasts.value.findIndex(t => t.id === id)
      if (idx !== -1) toasts.value.splice(idx, 1)
    }, 3000)
  }

  return { toasts, showToast }
}