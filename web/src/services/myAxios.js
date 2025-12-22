// src/src/services/myAxios.js 优化版
import axios from "axios";

// 注意：由于配置了自动导入，组件内不需要手动 import ElMessage
// 但在纯 JS 文件中仍需手动引入样式和组件
import {ElLoading, ElMessage} from "element-plus";
import "element-plus/theme-chalk/el-loading.css";
import "element-plus/theme-chalk/el-message.css";

const myAxios = axios.create({
    baseURL: import.meta.env.MODE === "production" ? import.meta.env.VITE_HTTP : import.meta.env.VITE_BASE_URL,
    timeout: 10000,
    withCredentials: false // 修复：去掉了原代码中末尾多余的 'L'
});

myAxios.defaults.headers.post["Content-Type"] = "application/json;charset=UTF-8";

let loadingInstance = null;

myAxios.interceptors.request.use(
    config => {
        loadingInstance = ElLoading.service({ 
            fullscreen: true,
            background: 'rgba(0, 0, 0, 0.1)',
            text: '数据加载中...'
        });
        return config;
    },
    err => {
        if (loadingInstance) loadingInstance.close();
        return Promise.reject(err);
    }
);

myAxios.interceptors.response.use(
    response => {
        if (loadingInstance) loadingInstance.close();
        
        // 统一处理后端返回格式
        const { code, message, data } = response.data;
        
        // 假设 code 200 为成功
        if (code === 0 || response.status === 200) {
            return response.data; // 返回全量数据包含 code 和 data
        }
        
        ElMessage.error(message || "系统业务异常");
        return Promise.reject(new Error(message || "Error"));
    },
    err => {
        if (loadingInstance) loadingInstance.close();
        const errorMsg = err.message.includes("timeout") ? "请求超时，请检查后端服务" : "网络连接异常";
        ElMessage.error(errorMsg);
        return Promise.reject(err);
    }
);

export default myAxios;