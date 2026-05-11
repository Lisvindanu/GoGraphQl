package service

import (
	"context"
	"strings"

	"github.com/lisvindanuu/indonesiaql/internal/staticdata"
)

type GunungBerapiResult struct {
	Nama                    string
	Bentuk                  string
	TinggiMeter             string
	EstimasiLetusanTerakhir string
	Geolokasi               string
}

type GunungBerapiService struct{}

func NewGunungBerapiService() *GunungBerapiService { return &GunungBerapiService{} }

func (s *GunungBerapiService) ListAll(_ context.Context) []GunungBerapiResult {
	result := make([]GunungBerapiResult, len(staticdata.GunungBerapi))
	for i, g := range staticdata.GunungBerapi {
		result[i] = GunungBerapiResult{
			Nama:                    g.Nama,
			Bentuk:                  g.Bentuk,
			TinggiMeter:             g.TinggiMeter,
			EstimasiLetusanTerakhir: g.EstimasiLetusanTerakhir,
			Geolokasi:               g.Geolokasi,
		}
	}
	return result
}

func (s *GunungBerapiService) Cari(_ context.Context, nama string) []GunungBerapiResult {
	lower := strings.ToLower(nama)
	var result []GunungBerapiResult
	for _, g := range staticdata.GunungBerapi {
		if strings.Contains(strings.ToLower(g.Nama), lower) {
			result = append(result, GunungBerapiResult{
				Nama:                    g.Nama,
				Bentuk:                  g.Bentuk,
				TinggiMeter:             g.TinggiMeter,
				EstimasiLetusanTerakhir: g.EstimasiLetusanTerakhir,
				Geolokasi:               g.Geolokasi,
			})
		}
	}
	return result
}
