package server

import (
	"net/http"
	"strconv"

	"github.com/empirefox/esecend/cerr"
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/esecend/utils"
	"github.com/empirefox/reform"
	"github.com/gin-gonic/gin"
)

func (s *Server) GetProductsBundle(c *gin.Context) {
	bundle := make(map[string][]reform.Struct)
	var products []reform.Struct
	var err error

	matrix, _ := utils.ParseMatrixPath(c.Param("matrix"))
	for path, values := range matrix {
		if len(values) > 0 && values[0] != "" {
			searcher := s.ProductResource.NewSearcherFromRaw(values[0])
			seg, e := searcher.FindMany()
			if e == nil {
				bundle[path] = seg
				products = append(products, seg...)
			} else if err == nil || err == reform.ErrNoRows {
				err = e
			}
		}
	}

	// no products
	if err != nil && len(products) == 0 {
		ResponseArray(c, nil, err)
		return
	}

	data, err := s.DB.ProductsFillResponse(products...)
	ResponseArray(c, &front.ProductsBundleResponse{
		Bundle: bundle,
		Skus:   data.Skus,
		Attrs:  data.Attrs,
	}, err)
}

// Get /:table/ls?
//		&q(query)=bmw
//		&st(start)=100&sz(size)=20&tl=1
//		&sp(scope)=2016style+white
//		&ft(filter)=Price:gteq:10+Price:lteq:20+Discount:true
//		&ob(order)=Price:desc
func (s *Server) GetProducts(c *gin.Context) {
	searcher := s.ProductResource.NewSearcher(c)
	products, err := searcher.FindMany()

	var data *front.ProductsResponse
	if err == nil {
		data, err = s.DB.ProductsFillResponse(products...)
	}

	ResponseArray(c, data, err)
}

func (s *Server) GetProduct(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		front.NewCodev(cerr.InvalidUrlParam).Abort(c, http.StatusBadRequest)
		return
	}

	product, err := s.DB.GetDB().FindByPrimaryKeyFrom(front.ProductTable, id)

	var data *front.ProductsResponse
	if err == nil {
		data, err = s.DB.ProductsFillResponse(product)
	}

	ResponseArray(c, data, err)
}
