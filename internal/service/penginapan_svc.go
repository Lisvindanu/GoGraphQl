package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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

func (s *PenginapanService) GetByKota(ctx context.Context, kota string) ([]PenginapanItem, error) {
	kota = strings.TrimSpace(kota)
	if kota == "" {
		return nil, fmt.Errorf("nama kota tidak boleh kosong")
	}

	// Overpass QL: cari hotel/guest_house/hostel dalam area kota
	// Gunakan regex case-insensitive agar "bandung" cocok dengan "Kota Bandung"
	query := fmt.Sprintf(`[out:json][timeout:30];
(
  area["name"~"%s","i"]["boundary"="administrative"];
  area["name"~"%s","i"]["place"~"^(city|town|municipality)$"];
)->.searchArea;
(
  node["tourism"~"^(hotel|guest_house|hostel|motel|inn)$"](area.searchArea);
  way["tourism"~"^(hotel|guest_house|hostel|motel|inn)$"](area.searchArea);
);
out center tags 100;`, kota, kota)

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
