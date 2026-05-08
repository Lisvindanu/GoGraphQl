package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/lisvindanuu/indonesiaql/internal/cache"
)

type KursService struct {
	cache  *cache.Cache
	client *http.Client
}

func NewKursService(c *cache.Cache) *KursService {
	return &KursService{
		cache:  c,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

type KursResult struct {
	MataUang   string
	KursBeli   float64
	KursJual   float64
	KursTengah float64
	Tanggal    string
}

type erAPIResponse struct {
	Result         string             `json:"result"`
	TimeLastUpdate string             `json:"time_last_update_utc"`
	BaseCode       string             `json:"base_code"`
	Rates          map[string]float64 `json:"rates"`
}

// erAPIURL fetches all rates with USD as base (free tier, no API key required).
const erAPIURL = "https://open.er-api.com/v6/latest/USD"

// defaultCurrencies is the list of currencies shown when no filter is applied.
var defaultCurrencies = []string{
	"USD", "EUR", "SGD", "MYR", "JPY", "GBP", "AUD",
	"SAR", "CNY", "KRW", "HKD", "THB", "CHF", "CAD",
}

func (s *KursService) GetKurs(ctx context.Context, mataUang *string) ([]KursResult, error) {
	cacheKey := "kurs:all"
	filterCur := ""
	if mataUang != nil && *mataUang != "" {
		filterCur = strings.ToUpper(*mataUang)
		cacheKey = "kurs:" + filterCur
	}
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached.([]KursResult), nil
	}

	req, err := http.NewRequestWithContext(ctx, "GET", erAPIURL, nil)
	if err != nil {
		return nil, fmt.Errorf("buat request kurs: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch kurs: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1*1024*1024))
	if err != nil {
		return nil, fmt.Errorf("baca response kurs: %w", err)
	}

	var erResp erAPIResponse
	if err := json.Unmarshal(body, &erResp); err != nil {
		return nil, fmt.Errorf("parse response kurs: %w", err)
	}

	if erResp.Result != "success" {
		return nil, fmt.Errorf("kurs API tidak tersedia")
	}

	idrRate, ok := erResp.Rates["IDR"]
	if !ok || idrRate == 0 {
		return nil, fmt.Errorf("rate IDR tidak tersedia")
	}

	tanggal := parseTanggalKurs(erResp.TimeLastUpdate)

	// Determine which currencies to return
	currencies := defaultCurrencies
	if filterCur != "" {
		currencies = []string{filterCur}
	}

	var results []KursResult
	for _, cur := range currencies {
		curRate, ok := erResp.Rates[cur]
		if !ok || curRate == 0 {
			continue
		}
		// 1 CUR = (idrRate / curRate) IDR
		mid := idrRate / curRate
		results = append(results, KursResult{
			MataUang:   cur,
			KursBeli:   math.Round(mid * 0.995),
			KursJual:   math.Round(mid * 1.005),
			KursTengah: math.Round(mid),
			Tanggal:    tanggal,
		})
	}

	s.cache.Set(cacheKey, results, 1*time.Hour)
	return results, nil
}

// parseTanggalKurs parses the RFC1123Z date string from ExchangeRate-API
// into YYYY-MM-DD format. Falls back to today's date on error.
func parseTanggalKurs(s string) string {
	t, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", s)
	if err != nil {
		return time.Now().Format("2006-01-02")
	}
	return t.Format("2006-01-02")
}
