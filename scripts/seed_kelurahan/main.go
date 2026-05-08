package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

const emsifaBase = "https://emsifa.github.io/api-wilayah-indonesia/api"

type village struct {
	ID         string `json:"id"`
	DistrictID string `json:"district_id"`
	Name       string `json:"name"`
	PostalCode string `json:"postal_code"`
}

var client = &http.Client{Timeout: 30 * time.Second}

func fetchJSON(url string, out any) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, url)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL env var is required")
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		log.Fatal("connect db:", err)
	}
	defer conn.Close(ctx)

	rows, err := conn.Query(ctx, "SELECT kode FROM kecamatan ORDER BY kode")
	if err != nil {
		log.Fatal("query kecamatan:", err)
	}
	var kecamatanList []string
	for rows.Next() {
		var kode string
		if err := rows.Scan(&kode); err != nil {
			log.Fatal("scan:", err)
		}
		kecamatanList = append(kecamatanList, kode)
	}
	rows.Close()
	log.Printf("Found %d kecamatan\n", len(kecamatanList))

	batch := make([][]any, 0, 1000)
	total := 0
	flush := func() error {
		if len(batch) == 0 {
			return nil
		}
		_, err := conn.CopyFrom(ctx, pgx.Identifier{"kelurahan"},
			[]string{"kode", "kecamatan_kode", "nama", "kode_pos"},
			pgx.CopyFromRows(batch))
		if err != nil {
			return err
		}
		total += len(batch)
		batch = batch[:0]
		return nil
	}

	for i, kode := range kecamatanList {
		var villages []village
		url := fmt.Sprintf("%s/villages/%s.json", emsifaBase, kode)
		if err := fetchJSON(url, &villages); err != nil {
			log.Printf("  WARN: kecamatan %s: %v", kode, err)
			continue
		}
		for _, v := range villages {
			var kodePos *string
			if v.PostalCode != "" {
				kodePos = &v.PostalCode
			}
			batch = append(batch, []any{v.ID, v.DistrictID, v.Name, kodePos})
		}
		if len(batch) >= 1000 {
			if err := flush(); err != nil {
				log.Fatal("insert batch:", err)
			}
		}
		if (i+1)%500 == 0 {
			log.Printf("  %d/%d kecamatan, %d kelurahan inserted...", i+1, len(kecamatanList), total)
		}
	}
	if err := flush(); err != nil {
		log.Fatal("insert final:", err)
	}
	log.Printf("Done. %d kelurahan total.\n", total)
}
