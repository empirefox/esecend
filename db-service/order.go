package dbsrv

import (
	"strconv"
	"strings"
	"time"

	"gopkg.in/doug-martin/goqu.v3"

	"github.com/Sirupsen/logrus"
	"github.com/cznic/sortutil"
	"github.com/dchest/uniuri"
	"github.com/empirefox/esecend/admin"
	"github.com/empirefox/esecend/cerr"
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/esecend/l"
	"github.com/empirefox/esecend/models"
	"github.com/empirefox/reform"
	"github.com/golang/glog"
)

func (dbs *DbService) GetOrderItems(o *front.Order) ([]*front.OrderItem, error) {
	if o.Items == nil {
		items, err := dbs.GetDB().FindAllFrom(front.OrderItemTable, "$OrderID", o.ID)
		if err != nil {
			return nil, err
		}
		for _, item := range items {
			o.Items = append(o.Items, item.(*front.OrderItem))
		}
	}
	return o.Items, nil
}

// only abc or points
func (dbs *DbService) CheckoutOrderOne(
	tokUsr *models.User, payload *front.CheckoutOnePayload,
) (*front.Order, error) {

	if payload.SkuID == 0 {
		return nil, cerr.InvalidSkuId
	}

	tx, err := dbs.Tx()
	if err != nil {
		return nil, err
	}
	defer tx.RollbackIfNeeded()
	db := tx.GetDB()

	var sku front.Sku
	if err = db.FindByPrimaryKeyTo(&sku, payload.SkuID); err != nil {
		return nil, err
	}
	if sku.Stock < 1 {
		return nil, cerr.InvalidSkuStock
	}

	var product front.Product
	if err = db.FindByPrimaryKeyTo(&product, sku.ProductID); err != nil {
		return nil, err
	}

	if product.Vpn != front.TVpnVip && product.Vpn != front.TVpnPoints {
		return nil, cerr.OnlyAbcOrPoints
	}

	// Attrs: find all inter table data
	skuToAttrs, err := db.FindAllFrom(front.ProductAttrIdTable, "$SkuID", sku.ID)
	if err != nil {
		return nil, err
	}
	if len(skuToAttrs) == 0 {
		return nil, cerr.InvalidSkuId
	}

	// Attrs: query attrs
	var attrIds []interface{}
	for _, s2a := range skuToAttrs {
		attrIds = append(attrIds, s2a.(*front.ProductAttrId).AttrID)
	}
	if len(attrIds) == 0 {
		return nil, cerr.InvalidAttrId
	}
	attrs, err := db.FindAllFromPK(front.ProductAttrTable, attrIds...)
	if err != nil {
		return nil, err
	}
	if len(attrs) != len(attrIds) {
		return nil, cerr.InvalidAttrId
	}
	if len(attrs) != len(payload.Attrs) {
		return nil, cerr.InvalidAttrLen
	}

	attrMap := make(map[uint]*front.ProductAttr)
	var aIds []uint
	for _, attri := range attrs {
		attr := attri.(*front.ProductAttr)
		attrMap[attr.ID] = attr
		aIds = append(aIds, attr.ID)
	}

	// Attrs: check equal and load values
	var attrsCopy []uint
	attrsCopy = payload.Attrs[:]
	sortutil.UintSlice(attrsCopy).Sort()
	sortutil.UintSlice(aIds).Sort()

	for i, id := range aIds {
		if attrsCopy[i] != id {
			return nil, cerr.InvalidAttrId
		}
	}

	// load values
	var attrSnapshot []string
	for _, attrId := range payload.Attrs {
		attrSnapshot = append(attrSnapshot, attrMap[attrId].Value)
	}

	now := time.Now().Unix()

	// save order
	order := front.Order{
		Remark: payload.Remark,
		UserID: tokUsr.ID,

		// study http://help.vipshop.com/themelist.php?type=detail&id=330
		State:     front.TOrderStateNopay,
		CreatedAt: now,

		// OrderAddress
		Contact:        payload.Contact,
		Phone:          payload.Phone,
		DeliverAddress: payload.DeliverAddress,
		User1:          tokUsr.User1,
	}

	if product.Vpn == front.TVpnPoints {
		order.PayPoints = sku.SalePrice
	} else {
		// Invoice
		order.InvoiceTo = payload.InvoiceTo
		order.InvoiceToCom = payload.InvoiceToCom
		order.PayAmount = sku.SalePrice
	}

	if err = db.Insert(&order); err != nil {
		return nil, err
	}

	// store1
	var store1 uint
	if product.StoreID != 0 {
		var store front.Store
		err = db.FindByPrimaryKeyTo(&store, product.StoreID)
		if err != nil && err != reform.ErrNoRows {
			return nil, err
		}
		err = nil
		store1 = store.User1
	}

	item := front.OrderItem{
		ProductID: sku.ProductID,
		SkuID:     sku.ID,
		Vpn:       product.Vpn,
		Quantity:  1,
		Price:     sku.SalePrice,
		CreatedAt: now,
		Img:       sku.Img,
		UserID:    tokUsr.ID,
		OrderID:   order.ID,
		StoreID:   product.StoreID,
		Name:      product.Name,
		Attrs:     strings.Join(attrSnapshot, " "),
		Store1:    store1,
	}
	if item.Img == "" {
		item.Img = product.Img
	}

	// update sku stock
	sku.Stock--
	if err = db.UpdateColumns(&sku, "Stock"); err != nil {
		return nil, err
	}
	if err = db.Insert(&item); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	order.Items = []*front.OrderItem{&item}
	return &order, nil
}

