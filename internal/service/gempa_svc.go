package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type GempaItem struct {
	Tanggal  string
	Jam      string
	Magnitude string
	Kedalaman string
	Lintang  string
	Bujur    string
	Wilayah  string
	Potensi  string
	Dirasakan string
}

type GempaService struct {
	client *http.Client
}

func NewGempaService() *GempaService {
	return &GempaService{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *GempaService) parseItem(raw map[string]interface{}) GempaItem {
	str := func(key string) string {
		if v, ok := raw[key].(string); ok {
			return v
		}
		return ""
	}
	return GempaItem{
		Tanggal:   str("Tanggal"),
		Jam:       str("Jam"),
		Magnitude: str("Magnitude"),
		Kedalaman: str("Kedalaman"),
		Lintang:   str("Lintang"),
		Bujur:     str("Bujur"),
		Wilayah:   str("Wilayah"),
		Potensi:   str("Potensi"),
		Dirasakan: str("Dirasakan"),
	}
}

func (s *GempaService) fetch(url string) ([]byte, error) {
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (s *GempaService) GetTerbaru() (*GempaItem, error) {
	body, err := s.fetch("https://data.bmkg.go.id/DataMKG/TEWS/autogempa.json")
	if err != nil {
		return nil, fmt.Errorf("fetch BMKG: %w", err)
	}

	var root struct {
		Infogempa struct {
			Gempa map[string]interface{} `json:"gempa"`
		} `json:"Infogempa"`
	}
	if err := json.Unmarshal(body, &root); err != nil {
		return nil, err
	}

	item := s.parseItem(root.Infogempa.Gempa)
	return &item, nil
}

func (s *GempaService) GetList() ([]GempaItem, error) {
	body, err := s.fetch("https://data.bmkg.go.id/DataMKG/TEWS/gempaterkini.json")
	if err != nil {
		return nil, fmt.Errorf("fetch BMKG: %w", err)
	}

	var root struct {
		Infogempa struct {
			Gempa []map[string]interface{} `json:"gempa"`
		} `json:"Infogempa"`
	}
	if err := json.Unmarshal(body, &root); err != nil {
		return nil, err
	}

	result := make([]GempaItem, 0, len(root.Infogempa.Gempa))
	for _, raw := range root.Infogempa.Gempa {
		result = append(result, s.parseItem(raw))
	}
	return result, nil
}
