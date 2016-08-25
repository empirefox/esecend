package cerr

import "fmt"

type CodedError int

func (ce CodedError) Error() string {
	return fmt.Sprintf("Error code: %d", ce)
}

// all are non-StatusForbidden
const (
	Error CodedError = iota
	Unauthorized
	SonyFlakeTimeout
	InvalidUrlParam
	InvalidPostBody
	InvalidPhoneFormat
	RebindSamePhone
	PhoneOccupied
	PhoneBindRequired
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
	InvalidUserID
	Forbidden
	RetrySmsFailed
	SendSmsError
	SendSmsFailed
	SmsVerifyFailed

	InvalidProductId
	InvalidSkuStock
	InvalidSkuId
	InvalidAttrId
	InvalidAttrLen
	InvalidGroupbuyId
	InvalidCheckoutTotal
	InvalidCheckoutFreight
	InvalidPaykey
	PaykeyNeedBeSet
	NotEnoughMoney
	NotEnoughPoints
	InvalidPayAmount
	NotNopayState
	NoWayToPaidState
	NoWayToTargetState
	NoPermToState
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
	ParseWxTotalFeeFailed
)
