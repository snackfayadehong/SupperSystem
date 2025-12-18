<template>
    <div class="stats-grid">
        <div v-for="item in stats" :key="item.title" class="stat-glass-card" :style="{ '--color': item.color }">
            <div class="icon-box" :style="{ background: item.color }">
                <el-icon><component :is="item.icon" /></el-icon>
            </div>
            <div class="content">
                <p class="label">{{ item.title }}</p>
                <p class="value">{{ item.value }}</p>
            </div>
        </div>
    </div>
</template>

<script setup>
import { computed } from "vue";
import { Download, Upload, RefreshLeft, User } from "@element-plus/icons-vue";

const props = defineProps(['data']);

const stats = computed(() => {
    const sum = (type) => props.data.reduce((s, row) => s + (row[type]?.reduce((ss, i) => ss + i.totalAmount, 0) || 0), 0);
    const format = (val) => new Intl.NumberFormat('zh-CN', { style: 'currency', currency: 'CNY' }).format(val);

    return [
        { title: "总入库金额", value: format(sum('inbound')), icon: Download, color: "#67C23A" },
        { title: "总出库金额", value: format(sum('outbound')), icon: Upload, color: "#F56C6C" },
        { title: "退还总额", value: format(sum('return')), icon: RefreshLeft, color: "#E6A23C" },
        { title: "操作员总数", value: props.data.length, icon: User, color: "#409EFF" }
    ];
});
</script>

<style scoped>
.stats-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(260px, 1fr)); gap: 20px; }
.stat-glass-card {
    background: rgba(255, 255, 255, 0.8);
    backdrop-filter: blur(10px);
    border-radius: 14px;
    padding: 20px;
    display: flex;
    align-items: center;
    gap: 16px;
    border: 1px solid rgba(255, 255, 255, 0.5);
    border-left: 5px solid var(--color);
    box-shadow: 0 4px 12px rgba(0,0,0,0.05);
}
.icon-box { width: 48px; height: 48px; border-radius: 12px; display: flex; align-items: center; justify-content: center; color: white; font-size: 22px; }
.label { margin: 0; font-size: 13px; color: #909399; }
.value { margin: 4px 0 0; font-size: 20px; font-weight: bold; color: #303133; }
</style>