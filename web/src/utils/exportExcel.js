import * as XLSXStyle from "xlsx-js-style";

/**
 * 工作量报表导出工具 (Excel)
 * 基于最新 Vue 组件逻辑适配
 * * 导出顺序 & 结构 (总计 29 列):
 * 1. 采购订单 (0-1): 操作员 | 单据数 (无金额)
 * 2. 催货记录 (2-3): 操作员 | 单据数 (无金额)
 * 3. 入库登记 (4-8): 操作员 | 分类 | 品规 | 单据 | 金额
 * 4. 入库验收 (9-13): 操作员 | 分类 | 品规 | 单据 | 金额
 * 5. 出库发放 (14-18): 操作员 | 分类 | 品规 | 单据 | 金额
 * 6. 二级库退 (19-23): 操作员 | 分类 | 品规 | 单据 | 金额
 * 7. 退货汇总 (24-28): 操作员 | 分类 | 品规 | 单据 | 金额
 */

// 1. 基础行背景色 (交替显示，区分人员)
const opBgColors = [
    "BDD7EE", // 蓝
    "C6E0B4", // 绿
    "FFE699", // 金
    "F8CBAD", // 橘
    "E2E2E2", // 灰
    "D9E1F2"  // 淡蓝
];

// 2. 业务表头颜色 (同步 WorkloadStats.vue 的 UI 配色)
const sectionColors = {
    purchase: "9C27B0",  // 紫色 (采购/催货)
    push: "D485C9",      // 浅紫 (催货区分)
    inReg: "409EFF",     // 蓝色 (登记)
    inbound: "67C23A",   // 绿色 (验收)
    outbound: "F56C6C",  // 红色 (出库)
    secondary: "909399", // 灰色 (二级退)
    returns: "E6A23C"    // 黄色 (退货)
};

// 3. 基础单元格样式
const styleBase = {
    font: { name: "宋体", sz: 10.5 },
    alignment: { vertical: "center", horizontal: "center" },
    border: {
        top: { style: "thin", color: { rgb: "888888" } },
        bottom: { style: "thin", color: { rgb: "888888" } },
        left: { style: "thin", color: { rgb: "888888" } },
        right: { style: "thin", color: { rgb: "888888" } }
    }
};

