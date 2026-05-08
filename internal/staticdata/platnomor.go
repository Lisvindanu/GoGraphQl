package staticdata

type PlatInfo struct { Wilayah string; Provinsi string }

// PlatNomor maps plat prefix to wilayah info
var PlatNomor = map[string]PlatInfo{
	// Jawa
	"A":  {Wilayah: "Banten", Provinsi: "Banten"},
	"B":  {Wilayah: "Jakarta", Provinsi: "DKI Jakarta"},
	"D":  {Wilayah: "Bandung", Provinsi: "Jawa Barat"},
	"E":  {Wilayah: "Cirebon", Provinsi: "Jawa Barat"},
	"F":  {Wilayah: "Bogor", Provinsi: "Jawa Barat"},
	"G":  {Wilayah: "Pekalongan", Provinsi: "Jawa Tengah"},
	"H":  {Wilayah: "Semarang", Provinsi: "Jawa Tengah"},
	"K":  {Wilayah: "Pati", Provinsi: "Jawa Tengah"},
	"L":  {Wilayah: "Surabaya", Provinsi: "Jawa Timur"},
	"M":  {Wilayah: "Madura", Provinsi: "Jawa Timur"},
	"N":  {Wilayah: "Malang", Provinsi: "Jawa Timur"},
	"P":  {Wilayah: "Besuki / Banyuwangi", Provinsi: "Jawa Timur"},
	"R":  {Wilayah: "Banyumas", Provinsi: "Jawa Tengah"},
	"S":  {Wilayah: "Bojonegoro", Provinsi: "Jawa Timur"},
	"T":  {Wilayah: "Karawang", Provinsi: "Jawa Barat"},
	"W":  {Wilayah: "Mojokerto / Sidoarjo / Gresik", Provinsi: "Jawa Timur"},
	"Z":  {Wilayah: "Tasikmalaya / Garut", Provinsi: "Jawa Barat"},
	"AA": {Wilayah: "Kedu / Magelang", Provinsi: "Jawa Tengah"},
	"AB": {Wilayah: "Yogyakarta", Provinsi: "DI Yogyakarta"},
	"AD": {Wilayah: "Surakarta / Solo", Provinsi: "Jawa Tengah"},
	"AE": {Wilayah: "Madiun", Provinsi: "Jawa Timur"},
	"AG": {Wilayah: "Kediri", Provinsi: "Jawa Timur"},
	// Bali & Nusa Tenggara
	"DK": {Wilayah: "Bali", Provinsi: "Bali"},
	"DR": {Wilayah: "Mataram / NTB", Provinsi: "Nusa Tenggara Barat"},
	"DH": {Wilayah: "Kupang / NTT", Provinsi: "Nusa Tenggara Timur"},
	// Sumatra
	"BA": {Wilayah: "Padang", Provinsi: "Sumatera Barat"},
	"BB": {Wilayah: "Tapanuli / Sumut Bagian Selatan", Provinsi: "Sumatera Utara"},
	"BD": {Wilayah: "Bengkulu", Provinsi: "Bengkulu"},
	"BE": {Wilayah: "Bandar Lampung", Provinsi: "Lampung"},
	"BG": {Wilayah: "Palembang", Provinsi: "Sumatera Selatan"},
	"BH": {Wilayah: "Jambi", Provinsi: "Jambi"},
	"BK": {Wilayah: "Medan", Provinsi: "Sumatera Utara"},
	"BL": {Wilayah: "Banda Aceh", Provinsi: "Aceh"},
	"BM": {Wilayah: "Pekanbaru", Provinsi: "Riau"},
	"BN": {Wilayah: "Pangkal Pinang", Provinsi: "Kepulauan Bangka Belitung"},
	"BP": {Wilayah: "Batam / Tanjung Pinang", Provinsi: "Kepulauan Riau"},
	// Kalimantan
	"DA": {Wilayah: "Banjarmasin", Provinsi: "Kalimantan Selatan"},
	"KB": {Wilayah: "Pontianak", Provinsi: "Kalimantan Barat"},
	"KH": {Wilayah: "Palangka Raya", Provinsi: "Kalimantan Tengah"},
	"KT": {Wilayah: "Samarinda / Balikpapan", Provinsi: "Kalimantan Timur"},
	"KU": {Wilayah: "Tanjung Selor", Provinsi: "Kalimantan Utara"},
	// Sulawesi
	"DL": {Wilayah: "Manado", Provinsi: "Sulawesi Utara"},
	"DM": {Wilayah: "Gorontalo", Provinsi: "Gorontalo"},
	"DN": {Wilayah: "Palu", Provinsi: "Sulawesi Tengah"},
	"DT": {Wilayah: "Kendari", Provinsi: "Sulawesi Tenggara"},
	"DW": {Wilayah: "Makassar", Provinsi: "Sulawesi Selatan"},
	"DX": {Wilayah: "Mamuju", Provinsi: "Sulawesi Barat"},
	// Maluku & Papua
	"DE": {Wilayah: "Ambon", Provinsi: "Maluku"},
	"DG": {Wilayah: "Ternate", Provinsi: "Maluku Utara"},
	"PA": {Wilayah: "Jayapura", Provinsi: "Papua"},
	"PB": {Wilayah: "Manokwari", Provinsi: "Papua Barat"},
	"PK": {Wilayah: "Sorong", Provinsi: "Papua Barat Daya"},
	"PP": {Wilayah: "Nabire / Papua Tengah", Provinsi: "Papua Tengah"},
	"PD": {Wilayah: "Wamena / Papua Pegunungan", Provinsi: "Papua Pegunungan"},
}
