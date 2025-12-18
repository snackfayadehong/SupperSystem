import { createApp } from "vue";
import { createPinia } from "pinia";
import App from "./App.vue";
import router from "./router";
import "@/styles/global.css";
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'
// 注意：由于 vite 配置了自动导入，这里不需要手动 use(ElementPlus) 也能使用组件
// 但如果你需要全局配置（如中文），仍可保留

const app = createApp(App);
const pinia = createPinia();
pinia.use(piniaPluginPersistedstate)

app.use(pinia);
app.use(router);
app.use(ElementPlus, {
    locale: zhCn
});
app.mount("#app");