package graph

import (
	"context"

	"github.com/lisvindanuu/indonesiaql/graph/model"
)

func (r *queryResolver) KodePos(ctx context.Context, kode string) ([]*model.KodePosResult, error) {
	list, err := r.kodePosSvc.Search(ctx, kode)
	if err != nil {
		return nil, err
	}
	result := make([]*model.KodePosResult, len(list))
	for i, item := range list {
		item := item
		result[i] = &model.KodePosResult{
			KodePos:   item.KodePos,
			Kelurahan: item.Kelurahan,
			Kecamatan: item.Kecamatan,
			Kota:      item.Kota,
			Provinsi:  item.Provinsi,
		}
	}
	return result, nil
}
