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

var (
	reGram  = regexp.MustCompile(`(?i)(\d+(?:[.,]\d+)?)\s*gr`)
	rePrice = regexp.MustCompile(`[\d.,]+`)
)

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
	req, err := http.NewRequest("GET", "https://www.logammulia.com/id/harga-emas-hari-ini", nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch logammulia: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	return parseEmasHTML(string(body))
}

// parseEmasHTML extracts gold price rows from logammulia.com HTML.
// Each row contains: Gram | Harga Jual | Harga Buyback
func parseEmasHTML(html string) ([]HargaEmasItem, error) {
	reRow := regexp.MustCompile(`(?s)<tr[^>]*>.*?</tr>`)
	reCell := regexp.MustCompile(`(?s)<td[^>]*>(.*?)</td>`)
	reStripTag := regexp.MustCompile(`<[^>]+>`)

	rows := reRow.FindAllString(html, -1)
	result := make([]HargaEmasItem, 0, 10)

	for _, row := range rows {
		cells := reCell.FindAllStringSubmatch(row, -1)
		if len(cells) < 3 {
			continue
		}

		cell0 := reStripTag.ReplaceAllString(cells[0][1], "")
		gramMatch := reGram.FindStringSubmatch(cell0)
		if gramMatch == nil {
			continue
		}
		gram := parseGram(gramMatch[1])
		if gram <= 0 {
			continue
		}

		cell1 := reStripTag.ReplaceAllString(cells[1][1], "")
		cell2 := reStripTag.ReplaceAllString(cells[2][1], "")

		priceMatch1 := rePrice.FindString(strings.TrimSpace(cell1))
		priceMatch2 := rePrice.FindString(strings.TrimSpace(cell2))
		if priceMatch1 == "" || priceMatch2 == "" {
			continue
		}

		jual := parseIDRInt(priceMatch1)
		buyback := parseIDRInt(priceMatch2)
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
