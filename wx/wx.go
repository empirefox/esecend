package wx

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"time"

	"gopkg.in/doug-martin/goqu.v3"

	"github.com/Sirupsen/logrus"
	"github.com/chanxuehong/util"
	"github.com/chanxuehong/wechat.v2/mch/core"
	"github.com/chanxuehong/wechat.v2/mch/pay"
	"github.com/dchest/uniuri"
	"github.com/empirefox/esecend/cerr"
	"github.com/empirefox/esecend/config"
	"github.com/empirefox/esecend/db-service"
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/esecend/l"
	"github.com/empirefox/esecend/lok"
	"github.com/empirefox/esecend/models"
)

var log = logrus.New()

type WxClient struct {
	*core.Client
	notifyUrl string
	wx        *config.Weixin
	dbs       *dbsrv.DbService
}

func NewWxClient(config *config.Config, dbs *dbsrv.DbService) (*WxClient, error) {
	weixin := &config.Weixin
	httpClient, err := core.NewTLSHttpClient(weixin.CertFile, weixin.KeyFile)
	if err != nil {
		return nil, err
	}
	return &WxClient{
		wx:        weixin,
		notifyUrl: config.Security.SecendOrigin + config.Security.PayNotifyPath,
		Client:    core.NewClient(weixin.AppId, weixin.MchId, weixin.ApiKey, httpClient),
		dbs:       dbs,
	}, nil
}

// only can be called by PrepayOrder
func (wc *WxClient) UnifiedOrder(tokUsr *models.User, order *front.Order, ip string, attach *models.UnifiedOrderAttach) (*front.WxPayArgs, error) {
	req := &pay.UnifiedOrderRequest{
		DeviceInfo:     "WEB",
		Body:           wc.wx.PayBody,
		OutTradeNo:     order.TrackingNumber(),
		TotalFee:       int64(order.PayAmount),
		SpbillCreateIP: ip,
		NotifyURL:      wc.wx.PayNotifyURL,
		TradeType:      "JSAPI",
		OpenId:         tokUsr.OpenId,
	}
	if attach != nil && attach.CashPaid != 0 && attach.PointsPaid != 0 {
		attachBs, _ := json.Marshal(attach)
		req.Attach = base64.URLEncoding.EncodeToString(attachBs)
	}

	res, err := pay.UnifiedOrder2(wc.Client, req)
	if err != nil {
		return nil, err
	}

	args := &front.WxPayArgs{
		AppId:     wc.wx.AppId,
		TimeStamp: strconv.FormatInt(time.Now().Unix(), 10),
		NonceStr:  uniuri.NewLen(32),
		Package:   "prepay_id=" + res.PrepayId, // 2hour
		SignType:  "MD5",
	}
	args.PaySign = core.JsapiSign(args.AppId, args.TimeStamp, args.NonceStr, args.Package, args.SignType, wc.wx.ApiKey)
	return args, nil
}

func (wc *WxClient) OnWxPayNotify(r io.Reader) interface{} {
	m, err := util.DecodeXMLToMap(r)
	if err != nil {
		return NewWxResponse("FAIL", "failed to parse request body")
	}
	if m["return_code"] != "SUCCESS" {
		return NewWxResponse(m["return_code"], m["return_msg"])
	}

	sign := core.Sign(m, wc.wx.ApiKey, md5.New)
	if sign != m["sign"] {
		return NewWxResponse("FAIL", "failed to validate md5")
	}

	var at int64
	var id uint
	_, err = fmt.Sscanf(m["out_trade_no"], "%d-%d", &at, &id)
	if err != nil {
		return NewWxResponse("FAIL", "failed to parse out_trade_no")
	}

	if m["result_code"] == "SUCCESS" {
		m["trade_state"] = "SUCCESS"
	} else {
		m["trade_state"] = "NOPAY"
	}

	if !lok.OrderLok.Lock(id) {
		return NewWxResponse("FAIL", "order is locked temporally")
	}
	defer lok.OrderLok.Unlock(id)

	var userId uint
	var cashLocked, pointsLocked bool
	defer func() {
		if cashLocked {
			lok.CashLok.Unlock(userId)
		}
		if pointsLocked {
			lok.PointsLok.Unlock(userId)
		}
	}()

	err = wc.dbs.InTx(func(dbs *dbsrv.DbService) error {
		order, err := dbs.GetBareOrder(nil, id)
		if err != nil {
			return err
		}
		userId, cashLocked, pointsLocked, err = wc.updateWxOrderSate(dbs, order, m, nil)
		return err
	})

	// must trans to WxReponse
	if err != nil {
		return NewWxResponse("FAIL", "failed to update trade state")
	}

	return NewWxResponse("SUCCESS", "")
}

