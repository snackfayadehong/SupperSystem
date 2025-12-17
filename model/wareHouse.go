package model

// 库房工作量查询模型

type WorkloadRequest struct {
	StartTime string `json:"startTime" binding:"required"`
	EndTime   string `json:"endTime" binding:"required"`
}
type WorkloadData struct {
	OperatorName  string  `gorm:"column:OperatorName" json:"operator_name"`
	OperationType string  `gorm:"column:OperationType" json:"operation_type"`
	Category      string  `gorm:"column:Category" json:"category"`
	SpecCount     int     `gorm:"column:SpecCount" json:"spec_count"`
	BillCount     int     `gorm:"column:BillCount" json:"bill_count"`
	TotalAmount   float64 `gorm:"column:TotalAmount" json:"total_amount"`
}

// UserWorkloadInfoNew 2024-07-18 加入普耗工作量数据，重新调整工作量查询逻辑
type UserWorkloadInfoNew struct {
	MenName          string  `json:"name" gorm:"column:MEnName"`
	In_BillNum       int     `json:"prodAcBill" gorm:"column:in_BillNum"`
	In_ProdSpecNum   int     `json:"prodAcSpec" gorm:"column:in_ProdSpecNum"`
	In_TotalAmount   float32 `json:"prodAcTotal" gorm:"column:in_TotalAmount"`
	Out_BillNum      int     `json:"prodDpBill" gorm:"column:out_BillNum"`
	Out_ProdSpecNum  int     `json:"prodDpSpec" gorm:"column:out_ProdSpecNum"`
	Out_TotalAmount  float32 `json:"prodDpTotal" gorm:"column:out_TotalAmount"`
	Back_BillNum     int     `json:"refBill" gorm:"column:back_BillNum"`
	Back_ProdSpecNum int     `json:"refSpec" gorm:"column:back_ProdSpecNum"`
	Back_TotalAmount float32 `json:"refTotal" gorm:"column:back_TotalAmount"`
	TotalBillAmount  int     `json:"total_bill_amount"`
	TotalSpecAmount  int     `json:"total_spec_amount"`
	TotalTotalAmount float32 `json:"total_total_amount"`
}
