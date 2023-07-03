<template>
    <el-table :data="tableData" style="width: 100%" :border="true">
        <el-table-column fixed prop="date" label="单号" />
        <el-table-column prop="name" label="制单人" />
        <el-table-column prop="state" label="领用科室" />
        <el-table-column prop="city" label="领用人" />
        <el-table-column prop="address" label="出库科室" />
        <el-table-column prop="zip" label="库房审核时间" />
    </el-table>
</template>

<script setup name="departmentCollar">
import { ref } from "vue";
import bus from "../../eventBus";
const tableData = ref();
bus.on("getData", res => {
    tableData.value = res;
});
</script>
