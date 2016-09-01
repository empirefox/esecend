//go:generate reform
package front

import "github.com/empirefox/reform"

type UserCashType int

const (
	TUserCashUnknown UserCashType = iota
	TUserCashPrepay
	TUserCashPrepayBack
	TUserCashTrade
	TUserCashRefund
	TUserCashWithdraw
	TUserCashRecharge
	TUserCashRebate
	TUserCashStoreRebate
)

//reform:cc_member_account_log
type UserCash struct {
	ID        uint         `reform:"id,pk"`
	UserID    uint         `reform:"user_id" json:"-"`
	CreatedAt int64        `reform:"create_time"`
	Type      UserCashType `reform:"log_type"`
	Remark    string       `reform:"remark"`
	Amount    int          `reform:"amount"`
	Balance   int          `reform:"balance"`
	OrderID   uint         `reform:"order_id"`
}

//reform:cc_user_cash_unfrozen
type UserCashUnfrozen struct {
	ID        uint  `reform:"id,pk"`
	FrozenID  uint  `reform:"frozen_id"`
	CreatedAt int64 `reform:"create_at"`
	Amount    uint  `reform:"amount"`
}

//reform:cc_user_cash_frozen
type UserCashFrozen struct {
	ID        uint   `reform:"id,pk"`
	UserID    uint   `reform:"user_id" json:"-"`
	CreatedAt int64  `reform:"create_at"`
	Remark    string `reform:"remark"`
	Total     uint   `reform:"total"`
	Stages    uint   `reform:"stages"`
	DoneAt    int64  `reform:"done_at"`
}

//reform:cc_member_credit_log
type PointsItem struct {
	ID        uint  `reform:"id,pk"`
	UserID    uint  `reform:"user_id" json:"-"`
	CreatedAt int64 `reform:"create_time"`
	Amount    int   `reform:"amount"`
	Balance   int   `reform:"balance"`
	OrderID   uint  `reform:"order_id"`
}

type Wallet struct {
	Cashes   []reform.Struct // UserCash
	Frozen   []reform.Struct // UserCashFrozen, exclude DoneAt
	Unfrozen []reform.Struct // UserCashUnfrozen
	Points   []reform.Struct // PointsItem
}

type WxPayArgs struct {
	AppId     string `json:"appId"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
}
