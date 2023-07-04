<template>
    <el-table :data="tableData" style="width: 100%" :border="true" class="elTable">
        <el-table-column fixed prop="DepartmentCollarCode" label="单号" align="center" />
        <el-table-column prop="BLMakerName" label="制单人" align="center" />
        <el-table-column prop="LeadingDepartmentName" label="领用科室" align="center" />
        <el-table-column prop="LeaderName" label="领用人" align="center" />
        <el-table-column prop="TreasuryDepartmentName" label="出库科室" align="center" />
        <el-table-column prop="BLDate" label="库房审核时间" align="center" />
    </el-table>
</template>

<script setup name="departmentCollar">
import { ref } from "vue";
import bus from "../../eventBus";
const tableData = ref([]);
bus.on("getDepartmenCollarData", res => {
    tableData.value = res;
    console.log(res);
});
</script>

<style scoped></style>
