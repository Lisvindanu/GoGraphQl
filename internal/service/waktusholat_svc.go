package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/lisvindanuu/indonesiaql/internal/cache"
)

type WaktuSholatResult struct {
	Kota    string
	Tanggal string
	Subuh   string
	Terbit  string
	Dzuhur  string
	Ashar   string
	Maghrib string
	Isya    string
}

type WaktuSholatService struct {
	cache     *cache.Cache
	client    *http.Client
	cuacaSvc  *CuacaService
}

func NewWaktuSholatService(c *cache.Cache, cuacaSvc *CuacaService) *WaktuSholatService {
	return &WaktuSholatService{
		cache:    c,
		client:   &http.Client{Timeout: 10 * time.Second},
		cuacaSvc: cuacaSvc,
	}
}

type aladhanResponse struct {
	Code int `json:"code"`
	Data struct {
		Timings struct {
			Fajr    string `json:"Fajr"`
			Sunrise string `json:"Sunrise"`
			Dhuhr   string `json:"Dhuhr"`
			Asr     string `json:"Asr"`
			Maghrib string `json:"Maghrib"`
			Isha    string `json:"Isha"`
		} `json:"timings"`
		Date struct {
			Readable string `json:"readable"`
		} `json:"date"`
	} `json:"data"`
}

func (s *WaktuSholatService) Get(ctx context.Context, kota string) (*WaktuSholatResult, error) {
	cacheKey := "waktusholat:" + strings.ToLower(kota) + ":" + time.Now().Format("2006-01-02")
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached.(*WaktuSholatResult), nil
	}

	geo, err := s.cuacaSvc.geocode(ctx, kota)
	if err != nil {
		return nil, fmt.Errorf("lokasi '%s' tidak ditemukan", kota)
	}

	url := fmt.Sprintf(
		"https://api.aladhan.com/v1/timings?latitude=%.4f&longitude=%.4f&method=11",
		geo.Lat, geo.Lon,
	)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("buat request waktu sholat: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch waktu sholat: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 512*1024))
	var data aladhanResponse
	if err := json.Unmarshal(body, &data); err != nil || data.Code != 200 {
		return nil, fmt.Errorf("gagal parse data waktu sholat")
	}

	t := data.Data.Timings
	result := &WaktuSholatResult{
		Kota:    geo.Name,
		Tanggal: data.Data.Date.Readable,
		Subuh:   t.Fajr,
		Terbit:  t.Sunrise,
		Dzuhur:  t.Dhuhr,
		Ashar:   t.Asr,
		Maghrib: t.Maghrib,
		Isya:    t.Isha,
	}

	// Cache until end of day
	s.cache.Set(cacheKey, result, 6*time.Hour)
	return result, nil
}
