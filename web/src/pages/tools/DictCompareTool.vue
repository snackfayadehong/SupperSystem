<template>
    <div class="dict-compare-container">
        <section class="unified-hero-card">
            <div class="hero-content">
                <div class="text-area">
                    <h1 class="hero-title">字典数据实时对比工具</h1>
                    <p class="hero-desc">怡道系统与 HIS 系统字典信息校对</p>
                </div>

                <div class="search-area glass-effect">
                    <el-input v-model="keyword" placeholder="请输入 14位材料代码 或 6位产品ID" class="integrated-input" clearable maxlength="14" @input="handleInput" @keyup.enter="handleCompare">
                        <template #prefix
                            ><el-icon>
                                <Search /> </el-icon
                        ></template>
                        <template #suffix>
                            <el-tag :type="keywordType === '产品ID' ? 'success' : 'primary'" size="small" effect="dark" class="mode-tag">
                                {{ keywordType }}
                            </el-tag>
                        </template>
                    </el-input>
                    <el-button type="primary" :loading="loading" class="audit-btn" @click="handleCompare"> 开始对比 </el-button>
                </div>
            </div>
        </section>

        <el-dialog v-model="choiceVisible" title="核对冲突：请确认具体条目" width="850px" append-to-body destroy-on-close class="selection-dialog">
            <div class="selection-header">
                <el-alert
                    title="检测到多个怡道系统产品信息"
                    type="warning"
                    description="该材料代码在怡道系统对应多个产品 ID，请通过下方的名称与规格进行确认，点击对应行进行对比。"
                    show-icon
                    :closable="false"
                />
            </div>

            <div class="selection-scroller">
                <div v-for="item in multiOptions" :key="item.ProductInfoID" class="choice-card" @click="confirmSelection(item)">
                    <div class="card-meta">
                        <span class="meta-id">产品ID: {{ item.ProductInfoID }}</span>
                        <span class="meta-code">材料代码：{{ item.ypdm }}</span>
                    </div>

                    <div class="card-body">
                        <div class="item-name">{{ item.ypmc }}</div>
                        <div class="item-spec">{{ item.ypgg || "规格型号未维护" }}</div>
                    </div>

                    <div class="card-arrow">
                        <el-icon>
                            <ArrowRightBold />
                        </el-icon>
                    </div>
                </div>
            </div>
        </el-dialog>

        <transition name="list-fade">
            <div v-if="compareResults.length" class="results-container">
                <div class="summary-bar" :class="{ 'has-mismatch': mismatchCount > 0 }">
                    <div class="summary-left">
                        <el-icon class="status-icon">
                            <InfoFilled />
                        </el-icon>
                        <span class="summary-text"
                            >比对完成：共 {{ compareResults.length }} 项，发现 <b class="danger">{{ mismatchCount }}</b> 项差异</span
                        >
                    </div>
                    <el-tag :type="mismatchCount === 0 ? 'success' : 'danger'" effect="dark" class="status-pill">
                        {{ mismatchCount === 0 ? "数据一致" : "检测到冲突" }}
                    </el-tag>
                </div>

                <div class="diff-view">
                    <div class="diff-header">
                        <div class="col-name">属性</div>
                        <div class="col-data">怡道系统</div>
                        <div class="col-status">状态</div>
                        <div class="col-data">HIS 系统</div>
                    </div>

                    <div class="diff-card identity-row">
                        <div class="col-name">
                            <span class="label">产品 ID</span>
                            <span class="key">PRODUCTINFOID</span>
                        </div>
                        <div class="col-data val-box identity-box">{{ currentSelection.ProductInfoID || "-" }}</div>
                        <div class="col-status">
                            <el-icon class="icon-match">
                                <CircleCheckFilled />
                            </el-icon>
                        </div>
                        <div class="col-data val-box identity-box">{{ currentSelection.ProductInfoID || "-" }}</div>
                    </div>

                    <div class="diff-card identity-row">
                        <div class="col-name">
                            <span class="label">材料代码</span>
                            <span class="key">CODE</span>
                        </div>
                        <div class="col-data val-box identity-box">{{ currentSelection.ypdm || "-" }}</div>
                        <div class="col-status">
                            <el-icon class="icon-match">
                                <CircleCheckFilled />
                            </el-icon>
                        </div>
                        <div class="col-data val-box identity-box">{{ currentSelection.ypdm || "-" }}</div>
                    </div>

                    <div v-for="(item, index) in compareResults" :key="item.field" class="diff-card" :class="{ 'mismatch-card': !item.isMatch }" :style="{ animationDelay: index * 0.05 + 's' }">
                        <div class="col-name">
                            <span class="label">{{ item.label }}</span>
                            <span class="key">{{ item.field.toUpperCase() }}</span>
                        </div>
                        <div class="col-data val-box">{{ item.localValue || "-" }}</div>
                        <div class="col-status">
                            <el-icon v-if="item.isMatch" class="icon-match">
                                <CircleCheckFilled />
                            </el-icon>
                            <el-icon v-else class="icon-mismatch">
                                <WarningFilled />
                            </el-icon>
                        </div>
                        <div class="col-data val-box" :class="{ 'error-box': !item.isMatch }">
                            <span class="val-text">{{ item.hisValue || "-" }}</span>
                            <span v-if="!item.isMatch" class="error-label">冲突项</span>
                        </div>
                    </div>
                </div>
            </div>
        </transition>

        <el-empty v-if="!loading && hasSearched && compareResults.length === 0" :description="errorMsg || '未检索到该编码的比对数据'" />
    </div>
