//go:generate reform
package front

import "github.com/empirefox/reform"

type VipRebateType int

const (
	TVipRebateUnknown VipRebateType = iota
	TVipRebateRebate
	TVipRebateReward
)

//reform:cc_member
type VipIntro struct {
	ID           uint   `reform:"id,pk"`
	CreatedAt    int64  `reform:"create_date"`
	HeadImageURL string `reform:"avatar"`

	Nickname     string `reform:"name"`
	Sex          int    `reform:"sex"`
	City         string `reform:"city"`
	Province     string `reform:"province"`
	Birthday     int64  `reform:"birthday"`
	CarInsurance string `reform:"car_insurance"`
	InsuranceFee uint   `reform:"insurance_fee"`
	CarIntro     string `reform:"car_intro"`
	Hobby        string `reform:"hobby"`
	Career       string `reform:"career"`
	Demand       string `reform:"demand"`
	Intro        string `reform:"intro"`
}

//reform:cc_vip_rebate_origin
type VipRebateOrigin struct {
	ID        uint  `reform:"id,pk"`
	UserID    uint  `reform:"user_id"`
	CreatedAt int64 `reform:"created_at"`
	NotBefore int64 `reform:"nbf"`
	ExpiresAt int64 `reform:"exp"`
	OrderID   uint  `reform:"order_id"`
	ItemID    uint  `reform:"item_id"`
	Amount    uint  `reform:"amount"`
	Balance   uint  `reform:"balance"`
	User1     uint  `reform:"user1"`
	User1Used bool  `reform:"user1_used"`
}

func (vip *VipRebateOrigin) Valid(now int64) bool {
	return vip.NotBefore <= now && now < vip.ExpiresAt
}

type VipRebatePayload struct {
	Type   VipRebateType
	SubIDs []uint
}

//reform:cc_member
type VipName struct {
	ID       uint   `reform:"id,pk"`
	Nickname string `reform:"name"`
}

type QualificationsResponse struct {
	Items []reform.Struct // VipRebateOrigin
	Names []reform.Struct // VipName
}
