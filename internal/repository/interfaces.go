package repository

import (
	"context"

	"github.com/lisvindanuu/indonesiaql/internal/domain"
)

type WilayahRepository interface {
	ListProvinsi(ctx context.Context) ([]domain.Provinsi, error)
	GetProvinsiByKode(ctx context.Context, kode string) (*domain.Provinsi, error)
	ListKotaByProvinsi(ctx context.Context, provinsiKode string) ([]domain.Kota, error)
	GetKotaByKode(ctx context.Context, kode string) (*domain.Kota, error)
	ListKecamatanByKota(ctx context.Context, kotaKode string) ([]domain.Kecamatan, error)
	ListKelurahanByKecamatan(ctx context.Context, kecamatanKode string) ([]domain.Kelurahan, error)
	SearchWilayah(ctx context.Context, query string, limit int) ([]domain.WilayahSearchResult, error)
}

type HariLiburRepository interface {
	ListByTahun(ctx context.Context, tahun int) ([]domain.HariLibur, error)
	ListByBulan(ctx context.Context, tahun, bulan int) ([]domain.HariLibur, error)
	GetByTanggal(ctx context.Context, tanggal string) (*domain.HariLibur, error)
}
