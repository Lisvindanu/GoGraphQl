package service

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type KodePosResult struct {
	KodePos   string
	Kelurahan string
	Kecamatan string
	Kota      string
	Provinsi  string
}

type KodePosService struct {
	db *pgxpool.Pool
}

func NewKodePosService(db *pgxpool.Pool) *KodePosService {
	return &KodePosService{db: db}
}

func (s *KodePosService) Search(ctx context.Context, kode string) ([]KodePosResult, error) {
	kode = strings.TrimSpace(kode)

	rows, err := s.db.Query(ctx, `
		SELECT
			k.kode_pos,
			k.nama       AS kelurahan,
			kec.nama     AS kecamatan,
			kot.nama     AS kota,
			p.nama       AS provinsi
		FROM kelurahan k
		JOIN kecamatan kec ON kec.kode = k.kecamatan_kode
		JOIN kota kot      ON kot.kode = kec.kota_kode
		JOIN provinsi p    ON p.kode  = kot.provinsi_kode
		WHERE k.kode_pos = $1
		ORDER BY p.nama, kot.nama, kec.nama, k.nama
		LIMIT 50
	`, kode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []KodePosResult
	for rows.Next() {
		var r KodePosResult
		if err := rows.Scan(&r.KodePos, &r.Kelurahan, &r.Kecamatan, &r.Kota, &r.Provinsi); err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return result, nil
}
