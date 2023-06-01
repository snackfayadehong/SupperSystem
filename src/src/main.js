import { createApp } from "vue";
import "./style.css";
import App from "./App.vue";
import ElementPlus from "element-plus";
import zhCn from "element-plus/dist/locale/zh-cn.mjs"; // 组件中文化

import axios from "axios";
import { ElLoading } from "element-plus";

const app = createApp(App);

//配置请求根路径
axios.defaults.baseURL = "http://127.0.0.1";
//全局挂载
app.config.globalProperties.$http = axios;
//请求拦截器
let loadingInstance = null;
axios.interceptors.request.use(config => {
    // 创建Loading组件的实例，并全屏展示loading效果
    loadingInstance = ElLoading.service({ fullscreen: true });
    return config;
});

// 响应拦截器
axios.interceptors.response.use(response => {
    // 调用Loading实例的close（）方法关闭Loading效果
    loadingInstance.close();
    return response;
});
// 中文
app.use(ElementPlus, {
    locale: zhCn
});

app.mount("#app");

// createApp(App).mount("#app");
