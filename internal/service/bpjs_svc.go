package service

type IuranBPJS struct {
	Kelas      string
	Segmen     string
	Nominal    int
	Keterangan string
}

var iuranBPJSData = []IuranBPJS{
	{
		Kelas:      "Kelas I",
		Segmen:     "Peserta Mandiri",
		Nominal:    150000,
		Keterangan: "Iuran per orang per bulan untuk peserta mandiri kelas I",
	},
	{
		Kelas:      "Kelas II",
		Segmen:     "Peserta Mandiri",
		Nominal:    100000,
		Keterangan: "Iuran per orang per bulan untuk peserta mandiri kelas II",
	},
	{
		Kelas:      "Kelas III",
		Segmen:     "Peserta Mandiri",
		Nominal:    35000,
		Keterangan: "Iuran per orang per bulan (subsidi pemerintah Rp 7.000, total iuran Rp 42.000)",
	},
	{
		Kelas:      "Kelas III",
		Segmen:     "PBI (Penerima Bantuan Iuran)",
		Nominal:    0,
		Keterangan: "Ditanggung penuh oleh pemerintah (APBN/APBD), iuran Rp 42.000/bulan",
	},
	{
		Kelas:      "Sesuai Hak",
		Segmen:     "PPU (Pekerja Penerima Upah)",
		Nominal:    0,
		Keterangan: "5% dari gaji (4% ditanggung pemberi kerja, 1% ditanggung pekerja), maks gaji acuan Rp 12.000.000",
	},
	{
		Kelas:      "Sesuai Hak",
		Segmen:     "PBPU/Bukan Pekerja",
		Nominal:    0,
		Keterangan: "Iuran mengikuti pilihan kelas (I/II/III) sesuai pendaftaran peserta",
	},
}

type BPJSService struct{}

func NewBPJSService() *BPJSService { return &BPJSService{} }

func (s *BPJSService) GetIuran() []IuranBPJS {
	return iuranBPJSData
}
