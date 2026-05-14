package service

import (
	"context"
	"strings"

	"github.com/lisvindanuu/indonesiaql/internal/staticdata"
)

type PresidenResult struct {
	Urutan        int
	Nama          string
	WakilPresiden string
	MulaiJabatan  int
	AkhirJabatan  int
}

type PresidenService struct{}

func NewPresidenService() *PresidenService { return &PresidenService{} }

func (s *PresidenService) ListAll(_ context.Context) []PresidenResult {
	result := make([]PresidenResult, len(staticdata.DaftarPresiden))
	for i, p := range staticdata.DaftarPresiden {
		result[i] = PresidenResult{
			Urutan:        p.Urutan,
			Nama:          p.Nama,
			WakilPresiden: p.WakilPresiden,
			MulaiJabatan:  p.MulaiJabatan,
			AkhirJabatan:  p.AkhirJabatan,
		}
	}
	return result
}

func (s *PresidenService) Cari(_ context.Context, nama string) []PresidenResult {
	lower := strings.ToLower(nama)
	var result []PresidenResult
	for _, p := range staticdata.DaftarPresiden {
		if strings.Contains(strings.ToLower(p.Nama), lower) ||
			strings.Contains(strings.ToLower(p.WakilPresiden), lower) {
			result = append(result, PresidenResult{
				Urutan:        p.Urutan,
				Nama:          p.Nama,
				WakilPresiden: p.WakilPresiden,
				MulaiJabatan:  p.MulaiJabatan,
				AkhirJabatan:  p.AkhirJabatan,
			})
		}
	}
	return result
}
