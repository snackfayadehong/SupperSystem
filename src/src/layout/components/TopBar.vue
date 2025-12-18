<template>
    <div class="top-bar">
        <div class="left">
            <el-button text @click="$emit('toggle-menu')" class="fold-btn">
                <el-icon :size="20">
                    <component :is="isCollapse ? Expand : Fold" />
                </el-icon>
            </el-button>
            
            <el-breadcrumb separator="/">
                <transition-group name="breadcrumb">
                    <el-breadcrumb-item v-for="item in breadcrumb" :key="item.path">
                        {{ item.label }}
                    </el-breadcrumb-item>
                </transition-group>
            </el-breadcrumb>
        </div>

        <div class="right">
            <el-tooltip :content="statusText">
                <div class="status-card">
                    <span class="dot" :class="status"></span>
                    <span class="status-label">{{ statusText }}</span>
                </div>
            </el-tooltip>

            <div class="time-box">
                <el-icon><Clock /></el-icon>
                <span>{{ currentTime }}</span>
            </div>

            <el-tooltip content="全屏切换">
                <el-button text @click="toggleFullscreen" class="action-btn">
                    <el-icon><FullScreen /></el-icon>
                </el-button>
            </el-tooltip>

            <el-dropdown trigger="click">
                <span class="user-info">
                    <el-avatar :size="30" :src="user.avatar" />
                    <span class="name">{{ user.name }}</span>
                    <el-icon class="el-icon--right"><ArrowDown /></el-icon>
                </span>

                <template #dropdown>
                    <el-dropdown-menu>
                        <el-dropdown-item :icon="User">个人信息</el-dropdown-item>
                        <el-dropdown-item :icon="Setting">偏好设置</el-dropdown-item>
                        <el-dropdown-item divided @click="logout" :icon="SwitchButton" class="logout-item"> 
                            退出登录 
                        </el-dropdown-item>
                    </el-dropdown-menu>
                </template>
            </el-dropdown>
        </div>
    </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, defineProps, defineEmits } from 'vue';
import { useRoute } from 'vue-router';
import { getBreadcrumb } from "../utils/breadcrumb";
// 确保所有用到的图标都已导入
import { Fold, Expand, FullScreen, Clock, ArrowDown, User, Setting, SwitchButton } from "@element-plus/icons-vue";

const props = defineProps({
    isCollapse: Boolean
});
const emit = defineEmits(['toggle-menu']);

const route = useRoute();
const status = ref("ok"); // ok | warn | error
const currentTime = ref("");
let timer = null;

const user = ref({
    name: "Admin",
    avatar: "https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png"
});

// 计算面包屑
const breadcrumb = computed(() => getBreadcrumb(route.path));

// 系统状态文字映射
const statusText = computed(() => ({
    ok: "系统正常",
    warn: "存在告警",
    error: "系统异常"
}[status.value]));

// 更新时间
const updateTime = () => {
    const now = new Date();
    const pad = n => String(n).padStart(2, "0");
    currentTime.value = `${now.getFullYear()}-${pad(now.getMonth() + 1)}-${pad(now.getDate())} ${pad(now.getHours())}:${pad(now.getMinutes())}:${pad(now.getSeconds())}`;
};

onMounted(() => {
    updateTime();
    timer = setInterval(updateTime, 1000);
});

onUnmounted(() => {
    if (timer) clearInterval(timer); // 记得清除定时器，防止内存泄漏
});

const toggleFullscreen = () => {
    if (!document.fullscreenElement) {
        document.documentElement.requestFullscreen();
    } else {
        document.exitFullscreen();
    }
};

const logout = () => {
    console.log("执行登出逻辑...");
};
</script>

<style scoped>
/* 保持原有样式，仅做微调以支持动画 */
.top-bar {
    height: 100%;
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0 16px;
}

.left, .right {
    display: flex;
    align-items: center;
    gap: 20px;
}

.fold-btn:hover {
    color: #409eff;
    background: transparent;
}

.status-card {
    display: flex;
    align-items: center;
    padding: 4px 12px;
    background: #f5f7fa;
    border-radius: 20px;
    font-size: 13px;
}

.dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    margin-right: 8px;
}

.dot.ok { background: #67c23a; box-shadow: 0 0 4px #67c23a; }
.dot.warn { background: #e6a23c; }
.dot.error { background: #f56c6c; }

.time-box {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 14px;
    color: #606266;
    font-family: monospace;
}

.user-info {
    display: flex;
    align-items: center;
    cursor: pointer;
    padding: 2px 8px;
    border-radius: 4px;
    transition: background 0.2s;
}

.user-info:hover {
    background: #f0f2f5;
}

.name {
    margin: 0 8px;
    font-weight: 500;
}

.logout-item {
    color: #f56c6c;
}

/* 面包屑切换动画 */
.breadcrumb-enter-active,
.breadcrumb-leave-active {
  transition: all 0.5s;
}

.breadcrumb-enter-from,
.breadcrumb-leave-active {
  opacity: 0;
  transform: translateX(20px);
}

.breadcrumb-move {
  transition: all 0.5s;
}

.breadcrumb-leave-active {
  position: absolute;
}
</style>