const ExportExcelUtity = (res, workBookName, type) => {
    let workBook = XLSXStyle.utils.book_new();

    if (type === "Workload") {
        const sourceData = Array.isArray(res) ? res : [res];
        const rows = [];
        const totalCols = 29;

        // --- 1. 主标题行 ---
        const mainTitle = Array(totalCols).fill("");
        mainTitle[0] = {
            v: "中心库房全业务量统计明细报表",
            t: "s",
            s: {
                ...styleBase,
                font: { bold: true, sz: 16, color: { rgb: "FFFFFF" } },
                fill: { fgColor: { rgb: "409EFF" } } // 使用主色调蓝
            }
        };
        rows.push(mainTitle);

        // --- 2. 业务分组表头 (带颜色) ---
        const sectionHeader = Array(totalCols).fill("");
        // 辅助函数：生成带背景色的表头样式 (字体白色加粗)
        const getSecStyle = (rgb) => ({
            ...styleBase,
            font: { bold: true, sz: 11, color: { rgb: "FFFFFF" } },
            fill: { fgColor: { rgb } }
        });

        sectionHeader[0] = { v: "发送采购订单", t: "s", s: getSecStyle(sectionColors.purchase) };
        sectionHeader[2] = { v: "催货记录", t: "s", s: getSecStyle(sectionColors.push) };
        sectionHeader[4] = { v: "入库登记汇总", t: "s", s: getSecStyle(sectionColors.inReg) };
        sectionHeader[9] = { v: "入库验收汇总", t: "s", s: getSecStyle(sectionColors.inbound) };
        sectionHeader[14] = { v: "出库发放汇总", t: "s", s: getSecStyle(sectionColors.outbound) };
        sectionHeader[19] = { v: "二级库退库汇总", t: "s", s: getSecStyle(sectionColors.secondary) };
        sectionHeader[24] = { v: "退货汇总", t: "s", s: getSecStyle(sectionColors.returns) };
        rows.push(sectionHeader);

        // --- 3. 字段列名表头 ---
        const simpleCols = ["操作人员", "单据数"];
        const fullCols = ["操作人员", "材料分类", "品规数", "单据数", "总金额"];

        const subHeaderRow = [
            ...simpleCols, // 采购
            ...simpleCols, // 催货
            ...fullCols,   // 登记
            ...fullCols,   // 验收
            ...fullCols,   // 出库
            ...fullCols,   // 二级退
            ...fullCols    // 退货
        ].map(name => ({
            v: name,
            t: "s",
            s: { ...styleBase, font: { bold: true }, fill: { fgColor: { rgb: "F2F2F2" } } }
        }));
        rows.push(subHeaderRow);

        // --- 4. 数据填充 ---
        sourceData.forEach((opGroup, opIdx) => {
            const bgColor = opBgColors[opIdx % opBgColors.length];
            const currentStyle = { ...styleBase, fill: { fgColor: { rgb: bgColor } } };
            const nameStyle = { ...currentStyle, font: { ...currentStyle.font, bold: true } }; // 姓名加粗
            const amountStyle = { ...currentStyle, alignment: { horizontal: "right" } }; // 金额右对齐

            // 提取数据 (确保与 Vue 组件字段一致)
            const pur = opGroup.purchase || [];
            const push = opGroup.push || [];
            const inReg = opGroup.inReg || [];
            const inb = opGroup.inbound || [];
            const outb = opGroup.outbound || [];
            const sec = opGroup.secondaryRefund || []; // 注意字段名: secondaryRefund
            const ret = opGroup.return || [];

            // 计算最大行数，决定该人员占几行
            const maxLen = Math.max(pur.length, push.length, inReg.length, inb.length, outb.length, sec.length, ret.length);

            for (let i = 0; i < maxLen; i++) {
                const row = Array(totalCols).fill({ v: "", t: "s", s: currentStyle });

                // 辅助函数1：填充简单列 (采购/催货)
                const fillSimple = (data, offset) => {
                    if (data[i]) {
                        row[offset] = { v: opGroup.operator, s: nameStyle };
                        row[offset + 1] = { v: data[i].billCount, s: currentStyle };
                    } else if (i === 0) {
                        // 第一行必须显示名字，即使该业务无数据
                        row[offset] = { v: opGroup.operator, s: nameStyle };
                    }
                };

                // 辅助函数2：填充完整列 (带金额的业务)
                const fillFull = (data, offset) => {
                    if (data[i]) {
                        row[offset] = { v: opGroup.operator, s: nameStyle };
                        row[offset + 1] = { v: data[i].category, s: currentStyle };
                        row[offset + 2] = { v: data[i].specCount, s: currentStyle };
                        row[offset + 3] = { v: data[i].billCount, s: currentStyle };
                        row[offset + 4] = { v: data[i].totalAmount ? data[i].totalAmount.toFixed(2) : "0.00", s: amountStyle };
                    } else if (i === 0) {
                        row[offset] = { v: opGroup.operator, s: nameStyle };
                    }
                };

                // 按指定顺序填充
                fillSimple(pur, 0);   // 1. 采购
                fillSimple(push, 2);  // 2. 催货
                fillFull(inReg, 4);   // 3. 登记
                fillFull(inb, 9);     // 4. 验收
                fillFull(outb, 14);   // 5. 出库
                fillFull(sec, 19);    // 6. 二级退
                fillFull(ret, 24);    // 7. 退货

                rows.push(row);
            }

            // 人员之间插入空行 (视觉分隔)
            if (opIdx < sourceData.length - 1) {
                const emptyRow = Array(totalCols).fill({ v: "", t: "s", s: { fill: { fgColor: { rgb: "FFFFFF" } } } });
                rows.push(emptyRow);
            }
        });

        const ws = XLSXStyle.utils.aoa_to_sheet(rows);

        // --- 5. 设置单元格合并 ---
        ws["!merges"] = [
            { s: { r: 0, c: 0 }, e: { r: 0, c: totalCols - 1 } }, // 总标题
            { s: { r: 1, c: 0 }, e: { r: 1, c: 1 } },   // 采购
            { s: { r: 1, c: 2 }, e: { r: 1, c: 3 } },   // 催货
            { s: { r: 1, c: 4 }, e: { r: 1, c: 8 } },   // 登记
            { s: { r: 1, c: 9 }, e: { r: 1, c: 13 } },  // 验收
            { s: { r: 1, c: 14 }, e: { r: 1, c: 18 } }, // 出库
            { s: { r: 1, c: 19 }, e: { r: 1, c: 23 } }, // 二级退
            { s: { r: 1, c: 24 }, e: { r: 1, c: 28 } }  // 退货
        ];

        // --- 6. 智能列宽设置 ---
        // 逻辑：所有 "材料分类" 列 (索引 5, 10, 15, 20, 25) 设宽 (180px)，名字列 (0,2,4,9...) 适中 (100px)，其他窄一点
        const catCols = [5, 10, 15, 20, 25];
        const nameCols = [0, 2, 4, 9, 14, 19, 24];

        ws["!cols"] = Array(totalCols).fill(0).map((_, i) => {
            if (catCols.includes(i)) return { wpx: 180 };
            if (nameCols.includes(i)) return { wpx: 100 };
            return { wpx: 80 };
        });

        XLSXStyle.utils.book_append_sheet(workBook, ws, "全院业务量统计");
    }

    return XLSXStyle.writeFile(workBook, workBookName);
};

export default ExportExcelUtity;