</template>

<script setup>
import { ref, computed, h } from "vue";
import { Search, InfoFilled, CircleCheckFilled, WarningFilled, ArrowRightBold } from "@element-plus/icons-vue";
import myAxios from "@/services/myAxios";
import { ElMessage } from "element-plus";

// --- 状态定义 ---
const keyword = ref("");
const loading = ref(false);
const hasSearched = ref(false);
const compareResults = ref([]);
const errorMsg = ref("");
const choiceVisible = ref(false);
const multiOptions = ref([]);
const currentSelection = ref({ ProductInfoID: "", ypdm: "" });

// --- 计算属性 ---
const keywordType = computed(() => {
    if (!keyword.value) return "产品ID/材料代码";
    const isPureNumber = /^\d+$/.test(keyword.value);
    return isPureNumber && parseInt(keyword.value) <= 2147483647 ? "产品ID" : "材料代码";
});

const mismatchCount = computed(() => {
    return (compareResults.value || []).filter(item => item && !item.isMatch).length;
});

// --- 逻辑函数 ---
const handleInput = val => {
    keyword.value = val.replace(/[^a-zA-Z0-9]/g, "");
};

let msgInstance = null;
let warnCount = 0;
let resetTimer = null;

const handleCompare = async () => {
    if (!keyword.value) {
        warnCount = Math.min(warnCount + 1, 99);
        const tipContent = h("div", { style: "display: flex; align-items: center; gap: 8px;" }, [
            h("span", null, warnCount > 1 ? "请勿重复点击！" : "请输入对账关键字"),
            warnCount > 1 ? h("span", { class: "singleton-badge" }, warnCount) : null
        ]);
        if (msgInstance) msgInstance.close();
        msgInstance = ElMessage({ message: tipContent, type: "warning" });
        clearTimeout(resetTimer);
        resetTimer = setTimeout(() => {
            warnCount = 0;
        }, 3000);
        return;
    }

    loading.value = true;
    compareResults.value = [];
    hasSearched.value = true;

    try {
        const res = await myAxios.post("/dict/compare", { keyword: keyword.value });
        if (res.code === 201) {
            multiOptions.value = res.data || [];
            choiceVisible.value = true;
        } else if (res.code === 0) {
            compareResults.value = Array.isArray(res.data.results) ? res.data.results : [];
            currentSelection.value = {
                ProductInfoID: res.data.ProductInfoID,
                ypdm: res.data.ypdm
            };
            console.log("比对结果：", currentSelection.value);
        } else {
            errorMsg.value = res.message;
        }
    } catch (err) {
        errorMsg.value = "连接 HIS 接口异常";
    } finally {
        loading.value = false;
    }
};

const confirmSelection = row => {
    currentSelection.value = { ProductInfoID: row.ProductInfoID, ypdm: row.ypdm };
    keyword.value = row.ProductInfoID.toString();
    choiceVisible.value = false;
    handleCompare();
};
</script>

<style scoped>
/* 1. 基础布局 & Hero 卡片 */
.dict-compare-container {
    padding: 24px;
    background-color: #f5f7fa;
    min-height: 100vh;
    display: flex;
    flex-direction: column;
    gap: 20px;
}

.unified-hero-card {
    padding: 40px;
    background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
    border-radius: 24px;
    color: white;
    box-shadow: 0 12px 24px rgba(79, 172, 254, 0.2);
}

.hero-content {
    max-width: 900px;
    margin: 0 auto;
    text-align: center;
}

.hero-title {
    font-size: 28px;
    font-weight: 800;
    margin: 0;
}

.hero-desc {
    font-size: 15px;
    opacity: 0.9;
    margin: 10px 0 30px;
}

.glass-effect {
    background: rgba(255, 255, 255, 0.2);
    backdrop-filter: blur(10px);
    padding: 12px;
    border-radius: 16px;
    display: flex;
    gap: 12px;
    border: 1px solid rgba(255, 255, 255, 0.3);
}

