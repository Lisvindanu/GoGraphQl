package service

import (
	"context"
	"time"

	"github.com/lisvindanuu/indonesiaql/internal/domain"
	"github.com/lisvindanuu/indonesiaql/internal/repository"
)

type HariLiburService struct {
	repo repository.HariLiburRepository
}

func NewHariLiburService(repo repository.HariLiburRepository) *HariLiburService {
	return &HariLiburService{repo: repo}
}

func (s *HariLiburService) ListByTahun(ctx context.Context, tahun int) ([]domain.HariLibur, error) {
	return s.repo.ListByTahun(ctx, tahun)
}

func (s *HariLiburService) ListByBulan(ctx context.Context, tahun, bulan int) ([]domain.HariLibur, error) {
	return s.repo.ListByBulan(ctx, tahun, bulan)
}

func (s *HariLiburService) GetHariIni(ctx context.Context) (*domain.HariLibur, error) {
	today := time.Now().Format("2006-01-02")
	h, err := s.repo.GetByTanggal(ctx, today)
	if err != nil {
		return nil, nil // bukan hari libur
	}
	return h, nil
}
