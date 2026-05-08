package domain

import "time"

type HariLibur struct {
	ID      int
	Tanggal time.Time
	Nama    string
	Tahun   int
	Jenis   string
}