func (wc *WxClient) updateWxOrderSate(
	dbs *dbsrv.DbService, order *front.Order, src map[string]string, tokUsr *models.User,
) (userId uint, cashLocked, pointsLocked bool, err error) {

	tradeState := front.TradeStateNameToValue[src["trade_state"]]
	tid := src["transaction_id"]
	if order.TradeState == tradeState && order.TransactionId == tid {
		// no need update
		return
	}

	totalFee64, _ := strconv.ParseUint(src["total_fee"], 10, 64)
	if totalFee64 == 0 {
		err = cerr.ParseWxTotalFeeFailed
		return
	}

	var attach models.UnifiedOrderAttach
	attachB64, hasAttach := src["attach"]
	if hasAttach {
		var attachBs []byte
		attachBs, err = base64.URLEncoding.DecodeString(attachB64)
		if err != nil {
			return
		}

		err = json.Unmarshal(attachBs, &attach)
		if err != nil {
			return
		}

		if attach.CashPaid != order.CashPaid || attach.PointsPaid != order.PointsPaid {
			err = cerr.InvalidPayAmount
			return
		}
		if attach.UserID != tokUsr.ID {
			err = cerr.InvalidUserID
			return
		}
		userId = attach.UserID
	}

	switch tradeState {
	case front.SUCCESS:
		if err = dbsrv.PermitOrderState(order, front.TOrderStatePaid); err != nil {
			log.WithFields(l.Locate(logrus.Fields{
				"OrderID": order.ID,
				"State":   order.State,
			})).Info("Got wxpay with SUCCESS")
			return
		}
		if hasAttach && order.State == front.TOrderStatePrepaid {
			now := time.Now().Unix()
			if attach.PreCashID != 0 {
				if cashLocked = lok.CashLok.Lock(attach.UserID); !cashLocked {
					err = cerr.CashTmpLocked
					return
				}

				ds := dbs.DS.Where(goqu.I(front.CapitalFlowTable.PK()).Eq(attach.PreCashID)).Where(goqu.I("$UserID").Eq(attach.UserID))
				_, err = dbs.GetDB().DsUpdateColumns(&front.CapitalFlow{
					Type:      front.TCapitalFlowTrade,
					CreatedAt: now,
				}, ds, "Type", "CreatedAt")
				if err != nil {
					return
				}
			}

			if attach.PrePointsID != 0 {
				if pointsLocked = lok.PointsLok.Lock(attach.UserID); !pointsLocked {
					err = cerr.CashTmpLocked
					return
				}

				ds := dbs.DS.Where(goqu.I(front.PointsItemTable.PK()).Eq(attach.PrePointsID)).Where(goqu.I("$UserID").Eq(attach.UserID))
				_, err = dbs.GetDB().DsUpdateColumns(&front.PointsItem{
					Type:      front.TPointsTrade,
					CreatedAt: now,
				}, ds, "Type", "CreatedAt")
				if err != nil {
					return
				}
			}
		}

		timeEnd, errTime := time.Parse("20060102150405", src["time_end"])
		if errTime != nil {
			timeEnd = time.Now()
		}

		data := front.Order{
			ID:            order.ID,
			WxPaid:        uint(totalFee64),
			TransactionId: tid,
			TradeState:    tradeState,
			State:         front.TOrderStatePaid,
			PaidAt:        timeEnd.Unix(),
		}
		err = dbs.GetDB().UpdateColumns(&data, "TransactionId", "TradeState", "State", "PaidAt")
		if err != nil {
			return
		}
		order.WxPaid = data.WxPaid
		order.TransactionId = data.TransactionId
		order.TradeState = data.TradeState
		order.State = data.State
		order.PaidAt = data.PaidAt

	case front.REFUND, front.USERPAYING, front.PAYERROR:
		err = dbs.GetDB().UpdateColumns(&front.Order{ID: order.ID, TradeState: tradeState}, "TradeState")
		if err == nil {
			order.TradeState = tradeState
		}

	case front.CLOSED:
		// must be an expired order, ignore
		// refound must be called with SUCCESS on closeorder
		// web must refresh!
		err = cerr.OrderClosed
	}

	return
}

func (wc *WxClient) UpdateWxOrderSate(tokUsr *models.User, dbs *dbsrv.DbService, order *front.Order) (cashLocked, pointsLocked bool, err error) {
	var res map[string]string
	res, err = wc.OrderQuery(order)
	if err != nil {
		return
	}

	_, cashLocked, pointsLocked, err = wc.updateWxOrderSate(dbs, order, res, tokUsr)
	return
}

func (wc *WxClient) OrderQuery(order *front.Order) (map[string]string, error) {
	req := map[string]string{
		"appid":        wc.wx.AppId,
		"mch_id":       wc.wx.MchId,
		"out_trade_no": order.TrackingNumber(),
		"nonce_str":    uniuri.NewLen(32),
		"notify_url":   wc.notifyUrl,
	}
	if order.TransactionId != "" {
		req["transaction_id"] = order.TransactionId
	}
	req["sign"] = core.Sign(req, wc.wx.ApiKey, md5.New)
	return pay.OrderQuery(wc.Client, req)
}

func (wc *WxClient) OrderClose(order *front.Order) (map[string]string, error) {
	req := map[string]string{
		"appid":        wc.wx.AppId,
		"mch_id":       wc.wx.MchId,
		"out_trade_no": order.TrackingNumber(),
		"nonce_str":    uniuri.NewLen(32),
	}
	req["sign"] = core.Sign(req, wc.wx.ApiKey, md5.New)
	return pay.CloseOrder(wc.Client, req)
}

func (wc *WxClient) OrderRefund(order *front.Order, opUserId string) (map[string]string, error) {
	req := map[string]string{
		"appid":         wc.wx.AppId,
		"mch_id":        wc.wx.MchId,
		"nonce_str":     uniuri.NewLen(32),
		"out_trade_no":  order.TrackingNumber(),
		"out_refund_no": order.TrackingNumber(),
		"total_fee":     strconv.Itoa(int(order.WxPaid)),
		"refund_fee":    strconv.Itoa(int(order.WxRefund)),
		"op_user_id":    opUserId,
	}
	req["sign"] = core.Sign(req, wc.wx.ApiKey, md5.New)
	return pay.Refund(wc.Client, req)
}

type WxResponse struct {
	XMLName    xml.Name  `xml:"xml"`
	ReturnCode CDATAText `xml:"return_code"`
	ReturnMsg  CDATAText `xml:"return_msg,omitempty"`
}

type CDATAText struct {
	Text string `xml:",cdata"`
}

func NewWxResponse(code, msg string) *WxResponse {
	return &WxResponse{
		ReturnCode: CDATAText{code},
		ReturnMsg:  CDATAText{msg},
	}
}
