//go:generate jsonconst -w=c -type=CodedError ./cerr
//go:generate jsonconst -w=u -type=TradeState,CapitalFlowType,PointsType ./front
//go:generate mapconst -type=TradeState,CapitalFlowType,PointsType ./front
package esecend
