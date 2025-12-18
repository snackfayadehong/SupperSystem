<template>
    <div class="top-bar">
        <div class="left-box">
            <el-icon class="collapse-btn" @click="appStore.toggleSidebar">
                <component :is="appStore.isCollapse ? 'Expand' : 'Fold'" />
            </el-icon>
            <el-breadcrumb separator="/">
                <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
                <el-breadcrumb-item v-if="$route.meta.title">{{ $route.meta.title }}</el-breadcrumb-item>
            </el-breadcrumb>
        </div>

        <div class="right-box">
            <el-tooltip :content="isDark ? '切换明亮模式' : '切换暗黑模式'">
                <el-switch
                    v-model="isDark"
                    inline-prompt
                    active-icon="Moon"
                    inactive-icon="Sunny"
                    @change="toggleDark"
                    class="dark-switch"
                />
            </el-tooltip>

            <el-button link @click="toggleFullscreen">
                <el-icon :size="18"><FullScreen /></el-icon>
            </el-button>

            <el-dropdown trigger="click">
                <div class="user-avatar">
                    <el-avatar :size="30" src="https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png" />
                    <span class="username">管理员</span>
                </div>
                <template #dropdown>
                    <el-dropdown-menu>
                        <el-dropdown-item icon="User">个人中心</el-dropdown-item>
                        <el-dropdown-item divided icon="SwitchButton">退出登录</el-dropdown-item>
                    </el-dropdown-menu>
                </template>
            </el-dropdown>
        </div>
    </div>
</template>

<script setup>
import { useAppStore } from "@/stores/app";
import { useDark, useToggle, useFullscreen } from "@vueuse/core";
import { Expand, Fold, Moon, Sunny, FullScreen, User, SwitchButton } from "@element-plus/icons-vue";

const appStore = useAppStore();
const isDark = useDark(); // 自动管理 html 的 .dark 类名
const toggleDark = useToggle(isDark);
const { isFullscreen, toggle: toggleFullscreen } = useFullscreen();
</script>

<style scoped>
.top-bar { display: flex;
  justify-content: space-between;
  align-items: center;
  height: 100%; /* 确保填满 56px 的 el-header */
  padding: 0 20px; }
.left-box, .right-box { display: flex; align-items: center; gap: 15px; }
.collapse-btn { font-size: 20px; cursor: pointer; transition: color 0.3s; }
.collapse-btn:hover { color: var(--el-color-primary); }
.user-avatar { display: flex; align-items: center; gap: 8px; cursor: pointer; }
.dark-switch { --el-switch-on-color: #333; }
</style>