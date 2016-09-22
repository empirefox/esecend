//go:generate reform
package front

//reform:cc_member_address
type Address struct {
	ID       uint   `reform:"id,pk"`
	UserID   uint   `reform:"uid" json:"-"`
	Contact  string `reform:"contactor"`
	Phone    string `reform:"tel_num"`
	Province string `reform:"province"`
	City     string `reform:"city"`
	District string `reform:"district"`
	House    string `reform:"address"`
	Pos      int64  `reform:"pos"`
}
