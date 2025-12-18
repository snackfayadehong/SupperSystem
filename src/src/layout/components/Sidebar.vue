<template>
    <el-menu router :default-active="$route.path">
        <template v-for="item in menu" :key="item.path">
            <el-sub-menu v-if="item.children && item.children.length" :index="item.path">
                <template #title>
                    {{ item.label }}
                </template>

                <el-menu-item v-for="child in item.children" :key="child.path" :index="child.path">
                    {{ child.label }}
                </el-menu-item>
            </el-sub-menu>

            <el-menu-item v-else :index="item.path">
                {{ item.label }}
            </el-menu-item>
        </template>
    </el-menu>
</template>

<script>
import menu from "@/router/menu";

export default {
    name: "Sidebar",
    data() {
        return {
            menu
        };
    }
};
</script>
