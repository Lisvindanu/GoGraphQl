package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

const emsifaBase = "https://emsifa.github.io/api-wilayah-indonesia/api"

type province struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type regency struct {
	ID         string `json:"id"`
	ProvinceID string `json:"province_id"`
	Name       string `json:"name"`
}

type district struct {
	ID        string `json:"id"`
	RegencyID string `json:"regency_id"`
	Name      string `json:"name"`
}

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
	withKecamatan := flag.Bool("with-kecamatan", true, "seed kecamatan data")
	withKelurahan := flag.Bool("with-kelurahan", false, "seed kelurahan data (slow, ~7000 HTTP requests)")
	flag.Parse()

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

	// --- Provinsi ---
	log.Println("Fetching provinces...")
	var provinces []province
	if err := fetchJSON(emsifaBase+"/provinces.json", &provinces); err != nil {
		log.Fatal("fetch provinces:", err)
	}
	provRows := make([][]any, len(provinces))
	for i, p := range provinces {
		provRows[i] = []any{p.ID, p.Name}
	}
	n, err := conn.CopyFrom(ctx, pgx.Identifier{"provinsi"}, []string{"kode", "nama"},
		pgx.CopyFromRows(provRows))
	if err != nil {
		log.Fatal("insert provinsi:", err)
	}
	log.Printf("  inserted %d provinsi\n", n)

	// --- Kota ---
	log.Println("Fetching regencies (kota)...")
	var allRegencies []regency
	for _, p := range provinces {
		var regs []regency
		if err := fetchJSON(fmt.Sprintf("%s/regencies/%s.json", emsifaBase, p.ID), &regs); err != nil {
			log.Printf("  WARN: fetch regencies for province %s: %v", p.ID, err)
			continue
		}
		allRegencies = append(allRegencies, regs...)
	}
	regRows := make([][]any, len(allRegencies))
	for i, r := range allRegencies {
		regRows[i] = []any{r.ID, r.ProvinceID, r.Name}
	}
	n, err = conn.CopyFrom(ctx, pgx.Identifier{"kota"}, []string{"kode", "provinsi_kode", "nama"},
		pgx.CopyFromRows(regRows))
	if err != nil {
		log.Fatal("insert kota:", err)
	}
	log.Printf("  inserted %d kota\n", n)

	if !*withKecamatan {
		log.Println("Done. Skipping kecamatan (use --with-kecamatan).")
		return
	}

	// --- Kecamatan ---
	log.Printf("Fetching districts (kecamatan) for %d kota (a few minutes)...\n", len(allRegencies))
	var allDistricts []district
	for i, r := range allRegencies {
		var dists []district
		if err := fetchJSON(fmt.Sprintf("%s/districts/%s.json", emsifaBase, r.ID), &dists); err != nil {
			log.Printf("  WARN: fetch districts for kota %s: %v", r.ID, err)
			continue
		}
		allDistricts = append(allDistricts, dists...)
		if (i+1)%50 == 0 {
			log.Printf("  progress: %d/%d kota processed...", i+1, len(allRegencies))
		}
	}
	distRows := make([][]any, len(allDistricts))
	for i, d := range allDistricts {
		distRows[i] = []any{d.ID, d.RegencyID, d.Name}
	}
	n, err = conn.CopyFrom(ctx, pgx.Identifier{"kecamatan"}, []string{"kode", "kota_kode", "nama"},
		pgx.CopyFromRows(distRows))
	if err != nil {
		log.Fatal("insert kecamatan:", err)
	}
	log.Printf("  inserted %d kecamatan\n", n)

	if !*withKelurahan {
		log.Println("Done. Skip kelurahan (use --with-kelurahan for full data).")
		return
	}

	// --- Kelurahan ---
	log.Printf("Fetching villages (kelurahan) for %d kecamatan (very slow)...\n", len(allDistricts))
	batch := make([][]any, 0, 1000)
	flush := func() error {
		if len(batch) == 0 {
			return nil
		}
		_, err := conn.CopyFrom(ctx, pgx.Identifier{"kelurahan"},
			[]string{"kode", "kecamatan_kode", "nama", "kode_pos"},
			pgx.CopyFromRows(batch))
		batch = batch[:0]
		return err
	}
	for i, d := range allDistricts {
		var villages []village
		if err := fetchJSON(fmt.Sprintf("%s/villages/%s.json", emsifaBase, d.ID), &villages); err != nil {
			log.Printf("  WARN: fetch villages for kecamatan %s: %v", d.ID, err)
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
				log.Fatal("insert kelurahan batch:", err)
			}
		}
		if (i+1)%500 == 0 {
			log.Printf("  progress: %d/%d kecamatan processed...", i+1, len(allDistricts))
		}
	}
	if err := flush(); err != nil {
		log.Fatal("insert kelurahan final:", err)
	}
	log.Println("Done. All wilayah data seeded.")
}