// no ABC or points
func (dbs *DbService) CheckoutOrder(tokUsr *models.User, payload *front.CheckoutPayload) (*front.Order, error) {
	// prepare skuIds, groupbuyIds to query
	var skuIds []interface{}
	var groupbuyIds []interface{}
	skuidToPayloadItem := make(map[uint]*front.CheckoutPayloadItem)
	for i := range payload.Items {
		skuid := payload.Items[i].SkuID
		skuIds = append(skuIds, skuid)
		if gbid := payload.Items[i].GroupBuyID; gbid != 0 {
			groupbuyIds = append(groupbuyIds, gbid)
		}
		skuidToPayloadItem[skuid] = &payload.Items[i]
	}
	if len(skuIds) == 0 {
		return nil, cerr.InvalidSkuId
	}

	tx, err := dbs.Tx()
	if err != nil {
		return nil, err
	}
	defer tx.RollbackIfNeeded()
	db := tx.GetDB()

	// query skus, groupbuys
	skus, err := db.FindAllFromPK(front.SkuTable, skuIds...)
	if err != nil {
		return nil, err
	}
	lskus := len(skus)
	if lskus != len(skuIds) {
		return nil, cerr.InvalidSkuId
	}

	gbMap := make(map[uint]*front.GroupBuyItem)
	if len(groupbuyIds) != 0 {
		gbs, err := db.FindAllFromPK(front.GroupBuyItemTable, groupbuyIds...)
		if err != nil {
			return nil, err
		}
		if len(gbs) != len(groupbuyIds) {
			return nil, cerr.InvalidGroupbuyId
		}
		for _, gbi := range gbs {
			gb := gbi.(*front.GroupBuyItem)
			gbMap[gb.SkuID] = gb
		}

	}

	// compute part of order item, prepare for query products
	var total uint
	var price uint
	var freight uint
	var items []*front.OrderItem
	productIdMap := make(map[uint]bool)
	skuMap := make(map[uint]*front.Sku)
	now := time.Now().Unix()
	for _, skui := range skus {
		sku := skui.(*front.Sku)
		num := skuidToPayloadItem[sku.ID].Quantity
		if num == 0 || num > sku.Stock {
			return nil, cerr.InvalidSkuStock
		}

		price = sku.SalePrice
		if gb, ok := gbMap[sku.ID]; ok {
			price = gb.Price
		}
		total += price * num
		if sku.Freight > freight {
			freight = sku.Freight
		}
		productIdMap[sku.ProductID] = true
		skuMap[sku.ID] = sku
		items = append(items, &front.OrderItem{
			ProductID:  sku.ProductID,
			SkuID:      sku.ID,
			Quantity:   num,
			Price:      price,
			CreatedAt:  now,
			Img:        sku.Img,
			DeliverFee: sku.Freight,
		})
	}

	if len(productIdMap) == 0 {
		return nil, cerr.InvalidProductId
	}
	// query products
	var productIds []interface{}
	for pid := range productIdMap {
		productIds = append(productIds, pid)
	}
	products, err := db.FindAllFromPK(front.ProductTable, productIds...)
	if err != nil {
		return nil, err
	}
	if len(products) != len(productIds) {
		return nil, cerr.InvalidProductId
	}

	productMap := make(map[uint]*front.Product)
	for _, producti := range products {
		product := producti.(*front.Product)
		productMap[product.ID] = product
		if product.Vpn != front.TVpnNormal {
			return nil, cerr.NoAbcOrPoints
		}
	}

	// check order money
	if total >= dbs.config.Order.FreeDeliverLine {
		freight = 0
	}
	total += freight

	if freight != payload.DeliverFee {
		return nil, cerr.InvalidCheckoutFreight
	}

	if total != payload.Total {
		return nil, cerr.InvalidCheckoutTotal
	}

	// Attrs: find all inter table data
	skuToAttrs, err := db.FindAllFrom(front.ProductAttrIdTable, "$SkuID", skuIds...)
	if err != nil {
		return nil, err
	}
	if len(skuToAttrs) == 0 {
		return nil, cerr.InvalidSkuId
	}
	// Attrs: prepare attr ids
	attrIdMap := make(map[uint]bool)
	for _, skuToAttri := range skuToAttrs {
		skuToAttr := skuToAttri.(*front.ProductAttrId)
		attrIdMap[skuToAttr.AttrID] = true

		// Attrs: save ids to payload item
		payloadItem := skuidToPayloadItem[skuToAttr.SkuID]
		payloadItem.AttrIds = append(payloadItem.AttrIds, skuToAttr.AttrID)
	}
	if len(attrIdMap) == 0 {
		return nil, cerr.InvalidAttrId
	}

	// Attrs: query attrs
	var attrIds []interface{}
	for aid := range attrIdMap {
		attrIds = append(attrIds, aid)
	}
	attrs, err := db.FindAllFromPK(front.ProductAttrTable, attrIds...)
	if err != nil {
		return nil, err
	}
	if len(attrs) != len(attrIds) {
		return nil, cerr.InvalidAttrId
	}
	attrMap := make(map[uint]*front.ProductAttr)
	for _, attri := range attrs {
		attr := attri.(*front.ProductAttr)
		attrMap[attr.ID] = attr
	}

	// Attrs: check equal and load values
	for _, payloadItem := range skuidToPayloadItem {
		// check equal
		if len(payloadItem.Attrs) != len(payloadItem.AttrIds) {
			return nil, cerr.InvalidAttrLen
		}
		var attrsCopy []uint
		attrsCopy = append(attrsCopy, payloadItem.Attrs...)
		sortutil.UintSlice(attrsCopy).Sort()

		sortutil.UintSlice(payloadItem.AttrIds).Sort()

		for i, id := range payloadItem.AttrIds {
			if attrsCopy[i] != id {
				return nil, cerr.InvalidAttrId
			}
		}

		// load values
		var values []string
		for _, attrId := range payloadItem.Attrs {
			values = append(values, attrMap[attrId].Value)
		}
		payloadItem.AttrValues = strings.Join(values, " ")
	}

	// save order
	order := front.Order{
		PayAmount: total,
		Remark:    payload.Remark,
		UserID:    tokUsr.ID,
		User1:     tokUsr.User1,

		IsDeliverPay: payload.IsDeliverPay,
		DeliverFee:   payload.DeliverFee,

		// study http://help.vipshop.com/themelist.php?type=detail&id=330
		State:     front.TOrderStateNopay,
		CreatedAt: now,

		// Invoice
		InvoiceTo:    payload.InvoiceTo,
		InvoiceToCom: payload.InvoiceToCom,

		// OrderAddress
		Contact:        payload.Contact,
		Phone:          payload.Phone,
		DeliverAddress: payload.DeliverAddress,
	}

	if err = db.Insert(&order); err != nil {
		return nil, err
	}

	store1Map := make(map[uint]uint)
	// complete order items after order saved
	for _, item := range items {
		product := productMap[item.ProductID]
		item.UserID = tokUsr.ID
		item.OrderID = order.ID
		item.StoreID = product.StoreID
		item.Name = product.Name
		if item.Img == "" {
			item.Img = product.Img
		}
		item.Attrs = skuidToPayloadItem[item.SkuID].AttrValues

		// store1
		if product.StoreID != 0 {
			if store1, ok := store1Map[product.StoreID]; ok {
				item.Store1 = store1
			} else {
				var store front.Store
				err = db.FindByPrimaryKeyTo(&store, product.StoreID)
				if err != nil && err != reform.ErrNoRows {
					return nil, err
				}
				err = nil
				store1Map[product.StoreID] = store.User1
				item.Store1 = store.User1
			}
		}

		// update sku stock
		sku := skuMap[item.SkuID]
		sku.Stock -= item.Quantity
		if err = db.UpdateColumns(sku, "Stock"); err != nil {
			return nil, err
		}
		if err = db.Insert(item); err != nil {
			return nil, err
		}
	}
	order.Items = items

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &order, nil
}

