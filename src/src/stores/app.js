// src/stores/app.js
import { defineStore } from 'pinia'

export const useAppStore = defineStore('app', {
  state: () => ({
    // 侧边栏折叠状态
    isCollapse: false,
    // 系统标题
    systemTitle: 'WorkloadQuery',
    // 当前系统状态: 'ok' | 'warn' | 'error'
    systemStatus: 'ok'
  }),
  actions: {
    toggleSidebar() {
      this.isCollapse = !this.isCollapse
    }
  },
  // 开启持久化：刷新页面后折叠状态依然保留
  persist: {
    key: 'app-state',
    storage: localStorage,
    paths: ['isCollapse'] 
  }
})