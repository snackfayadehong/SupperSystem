<template>
    <div class="home-dashboard">
        <header class="hero-section">
            <div class="hero-content">
                <h1 class="welcome-text">{{ welcomeMessage }}</h1>
                <p class="date-text">{{ currentTime.format('YYYY年MM月DD日 dddd HH:mm:ss') }}</p>
                <div class="system-status-pill">
                    <span class="pulse-dot"></span>
                    后端 Go 服务：已连接 (172.21.1.75)
                </div>
            </div>
        </header>

        <div class="dashboard-grid">
            
            <section class="stats-row">
                <div v-for="item in stats" :key="item.title" class="modern-stat-card">
                    <div class="card-icon" :style="{ background: item.color + '15', color: item.color }">
                        <el-icon :size="24"><component :is="item.icon" /></el-icon>
                    </div>
                    <div class="card-body">
                        <span class="label">{{ item.title }}</span>
                        <h2 class="value">{{ item.value }}</h2>
                        <div class="trend-tag" :style="{ background: item.color + '10', color: item.color }">
                            {{ item.trend }}
                        </div>
                    </div>
                </div>
            </section>

            <div class="bottom-layout">
                
                <el-card shadow="never" class="glass-panel quick-panel">
                    <template #header>
                        <div class="panel-header">
                            <el-icon><Menu /></el-icon>
                            <span>快捷入口</span>
                        </div>
                    </template>
                    <div class="entry-grid">
                        <div class="entry-item" @click="$router.push('/workload')">
                            <div class="icon-circle blue"><el-icon><Monitor /></el-icon></div>
                            <span class="entry-label">工作量查询</span>
                        </div>
                        <div class="entry-item">
                            <div class="icon-circle green"><el-icon><Document /></el-icon></div>
                            <span class="entry-label">账单对账</span>
                        </div>
                        <div class="entry-item">
                            <div class="icon-circle orange"><el-icon><Management /></el-icon></div>
                            <span class="entry-label">采购汇总</span>
                        </div>
                    </div>
                </el-card>

                <el-card shadow="never" class="glass-panel activity-panel">
                    <template #header>
                        <div class="panel-header">
                            <el-icon><Timer /></el-icon>
                            <span>系统日志</span>
                        </div>
                    </template>
                    <el-scrollbar height="220px">
                        <el-timeline>
                            <el-timeline-item timestamp="2025-12-18" type="primary" :hollow="true">
                                完成 WorkloadPage 模块组件化重构
                            </el-timeline-item>
                            <el-timeline-item timestamp="2025-12-17">
                                优化后端 WorkloadQuery 接口查询性能
                            </el-timeline-item>
                            <el-timeline-item timestamp="2025-12-16">
                                财务对账系统数据导入完成
                            </el-timeline-item>
                        </el-timeline>
                    </el-scrollbar>
                </el-card>

            </div>
        </div>
    </div>
</template>

<script setup>
// 按需导入所有用到的图标，确保 :is 绑定正常
import { 
    Monitor, Document, Menu, Timer, Management, 
    Search, DataAnalysis, WarnTriangleFilled, User 
} from "@element-plus/icons-vue";
import { useHomeData } from "./composables/useHomeData";

const { stats, welcomeMessage, currentTime } = useHomeData();
</script>

<style scoped>
/* 容器布局 */
.home-dashboard {
    padding: 32px;
    background: #f6f8fb;
    min-height: calc(100vh - 56px);
    display: flex;
    flex-direction: column;
    gap: 32px; /* 基础间距 */
}

/* 欢迎看板 */
.hero-section {
    padding: 48px;
    background: linear-gradient(135deg, #4158D0 0%, #C850C0 46%, #FFCC70 100%);
    border-radius: 24px;
    color: white;
    box-shadow: 0 20px 40px rgba(65, 88, 208, 0.2);
}
.welcome-text { font-size: 36px; font-weight: 800; margin: 0; }
.date-text { font-size: 16px; opacity: 0.9; margin: 12px 0 24px; font-family: 'PingFang SC', monospace; }
.system-status-pill {
    display: inline-flex;
    align-items: center;
    background: rgba(255, 255, 255, 0.2);
    padding: 8px 16px;
    border-radius: 50px;
    font-size: 14px;
    backdrop-filter: blur(10px);
}

/* 指标卡片 */
.dashboard-grid {
    display: flex;
    flex-direction: column;
    gap: 40px; /* 增加指标行与下方布局的间距 */
}
.stats-row {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
    gap: 24px;
}
.modern-stat-card {
    background: white;
    padding: 24px;
    border-radius: 20px;
    display: flex;
    align-items: center;
    gap: 20px;
    box-shadow: 0 4px 20px rgba(0,0,0,0.02);
    transition: all 0.3s ease;
}
.modern-stat-card:hover {
    transform: translateY(-10px);
    box-shadow: 0 15px 30px rgba(0,0,0,0.08);
}
.card-icon { width: 60px; height: 60px; border-radius: 16px; display: flex; align-items: center; justify-content: center; }
.label { font-size: 14px; color: #8c8c8c; }
.value { font-size: 28px; font-weight: 700; color: #262626; margin: 4px 0; }
.trend-tag { display: inline-block; padding: 2px 8px; border-radius: 4px; font-size: 12px; font-weight: 600; }

/* 底部布局 & 玻璃拟态修复 */
.bottom-layout {
    display: grid;
    grid-template-columns: 1fr 1.5fr;
    gap: 32px;
}

/* 修复后的玻璃面板样式 */
.glass-panel {
    background: rgba(255, 255, 255, 0.7) !important;
    backdrop-filter: blur(12px);
    -webkit-backdrop-filter: blur(12px);
    border: 1px solid rgba(255, 255, 255, 0.4) !important;
    border-radius: 20px;
}

.panel-header { display: flex; align-items: center; gap: 8px; font-weight: 700; font-size: 16px; }

/* 快捷入口优化 */
.entry-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 20px; padding: 10px; }
.entry-item {
    display: flex; flex-direction: column; align-items: center; gap: 12px;
    cursor: pointer; padding: 20px; border-radius: 16px; transition: all 0.2s;
}
.entry-item:hover { background: rgba(64, 158, 255, 0.08); transform: scale(1.05); }
.icon-circle {
    width: 54px; height: 54px; border-radius: 50%;
    display: flex; align-items: center; justify-content: center; color: white; font-size: 22px;
    box-shadow: 0 8px 15px rgba(0,0,0,0.1);
}
.icon-circle.blue { background: linear-gradient(135deg, #667eea, #764ba2); }
.icon-circle.green { background: linear-gradient(135deg, #42e695, #3bb2b8); }
.icon-circle.orange { background: linear-gradient(135deg, #f6d365, #fda085); }
.entry-label { font-size: 14px; color: #595959; font-weight: 500; }

/* 动画 */
.pulse-dot {
    width: 10px; height: 10px; background: #52c41a; border-radius: 50%;
    margin-right: 8px; animation: pulse 2s infinite;
}
@keyframes pulse {
    0% { box-shadow: 0 0 0 0 rgba(82, 196, 26, 0.7); }
    70% { box-shadow: 0 0 0 10px rgba(82, 196, 26, 0); }
    100% { box-shadow: 0 0 0 0 rgba(82, 196, 26, 0); }
}

@media (max-width: 1100px) {
    .bottom-layout { grid-template-columns: 1fr; }
}
</style>