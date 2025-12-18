// src/router/index.js
import { createRouter, createWebHistory } from "vue-router";
import MainLayout from "@/layout/MainLayout.vue";
import menuList from "./menu";
import NProgress from "nprogress"; // 引入进度条
import "nprogress/nprogress.css"; // 引入进度条样式

// 配置 NProgress (去掉右上角的螺旋加载圈)
NProgress.configure({ showSpinner: false });

const pages = import.meta.glob("../pages/**/*.vue");

function loadView(view) {
    const key = `../pages/${view}.vue`;
    const loader = pages[key];
    if (!loader) {
        console.error(`❌ 页面不存在: ${key}`);
    }
    return loader;
}

const routes = [
    {
        path: "/",
        component: MainLayout,
        redirect: "/home",
        children: [
            ...menuList.map(item => ({
                path: item.path,
                name: item.name,
                component: loadView(item.component),
                meta: { title: item.label }
            })),
            {
                path: "/:pathMatch(.*)*",
                name: "NotFound",
                component: () => import("@/pages/error/NotFound.vue"),
                meta: { title: "404 - 页面不存在" }
            }
        ]
    }
];

const router = createRouter({
    history: createWebHistory(),
    routes,
    // 优化 1: 切换路由时自动滚动到页面顶部
    scrollBehavior() {
        return { top: 0 };
    }
});

// 优化 2: 增加前置守卫开启进度条
router.beforeEach((to, from, next) => {
    NProgress.start();
    next();
});

router.afterEach((to) => {
    // 优化 3: 进度条结束
    NProgress.done();
    document.title = to.meta.title ? `${to.meta.title} | WorkloadQuery` : 'WorkloadQuery';
});

export default router;