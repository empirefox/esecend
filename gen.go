//go:generate jsonconst -w=c -type=CodedError ./cerr
//go:generate jsonconst -w=u -type=TradeState,UserCashType,OrderState ./front
//go:generate mapconst -type=TradeState,UserCashType,OrderState ./front
package esecend
