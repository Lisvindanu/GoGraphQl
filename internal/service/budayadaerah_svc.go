package service

import (
	"context"
	"strings"

	"github.com/lisvindanuu/indonesiaql/internal/staticdata"
)

type BudayaDaerahResult struct {
	Provinsi    string
	RumahAdat   string
	PakaianAdat string
	TariDaerah  string
}

type BudayaDaerahService struct{}

func NewBudayaDaerahService() *BudayaDaerahService { return &BudayaDaerahService{} }

func (s *BudayaDaerahService) ListAll(_ context.Context) []BudayaDaerahResult {
	result := make([]BudayaDaerahResult, len(staticdata.BudayaDaerah))
	for i, b := range staticdata.BudayaDaerah {
		result[i] = BudayaDaerahResult{
			Provinsi:    b.Provinsi,
			RumahAdat:   b.RumahAdat,
			PakaianAdat: b.PakaianAdat,
			TariDaerah:  b.TariDaerah,
		}
	}
	return result
}

func (s *BudayaDaerahService) Cari(_ context.Context, provinsi string) []BudayaDaerahResult {
	lower := strings.ToLower(provinsi)
	var result []BudayaDaerahResult
	for _, b := range staticdata.BudayaDaerah {
		if strings.Contains(strings.ToLower(b.Provinsi), lower) {
			result = append(result, BudayaDaerahResult{
				Provinsi:    b.Provinsi,
				RumahAdat:   b.RumahAdat,
				PakaianAdat: b.PakaianAdat,
				TariDaerah:  b.TariDaerah,
			})
		}
	}
	return result
}
