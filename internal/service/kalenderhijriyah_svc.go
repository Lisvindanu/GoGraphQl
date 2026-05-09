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

func islamicToJD(year, month, day int) float64 {
	return float64(day) +
		math.Ceil(29.5*float64(month-1)) +
		float64(year-1)*354.0 +
		math.Floor((11.0*float64(year)+3.0)/30.0) +
		1948439.5
}

func gregorianToHijri(year, month, day int) (int, int, int) {
	// Gregorian ke Julian Day Number
	a := (14 - month) / 12
	y := year + 4800 - a
	m := month + 12*a - 3
	jdn := day + (153*m+2)/5 + 365*y + y/4 - y/100 + y/400 - 32045

	// JDN ke Hijriyah (Fourmilab algorithm, epoch = 1948439.5)
	jd := float64(jdn) + 0.5

	hYear := int(math.Floor(((30.0*(jd-1948439.5)) + 10646.0) / 10631.0))

	hMonth := int(math.Ceil((jd-(29.0+islamicToJD(hYear, 1, 1)))/29.5)) + 1
	if hMonth < 1 {
		hMonth = 1
	}
	if hMonth > 12 {
		hMonth = 12
	}

	hDay := int(jd-islamicToJD(hYear, hMonth, 1)) + 1

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
