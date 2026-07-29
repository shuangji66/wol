import { apiUrl } from '@/utils/api'
import { useToast } from './useToast'

export function useWol() {
  const { showToast } = useToast()

  const wakeDevice = async (mac) => {
    try {
      const res = await fetch(apiUrl('/api/wake'), {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ mac })
      })
      if (res.ok) {
        showToast('唤醒指令已发送！', 'success')
        return true
      } else {
        showToast('发送失败，请重试。', 'error')
        return false
      }
    } catch {
      showToast('唤醒请求失败', 'error')
      return false
    }
  }

  const wakeGroup = async (group) => {
    try {
      const res = await fetch(apiUrl('/api/wake'), {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ group })
      })
      if (res.ok) {
        showToast(`已发送 ${group} 分组唤醒指令！`, 'success')
        return true
      } else {
        showToast('唤醒分组失败', 'error')
        return false
      }
    } catch {
      showToast('请求失败', 'error')
      return false
    }
  }

  const shutdownDevice = async (ip, user, pass) => {
    try {
      const res = await fetch(apiUrl('/api/shutdown'), {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ ip, user, pass })
      })
      if (res.ok) {
        showToast('关机指令已发送！', 'success')
        return true
      } else {
        const err = await res.text()
        showToast('关机失败: ' + err, 'error')
        return false
      }
    } catch {
      showToast('关机请求失败', 'error')
      return false
    }
  }

  return { wakeDevice, wakeGroup, shutdownDevice }
}