<template>
  <div v-if="visible" class="fixed inset-0 z-50 flex items-center justify-center backdrop-blur-sm">
    <div class="glass-modal w-full max-w-md p-6 text-center">
      <h2 class="text-lg font-semibold text-slate-800 mb-2">{{ title }}</h2>
      <p class="text-sm text-slate-600 mb-6">{{ message }}</p>
      <div class="flex gap-3 justify-center">
        <button @click="cancel" class="btn-secondary px-6 py-2 rounded-lg text-sm font-medium">取消</button>
        <button @click="confirm" class="px-6 py-2 rounded-lg text-sm font-medium" :class="confirmClass">
          {{ confirmText }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'   // 修复：添加 computed 导入

const props = defineProps({
  visible: Boolean,
  title: { type: String, default: '提示' },
  message: { type: String, default: '确认继续吗？' },
  confirmText: { type: String, default: '确认' },
  confirmType: { type: String, default: 'danger' } // 'danger' | 'primary'
})

const emit = defineEmits(['confirm', 'cancel'])

const confirmClass = computed(() => {
  return props.confirmType === 'danger' ? 'btn-danger' : 'btn-primary'
})

const confirm = () => { emit('confirm') }
const cancel = () => { emit('cancel') }
</script>