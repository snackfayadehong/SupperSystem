import * as XLSXStyle from "xlsx-js-style";

/**
 * 工作量报表导出工具
 * 结构：入库(0-4列) | 出库(5-9列) | 退还(10-14列)
 */

// 1. 调深后的商务配色方案 (饱和度更高，清晰区分人员)
const opBgColors = [
    "BDD7EE", // 商务蓝
    "C6E0B4", // 雅致绿
    "FFE699", // 暖金
    "F8CBAD", // 浅橘
    "D9E1F2", // 冰蓝
    "E2E2E2"  // 浅灰
];

// 2. 业务分组颜色 (表头专用，饱和度最高)
const sectionColors = {
    inbound: "A9D08E",  // 绿色 (入库)
    outbound: "F4B084", // 橙色 (出库)
    returns: "FFD966"   // 黄色 (退还)
};

// 3. 基础样式定义
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

        // --- 标题行 (第1行) ---
        const mainTitle = Array(15).fill("");
        mainTitle[0] = { 
            v: "中心库房业务量统计明细报表", 
            t: "s", 
            s: { 
                ...styleBase, 
                font: { bold: true, sz: 14, color: { rgb: "FFFFFF" } }, 
                fill: { fgColor: { rgb: "4472C4" } } 
            } 
        };
        rows.push(mainTitle);

        // --- 业务分组头 (第2行) ---
        const sectionHeader = Array(15).fill("");
        const getSecStyle = (rgb) => ({ ...styleBase, font: { bold: true, sz: 11 }, fill: { fgColor: { rgb } } });
        
        sectionHeader[0] = { v: "入库信息汇总", t: "s", s: getSecStyle(sectionColors.inbound) };
        sectionHeader[5] = { v: "出库信息汇总", t: "s", s: getSecStyle(sectionColors.outbound) };
        sectionHeader[10] = { v: "退还信息汇总", t: "s", s: getSecStyle(sectionColors.returns) };
        rows.push(sectionHeader);

        // --- 明细表头 (第3行) ---
        const subNames = ["操作人员", "材料分类", "品规数", "单据数", "总金额"];
        const subHeaderRow = [...subNames, ...subNames, ...subNames].map(name => ({
            v: name, 
            t: "s", 
            s: { ...styleBase, font: { bold: true }, fill: { fgColor: { rgb: "F2F2F2" } } }
        }));
        rows.push(subHeaderRow);

        // --- 数据行处理 (从第4行开始) ---
        sourceData.forEach((opGroup, opIdx) => {
            const bgColor = opBgColors[opIdx % opBgColors.length];
            const currentStyle = { ...styleBase, fill: { fgColor: { rgb: bgColor } } };
            const nameStyle = { ...currentStyle, font: { ...currentStyle.font, bold: true } }; // 操作员姓名加粗
            const amountStyle = { ...currentStyle, alignment: { horizontal: "right" } };    // 金额右对齐

            const inb = opGroup.inbound || [];
            const outb = opGroup.outbound || [];
            const ret = opGroup.return || [];
            
            // 计算当前操作员需要的最大行数
            const maxLen = Math.max(inb.length, outb.length, ret.length);

            for (let i = 0; i < maxLen; i++) {
                const row = Array(15).fill({ v: "", t: "s", s: currentStyle });

                // 1. 填充入库列 (0-4列)
                if (inb[i]) {
                    row[0] = { v: opGroup.operator, s: nameStyle };
                    row[1] = { v: inb[i].category, s: currentStyle };
                    row[2] = { v: inb[i].specCount, s: currentStyle };
                    row[3] = { v: inb[i].billCount, s: currentStyle };
                    row[4] = { v: inb[i].totalAmount.toFixed(2), s: amountStyle };
                } else if (i === 0) { row[0] = { v: opGroup.operator, s: nameStyle }; }

                // 2. 填充出库列 (5-9列)
                if (outb[i]) {
                    row[5] = { v: opGroup.operator, s: nameStyle };
                    row[6] = { v: outb[i].category, s: currentStyle };
                    row[7] = { v: outb[i].specCount, s: currentStyle };
                    row[8] = { v: outb[i].billCount, s: currentStyle };
                    row[9] = { v: outb[i].totalAmount.toFixed(2), s: amountStyle };
                } else if (i === 0) { row[5] = { v: opGroup.operator, s: nameStyle }; }

                // 3. 填充退还列 (10-14列)
                if (ret[i]) {
                    row[10] = { v: opGroup.operator, s: nameStyle };
                    row[11] = { v: ret[i].category, s: currentStyle };
                    row[12] = { v: ret[i].specCount, s: currentStyle };
                    row[13] = { v: ret[i].billCount, s: currentStyle };
                    row[14] = { v: ret[i].totalAmount.toFixed(2), s: amountStyle };
                } else if (i === 0) { row[10] = { v: opGroup.operator, s: nameStyle }; }

                rows.push(row);
            }

            // --- 核心修改：在不同人员之间空一行 ---
            // 如果不是最后一个人，增加一个空白行作为分隔
            if (opIdx < sourceData.length - 1) {
                const emptyRow = Array(15).fill({ v: "", t: "s", s: { fill: { fgColor: { rgb: "FFFFFF" } } } });
                rows.push(emptyRow);
            }
        });

        const ws = XLSXStyle.utils.aoa_to_sheet(rows);

        // --- 设置单元格合并 ---
        ws["!merges"] = [
            { s: { r: 0, c: 0 }, e: { r: 0, c: 14 } }, // 总标题
            { s: { r: 1, c: 0 }, e: { r: 1, c: 4 } },  // 入库分组
            { s: { r: 1, c: 5 }, e: { r: 1, c: 9 } },  // 出库分组
            { s: { r: 1, c: 10 }, e: { r: 1, c: 14 } } // 退还分组
        ];

        // --- 设置列宽 ---
        ws["!cols"] = Array(15).fill(0).map((_, i) => ({ wpx: (i % 5 === 1) ? 180 : 85 }));

        XLSXStyle.utils.book_append_sheet(workBook, ws, "操作员业务量统计");
    }

    return XLSXStyle.writeFile(workBook, workBookName);
};

export default ExportExcelUtity;