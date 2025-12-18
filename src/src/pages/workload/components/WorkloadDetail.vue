<!-- # 交互层：详细数据的弹窗展示 -->
<template>
    <el-dialog 
        :model-value="visible" 
        @update:model-value="$emit('update:visible', $event)"
        :title="`🔍 详细数据 - ${userData?.operator}`" 
        width="80%" 
        destroy-on-close
    >
        <div v-if="userData" class="detail-container">
            <el-descriptions :column="3" border>
                <el-descriptions-item label="操作员">{{ userData.operator }}</el-descriptions-item>
                <el-descriptions-item label="入库总计">{{ formatCurrency(calculateTotal(userData.inbound)) }}</el-descriptions-item>
                <el-descriptions-item label="出库总计">{{ formatCurrency(calculateTotal(userData.outbound)) }}</el-descriptions-item>
            </el-descriptions>

            <el-tabs type="border-card" class="detail-tabs">
                <el-tab-pane label="入库明细">
                    <el-table :data="userData.inbound" border height="300">
                        <el-table-column prop="category" label="分类" />
                        <el-table-column prop="specCount" label="品规" />
                        <el-table-column prop="totalAmount" label="金额">
                            <template #default="{ row }">{{ formatCurrency(row.totalAmount) }}</template>
                        </el-table-column>
                    </el-table>
                </el-tab-pane>
                </el-tabs>
        </div>
        
        <template #footer>
            <el-button @click="$emit('update:visible', false)">关闭</el-button>
            <el-button type="primary" @click="$emit('export-detail', userData)">导出此条详情</el-button>
        </template>
    </el-dialog>
</template>

<script setup>
defineProps(['visible', 'userData']);
defineEmits(['update:visible', 'export-detail']);

const formatCurrency = val => new Intl.NumberFormat("zh-CN", { style: "currency", currency: "CNY" }).format(val);
const calculateTotal = items => (items || []).reduce((sum, i) => sum + i.totalAmount, 0);
</script>

<style scoped>
.detail-tabs { margin-top: 20px; }
</style>