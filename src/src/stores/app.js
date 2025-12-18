// src/src/stores/app.js (新创建)
import { defineStore } from 'pinia'

export const useAppStore = defineStore('app', {
  state: () => ({
    isCollapse: false,
    systemStatus: 'ok', // 'ok' | 'warn' | 'error'
    userInfo: {
      name: 'Admin',
      avatar: 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png'
    }
  }),
  persist:true,
  actions: {
    toggleSidebar() {
      this.isCollapse = !this.isCollapse
    },
    updateStatus(status) {
      this.systemStatus = status
    }
  }
})