func (dbs *DbService) PrepayOrder(tokUsr *models.User, orderId uint, ip *string) (o *front.Order, args *front.WxPayArgs, err error) {
	db := dbs.GetDB()

	var order front.Order
	ds := dbs.DS.Where(goqu.I(front.OrderTable.PK()).Eq(orderId)).Where(goqu.I("$UserID").Eq(tokUsr.ID))
	if err = db.DsSelectOneTo(&order, ds); err != nil {
		return
	}

	// exclude points pay
	if order.PayPoints != 0 {
		err = cerr.InvalidPayType
		return
	}

	now := time.Now().Unix()

	// prepay after prepaid
	if order.State == front.TOrderStatePrepaid && order.PrepaidAt != 0 {
		glog.Errorln("order.State == front.TOrderStatePrepaid")
		durPrepaid := now - order.PrepaidAt
		if durPrepaid < 290 {
			o = &order
			args = dbs.wc.NewWxPayArgs(&order.WxPrepayID)
			return
		}
		if durPrepaid < 300 {
			err = cerr.WxOrderCloseIn5Min
			return
		}
		// close transaction_id first
		var res map[string]string
		res, err = dbs.wc.OrderClose(&order)
		glog.Errorln(err)
		if err != nil {
			return
		}

		if res["result_code"] != "SUCCESS" {
			// colse FAIL
			log.WithFields(l.Locate(logrus.Fields{
				"out_trade_no": order.TrackingNumber(),
				"err_code":     res["err_code"],
			})).Info("Failed to close")

			switch res["err_code"] {
			case "ORDERPAID":
				err = cerr.WxOrderAlreadyPaid
			case "SYSTEMERROR":
				err = cerr.WxSystemFailed
			case "ORDERNOTEXIST", "ORDERCLOSED": // pass
			default:
				err = cerr.ApiImplementFailed
			}
			if err != nil {
				return
			}
		}
	}

	if err = PermitOrderState(&order, front.TOrderStatePrepaid); err != nil {
		err = cerr.NotNopayState
		return
	}

	order.PrepaidAt = now
	order.WxTradeNo = uniuri.NewLen(16)
	var preid *string
	preid, args, err = dbs.wc.UnifiedOrder(tokUsr, &order, ip)
	glog.Errorln(err)
	if err != nil {
		return
	}

	order.WxPrepayID = *preid
	order.State = front.TOrderStatePrepaid
	order.WxTransactionId = ""
	order.WxTradeState = front.UNKNOWN

	err = db.UpdateColumns(&order, "WxTradeNo", "State", "PrepaidAt", "WxTransactionId", "WxTradeState", "WxPrepayID")
	if err == nil {
		o = &order
	}
	return
}

