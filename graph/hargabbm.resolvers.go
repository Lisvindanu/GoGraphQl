package graph

import (
	"context"

	"github.com/lisvindanuu/indonesiaql/graph/model"
)

func (r *queryResolver) HargaBbm(ctx context.Context) ([]*model.HargaBBMItem, error) {
	list := r.hargaBBMSvc.GetAll()
	result := make([]*model.HargaBBMItem, len(list))
	for i, item := range list {
		item := item
		result[i] = &model.HargaBBMItem{
			Nama:   item.Nama,
			Harga:  item.Harga,
			Satuan: item.Satuan,
			Jenis:  item.Jenis,
		}
	}
	return result, nil
}
