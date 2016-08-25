package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/empirefox/esecend/admin"
	"github.com/empirefox/esecend/front"
	"github.com/stretchr/testify/require"
)

func TestOrder(t *testing.T) {
	now := time.Now().Unix()
	require := require.New(t)

	// 1. checkout
	req, _ := http.NewRequest("POST", "/checkout", strings.NewReader(`
{
  "Contact": "gang",
  "Phone": "13122223333",
  "DeliverAddress": "sichuan chengdu",
  "InvoiceTo": "gangge",
  "InvoiceToCom": false,
  "Remark": "ASAP",
  "Total": 50000,
  "DeliverFee": 0,
  "Items": [{
    "SkuID": 1,
    "Quantity": 10,
    "Attrs": [1,3]
  }]
}`))
	res := httptest.NewRecorder()
	server.ServeHTTP(res, req)

	var order front.Order
	err := json.NewDecoder(res.Body).Decode(&order)
	require.Nil(err)
	require.Equal(200, res.Code)

	// 2. pay without enough money
	orderPayPayload, _ := json.Marshal(&front.OrderPayPayload{
		Key:     "123456",
		Amount:  50000,
		OrderID: order.ID,
		Cash:    50000,
	})
	req, _ = http.NewRequest("POST", "/order_pay", bytes.NewReader(orderPayPayload))
	res = httptest.NewRecorder()
	server.ServeHTTP(res, req)
	require.NotEqual(200, res.Code)

	// 3. recharge money
	db := server.DB.GetDB()
	err = db.Insert(&front.CapitalFlow{
		UserID:    1,
		CreatedAt: now,
		Type:      front.TCapitalFlowRecharge,
		Amount:    60000,
		Balance:   60000,
	})
	require.Nil(err)

	// 4. pay with enough money
	req, _ = http.NewRequest("POST", "/order_pay", bytes.NewReader(orderPayPayload))
	res = httptest.NewRecorder()
	server.ServeHTTP(res, req)
	// check state
	err = json.NewDecoder(res.Body).Decode(&order)
	require.Nil(err)
	require.Equal(front.TOrderStatePaid, order.State)
	require.Equal(200, res.Code)

	// 5. admin state to picking
	adminToken, err := jwt.NewWithClaims(jwt.GetSigningMethod(server.Config.Security.AdminSignType), &admin.Claims{
		AdminId: 111,
		UserId:  1,
		OrderID: order.ID,
		State:   front.TOrderStatePicking,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now,
			ExpiresAt: now + 10,
		},
	}).SignedString([]byte(server.Config.Security.AdminKey))
	require.Nil(err)

	req, _ = http.NewRequest("GET", "/admin/order_state", nil)
	req.Header.Set("Authorization", "BEARER "+adminToken)
	res = httptest.NewRecorder()
	server.ServeHTTP(res, req)
	// check state
	err = json.NewDecoder(res.Body).Decode(&order)
	require.Nil(err)
	require.Equal(front.TOrderStatePicking, order.State)
	require.Equal(200, res.Code)

	// 6. admin state to deliver
	adminToken, err = jwt.NewWithClaims(jwt.GetSigningMethod(server.Config.Security.AdminSignType), &admin.Claims{
		AdminId: 111,
		UserId:  1,
		OrderID: order.ID,
		State:   front.TOrderStateDelivered,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now,
			ExpiresAt: now + 10,
		},
		DeliverCom: "shunfeng",
		DeliverNo:  "924383289093",
	}).SignedString([]byte(server.Config.Security.AdminKey))
	require.Nil(err)

	req, _ = http.NewRequest("GET", "/admin/order_state", nil)
	req.Header.Set("Authorization", "BEARER "+adminToken)
	res = httptest.NewRecorder()
	server.ServeHTTP(res, req)
	// check state
	err = json.NewDecoder(res.Body).Decode(&order)
	require.Nil(err)
	require.Equal(front.TOrderStateDelivered, order.State)
	require.Equal(200, res.Code)

	// 7. user complete the order
	orderChangeStatePayload, _ := json.Marshal(front.OrderChangeStatePayload{
		ID:    order.ID,
		State: front.TOrderStateCompleted,
	})
	req, _ = http.NewRequest("POST", "/order_state", bytes.NewReader(orderChangeStatePayload))
	res = httptest.NewRecorder()
	server.ServeHTTP(res, req)
	require.Equal(200, res.Code)

	// 8. user eval the order
	evalItemPayload, _ := json.Marshal(front.EvalItem{
		Eval:        "good",
		RateStar:    1,
		RateFit:     2,
		RateServe:   3,
		RateDeliver: 4,
	})
	req, _ = http.NewRequest("POST", fmt.Sprintf("/eval/%d", order.ID), bytes.NewReader(evalItemPayload))
	res = httptest.NewRecorder()
	server.ServeHTTP(res, req)
	require.Equal(200, res.Code)

	// 9. veriry order from database
	req, _ = http.NewRequest("GET", "/orders", nil)
	res = httptest.NewRecorder()
	server.ServeHTTP(res, req)
	var orders []front.Order
	err = json.NewDecoder(res.Body).Decode(&orders)
	require.Nil(err)
	require.Equal(200, res.Code)

	require.Equal(1, len(orders))
	require.Equal(front.TOrderStateEvaled, orders[0].State)
	require.Equal("924383289093", orders[0].DeliverNo)
}
