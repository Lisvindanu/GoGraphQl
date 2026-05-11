package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type PenginapanItem struct {
	ID      string
	Nama    string
	Tipe    string
	Alamat  *string
	Bintang *string
	Telepon *string
	Website *string
	Lat     float64
	Lon     float64
}

type overpassResponse struct {
	Elements []overpassElement `json:"elements"`
}

type overpassElement struct {
	Type   string            `json:"type"`
	ID     int64             `json:"id"`
	Lat    float64           `json:"lat"`
	Lon    float64           `json:"lon"`
	Center *overpassCenter   `json:"center"`
	Tags   map[string]string `json:"tags"`
}

type overpassCenter struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type PenginapanService struct {
	client *http.Client
}

func NewPenginapanService() *PenginapanService {
	return &PenginapanService{
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

type nominatimResult struct {
	BoundingBox []string `json:"boundingbox"`
}

func (s *PenginapanService) geocodeBbox(kota string) (south, west, north, east float64, err error) {
	reqURL := fmt.Sprintf(
		"https://nominatim.openstreetmap.org/search?q=%s&countrycodes=id&format=json&limit=1",
		url.QueryEscape(kota),
	)
	req, _ := http.NewRequest("GET", reqURL, nil)
	req.Header.Set("User-Agent", "IndonesiaQL/1.0")

	resp, err := s.client.Do(req)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("nominatim: %w", err)
	}
	defer resp.Body.Close()

	var results []nominatimResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil || len(results) == 0 {
		return 0, 0, 0, 0, fmt.Errorf("kota tidak ditemukan: %s", kota)
	}

	bb := results[0].BoundingBox // [south, north, west, east]
	if len(bb) < 4 {
		return 0, 0, 0, 0, fmt.Errorf("bbox tidak valid")
	}
	parse := func(s string) float64 {
		v, _ := strconv.ParseFloat(s, 64)
		return v
	}
	return parse(bb[0]), parse(bb[2]), parse(bb[1]), parse(bb[3]), nil
}

func (s *PenginapanService) GetByKota(ctx context.Context, kota string) ([]PenginapanItem, error) {
	kota = strings.TrimSpace(kota)
	if kota == "" {
		return nil, fmt.Errorf("nama kota tidak boleh kosong")
	}

	south, west, north, east, err := s.geocodeBbox(kota)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`[out:json][timeout:30];
(
  node["tourism"~"^(hotel|guest_house|hostel|motel|inn)$"](%f,%f,%f,%f);
  way["tourism"~"^(hotel|guest_house|hostel|motel|inn)$"](%f,%f,%f,%f);
);
out center tags 100;`, south, west, north, east, south, west, north, east)

	apiURL := "https://overpass-api.de/api/interpreter"
	resp, err := s.client.PostForm(apiURL, url.Values{"data": {query}})
	if err != nil {
		return nil, fmt.Errorf("gagal menghubungi Overpass API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Overpass API error: status %d", resp.StatusCode)
	}

	var result overpassResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("gagal parse respons Overpass: %w", err)
	}

	items := make([]PenginapanItem, 0, len(result.Elements))
	for _, el := range result.Elements {
		nama, ok := el.Tags["name"]
		if !ok || strings.TrimSpace(nama) == "" {
			continue
		}

		lat, lon := el.Lat, el.Lon
		if el.Center != nil {
			lat = el.Center.Lat
			lon = el.Center.Lon
		}
		if lat == 0 && lon == 0 {
			continue
		}

		tipe := el.Tags["tourism"]

		item := PenginapanItem{
			ID:   fmt.Sprintf("%s/%d", el.Type, el.ID),
			Nama: nama,
			Tipe: tipe,
			Lat:  lat,
			Lon:  lon,
		}

		if v := buildAlamat(el.Tags); v != "" {
			item.Alamat = &v
		}
		if v, ok := el.Tags["stars"]; ok && v != "" {
			item.Bintang = &v
		}
		if v, ok := el.Tags["phone"]; ok && v != "" {
			item.Telepon = &v
		} else if v, ok := el.Tags["contact:phone"]; ok && v != "" {
			item.Telepon = &v
		}
		if v, ok := el.Tags["website"]; ok && v != "" {
			item.Website = &v
		} else if v, ok := el.Tags["contact:website"]; ok && v != "" {
			item.Website = &v
		}

		items = append(items, item)
	}

	return items, nil
}

func buildAlamat(tags map[string]string) string {
	if v, ok := tags["addr:full"]; ok && v != "" {
		return v
	}
	parts := []string{}
	if v := tags["addr:street"]; v != "" {
		if no := tags["addr:housenumber"]; no != "" {
			parts = append(parts, v+" "+no)
		} else {
			parts = append(parts, v)
		}
	}
	if v := tags["addr:city"]; v != "" {
		parts = append(parts, v)
	}
	return strings.Join(parts, ", ")
}
