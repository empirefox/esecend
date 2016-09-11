//go:generate reform
package front

//reform:cc_vip_rebate_origin
type VipRebateOrigin struct {
	ID        uint  `reform:"id,pk"`
	UserID    uint  `reform:"user_id" json:"-"`
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

type VipRebateRequest struct {
	Type   string
	SubIDs []uint
}
