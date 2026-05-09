package service

import (
	"fmt"
	"math"
	"time"
)

type KalenderHijriyahResult struct {
	TanggalMasehi   string
	TanggalHijriyah string
	Hari            string
	HariArab        string
	Bulan           string
	BulanArab       string
	Tahun           int
}

var bulanHijriyah = []string{
	"", "Muharram", "Safar", "Rabi'ul Awal", "Rabi'ul Akhir",
	"Jumadil Awal", "Jumadil Akhir", "Rajab", "Sya'ban",
	"Ramadan", "Syawal", "Dzulqa'dah", "Dzulhijjah",
}

var bulanHijriyahArab = []string{
	"", "مُحَرَّم", "صَفَر", "رَبِيعٌ الأوَّل", "رَبِيعٌ الثَّانِي",
	"جُمَادَى الأُولَى", "جُمَادَى الثَّانِيَة", "رَجَب", "شَعْبَان",
	"رَمَضَان", "شَوَّال", "ذُو الْقَعْدَة", "ذُو الْحِجَّة",
}

var hariIndonesia = []string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}
var hariArab = []string{"الأحد", "الإثنين", "الثلاثاء", "الأربعاء", "الخميس", "الجمعة", "السبت"}

func gregorianToHijri(year, month, day int) (int, int, int) {
	// Julian Day Number untuk tanggal Gregorian
	a := (14 - month) / 12
	y := year + 4800 - a
	m := month + 12*a - 3
	jdn := day + (153*m+2)/5 + 365*y + y/4 - y/100 + y/400 - 32045

	// Konversi JDN ke Hijriyah
	l := jdn - 1948440 + 10632
	n := (l-1)/10631
	l = l - 10631*n + 354
	j := (10985-l)/5316
	k := (50*l - 119) / 1768
	l = l - (1317*j)/900 - (j * 10) + (44*k)/1800 - k
	j = (11*n + 3) / 30
	i := 30 * (n + 1)
	i = i - j
	i = i + int(math.Floor(float64(l-1)/29.5001)) + 1
	l = l - int(math.Floor(float64(i-1)*29.5001))

	hYear := 30*(n-1) + int(math.Ceil(float64(i)/1.05263)) + j
	hMonth := i
	hDay := l

	return hYear, hMonth, hDay
}

type KalenderHijriyahService struct{}

func NewKalenderHijriyahService() *KalenderHijriyahService {
	return &KalenderHijriyahService{}
}

func (s *KalenderHijriyahService) Convert(tanggal string) (*KalenderHijriyahResult, error) {
	t, err := time.Parse("2006-01-02", tanggal)
	if err != nil {
		return nil, fmt.Errorf("format tanggal tidak valid, gunakan YYYY-MM-DD")
	}

	y, m, d := t.Date()
	hYear, hMonth, hDay := gregorianToHijri(y, int(m), d)

	weekday := int(t.Weekday())

	return &KalenderHijriyahResult{
		TanggalMasehi:   t.Format("02 January 2006"),
		TanggalHijriyah: fmt.Sprintf("%d %s %d H", hDay, bulanHijriyah[hMonth], hYear),
		Hari:            hariIndonesia[weekday],
		HariArab:        hariArab[weekday],
		Bulan:           bulanHijriyah[hMonth],
		BulanArab:       bulanHijriyahArab[hMonth],
		Tahun:           hYear,
	}, nil
}
