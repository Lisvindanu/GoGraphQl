package service

import (
	"context"
	"strings"

	"github.com/lisvindanuu/indonesiaql/internal/staticdata"
)

type PahlawanResult struct {
	Nama          string
	TahunLahir    int
	TahunWafat    int
	Deskripsi     string
	TahunDiangkat int
}

type PahlawanService struct{}

func NewPahlawanService() *PahlawanService { return &PahlawanService{} }

func (s *PahlawanService) ListAll(_ context.Context) []PahlawanResult {
	result := make([]PahlawanResult, len(staticdata.PahlawanNasional))
	for i, p := range staticdata.PahlawanNasional {
		result[i] = PahlawanResult{
			Nama:          p.Nama,
			TahunLahir:    p.TahunLahir,
			TahunWafat:    p.TahunWafat,
			Deskripsi:     p.Deskripsi,
			TahunDiangkat: p.TahunDiangkat,
		}
	}
	return result
}

func (s *PahlawanService) Cari(_ context.Context, nama string) []PahlawanResult {
	lower := strings.ToLower(nama)
	var result []PahlawanResult
	for _, p := range staticdata.PahlawanNasional {
		if strings.Contains(strings.ToLower(p.Nama), lower) {
			result = append(result, PahlawanResult{
				Nama:          p.Nama,
				TahunLahir:    p.TahunLahir,
				TahunWafat:    p.TahunWafat,
				Deskripsi:     p.Deskripsi,
				TahunDiangkat: p.TahunDiangkat,
			})
		}
	}
	return result
}
