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
	TUserCashReward
	TUserCashRebate
	TUserCashStoreRebate
)

//reform:cc_member_account_log
type UserCash struct {
	ID        uint         `reform:"id,pk"`
	UserID    uint         `reform:"user_id" json:"-"`
	OrderID   uint         `reform:"order_id"`
	CreatedAt int64        `reform:"create_time"`
	Amount    int          `reform:"amount"`
	Remark    string       `reform:"remark"`
	Type      UserCashType `reform:"log_type"`
	Balance   int          `reform:"balance"`
}

//reform:cc_user_cash_frozen
type UserCashFrozen struct {
	ID        uint         `reform:"id,pk"`
	UserID    uint         `reform:"user_id" json:"-"`
	OrderID   uint         `reform:"order_id"`
	CreatedAt int64        `reform:"created_at"`
	Type      UserCashType `reform:"typ"`
	Amount    uint         `reform:"amount"`
	Remark    string       `reform:"remark"`
	ThawedAt  int64        `reform:"thawed_at"`
}

//reform:cc_user_cash_rebate_item
type UserCashRebateItem struct {
	ID        uint  `reform:"id,pk"`
	RebateID  uint  `reform:"rebate_id"`
	CreatedAt int64 `reform:"created_at"`
	Amount    uint  `reform:"amount"`
}

//reform:cc_user_cash_rebate
type UserCashRebate struct {
	ID        uint         `reform:"id,pk"`
	UserID    uint         `reform:"user_id" json:"-"`
	OrderID1  uint         `reform:"order_id1"`
	OrderID2  uint         `reform:"order_id2"`
	CreatedAt int64        `reform:"created_at"`
	Type      UserCashType `reform:"typ"`
	Amount    uint         `reform:"amount"`
	Remark    string       `reform:"remark"`
	Stages    uint         `reform:"stages"`
	DoneAt    int64        `reform:"done_at"`
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
	Cashes      []reform.Struct // UserCash
	Frozen      []reform.Struct // UserCashFrozen, exclude ThawedAt
	Rebates     []reform.Struct // UserCashRebate, exclude DoneAt
	RebateItems []reform.Struct // UserCashRebateItem
	Points      []reform.Struct // PointsItem
}

type WxPayArgs struct {
	AppId     string `json:"appId"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
}
