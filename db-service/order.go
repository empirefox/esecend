package dbsrv

import (
	"strconv"
	"strings"
	"time"

	"gopkg.in/doug-martin/goqu.v3"

	"github.com/Sirupsen/logrus"
	"github.com/cznic/sortutil"
	"github.com/empirefox/esecend/admin"
	"github.com/empirefox/esecend/cerr"
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/esecend/l"
	"github.com/empirefox/esecend/lok"
	"github.com/empirefox/esecend/models"
	"github.com/empirefox/reform"
)

type WxPaier interface {
	UnifiedOrder(tokUsr *models.User, order *front.Order, ip string, attach *models.UnifiedOrderAttach) (*front.WxPayArgs, error)
	OrderClose(order *front.Order) (map[string]string, error)
	OrderRefund(order *front.Order, opUserId string) (map[string]string, error)
}

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

// TODO change to single thread
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
	if len(skus) != len(skuIds) {
		return nil, cerr.InvalidSkuId
	}

	gbs, err := db.FindAllFromPK(front.GroupBuyItemTable, groupbuyIds...)
	if err != nil {
		return nil, err
	}
	if len(gbs) != len(groupbuyIds) {
		return nil, cerr.InvalidGroupbuyId
	}
	gbMap := make(map[uint]*front.GroupBuyItem)
	for _, gbi := range gbs {
		gb := gbi.(*front.GroupBuyItem)
		gbMap[gb.SkuID] = gb
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
	// Attrs: prepare attr ids
	attrIdMap := make(map[uint]bool)
	for _, skuToAttri := range skuToAttrs {
		skuToAttr := skuToAttri.(*front.ProductAttrId)
		attrIdMap[skuToAttr.AttrID] = true

		// Attrs: save ids to payload item
		payloadItem := skuidToPayloadItem[skuToAttr.SkuID]
		payloadItem.AttrIds = append(payloadItem.AttrIds, skuToAttr.AttrID)
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

	// complete order items after order saved
	for _, item := range items {
		product := productMap[item.ProductID]
		item.UserID = tokUsr.ID
		item.OrderID = order.ID
		item.Name = product.Name
		if item.Img == "" {
			item.Img = product.Img
		}
		item.Attrs = skuidToPayloadItem[item.SkuID].AttrValues

		// update sku stock
		sku := skuMap[item.SkuID]
		sku.Stock -= item.Quantity
		if err := db.UpdateColumns(sku, "Stock"); err != nil {
			return nil, err
		}
		if err := db.Insert(item); err != nil {
			return nil, err
		}
	}
	order.Items = items

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &order, nil
}

// only used by PrepayOrder
func (tx *DbService) prepayOrderAfterClosedWxOrder(
	tokUsr *models.User, order *front.Order, paier WxPaier, ip string,
) (args *front.WxPayArgs, cashLocked, pointsLocked bool, err error) {

	db := tx.GetDB()

	if cashLocked = lok.CashLok.Lock(tokUsr.ID); !cashLocked {
		err = cerr.CashTmpLocked
		return
	}

	var flow front.CapitalFlow
	ds := tx.DS.Where(goqu.I("$OrderID").Eq(order.ID)).Where(goqu.I("$Type").Eq(front.TCapitalFlowPrepay))
	if err = db.DsSelectOneTo(&flow, ds); err != nil && err != reform.ErrNoRows {
		return
	}

	if pointsLocked = lok.PointsLok.Lock(tokUsr.ID); !pointsLocked {
		err = cerr.PointsTmpLocked
		return
	}

	var points front.PointsItem
	ds = tx.DS.Where(goqu.I("$OrderID").Eq(order.ID)).Where(goqu.I("$Type").Eq(front.TPointsPrepay))
	if err = db.DsSelectOneTo(&points, ds); err != nil && err != reform.ErrNoRows {
		return
	}

	if uint(-flow.Amount) == order.CashPaid && uint(-points.Amount)*tx.config.Order.Point2Cent == order.PointsPaid {
		// no need prepay again

		if err = db.UpdateColumns(&front.Order{ID: order.ID}, "TransactionId", "TradeState"); err != nil {
			return
		}

		attach := &models.UnifiedOrderAttach{
			PreCashID:   flow.ID,
			CashPaid:    order.CashPaid,
			PrePointsID: points.ID,
			PointsPaid:  order.PointsPaid,
			UserID:      tokUsr.ID,
		}

		args, err = paier.UnifiedOrder(tokUsr, order, ip, attach)
		return
	}

	// refund is allowed only when close order
	if flow.ID != 0 || points.ID != 0 {
		err = cerr.OrderCloseNeeded
	}

	// args==nil err==nil means go on
	return
}

// prepay with cash and points, then get wx prepay_id
// cannot change prepaid when do the 2nd time
func (dbs *DbService) PrepayOrder(
	tokUsr *models.User, payload *front.OrderPrepayPayload, paier WxPaier,
) (args *front.WxPayArgs, cashLocked, pointsLocked bool, err error) {

	if payload.Wx == 0 {
		err = cerr.NotPrepayOrder
		return
	}
	pointsPaid := payload.Points * dbs.config.Order.Point2Cent
	if payload.Amount != payload.Cash+payload.Wx+pointsPaid || payload.Amount == 0 {
		err = cerr.InvalidPayAmount
		return
	}

	db := dbs.GetDB()

	var order front.Order
	ds := dbs.DS.Where(goqu.I(front.OrderTable.PK()).Eq(payload.OrderID)).Where(goqu.I("$UserID").Eq(tokUsr.ID))
	if err = db.DsSelectOneTo(&order, ds); err != nil {
		return
	}

	// prepay after prepaid
	if order.State == front.TOrderStatePrepaid {
		if order.CashPaid != payload.Cash || pointsPaid != order.PointsPaid {
			err = cerr.InvalidPrepayPayload
			return
		}
		// close transaction_id firstly
		var res map[string]string
		res, err = paier.OrderClose(&order)
		if err != nil {
			return
		}

		if res["result_code"] == "SUCCESS" {
			var flow front.CapitalFlow
			if order.CashPaid != 0 {
				if cashLocked = lok.CashLok.Lock(tokUsr.ID); !cashLocked {
					err = cerr.CashTmpLocked
					return
				}

				ds = dbs.DS.Where(goqu.I("$OrderID").Eq(order.ID)).Where(goqu.I("$Type").Eq(front.TCapitalFlowPrepay))
				if err = db.DsSelectOneTo(&flow, ds); err != nil && err != reform.ErrNoRows {
					return
				}
				if uint(-flow.Amount) != order.CashPaid {
					err = cerr.InvalidCashPrepaid
					return
				}

			}

			var points front.PointsItem
			if order.PointsPaid != 0 {
				if pointsLocked = lok.PointsLok.Lock(tokUsr.ID); !pointsLocked {
					err = cerr.PointsTmpLocked
					return
				}

				ds = dbs.DS.Where(goqu.I("$OrderID").Eq(order.ID)).Where(goqu.I("$Type").Eq(front.TPointsPrepay))
				if err = db.DsSelectOneTo(&points, ds); err != nil && err != reform.ErrNoRows {
					return
				}
				if uint(-points.Amount)*dbs.config.Order.Point2Cent != order.PointsPaid {
					err = cerr.InvalidPointsPrepaid
					return
				}

				if err = db.UpdateColumns(&front.Order{ID: order.ID}, "TransactionId", "TradeState"); err != nil {
					return
				}
			}

			attach := &models.UnifiedOrderAttach{
				PreCashID:   flow.ID,
				CashPaid:    order.CashPaid,
				PrePointsID: points.ID,
				PointsPaid:  order.PointsPaid,
				UserID:      tokUsr.ID,
			}

			args, err = paier.UnifiedOrder(tokUsr, &order, payload.Ip, attach)
			return
		} else {
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
			case "ORDERNOTEXIST", "ORDERCLOSED":
				// TODO remove this?
				// check prepaid
				args, cashLocked, pointsLocked, err = dbs.prepayOrderAfterClosedWxOrder(tokUsr, &order, paier, payload.Ip)
				if args != nil {
					// return OK!
					return
				}
				// to be continued, do DO NOT return!

			default:
				err = cerr.ApiImplementFailed
			}
			if err != nil {
				return
			}
		}
	} else if err = PermitOrderState(&order, front.TOrderStatePrepaid); err != nil {
		err = cerr.NotNopayState
		return
	}

	if order.PayAmount != payload.Amount {
		err = cerr.InvalidPayAmount
		return
	}

	now := time.Now().Unix()
	attach := models.UnifiedOrderAttach{
		UserID: tokUsr.ID,
	}

	if payload.Cash != 0 {
		if !cashLocked {
			cashLocked = lok.CashLok.Lock(tokUsr.ID)
		}
		if !cashLocked {
			err = cerr.CashTmpLocked
			return
		}

		var flow front.CapitalFlow
		ds = dbs.DS.Where(goqu.I("$UserID").Eq(tokUsr.ID)).Order(goqu.I("$CreatedAt").Desc().NullsLast())
		if err = db.DsSelectOneTo(&flow, ds); err != nil {
			return
		}

		if flow.Balance < int(payload.Cash) {
			err = cerr.NotEnoughMoney
			return
		}

		err = db.Save(&front.CapitalFlow{
			UserID:    tokUsr.ID,
			CreatedAt: now,
			Type:      front.TCapitalFlowPrepay,
			Amount:    -int(payload.Cash),
			Balance:   flow.Balance - int(payload.Cash),
			OrderID:   order.ID,
		})
		if err != nil {
			return
		}
		attach.CashPaid = payload.Cash
		attach.PreCashID = flow.ID
	}

	if payload.Points != 0 {
		if !pointsLocked {
			pointsLocked = lok.PointsLok.Lock(tokUsr.ID)
		}
		if !pointsLocked {
			err = cerr.PointsTmpLocked
			return
		}

		var points front.PointsItem
		ds = dbs.DS.Where(goqu.I("$UserID").Eq(tokUsr.ID)).Order(goqu.I("$CreatedAt").Desc().NullsLast())
		if err = db.DsSelectOneTo(&points, ds); err != nil {
			return
		}

		if points.Balance < int(payload.Points) {
			err = cerr.NotEnoughPoints
			return
		}

		err = db.Save(&front.PointsItem{
			UserID:    tokUsr.ID,
			CreatedAt: now,
			Type:      front.TPointsPrepay,
			Amount:    -int(payload.Points),
			Balance:   points.Balance - int(payload.Points),
			OrderID:   order.ID,
		})
		if err != nil {
			return
		}
		attach.PointsPaid = pointsPaid
		attach.PrePointsID = points.ID
	}

	order.CashPaid = attach.CashPaid
	order.PointsPaid = attach.PointsPaid
	order.State = front.TOrderStatePrepaid
	order.PrepaidAt = now
	order.TransactionId = ""
	order.TradeState = front.UNKNOWN

	err = db.UpdateColumns(&order, "CashPaid", "PointsPaid", "State", "PrepaidAt", "TransactionId", "TradeState")

	if err == nil {
		args, err = paier.UnifiedOrder(tokUsr, &order, payload.Ip, &attach)
	}

	return
}

