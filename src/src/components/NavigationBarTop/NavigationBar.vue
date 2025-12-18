<template>
    <el-tabs v-model="activeName" type="card" class="demo-tabs" @tab-click="handleClick">
        <el-tab-pane label="库房工作量统计" name="workload">
            <query-bar parent-name="workload" />
            <workload-data-show />
        </el-tab-pane>
        <el-tab-pane label="123" name="WorkLoadDataShow2">
            <query-bar parent-name="WorkLoadDataShow2" />
            <work-load-data-show2 />
        </el-tab-pane>
    </el-tabs>
</template>
<script setup>
import { ref } from "vue";
import QueryBar from "../QueryBarTop/QueryBar.vue";
import WorkLoadDataShow2 from "@/pages/workload/components/WorkLoadDataShow2.vue";
import { ElTabPane } from "element-plus";

const activeName = ref("workload");
const handleClick = () => {};
</script>
<style>
.demo-tabs > .el-tabs__content {
    padding-left: 32px;
    font-size: 32px;
    font-weight: 600;
}
.elTable {
    font-size: 16px;
    --el-table-border-color: #a4b6e1;
    --el-table-border: 1px solid #a4b6e1;
    --el-table-header-text-color: black;
    --el-text-color-secondary: black;
}
</style>
