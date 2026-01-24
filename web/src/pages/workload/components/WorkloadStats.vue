<template>
    <div class="stats-container">
        <div class="stats-cards-grid">
            <div v-for="item in stats" :key="item.title" class="modern-stat-card">
                <div class="card-icon" :style="{ background: item.color + '15', color: item.color }">
                    <el-icon :size="24">
                        <component :is="item.icon" />
                    </el-icon>
                </div>
                <div class="card-body">
                    <span class="label">{{ item.title }}</span>
                    <h2 class="value">{{ item.value }}</h2>
                    <span v-if="item.sub" class="sub-text">{{ item.sub }}</span>
                </div>
            </div>
        </div>

        <div class="analysis-glass-card">
            <div class="chart-box">
                <svg viewBox="0 0 36 36" class="circular-chart">
                    <path class="circle-bg"
                        d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831" />

                    <path v-for="segment in chartSegments" :key="segment.key" class="circle" :stroke="segment.color"
                        :style="{
                            strokeDasharray: `${segment.rawPercent}, 100`,
                            strokeDashoffset: `-${segment.offset}`
                        }" d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831" />
                </svg>
                <div class="chart-center">
                    <span class="center-label">金额分布</span>
                </div>
            </div>

            <div class="chart-legend">
                <div v-for="seg in chartSegments" :key="seg.key" class="legend-item">
                    <span class="dot" :style="{ background: seg.color }"></span>
                    {{ seg.label }} {{ seg.displayLabel }}
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { computed } from "vue";
import { Download, Upload, RefreshLeft, User, List, Money } from "@element-plus/icons-vue";

const props = defineProps(["data"]);

// 辅助计算总金额
const sumAmount = list => list?.reduce((acc, item) => acc + (item.totalAmount || 0), 0) || 0;
// 辅助计算单据数
const sumBills = (rows, type) => rows.reduce((acc, row) => acc + (row[type]?.reduce((s, i) => s + (i.billCount || 0), 0) || 0), 0);

// 计算环形图数据 (只计算有金额的业务)
const chartSegments = computed(() => {
    const totalMap = {
        inbound: props.data.reduce((s, r) => s + sumAmount(r.inbound), 0),
        outbound: props.data.reduce((s, r) => s + sumAmount(r.outbound), 0),
        return: props.data.reduce((s, r) => s + sumAmount(r.return), 0),
        inReg: props.data.reduce((s, r) => s + sumAmount(r.inReg), 0),
        secondary: props.data.reduce((s, r) => s + sumAmount(r.secondaryRefund), 0)
    };

    const total = Object.values(totalMap).reduce((a, b) => a + b, 0) || 1;

    // 定义配置
    const configs = [
        { key: "inbound", label: "入库", color: "#67C23A" },
        { key: "outbound", label: "出库", color: "#F56C6C" },
        { key: "inReg", label: "登记", color: "#409EFF" },
        { key: "return", label: "退还供方", color: "#E6A23C" },
        { key: "secondary", label: "二级库退库", color: "#909399" }
    ];

    let currentOffset = 0;
    return configs
        .map(cfg => {
            const val = totalMap[cfg.key];
            // 核心修改 1: 使用浮点数 rawPercent 保证绘图精度，不进行 Math.round
            const rawPercent = (val / total) * 100;

            // 核心修改 2: 智能格式化显示文本
            let displayLabel = '0%';
            if (rawPercent > 0) {
                if (rawPercent < 0.01) {
                    displayLabel = '<0.01%'; // 极小值
                } else if (rawPercent < 1) {
                    displayLabel = rawPercent.toFixed(2) + '%'; // 小于1%显示两位小数，如 0.02%
                } else {
                    displayLabel = Math.round(rawPercent) + '%'; // 正常值保持整数
                }
            }

            const segment = {
                ...cfg,
                value: val,
                rawPercent: rawPercent, // 用于 style 绑定
                displayLabel: displayLabel, // 用于页面文字显示
                offset: currentOffset
            };

            currentOffset += rawPercent; // 累加也使用精确值
            return segment;
        })
        .filter(s => s.value > 0); // 只显示有数据的段
});

