package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lisvindanuu/indonesiaql/internal/domain"
)

type hariLiburRepo struct {
	db *pgxpool.Pool
}

func NewHariLiburRepository(db *pgxpool.Pool) HariLiburRepository {
	return &hariLiburRepo{db: db}
}

func (r *hariLiburRepo) ListByTahun(ctx context.Context, tahun int) ([]domain.HariLibur, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, tanggal, nama, tahun, jenis FROM hari_libur WHERE tahun = $1 ORDER BY tanggal`,
		tahun)
	if err != nil {
		return nil, fmt.Errorf("list hari libur tahun %d: %w", tahun, err)
	}
	defer rows.Close()
	return scanHariLibur(rows)
}

func (r *hariLiburRepo) ListByBulan(ctx context.Context, tahun, bulan int) ([]domain.HariLibur, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, tanggal, nama, tahun, jenis FROM hari_libur
		 WHERE tahun = $1 AND EXTRACT(MONTH FROM tanggal) = $2 ORDER BY tanggal`,
		tahun, bulan)
	if err != nil {
		return nil, fmt.Errorf("list hari libur %d/%d: %w", bulan, tahun, err)
	}
	defer rows.Close()
	return scanHariLibur(rows)
}

func (r *hariLiburRepo) GetByTanggal(ctx context.Context, tanggal string) (*domain.HariLibur, error) {
	var h domain.HariLibur
	err := r.db.QueryRow(ctx,
		`SELECT id, tanggal, nama, tahun, jenis FROM hari_libur WHERE tanggal = $1 LIMIT 1`,
		tanggal).Scan(&h.ID, &h.Tanggal, &h.Nama, &h.Tahun, &h.Jenis)
	if err != nil {
		return nil, fmt.Errorf("get hari libur %s: %w", tanggal, err)
	}
	return &h, nil
}

type scannable interface {
	Next() bool
	Scan(...any) error
	Err() error
}

func scanHariLibur(rows scannable) ([]domain.HariLibur, error) {
	var result []domain.HariLibur
	for rows.Next() {
		var h domain.HariLibur
		if err := rows.Scan(&h.ID, &h.Tanggal, &h.Nama, &h.Tahun, &h.Jenis); err != nil {
			return nil, fmt.Errorf("scan hari libur: %w", err)
		}
		result = append(result, h)
	}
	return result, rows.Err()
}
