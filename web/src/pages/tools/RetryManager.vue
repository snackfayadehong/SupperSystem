<template>
    <div class="retry-manager-container">
        <section class="unified-hero-card">
            <div class="hero-content">
                <h1 class="hero-title">接口重试管理中心</h1>
                <p class="hero-desc">监控并补偿发送失败的出库单与退库单据</p>

                <div class="search-area glass-effect">
                    <el-select v-model="queryType" placeholder="单据类型" class="type-select">
                        <el-option label="领用/其他出库单" value="delivery" />
                        <el-option label="科室退库单" value="refund" />
                    </el-select>
                    <el-date-picker
                        v-model="dateRange"
                        type="datetimerange"
                        range-separator="至"
                        start-placeholder="开始时间"
                        end-placeholder="结束时间"
                        value-format="YYYY-MM-DD HH:mm:ss"
                        class="integrated-picker"
                    />
                    <el-button type="primary" :loading="loading" class="audit-btn" @click="handleQuery"> 查询待重试单据 </el-button>
                </div>
            </div>
        </section>

        <transition name="list-fade">
            <div v-if="tableData.length" class="results-container">
                <el-card shadow="never" class="table-card glass-panel">
                    <el-table :data="tableData" style="width: 100%" border stripe>
                        <el-table-column type="expand">
                            <template #default="props">
                                <div class="detail-wrapper">
                                    <el-descriptions title="单据异常详情" :column="2" border>
                                        <el-descriptions-item label="明细序号">{{ props.row.detailSort }}</el-descriptions-item>
                                        <el-descriptions-item label="出库方式">{{ props.row.ckfs }}</el-descriptions-item>
                                        <el-descriptions-item label="错误描述" :span="2">
                                            <span class="error-text">{{ props.row.scsm || "暂无详细错误日志" }}</span>
                                        </el-descriptions-item>
                                    </el-descriptions>
                                </div>
                            </template>
                        </el-table-column>

                        <el-table-column prop="yddh" label="单据编号" width="220" />
                        <el-table-column prop="sczt" label="状态">
                            <template #default="scope">
                                <el-tag :type="scope.row.sczt === '0' ? 'warning' : 'danger'" effect="light">
                                    {{ scope.row.sczt === "0" ? "等待重试" : "发送失败" }}
                                </el-tag>
                            </template>
                        </el-table-column>

                        <el-table-column label="操作" width="150" align="center">
                            <template #default="scope">
                                <el-button type="primary" size="small" icon="Refresh" @click="executeRetry(scope.row)"> 触发重试 </el-button>
                            </template>
                        </el-table-column>
                    </el-table>
                </el-card>
            </div>
        </transition>

        <el-empty v-if="!loading && hasSearched && tableData.length === 0" description="暂无符合条件的待重试单据" />
    </div>
</template>

<script setup>
import { ref } from "vue";
import { Refresh } from "@element-plus/icons-vue";
import myAxios from "@/services/myAxios"; // 使用项目统一 Axios
import { ElMessage, ElMessageBox } from "element-plus";

const queryType = ref("delivery");
const dateRange = ref([]);
const loading = ref(false);
const hasSearched = ref(false);
const tableData = ref([]);

// 查询逻辑
const handleQuery = async () => {
    if (!dateRange.value || dateRange.value.length === 0) {
        return ElMessage.warning("请选择查询时间范围");
    }

    loading.value = true;
    hasSearched.value = true;
    try {
        // 接口路径适配后端控制器
        const res = await myAxios.post(`retry/list`, {
            queryType: queryType.value,
            startTime: dateRange.value[0],
            endTime: dateRange.value[1]
        });
        tableData.value = res.data || [];
    } catch (err) {
        ElMessage.error("查询失败");
    } finally {
        loading.value = false;
    }
};

// 触发重试逻辑 [Requirement 3]
const executeRetry = row => {
    ElMessageBox.confirm(`确认重新推送单据 ${row.ckdh} 到 HIS 系统吗？`, "操作确认", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning"
    }).then(async () => {
        try {
            const res = await myAxios.post(`/api/${queryType.value}/retry`, {
                ckdh: row.ckdh,
                detailSort: row.detailSort
            });
            if (res.code === 0) {
                ElMessage.success("重试指令已下发");
                handleQuery(); // 刷新列表
            }
        } catch (err) {
            ElMessage.error("触发重试异常");
        }
    });
};
</script>

<style scoped>
/* 复用 DictCompareTool.vue 样式规范 */
.retry-manager-container {
    padding: 24px;
    background-color: #f5f7fa;
    min-height: 100vh;
}

.unified-hero-card {
    padding: 40px;
    background: linear-gradient(135deg, #67c23a 0%, #4facfe 100%); /* 换成绿色系区分功能 */
    border-radius: 24px;
    color: white;
    box-shadow: 0 12px 24px rgba(103, 194, 58, 0.2);
    margin-bottom: 24px;
}

.hero-content {
    max-width: 1000px;
    margin: 0 auto;
    text-align: center;
}

.hero-title {
    font-size: 28px;
    font-weight: 800;
}

.hero-desc {
    font-size: 15px;
    opacity: 0.9;
    margin: 10px 0 30px;
}

.glass-effect {
    background: rgba(255, 255, 255, 0.2);
    backdrop-filter: blur(10px);
    padding: 12px;
    border-radius: 16px;
    display: flex;
    gap: 12px;
    border: 1px solid rgba(255, 255, 255, 0.3);
}

.type-select {
    width: 160px;
}
.integrated-picker {
    flex: 1;
}

.table-card {
    border-radius: 16px;
    border: none;
}

.detail-wrapper {
    padding: 20px 40px;
    background-color: #fcfdfe;
}

.error-text {
    color: #f56c6c;
    font-weight: bold;
}

:deep(.el-table__expanded-cell) {
    padding: 0 !important;
}
</style>
