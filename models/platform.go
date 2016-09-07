//go:generate reform
package models

import "github.com/empirefox/esecend/front"

//reform:cc_platform_cash
type PlatformCash struct {
	ID         uint               `reform:"id,pk"`
	CreatedAt  int64              `reform:"create_at"`
	Type       front.UserCashType `reform:"typ"`
	OrderID    uint               `reform:"order_id"`
	OrderTotal uint               `reform:"total"`
	Amount     uint               `reform:"amount"`
	Balance    uint               `reform:"balance"`
	Remark     string             `reform:"remark"`
}
