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

	"github.com/Sirupsen/logrus"
	"github.com/chanxuehong/util"
	"github.com/chanxuehong/wechat.v2/mch/core"
	"github.com/chanxuehong/wechat.v2/mch/mmpaymkttransfers/promotion"
	"github.com/chanxuehong/wechat.v2/mch/pay"
	"github.com/dchest/uniuri"
	"github.com/empirefox/esecend/config"
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/esecend/models"
)

var log = logrus.New()

type WxClient struct {
	*core.Client
	notifyUrl string
	wx        *config.Weixin
}

func NewWxClient(config *config.Config) (*WxClient, error) {
	weixin := &config.Weixin
	httpClient, err := core.NewTLSHttpClient(weixin.CertFile, weixin.KeyFile)
	if err != nil {
		return nil, err
	}
	return &WxClient{
		wx:        weixin,
		notifyUrl: config.Security.SecendOrigin + config.Security.PayNotifyPath,
		Client:    core.NewClient(weixin.AppId, weixin.MchId, weixin.ApiKey, httpClient),
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

func (wc *WxClient) OnWxPayNotify(r io.Reader) (*WxResponse, map[string]string) {
	m, err := util.DecodeXMLToMap(r)
	if err != nil {
		return NewWxResponse("FAIL", "failed to parse request body"), nil
	}
	if m["return_code"] != "SUCCESS" {
		return NewWxResponse(m["return_code"], m["return_msg"]), nil
	}

	sign := core.Sign(m, wc.wx.ApiKey, md5.New)
	if sign != m["sign"] {
		return NewWxResponse("FAIL", "failed to validate md5"), nil
	}

	var at int64
	var id uint
	_, err = fmt.Sscanf(m["out_trade_no"], "%d-%d", &at, &id)
	if err != nil {
		return NewWxResponse("FAIL", "failed to parse out_trade_no"), nil
	}

	if m["result_code"] == "SUCCESS" {
		m["trade_state"] = "SUCCESS"
	} else {
		m["trade_state"] = "NOPAY"
	}

	return nil, m
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

type TransfersArgs struct {
	TradeNo string
	OpenID  string
	Amount  uint
	Desc    string
	Ip      string
}

func (wc *WxClient) Transfers(args *TransfersArgs) (map[string]string, error) {
	req := map[string]string{
		"mch_appid":        wc.wx.AppId,
		"mchid":            wc.wx.MchId,
		"nonce_str":        uniuri.NewLen(32),
		"partner_trade_no": args.TradeNo,
		"openid":           args.OpenID,
		"check_name":       "NO_CHECK",
		"amount":           strconv.FormatUint(uint64(args.Amount), 10),
		"desc":             args.Desc,
		"spbill_create_ip": args.Ip,
	}
	req["sign"] = core.Sign(req, wc.wx.ApiKey, md5.New)
	return promotion.Transfers(wc.Client, req)
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
