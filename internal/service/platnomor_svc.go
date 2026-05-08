package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/lisvindanuu/indonesiaql/internal/staticdata"
)

type PlatNomorResult struct {
	Kode     string
	Wilayah  string
	Provinsi string
}

type PlatNomorService struct{}

func NewPlatNomorService() *PlatNomorService { return &PlatNomorService{} }

func (s *PlatNomorService) GetByKode(_ context.Context, kode string) (*PlatNomorResult, error) {
	kode = strings.TrimSpace(strings.ToUpper(kode))
	if info, ok := staticdata.PlatNomor[kode]; ok {
		return &PlatNomorResult{Kode: kode, Wilayah: info.Wilayah, Provinsi: info.Provinsi}, nil
	}
	return nil, fmt.Errorf("kode plat '%s' tidak ditemukan", kode)
}