// no wx pay, need paykey
func (dbs *DbService) PayOrder(tokUsr *models.User, payload *front.OrderPayPayload) (o *front.Order, err error) {
	db := dbs.GetDB()

	if err = db.Reload(tokUsr); err != nil {
		return
	}

	if tokUsr.Paykey == nil {
		err = cerr.PaykeyNeedBeSet
		return
	}

	if err = models.ComparePaykey(*tokUsr.Paykey, []byte(payload.Key)); err != nil {
		err = cerr.InvalidPaykey
		return
	}

	var order front.Order
	ds := dbs.DS.Where(goqu.I(front.OrderTable.PK()).Eq(payload.OrderID), goqu.I("$UserID").Eq(tokUsr.ID))
	if err = db.DsSelectOneTo(&order, ds); err != nil {
		return
	}
	if err = PermitOrderState(&order, front.TOrderStatePaid); err != nil {
		err = cerr.NoWayToPaidState
		return
	}

	now := time.Now().Unix()

	if payload.IsPoints && order.PayPoints != 0 {
		if order.PayPoints != payload.Amount {
			err = cerr.InvalidPayAmount
			return
		}

		var top front.PointsItem
		ds = dbs.DS.Where(goqu.I("$UserID").Eq(tokUsr.ID)).Order(goqu.I("$CreatedAt").Desc())
		if err = db.DsSelectOneTo(&top, ds); err != nil {
			glog.Errorln(err)
			return
		}

		if top.Balance < int(payload.Amount) {
			err = cerr.NotEnoughPoints
			return
		}

		err = db.Save(&front.PointsItem{
			UserID:    tokUsr.ID,
			CreatedAt: now,
			Amount:    -int(payload.Amount),
			Balance:   top.Balance - int(payload.Amount),
			OrderID:   order.ID,
		})
		if err != nil {
			glog.Errorln(err)
			return
		}

		order.PointsPaid = payload.Amount
		order.State = front.TOrderStatePaid
		order.PaidAt = now
		err = db.UpdateColumns(&order, "PointsPaid", "State", "PaidAt")
	} else if !payload.IsPoints && order.PayPoints == 0 {
		if order.PayAmount != payload.Amount {
			err = cerr.InvalidPayAmount
			return
		}

		var top front.UserCash
		ds = dbs.DS.Where(goqu.I("$UserID").Eq(tokUsr.ID)).Order(goqu.I("$CreatedAt").Desc())
		if err = db.DsSelectOneTo(&top, ds); err != nil {
			glog.Errorln(err)
			return
		}

		if top.Balance < int(payload.Amount) {
			err = cerr.NotEnoughMoney
			return
		}

		err = db.Save(&front.UserCash{
			UserID:    tokUsr.ID,
			CreatedAt: now,
			Type:      front.TUserCashTrade, // Trade type
			Amount:    -int(payload.Amount),
			Balance:   top.Balance - int(payload.Amount),
			OrderID:   order.ID,
		})
		if err != nil {
			glog.Errorln(err)
			return
		}

		order.CashPaid = payload.Amount
		order.State = front.TOrderStatePaid
		order.PaidAt = now
		err = db.UpdateColumns(&order, "CashPaid", "State", "PaidAt")
	} else {
		err = cerr.InvalidPayType
	}

	if err == nil {
		o = &order
	}

	return
}

