package model

type ReturnNo struct {
	Ckdh           string `json:"yddh" gorm:"column:ckdh"`
	Ckfs           string `json:"ckfs" gorm:"column:ckfs"`
	Sczt           string `json:"sczt" gorm:"column:sczt;default:''"`
	Scsm           string `json:"scsm" gorm:"column:scsm;default:''"`
	StoreHouseName string `json:"storeHouseName" gorm:"column:storeHouseName"` // 库房
	SupplierName   string `json:"supplierName" gorm:"column:supplierName"`
}

type ReturnSerializer interface {
	Serializer() interface{}
}

type ReturnFullSerializer struct {
	*ReturnNo
}

func (r *ReturnFullSerializer) ReturnSerializer() interface{} {
	return struct {
		Ckdh string `json:"yddh" gorm:"column:ckdh"`
		Ckfs string `json:"ckfs" gorm:"column:ckfs"`
		Sczt string `json:"sczt" gorm:"column:sczt;default:''"`
		Scsm string `json:"scsm" gorm:"column:scsm;default:''"`
	}{
		Ckdh: r.Ckdh,
		Ckfs: r.Ckfs,
		Sczt: r.Sczt,
		Scsm: r.Scsm,
	}
}
