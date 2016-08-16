package cerr

import "fmt"

type CodedError int

func (ce CodedError) Error() string {
	return fmt.Sprintf("Error code: %d", ce)
}

// all are non-StatusForbidden
const (
	Error CodedError = iota
	SonyFlakeTimeout
	InvalidUrlParam
	InvalidPostBody
	InvalidPhoneFormat
	RebindSamePhone
	PhoneOccupied
	UserNotFound
	DbFailed
	CaptchaRejected
	UpdateWxOrderStateFailed
	SystemModeNotAllowed
	InvalidRefreshToken
	NoRefreshToken
	NoNeedRefreshToken
	NoAccessToken
	InvalidAccessToken
	GenCaptchaFailed
	RemoteHTTPFailed
	InvalidTokenSubject
	InvalidTokenExpires
	InvalidSignAlg
	InvalidClaimId
	Forbidden

	InvalidProductId
	InvalidSkuStock
	InvalidSkuId
	InvalidAttrId
	InvalidAttrLen
	InvalidGroupbuyId
	InvalidCheckoutTotal
	InvalidCheckoutFreight
	InvalidPaykey
	NotEnoughMoney
	NotEnoughPoints
	InvalidPayAmount
	NotNopayState
	NoWayToPaidState
	NoWayToTargetState
	OrderClosed
	OrderCloseNeeded
	OrderTmpLocked
	OrderCompleteTimeout
	OrderEvalTimeout
	CashTmpLocked
	PointsTmpLocked
	OrderItemNotFound
	NotPrepayOrder
	WxPayNotCompleted
	WxRefundNotCompleted
	WxOrderNotExist
	WxOrderAlreadyClosed
	WxOrderCloseFailed
	WxOrderAlreadyPaid
	WxSystemFailed
	InvalidCashPrepaid
	InvalidPointsPrepaid
	InvalidPrepayPayload
	ApiImplementFailed
	WxUserAbnormal
)
