package dbsrv

import (
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/reform"
)

func (dbs *DbService) ProductsFillResponse(products ...reform.Struct) (*front.ProductsResponse, error) {
	if len(products) == 0 {
		return nil, reform.ErrNoRows
	}

	var ids []interface{}
	for _, p := range products {
		ids = append(ids, p.(*front.Product).ID)
	}
	skus, err := dbs.GetDB().FindAllFrom(front.SkuTable, "$ProductID", ids...)
	if err != nil {
		return nil, err
	}

	var attrIds []reform.Struct
	if len(skus) != 0 {
		ids = nil
		for _, sku := range skus {
			ids = append(ids, sku.(*front.Sku).ID)
		}
		attrIds, err = dbs.GetDB().FindAllFrom(front.ProductAttrIdTable, "$SkuID", ids...)
		if err != nil {
			return nil, err
		}
	}

	return &front.ProductsResponse{
		Products: products,
		Skus:     skus,
		Attrs:    attrIds,
	}, nil
}
