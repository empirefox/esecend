package dbsrv

import (
	"strconv"
	"strings"
	"time"

	"gopkg.in/doug-martin/goqu.v3"
	"gopkg.in/reform.v1"

	"github.com/Sirupsen/logrus"
	"github.com/cznic/sortutil"
	"github.com/empirefox/esecend/cerr"
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/esecend/l"
	"github.com/empirefox/esecend/lok"
	"github.com/empirefox/esecend/models"
)

type WxPaier interface {
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
	profile, err := db.SelectOneFrom(front.ProfileView, "LIMIT 1")
	if err != nil {
		return nil, err
	}
	if total >= profile.(*front.Profile).FreeDeliverLine {
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
	// TODO save user_id?
	for _, item := range items {
		product := productMap[item.ProductID]
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
	tokUsr *models.User,
	order *front.Order,
	unifiedOrder func(order *front.Order, attach *models.UnifiedOrderAttach) (*front.WxPayArgs, error),
) (*front.WxPayArgs, error) {

	db := tx.GetDB()

	if !lok.CashLok.Lock(tokUsr.ID) {
		return nil, cerr.CashTmpLocked
	}
	defer lok.CashLok.Unlock(tokUsr.ID)

	var flow front.CapitalFlow
	ds := tx.DS.Where(goqu.I("$OrderID").Eq(order.ID)).Where(goqu.I("$Type").Eq(front.TCapitalFlowPrepay))
	if err := db.DsSelectOneTo(&flow, ds); err != nil && err != reform.ErrNoRows {
		return nil, err
	}

	if !lok.PointsLok.Lock(tokUsr.ID) {
		return nil, cerr.CashTmpLocked
	}
	defer lok.PointsLok.Unlock(tokUsr.ID)

	var points front.PointsItem
	ds = tx.DS.Where(goqu.I("$OrderID").Eq(order.ID)).Where(goqu.I("$Type").Eq(front.TPointsPrepay))
	if err := db.DsSelectOneTo(&points, ds); err != nil && err != reform.ErrNoRows {
		return nil, err
	}

	if uint(-flow.Amount) == order.CashPaid && uint(-points.Amount)*dbs.config.Order.Point2Cent == order.PointsPaid {
		// no need prepay again

		if err = db.UpdateColumns(&front.Order{ID: order.ID}, "TransactionId", "TradeState"); err != nil {
			return nil, err
		}

		attach := models.UnifiedOrderAttach{
			PreCashID:   flow.ID,
			CashPaid:    order.CashPaid,
			PrePointsID: points.ID,
			PointsPaid:  order.PointsPaid,
			UserID:      tokUsr.ID,
		}

		args, err := unifiedOrder(&order, &attach)
		if err != nil {
			return nil, err
		}

		if err := tx.Commit(); err != nil {
			return nil, err
		}

		order.TransactionId = ""
		order.TradeState = front.UNKNOWN
		return args, nil
	}

	// refund is allowed only when close order
	if flow.ID != 0 || points.ID != 0 {
		return nil, cerr.OrderCloseNeeded
	}

	// means go on
	return nil, nil
}

// prepay with cash and points, then get wx prepay_id
// cannot change prepaid when do the 2nd time
func (dbs *DbService) PrepayOrder(
	tokUsr *models.User,
	payload *front.OrderPrepayPayload,
	unifiedOrder func(order *front.Order, attach *models.UnifiedOrderAttach) (*front.WxPayArgs, error),
	closeWxOrder func(order *front.Order) (map[string]string, error),
) (*front.WxPayArgs, error) {

	if payload.Wx == 0 {
		return nil, cerr.NotPrepayOrder
	}
	pointsPaid := payload.Points * dbs.config.Order.Point2Cent
	if payload.Amount != payload.Cash+payload.Wx+pointsPaid || payload.Amount == 0 {
		return nil, cerr.InvalidPayAmount
	}

	if !lok.OrderLok.Lock(payload.OrderID) {
		return nil, cerr.OrderTmpLocked
	}
	defer lok.OrderLok.Unlock(payload.OrderID)

	tx, err := dbs.Tx()
	if err != nil {
		return nil, err
	}
	var wxOrderClose bool
	defer tx.RollbackIfNeeded()
	db := tx.GetDB()

	var order front.Order
	ds := dbs.DS.Where(goqu.I(front.OrderTable.PK()).Eq(payload.OrderID)).
		Where(goqu.I("$CreatedAt").Eq(payload.CreatedAt)).
		Where(goqu.I("$UserID").Eq(tokUsr.ID))
	if err = db.DsSelectOneTo(&order, ds); err != nil {
		return nil, err
	}

	// prepay after prepaid
	if order.State == front.TOrderStatePrepaid {
		if order.CashPaid != payload.Cash || pointsPaid != order.PointsPaid {
			return nil, InvalidPrepayPayload
		}
		// close transaction_id firstly
		res, err := closeWxOrder(order)
		if err != nil {
			return nil, err
		}

		if res["result_code"] == "SUCCESS" {
			if !lok.CashLok.Lock(tokUsr.ID) {
				return nil, cerr.CashTmpLocked
			}
			defer lok.CashLok.Unlock(tokUsr.ID)

			var flow front.CapitalFlow
			ds = dbs.DS.Where(goqu.I("$OrderID").Eq(order.ID)).Where(goqu.I("$Type").Eq(front.TCapitalFlowPrepay))
			if err = db.DsSelectOneTo(&flow, ds); err != nil && err != reform.ErrNoRows {
				return nil, err
			}
			if uint(-flow.Amount) != order.CashPaid {
				return nil, cerr.InvalidCashPrepaid
			}

			if !lok.PointsLok.Lock(tokUsr.ID) {
				return nil, cerr.CashTmpLocked
			}
			defer lok.PointsLok.Unlock(tokUsr.ID)

			var points front.PointsItem
			ds = dbs.DS.Where(goqu.I("$OrderID").Eq(order.ID)).Where(goqu.I("$Type").Eq(front.TPointsPrepay))
			if err = db.DsSelectOneTo(&points, ds); err != nil && err != reform.ErrNoRows {
				return nil, err
			}
			if uint(-points.Amount)*dbs.config.Order.Point2Cent != order.PointsPaid {
				return nil, cerr.InvalidPointsPrepaid
			}

			if err = db.UpdateColumns(&front.Order{ID: order.ID}, "TransactionId", "TradeState"); err != nil {
				return nil, err
			}

			attach := models.UnifiedOrderAttach{
				PreCashID:   flow.ID,
				CashPaid:    order.CashPaid,
				PrePointsID: points.ID,
				PointsPaid:  order.PointsPaid,
				UserID:      tokUsr.ID,
			}

			args, err := unifiedOrder(&order, &attach)
			if err != nil {
				return nil, err
			}

			if err := tx.Commit(); err != nil {
				return nil, err
			}

			order.TransactionId = ""
			order.TradeState = front.UNKNOWN
			return args, nil
		} else {
			// colse FAIL
			log.WithFields(l.Locate(logrus.Fields{
				"out_trade_no": outTradeNo,
				"err_code":     res["err_code"],
			})).Info("Failed to close")

			switch res["err_code"] {
			case "ORDERPAID":
				return nil, cerr.WxOrderAlreadyPaid
			case "SYSTEMERROR":
				return nil, cerr.WxSystemFailed
			case "ORDERNOTEXIST", "ORDERCLOSED":
				// TODO remove this?
				// check prepaid
				args, err := tx.prepayOrderAfterClosedWxOrder(tokUsr, order, unifiedOrder)
				if err != nil {
					return nil, err
				}
				if args != nil {
					return args, nil
				}

			default:
				return nil, cerr.ApiImplementFailed
			}
		}
	} else if err = PermitOrderState(&order, front.TOrderStatePrepaid); err != nil {
		return nil, cerr.NotNopayState
	}

	if order.PayAmount != payload.Amount {
		return nil, cerr.InvalidPayAmount
	}

	now := time.Now().Unix()
	attach := models.UnifiedOrderAttach{
		UserID: tokUsr.ID,
	}

	if payload.Cash != 0 {
		if !lok.CashLok.Lock(tokUsr.ID) {
			return nil, cerr.CashTmpLocked
		}
		defer lok.CashLok.Unlock(tokUsr.ID)

		var flow front.CapitalFlow
		ds = dbs.DS.Where(goqu.I("$UserID").Eq(tokUsr.ID)).Order(goqu.I("$CreatedAt").Desc().NullsLast())
		if err = db.DsSelectOneTo(&flow, ds); err != nil {
			return nil, err
		}

		if flow.Balance < int(payload.Cash) {
			return nil, cerr.NotEnoughMoney
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
			return nil, err
		}
		attach.CashPaid = payload.Cash
		attach.PreCashID = flow.ID
	}

	if payload.Points != 0 {
		if !lok.PointsLok.Lock(tokUsr.ID) {
			return nil, cerr.CashTmpLocked
		}
		defer lok.PointsLok.Unlock(tokUsr.ID)

		var points front.PointsItem
		ds = dbs.DS.Where(goqu.I("$UserID").Eq(tokUsr.ID)).Order(goqu.I("$CreatedAt").Desc().NullsLast())
		if err = db.DsSelectOneTo(&points, ds); err != nil {
			return nil, err
		}

		if points.Balance < int(payload.Points) {
			return nil, cerr.NotEnoughPoints
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
			return nil, err
		}
		attach.PointsPaid = pointsPaid
		attach.PrePointsID = points.ID
	}

	err = db.UpdateColumns(&front.Order{
		ID:         order.ID,
		CashPaid:   attach.CashPaid,
		PointsPaid: attach.PointsPaid,
		State:      front.TOrderStatePrepaid,
		PrepaidAt:  now,
	}, "CashPaid", "PointsPaid", "State", "PrepaidAt", "TransactionId", "TradeState")
	if err != nil {
		return nil, err
	}

	args, err := unifiedOrder(&order, &attach)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	order.CashPaid = attach.CashPaid
	order.PointsPaid = attach.PointsPaid
	order.State = front.TOrderStatePrepaid
	order.PrepaidAt = now
	order.TransactionId = ""
	order.TradeState = front.UNKNOWN

	return args, nil
}

// no wx pay, need paykey
func (dbs *DbService) PayOrder(tokUsr *models.User, payload *front.OrderPayPayload) (*front.Order, error) {
	pointsPaid := payload.Points * dbs.config.Order.Point2Cent
	if payload.Amount != payload.Cash+pointsPaid || payload.Amount == 0 {
		return nil, cerr.InvalidPayAmount
	}

	if !lok.OrderLok.Lock(payload.OrderID) {
		return nil, cerr.OrderTmpLocked
	}
	defer lok.OrderLok.Unlock(payload.OrderID)

	tx, err := dbs.Tx()
	if err != nil {
		return nil, err
	}
	defer tx.RollbackIfNeeded()
	db := tx.GetDB()

	if err = db.Reload(tokUsr); err != nil {
		return nil, err
	}

	if err = models.ComparePaykey([]byte(tokUsr.Paykey), []byte(payload.Key)); err != nil {
		return nil, err
	}

	var order front.Order
	ds := dbs.DS.Where(goqu.I(front.OrderTable.PK()).Eq(payload.OrderID)).
		Where(goqu.I("$CreatedAt").Eq(payload.CreatedAt)).
		Where(goqu.I("$UserID").Eq(tokUsr.ID))
	if err = db.DsSelectOneTo(&order, ds); err != nil {
		return nil, err
	}
	if err = PermitOrderState(&order, front.TOrderStatePaid); err != nil {
		return nil, cerr.NoWayToPaidState
	}
	if order.PayAmount != payload.Amount {
		return nil, cerr.InvalidPayAmount
	}

	now := time.Now().Unix()

	if payload.Cash != 0 {
		if !lok.CashLok.Lock(tokUsr.ID) {
			return nil, cerr.CashTmpLocked
		}
		defer lok.CashLok.Unlock(tokUsr.ID)

		var flow front.CapitalFlow
		ds = dbs.DS.Where(goqu.I("$UserID").Eq(tokUsr.ID)).Order(goqu.I("$CreatedAt").Desc().NullsLast())
		if err = db.DsSelectOneTo(&flow, ds); err != nil {
			return nil, err
		}

		if flow.Balance < int(payload.Cash) {
			return nil, cerr.NotEnoughMoney
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
			return nil, err
		}
	}

	if payload.Points != 0 {
		if !lok.PointsLok.Lock(tokUsr.ID) {
			return nil, cerr.CashTmpLocked
		}
		defer lok.PointsLok.Unlock(tokUsr.ID)

		var points front.PointsItem
		ds = dbs.DS.Where(goqu.I("$UserID").Eq(tokUsr.ID)).Order(goqu.I("$CreatedAt").Desc().NullsLast())
		if err = db.DsSelectOneTo(&points, ds); err != nil {
			return nil, err
		}

		if points.Balance < int(payload.Points) {
			return nil, cerr.NotEnoughPoints
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
			return nil, err
		}
	}

	order.CashPaid = payload.Cash
	order.PointsPaid = pointsPaid
	order.State = front.TOrderStatePaid
	order.PaidAt = now
	if err = db.UpdateColumns(&order, "CashPaid", "PointsPaid", "State", "PaidAt"); err != nil {
		return nil, err
	}

	return &order, tx.Commit()
}

// TOrderStateEvaled is standalone
func (dbs *DbService) OrderChangeState(
	order *front.Order, tokUsr *models.User, payload *front.OrderChangeStatePayload, paier WxPaier,
) error {

	db := dbs.GetDB()

	// lock order before tx
	ds := dbs.DS.Where(goqu.I(front.OrderTable.PK()).Eq(payload.ID)).
		Where(goqu.I("$CreatedAt").Eq(payload.CreatedAt)).
		Where(goqu.I("$UserID").Eq(tokUsr.ID))
	err := db.DsSelectOneTo(order, ds)
	if err != nil {
		return err
	}

	if err = PermitOrderState(order, payload.State); err != nil {
		return cerr.NoWayToTargetState
	}

	now := time.Now().Unix()
	switch payload.State {
	case front.TOrderStateCanceled:
		switch order.State {
		case front.TOrderStateNopay:
			return db.UpdateColumns(&front.Order{
				ID:         order.ID,
				State:      payload.State,
				CanceledAt: now,
			})
		case front.TOrderStatePrepaid:
			res, err := paier.OrderClose(order)
			if err != nil {
				return err
			}
			if res["result_code"] == "SUCCESS" {
				if !lok.CashLok.Lock(tokUsr.ID) {
					return cerr.CashTmpLocked
				}
				defer lok.CashLok.Unlock(tokUsr.ID)

				var flow front.CapitalFlow
				ds = dbs.DS.Where(goqu.I("$OrderID").Eq(order.ID)).Where(goqu.I("$Type").Eq(front.TCapitalFlowPrepay))
				if err = db.DsSelectOneTo(&flow, ds); err != nil && err != reform.ErrNoRows {
					return err
				}
				if flow.ID != 0 {
					// repaid
					var flow2 front.CapitalFlow
					ds = dbs.DS.Where(goqu.I("$OrderID").Eq(order.ID)).Where(goqu.I("$Type").Eq(front.TCapitalFlowPrepayBack))
					if err = db.DsSelectOneTo(&flow2, ds); err != nil && err != reform.ErrNoRows {
						return err
					}
					if flow2.ID == 0 {
						// not refund yet
						var flow1 front.CapitalFlow
						ds = dbs.DS.Where(goqu.I("$UserID").Eq(tokUsr.ID)).Order(goqu.I("$CreatedAt").Desc().NullsLast())
						if err = db.DsSelectOneTo(&flow1, ds); err != nil {
							return err
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
							return err
						}
					}
				}

				if !lok.PointsLok.Lock(tokUsr.ID) {
					return cerr.CashTmpLocked
				}
				defer lok.PointsLok.Unlock(tokUsr.ID)

				var points front.PointsItem
				ds = dbs.DS.Where(goqu.I("$OrderID").Eq(order.ID)).Where(goqu.I("$Type").Eq(front.TPointsPrepay))
				if err = db.DsSelectOneTo(&points, ds); err != nil && err != reform.ErrNoRows {
					return err
				}
				if points.ID != 0 {
					// repaid
					var points2 front.PointsItem
					ds = dbs.DS.Where(goqu.I("$OrderID").Eq(order.ID)).Where(goqu.I("$Type").Eq(front.TPointsPrepayBack))
					if err = db.DsSelectOneTo(&points2, ds); err != nil && err != reform.ErrNoRows {
						return err
					}
					if points2.ID == 0 {
						// not refund yet
						var points1 front.PointsItem
						ds = dbs.DS.Where(goqu.I("$UserID").Eq(tokUsr.ID)).Order(goqu.I("$CreatedAt").Desc().NullsLast())
						if err = db.DsSelectOneTo(&points1, ds); err != nil {
							return err
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
							return err
						}
					}
				}

				err = db.UpdateColumns(&front.Order{
					ID:         order.ID,
					CanceledAt: now,
					State:      front.TOrderStateCanceled,
				}, "CanceledAt", "State", "TransactionId", "TradeState")
				if err != nil {
					return err
				}

				res, err := closeWxOrder(order)
				if err != nil {
					return err
				}

				if res["result_code"] != "SUCCESS" {
					switch res["err_code"] {
					case "ORDERPAID":
						return cerr.WxOrderAlreadyPaid
					case "SYSTEMERROR":
						return cerr.WxSystemFailed
					case "SIGNERROR", "REQUIRE_POST_METHOD", "XML_FORMAT_ERROR":
						return cerr.ApiImplementFailed
					}
				}
				return nil
			}
		case front.TOrderStatePaid, front.TOrderStatePicking:
			if order.CashRefund != 0 {
				if !lok.CashLok.Lock(tokUsr.ID) {
					return cerr.CashTmpLocked
				}
				defer lok.CashLok.Unlock(tokUsr.ID)

				var flow front.CapitalFlow
				ds = dbs.DS.Where(goqu.I("$OrderID").Eq(order.ID)).Where(goqu.I("$Type").Eq(front.TCapitalFlowRefund))
				if err = db.DsSelectOneTo(&flow, ds); err != nil && err != reform.ErrNoRows {
					return err
				}
				if flow.ID == 0 {
					// not refund yet
					var flow1 front.CapitalFlow
					ds = dbs.DS.Where(goqu.I("$UserID").Eq(tokUsr.ID)).Order(goqu.I("$CreatedAt").Desc().NullsLast())
					if err = db.DsSelectOneTo(&flow1, ds); err != nil {
						return err
					}
					// refund
					err = db.Save(&front.CapitalFlow{
						UserID:    tokUsr.ID,
						CreatedAt: now,
						Type:      front.TCapitalFlowRefund, // Trade type
						Amount:    int(order.CashRefund),
						Balance:   flow1.Balance + int(order.CashRefund),
						OrderID:   order.ID,
					})
					if err != nil {
						return err
					}
				}
			}

			if order.PointsRefund != 0 {
				if !lok.PointsLok.Lock(tokUsr.ID) {
					return cerr.CashTmpLocked
				}
				defer lok.PointsLok.Unlock(tokUsr.ID)

				var points front.PointsItem
				ds = dbs.DS.Where(goqu.I("$OrderID").Eq(order.ID)).Where(goqu.I("$Type").Eq(front.TPointsRefund))
				if err = db.DsSelectOneTo(&points, ds); err != nil && err != reform.ErrNoRows {
					return err
				}
				if points.ID == 0 {
					// not refund yet
					var points1 front.PointsItem
					ds = dbs.DS.Where(goqu.I("$UserID").Eq(tokUsr.ID)).Order(goqu.I("$CreatedAt").Desc().NullsLast())
					if err = db.DsSelectOneTo(&points1, ds); err != nil {
						return err
					}
					// refund
					err = db.Save(&front.PointsItem{
						UserID:    tokUsr.ID,
						CreatedAt: now,
						Type:      front.TPointsRefund, // Trade type
						Amount:    int(order.PointsRefund),
						Balance:   points1.Balance + int(order.PointsRefund),
						OrderID:   order.ID,
					})
					if err != nil {
						return err
					}
				}
			}

			if order.WxRefund != 0 {
				res, err := paier.OrderRefund(order, strconv.Itoa(int(tokUsr.ID)))
				if err != nil {
					return err
				}

				if res["result_code"] == "SUCCESS" {
					err = db.UpdateColumns(&front.Order{
						ID:         order.ID,
						CanceledAt: now,
						State:      front.TOrderStateCanceled,
					}, "CanceledAt", "State")
					if err != nil {
						return err
					}
				} else {
					switch res["err_code"] {
					case "SYSTEMERROR":
						return cerr.WxSystemFailed
					case "USER_ACCOUNT_ABNORMAL":
						return cerr.WxUserAbnormal
						//					case "INVALID_TRANSACTIONID":
						//						return cerr.WxSystemFailed
						//					case "SIGNERROR", "PARAM_ERROR", "APPID_NOT_EXIST", "MCHID_NOT_EXIST",
						//						"APPID_MCHID_NOT_MATCH", "REQUIRE_POST_METHOD", "XML_FORMAT_ERROR":
						//						return cerr.ApiImplementFailed
					}
					return cerr.ApiImplementFailed
				}
			}

			return nil
		} // to cancel end

	case front.TOrderStateCompleted, front.TOrderStateReturnStarted:
		if dbs.IsOrderCompleted(order) {
			return cerr.OrderCompleteTimeout
		}
		return db.UpdateColumns(&front.Order{
			ID:          order.ID,
			State:       payload.State,
			CompletedAt: now,
		})
		return nil
	}

	return cerr.Forbidden
}

func (dbs *DbService) MgrOrderState() error {}

func (dbs *DbService) GetBareOrder(id uint, at int64) (*front.Order, error) {
	var order front.Order
	ds := dbs.DS.Where(goqu.I("$CreatedAt").Eq(at), goqu.I(front.OrderTable.PK()).Eq(id))
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
