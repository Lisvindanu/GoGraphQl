package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type IHSGResult struct {
	Symbol              string
	Nama                string
	Harga               float64
	Perubahan           float64
	PersentasePerubahan float64
	Open                float64
	High                float64
	Low                 float64
	Volume              int
	Waktu               string
}

type IHSGService struct {
	client *http.Client
}

func NewIHSGService() *IHSGService {
	return &IHSGService{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *IHSGService) fetchChart(symbol string) (*IHSGResult, error) {
	encoded := url.PathEscape(symbol)
	apiURL := "https://query1.finance.yahoo.com/v8/finance/chart/" + encoded + "?range=1d&interval=1d"

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch Yahoo Finance: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var root struct {
		Chart struct {
			Result []struct {
				Meta struct {
					Symbol               string  `json:"symbol"`
					ShortName            string  `json:"shortName"`
					LongName             string  `json:"longName"`
					RegularMarketPrice   float64 `json:"regularMarketPrice"`
					ChartPreviousClose   float64 `json:"chartPreviousClose"`
					RegularMarketDayHigh float64 `json:"regularMarketDayHigh"`
					RegularMarketDayLow  float64 `json:"regularMarketDayLow"`
					RegularMarketVolume  int     `json:"regularMarketVolume"`
					RegularMarketTime    int64   `json:"regularMarketTime"`
					RegularMarketOpen    float64 `json:"regularMarketOpen"`
				} `json:"meta"`
			} `json:"result"`
			Error *struct {
				Code        string `json:"code"`
				Description string `json:"description"`
			} `json:"error"`
		} `json:"chart"`
	}

	if err := json.Unmarshal(body, &root); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	if root.Chart.Error != nil {
		return nil, fmt.Errorf("Yahoo Finance error: %s - %s", root.Chart.Error.Code, root.Chart.Error.Description)
	}
	if len(root.Chart.Result) == 0 {
		return nil, fmt.Errorf("data tidak ditemukan untuk simbol '%s'", symbol)
	}

	meta := root.Chart.Result[0].Meta
	perubahan := meta.RegularMarketPrice - meta.ChartPreviousClose
	persen := 0.0
	if meta.ChartPreviousClose != 0 {
		persen = (perubahan / meta.ChartPreviousClose) * 100
	}

	nama := meta.ShortName
	if nama == "" {
		nama = meta.LongName
	}
	if nama == "" {
		nama = strings.TrimSuffix(meta.Symbol, "=X")
	}

	waktu := ""
	if meta.RegularMarketTime > 0 {
		t := time.Unix(meta.RegularMarketTime, 0).In(time.FixedZone("WIB", 7*3600))
		waktu = t.Format("2006-01-02 15:04:05 WIB")
	}

	return &IHSGResult{
		Symbol:              meta.Symbol,
		Nama:                nama,
		Harga:               meta.RegularMarketPrice,
		Perubahan:           perubahan,
		PersentasePerubahan: persen,
		Open:                meta.RegularMarketOpen,
		High:                meta.RegularMarketDayHigh,
		Low:                 meta.RegularMarketDayLow,
		Volume:              meta.RegularMarketVolume,
		Waktu:               waktu,
	}, nil
}

func (s *IHSGService) GetIHSG() (*IHSGResult, error) {
	return s.fetchChart("^JKSE")
}

func (s *IHSGService) GetSaham(kode string) (*IHSGResult, error) {
	kode = strings.ToUpper(strings.TrimSpace(kode))
	return s.fetchChart(kode)
}
