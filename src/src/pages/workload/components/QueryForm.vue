<template>
    <el-card shadow="never">
        <el-form :inline="true" :model="form">
            <el-form-item label="产品">
                <el-input v-model="form.product" placeholder="产品名" />
            </el-form-item>

            <el-form-item label="日期">
                <el-date-picker v-model="form.date" type="daterange" />
            </el-form-item>

            <el-form-item>
                <el-button type="primary" @click="$emit('search', form)"> 查询 </el-button>
                <el-button @click="$emit('reset')"> 重置 </el-button>
            </el-form-item>
        </el-form>
    </el-card>
</template>

<script setup>
import { reactive } from "vue";

const form = reactive({
    product: "",
    date: []
});
</script>
