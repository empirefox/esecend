//go:generate reform
package front

//reform:cc_store
type Store struct {
	ID   uint   `reform:"id,pk"`
	Name string `reform:"name"`
}

//reform:cc_store_cash
type StoreCash struct {
	ID        uint   `reform:"id,pk"`
	StoreID   uint   `reform:"store_id" json:"-"`
	CreatedAt int64  `reform:"create_at"`
	OrderID   uint   `reform:"order_id"`
	Amount    uint   `reform:"amount"`
	Balance   uint   `reform:"balance"`
	Remark    string `reform:"remark"`
}
