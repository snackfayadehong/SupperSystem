<template>
    <div class="demo-date-picker">
        <div class="block">
            <p class="dateTime-select">时间区间:</p>
            <el-config-provider>
                <el-date-picker v-model="queryDate" type="daterange" value-format="YYYY-MM-DD" start-placeholder="开始日期" end-placeholder="结束日期" />
            </el-config-provider>
        </div>
        <el-button type="primary" class="qry-btn" @click="query">查询</el-button>
    </div>
</template>

<script setup>
import myAxios from "../../api/myAxios.js";
import { ref } from "vue";
import bus from "../../eventBus";
import dayjs from "dayjs";
import { getCurrentInstance } from "vue";
// 格式化时间
let end = dayjs().format("YYYY-MM-DD");
let start = dayjs().subtract(6, "day").format("YYYY-MM-DD");
// 时间
const queryDate = ref("");
queryDate.value = [start, end];

// 调用组件名
let commName = "";
let currentCpn = getCurrentInstance();

// 查询工作量数据
const query = async () => {
    // 获取父组件名称,根据夫组件名称查询不同数据
    commName = currentCpn.parent.proxy.name;
    console.log(commName);
    let res = "";
    let value = JSON.parse(JSON.stringify(queryDate.value));
    if (commName == "") {
        alert("组件名为空");
        return;
    }
    if (commName == "workload") {
        // 工作量查询
        let res = await myAxios.post("/getWorkload", { startTime: value[0], endTime: value[1] });
        if (res == "") {
            alert("无数据");
        }
    } else if (commName == "departmentCollar") {
        // 未上账单据查询
        let res = await myAxios.post("/getNoAccountEntry", { startTime: value[0], endTime: value[1] });
        if (res == "") {
            alert("无数据");
        }
    }
    bus.emit("getData", res.data.Data);
};
</script>
<style scoped>
.dateTime-select {
    margin: 0;
}
.demo-date-picker {
    display: flex;
    align-items: center;
    width: 100%;
    padding: 0;
    flex-wrap: wrap;
    font-size: 14px;
}
.demo-date-picker .block {
    padding-bottom: 15px;
    text-align: left;
    /* border-right: solid 1px var(--el-border-color); */
    flex: 1;
}
.demo-date-picker .qry-btn {
    margin-right: 80px;
}
.demo-date-picker .block:last-child {
    border-right: none;
}
</style>
