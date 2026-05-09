package service

type HargaBBMItem struct {
	Nama   string
	Harga  int
	Satuan string
	Jenis  string
}

var hargaBBMData = []HargaBBMItem{
	{Nama: "Pertalite", Harga: 10000, Satuan: "per liter", Jenis: "Subsidi"},
	{Nama: "Pertamax", Harga: 13500, Satuan: "per liter", Jenis: "Non-Subsidi"},
	{Nama: "Pertamax Green 95", Harga: 13500, Satuan: "per liter", Jenis: "Non-Subsidi"},
	{Nama: "Pertamax Turbo", Harga: 14000, Satuan: "per liter", Jenis: "Non-Subsidi"},
	{Nama: "Pertamax Racing", Harga: 38000, Satuan: "per liter", Jenis: "Non-Subsidi"},
	{Nama: "Solar Subsidi (Biosolar)", Harga: 6800, Satuan: "per liter", Jenis: "Subsidi"},
	{Nama: "Dexlite", Harga: 13500, Satuan: "per liter", Jenis: "Non-Subsidi"},
	{Nama: "Pertamina Dex", Harga: 13950, Satuan: "per liter", Jenis: "Non-Subsidi"},
}

type HargaBBMService struct{}

func NewHargaBBMService() *HargaBBMService { return &HargaBBMService{} }

func (s *HargaBBMService) GetAll() []HargaBBMItem {
	return hargaBBMData
}
