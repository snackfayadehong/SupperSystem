package controller

import (
	clientDb "WorkloadQuery/db"
	"WorkloadQuery/model"
	"time"
)

type PurchaseSummaryTimeInterval struct {
	StartTime string `json:"startTime" binding:"required"`
	EndTime   string `json:"endTime" binding:"required"`
}

// NotDeliveredPurchaseSummary 未到货订单信息
func (pr *PurchaseSummaryTimeInterval) NotDeliveredPurchaseSummary() (purchaseSummary []model.AllPurchaseSummary) {
	purchaseSummary = NoDeliveredPurchaseSummary(&pr.StartTime, &pr.EndTime)
	return purchaseSummary
}

// NoDeliveredPurchaseSummary 根据未到货信息找到 未到货订单明细数据
func NoDeliveredPurchaseSummary(startTime *string, endTime *string) (purchaseSummaryHasChildren []model.AllPurchaseSummary) {
	var id []string // purchaseSummaryId 订单ID 用于查找明细
	clientDb.DB.Raw(clientDb.NotDeliveredPurchaseSummarySql, *startTime, *endTime).Scan(&purchaseSummaryHasChildren)
	if len(purchaseSummaryHasChildren) < 1 {
		return nil
	}
	// i 时间格式化
	// v 写入id 根据 purchaseSummaryId 查找明细
	for i, v := range purchaseSummaryHasChildren {
		tempTime, _ := time.Parse("2006-01-02T15:04:05Z", (purchaseSummaryHasChildren)[i].AuditorDate)
		(purchaseSummaryHasChildren)[i].AuditorDate = tempTime.Format("2006-01-02 15:04:05")
		id = append(id, v.PurchaseSummaryID) // 采购订单单据ID
	}
	// 根据 purchaseSummaryId 查找明细
	var purchaseSummaryDetail *[]model.PurchaseSummaryDetail
	clientDb.DB.Raw(clientDb.NotDeliveredPurchaseSummaryDetailSql, id).Find(&purchaseSummaryDetail)
	// 循环 根据 purchaseSummaryId 匹配明细
	// for i1, v1 := range *purchaseSummaryHasChildren {
	// 	for i2, v2 := range *purchaseSummaryDetail {
	// 		if v1.PurchaseSummaryID == v2.PurchaseSummaryID {
	// 			(*purchaseSummaryHasChildren)[i1].Children = append((*purchaseSummaryHasChildren)[i1].Children, v2)
	// 		}
	// 	}
	// }
	// 明细map 方便根据PurchaseSummaryId查找匹配记录后,写入 AllPurchaseSummary 嵌套结构体
	detailMap := make(map[string][]model.PurchaseSummaryDetail)
	for _, v2 := range *purchaseSummaryDetail {
		detailMap[v2.PurchaseSummaryID] = append(detailMap[v2.PurchaseSummaryID], v2)
	}
	// 循环 purchaseSummaryHasChildren
	for i1, v1 := range purchaseSummaryHasChildren {
		if details, ok := detailMap[v1.PurchaseSummaryID]; ok {
			(purchaseSummaryHasChildren)[i1].Children = &details
		}
	}
	return
}
