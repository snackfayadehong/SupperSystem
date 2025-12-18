// src/pages/workload/composables/useWorkload.js
import { ref, computed } from "vue";
import { ElMessage } from "element-plus";
import myAxios from "@/services/myAxios"; //
import dayjs from "dayjs"; 

export function useWorkload() {
    const rawData = ref([]);
    const loading = ref(false);
    const searchQuery = ref("");
    const filterType = ref("");
    const dateRange = ref([]); // 存储选中的 [Date, Date]
    const currentPage = ref(1);
    const pageSize = ref(10);

    // 手动查询函数：对接 /api/getWorkload
    const fetchList = async () => {
        if (!dateRange.value || dateRange.value.length !== 2) {
            return ElMessage.warning({
                message: "请先选择查询时间范围",
                grouping: true, // 合并相同消息，防止撑开页面
            });
        }

        loading.value = true;
        try {
            // 将日期格式化为后端识别的 YYYY-MM-DD
            const payload = {
                startDate: dayjs(dateRange.value[0]).format("YYYY-MM-DD"),
                endDate: dayjs(dateRange.value[1]).format("YYYY-MM-DD")
            };
            
            const res = await myAxios.post("/getWorkload", payload);
            rawData.value = res.data || [];
            ElMessage.success("查询成功");
        } catch (error) {
            console.error(error);
            ElMessage.error("获取数据失败，请检查后端服务");
        } finally {
            loading.value = false;
        }
    };

    // 前端过滤：基于已获取的数据进行检索
    const filteredData = computed(() => {
        return rawData.value.filter(item => {
            const matchSearch = item.operator.toLowerCase().includes(searchQuery.value.toLowerCase());
            const matchType = !filterType.value || (item[filterType.value] && item[filterType.value].length > 0);
            return matchSearch && matchType;
        });
    });

    // 分页切片
    const paginatedData = computed(() => {
        const start = (currentPage.value - 1) * pageSize.value;
        return filteredData.value.slice(start, start + pageSize.value);
    });

    return {
        rawData, loading, searchQuery, filterType, dateRange,
        currentPage, pageSize, filteredData, paginatedData, fetchList
    };
}