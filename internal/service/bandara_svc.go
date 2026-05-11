package service

import "strings"

type BandaraItem struct {
	KodeIATA string
	Nama     string
	Kota     string
	Provinsi string
}

var bandaraData = []BandaraItem{
	// Jawa
	{KodeIATA: "CGK", Nama: "Soekarno-Hatta International", Kota: "Tangerang", Provinsi: "Banten"},
	{KodeIATA: "HLP", Nama: "Halim Perdanakusuma International", Kota: "Jakarta", Provinsi: "DKI Jakarta"},
	{KodeIATA: "SUB", Nama: "Juanda International", Kota: "Surabaya", Provinsi: "Jawa Timur"},
	{KodeIATA: "BDO", Nama: "Husein Sastranegara International", Kota: "Bandung", Provinsi: "Jawa Barat"},
	{KodeIATA: "SRG", Nama: "Ahmad Yani International", Kota: "Semarang", Provinsi: "Jawa Tengah"},
	{KodeIATA: "JOG", Nama: "Adisucipto International", Kota: "Yogyakarta", Provinsi: "DI Yogyakarta"},
	{KodeIATA: "YIA", Nama: "Yogyakarta International", Kota: "Kulon Progo", Provinsi: "DI Yogyakarta"},
	{KodeIATA: "SOC", Nama: "Adisumarmo International", Kota: "Solo", Provinsi: "Jawa Tengah"},
	{KodeIATA: "MLG", Nama: "Abdul Rachman Saleh", Kota: "Malang", Provinsi: "Jawa Timur"},
	{KodeIATA: "JBB", Nama: "Notohadinegoro", Kota: "Jember", Provinsi: "Jawa Timur"},
	// Bali & Nusa Tenggara
	{KodeIATA: "DPS", Nama: "Ngurah Rai International", Kota: "Denpasar", Provinsi: "Bali"},
	{KodeIATA: "LOP", Nama: "Lombok International", Kota: "Praya", Provinsi: "Nusa Tenggara Barat"},
	{KodeIATA: "KOE", Nama: "El Tari International", Kota: "Kupang", Provinsi: "Nusa Tenggara Timur"},
	{KodeIATA: "MOF", Nama: "Frans Xavier Seda", Kota: "Maumere", Provinsi: "Nusa Tenggara Timur"},
	// Sumatera
	{KodeIATA: "KNO", Nama: "Kualanamu International", Kota: "Deli Serdang", Provinsi: "Sumatera Utara"},
	{KodeIATA: "BTJ", Nama: "Sultan Iskandar Muda International", Kota: "Banda Aceh", Provinsi: "Aceh"},
	{KodeIATA: "PDG", Nama: "Minangkabau International", Kota: "Padang Pariaman", Provinsi: "Sumatera Barat"},
	{KodeIATA: "PKU", Nama: "Sultan Syarif Kasim II International", Kota: "Pekanbaru", Provinsi: "Riau"},
	{KodeIATA: "BTH", Nama: "Hang Nadim International", Kota: "Batam", Provinsi: "Kepulauan Riau"},
	{KodeIATA: "TNJ", Nama: "Raja Haji Fisabilillah International", Kota: "Tanjungpinang", Provinsi: "Kepulauan Riau"},
	{KodeIATA: "DJB", Nama: "Sultan Thaha", Kota: "Jambi", Provinsi: "Jambi"},
	{KodeIATA: "PLM", Nama: "Sultan Mahmud Badaruddin II International", Kota: "Palembang", Provinsi: "Sumatera Selatan"},
	{KodeIATA: "BKS", Nama: "Fatmawati Soekarno", Kota: "Bengkulu", Provinsi: "Bengkulu"},
	{KodeIATA: "TKG", Nama: "Radin Inten II International", Kota: "Bandar Lampung", Provinsi: "Lampung"},
	{KodeIATA: "PGK", Nama: "Depati Amir", Kota: "Pangkalpinang", Provinsi: "Bangka Belitung"},
	// Kalimantan
	{KodeIATA: "PNK", Nama: "Supadio International", Kota: "Pontianak", Provinsi: "Kalimantan Barat"},
	{KodeIATA: "PKY", Nama: "Tjilik Riwut", Kota: "Palangka Raya", Provinsi: "Kalimantan Tengah"},
	{KodeIATA: "BDJ", Nama: "Syamsuddin Noor International", Kota: "Banjarbaru", Provinsi: "Kalimantan Selatan"},
	{KodeIATA: "BPN", Nama: "Sultan Aji Muhammad Sulaiman Sepinggan International", Kota: "Balikpapan", Provinsi: "Kalimantan Timur"},
	{KodeIATA: "AAP", Nama: "Samarinda Baru", Kota: "Samarinda", Provinsi: "Kalimantan Timur"},
	{KodeIATA: "TRK", Nama: "Juwata International", Kota: "Tarakan", Provinsi: "Kalimantan Utara"},
	// Sulawesi
	{KodeIATA: "UPG", Nama: "Sultan Hasanuddin International", Kota: "Makassar", Provinsi: "Sulawesi Selatan"},
	{KodeIATA: "MDC", Nama: "Sam Ratulangi International", Kota: "Manado", Provinsi: "Sulawesi Utara"},
	{KodeIATA: "PLW", Nama: "Mutiara SIS Al-Jufrie", Kota: "Palu", Provinsi: "Sulawesi Tengah"},
	{KodeIATA: "KDI", Nama: "Haluoleo International", Kota: "Kendari", Provinsi: "Sulawesi Tenggara"},
	{KodeIATA: "GTO", Nama: "Djalaluddin", Kota: "Gorontalo", Provinsi: "Gorontalo"},
	{KodeIATA: "MJU", Nama: "Tampa Padang", Kota: "Mamuju", Provinsi: "Sulawesi Barat"},
	// Maluku
	{KodeIATA: "AMQ", Nama: "Pattimura International", Kota: "Ambon", Provinsi: "Maluku"},
	{KodeIATA: "TTE", Nama: "Sultan Baabullah", Kota: "Ternate", Provinsi: "Maluku Utara"},
	// Papua
	{KodeIATA: "DJJ", Nama: "Sentani International", Kota: "Jayapura", Provinsi: "Papua"},
	{KodeIATA: "BIK", Nama: "Frans Kaisiepo International", Kota: "Biak", Provinsi: "Papua"},
	{KodeIATA: "TIM", Nama: "Moses Kilangin", Kota: "Timika", Provinsi: "Papua Tengah"},
	{KodeIATA: "WMX", Nama: "Rendani", Kota: "Manokwari", Provinsi: "Papua Barat"},
	{KodeIATA: "SOQ", Nama: "Domine Eduard Osok", Kota: "Sorong", Provinsi: "Papua Barat Daya"},
	{KodeIATA: "WAR", Nama: "Waris", Kota: "Wamena", Provinsi: "Papua Pegunungan"},
	{KodeIATA: "MKW", Nama: "Marinda", Kota: "Raja Ampat", Provinsi: "Papua Barat Daya"},
}

type BandaraService struct{}

func NewBandaraService() *BandaraService { return &BandaraService{} }

func (s *BandaraService) GetAll() []BandaraItem {
	return bandaraData
}

func (s *BandaraService) GetByKode(kode string) *BandaraItem {
	kode = strings.ToUpper(strings.TrimSpace(kode))
	for _, item := range bandaraData {
		if item.KodeIATA == kode {
			return &item
		}
	}
	return nil
}
