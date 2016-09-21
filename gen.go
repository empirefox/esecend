//go:generate jsonconst -w=c -type=CodedError ./cerr
//go:generate jsonconst -w=u -type=TradeState,UserCashType,OrderState,BillboardType ./front
//go:generate mapconst -type=TradeState,UserCashType,OrderState,BillboardType ./front
package esecend
