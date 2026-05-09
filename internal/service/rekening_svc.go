package service

import (
	"fmt"
	"strings"
	"unicode"
)

type ValidasiRekeningResult struct {
	Valid      bool
	Bank       string
	NoRekening string
	Panjang    int
	Keterangan string
}

type bankInfo struct {
	Nama    string
	Panjang []int
}

var bankDatabase = map[string]bankInfo{
	"BCA":      {Nama: "Bank Central Asia (BCA)", Panjang: []int{10}},
	"MANDIRI":  {Nama: "Bank Mandiri", Panjang: []int{13}},
	"BRI":      {Nama: "Bank Rakyat Indonesia (BRI)", Panjang: []int{15}},
	"BNI":      {Nama: "Bank Negara Indonesia (BNI)", Panjang: []int{10}},
	"BTN":      {Nama: "Bank Tabungan Negara (BTN)", Panjang: []int{15}},
	"CIMB":     {Nama: "CIMB Niaga", Panjang: []int{13}},
	"DANAMON":  {Nama: "Bank Danamon", Panjang: []int{10}},
	"PERMATA":  {Nama: "Bank Permata", Panjang: []int{10}},
	"MAYBANK":  {Nama: "Maybank Indonesia", Panjang: []int{12}},
	"BSI":      {Nama: "Bank Syariah Indonesia (BSI)", Panjang: []int{10}},
	"MUAMALAT": {Nama: "Bank Muamalat", Panjang: []int{16}},
	"OCBC":     {Nama: "OCBC NISP", Panjang: []int{10, 16}},
	"PANIN":    {Nama: "Bank Panin", Panjang: []int{13}},
	"MEGA":     {Nama: "Bank Mega", Panjang: []int{15}},
	"JAGO":     {Nama: "Bank Jago", Panjang: []int{10}},
}

type RekeningService struct{}

func NewRekeningService() *RekeningService { return &RekeningService{} }

func (s *RekeningService) Validasi(bank, noRekening string) (*ValidasiRekeningResult, error) {
	bankKode := strings.ToUpper(strings.TrimSpace(bank))
	noRek := strings.TrimSpace(noRekening)

	for _, c := range noRek {
		if !unicode.IsDigit(c) {
			return &ValidasiRekeningResult{
				Valid:      false,
				Bank:       bankKode,
				NoRekening: noRek,
				Panjang:    len(noRek),
				Keterangan: "Nomor rekening hanya boleh berisi angka",
			}, nil
		}
	}

	info, found := bankDatabase[bankKode]
	if !found {
		return nil, fmt.Errorf("kode bank '%s' tidak dikenali. Gunakan: BCA, MANDIRI, BRI, BNI, BTN, CIMB, DANAMON, PERMATA, MAYBANK, BSI, MUAMALAT, OCBC, PANIN, MEGA, JAGO", bankKode)
	}

	panjang := len(noRek)
	valid := false
	for _, p := range info.Panjang {
		if panjang == p {
			valid = true
			break
		}
	}

	keterangan := ""
	if valid {
		keterangan = fmt.Sprintf("Format nomor rekening %s valid (%d digit)", info.Nama, panjang)
	} else {
		expected := make([]string, len(info.Panjang))
		for i, p := range info.Panjang {
			expected[i] = fmt.Sprintf("%d", p)
		}
		keterangan = fmt.Sprintf("Format tidak valid: %s membutuhkan %s digit, diterima %d digit",
			info.Nama, strings.Join(expected, " atau "), panjang)
	}

	return &ValidasiRekeningResult{
		Valid:      valid,
		Bank:       info.Nama,
		NoRekening: noRek,
		Panjang:    panjang,
		Keterangan: keterangan,
	}, nil
}
