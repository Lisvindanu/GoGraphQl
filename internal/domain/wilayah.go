package domain

type Provinsi struct {
	Kode string
	Nama string
}

type Kota struct {
	Kode         string
	ProvinsiKode string
	Nama         string
}

type Kecamatan struct {
	Kode     string
	KotaKode string
	Nama     string
}

type Kelurahan struct {
	Kode          string
	KecamatanKode string
	Nama          string
	KodePos       *string
}

type WilayahSearchResult struct {
	Kode     string
	Nama     string
	Tipe     string
	Kota     *string
	Provinsi *string
}
