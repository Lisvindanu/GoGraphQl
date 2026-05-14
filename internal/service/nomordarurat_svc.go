package service

import (
	"context"
	"strings"

	"github.com/lisvindanuu/indonesiaql/internal/staticdata"
)

type NomorDaruratResult struct {
	Nomor    string
	Layanan  string
	Kategori string
}

type NomorDaruratService struct{}

func NewNomorDaruratService() *NomorDaruratService { return &NomorDaruratService{} }

func (s *NomorDaruratService) ListAll(_ context.Context) []NomorDaruratResult {
	result := make([]NomorDaruratResult, len(staticdata.NomorDarurat))
	for i, n := range staticdata.NomorDarurat {
		result[i] = NomorDaruratResult{
			Nomor:    n.Nomor,
			Layanan:  n.Layanan,
			Kategori: n.Kategori,
		}
	}
	return result
}

func (s *NomorDaruratService) Cari(_ context.Context, query string) []NomorDaruratResult {
	lower := strings.ToLower(query)
	var result []NomorDaruratResult
	for _, n := range staticdata.NomorDarurat {
		if strings.Contains(strings.ToLower(n.Layanan), lower) ||
			strings.Contains(strings.ToLower(n.Kategori), lower) ||
			strings.Contains(strings.ToLower(n.Nomor), lower) {
			result = append(result, NomorDaruratResult{
				Nomor:    n.Nomor,
				Layanan:  n.Layanan,
				Kategori: n.Kategori,
			})
		}
	}
	return result
}
