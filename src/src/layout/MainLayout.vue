<template>
    <el-container class="layout-wrapper">
        <el-header height="56px" class="header-container">
            <TopBar />
        </el-header>

        <el-container class="sub-container">
            <el-aside :width="appStore.isCollapse ? '64px' : '220px'" class="aside-container">
                <Sidebar :collapse="appStore.isCollapse" />
            </el-aside>

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

<script setup name="MainLayout">
import { useAppStore } from "@/stores/app";
import TopBar from "./components/TopBar.vue";
import Sidebar from "./components/Sidebar.vue";

const appStore = useAppStore();
</script>

<style scoped>
.layout-wrapper {
    height: 100vh;
    overflow: hidden;
}

.aside-container {
    background-color: #001529;
    transition: width 0.3s cubic-bezier(0.645, 0.045, 0.355, 1);
    overflow: hidden;
    box-shadow: 2px 0 8px rgba(0, 0, 0, 0.15);
}

.main-content {
    background-color: var(--el-bg-color-page); /* 使用 Element Plus 的背景变量以适配暗黑模式 */
    padding: 20px;
    overflow-y: auto;
}

/* 页面切换动画：平滑淡入并轻微横移 */
.fade-transform-enter-active,
.fade-transform-leave-active {
    transition: all 0.3s ease;
}

.fade-transform-enter-from {
    opacity: 0;
    transform: translateX(-20px);
}

.fade-transform-leave-to {
    opacity: 0;
    transform: translateX(20px);
}
</style>