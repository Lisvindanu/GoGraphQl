package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

type hari struct {
	tanggal string
	nama    string
	jenis   string
}

// Data resmi berdasarkan SKB 3 Menteri.
var data = []hari{
	// ====== 2024 ======
	{"2024-01-01", "Tahun Baru Masehi 2024", "nasional"},
	{"2024-02-08", "Isra Mikraj Nabi Muhammad SAW", "nasional"},
	{"2024-02-09", "Cuti Bersama Isra Mikraj", "cuti_bersama"},
	{"2024-03-11", "Hari Suci Nyepi (Tahun Baru Saka 1946)", "nasional"},
	{"2024-03-29", "Wafat Isa Al Masih", "nasional"},
	{"2024-04-08", "Cuti Bersama Idul Fitri 1445 H", "cuti_bersama"},
	{"2024-04-09", "Cuti Bersama Idul Fitri 1445 H", "cuti_bersama"},
	{"2024-04-10", "Hari Raya Idul Fitri 1445 H", "nasional"},
	{"2024-04-11", "Hari Raya Idul Fitri 1445 H (Hari Kedua)", "nasional"},
	{"2024-04-12", "Cuti Bersama Idul Fitri 1445 H", "cuti_bersama"},
	{"2024-04-15", "Cuti Bersama Idul Fitri 1445 H", "cuti_bersama"},
	{"2024-05-01", "Hari Buruh Internasional", "nasional"},
	{"2024-05-09", "Kenaikan Isa Al Masih", "nasional"},
	{"2024-05-10", "Cuti Bersama Kenaikan Isa Al Masih", "cuti_bersama"},
	{"2024-05-23", "Hari Raya Waisak 2568 BE", "nasional"},
	{"2024-06-01", "Hari Lahir Pancasila", "nasional"},
	{"2024-06-17", "Hari Raya Idul Adha 1445 H", "nasional"},
	{"2024-07-07", "Tahun Baru Islam 1446 H", "nasional"},
	{"2024-08-17", "Hari Kemerdekaan Republik Indonesia", "nasional"},
	{"2024-09-16", "Maulid Nabi Muhammad SAW 1446 H", "nasional"},
	{"2024-12-25", "Hari Raya Natal", "nasional"},
	{"2024-12-26", "Cuti Bersama Natal", "cuti_bersama"},

	// ====== 2025 ======
	{"2025-01-01", "Tahun Baru Masehi 2025", "nasional"},
	{"2025-01-27", "Isra Mikraj Nabi Muhammad SAW 1446 H", "nasional"},
	{"2025-01-28", "Cuti Bersama Isra Mikraj", "cuti_bersama"},
	{"2025-01-29", "Tahun Baru Imlek 2576", "nasional"},
	{"2025-03-28", "Hari Suci Nyepi (Tahun Baru Saka 1947)", "nasional"},
	{"2025-03-30", "Hari Raya Idul Fitri 1446 H", "nasional"},
	{"2025-03-31", "Hari Raya Idul Fitri 1446 H (Hari Kedua)", "nasional"},
	{"2025-04-01", "Cuti Bersama Idul Fitri 1446 H", "cuti_bersama"},
	{"2025-04-02", "Cuti Bersama Idul Fitri 1446 H", "cuti_bersama"},
	{"2025-04-03", "Cuti Bersama Idul Fitri 1446 H", "cuti_bersama"},
	{"2025-04-04", "Cuti Bersama Idul Fitri 1446 H", "cuti_bersama"},
	{"2025-04-07", "Cuti Bersama Idul Fitri 1446 H", "cuti_bersama"},
	{"2025-04-18", "Wafat Isa Al Masih", "nasional"},
	{"2025-05-01", "Hari Buruh Internasional", "nasional"},
	{"2025-05-12", "Hari Raya Waisak 2569 BE", "nasional"},
	{"2025-05-29", "Kenaikan Isa Al Masih", "nasional"},
	{"2025-05-30", "Cuti Bersama Kenaikan Isa Al Masih", "cuti_bersama"},
	{"2025-06-01", "Hari Lahir Pancasila", "nasional"},
	{"2025-06-06", "Hari Raya Idul Adha 1446 H", "nasional"},
	{"2025-06-27", "Tahun Baru Islam 1447 H", "nasional"},
	{"2025-08-17", "Hari Kemerdekaan Republik Indonesia", "nasional"},
	{"2025-09-05", "Maulid Nabi Muhammad SAW 1447 H", "nasional"},
	{"2025-12-25", "Hari Raya Natal", "nasional"},
	{"2025-12-26", "Cuti Bersama Natal", "cuti_bersama"},
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

	inserted := 0
	skipped := 0
	for _, h := range data {
		t, err := time.Parse("2006-01-02", h.tanggal)
		if err != nil {
			log.Fatalf("invalid date %s: %v", h.tanggal, err)
		}

		tag, err := conn.Exec(ctx,
			`INSERT INTO hari_libur (tanggal, nama, tahun, jenis)
			 VALUES ($1, $2, $3, $4)
			 ON CONFLICT (tanggal, jenis) DO NOTHING`,
			t, h.nama, t.Year(), h.jenis,
		)
		if err != nil {
			log.Fatalf("insert %s: %v", h.tanggal, err)
		}
		if tag.RowsAffected() == 0 {
			skipped++
		} else {
			inserted++
		}
	}

	log.Printf("Done: %d inserted, %d skipped (already exists).", inserted, skipped)
}
