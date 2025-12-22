package model

type RefundNo struct {
	Yddh string `json:"yddh" gorm:"column:yddh"`
	Rkfs string `json:"rkfs" gorm:"column:rkfs"`
	Sczt string `json:"sczt" gorm:"column:sczt;default:''"`
	Scsm string `json:"scsm" gorm:"column:scsm;default:''"`
}

type GetRefund interface {
	GetRefundNo() error
}

type RetryRefund interface {
	RetryRefundToHis() error
}