// TOrderStateEvaled is standalone
func (dbs *DbService) OrderChangeState(
	order *front.Order, tokUsr *models.User, payload *front.OrderChangeStatePayload,
) (err error) {

	db := dbs.GetDB()

	ds := dbs.DS.Where(goqu.I(front.OrderTable.PK()).Eq(payload.ID)).Where(goqu.I("$UserID").Eq(tokUsr.ID))
	err = db.DsSelectOneTo(order, ds)
	if err != nil {
		return
	}

	if err = PermitOrderState(order, payload.State); err != nil {
		err = cerr.NoWayToTargetState
		return
	}

	now := time.Now().Unix()
	switch payload.State {
	case front.TOrderStateCanceled:
		switch order.State {
		case front.TOrderStateNopay:
			order.State = front.TOrderStateCanceled
			order.CanceledAt = now
			err = db.UpdateColumns(order, "CanceledAt", "State")
			return
		case front.TOrderStatePrepaid:
			var res map[string]string
			res, err = dbs.wc.OrderClose(order)
			if err != nil {
				return
			}

			if res["result_code"] != "SUCCESS" {
				switch res["err_code"] {
				case "ORDERPAID":
					//					err = cerr.WxOrderAlreadyPaid
					order.State = front.TOrderStatePaid
					order.PaidAt = now
					order.WxTradeState = front.SUCCESS
					err = db.UpdateColumns(order, "State", "PaidAt", "WxTradeState")
					return
				case "SYSTEMERROR":
					err = cerr.WxSystemFailed
				case "SIGNERROR", "REQUIRE_POST_METHOD", "XML_FORMAT_ERROR":
					err = cerr.ApiImplementFailed
				}
			}
			if err != nil {
				return
			}

			order.State = front.TOrderStateCanceled
			order.CanceledAt = now
			err = db.UpdateColumns(order, "CanceledAt", "State")
			return

		case front.TOrderStatePaid, front.TOrderStatePicking:
			if order.CashRefund == 0 {
				order.CashRefund = order.CashPaid
			}
			if order.WxRefund == 0 {
				order.WxRefund = order.WxPaid
			}
			err = dbs.orderRefund(tokUsr.ID, order)
			if err != nil {
				return
			}
			order.State = front.TOrderStateCanceled
			order.CanceledAt = now
			err = db.UpdateColumns(order, "State", "CanceledAt", "WxRefundID")
			return
		} // to cancel end

	case front.TOrderStateCompleted:
		if dbs.IsOrderCompleted(order) {
			err = cerr.OrderCompleteTimeout
			return
		}
		order.State = front.TOrderStateCompleted
		order.CompletedAt = now

		// TODO prove it!
		var cols []string
		cols, err = dbs.OrderMaintanence(order)
		if err != nil {
			return
		}
		err = db.UpdateColumns(order, append(cols, "CompletedAt", "State")...)
		return

	case front.TOrderStateReturnStarted:
		if order.PointsPaid != 0 {
			return cerr.NoAbcOrPoints
		}
		if dbs.IsOrderCompleted(order) {
			err = cerr.OrderCompleteTimeout
			return
		}
		order.State = front.TOrderStateReturnStarted
		order.ReturnStaredtAt = now
		err = db.UpdateColumns(order, "ReturnStaredtAt", "State")
		return

	}

	err = cerr.Forbidden
	return
}

