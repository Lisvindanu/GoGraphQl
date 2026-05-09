package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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

const comulineBaseURL = "https://www.api.comuline.com"

func (s *KRLService) GetStasiun() ([]StasiunKRL, error) {
	body, err := s.fetch(comulineBaseURL + "/v1/station")
	if err != nil {
		return nil, err
	}

	var root struct {
		Metadata struct {
			Success bool `json:"success"`
		} `json:"metadata"`
		Data []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &root); err != nil {
		return nil, fmt.Errorf("parse stasiun: %w", err)
	}

	result := make([]StasiunKRL, 0, len(root.Data))
	for _, d := range root.Data {
		if d.Type != "KRL" {
			continue
		}
		result = append(result, StasiunKRL{
			StasiunId:   d.ID,
			StasiunNama: d.Name,
			StasiunKode: d.ID,
		})
	}
	return result, nil
}

var wib = time.FixedZone("WIB", 7*3600)

func parseDestTime(departsAt string) string {
	// Try RFC3339Nano first (handles fractional seconds like "2024-03-10T09:55:07.213Z")
	if t, err := time.Parse(time.RFC3339Nano, departsAt); err == nil {
		return t.In(wib).Format("15:04")
	}
	if t, err := time.Parse(time.RFC3339, departsAt); err == nil {
		return t.In(wib).Format("15:04")
	}
	// Fallback: extract HH:mm from "...THH:mm:ss..."
	if idx := strings.Index(departsAt, "T"); idx >= 0 && len(departsAt) >= idx+6 {
		return departsAt[idx+1 : idx+6]
	}
	return departsAt
}

func (s *KRLService) GetJadwal(stasiunId, timeFrom, timeTo string) ([]JadwalKRL, error) {
	apiURL := fmt.Sprintf("%s/v1/schedule/%s", comulineBaseURL, stasiunId)

	body, err := s.fetch(apiURL)
	if err != nil {
		return nil, err
	}

	var root struct {
		Metadata struct {
			Success bool `json:"success"`
		} `json:"metadata"`
		Data []struct {
			TrainID              string `json:"train_id"`
			Line                 string `json:"line"`
			Route                string `json:"route"`
			DepartsAt            string `json:"departs_at"`
			StationDestinationID string `json:"station_destination_id"`
			Metadata             struct {
				Origin struct {
					Color *string `json:"color"`
				} `json:"origin"`
			} `json:"metadata"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &root); err != nil {
		return nil, fmt.Errorf("parse jadwal: %w", err)
	}

	result := make([]JadwalKRL, 0, len(root.Data))
	for _, d := range root.Data {
		destTime := parseDestTime(d.DepartsAt)

		if timeFrom != "" && destTime < timeFrom {
			continue
		}
		if timeTo != "" && destTime > timeTo {
			continue
		}

		colorCode := ""
		if d.Metadata.Origin.Color != nil {
			colorCode = *d.Metadata.Origin.Color
		}

		result = append(result, JadwalKRL{
			TrainId:     d.TrainID,
			KaName:      d.Line,
			RouteName:   d.Route,
			DestTime:    destTime,
			DestStasiun: d.StationDestinationID,
			ColorCode:   colorCode,
		})
	}
	return result, nil
}