// no wx pay, need paykey
func (dbs *DbService) PayOrder(order *front.Order, tokUsr *models.User, payload *front.OrderPayPayload) (cashLocked, pointsLocked bool, err error) {
	pointsPaid := payload.Points * dbs.config.Order.Point2Cent
	if payload.Amount != payload.Cash+pointsPaid || payload.Amount == 0 {
		err = cerr.InvalidPayAmount
		return
	}

	db := dbs.GetDB()

	if err = db.Reload(tokUsr); err != nil {
		return
	}

	if err = models.ComparePaykey([]byte(tokUsr.Paykey), []byte(payload.Key)); err != nil {
		err = cerr.InvalidPaykey
		return
	}

	ds := dbs.DS.Where(goqu.I(front.OrderTable.PK()).Eq(payload.OrderID)).Where(goqu.I("$UserID").Eq(tokUsr.ID))
	if err = db.DsSelectOneTo(order, ds); err != nil {
		return
	}
	if err = PermitOrderState(order, front.TOrderStatePaid); err != nil {
		err = cerr.NoWayToPaidState
		return
	}
	if order.PayAmount != payload.Amount {
		err = cerr.InvalidPayAmount
		return
	}

	now := time.Now().Unix()

	if payload.Cash != 0 {
		if cashLocked = lok.CashLok.Lock(tokUsr.ID); !cashLocked {
			err = cerr.CashTmpLocked
			return
		}

		var flow front.CapitalFlow
		ds = dbs.DS.Where(goqu.I("$UserID").Eq(tokUsr.ID)).Order(goqu.I("$CreatedAt").Desc().NullsLast())
		if err = db.DsSelectOneTo(&flow, ds); err != nil {
			return
		}

		if flow.Balance < int(payload.Cash) {
			err = cerr.NotEnoughMoney
			return
		}

		err = db.Save(&front.CapitalFlow{
			UserID:    tokUsr.ID,
			CreatedAt: now,
			Type:      front.TCapitalFlowTrade, // Trade type
			Amount:    -int(payload.Cash),
			Balance:   flow.Balance - int(payload.Cash),
			OrderID:   order.ID,
		})
		if err != nil {
			return
		}
	}

	if payload.Points != 0 {
		if pointsLocked = lok.PointsLok.Lock(tokUsr.ID); !pointsLocked {
			err = cerr.PointsTmpLocked
			return
		}

		var points front.PointsItem
		ds = dbs.DS.Where(goqu.I("$UserID").Eq(tokUsr.ID)).Order(goqu.I("$CreatedAt").Desc().NullsLast())
		if err = db.DsSelectOneTo(&points, ds); err != nil {
			return
		}

		if points.Balance < int(payload.Points) {
			err = cerr.NotEnoughPoints
			return
		}

		err = db.Save(&front.PointsItem{
			UserID:    tokUsr.ID,
			CreatedAt: now,
			Amount:    -int(payload.Points),
			Balance:   points.Balance - int(payload.Points),
			Type:      front.TPointsTrade, // Trade type
			OrderID:   order.ID,
		})
		if err != nil {
			return
		}
	}

	order.CashPaid = payload.Cash
	order.PointsPaid = pointsPaid
	order.State = front.TOrderStatePaid
	order.PaidAt = now
	err = db.UpdateColumns(order, "CashPaid", "PointsPaid", "State", "PaidAt")
	return
}

