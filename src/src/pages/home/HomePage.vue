<template>
    <div class="home-page">
        <StatCards />
        <ChartPanel />
        <QuickEntry />
    </div>
</template>

<script setup>
import StatCards from "./components/StatCards.vue";
import ChartPanel from "./components/ChartPanel.vue";
import QuickEntry from "./components/QuickEntry.vue";
</script>

<style scoped>
.home-page {
    display: flex;
    flex-direction: column;
    gap: 16px;
}
</style>
