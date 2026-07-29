<template>
  <div class="fixed bottom-6 left-1/2 -translate-x-1/2 z-50 flex flex-col items-center gap-2">
    <div
      v-for="toast in toasts"
      :key="toast.id"
      class="px-5 py-2.5 rounded-xl shadow-lg border-l-4 transition-all duration-500 animate-fade-up"
      :class="{
        'bg-white border-success': toast.type === 'success',
        'bg-white border-danger': toast.type === 'error',
        'bg-white border-primary-500': toast.type === 'info',
      }"
    >
      <span class="text-sm text-slate-700">{{ toast.message }}</span>
    </div>
  </div>
</template>

<script setup>
import { useToast } from '@/composables/useToast'
const { toasts } = useToast()
</script>

<style scoped>
@keyframes fade-up {
  0% { opacity: 0; transform: translateY(20px); }
  100% { opacity: 1; transform: translateY(0); }
}
.animate-fade-up {
  animation: fade-up 0.3s ease-out;
}
</style>