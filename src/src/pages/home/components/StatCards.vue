<template>
    <el-row :gutter="16">
        <el-col :span="6" v-for="item in stats" :key="item.title">
            <el-card shadow="hover">
                <div class="stat-card">
                    <div class="title">{{ item.title }}</div>
                    <div class="value">{{ item.value }}</div>
                </div>
            </el-card>
        </el-col>
    </el-row>
</template>

<script setup>
const stats = [
    { title: "总记录数", value: 1234 },
    { title: "今日查询", value: 56 },
    { title: "异常数据", value: 3 },
    { title: "系统状态", value: "正常" }
];
</script>

<style scoped>
.stat-card {
    text-align: center;
}
.title {
    color: #666;
    font-size: 14px;
}
.value {
    margin-top: 8px;
    font-size: 24px;
    font-weight: bold;
}
</style>
