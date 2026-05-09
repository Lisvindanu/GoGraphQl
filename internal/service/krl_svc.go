package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type StasiunKRL struct {
	StasiunId   string
	StasiunNama string
	StasiunKode string
}

type JadwalKRL struct {
	TrainId     string
	KaName      string
	RouteName   string
	DestTime    string
	DestStasiun string
	ColorCode   string
}

type KRLService struct {
	client *http.Client
}

func NewKRLService() *KRLService {
	return &KRLService{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *KRLService) fetch(apiURL string) ([]byte, error) {
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch KRL API: %w", err)
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (s *KRLService) GetStasiun() ([]StasiunKRL, error) {
	body, err := s.fetch("https://api-partner.krl.co.id/krlweb/v1/stasiun?status=all")
	if err != nil {
		return nil, err
	}

	var root struct {
		Status int `json:"status"`
		Data   []struct {
			StasiunID   string `json:"stasiun_id"`
			StasiunNama string `json:"stasiun_nama"`
			StasiunKode string `json:"stasiun_kode"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &root); err != nil {
		return nil, fmt.Errorf("parse stasiun: %w", err)
	}

	result := make([]StasiunKRL, 0, len(root.Data))
	for _, d := range root.Data {
		result = append(result, StasiunKRL{
			StasiunId:   d.StasiunID,
			StasiunNama: d.StasiunNama,
			StasiunKode: d.StasiunKode,
		})
	}
	return result, nil
}

func (s *KRLService) GetJadwal(stasiunId, timeFrom, timeTo string) ([]JadwalKRL, error) {
	if timeFrom == "" {
		timeFrom = "00:00"
	}
	if timeTo == "" {
		timeTo = "23:59"
	}

	apiURL := fmt.Sprintf(
		"https://api-partner.krl.co.id/krlweb/v1/schedule?stasiunid=%s&timeFrom=%s&timeTo=%s",
		stasiunId, timeFrom, timeTo,
	)

	body, err := s.fetch(apiURL)
	if err != nil {
		return nil, err
	}

	var root struct {
		Status int `json:"status"`
		Data   []struct {
			TrainID         string `json:"train_id"`
			KaName          string `json:"ka_name"`
			RouteName       string `json:"route_name"`
			DestTime        string `json:"dest_time"`
			DestStationName string `json:"dest_station_name"`
			ColorCode       string `json:"color_code"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &root); err != nil {
		return nil, fmt.Errorf("parse jadwal: %w", err)
	}

	result := make([]JadwalKRL, 0, len(root.Data))
	for _, d := range root.Data {
		result = append(result, JadwalKRL{
			TrainId:     d.TrainID,
			KaName:      d.KaName,
			RouteName:   d.RouteName,
			DestTime:    d.DestTime,
			DestStasiun: d.DestStationName,
			ColorCode:   d.ColorCode,
		})
	}
	return result, nil
}
