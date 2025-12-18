<template>
    <el-card shadow="never">
        <el-table :data="data" border>
            <el-table-column prop="name" label="名称" />
            <el-table-column prop="count" label="数量" />
            <el-table-column prop="date" label="日期" />
        </el-table>
    </el-card>
</template>

<script setup>
defineProps({
    data: {
        type: Array,
        default: () => []
    }
});
</script>
