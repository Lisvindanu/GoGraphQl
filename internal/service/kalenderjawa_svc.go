package service

import (
	"fmt"
	"time"
)

const epochJD = 2427108 // JDN untuk 1 Januari 1933 (Selasa Wage, Wuku Sinta)
const epochJawaJD = 2317746 // JDN untuk 8 Juli 1633 (1 Suro 1555 AJ)

var namaHari = []string{"Ahad", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}
var namaPasaran = []string{"Legi", "Pahing", "Pon", "Wage", "Kliwon"}
var namaWuku = []string{
	"Sinta", "Landep", "Wukir", "Kurantil", "Tolu", "Gumbreg",
	"Warigalit", "Wariga", "Julungwangi", "Sungsang", "Galungan", "Kuningan",
	"Langkir", "Mandhasiya", "Julungpujut", "Pahang", "Kuruwelut", "Marakeh",
	"Tambir", "Medangkungan", "Matal", "Uye", "Menail", "Prangbakat",
	"Bala", "Wugu", "Wayang", "Kelawu", "Dukut", "Watugunung",
}
var namaWindu = []string{"Adi", "Kuntara", "Sengara", "Sancaya"}
var namaTahunJawa = []string{"Alip", "Ehe", "Jimawal", "Je", "Dal", "Be", "Wawu", "Jimakir"}

type KalenderJawaResult struct {
	TanggalMasehi   string
	Hari            string
	Pasaran         string
	Wuku            string
	TahunJawa       int
	NamaWindu       string
	TahunDalamWindu string
}

func dateToJD(t time.Time) int {
	y, m, d := t.Date()
	a := (14 - int(m)) / 12
	yr := y + 4800 - a
	mo := int(m) + 12*a - 3
	return d + (153*mo+2)/5 + 365*yr + yr/4 - yr/100 + yr/400 - 32045
}

func KalenderJawa(tanggal string) (*KalenderJawaResult, error) {
	t, err := time.Parse("2006-01-02", tanggal)
	if err != nil {
		return nil, fmt.Errorf("format tanggal tidak valid, gunakan YYYY-MM-DD: %w", err)
	}

	jd := dateToJD(t)

	hari := namaHari[t.Weekday()]

	// Pasaran: epoch 1 Jan 1933 = Wage (index 3)
	pasaranIdx := ((jd-epochJD)%5 + 3 + 5) % 5
	pasaran := namaPasaran[pasaranIdx]

	// Wuku: siklus 210 hari (30 wuku × 7 hari)
	daysSinceEpoch := jd - epochJD
	wukuDayOffset := ((daysSinceEpoch % 210) + 210) % 210
	wuku := namaWuku[wukuDayOffset/7]

	// Tahun Jawa
	daysSinceJawaEpoch := jd - epochJawaJD
	if daysSinceJawaEpoch < 0 {
		return nil, fmt.Errorf("tanggal terlalu lama, minimum 8 Juli 1633 M")
	}

	winduIdx := (daysSinceJawaEpoch / 2835) % 4
	tahunDalamWinduIdx := (daysSinceJawaEpoch / 354) % 8
	tahunJawa := 1555 + daysSinceJawaEpoch/354

	return &KalenderJawaResult{
		TanggalMasehi:   tanggal,
		Hari:            hari,
		Pasaran:         pasaran,
		Wuku:            wuku,
		TahunJawa:       tahunJawa,
		NamaWindu:       namaWindu[winduIdx],
		TahunDalamWindu: namaTahunJawa[tahunDalamWinduIdx],
	}, nil
}
