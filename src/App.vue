<template>
  <div class="max-w-4xl mx-auto space-y-5">
    <!-- Header -->
    <header class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-3">
      <div class="flex flex-wrap gap-2">
        <button @click="openScan" class="btn-primary px-4 py-2 rounded-lg text-sm flex items-center gap-1.5">
          <MagnifyingGlassIcon class="w-4 h-4" /> 扫描局域网
        </button>
        <button @click="openAddDevice" class="btn-secondary px-4 py-2 rounded-lg text-sm flex items-center gap-1.5">
          <PlusIcon class="w-4 h-4" /> 手动添加
        </button>
        <button @click="helpVisible = true" class="p-2 text-slate-500 hover:text-slate-700 rounded-lg hover:bg-slate-100 transition">
          <QuestionMarkCircleIcon class="w-5 h-5" />
        </button>
      </div>
    </header>

    <!-- Device Grid -->
    <div v-if="deviceStore.devices.length === 0" class="text-center py-12 text-slate-400 text-sm">
      暂无设备，点击“扫描局域网”开始发现，或点击“手动添加”。
    </div>
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
      <!-- Group buttons -->
      <div v-if="deviceStore.groups.length > 0" class="col-span-full flex flex-wrap gap-2 mb-1">
        <button
          v-for="g in deviceStore.groups"
          :key="g"
          @click="handleWakeGroup(g)"
          class="btn-purple px-3 py-1.5 rounded-lg text-xs"
        >
          唤醒分组: {{ g }}
        </button>
      </div>

      <DeviceCard
        v-for="device in deviceStore.devices"
        :key="device.mac"
        :device="device"
        @edit="() => openEditDevice(device)"
        @delete="() => confirmDelete(device.mac)"
        @wake="() => handleWake(device.mac)"
        @shutdown="() => openShutdown(device.ip)"
      />
    </div>

    <!-- Modals -->
    <ScanModal :visible="scanVisible" @close="scanVisible = false" @added="scanVisible = false" />
    <DeviceModal :visible="deviceModalVisible" :editing-device="editingDevice" @close="deviceModalVisible = false; editingDevice = null" @saved="deviceModalVisible = false; editingDevice = null" />
    <ShutdownModal :visible="shutdownVisible" :ip="shutdownIp" @close="shutdownVisible = false" @done="shutdownVisible = false" />
    <HelpModal :visible="helpVisible" @close="helpVisible = false" />
    <ConfirmModal :visible="confirmVisible" :message="confirmMessage" :title="confirmTitle" :confirm-type="confirmType" @confirm="confirmResolve(true)" @cancel="confirmResolve(false)" />
    <ToastContainer />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useDeviceStore } from '@/stores/deviceStore'
import { useWol } from '@/composables/useWol'
import { useToast } from '@/composables/useToast'

import DeviceCard from '@/components/DeviceCard.vue'
import ScanModal from '@/components/ScanModal.vue'
import DeviceModal from '@/components/DeviceModal.vue'
import ShutdownModal from '@/components/ShutdownModal.vue'
import HelpModal from '@/components/HelpModal.vue'
import ConfirmModal from '@/components/ConfirmModal.vue'
import ToastContainer from '@/components/ToastContainer.vue'

import {
  MagnifyingGlassIcon,
  PlusIcon,
  QuestionMarkCircleIcon,
} from '@heroicons/vue/24/outline'

const deviceStore = useDeviceStore()
const { wakeDevice, wakeGroup } = useWol()
const { showToast } = useToast()

// 状态
const scanVisible = ref(false)
const deviceModalVisible = ref(false)
const editingDevice = ref(null)
const shutdownVisible = ref(false)
const shutdownIp = ref('')
const helpVisible = ref(false)

// Confirm modal
const confirmVisible = ref(false)
const confirmMessage = ref('')
const confirmTitle = ref('提示')
const confirmType = ref('danger')
let confirmResolveFn = null

const openScan = () => { scanVisible.value = true }
const openAddDevice = () => {
  editingDevice.value = null
  deviceModalVisible.value = true
}
const openEditDevice = (device) => {
  editingDevice.value = device
  deviceModalVisible.value = true
}
const openShutdown = (ip) => {
  shutdownIp.value = ip
  shutdownVisible.value = true
}

// 唤醒后立即刷新
const handleWake = async (mac) => {
  await wakeDevice(mac)
  await deviceStore.fetchDevices()
}

const handleWakeGroup = async (g) => {
  await wakeGroup(g)
  await deviceStore.fetchDevices()
}

// 删除确认
const confirmDelete = (mac) => {
  showConfirm('确定要删除该设备吗？', '删除确认', 'danger')
    .then((ok) => {
      if (ok) {
        deviceStore.deleteDevice(mac)
        showToast('设备已删除', 'info')
      }
    })
}

const showConfirm = (message, title = '提示', type = 'danger') => {
  return new Promise((resolve) => {
    confirmMessage.value = message
    confirmTitle.value = title
    confirmType.value = type
    confirmVisible.value = true
    confirmResolveFn = resolve
  })
}

const confirmResolve = (value) => {
  confirmVisible.value = false
  if (confirmResolveFn) {
    confirmResolveFn(value)
    confirmResolveFn = null
  }
}

onMounted(() => {
  // 初始化加载设备
  deviceStore.fetchDevices()
})
</script>