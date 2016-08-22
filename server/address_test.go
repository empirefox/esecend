package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/empirefox/esecend/front"
	"github.com/stretchr/testify/require"
)

func TestAddr(t *testing.T) {
	require := require.New(t)
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/addr", strings.NewReader(`
{
  "Contact": "gang",
  "Phone": "13011112222",
  "Province": "sichuan",
  "City": "chengdu",
  "District": "jinniu",
  "House": "wandan",
  "Pos": 1
}`))
	server.ServeHTTP(res, req)

	require.Equal(200, res.Code)
	var addr front.Address
	err := json.NewDecoder(res.Body).Decode(&addr)
	require.Nil(err)
	require.NotEqual(0, addr.ID)

	req, _ = http.NewRequest("GET", "/addrs", nil)
	res = httptest.NewRecorder()
	server.ServeHTTP(res, req)
	require.Equal(200, res.Code)
	var addrs []front.Address
	err = json.NewDecoder(res.Body).Decode(&addrs)
	require.Nil(err)
	require.Equal(1, len(addrs))

	req, _ = http.NewRequest("DELETE", fmt.Sprintf("/addr/%d", addr.ID), nil)
	res = httptest.NewRecorder()
	server.ServeHTTP(res, req)
	require.Equal(200, res.Code)

	req, _ = http.NewRequest("GET", "/addrs", nil)
	res = httptest.NewRecorder()
	server.ServeHTTP(res, req)
	require.Equal(200, res.Code)
	addrs = nil
	err = json.NewDecoder(res.Body).Decode(&addrs)
	require.Nil(err)
	require.Equal(0, len(addrs))
}
