import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { apiUrl } from '@/utils/api'
import { useToast } from '@/composables/useToast'

export const useDeviceStore = defineStore('device', () => {
  const devices = ref([])
  const { showToast } = useToast()

  const fetchDevices = async () => {
    try {
      const res = await fetch(apiUrl('/api/devices'))
      if (!res.ok) throw new Error('Network error')
      const data = await res.json()
      // 过滤掉无效设备（MAC 为空）
      devices.value = (data || []).filter(d => d.mac && d.mac.trim() !== '')
    } catch {
      showToast('无法连接至服务器，请检查后端状态', 'error')
    }
  }

  const saveDevice = async (payload) => {
    try {
      await fetch(apiUrl('/api/devices'), {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      })
      await fetchDevices()
    } catch {
      showToast('保存失败', 'error')
    }
  }

  const deleteDevice = async (mac) => {
    try {
      await fetch(apiUrl(`/api/devices?mac=${mac}`), { method: 'DELETE' })
      await fetchDevices()
    } catch {
      showToast('删除失败', 'error')
    }
  }

  // 计算分组
  const groups = computed(() => {
    const gs = devices.value.map(d => d.group).filter(Boolean)
    return [...new Set(gs)]
  })

  return { devices, groups, fetchDevices, saveDevice, deleteDevice }
})