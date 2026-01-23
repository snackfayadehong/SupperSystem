<template>
    <div class="workload-detail-wrapper">
        <div class="detail-summary-scroll">
            <div class="summary-item success">
                <div class="s-info">
                    <span class="s-label">入库验收</span><span class="s-value">￥{{ format(calculateTotal(data?.inbound)) }}</span>
                </div>
            </div>
            <div class="summary-item primary">
                <div class="s-info">
                    <span class="s-label">入库登记</span><span class="s-value">￥{{ format(calculateTotal(data?.inReg)) }}</span>
                </div>
            </div>
            <div class="summary-item danger">
                <div class="s-info">
                    <span class="s-label">出库发放</span><span class="s-value">￥{{ format(calculateTotal(data?.outbound)) }}</span>
                </div>
            </div>
            <div class="summary-item warning">
                <div class="s-info">
                    <span class="s-label">所有退还</span><span class="s-value">￥{{ format(calculateTotal(data?.return) + calculateTotal(data?.secondaryRefund)) }}</span>
                </div>
            </div>
        </div>

        <el-card shadow="never" class="modern-detail-card">
            <template #header>
                <div class="detail-header">
                    <div class="user-info">
                        <el-avatar :size="32" :icon="UserFilled" class="user-avatar" />
                        <span class="operator-name">{{ data?.operator || "未知人员" }}</span>
                        <el-tag size="small" type="info" effect="plain" round class="status-tag">全业务视图</el-tag>
                    </div>
                    <div class="header-actions">
                        <el-button type="success" size="small" icon="Download" plain @click="$emit('export-current', data)">导出 Excel</el-button>
                        <el-button type="info" link @click="$emit('close')">关闭</el-button>
                    </div>
                </div>
            </template>

            <div class="sections-stack">
                <template v-for="(section, key) in sections" :key="key">
                    <div v-if="data?.[key]?.length" class="detail-section">
                        <h4 class="modern-section-title" :class="section.class">
                            <span class="indicator"></span> {{ section.title }}
                            <span class="count-badge">{{ data[key].length }} 类/项</span>
                        </h4>
                        <el-table :data="data[key]" border stripe size="small" class="glass-table">
                            <el-table-column prop="category" label="业务/物料分类" min-width="160" show-overflow-tooltip />
                            <el-table-column prop="billCount" label="单据数" align="center" width="90" />
                            <el-table-column prop="specCount" label="品规数" align="center" width="90" />
                            <el-table-column v-if="!section.hideMoney" prop="totalAmount" label="合计金额" align="right" width="140">
                                <template #default="{ row }">
                                    <span class="amount-font">￥{{ row.totalAmount?.toLocaleString(undefined, { minimumFractionDigits: 2 }) }}</span>
                                </template>
                            </el-table-column>
                        </el-table>
                    </div>
                </template>
            </div>

            <el-empty v-if="!hasAnyData" :image-size="100" description="该人员在此时间段内无任何业务记录" />
        </el-card>
    </div>
</template>

<script setup>
import { computed } from "vue";
import { UserFilled, Download } from "@element-plus/icons-vue";

const props = defineProps({
    data: { type: Object, default: () => ({}) }
});

defineEmits(["close", "export-current"]);

// 定义各板块配置，方便循环渲染
const sections = {
    inbound: { title: "入库验收 (Inbound)", class: "in", hideMoney: false },
    inReg: { title: "入库登记 (Registration)", class: "reg", hideMoney: false },
    outbound: { title: "出库发放 (Outbound)", class: "out", hideMoney: false },
    return: { title: "退货业务 (Return)", class: "ret", hideMoney: false },
    secondaryRefund: { title: "二级库退库 (Secondary Refund)", class: "sec", hideMoney: false },
    purchase: { title: "采购订单 (Purchase)", class: "pur", hideMoney: true }, // 隐藏金额列
    push: { title: "催货记录 (Push)", class: "push", hideMoney: true }
};

const hasAnyData = computed(() => {
    return Object.keys(sections).some(key => props.data?.[key]?.length > 0);
});

const calculateTotal = list => list?.reduce((s, i) => s + (i.totalAmount || 0), 0) || 0;
const format = v => v.toLocaleString(undefined, { minimumFractionDigits: 0, maximumFractionDigits: 0 });
</script>

<style scoped>
.workload-detail-wrapper {
    padding-bottom: 20px;
}

/* 滚动汇总栏 */
.detail-summary-scroll {
    display: flex;
    gap: 12px;
    overflow-x: auto;
    padding-bottom: 12px;
    margin-bottom: 16px;
}
.summary-item {
    flex: 0 0 160px; /* 固定宽度 */
    padding: 12px 16px;
    border-radius: 10px;
    border: 1px solid transparent;
}
.summary-item.success {
    background: #f0f9eb;
    border-color: #b3e19d;
    color: #67c23a;
}
.summary-item.primary {
    background: #ecf5ff;
    border-color: #b3d8ff;
    color: #409eff;
}
.summary-item.danger {
    background: #fef0f0;
    border-color: #fbc4c4;
    color: #f56c6c;
}
.summary-item.warning {
    background: #fdf6ec;
    border-color: #f5dab1;
    color: #e6a23c;
}

.s-label {
    font-size: 12px;
    display: block;
    opacity: 0.8;
    margin-bottom: 4px;
}
.s-value {
    font-size: 16px;
    font-weight: 800;
    font-family: "Monaco", monospace;
}

.modern-detail-card {
    border-radius: 16px !important;
    border: 1px solid #ebeef5 !important;
}
.detail-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
}
.user-info {
    display: flex;
    align-items: center;
    gap: 10px;
}
.user-avatar {
    background: #409eff;
}
.operator-name {
    font-weight: 700;
    font-size: 16px;
    color: #303133;
}

.sections-stack {
    display: flex;
    flex-direction: column;
    gap: 24px;
    margin-top: 10px;
}

.modern-section-title {
    display: flex;
    align-items: center;
    font-size: 15px;
    font-weight: 700;
    margin-bottom: 12px;
    color: #303133;
}
.indicator {
    width: 4px;
    height: 14px;
    border-radius: 2px;
    margin-right: 8px;
}

/* 各板块颜色 */
.in .indicator {
    background: #67c23a;
}
.reg .indicator {
    background: #409eff;
}
.out .indicator {
    background: #f56c6c;
}
.ret .indicator {
    background: #e6a23c;
}
.sec .indicator {
    background: #909399;
}
.pur .indicator,
.push .indicator {
    background: #9c27b0;
}

.count-badge {
    margin-left: auto;
    font-size: 12px;
    font-weight: normal;
    color: #909399;
    background: #f4f4f5;
    padding: 2px 8px;
    border-radius: 10px;
}
.amount-font {
    font-family: "Monaco", monospace;
    font-weight: 700;
    color: #303133;
}
</style>
