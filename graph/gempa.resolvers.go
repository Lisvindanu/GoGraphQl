package graph

import (
	"context"

	"github.com/lisvindanuu/indonesiaql/graph/model"
)

func (r *queryResolver) GempaTerbaru(ctx context.Context) (*model.GempaItem, error) {
	item, err := r.gempaSvc.GetTerbaru()
	if err != nil {
		return nil, err
	}
	return &model.GempaItem{
		Tanggal:   item.Tanggal,
		Jam:       item.Jam,
		Magnitude: item.Magnitude,
		Kedalaman: item.Kedalaman,
		Lintang:   item.Lintang,
		Bujur:     item.Bujur,
		Wilayah:   item.Wilayah,
		Potensi:   item.Potensi,
		Dirasakan: item.Dirasakan,
	}, nil
}

func (r *queryResolver) GempaList(ctx context.Context) ([]*model.GempaItem, error) {
	list, err := r.gempaSvc.GetList()
	if err != nil {
		return nil, err
	}
	result := make([]*model.GempaItem, len(list))
	for i, item := range list {
		item := item
		result[i] = &model.GempaItem{
			Tanggal:   item.Tanggal,
			Jam:       item.Jam,
			Magnitude: item.Magnitude,
			Kedalaman: item.Kedalaman,
			Lintang:   item.Lintang,
			Bujur:     item.Bujur,
			Wilayah:   item.Wilayah,
			Potensi:   item.Potensi,
			Dirasakan: item.Dirasakan,
		}
	}
	return result, nil
}
