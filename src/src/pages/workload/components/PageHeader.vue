<template>
    <el-card shadow="never" class="page-header">
        <div class="header">
            <div class="title">{{ title }}</div>
            <div class="actions">
                <slot />
            </div>
        </div>
    </el-card>
</template>

<script setup>
defineProps({
    title: String
});
</script>

<style scoped>
.header {
    display: flex;
    justify-content: space-between;
    align-items: center;
}
.title {
    font-size: 18px;
    font-weight: bold;
}
</style>
