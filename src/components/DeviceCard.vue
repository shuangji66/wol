<template>
  <div class="glass-card p-4 relative group transition-all hover:border-slate-300">
    <!-- 编辑 & 删除按钮 - 始终显示 -->
    <div class="absolute top-3 right-3 flex gap-1">
      <button @click="$emit('edit')" class="p-1 text-slate-400 hover:text-primary-500 rounded transition-colors" title="编辑">
        <PencilSquareIcon class="w-4 h-4" />
      </button>
      <button @click="$emit('delete')" class="p-1 text-slate-400 hover:text-danger rounded transition-colors" title="删除">
        <TrashIcon class="w-4 h-4" />
      </button>
    </div>

    <div>
      <h3 class="text-sm font-semibold text-slate-800 truncate">{{ device.hostname || '未识别主机' }}</h3>
      <div class="flex items-center gap-2 mt-1">
        <span class="status-dot" :class="device.online ? 'status-online' : 'status-offline'"></span>
        <span class="text-xs text-slate-500">{{ device.online ? '在线' : '离线' }}</span>
      </div>
      <p class="text-xs text-slate-500 mt-1">IP: {{ device.ip }}</p>
      <p class="text-xs text-slate-500">MAC: {{ device.mac }}</p>
      <p class="text-xs text-slate-500 mt-1">
        分组: <span class="tag">{{ device.group || '默认' }}</span>
      </p>
      <p class="text-xs text-slate-500">
        定时: <span class="tag">{{ device.schedule || '无' }}</span>
      </p>
    </div>

    <div class="mt-3 space-y-1.5">
      <button
        @click="$emit('wake')"
        :disabled="device.online"
        class="w-full py-1.5 rounded-lg text-sm font-medium transition"
        :class="device.online ? 'bg-slate-100 text-slate-400 cursor-not-allowed' : 'bg-primary-50 text-primary-500 hover:bg-primary-500 hover:text-white'"
      >
        {{ device.online ? '已上线' : '唤醒设备' }}
      </button>
      <button
        v-if="device.online"
        @click="$emit('shutdown')"
        class="w-full py-1.5 rounded-lg text-sm font-medium bg-red-50 text-danger hover:bg-danger hover:text-white transition"
      >
        远程关机
      </button>
    </div>
  </div>
</template>

<script setup>
import { PencilSquareIcon, TrashIcon } from '@heroicons/vue/24/outline'

defineProps({
  device: {
    type: Object,
    required: true
  }
})

defineEmits(['edit', 'delete', 'wake', 'shutdown'])
</script>