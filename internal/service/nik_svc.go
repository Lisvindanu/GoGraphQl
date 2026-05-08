package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	nikdata "github.com/lisvindanuu/indonesiaql/internal/nik"
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
	kodeKab := nik[0:4]
	kodeKec := nik[0:6]

	// Provinsi from DB (BPS and Kemendagri use same 2-digit codes for provinces)
	provinsi, err := s.wilayahRepo.GetProvinsiByKode(ctx, kodeProvinsi)
	if err == nil && provinsi != nil {
		result.Provinsi = &provinsi.Nama
	}

	// Kabupaten from static Kemendagri map (matches NIK standard)
	if nama, ok := nikdata.NIKKabupaten[kodeKab]; ok {
		result.Kota = &nama
	}

	// Kecamatan from static Kemendagri map
	if nama, ok := nikdata.NIKKecamatan[kodeKec]; ok {
		result.Kecamatan = &nama
	}

	ddStr := nik[6:8]
	mmStr := nik[8:10]
	yyStr := nik[10:12]

	dd, _ := strconv.Atoi(ddStr)
	mm, _ := strconv.Atoi(mmStr)
	yy, _ := strconv.Atoi(yyStr)

	jenisKelamin := "LAKI-LAKI"
	if dd > 40 {
		dd -= 40
		jenisKelamin = "PEREMPUAN"
	}
	result.JenisKelamin = &jenisKelamin

	tahun := 1900 + yy
	if yy < 22 {
		tahun = 2000 + yy
	}

	if dd >= 1 && dd <= 31 && mm >= 1 && mm <= 12 {
		tgl := fmt.Sprintf("%02d/%02d/%04d", dd, mm, tahun)
		result.TanggalLahir = &tgl
	} else {
		result.Errors = append(result.Errors, "Format tanggal lahir tidak valid")
	}

	result.Valid = len(result.Errors) == 0
	return result, nil
}
