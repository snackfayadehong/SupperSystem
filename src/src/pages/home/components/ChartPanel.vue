<template>
    <el-card shadow="never">
        <template #header>
            <span>数据趋势</span>
        </template>

        <div class="chart-placeholder">图表区域（后续接 ECharts）</div>
    </el-card>
</template>

<style scoped>
.chart-placeholder {
    height: 300px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #999;
    background: #f5f7fa;
}
</style>
