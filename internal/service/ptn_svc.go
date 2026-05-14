package service

import (
	"context"
	"strings"

	"github.com/lisvindanuu/indonesiaql/internal/staticdata"
)

type PTNResult struct {
	Nama       string
	Singkatan  string
	Kota       string
	Provinsi   string
	Akreditasi string
	Jenis      string
}

type PTNService struct{}

func NewPTNService() *PTNService { return &PTNService{} }

func (s *PTNService) ListAll(_ context.Context) []PTNResult {
	result := make([]PTNResult, len(staticdata.DaftarPTN))
	for i, p := range staticdata.DaftarPTN {
		result[i] = PTNResult{
			Nama:       p.Nama,
			Singkatan:  p.Singkatan,
			Kota:       p.Kota,
			Provinsi:   p.Provinsi,
			Akreditasi: p.Akreditasi,
			Jenis:      p.Jenis,
		}
	}
	return result
}

func (s *PTNService) Cari(_ context.Context, query string) []PTNResult {
	lower := strings.ToLower(query)
	var result []PTNResult
	for _, p := range staticdata.DaftarPTN {
		if strings.Contains(strings.ToLower(p.Nama), lower) ||
			strings.Contains(strings.ToLower(p.Singkatan), lower) ||
			strings.Contains(strings.ToLower(p.Kota), lower) ||
			strings.Contains(strings.ToLower(p.Provinsi), lower) {
			result = append(result, PTNResult{
				Nama:       p.Nama,
				Singkatan:  p.Singkatan,
				Kota:       p.Kota,
				Provinsi:   p.Provinsi,
				Akreditasi: p.Akreditasi,
				Jenis:      p.Jenis,
			})
		}
	}
	return result
}
