package service

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/lisvindanuu/indonesiaql/internal/cache"
)

type CuacaService struct {
	cache  *cache.Cache
	client *http.Client
}

func NewCuacaService(c *cache.Cache) *CuacaService {
	return &CuacaService{
		cache:  c,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

type Prakiraan struct {
	Waktu          string
	Suhu           string
	Kelembapan     string
	Cuaca          string
	KecepatanAngin string
	ArahAngin      string
}

type CuacaResult struct {
	Provinsi  string
	Kota      string
	Prakiraan []Prakiraan
}

type bmkgData struct {
	XMLName  xml.Name     `xml:"data"`
	Forecast bmkgForecast `xml:"forecast"`
}

type bmkgForecast struct {
	Area []bmkgArea `xml:"area"`
}

type bmkgArea struct {
	Description string      `xml:"description,attr"`
	Parameters  []bmkgParam `xml:"parameter"`
}

type bmkgParam struct {
	ID        string          `xml:"id,attr"`
	TimeRange []bmkgTimeRange `xml:"timerange"`
}

type bmkgTimeRange struct {
	Datetime string      `xml:"datetime,attr"`
	Value    []bmkgValue `xml:"value"`
}

type bmkgValue struct {
	Text string `xml:",chardata"`
}

var namaProvinsiURL = map[string]string{
	"11": "Aceh", "12": "SumateraUtara", "13": "SumateraBarat",
	"14": "Riau", "15": "Jambi", "16": "SumateraSelatan",
	"17": "Bengkulu", "18": "Lampung", "19": "KepulauanBangkaBelitung",
	"21": "KepulauanRiau", "31": "DKIJakarta", "32": "JawaBarat",
	"33": "JawaTengah", "34": "DIYogyakarta", "35": "JawaTimur",
	"36": "Banten", "51": "Bali", "52": "NusaTenggaraBarat",
	"53": "NusaTenggaraTimur", "61": "KalimantanBarat", "62": "KalimantanTengah",
	"63": "KalimantanSelatan", "64": "KalimantanTimur", "65": "KalimantanUtara",
	"71": "SulawesiUtara", "72": "SulawesiTengah", "73": "SulawesiSelatan",
	"74": "SulawesiTenggara", "75": "Gorontalo", "76": "SulawesiBarat",
	"81": "Maluku", "82": "MalukuUtara", "91": "PapuaBarat", "94": "Papua",
}

func (s *CuacaService) GetCuaca(ctx context.Context, provinsiKode, kotaQuery string) (*CuacaResult, error) {
	cacheKey := "cuaca:" + provinsiKode + ":" + strings.ToLower(kotaQuery)
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached.(*CuacaResult), nil
	}

	namaProvinsi, ok := namaProvinsiURL[provinsiKode]
	if !ok {
		return nil, fmt.Errorf("kode provinsi tidak dikenali: %s", provinsiKode)
	}

	url := fmt.Sprintf("https://data.bmkg.go.id/DataMKG/MEWS/DigitalForecast/DigitalForecast-%s.xml", namaProvinsi)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("buat request BMKG: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch BMKG: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 5*1024*1024))
	if err != nil {
		return nil, fmt.Errorf("baca response BMKG: %w", err)
	}

	var data bmkgData
	if err := xml.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("parse XML BMKG: %w", err)
	}

	result := parseBMKG(&data, namaProvinsi, kotaQuery)
	if result == nil {
		return nil, fmt.Errorf("kota '%s' tidak ditemukan dalam data BMKG", kotaQuery)
	}

	s.cache.Set(cacheKey, result, 30*time.Minute)
	return result, nil
}

func parseBMKG(data *bmkgData, provinsi, kotaQuery string) *CuacaResult {
	kotaLower := strings.ToLower(kotaQuery)
	for _, area := range data.Forecast.Area {
		if !strings.Contains(strings.ToLower(area.Description), kotaLower) {
			continue
		}
		paramMap := make(map[string][]bmkgTimeRange)
		for _, p := range area.Parameters {
			paramMap[p.ID] = p.TimeRange
		}
		var prakiraans []Prakiraan
		trs := paramMap["t"]
		for i := 0; i < len(trs) && i < 8; i++ {
			p := Prakiraan{Waktu: trs[i].Datetime}
			if len(trs[i].Value) > 0 {
				p.Suhu = trs[i].Value[0].Text + "°C"
			}
			if hTrs := paramMap["hu"]; i < len(hTrs) && len(hTrs[i].Value) > 0 {
				p.Kelembapan = hTrs[i].Value[0].Text + "%"
			}
			if wTrs := paramMap["weather"]; i < len(wTrs) && len(wTrs[i].Value) > 0 {
				p.Cuaca = wTrs[i].Value[0].Text
			}
			if wsTrs := paramMap["ws"]; i < len(wsTrs) && len(wsTrs[i].Value) > 0 {
				p.KecepatanAngin = wsTrs[i].Value[0].Text + " km/h"
			}
			if wdTrs := paramMap["wd"]; i < len(wdTrs) && len(wdTrs[i].Value) > 0 {
				p.ArahAngin = wdTrs[i].Value[0].Text
			}
			prakiraans = append(prakiraans, p)
		}
		return &CuacaResult{Provinsi: provinsi, Kota: area.Description, Prakiraan: prakiraans}
	}
	return nil
}
