package service

import (
	"context"

	"github.com/lisvindanuu/indonesiaql/internal/domain"
	"github.com/lisvindanuu/indonesiaql/internal/repository"
)

type WilayahService struct {
	repo repository.WilayahRepository
}

func NewWilayahService(repo repository.WilayahRepository) *WilayahService {
	return &WilayahService{repo: repo}
}

func (s *WilayahService) ListProvinsi(ctx context.Context) ([]domain.Provinsi, error) {
	return s.repo.ListProvinsi(ctx)
}

func (s *WilayahService) GetProvinsi(ctx context.Context, kode string) (*domain.Provinsi, error) {
	return s.repo.GetProvinsiByKode(ctx, kode)
}

func (s *WilayahService) ListKota(ctx context.Context, provinsiKode string) ([]domain.Kota, error) {
	return s.repo.ListKotaByProvinsi(ctx, provinsiKode)
}

func (s *WilayahService) GetKota(ctx context.Context, kode string) (*domain.Kota, error) {
	return s.repo.GetKotaByKode(ctx, kode)
}

func (s *WilayahService) ListKecamatan(ctx context.Context, kotaKode string) ([]domain.Kecamatan, error) {
	return s.repo.ListKecamatanByKota(ctx, kotaKode)
}

func (s *WilayahService) ListKelurahan(ctx context.Context, kecamatanKode string) ([]domain.Kelurahan, error) {
	return s.repo.ListKelurahanByKecamatan(ctx, kecamatanKode)
}

func (s *WilayahService) Search(ctx context.Context, query string, limit int) ([]domain.WilayahSearchResult, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}
	return s.repo.SearchWilayah(ctx, query, limit)
}
