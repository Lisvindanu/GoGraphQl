package service

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/lisvindanuu/indonesiaql/internal/staticdata"
)

type KodeBankResult struct {
	Kode string
	Nama string
}

type KodeBankService struct{}

func NewKodeBankService() *KodeBankService { return &KodeBankService{} }

func (s *KodeBankService) GetByKode(_ context.Context, kode string) (*KodeBankResult, error) {
	kode = strings.TrimSpace(kode)
	if len(kode) < 3 {
		kode = fmt.Sprintf("%03s", kode)
	}
	if info, ok := staticdata.KodeBank[kode]; ok {
		return &KodeBankResult{Kode: kode, Nama: info.Nama}, nil
	}
	return nil, fmt.Errorf("kode bank '%s' tidak ditemukan", kode)
}

func (s *KodeBankService) ListAll(_ context.Context) []KodeBankResult {
	result := make([]KodeBankResult, 0, len(staticdata.KodeBank))
	for kode, info := range staticdata.KodeBank {
		result = append(result, KodeBankResult{Kode: kode, Nama: info.Nama})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Kode < result[j].Kode })
	return result
}
