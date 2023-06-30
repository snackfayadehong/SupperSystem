package model

// 高值库房工作量查询模型

// UserProdAccept 验收
type UserProdAccept struct {
	Name        string `gorm:"column:MEnName"`
	ProdAcSpec  int    `gorm:"column:prodSpecNum"`
	ProdAcBill  int    `gorm:"column:billNum"`
	ProdAcTotal string `gorm:"column:totalAmount"`
}

// DepartmentCollar 出库
type DepartmentCollar struct {
	Name        string `gorm:"column:BLMakerName"`
	ProdDpSpec  int    `gorm:"column:DpSpecNum"`
	ProdDpBill  int    `gorm:"column:DpBillNum"`
	ProdDpTotal string `gorm:"column:DpTotalAmount"`
}

// RefundProd 退货
type RefundProd struct {
	Name     string `gorm:"column:EmployeeName"`
	RefSpec  int    `gorm:"column:ReFSpecNum"`
	RefBill  int    `gorm:"column:RefBillNum"`
	RefTotal string `gorm:"column:RefTotalAmount"`
}
type UserWorkloadInfo struct {
	Name             string `json:"name"`
	ProdAcSpec       int    `json:"prodAcSpec"`
	ProdAcBill       int    `json:"prodAcBill"`
	ProdAcTotal      string `json:"prodAcTotal"`
	ProdDpSpec       int    `json:"prodDpSpec"`
	ProdDpBill       int    `json:"prodDpBill"`
	ProdDpTotal      string `json:"prodDpTotal"`
	RefSpec          int    `json:"refSpec"`
	RefBill          int    `json:"refBill"`
	RefTotal         string `json:"refTotal"`
	TotalBillAmount  int    `json:"total_bill_amount"`
	TotalSpecAmount  int    `json:"total_spec_amount"`
	TotalTotalAmount string `json:"total_total_amount"`
}
