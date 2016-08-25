//go:generate reform
package front

import "github.com/empirefox/reform"

type CapitalFlowType int

const (
	TCapitalFlowUnknown CapitalFlowType = iota
	TCapitalFlowPrepay
	TCapitalFlowPrepayBack
	TCapitalFlowTrade
	TCapitalFlowRefund
	TCapitalFlowWithdraw
	TCapitalFlowRecharge
)

type PointsType int

const (
	TPointsUnkonwn PointsType = iota
	TPointsReward
	TPointsPrepay
	TPointsPrepayBack
	TPointsTrade
	TPointsRefund
)

//reform:cc_member_account_log
type CapitalFlow struct {
	ID        uint            `reform:"id,pk"`
	UserID    uint            `reform:"user_id" json:"-"`
	CreatedAt int64           `reform:"create_time"`
	Type      CapitalFlowType `reform:"log_type"`
	Reason    string          `reform:"reason"`
	Amount    int             `reform:"amount"`
	Balance   int             `reform:"balance"`
	OrderID   uint            `reform:"order_id"`
}

//reform:cc_member_credit_log
type PointsItem struct {
	ID        uint       `reform:"id,pk"`
	UserID    uint       `reform:"user_id" json:"-"`
	CreatedAt int64      `reform:"create_time"`
	Type      PointsType `reform:"log_type"`
	Reason    string     `reform:"reason"`
	Amount    int        `reform:"amount"`
	Balance   int        `reform:"balance"`
	OrderID   uint       `reform:"order_id"`
}

type Wallet struct {
	CapitalFlows []reform.Struct // CapitalFlow
	PointsList   []reform.Struct // PointsItem
}

type WxPayArgs struct {
	AppId     string `json:"appId"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
}
