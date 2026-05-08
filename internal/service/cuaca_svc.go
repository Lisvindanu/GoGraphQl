package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
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

type geoResult struct {
	Lat  float64
	Lon  float64
	Name string
}

type openMeteoGeoResponse struct {
	Results []struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Name      string  `json:"name"`
	} `json:"results"`
}

type openMeteoResponse struct {
	Hourly struct {
		Time             []string  `json:"time"`
		Temperature2m    []float64 `json:"temperature_2m"`
		RelativeHumidity []float64 `json:"relative_humidity_2m"`
		WeatherCode      []int     `json:"weather_code"`
		WindSpeed10m     []float64 `json:"wind_speed_10m"`
		WindDirection10m []float64 `json:"wind_direction_10m"`
	} `json:"hourly"`
}

func wmoToIndonesian(code int) string {
	switch {
	case code == 0:
		return "Cerah"
	case code <= 3:
		return "Sebagian Berawan"
	case code <= 48:
		return "Berkabut"
	case code <= 55:
		return "Gerimis"
	case code <= 65:
		return "Hujan"
	case code <= 67:
		return "Hujan Bercampur Es"
	case code <= 77:
		return "Salju"
	case code <= 82:
		return "Hujan Lebat"
	case code <= 86:
		return "Hujan Salju Lebat"
	case code <= 99:
		return "Badai Petir"
	default:
		return "Tidak Diketahui"
	}
}

func windDegToDir(deg float64) string {
	dirs := []string{"U", "TL", "T", "TG", "S", "BD", "B", "BL"}
	idx := int(math.Round(deg/45.0)) % 8
	return dirs[idx]
}

func (s *CuacaService) geocode(ctx context.Context, query string) (*geoResult, error) {
	cacheKey := "geocode:" + strings.ToLower(query)
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached.(*geoResult), nil
	}

	geoURL := "https://geocoding-api.open-meteo.com/v1/search?name=" +
		url.QueryEscape(query+" Indonesia") + "&count=1&language=id&format=json"

	req, err := http.NewRequestWithContext(ctx, "GET", geoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("buat request geocode: %w", err)
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch geocode: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 512*1024))
	var geo openMeteoGeoResponse
	if err := json.Unmarshal(body, &geo); err != nil || len(geo.Results) == 0 {
		return nil, fmt.Errorf("lokasi '%s' tidak ditemukan", query)
	}

	r := &geoResult{
		Lat:  geo.Results[0].Latitude,
		Lon:  geo.Results[0].Longitude,
		Name: geo.Results[0].Name,
	}
	s.cache.Set(cacheKey, r, 24*time.Hour)
	return r, nil
}

func (s *CuacaService) GetCuaca(ctx context.Context, provinsiKode, kotaQuery string) (*CuacaResult, error) {
	cacheKey := "cuaca:" + strings.ToLower(kotaQuery)
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached.(*CuacaResult), nil
	}

	info, err := s.geocode(ctx, kotaQuery)
	if err != nil {
		return nil, err
	}

	weatherURL := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f"+
			"&hourly=temperature_2m,relative_humidity_2m,weather_code,wind_speed_10m,wind_direction_10m"+
			"&timezone=Asia%%2FJakarta&forecast_days=2",
		info.Lat, info.Lon,
	)

	req, err := http.NewRequestWithContext(ctx, "GET", weatherURL, nil)
	if err != nil {
		return nil, fmt.Errorf("buat request cuaca: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch cuaca: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 2*1024*1024))
	if err != nil {
		return nil, fmt.Errorf("baca response cuaca: %w", err)
	}

	var om openMeteoResponse
	if err := json.Unmarshal(body, &om); err != nil {
		return nil, fmt.Errorf("parse response cuaca: %w", err)
	}

	var prakiraans []Prakiraan
	for i := 0; i < len(om.Hourly.Time) && i < 8; i++ {
		prakiraans = append(prakiraans, Prakiraan{
			Waktu:          om.Hourly.Time[i],
			Suhu:           fmt.Sprintf("%.0f°C", om.Hourly.Temperature2m[i]),
			Kelembapan:     fmt.Sprintf("%.0f%%", om.Hourly.RelativeHumidity[i]),
			Cuaca:          wmoToIndonesian(om.Hourly.WeatherCode[i]),
			KecepatanAngin: fmt.Sprintf("%.1f km/h", om.Hourly.WindSpeed10m[i]),
			ArahAngin:      windDegToDir(om.Hourly.WindDirection10m[i]),
		})
	}

	result := &CuacaResult{
		Provinsi:  provinsiKode,
		Kota:      info.Name,
		Prakiraan: prakiraans,
	}

	s.cache.Set(cacheKey, result, 30*time.Minute)
	return result, nil
}