.integrated-input {
    flex: 1;
}

:deep(.el-input__wrapper) {
    border-radius: 10px;
    padding: 0 15px;
}

/* 2. 多项目选择弹窗美化 */
.selection-scroller {
    display: flex;
    flex-direction: column;
    gap: 16px;
    /* 增加条目间距 */
    margin-top: 10px;
}

.choice-card {
    display: flex;
    align-items: center;
    padding: 20px 25px;
    background: #fff;
    border: 1px solid #e2e8f0;
    border-radius: 16px;
    cursor: pointer;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.choice-card:hover {
    transform: translateX(10px);
    border-color: #4facfe;
    background: #f0f9ff;
}

.card-meta {
    width: 150px;
    display: flex;
    flex-direction: column;
    gap: 4px;
    border-right: 1px solid #edf2f7;
    margin-right: 20px;
}

.meta-id {
    font-family: "JetBrains Mono", monospace;
    font-weight: 800;
    color: #409eff;
    font-size: 17px;
}

.meta-code {
    font-size: 12px;
    color: #94a3b8;
}

.card-body {
    flex: 1;
}

.item-name {
    font-weight: 800;
    color: #1e293b;
    font-size: 16px;
}

.item-spec {
    font-size: 13px;
    color: #64748b;
}

/* 3. 比对结果区核心优化 */
.results-container {
    max-width: 1200px;
    margin: 0 auto;
    width: 100%;
}

.summary-bar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    background: #fff;
    padding: 16px 25px;
    border-radius: 12px;
    margin-bottom: 16px;
    border-left: 6px solid #409eff;
}

.summary-bar.has-mismatch {
    border-left-color: #f56c6c;
    background: #fffbfa;
}

.summary-text {
    font-weight: bold;
    color: #475569;
}

.danger {
    color: #f56c6c;
    font-size: 20px;
    margin: 0 4px;
}

.status-pill {
    font-weight: 900;
    padding: 0 14px;
}

/* 扁平化对齐表头并设置居中 */
.diff-header {
    display: grid;
    grid-template-columns: 180px 1fr 100px 1fr;
    padding: 0 24px 12px;
    font-size: 13px;
    color: #3966a5;
    font-weight: 800;
    text-transform: uppercase;
    letter-spacing: 1px;
}

/* 居中比对系统的标题 */
.diff-header .col-data {
    text-align: center;
    font-size: 20px;
}

.diff-header .col-status {
    text-align: center;
}

.diff-card {
    display: grid;
    grid-template-columns: 180px 1fr 100px 1fr;
    align-items: center;
    background: #fff;
    padding: 24px;
    border-radius: 16px;
    margin-bottom: 12px;
    transition: all 0.3s ease;
    border: 1px solid transparent;
}

/* 身份标识行特殊样式 */
.identity-row {
    background: #f0f7ff;
    border-left: 4px solid #38bdf8;
}

.identity-box {
    background: #e6f1fe !important;
    color: #0284c7 !important;
    font-weight: 800;
}

.diff-card.mismatch-card {
    border: 1px solid #ffccc7;
    background: #fffbfa;
}

/* 属性标签增强 */
.label {
    font-weight: 900;
    color: #1e293b;
    display: block;
    font-size: 18px;
    /* 加大字段名 */
    margin-bottom: 2px;
}

.key {
    font-size: 12px;
    color: #3074c7;
    font-family: "JetBrains Mono", monospace;
}

.val-box {
    padding: 14px 20px;
    background: #f8fafc;
    border-radius: 10px;
    font-family: "JetBrains Mono", "Consolas", monospace;
    font-size: 15px;
    font-weight: 600;
    color: #475569;
}

/* 冲突高亮 */
.error-box {
    background: #fff1f0;
    color: #ff4d4f;
    position: relative;
    border: 1px dashed #ffa39e;
}

.error-label {
    position: absolute;
    top: -10px;
    right: 10px;
    background: #ff4d4f;
    color: white;
    padding: 2px 8px;
    border-radius: 6px;
    font-size: 10px;
    font-weight: 900;
    box-shadow: 0 2px 6px rgba(255, 77, 79, 0.2);
}

.col-status {
    display: flex;
    justify-content: center;
    font-size: 26px;
}

.icon-match {
    color: #52c41a;
}

.icon-mismatch {
    color: #ff4d4f;
    filter: drop-shadow(0 0 4px rgba(255, 77, 79, 0.2));
}

.singleton-badge {
    background: #f56c6c;
    color: #fff;
    padding: 0 6px;
    border-radius: 10px;
    font-size: 10px;
}
</style>
