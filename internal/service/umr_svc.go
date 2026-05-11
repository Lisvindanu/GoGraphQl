package service

import "strings"

type UMRProvinsi struct {
	Provinsi string
	Kode     string
	Upah     int
	Tahun    int
}

// Data UMP 2025 berdasarkan Peraturan Pemerintah dan Keputusan Gubernur masing-masing provinsi.
var umrData2025 = []UMRProvinsi{
	{Provinsi: "Aceh", Kode: "11", Upah: 3685615, Tahun: 2025},
	{Provinsi: "Sumatera Utara", Kode: "12", Upah: 2890000, Tahun: 2025},
	{Provinsi: "Sumatera Barat", Kode: "13", Upah: 2994193, Tahun: 2025},
	{Provinsi: "Riau", Kode: "14", Upah: 3508775, Tahun: 2025},
	{Provinsi: "Jambi", Kode: "15", Upah: 3234533, Tahun: 2025},
	{Provinsi: "Sumatera Selatan", Kode: "16", Upah: 3681571, Tahun: 2025},
	{Provinsi: "Bengkulu", Kode: "17", Upah: 2670011, Tahun: 2025},
	{Provinsi: "Lampung", Kode: "18", Upah: 2893688, Tahun: 2025},
	{Provinsi: "Bangka Belitung", Kode: "19", Upah: 3872983, Tahun: 2025},
	{Provinsi: "Kepulauan Riau", Kode: "21", Upah: 3623654, Tahun: 2025},
	{Provinsi: "DKI Jakarta", Kode: "31", Upah: 5396761, Tahun: 2025},
	{Provinsi: "Jawa Barat", Kode: "32", Upah: 2191232, Tahun: 2025},
	{Provinsi: "Jawa Tengah", Kode: "33", Upah: 2169348, Tahun: 2025},
	{Provinsi: "DI Yogyakarta", Kode: "34", Upah: 2264080, Tahun: 2025},
	{Provinsi: "Jawa Timur", Kode: "35", Upah: 2305984, Tahun: 2025},
	{Provinsi: "Banten", Kode: "36", Upah: 2727812, Tahun: 2025},
	{Provinsi: "Bali", Kode: "51", Upah: 2996560, Tahun: 2025},
	{Provinsi: "Nusa Tenggara Barat", Kode: "52", Upah: 2628780, Tahun: 2025},
	{Provinsi: "Nusa Tenggara Timur", Kode: "53", Upah: 2328555, Tahun: 2025},
	{Provinsi: "Kalimantan Barat", Kode: "61", Upah: 2878285, Tahun: 2025},
	{Provinsi: "Kalimantan Tengah", Kode: "62", Upah: 3473621, Tahun: 2025},
	{Provinsi: "Kalimantan Selatan", Kode: "63", Upah: 3496194, Tahun: 2025},
	{Provinsi: "Kalimantan Timur", Kode: "64", Upah: 3579313, Tahun: 2025},
	{Provinsi: "Kalimantan Utara", Kode: "65", Upah: 3580160, Tahun: 2025},
	{Provinsi: "Sulawesi Utara", Kode: "71", Upah: 3775425, Tahun: 2025},
	{Provinsi: "Sulawesi Tengah", Kode: "72", Upah: 2914583, Tahun: 2025},
	{Provinsi: "Sulawesi Selatan", Kode: "73", Upah: 3657527, Tahun: 2025},
	{Provinsi: "Sulawesi Tenggara", Kode: "74", Upah: 3073551, Tahun: 2025},
	{Provinsi: "Gorontalo", Kode: "75", Upah: 3221731, Tahun: 2025},
	{Provinsi: "Sulawesi Barat", Kode: "76", Upah: 2871794, Tahun: 2025},
	{Provinsi: "Maluku", Kode: "81", Upah: 3141699, Tahun: 2025},
	{Provinsi: "Maluku Utara", Kode: "82", Upah: 3408000, Tahun: 2025},
	{Provinsi: "Papua Barat", Kode: "91", Upah: 3837000, Tahun: 2025},
	{Provinsi: "Papua", Kode: "92", Upah: 4285848, Tahun: 2025},
	{Provinsi: "Papua Selatan", Kode: "93", Upah: 4285848, Tahun: 2025},
	{Provinsi: "Papua Tengah", Kode: "94", Upah: 4285848, Tahun: 2025},
	{Provinsi: "Papua Pegunungan", Kode: "95", Upah: 4285848, Tahun: 2025},
	{Provinsi: "Papua Barat Daya", Kode: "96", Upah: 3975000, Tahun: 2025},
}

type UMRService struct{}

func NewUMRService() *UMRService { return &UMRService{} }

func (s *UMRService) GetAll(tahun *int) []UMRProvinsi {
	if tahun == nil || *tahun == 2025 {
		return umrData2025
	}
	return []UMRProvinsi{}
}

func (s *UMRService) GetByProvinsi(query string) *UMRProvinsi {
	q := strings.ToLower(strings.TrimSpace(query))
	for _, item := range umrData2025 {
		if item.Kode == query || strings.ToLower(item.Provinsi) == q ||
			strings.Contains(strings.ToLower(item.Provinsi), q) {
			return &item
		}
	}
	return nil
}