// TOrderStateEvaled is standalone
func (dbs *DbService) OrderChangeState(
	order *front.Order, tokUsr *models.User, payload *front.OrderChangeStatePayload, paier WxPaier,
) (cashLocked, pointsLocked bool, err error) {

	db := dbs.GetDB()

	// lock order before tx
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
			if cashLocked = lok.CashLok.Lock(tokUsr.ID); !cashLocked {
				err = cerr.CashTmpLocked
				return
			}

			var flow front.CapitalFlow
			ds = dbs.DS.Where(goqu.I("$OrderID").Eq(order.ID)).Where(goqu.I("$Type").Eq(front.TCapitalFlowPrepay))
			if err = db.DsSelectOneTo(&flow, ds); err != nil && err != reform.ErrNoRows {
				return
			}
			if flow.ID != 0 {
				// repaid
				var flow2 front.CapitalFlow
				ds = dbs.DS.Where(goqu.I("$OrderID").Eq(order.ID)).Where(goqu.I("$Type").Eq(front.TCapitalFlowPrepayBack))
				if err = db.DsSelectOneTo(&flow2, ds); err != nil && err != reform.ErrNoRows {
					return
				}
				if flow2.ID == 0 {
					// not refund yet
					var flow1 front.CapitalFlow
					ds = dbs.DS.Where(goqu.I("$UserID").Eq(tokUsr.ID)).Order(goqu.I("$CreatedAt").Desc().NullsLast())
					if err = db.DsSelectOneTo(&flow1, ds); err != nil {
						return
					}
					// refund
					err = db.Save(&front.CapitalFlow{
						UserID:    tokUsr.ID,
						CreatedAt: now,
						Type:      front.TCapitalFlowPrepayBack, // Trade type
						Amount:    -flow.Amount,
						Balance:   flow1.Balance - flow.Amount,
						OrderID:   order.ID,
					})
					if err != nil {
						return
					}
				}
			}

			if pointsLocked = lok.PointsLok.Lock(tokUsr.ID); !pointsLocked {
				err = cerr.PointsTmpLocked
				return
			}

			var points front.PointsItem
			ds = dbs.DS.Where(goqu.I("$OrderID").Eq(order.ID)).Where(goqu.I("$Type").Eq(front.TPointsPrepay))
			if err = db.DsSelectOneTo(&points, ds); err != nil && err != reform.ErrNoRows {
				return
			}
			if points.ID != 0 {
				// repaid
				var points2 front.PointsItem
				ds = dbs.DS.Where(goqu.I("$OrderID").Eq(order.ID)).Where(goqu.I("$Type").Eq(front.TPointsPrepayBack))
				if err = db.DsSelectOneTo(&points2, ds); err != nil && err != reform.ErrNoRows {
					return
				}
				if points2.ID == 0 {
					// not refund yet
					var points1 front.PointsItem
					ds = dbs.DS.Where(goqu.I("$UserID").Eq(tokUsr.ID)).Order(goqu.I("$CreatedAt").Desc().NullsLast())
					if err = db.DsSelectOneTo(&points1, ds); err != nil {
						return
					}
					// refund
					err = db.Save(&front.PointsItem{
						UserID:    tokUsr.ID,
						CreatedAt: now,
						Type:      front.TPointsPrepayBack, // Trade type
						Amount:    -points.Amount,
						Balance:   points1.Balance - points.Amount,
						OrderID:   order.ID,
					})
					if err != nil {
						return
					}
				}
			}

			order.State = front.TOrderStateCanceled
			order.CanceledAt = now
			err = db.UpdateColumns(order, "CanceledAt", "State")
			if err != nil {
				return
			}

			// TOrderStatePrepaid indecate there must be an incompleted wx order
			var res map[string]string
			res, err = paier.OrderClose(order)
			if err != nil {
				return
			}

			if res["result_code"] != "SUCCESS" {
				switch res["err_code"] {
				case "ORDERPAID":
					err = cerr.WxOrderAlreadyPaid
				case "SYSTEMERROR":
					err = cerr.WxSystemFailed
				case "SIGNERROR", "REQUIRE_POST_METHOD", "XML_FORMAT_ERROR":
					err = cerr.ApiImplementFailed
				}
			}
			return

		case front.TOrderStatePaid, front.TOrderStatePicking:
			if order.CashRefund == 0 {
				order.CashRefund = order.CashPaid
			}
			if order.PointsRefund == 0 {
				order.PointsRefund = order.PointsPaid
			}
			if order.WxRefund == 0 {
				order.WxRefund = order.WxPaid
			}
			cashLocked, pointsLocked, err = dbs.orderRefund(0, tokUsr.ID, order, paier)
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
		err = db.UpdateColumns(order, "CompletedAt", "State")
		return

	case front.TOrderStateReturnStarted:
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

func (dbs *DbService) MgrOrderState(order *front.Order, claims *admin.Claims, paier WxPaier) (cashLocked, pointsLocked bool, err error) {
	db := dbs.GetDB()

	// lock order before tx
	ds := dbs.DS.Where(goqu.I(front.OrderTable.PK()).Eq(claims.OrderID))
	//		Where(goqu.I("$CreatedAt").Eq(payload.CreatedAt)).
	//		Where(goqu.I("$UserID").Eq(tokUsr.ID))
	err = db.DsSelectOneTo(order, ds)
	if err != nil {
		return
	}

	if err = PermitOrderState(order, claims.State); err != nil {
		err = cerr.NoWayToTargetState
		return
	}

	var refundCol string
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
		err = cerr.NoWayToTargetState
	}

	if err != nil {
		return
	}

	order.CashRefund = claims.CashRefund
	order.PointsRefund = claims.PointsRefund
	order.WxRefund = claims.WxRefund

	cashLocked, pointsLocked, err = dbs.orderRefund(claims.AminId, claims.UserId, order, paier)
	if err == nil {
		err = db.UpdateColumns(order, refundCol, "State", "WxRefundID")
	}
	return
}

