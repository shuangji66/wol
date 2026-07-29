<template>
  <div v-if="visible" class="fixed inset-0 z-40 flex items-center justify-center backdrop-blur-sm p-4">
    <div class="glass-modal w-full max-w-md p-6">
      <h2 class="text-lg font-semibold text-slate-800 mb-4">远程关机</h2>
      <p class="text-xs text-slate-500 mb-4">将通过 SSH 发送关机指令。Windows 需要开启 OpenSSH 服务器。</p>

      <div class="space-y-4">
        <div>
          <label class="block text-xs font-medium text-slate-600 mb-1">SSH 用户名</label>
          <input v-model="user" type="text" class="w-full rounded-lg border border-slate-200 bg-slate-50 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500/20 focus:border-primary-500" placeholder="例如: root" />
        </div>
        <div>
          <label class="block text-xs font-medium text-slate-600 mb-1">SSH 密码</label>
          <input v-model="pass" type="password" class="w-full rounded-lg border border-slate-200 bg-slate-50 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500/20 focus:border-primary-500" placeholder="输入密码" />
        </div>
        <div class="flex items-center gap-2">
          <input v-model="remember" type="checkbox" id="remember-creds" class="w-4 h-4" />
          <label for="remember-creds" class="text-xs text-slate-600 cursor-pointer">安全记住密码 (存放于本地)</label>
        </div>
      </div>

      <div class="flex justify-end gap-3 mt-6 pt-4 border-t border-slate-200">
        <button @click="$emit('close')" class="btn-secondary px-5 py-2 rounded-lg text-sm">取消</button>
        <button @click="shutdown" class="btn-danger px-5 py-2 rounded-lg text-sm">确认关机</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useWol } from '@/composables/useWol'
import { useToast } from '@/composables/useToast'
import { useDeviceStore } from '@/stores/deviceStore'
import { encryptCreds, decryptCreds } from '@/utils/crypto'

const props = defineProps({
  visible: Boolean,
  ip: String,
})

const emit = defineEmits(['close', 'done'])

const { shutdownDevice } = useWol()
const { showToast } = useToast()
const deviceStore = useDeviceStore()

const user = ref('root')
const pass = ref('')
const remember = ref(false)

watch(() => props.visible, (val) => {
  if (val && props.ip) {
    const key = 'wol_sh_cred_' + props.ip
    const saved = localStorage.getItem(key)
    if (saved) {
      const dec = decryptCreds(saved)
      if (dec) {
        user.value = dec.user
        pass.value = dec.pass
        remember.value = true
      } else {
        user.value = 'root'
        pass.value = ''
        remember.value = false
      }
    } else {
      user.value = 'root'
      pass.value = ''
      remember.value = false
    }
  }
})

const shutdown = async () => {
  if (!user.value || !pass.value) {
    showToast('请输入用户名和密码', 'error')
    return
  }
  if (remember.value) {
    localStorage.setItem('wol_sh_cred_' + props.ip, encryptCreds(user.value, pass.value))
  } else {
    localStorage.removeItem('wol_sh_cred_' + props.ip)
  }

  await shutdownDevice(props.ip, user.value, pass.value)
  await deviceStore.fetchDevices() // 关机后刷新状态
  emit('done')
  emit('close')
}
</script>