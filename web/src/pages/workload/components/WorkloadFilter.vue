<template>
    <el-card shadow="never" class="filter-card">
        <div class="filter-wrapper">
            <el-date-picker
                :model-value="dateRange"
                @update:model-value="$emit('update:dateRange', $event)"
                type="daterange"
                range-separator="至"
                start-placeholder="开始时间"
                end-placeholder="结束时间"
                style="width: 280px"
                :shortcuts="shortcuts"
            />

            <el-button type="primary" :icon="Search" @click="$emit('query')">查询数据</el-button>

            <el-divider direction="vertical" />

            <el-input :model-value="search" @input="$emit('update:search', $event)" placeholder="搜索操作员" prefix-icon="Search" clearable style="width: 200px" />

            <el-select :model-value="type" @change="$emit('update:type', $event)" placeholder="业务筛选" clearable style="width: 160px">
                <el-option label="全部业务" value="all" />
                <el-option label="📥 入库验收" value="inbound" />
                <el-option label="📤 出库发放" value="outbound" />
                <el-option label="📝 入库登记" value="inReg" />
                <el-option label="🔙 退货业务" value="return" />
                <el-option label="♻️ 二级库退库" value="secondaryRefund" />
                <el-option label="🛒 采购下单" value="purchase" />
                <el-option label="🔔 催货跟进" value="push" />
            </el-select>
        </div>
    </el-card>
</template>

<script setup>
import { Search } from "@element-plus/icons-vue";
defineProps(["search", "type", "dateRange", "total"]);
defineEmits(["update:search", "update:type", "update:dateRange", "query"]);

const shortcuts = [
    {
        text: "最近一周",
        value: () => {
            const end = new Date();
            const start = new Date();
            start.setTime(start.getTime() - 3600 * 1000 * 24 * 7);
            return [start, end];
        }
    },
    {
        text: "最近15天",
        value: () => {
            const end = new Date();
            const start = new Date();
            start.setTime(start.getTime() - 3600 * 1000 * 24 * 15);
            return [start, end];
        }
    },
    {
        text: "最近一个月",
        value: () => {
            const end = new Date();
            const start = new Date();
            start.setMonth(start.getMonth() - 1);
            return [start, end];
        }
    }
];
</script>

<style scoped>
.filter-card {
    border-radius: 8px;
    margin-bottom: 16px;
    border: none;
}
.filter-wrapper {
    display: flex;
    align-items: center;
    gap: 12px;
    flex-wrap: wrap;
}
</style>