func (dbs *DbService) MgrOrderState(order *front.Order, claims *admin.Claims) (err error) {
	db := dbs.GetDB()

	ds := dbs.DS.Where(goqu.I(front.OrderTable.PK()).Eq(claims.OrderID))
	err = db.DsSelectOneTo(order, ds)
	if err != nil {
		return
	}

	if err = PermitOrderState(order, claims.State); err != nil {
		glog.Errorln(order.State)
		err = cerr.NoWayToTargetState
		return
	}

	var refundCol string // need refund
	now := time.Now().Unix()
	switch claims.State {
	case front.TOrderStatePicking:
		order.State = front.TOrderStatePicking
		order.PickingAt = now
		err = db.UpdateColumns(order, "PickingAt", "State")

	case front.TOrderStateDelivered:
		order.State = front.TOrderStateDelivered
		order.DeliveredAt = now
		order.DeliverCom = claims.DeliverCom
		order.DeliverNo = claims.DeliverNo
		err = db.UpdateColumns(order, "DeliveredAt", "State", "DeliverCom", "DeliverNo")

	case front.TOrderStateRejecting:
		order.State = front.TOrderStateRejecting
		order.RejectingAt = now
		err = db.UpdateColumns(order, "RejectingAt", "State")

	case front.TOrderStateRejectBack:
		order.State = front.TOrderStateRejectBack
		order.RejectBackAt = now
		err = db.UpdateColumns(order, "RejectBackAt", "State")

	case front.TOrderStateReturning:
		order.State = front.TOrderStateReturning
		order.ReturnEnsuredAt = now
		err = db.UpdateColumns(order, "ReturnEnsuredAt", "State")

	case front.TOrderStateRejectRefound:
		order.State = front.TOrderStateRejectRefound
		order.RejectRefoundAt = now
		refundCol = "RejectRefoundAt"
	case front.TOrderStateReturned:
		order.State = front.TOrderStateReturned
		order.ReturnedAt = now
		refundCol = "ReturnedAt"
	default:
		err = cerr.NoPermToState
	}

	if refundCol == "" || err != nil {
		return
	}

	// refund
	order.CashRefund = claims.CashRefund
	order.WxRefund = claims.WxRefund
	if order.CashRefund+order.WxRefund > order.PayAmount {
		err = cerr.InvalidPayAmount
		return
	}

	err = dbs.orderRefund(order.UserID, order)
	if err == nil {
		err = db.UpdateColumns(order, refundCol, "CashRefund", "WxRefund", "State", "WxRefundID")
	} else {
		glog.Errorln(err)
	}
	return
}

