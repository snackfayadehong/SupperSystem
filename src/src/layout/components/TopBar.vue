<template>
    <div class="top-bar">
        <div class="left">
            <el-button link @click="appStore.toggleSidebar">
                <el-icon :size="20">
                    <component :is="appStore.isCollapse ? 'Expand' : 'Fold'" />
                </el-icon>
            </el-button>
            <Breadcrumb /> </div>

        <div class="right">
            <el-tooltip :content="isDark ? '切换到明亮模式' : '切换到暗黑模式'" placement="bottom">
                <el-switch
                    v-model="isDark"
                    inline-prompt
                    :active-icon="Moon"
                    :inactive-icon="Sunny"
                    @change="toggleDark"
                    class="dark-switch"
                />
            </el-tooltip>
            
            </div>
    </div>
</template>

<script setup>
import { useAppStore } from "@/stores/app";
import { useDark, useToggle } from "@vueuse/core"; // 需要安装 @vueuse/core
import { Expand, Fold, Moon, Sunny } from "@element-plus/icons-vue";

const appStore = useAppStore();

// 使用 VueUse 提供的暗黑模式工具，它会自动在 html 标签上切换 'dark' 类名
const isDark = useDark();
const toggleDark = useToggle(isDark);
</script>

<style scoped>
.top-bar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    height: 100%;
    padding: 0 16px;
}
.right {
    display: flex;
    align-items: center;
    gap: 20px;
}
.dark-switch {
    --el-switch-on-color: #333;
}
</style>