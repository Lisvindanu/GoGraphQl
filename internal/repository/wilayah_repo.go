package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lisvindanuu/indonesiaql/internal/domain"
)

type wilayahRepo struct {
	db *pgxpool.Pool
}

func NewWilayahRepository(db *pgxpool.Pool) WilayahRepository {
	return &wilayahRepo{db: db}
}

func (r *wilayahRepo) ListProvinsi(ctx context.Context) ([]domain.Provinsi, error) {
	rows, err := r.db.Query(ctx, `SELECT kode, nama FROM provinsi ORDER BY nama`)
	if err != nil {
		return nil, fmt.Errorf("list provinsi: %w", err)
	}
	defer rows.Close()

	var result []domain.Provinsi
	for rows.Next() {
		var p domain.Provinsi
		if err := rows.Scan(&p.Kode, &p.Nama); err != nil {
			return nil, fmt.Errorf("scan provinsi: %w", err)
		}
		result = append(result, p)
	}
	return result, rows.Err()
}

func (r *wilayahRepo) GetProvinsiByKode(ctx context.Context, kode string) (*domain.Provinsi, error) {
	var p domain.Provinsi
	err := r.db.QueryRow(ctx, `SELECT kode, nama FROM provinsi WHERE kode = $1`, kode).
		Scan(&p.Kode, &p.Nama)
	if err != nil {
		return nil, fmt.Errorf("get provinsi %s: %w", kode, err)
	}
	return &p, nil
}

func (r *wilayahRepo) ListKotaByProvinsi(ctx context.Context, provinsiKode string) ([]domain.Kota, error) {
	rows, err := r.db.Query(ctx,
		`SELECT kode, provinsi_kode, nama FROM kota WHERE provinsi_kode = $1 ORDER BY nama`,
		provinsiKode)
	if err != nil {
		return nil, fmt.Errorf("list kota: %w", err)
	}
	defer rows.Close()

	var result []domain.Kota
	for rows.Next() {
		var k domain.Kota
		if err := rows.Scan(&k.Kode, &k.ProvinsiKode, &k.Nama); err != nil {
			return nil, fmt.Errorf("scan kota: %w", err)
		}
		result = append(result, k)
	}
	return result, rows.Err()
}

func (r *wilayahRepo) GetKotaByKode(ctx context.Context, kode string) (*domain.Kota, error) {
	var k domain.Kota
	err := r.db.QueryRow(ctx,
		`SELECT kode, provinsi_kode, nama FROM kota WHERE kode = $1`, kode).
		Scan(&k.Kode, &k.ProvinsiKode, &k.Nama)
	if err != nil {
		return nil, fmt.Errorf("get kota %s: %w", kode, err)
	}
	return &k, nil
}

func (r *wilayahRepo) ListKecamatanByKota(ctx context.Context, kotaKode string) ([]domain.Kecamatan, error) {
	rows, err := r.db.Query(ctx,
		`SELECT kode, kota_kode, nama FROM kecamatan WHERE kota_kode = $1 ORDER BY nama`,
		kotaKode)
	if err != nil {
		return nil, fmt.Errorf("list kecamatan: %w", err)
	}
	defer rows.Close()

	var result []domain.Kecamatan
	for rows.Next() {
		var k domain.Kecamatan
		if err := rows.Scan(&k.Kode, &k.KotaKode, &k.Nama); err != nil {
			return nil, fmt.Errorf("scan kecamatan: %w", err)
		}
		result = append(result, k)
	}
	return result, rows.Err()
}

func (r *wilayahRepo) ListKelurahanByKecamatan(ctx context.Context, kecamatanKode string) ([]domain.Kelurahan, error) {
	rows, err := r.db.Query(ctx,
		`SELECT kode, kecamatan_kode, nama, kode_pos FROM kelurahan WHERE kecamatan_kode = $1 ORDER BY nama`,
		kecamatanKode)
	if err != nil {
		return nil, fmt.Errorf("list kelurahan: %w", err)
	}
	defer rows.Close()

	var result []domain.Kelurahan
	for rows.Next() {
		var k domain.Kelurahan
		if err := rows.Scan(&k.Kode, &k.KecamatanKode, &k.Nama, &k.KodePos); err != nil {
			return nil, fmt.Errorf("scan kelurahan: %w", err)
		}
		result = append(result, k)
	}
	return result, rows.Err()
}

func (r *wilayahRepo) SearchWilayah(ctx context.Context, query string, limit int) ([]domain.WilayahSearchResult, error) {
	q := "%" + query + "%"
	rows, err := r.db.Query(ctx, `
		SELECT kode, nama, 'provinsi' as tipe, NULL::text as kota, NULL::text as provinsi
		FROM provinsi WHERE nama ILIKE $1
		UNION ALL
		SELECT k.kode, k.nama, 'kota', NULL, p.nama
		FROM kota k JOIN provinsi p ON p.kode = k.provinsi_kode WHERE k.nama ILIKE $1
		UNION ALL
		SELECT kc.kode, kc.nama, 'kecamatan', k.nama, p.nama
		FROM kecamatan kc
		JOIN kota k ON k.kode = kc.kota_kode
		JOIN provinsi p ON p.kode = k.provinsi_kode
		WHERE kc.nama ILIKE $1
		LIMIT $2
	`, q, limit)
	if err != nil {
		return nil, fmt.Errorf("search wilayah: %w", err)
	}
	defer rows.Close()

	var result []domain.WilayahSearchResult
	for rows.Next() {
		var s domain.WilayahSearchResult
		if err := rows.Scan(&s.Kode, &s.Nama, &s.Tipe, &s.Kota, &s.Provinsi); err != nil {
			return nil, fmt.Errorf("scan search result: %w", err)
		}
		result = append(result, s)
	}
	return result, rows.Err()
}