const stats = computed(() => {
    const format = val => new Intl.NumberFormat("zh-CN", { style: "currency", currency: "CNY", maximumFractionDigits: 0 }).format(val);

    // 计算各维度总额
    const totalInbound = props.data.reduce((s, r) => s + sumAmount(r.inbound), 0);
    const totalOutbound = props.data.reduce((s, r) => s + sumAmount(r.outbound), 0);
    const totalInReg = props.data.reduce((s, r) => s + sumAmount(r.inReg), 0);
    const totalReturn = props.data.reduce((s, r) => s + sumAmount(r.return) + sumAmount(r.secondaryRefund), 0);

    // 采购/催货 只算单数
    const countPurchase = sumBills(props.data, "purchase");
    const countPush = sumBills(props.data, "push");

    return [
        { title: "入库验收总额", value: format(totalInbound), icon: Download, color: "#67C23A" },
        { title: "出库发放总额", value: format(totalOutbound), icon: Upload, color: "#F56C6C" },
        { title: "入库登记总额", value: format(totalInReg), icon: List, color: "#409EFF" },
        { title: "各类退还总额", value: format(totalReturn), icon: RefreshLeft, color: "#E6A23C", sub: "(含二级库)" },
        { title: "采购与催货", value: `${countPurchase + countPush} 单`, icon: Money, color: "#9c27b0", sub: `采购:${countPurchase} 催货:${countPush}` },
        { title: "活跃操作员", value: props.data.length + " 人", icon: User, color: "#303133" }
    ];
});
</script>

<style scoped>
/* 保持原有布局，增加网格自适应 */
.stats-container {
    display: grid;
    grid-template-columns: 1fr 300px;
    gap: 24px;
    margin-bottom: 24px;
}

.stats-cards-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
    /* 自适应列宽 */
    gap: 20px;
}

.modern-stat-card {
    background: var(--bg-card, #fff);
    padding: 20px;
    border-radius: 16px;
    display: flex;
    align-items: center;
    gap: 16px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.03);
    transition: transform 0.3s;
    border: 1px solid transparent;
}

.modern-stat-card:hover {
    transform: translateY(-5px);
    border-color: rgba(0, 0, 0, 0.05);
}

/* 字体和颜色微调 */
.card-body {
    display: flex;
    flex-direction: column;
}

.label {
    font-size: 13px;
    color: #909399;
}

.value {
    margin: 4px 0 0;
    font-size: 20px;
    font-weight: 700;
    color: #303133;
}

.sub-text {
    font-size: 12px;
    color: #c0c4cc;
    margin-top: 2px;
}

/* 图表部分样式复用原版，增加颜色类 */
.analysis-glass-card {
    background: var(--bg-card, #fff);
    padding: 24px;
    border-radius: 24px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.03);
}

.circular-chart {
    width: 140px;
    height: 140px;
    transform: rotate(-90deg);
}

.circle-bg {
    fill: none;
    stroke: #f0f2f5;
    stroke-width: 3.5;
}

.circle {
    fill: none;
    stroke-width: 3.8;
    stroke-linecap: round;
    transition: all 0.5s ease;
}

.chart-box {
    position: relative;
}

.chart-center {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    text-align: center;
}

.center-label {
    font-size: 13px;
    color: #909399;
    font-weight: 500;
}

.chart-legend {
    margin-top: 20px;
    width: 100%;
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.legend-item {
    display: flex;
    align-items: center;
    gap: 10px;
    font-size: 13px;
    color: #606266;
}

.dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
}

:global(html.dark) .modern-stat-card,
:global(html.dark) .analysis-glass-card {
    background: #1e1e1e;
}

:global(html.dark) .value {
    color: #e5eaf3;
}
</style>