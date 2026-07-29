<template>
  <div v-if="visible" class="fixed inset-0 z-40 flex items-center justify-center backdrop-blur-sm p-4">
    <div class="glass-modal w-full max-w-lg p-6 max-h-[90vh] overflow-y-auto">
      <h2 class="text-lg font-semibold text-slate-800 mb-4 flex items-center gap-2">
        <LightBulbIcon class="w-5 h-5 text-primary-500" />
        WOL 使用说明
      </h2>

      <div class="space-y-4 text-sm text-slate-600">
        <div>
          <h3 class="font-semibold text-slate-800">1. 如何唤醒电脑 (WOL)</h3>
          <p>目标电脑必须在 BIOS 中开启 <strong>Wake on LAN</strong> (或 PCIe 唤醒)，并且在系统网卡设置中开启“允许此设备唤醒计算机”。必须使用<strong>有线网络</strong>连接。</p>
          <!-- 已删除“关键步骤”说明 -->
        </div>
        <div>
          <h3 class="font-semibold text-slate-800">2. Windows 如何设置远程关机</h3>
          <p>WOL 本身不支持关机，我们通过 SSH 协议安全实现一键关机：</p>
          <ul class="list-disc pl-5 space-y-1 mt-1">
            <li>点击开始菜单处的<strong>搜索栏</strong>，搜索并打开<strong>【可选功能】</strong>。</li>
            <li>在页面上方点击<strong>【添加可选功能/查看功能】</strong>，搜索并安装 <strong>OpenSSH 服务器</strong>。</li>
            <li>按 <code class="bg-slate-100 px-1.5 py-0.5 rounded text-xs">Win + R</code> 输入 <code class="bg-slate-100 px-1.5 py-0.5 rounded text-xs">services.msc</code> 打开服务管理器，找到并启动 <strong>OpenSSH SSH Server</strong>（右键属性，建议启动类型设置为“自动”）。</li>
            <li>在此应用的关机弹窗中，输入你电脑的登陆用户名和密码即可一键关机。</li>
          </ul>
        </div>
        <div>
          <h3 class="font-semibold text-slate-800">3. Linux / macOS 关机</h3>
          <p>大多数 Linux 系统默认已安装 SSH。在关机弹窗中输入具有 <code>sudo</code> 权限的用户名和密码或者 <code>root</code> 密码即可。macOS 需在“共享”中开启“远程登录”。</p>
        </div>
      </div>

      <div class="flex justify-end mt-6 pt-4 border-t border-slate-200">
        <button @click="$emit('close')" class="btn-primary px-6 py-2 rounded-lg text-sm">我知道了</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { LightBulbIcon } from '@heroicons/vue/24/outline'

defineProps({
  visible: Boolean
})

defineEmits(['close'])
</script>