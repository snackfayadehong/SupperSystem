// src/pages/home/composables/useHomeData.js
import {computed, onMounted, ref} from "vue";
import dayjs from "dayjs";

export function useHomeData() {
    const currentTime = ref(dayjs());
    const stats = ref([
        { title: "今日查询量", value: "128", icon: "Search", color: "#409EFF", trend: "+12%" },
        { title: "本月处理量", value: "3,542", icon: "DataAnalysis", color: "#67C23A", trend: "+5.4%" },
        { title: "系统拦截", value: "4", icon: "WarnTriangleFilled", color: "#F56C6C", trend: "正常" },
        { title: "活跃用户", value: "12", icon: "User", color: "#E6A23C", trend: "实时" }
    ]);

    // 动态问候语
    const welcomeMessage = computed(() => {
        const hour = currentTime.value.hour();
        if (hour < 9) return "早安，管理员";
        if (hour < 12) return "上午好，管理员";
        if (hour < 18) return "下午好，管理员";
        return "晚上好，管理员";
    });

    const updateTime = () => {
        currentTime.value = dayjs();
    };

    onMounted(() => {
        const timer = setInterval(updateTime, 1000);
        return () => clearInterval(timer);
    });

    return { stats, welcomeMessage, currentTime };
}