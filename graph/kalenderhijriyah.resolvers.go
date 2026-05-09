package graph

import (
	"context"

	"github.com/lisvindanuu/indonesiaql/graph/model"
)

func (r *queryResolver) KalenderHijriyah(ctx context.Context, tanggal string) (*model.KalenderHijriyahResult, error) {
	res, err := r.kalenderHijriyahSvc.Convert(tanggal)
	if err != nil {
		return nil, err
	}
	return &model.KalenderHijriyahResult{
		TanggalMasehi:   res.TanggalMasehi,
		TanggalHijriyah: res.TanggalHijriyah,
		Hari:            res.Hari,
		HariArab:        res.HariArab,
		Bulan:           res.Bulan,
		BulanArab:       res.BulanArab,
		Tahun:           res.Tahun,
	}, nil
}
