<template>
    <el-table :data="data" v-loading="loading" row-key="operator" border stripe>
        <el-table-column type="expand">
            <template #default="{ row }">
                <div class="table-expand-wrapper">
                    <el-descriptions title="业务分类明细" :column="3" border size="small">
                        <el-descriptions-item label="入库项">{{ row.inbound?.length || 0 }}</el-descriptions-item>
                        <el-descriptions-item label="出库项">{{ row.outbound?.length || 0 }}</el-descriptions-item>
                        <el-descriptions-item label="退还项">{{ row.return?.length || 0 }}</el-descriptions-item>
                    </el-descriptions>
                </div>
            </template>
        </el-table-column>

        <el-table-column prop="operator" label="操作人员" align="center" width="180" />
        <el-table-column label="入库统计" align="center">
            <template #default="{ row }">
                <span class="amount success">{{ formatCurrency(calculateTotal(row.inbound)) }}</span>
            </template>
        </el-table-column>
        <el-table-column label="出库统计" align="center">
            <template #default="{ row }">
                <span class="amount danger">{{ formatCurrency(calculateTotal(row.outbound)) }}</span>
            </template>
        </el-table-column>
        <el-table-column label="操作" align="center" width="120">
            <template #default="{ row }">
                <el-button type="primary" link @click="$emit('view-detail', row)">详情</el-button>
            </template>
        </el-table-column>
    </el-table>
</template>

<script setup>
const props = defineProps(['data', 'loading']);
const formatCurrency = (val) => new Intl.NumberFormat('zh-CN', { style: 'currency', currency: 'CNY' }).format(val);
const calculateTotal = (items) => (items || []).reduce((sum, i) => sum + i.totalAmount, 0);
</script>

<style scoped>
.table-expand-wrapper { padding: 15px 50px; background: #fafafa; }
.amount { font-weight: bold; }
.amount.success { color: #67c23a; }
.amount.danger { color: #f56c6c; }
</style>