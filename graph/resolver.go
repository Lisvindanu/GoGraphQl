package graph

import (
	"github.com/lisvindanuu/indonesiaql/internal/service"
)

// Resolver adalah root resolver yang menyimpan semua service dependencies.
type Resolver struct {
	wilayahSvc   *service.WilayahService
	hariLiburSvc *service.HariLiburService
	cuacaSvc     *service.CuacaService
	kursSvc      *service.KursService
	nikSvc       *service.NIKService
}

func NewResolver(
	wilayahSvc *service.WilayahService,
	hariLiburSvc *service.HariLiburService,
	cuacaSvc *service.CuacaService,
	kursSvc *service.KursService,
	nikSvc *service.NIKService,
) *Resolver {
	return &Resolver{
		wilayahSvc:   wilayahSvc,
		hariLiburSvc: hariLiburSvc,
		cuacaSvc:     cuacaSvc,
		kursSvc:      kursSvc,
		nikSvc:       nikSvc,
	}
}
