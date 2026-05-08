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

type kotaInfo struct {
	Nama     string
	Provinsi string
	Lat      float64
	Lon      float64
}

// kotaDatabase maps Indonesian city names to coordinates and provinces.
// Koordinat berdasarkan data BPS/OpenStreetMap.
var kotaDatabase = []kotaInfo{
	{"Jakarta", "DKI Jakarta", -6.2088, 106.8456},
	{"Jakarta Pusat", "DKI Jakarta", -6.1765, 106.8227},
	{"Jakarta Selatan", "DKI Jakarta", -6.2615, 106.8106},
	{"Jakarta Utara", "DKI Jakarta", -6.1167, 106.8833},
	{"Jakarta Barat", "DKI Jakarta", -6.1667, 106.7500},
	{"Jakarta Timur", "DKI Jakarta", -6.2250, 106.9004},
	{"Surabaya", "Jawa Timur", -7.2575, 112.7521},
	{"Bandung", "Jawa Barat", -6.9175, 107.6191},
	{"Medan", "Sumatera Utara", 3.5952, 98.6722},
	{"Semarang", "Jawa Tengah", -6.9932, 110.4203},
	{"Makassar", "Sulawesi Selatan", -5.1477, 119.4327},
	{"Palembang", "Sumatera Selatan", -2.9761, 104.7754},
	{"Tangerang", "Banten", -6.1702, 106.6400},
	{"Depok", "Jawa Barat", -6.4025, 106.7942},
	{"Bekasi", "Jawa Barat", -6.2383, 106.9756},
	{"Bogor", "Jawa Barat", -6.5971, 106.8060},
	{"Denpasar", "Bali", -8.6500, 115.2167},
	{"Yogyakarta", "DI Yogyakarta", -7.7972, 110.3688},
	{"Malang", "Jawa Timur", -7.9839, 112.6214},
	{"Padang", "Sumatera Barat", -0.9471, 100.4172},
	{"Pekanbaru", "Riau", 0.5333, 101.4500},
	{"Banjarmasin", "Kalimantan Selatan", -3.3242, 114.5919},
	{"Samarinda", "Kalimantan Timur", -0.5022, 117.1536},
	{"Balikpapan", "Kalimantan Timur", -1.2675, 116.8289},
	{"Manado", "Sulawesi Utara", 1.4748, 124.8421},
	{"Ambon", "Maluku", -3.6954, 128.1814},
	{"Jayapura", "Papua", -2.5337, 140.7181},
	{"Pontianak", "Kalimantan Barat", -0.0263, 109.3425},
	{"Kupang", "Nusa Tenggara Timur", -10.1772, 123.6070},
	{"Mataram", "Nusa Tenggara Barat", -8.5833, 116.1167},
	{"Banda Aceh", "Aceh", 5.5483, 95.3238},
	{"Bengkulu", "Bengkulu", -3.7928, 102.2608},
	{"Jambi", "Jambi", -1.6101, 103.6131},
	{"Bandar Lampung", "Lampung", -5.4292, 105.2613},
	{"Serang", "Banten", -6.1100, 106.1500},
	{"Pangkalpinang", "Kepulauan Bangka Belitung", -2.1347, 106.1167},
	{"Tanjungpinang", "Kepulauan Riau", 0.9167, 104.4500},
	{"Tarakan", "Kalimantan Utara", 3.3333, 117.6333},
	{"Palangkaraya", "Kalimantan Tengah", -2.2161, 113.9162},
	{"Gorontalo", "Gorontalo", 0.5333, 123.0600},
	{"Palu", "Sulawesi Tengah", -0.9003, 119.8779},
	{"Kendari", "Sulawesi Tenggara", -3.9721, 122.5139},
	{"Mamuju", "Sulawesi Barat", -2.6811, 118.8886},
	{"Ternate", "Maluku Utara", 0.7833, 127.3833},
	{"Sorong", "Papua Barat", -0.8667, 131.2500},
	{"Manokwari", "Papua Barat", -0.8667, 134.0667},
	{"Cirebon", "Jawa Barat", -6.7320, 108.5523},
	{"Solo", "Jawa Tengah", -7.5755, 110.8242},
	{"Surakarta", "Jawa Tengah", -7.5755, 110.8242},
	{"Madiun", "Jawa Timur", -7.6298, 111.5239},
	{"Kediri", "Jawa Timur", -7.8166, 112.0112},
	{"Batam", "Kepulauan Riau", 1.0456, 104.0305},
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

func (s *CuacaService) GetCuaca(ctx context.Context, provinsiKode, kotaQuery string) (*CuacaResult, error) {
	cacheKey := "cuaca:" + provinsiKode + ":" + strings.ToLower(kotaQuery)
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached.(*CuacaResult), nil
	}

	var info *kotaInfo
	queryLower := strings.ToLower(kotaQuery)
	for i := range kotaDatabase {
		if strings.Contains(strings.ToLower(kotaDatabase[i].Nama), queryLower) {
			info = &kotaDatabase[i]
			break
		}
	}
	if info == nil {
		return nil, fmt.Errorf("kota '%s' tidak ditemukan", kotaQuery)
	}

	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f"+
			"&hourly=temperature_2m,relative_humidity_2m,weather_code,wind_speed_10m,wind_direction_10m"+
			"&timezone=Asia%%2FJakarta&forecast_days=2",
		info.Lat, info.Lon,
	)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
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
		Provinsi:  info.Provinsi,
		Kota:      info.Nama,
		Prakiraan: prakiraans,
	}

	s.cache.Set(cacheKey, result, 30*time.Minute)
	return result, nil
}
