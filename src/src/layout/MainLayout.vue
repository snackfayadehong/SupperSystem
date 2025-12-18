<template>
    <el-container style="height: 100vh">
        <!-- 顶部导航栏 -->
        <el-header height="56px" class="header">
            <TopBar />
        </el-header>

        <el-container>
            <!-- 左侧菜单 -->
            <el-aside width="220px" class="sidebar">
                <Sidebar />
            </el-aside>

            <!-- 主内容区 -->
            <el-main class="content">
                <div class="page-container">
                    <router-view />
                </div>
            </el-main>
        </el-container>
    </el-container>
</template>

<script>
import TopBar from "./components/TopBar.vue";
import Sidebar from "./components/Sidebar.vue";

export default {
    name: "MainLayout",
    components: {
        TopBar,
        Sidebar
    }
};
</script>

<style>
.siderbar {
    background: #001529;
    color: #fff;
}
.header {
    background: #fff;
    box-shadow: 0 1px 4px rgba(0, 21, 41, 0, 08);
    z-index: 10;
}
.content {
    padding: 16px;
}
.page-container {
    background: #fff;
    border-radius: 6px;
    padding: 16px;
    min-height: calc(100vh - 120px);
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
}
</style>
