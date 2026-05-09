package graph

import (
	"github.com/lisvindanuu/indonesiaql/internal/service"
)

// Resolver adalah root resolver yang menyimpan semua service dependencies.
type Resolver struct {
	wilayahSvc          *service.WilayahService
	hariLiburSvc        *service.HariLiburService
	cuacaSvc            *service.CuacaService
	kursSvc             *service.KursService
	nikSvc              *service.NIKService
	kodeBankSvc         *service.KodeBankService
	platNomorSvc        *service.PlatNomorService
	waktuSholatSvc      *service.WaktuSholatService
	gempaSvc            *service.GempaService
	kodePosSvc          *service.KodePosService
	kalenderHijriyahSvc *service.KalenderHijriyahService
	hargaBBMSvc         *service.HargaBBMService
	ihsgSvc             *service.IHSGService
	krlSvc              *service.KRLService
	bpjsSvc             *service.BPJSService
	rekeningSvc         *service.RekeningService
	inflasiSvc          *service.InflasiService
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
	gempaSvc *service.GempaService,
	kodePosSvc *service.KodePosService,
	kalenderHijriyahSvc *service.KalenderHijriyahService,
	hargaBBMSvc *service.HargaBBMService,
	ihsgSvc *service.IHSGService,
	krlSvc *service.KRLService,
	bpjsSvc *service.BPJSService,
	rekeningSvc *service.RekeningService,
	inflasiSvc *service.InflasiService,
) *Resolver {
	return &Resolver{
		wilayahSvc:          wilayahSvc,
		hariLiburSvc:        hariLiburSvc,
		cuacaSvc:            cuacaSvc,
		kursSvc:             kursSvc,
		nikSvc:              nikSvc,
		kodeBankSvc:         kodeBankSvc,
		platNomorSvc:        platNomorSvc,
		waktuSholatSvc:      waktuSholatSvc,
		gempaSvc:            gempaSvc,
		kodePosSvc:          kodePosSvc,
		kalenderHijriyahSvc: kalenderHijriyahSvc,
		hargaBBMSvc:         hargaBBMSvc,
		ihsgSvc:             ihsgSvc,
		krlSvc:              krlSvc,
		bpjsSvc:             bpjsSvc,
		rekeningSvc:         rekeningSvc,
		inflasiSvc:          inflasiSvc,
	}
}
