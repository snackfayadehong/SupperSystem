package model

// 库房工作量查询模型

type WorkloadRequest struct {
	StartTime string `json:"startTime" binding:"required"`
	EndTime   string `json:"endTime" binding:"required"`
}

// RawWorkloadRow 对应 SQL 直接查询出的扁平行
type RawWorkloadRow struct {
	OperatorName   string  `gorm:"column:OperatorName"`
	OperationType  int     `gorm:"column:OperationType"`
	StorehouseCode string  `gorm:"column:StorehouseCode"`
	SpecCount      int     `gorm:"column:SpecCount"`
	BillCount      int     `gorm:"column:BillCount"`
	TotalAmount    float64 `gorm:"column:TotalAmount"`
	FallbackName   string  `gorm:"column:FallbackName"`
}

// WorkloadDetail 对应前端明细表格
type WorkloadDetail struct {
	Category    string  `json:"category"`
	SpecCount   int     `json:"specCount"`
	BillCount   int     `json:"billCount"`
	TotalAmount float64 `json:"totalAmount"`
}

// WorkloadGroup 最终聚合结构
type WorkloadGroup struct {
	Operator        string           `json:"operator"`
	Inbound         []WorkloadDetail `json:"inbound"`         // 入库验收
	Outbound        []WorkloadDetail `json:"outbound"`        // 出库
	Return          []WorkloadDetail `json:"return"`          // 退货
	InReg           []WorkloadDetail `json:"inReg"`           // 入库登记
	Purchase        []WorkloadDetail `json:"purchase"`        // 采购订单发送
	Push            []WorkloadDetail `json:"push"`            // 采购订单催货
	SecondaryRefund []WorkloadDetail `json:"secondaryRefund"` // 二级库退库
}
