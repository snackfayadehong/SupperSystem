<template>
  <el-container class="layout-wrapper">
    <el-aside :width="appStore.isCollapse ? '64px' : '220px'" class="aside-container">
      <Sidebar :collapse="appStore.isCollapse" />
    </el-aside>

    <el-container direction="vertical" class="right-container">
      <el-header height="56px" class="header-navbar">
        <TopBar />
      </el-header>

      <el-main class="main-content">
        <router-view v-slot="{ Component }">
          <transition name="fade-transform" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { useAppStore } from "@/stores/app";
import TopBar from "./components/TopBar.vue";
import Sidebar from "./components/Sidebar.vue";

const appStore = useAppStore();
</script>

<style scoped>
.layout-wrapper {
  height: 100vh;
  width: 100%;
  overflow: hidden;
}

.aside-container {
  background-color: #001529; /* 保持与菜单背景一致 */
  transition: width 0.3s cubic-bezier(0.645, 0.045, 0.355, 1);
  overflow: hidden;
  height: 100%;
}

.right-container {
  height: 100%;
  overflow: hidden;
}

.header-navbar {
  background: var(--el-bg-color);
  border-bottom: 1px solid var(--el-border-color-light);
  padding: 0; /* 让 TopBar 组件内部处理 padding */
  z-index: 10;
  flex-shrink: 0;
  overflow: hidden;
}

.main-content {
  flex: 1;
  background: var(--el-bg-color-page);
  padding: 0 !important;
  overflow-y: auto;
  overflow-x: hidden;
}

/* 页面切换动画 */
.fade-transform-enter-active, .fade-transform-leave-active { transition: all 0.3s; }
.fade-transform-enter-from { opacity: 0; transform: translateX(-20px); }
.fade-transform-leave-to { opacity: 0; transform: translateX(20px); }
</style>