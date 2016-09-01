//go:generate jsonconst -w=c -type=CodedError ./cerr
//go:generate jsonconst -w=u -type=TradeState,UserCashType ./front
//go:generate mapconst -type=TradeState,UserCashType ./front
package esecend
