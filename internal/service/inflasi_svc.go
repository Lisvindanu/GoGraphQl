package service

type InflasiItem struct {
	Periode        string
	Bulan          string
	Tahun          int
	InflasiBulanan float64
	InflasiTahunan float64
	IHK            float64
}

// Data inflasi bulanan dari BPS (Badan Pusat Statistik)
var inflasiData = []InflasiItem{
	// 2024
	{Periode: "Jan-2024", Bulan: "Januari", Tahun: 2024, InflasiBulanan: 0.04, InflasiTahunan: 2.57, IHK: 106.30},
	{Periode: "Feb-2024", Bulan: "Februari", Tahun: 2024, InflasiBulanan: 0.37, InflasiTahunan: 2.75, IHK: 106.69},
	{Periode: "Mar-2024", Bulan: "Maret", Tahun: 2024, InflasiBulanan: 0.52, InflasiTahunan: 3.05, IHK: 107.25},
	{Periode: "Apr-2024", Bulan: "April", Tahun: 2024, InflasiBulanan: 0.25, InflasiTahunan: 3.00, IHK: 107.52},
	{Periode: "Mei-2024", Bulan: "Mei", Tahun: 2024, InflasiBulanan: -0.03, InflasiTahunan: 2.84, IHK: 107.49},
	{Periode: "Jun-2024", Bulan: "Juni", Tahun: 2024, InflasiBulanan: 0.08, InflasiTahunan: 2.51, IHK: 107.57},
	{Periode: "Jul-2024", Bulan: "Juli", Tahun: 2024, InflasiBulanan: 0.08, InflasiTahunan: 2.13, IHK: 107.66},
	{Periode: "Ags-2024", Bulan: "Agustus", Tahun: 2024, InflasiBulanan: -0.03, InflasiTahunan: 2.12, IHK: 107.62},
	{Periode: "Sep-2024", Bulan: "September", Tahun: 2024, InflasiBulanan: 0.16, InflasiTahunan: 1.84, IHK: 107.79},
	{Periode: "Okt-2024", Bulan: "Oktober", Tahun: 2024, InflasiBulanan: 0.08, InflasiTahunan: 1.71, IHK: 107.87},
	{Periode: "Nov-2024", Bulan: "November", Tahun: 2024, InflasiBulanan: 0.30, InflasiTahunan: 1.55, IHK: 108.19},
	{Periode: "Des-2024", Bulan: "Desember", Tahun: 2024, InflasiBulanan: 0.44, InflasiTahunan: 1.57, IHK: 108.67},
	// 2025
	{Periode: "Jan-2025", Bulan: "Januari", Tahun: 2025, InflasiBulanan: 0.76, InflasiTahunan: 0.76, IHK: 109.50},
	{Periode: "Feb-2025", Bulan: "Februari", Tahun: 2025, InflasiBulanan: 0.48, InflasiTahunan: 1.06, IHK: 110.02},
	{Periode: "Mar-2025", Bulan: "Maret", Tahun: 2025, InflasiBulanan: 1.65, InflasiTahunan: 1.03, IHK: 111.83},
	{Periode: "Apr-2025", Bulan: "April", Tahun: 2025, InflasiBulanan: -1.11, InflasiTahunan: 0.95, IHK: 110.60},
	{Periode: "Mei-2025", Bulan: "Mei", Tahun: 2025, InflasiBulanan: 0.17, InflasiTahunan: 0.92, IHK: 110.79},
	{Periode: "Jun-2025", Bulan: "Juni", Tahun: 2025, InflasiBulanan: 0.19, InflasiTahunan: 0.97, IHK: 111.00},
	{Periode: "Jul-2025", Bulan: "Juli", Tahun: 2025, InflasiBulanan: 0.24, InflasiTahunan: 1.13, IHK: 111.27},
	{Periode: "Ags-2025", Bulan: "Agustus", Tahun: 2025, InflasiBulanan: 0.21, InflasiTahunan: 1.35, IHK: 111.50},
}

type InflasiService struct{}

func NewInflasiService() *InflasiService { return &InflasiService{} }

func (s *InflasiService) GetAll(tahun *int) []InflasiItem {
	if tahun == nil {
		return inflasiData
	}
	result := make([]InflasiItem, 0)
	for _, item := range inflasiData {
		if item.Tahun == *tahun {
			result = append(result, item)
		}
	}
	return result
}
