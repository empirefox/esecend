//go:generate reform
package models

import (
	"fmt"

	"github.com/empirefox/esecend/front"
)

//reform:cc_platform_cash
type PlatformCash struct {
	ID         uint               `reform:"id,pk"`
	CreatedAt  int64              `reform:"created_at"`
	Type       front.UserCashType `reform:"typ"`
	OrderID    uint               `reform:"order_id"`
	OrderTotal uint               `reform:"total"`
	Amount     uint               `reform:"amount"`
	Balance    uint               `reform:"balance"`
	Remark     string             `reform:"remark"`
}

//reform:cc_cash_withdraw
type CashWithdraw struct {
	ID        uint   `reform:"id,pk"`
	CreatedAt int64  `reform:"created_at"`
	UserID    uint   `reform:"user_id"`
	Amount    uint   `reform:"amount"`
	Desc      string `reform:"desc"`
	TradeNo   string `reform:"trade_no"`
	WxNo      string `reform:"wx_no"`
	WxTime    int64  `reform:"wx_time"`
}

func (c *CashWithdraw) TrackingNumber() string {
	return fmt.Sprintf("%d-%d", c.CreatedAt, c.ID)
}
