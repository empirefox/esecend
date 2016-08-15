package models

import (
	"encoding/json"
	"strconv"

	"github.com/empirefox/esecend/front"
)

type WxGoodsDetail struct {
	ID    string `json:"goods_id"`   // Product.ID
	Name  string `json:"goods_name"` // Product.Name
	Num   uint   `json:"goods_num"`
	Price uint   `json:"price"`
}

type WxOrderDetail struct {
	Goods []WxGoodsDetail `json:"goods_detail"`
}

func MarshalWxOrderDetail(items []*front.OrderItem) ([]byte, error) {

	var goods []WxGoodsDetail
	for i := range items {
		goods = append(goods, WxGoodsDetail{
			ID:    strconv.FormatUint(uint64(items[i].ID), 10),
			Name:  items[i].Name,
			Num:   items[i].Quantity,
			Price: items[i].Price,
		})
	}

	return json.Marshal(&WxOrderDetail{goods})
}

type UnifiedOrderAttach struct {
	UserID      uint
	CashPaid    uint `,omitempty`
	PointsPaid  uint `,omitempty`
	PreCashID   uint `,omitempty`
	PrePointsID uint `,omitempty`
}
