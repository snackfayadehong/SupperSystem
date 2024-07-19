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
