package controller

import (
	"WorkloadQuery/db"
	"WorkloadQuery/model"
)

// UserWorkloadQuery  工作量查询
func UserWorkloadQuery(startTime string, endTime string) []model.UserWorkloadInfo {
	var ProdAccept []model.UserProdAccept                      // 入库
	var DpProd []model.DepartmentCollar                        // 出库
	var RefProd []model.RefundProd                             // 退货
	UserWorkloadMap := make(map[string]model.UserWorkloadInfo) // 合并数据map
	var UserWorkload []model.UserWorkloadInfo                  // 合并后的数据切片
	// 入库
	clientDb.DB.Raw(clientDb.UserProdAcceptSql, startTime, endTime, startTime, endTime).Find(&ProdAccept)
	// 出库
	clientDb.DB.Raw(clientDb.UserProdDpcSql, startTime, endTime, startTime, endTime).Find(&DpProd)
	// 退货
	clientDb.DB.Raw(clientDb.UserProdRefundSql, startTime, endTime, startTime, endTime).Find(&RefProd)

	// 合并数据
	for i := 0; i < len(ProdAccept) || i < len(DpProd) || i < len(RefProd); i++ {
		if i < len(ProdAccept) {
			usm := UserWorkloadMap[ProdAccept[i].Name] // 找到对应Name的map并把value赋值给usm
			usm.Name = ProdAccept[i].Name
			usm.ProdAcBill = ProdAccept[i].ProdAcBill
			usm.ProdAcSpec = ProdAccept[i].ProdAcSpec
			usm.ProdAcTotal = ProdAccept[i].ProdAcTotal
			UserWorkloadMap[ProdAccept[i].Name] = usm
		}
		if i < len(DpProd) {
			usm := UserWorkloadMap[DpProd[i].Name]
			usm.Name = DpProd[i].Name
			usm.ProdDpBill = DpProd[i].ProdDpBill
			usm.ProdDpSpec = DpProd[i].ProdDpSpec
			usm.ProdDpTotal = DpProd[i].ProdDpTotal
			UserWorkloadMap[DpProd[i].Name] = usm
		}
		if i < len(RefProd) {
			usm := UserWorkloadMap[RefProd[i].Name]
			usm.Name = RefProd[i].Name
			usm.RefBill = RefProd[i].RefBill
			usm.RefSpec = RefProd[i].RefSpec
			usm.RefTotal = RefProd[i].RefTotal
			UserWorkloadMap[RefProd[i].Name] = usm
		}
	}
	// 将合并后的数据写入切片
	for _, v := range UserWorkloadMap {
		UserWorkload = append(UserWorkload, v)
	}
	return UserWorkload
}
