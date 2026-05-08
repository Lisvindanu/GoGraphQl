package graph

import (
	"github.com/lisvindanuu/indonesiaql/internal/service"
)

// Resolver adalah root resolver yang menyimpan semua service dependencies.
type Resolver struct {
	wilayahSvc     *service.WilayahService
	hariLiburSvc   *service.HariLiburService
	cuacaSvc       *service.CuacaService
	kursSvc        *service.KursService
	nikSvc         *service.NIKService
	kodeBankSvc    *service.KodeBankService
	platNomorSvc   *service.PlatNomorService
	waktuSholatSvc *service.WaktuSholatService
}

func NewResolver(
	wilayahSvc *service.WilayahService,
	hariLiburSvc *service.HariLiburService,
	cuacaSvc *service.CuacaService,
	kursSvc *service.KursService,
	nikSvc *service.NIKService,
	kodeBankSvc *service.KodeBankService,
	platNomorSvc *service.PlatNomorService,
	waktuSholatSvc *service.WaktuSholatService,
) *Resolver {
	return &Resolver{
		wilayahSvc:     wilayahSvc,
		hariLiburSvc:   hariLiburSvc,
		cuacaSvc:       cuacaSvc,
		kursSvc:        kursSvc,
		nikSvc:         nikSvc,
		kodeBankSvc:    kodeBankSvc,
		platNomorSvc:   platNomorSvc,
		waktuSholatSvc: waktuSholatSvc,
	}
}
