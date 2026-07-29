<template>
  <div v-if="visible" class="fixed inset-0 z-40 flex items-center justify-center backdrop-blur-sm p-4">
    <div class="glass-modal w-full max-w-md p-6">
      <h2 class="text-lg font-semibold text-slate-800 mb-4">{{ isEdit ? '编辑设备' : '添加设备' }}</h2>

      <div class="space-y-4">
        <div>
          <label class="block text-xs font-medium text-slate-600 mb-1">设备名</label>
          <input v-model="form.hostname" type="text" class="w-full rounded-lg border border-slate-200 bg-slate-50 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500/20 focus:border-primary-500" placeholder="示例: 我的电脑" />
        </div>
        <div>
          <label class="block text-xs font-medium text-slate-600 mb-1">IP 地址</label>
          <input v-model="form.ip" type="text" class="w-full rounded-lg border border-slate-200 bg-slate-50 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500/20 focus:border-primary-500" placeholder="192.168.1.x" />
        </div>
        <div>
          <label class="block text-xs font-medium text-slate-600 mb-1">MAC 地址</label>
          <input v-model="form.mac" type="text" class="w-full rounded-lg border border-slate-200 bg-slate-50 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500/20 focus:border-primary-500" placeholder="AA:BB:CC:DD:EE:FF" />
        </div>
        <div>
          <label class="block text-xs font-medium text-slate-600 mb-1">分组 (可选)</label>
          <input v-model="form.group" type="text" class="w-full rounded-lg border border-slate-200 bg-slate-50 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500/20 focus:border-primary-500" placeholder="例如: 办公" />
        </div>
        <div>
          <label class="block text-xs font-medium text-slate-600 mb-1">定时唤醒 (可选)</label>
          <input v-model="form.schedule" type="time" class="w-full rounded-lg border border-slate-200 bg-slate-50 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500/20 focus:border-primary-500" />
        </div>
      </div>

      <div class="flex justify-end gap-3 mt-6 pt-4 border-t border-slate-200">
        <button @click="$emit('close')" class="btn-secondary px-5 py-2 rounded-lg text-sm">取消</button>
        <button @click="save" class="btn-primary px-5 py-2 rounded-lg text-sm">保存设置</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import { useDeviceStore } from '@/stores/deviceStore'
import { useToast } from '@/composables/useToast'

const props = defineProps({
  visible: Boolean,
  editingDevice: Object,
})

const emit = defineEmits(['close', 'saved'])

const deviceStore = useDeviceStore()
const { showToast } = useToast()

const form = ref({
  hostname: '',
  ip: '',
  mac: '',
  group: '',
  schedule: '',
})

const isEdit = computed(() => !!props.editingDevice)

watch(() => props.visible, (val) => {
  if (val && props.editingDevice) {
    const d = props.editingDevice
    form.value = {
      hostname: d.hostname || '',
      ip: d.ip || '',
      mac: d.mac || '',
      group: d.group || '',
      schedule: d.schedule || '',
    }
  } else if (val) {
    form.value = { hostname: '', ip: '', mac: '', group: '', schedule: '' }
  }
}, { immediate: true })

const save = async () => {
  const { hostname, ip, mac, group, schedule } = form.value
  if (!hostname || !ip || !mac) {
    showToast('请填写完整信息', 'error')
    return
  }
  // 编辑时保留原有在线状态，新增时默认为 false
  const online = isEdit.value ? props.editingDevice.online : false
  await deviceStore.saveDevice({
    ip,
    mac,
    hostname,
    group,
    schedule,
    online,
    oldMac: isEdit.value ? props.editingDevice.mac : null
  })
  showToast('设备已保存', 'success')
  // saveDevice 内部已调用 fetchDevices，无需再次调用
  emit('saved')
  emit('close')
}
</script>