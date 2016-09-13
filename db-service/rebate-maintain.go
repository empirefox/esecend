package dbsrv

import (
	"time"

	"github.com/empirefox/esecend/front"
)

func (s *DbService) RebateMaintain() error {
	db := s.GetDB()

	irebates, err := db.FindAllFrom(front.UserCashRebateTable, "$DoneAt", 0)
	if err != nil {
		return err
	}
	if len(irebates) == 0 {
		return nil
	}

	rebates := make(map[uint]*front.UserCashRebate)
	var ids []interface{}
	for _, irebate := range irebates {
		rebate := irebate.(*front.UserCashRebate)
		rebates[rebate.ID] = rebate
		ids = append(ids, rebate.ID)
	}

	allItems, err := db.FindAllFrom(front.UserCashRebateItemTable, "$RebateID", ids...)
	if err != nil {
		return err
	}

	for _, iitem := range allItems {
		item := iitem.(*front.UserCashRebateItem)
		rebates[item.RebateID].Items = append(rebates[item.RebateID].Items, item)
	}

	now := time.Now().Unix()
	for _, rebate := range rebates {
		createdAt := time.Unix(rebate.CreatedAt, 0)
		litems := len(rebate.Items)
		next := createdAt.AddDate(0, litems+1, 0).Unix()
		if now >= next {
			item := front.UserCashRebateItem{
				RebateID:  rebate.ID,
				CreatedAt: now,
			}
			if litems+1 == int(rebate.Stages) {
				item.Amount = rebate.Amount - (rebate.Amount/rebate.Stages)*(rebate.Stages-1)
				rebate.DoneAt = now
				err = db.UpdateColumns(rebate, "DoneAt")
				if err != nil {
					return err
				}
			} else {
				item.Amount = rebate.Amount / rebate.Stages
			}
			err = db.Insert(&item)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
