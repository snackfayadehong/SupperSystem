<template>
    <div class="workload-page">
        <el-card shadow="never" class="custom-page-header">
            <div class="header-flex">
                <div class="title-group">
                    <h1 class="main-title">工作量看板</h1>
                    <span class="sub-title">请选择时间段后点击查询以获取最新数据</span>
                </div>
                <div class="header-actions">
                    <el-button type="primary" :icon="Download" plain @click="handleExport">导出报表</el-button>
                </div>
            </div>
        </el-card>

        <main class="dashboard-content">
            <WorkloadStats :data="rawData" />

            <WorkloadFilter 
                v-model:search="searchQuery" 
                v-model:type="filterType" 
                v-model:dateRange="dateRange"
                :total="filteredData.length"
                @query="fetchList"
            />

            <el-card shadow="hover" class="table-card">
                <WorkloadTable 
                    :data="paginatedData" 
                    :loading="loading" 
                    @view-detail="handleViewDetail"
                />
                
                <div class="pagination-container">
                    <el-pagination
                        v-model:current-page="currentPage"
                        v-model:page-size="pageSize"
                        :total="filteredData.length"
                        :page-sizes="[10, 20, 50]"
                        layout="total, sizes, prev, pager, next"
                        background
                    />
                </div>
            </el-card>
        </main>

        <WorkloadDetail v-model:visible="detailVisible" :user-data="currentDetail" />
    </div>
</template>

<script setup>
import { ref } from "vue";
import { Download } from "@element-plus/icons-vue";
import { useWorkload } from "./composables/useWorkload";
import WorkloadStats from "./components/WorkloadStats.vue";
import WorkloadFilter from "./components/WorkloadFilter.vue";
import WorkloadTable from "./components/WorkloadTable.vue";
import WorkloadDetail from "./components/WorkloadDetail.vue";

const {
    rawData, loading, searchQuery, filterType, dateRange,
    currentPage, pageSize, filteredData, paginatedData, fetchList
} = useWorkload();

const detailVisible = ref(false);
const currentDetail = ref(null);

const handleViewDetail = (row) => {
    currentDetail.value = row;
    detailVisible.value = true;
};
</script>

<style scoped>
.workload-page { padding: 24px; background: #f5f7fa; min-height: 100%; }
.custom-page-header { border-radius: 12px; margin-bottom: 20px; border: none; }
.header-flex { display: flex; justify-content: space-between; align-items: center; }
.main-title { margin: 0; font-size: 20px; color: #303133; font-weight: 700; }
.sub-title { font-size: 13px; color: #909399; margin-top: 4px; display: block; }
.dashboard-content { display: flex; flex-direction: column; gap: 20px; }
.pagination-container { margin-top: 20px; display: flex; justify-content: flex-end; }
</style>