<template>
    <el-container class="layout-container">
        <el-aside :width="isCollapse ? '64px' : '220px'" class="aside-menu">
            <el-scrollbar>
                <Sidebar :collapse="isCollapse" />
            </el-scrollbar>
        </el-aside>

        <el-container class="main-wrapper">
            <el-header height="56px" class="header-navbar">
                <TopBar @toggle-menu="isCollapse = !isCollapse" :is-collapse="isCollapse" />
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
import { ref } from "vue";
import TopBar from "./components/TopBar.vue";
import Sidebar from "./components/Sidebar.vue";

const isCollapse = ref(false);
</script>

<style scoped>
.layout-container { height: 100vh; overflow: hidden; }
.aside-menu { background: #001529; transition: width 0.3s ease-in-out; }
.main-wrapper { flex-direction: column; background: #f0f2f5; }
.header-navbar { background: #fff; box-shadow: 0 1px 4px rgba(0,21,41,0.08); z-index: 10; }

/* 页面切换动画 */
.fade-transform-enter-active, .fade-transform-leave-active { transition: all 0.3s; }
.fade-transform-enter-from { opacity: 0; transform: translateX(-30px); }
.fade-transform-leave-to { opacity: 0; transform: translateX(30px); }
</style>