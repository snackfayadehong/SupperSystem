<template>
    <div class="workload-table-wrapper">
        <el-skeleton :loading="loading" animated :rows="10">
            <template #default>
                <el-table :data="data" border stripe style="width: 100%" class="modern-dashboard-table"
                    header-cell-class-name="modern-table-header" row-key="operator">
                    <el-table-column type="expand">
                        <template #default="{ row }">
                            <div class="expand-content">
                                <el-row :gutter="40">
                                    <el-col :span="12">
                                        <h4 class="preview-title">非金额业务统计 (单据量)</h4>
                                        <div class="mini-stat-row">
                                            <div class="mini-stat">
                                                <span class="label">采购下单</span>
                                                <span class="val">{{ calculateCount(row.purchase) }} 单</span>
                                            </div>
                                            <div class="mini-stat">
                                                <span class="label">催货跟进</span>
                                                <span class="val">{{ calculateCount(row.push) }} 单</span>
                                            </div>
                                        </div>
                                    </el-col>
                                    <el-col :span="12">
                                        <h4 class="preview-title">业务 TOP1 概览</h4>
                                        <div class="top-items-grid">
                                            <div class="top-item success-bg" v-if="row.inbound?.length">
                                                <span class="type-label">入库最忙：</span>
                                                <span class="item-val">{{ getTopCategory(row.inbound) }}</span>
                                            </div>
                                            <div class="top-item info-bg" v-if="row.inReg?.length">
                                                <span class="type-label">登记最多：</span>
                                                <span class="item-val">{{ getTopCategory(row.inReg) }}</span>
                                            </div>
                                        </div>
                                    </el-col>
                                </el-row>
                            </div>
                        </template>
                    </el-table-column>

                    <el-table-column prop="operator" label="操作人员" align="center" min-width="110" fixed>
                        <template #default="{ row }">
                            <span class="operator-name-modern">{{ row.operator }}</span>
                        </template>
                    </el-table-column>

                    <el-table-column label="入库验收(金额)" align="center" min-width="150">
                        <template #default="{ row }">
                            <span class="amount-val success">{{ formatCurrency(calculateTotal(row.inbound)) }}</span>
                        </template>
                    </el-table-column>

                    <el-table-column label="入库登记(金额)" align="center" min-width="150">
                        <template #default="{ row }">
                            <span class="amount-val primary">{{ formatCurrency(calculateTotal(row.inReg)) }}</span>
                        </template>
                    </el-table-column>

                    <el-table-column label="出库发放(金额)" align="center" min-width="150">
                        <template #default="{ row }">
                            <span class="amount-val danger">{{ formatCurrency(calculateTotal(row.outbound)) }}</span>
                        </template>
                    </el-table-column>

                    <el-table-column label="退货退库" align="center" min-width="140">
                        <template #default="{ row }">
                            <span class="amount-val warning">{{ formatCurrency(calculateTotal(row.return) +
                                calculateTotal(row.secondaryRefund)) }}</span>
                        </template>
                    </el-table-column>

                    <el-table-column label="管理" align="center" width="150" fixed="right">
                        <template #default="{ row }">
                            <el-button type="primary" link icon="View" @click="$emit('view-detail', row)">明细</el-button>
                            <el-button type="success" link icon="Download"
                                @click="$emit('export-row', row)">导出</el-button>
                        </template>
                    </el-table-column>
                </el-table>
            </template>
        </el-skeleton>
    </div>
</template>

<script setup>
const props = defineProps({
    data: { type: Array, default: () => [] },
    loading: { type: Boolean, default: false }
});

defineEmits(["view-detail", "export-row"]);

const getTopCategory = items => {
    if (!items || !Array.isArray(items) || items.length === 0) return "无记录";
    const sorted = [...items].sort((a, b) => (b.totalAmount || 0) - (a.totalAmount || 0));
    const top = sorted[0];
    return `${top.category} (￥${Math.round(top.totalAmount / 10000)}万)`; // 简化显示
};

const formatCurrency = val => {
    if (!val) return "-";
    return new Intl.NumberFormat("zh-CN", { style: "decimal", minimumFractionDigits: 2 }).format(val);
};

const calculateTotal = items => {
    if (!items || !Array.isArray(items)) return 0;
    return items.reduce((sum, item) => sum + (item.totalAmount || 0), 0);
};

const calculateCount = items => {
    if (!items || !Array.isArray(items)) return 0;
    return items.reduce((sum, item) => sum + (item.billCount || 0), 0);
};
</script>

<style scoped>
.workload-table-wrapper {
    width: 100%;
}

.modern-dashboard-table {
    --el-table-border-color: #f0f0f0;
    border-radius: 12px;
    overflow: hidden;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.02);
}

.operator-name-modern {
    font-weight: 700;
    color: #262626;
    font-size: 15px;
}

.amount-val {
    font-family: "Monaco", monospace;
    font-weight: 600;
    font-size: 15px;
}

.success {
    color: #67c23a;
}

.danger {
    color: #f56c6c;
}

.warning {
    color: #e6a23c;
}

.primary {
    color: #409eff;
}

/* 展开行样式 */
.expand-content {
    padding: 20px 40px;
    background: #fcfcfc;
}

.preview-title {
    margin: 0 0 12px;
    font-size: 13px;
    color: #909399;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.5px;
}

.mini-stat-row {
    display: flex;
    gap: 24px;
}

.mini-stat {
    background: #fff;
    padding: 10px 16px;
    border-radius: 8px;
    border: 1px solid #ebeef5;
    display: flex;
    flex-direction: column;
}

.mini-stat .label {
    font-size: 12px;
    color: #909399;
    margin-bottom: 4px;
}

.mini-stat .val {
    font-size: 16px;
    font-weight: bold;
    color: #303133;
}

.top-items-grid {
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.top-item {
    padding: 8px 12px;
    border-radius: 6px;
    font-size: 13px;
    border: 1px solid transparent;
    display: flex;
    justify-content: space-between;
}

.success-bg {
    background: #f0f9eb;
    color: #67c23a;
    border-color: #e1f3d8;
}

.info-bg {
    background: #f4f4f5;
    color: #909399;
    border-color: #e9e9eb;
}

.item-val {
    font-weight: bold;
}

:deep(.modern-table-header) {
    background-color: #fafafa !important;
    color: #606266;
    height: 48px;
}

:global(html.dark) .expand-content {
    background: #1a1a1a;
}
</style>
