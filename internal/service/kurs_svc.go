package service

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
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

type biKursResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    biBody   `xml:"Body"`
}

type biBody struct {
	Response biKursListResponse `xml:"getKursResponse"`
}

type biKursListResponse struct {
	Return biKursReturn `xml:"return"`
}

type biKursReturn struct {
	Tanggal string   `xml:"Tanggal"`
	Kurs    []biKurs `xml:"kurs"`
}

type biKurs struct {
	MataUang string `xml:"mata_uang"`
	Beli     string `xml:"beli"`
	Jual     string `xml:"jual"`
	Tengah   string `xml:"tengah"`
}

const biKursURL = "https://www.bi.go.id/biwebservice/wskursbi.asmx/getKursLokal"

func (s *KursService) GetKurs(ctx context.Context, mataUang *string) ([]KursResult, error) {
	cacheKey := "kurs:all"
	if mataUang != nil && *mataUang != "" {
		cacheKey = "kurs:" + strings.ToUpper(*mataUang)
	}
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached.([]KursResult), nil
	}

	req, err := http.NewRequestWithContext(ctx, "GET", biKursURL, nil)
	if err != nil {
		return nil, fmt.Errorf("buat request BI: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch kurs BI: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 2*1024*1024))
	if err != nil {
		return nil, fmt.Errorf("baca response BI: %w", err)
	}

	var biResp biKursResponse
	if err := xml.Unmarshal(body, &biResp); err != nil {
		return nil, fmt.Errorf("parse XML BI: %w", err)
	}

	tanggal := biResp.Body.Response.Return.Tanggal
	var results []KursResult
	for _, k := range biResp.Body.Response.Return.Kurs {
		if mataUang != nil && *mataUang != "" && !strings.EqualFold(k.MataUang, *mataUang) {
			continue
		}
		results = append(results, KursResult{
			MataUang:   k.MataUang,
			KursBeli:   parseKurs(k.Beli),
			KursJual:   parseKurs(k.Jual),
			KursTengah: parseKurs(k.Tengah),
			Tanggal:    tanggal,
		})
	}

	s.cache.Set(cacheKey, results, 1*time.Hour)
	return results, nil
}

func parseKurs(s string) float64 {
	s = strings.ReplaceAll(s, ",", "")
	v, _ := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return v
}
