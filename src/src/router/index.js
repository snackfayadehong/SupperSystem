// src/router/index.js
import { createRouter, createWebHistory } from "vue-router";
import MainLayout from "@/layout/MainLayout.vue";
import menuList from "./menu";

//  1. 一次性声明所有 pages
const pages = import.meta.glob("../pages/**/*.vue");
console.log("pages", pages);

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
        children: menuList.map(item => ({
            path: item.path,
            name: item.name,
            component: loadView(item.component)
        }))
    }
];

const router = createRouter({
    history: createWebHistory(),
    routes
});

export default router;
