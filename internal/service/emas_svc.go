package service

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type HargaEmasItem struct {
	Gram         float64
	HargaJual    int
	HargaBuyback int
}

type EmasService struct {
	client *http.Client
}

func NewEmasService() *EmasService {
	return &EmasService{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}


func parseIDRInt(s string) int {
	cleaned := strings.ReplaceAll(s, ".", "")
	cleaned = strings.ReplaceAll(cleaned, ",", "")
	n, _ := strconv.Atoi(cleaned)
	return n
}

func parseGram(s string) float64 {
	s = strings.ReplaceAll(s, ",", ".")
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func (s *EmasService) GetAll() ([]HargaEmasItem, error) {
	req, err := http.NewRequest("GET", "https://hargaemas.id/harga/antam/", nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch hargaemas.id: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	return parseEmasHTML(string(body))
}

// parseEmasHTML extracts gold price rows from hargaemas.id HTML.
// Grid structure: col-span-1 (gram) | col-span-2 (jual Rp...) | col-span-2 (buyback Rp...)
func parseEmasHTML(html string) ([]HargaEmasItem, error) {
	reStripTag := regexp.MustCompile(`<[^>]+>`)
	reRow := regexp.MustCompile(`(?s)col-span-1[^>]*>\s*([\d.,]+)\s*</div>\s*<div[^>]*col-span-2[^>]*>\s*Rp([\d.,]+)\s*</div>\s*<div[^>]*col-span-2[^>]*>\s*Rp([\d.,]+)\s*</div>`)

	matches := reRow.FindAllStringSubmatch(html, -1)
	result := make([]HargaEmasItem, 0, 12)

	for _, m := range matches {
		gramStr := strings.TrimSpace(reStripTag.ReplaceAllString(m[1], ""))
		gram := parseGram(gramStr)
		if gram <= 0 {
			continue
		}
		jual := parseIDRInt(m[2])
		buyback := parseIDRInt(m[3])
		if jual <= 0 || buyback <= 0 {
			continue
		}
		result = append(result, HargaEmasItem{
			Gram:         gram,
			HargaJual:    jual,
			HargaBuyback: buyback,
		})
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("tidak ada data harga emas yang berhasil diparsing")
	}
	return result, nil
}
