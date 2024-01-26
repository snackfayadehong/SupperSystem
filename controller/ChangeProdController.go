package controller

import (
	clientDb "WorkloadQuery/db"
	"WorkloadQuery/model"
	"fmt"
)

/*
	根据老物资采购系统传过来的字典信息变更
*/

// ChangeInfoElement 字典信息入参
type ChangeInfoElement struct {
	Code         string  `json:"Code"`                   // 院内代码
	HospitalName *string `json:"HospitalName,omitempty"` // 院内产品名称
	HospitalSpec *string `json:"HospitalSpec,omitempty"` // 院内规格
	YGCGID       *string `json:"YGCGID,omitempty"`       // 网采平台产品ID
	TradeCode    *string `json:"TradeCode,omitempty"`    // 商品代码，商品代码或挂网流水号（平台供货商填写交易编码的内容）
	MedicareCode *string `json:"MedicareCode,omitempty"` // 医保代码
	CategoryCode *string `json:"CategoryCode,omitempty"` // 18类分类代码，匹配104分类三级目录,如果匹配成功不修改，否则修改为三级目录
	SysCode      *string `json:"SysCode,omitempty"`      // 系统编码
	SysID        *string `json:"SysID,omitempty"`        // 系统编号
}
type RequestInfo struct {
	C []ChangeInfoElement
}

/*
ChangeProductInfo
更改产品基本信息
*/
func (i *RequestInfo) ChangeProductInfo(prod *[]model.ProductInfo) {
	// 查找入参和物资系统返回查询结果中相同的记录 时间复杂度O(M+N)
	// 使用Map存储prod切片中的Code
	pMap := make(map[string]int, len(*prod))
	for index, v := range i.C {
		pMap[v.Code] = index
	}
	// 在入参中找到Code相同的
	for _, item := range i.C {
		if pIndex, ok := pMap[item.Code]; ok {
			// 修改字典信息的业务逻辑
			// 1. 检查104分类；
			// 入参中104分类为第3级,物资系统为第4级,查询物资系统第4级代码对应的第3级与入参是否相同
			// 如果不相同则修改物资为第3级,相同则不更新
			if *item.CategoryCode != (*prod)[pIndex].ParentCusCategoryCode {
				// 修改为第3级
			}
		}
	}
}

/*
GetProductInfo
获取物资产品字典信息,返回不重复的字典信息
string返回记录哪些字典信息是重复的
*/
func (i *RequestInfo) GetProductInfo(Where []string) (*[]model.ProductInfo, string) {
	var prod *[]model.ProductInfo         // 原始记录
	var NoRepeatProd []model.ProductInfo  // 保留不重复的记录
	var msg string                        // 返回重复记录信息
	var repeatMap = make(map[string]bool) // 重复记录
	clientDb.DB.Raw(clientDb.QueryProd, Where).Find(&prod)
	// 如果只有一行记录直接返回
	if len(*prod) == 1 {
		return prod, msg
	}
	// 检查 查询结果中同一院内编码是否存在多条记录,且非停用或停供产品
	seen := make(map[string]bool)
	for _, el := range *prod {
		// 非停用或者停供且不重复的记录添加到map,和切片中
		if !seen[el.Code] && el.PurState == 0 && el.IsVoid == 0 {
			seen[el.Code] = true
			NoRepeatProd = append(NoRepeatProd, el)
			continue
		}
		// 记录哪些是重复的,只记录一次
		if !repeatMap[el.Code] {
			repeatMap[el.Code] = true
		}
	}
	for key := range repeatMap {
		msg += fmt.Sprintf("%s有重复字典记录或供货关系异常;", key)
	}
	return &NoRepeatProd, msg
}