// called just before commit of tx
func (dbs *DbService) orderRefund(adminId, userId uint, order *front.Order, paier WxPaier) (cashLocked, pointsLocked bool, err error) {
	db := dbs.GetDB()
	now := time.Now().Unix()

	if order.CashRefund != 0 {
		if cashLocked = lok.CashLok.Lock(userId); !cashLocked {
			err = cerr.CashTmpLocked
			return
		}

		var flow front.CapitalFlow
		ds := dbs.DS.Where(goqu.I("$OrderID").Eq(order.ID)).Where(goqu.I("$Type").Eq(front.TCapitalFlowRefund))
		if err = db.DsSelectOneTo(&flow, ds); err != nil && err != reform.ErrNoRows {
			return
		}
		if flow.ID == 0 {
			// not refund yet
			var flow1 front.CapitalFlow
			ds = dbs.DS.Where(goqu.I("$UserID").Eq(userId)).Order(goqu.I("$CreatedAt").Desc().NullsLast())
			if err = db.DsSelectOneTo(&flow1, ds); err != nil {
				return
			}
			// refund
			err = db.Save(&front.CapitalFlow{
				UserID:    userId,
				CreatedAt: now,
				Type:      front.TCapitalFlowRefund, // Trade type
				Amount:    int(order.CashRefund),
				Balance:   flow1.Balance + int(order.CashRefund),
				OrderID:   order.ID,
			})
			if err != nil {
				return
			}
		}
	}

	if order.PointsRefund != 0 {
		if pointsLocked = lok.PointsLok.Lock(userId); !pointsLocked {
			err = cerr.PointsTmpLocked
			return
		}

		var points front.PointsItem
		ds := dbs.DS.Where(goqu.I("$OrderID").Eq(order.ID)).Where(goqu.I("$Type").Eq(front.TPointsRefund))
		if err = db.DsSelectOneTo(&points, ds); err != nil && err != reform.ErrNoRows {
			return
		}
		if points.ID == 0 {
			// not refund yet
			var points1 front.PointsItem
			ds = dbs.DS.Where(goqu.I("$UserID").Eq(userId)).Order(goqu.I("$CreatedAt").Desc().NullsLast())
			if err = db.DsSelectOneTo(&points1, ds); err != nil {
				return
			}
			// refund
			err = db.Save(&front.PointsItem{
				UserID:    userId,
				CreatedAt: now,
				Type:      front.TPointsRefund, // Trade type
				Amount:    int(order.PointsRefund),
				Balance:   points1.Balance + int(order.PointsRefund),
				OrderID:   order.ID,
			})
			if err != nil {
				return
			}
		}
	}

	if order.WxRefund != 0 {
		var res map[string]string
		if adminId == 0 {
			adminId = userId
		}
		res, err = paier.OrderRefund(order, strconv.Itoa(int(adminId)))
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

func (dbs *DbService) IsOrderCompleted(order *front.Order) bool {
	now := time.Now().Unix()
	return order.CompletedAt != 0 || order.EvalStartedAt != 0 || order.EvalAt != 0 ||
		time.Duration(now-order.DeliveredAt) > dbs.config.Order.CompleteTimeoutDay*3600*24
}

func (dbs *DbService) IsEvalTimeout(order *front.Order) bool {
	now := time.Now().Unix()
	if order.CompletedAt == 0 {
		if time.Duration(now-order.DeliveredAt) > (dbs.config.Order.CompleteTimeoutDay+dbs.config.Order.EvalTimeoutDay)*3600*24 {
			return true
		}
	} else if time.Duration(now-order.CompletedAt) > dbs.config.Order.EvalTimeoutDay*3600*24 {
		return true
	}
	return false
}
