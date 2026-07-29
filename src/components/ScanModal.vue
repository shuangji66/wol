<template>
  <div v-if="visible" class="fixed inset-0 z-40 flex items-center justify-center backdrop-blur-sm p-4">
    <div class="glass-modal w-full max-w-2xl p-6 max-h-[90vh] flex flex-col m-4">
      <h2 class="text-lg font-semibold text-slate-800 mb-4">扫描结果</h2>

      <div class="flex-1 overflow-y-auto">
        <div v-if="scanning" class="text-center py-8">
          <div class="spinner mx-auto"></div>
          <p class="text-sm text-slate-500 mt-4">扫描中...预计需要 15 秒</p>
        </div>
        <div v-else-if="scanResults.length === 0" class="text-center py-8 text-slate-500 text-sm">
          没有发现新设备。
        </div>
        <div v-else>
          <div class="flex justify-start mb-2">
            <label class="flex items-center gap-1.5 text-sm text-slate-600 cursor-pointer">
              <input type="checkbox" v-model="selectAll" /> 全选 / 反选
            </label>
          </div>
          <div class="border border-slate-200 rounded-lg overflow-y-auto max-h-64">
            <table class="w-full text-sm">
              <thead class="bg-slate-50 sticky top-0">
                <tr>
                  <th class="px-4 py-2 text-center font-medium text-slate-600 w-16">选择</th>
                  <th class="px-4 py-2 text-left font-medium text-slate-600">主机名 / IP</th>
                  <th class="px-4 py-2 text-left font-medium text-slate-600">MAC 地址</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="dev in scanResults" :key="dev.mac" class="border-t border-slate-100 hover:bg-primary-50">
                  <td class="px-4 py-2 text-center">
                    <input type="checkbox" v-model="selectedMacs" :value="dev.mac" class="w-4 h-4" />
                  </td>
                  <td class="px-4 py-2">
                    <div class="font-medium text-slate-700">{{ dev.hostname || '未识别主机' }}</div>
                    <div class="text-xs text-slate-400">{{ dev.ip }}</div>
                  </td>
                  <td class="px-4 py-2 text-slate-500 text-xs">{{ dev.mac }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <div class="flex justify-end gap-3 mt-4 pt-4 border-t border-slate-200">
        <button @click="$emit('close')" class="btn-secondary px-5 py-2 rounded-lg text-sm">取消</button>
        <button
          @click="addSelected"
          :disabled="selectedMacs.length === 0"
          class="btn-primary px-5 py-2 rounded-lg text-sm disabled:opacity-50 disabled:cursor-not-allowed"
        >
          添加所选 ({{ selectedMacs.length }})
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useDeviceStore } from '@/stores/deviceStore'
import { useToast } from '@/composables/useToast'
import { apiUrl } from '@/utils/api'

const props = defineProps({
  visible: Boolean,
})

const emit = defineEmits(['close', 'added'])

const deviceStore = useDeviceStore()
const { showToast } = useToast()
const scanning = ref(true)
const scanResults = ref([])
const selectedMacs = ref([])

const selectAll = computed({
  get: () => selectedMacs.value.length === scanResults.value.length && scanResults.value.length > 0,
  set: (val) => {
    selectedMacs.value = val ? scanResults.value.map(d => d.mac) : []
  }
})

watch(() => props.visible, async (visible) => {
  if (visible) {
    scanning.value = true
    scanResults.value = []
    selectedMacs.value = []
    try {
      const res = await fetch(apiUrl('/api/scan'), { method: 'POST' })
      const data = await res.json()
      scanResults.value = data || []
    } catch {
      showToast('扫描失败，请检查设置', 'error')
    } finally {
      scanning.value = false
    }
  }
})

const addSelected = async () => {
  const selectedDevices = scanResults.value.filter(d => selectedMacs.value.includes(d.mac))
  for (const dev of selectedDevices) {
    await deviceStore.saveDevice({
      ip: dev.ip,
      mac: dev.mac,
      hostname: dev.hostname || '',
      group: '',
      schedule: '',
      online: !!dev.online,
      oldMac: null
    })
  }
  showToast(`成功添加 ${selectedDevices.length} 个设备！`, 'success')
  // saveDevice 内部已调用 fetchDevices，无需再调用
  emit('added')
  emit('close')
}
</script>

<style scoped>
.spinner {
  width: 32px;
  height: 32px;
  border: 3px solid #e2e8f0;
  border-top: 3px solid #3b82f6;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}
@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
</style>