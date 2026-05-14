package service

import (
	"context"
	"strings"

	"github.com/lisvindanuu/indonesiaql/internal/staticdata"
)

type StasiunKeretaResult struct {
	Nama  string
	Kode  string
	Jalur string
	Tipe  string
	Kota  string
}

type StasiunKeretaService struct{}

func NewStasiunKeretaService() *StasiunKeretaService { return &StasiunKeretaService{} }

func (s *StasiunKeretaService) ListAll(_ context.Context) []StasiunKeretaResult {
	result := make([]StasiunKeretaResult, len(staticdata.DaftarStasiunKereta))
	for i, st := range staticdata.DaftarStasiunKereta {
		result[i] = StasiunKeretaResult{
			Nama:  st.Nama,
			Kode:  st.Kode,
			Jalur: st.Jalur,
			Tipe:  st.Tipe,
			Kota:  st.Kota,
		}
	}
	return result
}

func (s *StasiunKeretaService) Cari(_ context.Context, nama string) []StasiunKeretaResult {
	lower := strings.ToLower(nama)
	var result []StasiunKeretaResult
	for _, st := range staticdata.DaftarStasiunKereta {
		if strings.Contains(strings.ToLower(st.Nama), lower) ||
			strings.Contains(strings.ToLower(st.Kota), lower) ||
			strings.Contains(strings.ToLower(st.Kode), lower) {
			result = append(result, StasiunKeretaResult{
				Nama:  st.Nama,
				Kode:  st.Kode,
				Jalur: st.Jalur,
				Tipe:  st.Tipe,
				Kota:  st.Kota,
			})
		}
	}
	return result
}