// called just before commit of tx
func (dbs *DbService) orderRefund(userId uint, order *front.Order) (err error) {
	db := dbs.GetDB()
	now := time.Now().Unix()

	if order.CashRefund != 0 {
		var flow front.UserCash
		ds := dbs.DS.Where(goqu.I("$OrderID").Eq(userId)).Where(goqu.I("$Type").Eq(front.TUserCashRefund))
		if err = db.DsSelectOneTo(&flow, ds); err != nil && err != reform.ErrNoRows {
			return
		}
		err = nil
		if flow.ID == 0 {
			// not refund yet
			var flow1 front.UserCash
			ds = dbs.DS.Where(goqu.I("$UserID").Eq(order.UserID)).Order(goqu.I("$CreatedAt").Desc())
			if err = db.DsSelectOneTo(&flow1, ds); err != nil {
				return
			}
			// refund
			err = db.Save(&front.UserCash{
				UserID:    userId,
				CreatedAt: now,
				Type:      front.TUserCashRefund, // Trade type
				Amount:    int(order.CashRefund),
				Balance:   flow1.Balance + int(order.CashRefund),
				OrderID:   order.ID,
			})
			if err != nil {
				return
			}
		}
	}

	if order.WxRefund != 0 {
		var res map[string]string
		res, err = dbs.wc.OrderRefund(order)
		if err != nil {
			return
		}

		if res["result_code"] != "SUCCESS" {
			switch res["err_code"] {
			case "SYSTEMERROR":
				err = cerr.WxSystemFailed
			case "USER_ACCOUNT_ABNORMAL":
				err = cerr.WxUserAbnormal
			//					case "INVALID_TRANSACTIONID":
			//						return cerr.WxSystemFailed
			//					case "SIGNERROR", "PARAM_ERROR", "APPID_NOT_EXIST", "MCHID_NOT_EXIST",
			//						"APPID_MCHID_NOT_MATCH", "REQUIRE_POST_METHOD", "XML_FORMAT_ERROR":
			//						return cerr.ApiImplementFailed
			default:
				err = cerr.ApiImplementFailed
			}
			if err != nil {
				return
			}
		}

		order.WxRefundID = res["refund_id"]
	}

	err = nil
	return
}

