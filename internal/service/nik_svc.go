package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/lisvindanuu/indonesiaql/internal/repository"
)

type NIKService struct {
	wilayahRepo repository.WilayahRepository
}

func NewNIKService(wilayahRepo repository.WilayahRepository) *NIKService {
	return &NIKService{wilayahRepo: wilayahRepo}
}

type NIKResult struct {
	Valid        bool
	NIK          string
	Provinsi     *string
	Kota         *string
	Kecamatan    *string
	TanggalLahir *string
	JenisKelamin *string
	Errors       []string
}

func (s *NIKService) Validasi(ctx context.Context, nik string) (*NIKResult, error) {
	result := &NIKResult{NIK: nik}

	nik = strings.TrimSpace(nik)
	if len(nik) != 16 {
		result.Errors = append(result.Errors, fmt.Sprintf("NIK harus 16 digit, ditemukan %d karakter", len(nik)))
		return result, nil
	}
	for _, c := range nik {
		if !unicode.IsDigit(c) {
			result.Errors = append(result.Errors, "NIK hanya boleh berisi angka")
			return result, nil
		}
	}

	kodeProvinsi := nik[0:2]
	kodeKota := nik[0:4]

	provinsi, err := s.wilayahRepo.GetProvinsiByKode(ctx, kodeProvinsi)
	if err == nil && provinsi != nil {
		result.Provinsi = &provinsi.Nama
	}

	kota, err := s.wilayahRepo.GetKotaByKode(ctx, kodeKota)
	if err == nil && kota != nil {
		result.Kota = &kota.Nama
	}

	ddStr := nik[6:8]
	mmStr := nik[8:10]
	yyStr := nik[10:12]

	dd, _ := strconv.Atoi(ddStr)
	mm, _ := strconv.Atoi(mmStr)
	yy, _ := strconv.Atoi(yyStr)

	jenisKelamin := "Laki-laki"
	if dd > 40 {
		dd -= 40
		jenisKelamin = "Perempuan"
	}
	result.JenisKelamin = &jenisKelamin

	tahun := 1900 + yy
	if yy < 30 {
		tahun = 2000 + yy
	}

	if dd >= 1 && dd <= 31 && mm >= 1 && mm <= 12 {
		tgl := fmt.Sprintf("%04d-%02d-%02d", tahun, mm, dd)
		result.TanggalLahir = &tgl
	} else {
		result.Errors = append(result.Errors, "Format tanggal lahir tidak valid")
	}

	result.Valid = len(result.Errors) == 0
	return result, nil
}