func (dbs *DbService) OrderPaidState(tokUsr *models.User, orderId uint) (order *front.Order, err error) {
	order, err = dbs.GetBareOrder(tokUsr, orderId)
	if err != nil {
		return
	}
	var src map[string]string
	src, err = dbs.wc.OrderQuery(order)
	if err != nil {
		return
	}
	return order, dbs.UpdateWxOrderSate(order, src)
}

func (dbs *DbService) OnWxPayNotify(src map[string]string, orderId uint) error {
	order, err := dbs.GetBareOrder(nil, orderId)
	if err != nil {
		return err
	}
	return dbs.UpdateWxOrderSate(order, src)
}

func (dbs *DbService) UpdateWxOrderSate(order *front.Order, src map[string]string) (err error) {
	tradeState := front.TradeStateNameToValue[src["trade_state"]]
	tid := src["transaction_id"]
	if order.WxTradeState == tradeState && order.WxTransactionId == tid {
		// no need update
		return
	}

	totalFee64, _ := strconv.ParseUint(src["total_fee"], 10, 64)
	if totalFee64 == 0 {
		err = cerr.ParseWxTotalFeeFailed
		return
	}

	switch tradeState {
	case front.SUCCESS:
		if err = PermitOrderState(order, front.TOrderStatePaid); err != nil {
			log.WithFields(l.Locate(logrus.Fields{
				"OrderID": order.ID,
				"State":   order.State,
			})).Info("Got wxpay with SUCCESS")
			return
		}

		timeEnd, errTime := time.Parse("20060102150405", src["time_end"])
		if errTime != nil {
			timeEnd = time.Now()
		}

		data := front.Order{
			ID:              order.ID,
			WxPaid:          uint(totalFee64),
			WxTransactionId: tid,
			WxTradeState:    tradeState,
			State:           front.TOrderStatePaid,
			PaidAt:          timeEnd.Unix(),
		}
		err = dbs.GetDB().UpdateColumns(&data, "WxPaid", "WxTransactionId", "WxTradeState", "State", "PaidAt")
		if err != nil {
			return
		}
		order.WxPaid = data.WxPaid
		order.WxTransactionId = data.WxTransactionId
		order.WxTradeState = data.WxTradeState
		order.State = data.State
		order.PaidAt = data.PaidAt

	case front.REFUND, front.USERPAYING, front.PAYERROR:
		err = dbs.GetDB().UpdateColumns(&front.Order{ID: order.ID, WxTradeState: tradeState}, "WxTradeState")
		if err == nil {
			order.WxTradeState = tradeState
		}

	case front.CLOSED:
		// must be an expired order, ignore
		// refound must be called with SUCCESS on closeorder
		// web must refresh!
		err = cerr.OrderClosed
	}

	return
}

func (dbs *DbService) GetBareOrder(tokUsr *models.User, id uint) (*front.Order, error) {
	var order front.Order
	ds := dbs.DS.Where(goqu.I(front.OrderTable.PK()).Eq(id))
	if tokUsr != nil {
		ds = ds.Where(goqu.I("$UserID").Eq(tokUsr.ID))
	}
	return &order, dbs.GetDB().DsSelectOneTo(&order, ds)
}

func (dbs *DbService) GetOrder1(tokUsr *models.User, id uint) (*front.Order, error) {
	var order front.Order
	ds := dbs.DS.Where(goqu.I("$UserID").Eq(tokUsr.ID), goqu.I(front.OrderTable.PK()).Eq(id))
	err := dbs.GetDB().DsSelectOneTo(&order, ds)
	if err != nil {
		return nil, err
	}

	if _, err = dbs.GetOrderItems(&order); err != nil {
		return nil, err
	}

	return &order, nil
}
