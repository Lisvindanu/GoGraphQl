package staticdata

type NomorDaruratInfo struct {
	Nomor    string
	Layanan  string
	Kategori string
}

var NomorDarurat = []NomorDaruratInfo{
	{"112", "Layanan Darurat Terpadu", "Darurat Utama"},
	{"110", "Kepolisian RI", "Keamanan"},
	{"113", "Pemadam Kebakaran", "Keamanan"},
	{"115", "Basarnas / SAR", "Kebencanaan"},
	{"117", "BNPB (Bencana Nasional)", "Kebencanaan"},
	{"118", "Ambulans (Kemenkes)", "Medis"},
	{"119", "SPGDT / Ambulans DKI Jakarta", "Medis"},
	{"122", "Posko Kewaspadaan Nasional", "Keamanan"},
	{"123", "PLN (Gangguan Listrik)", "Utilitas"},
	{"129", "Posko Bencana Alam", "Kebencanaan"},
	{"151", "Kementerian Perhubungan", "Transportasi"},
	{"14080", "Jasa Marga (Jalan Tol)", "Transportasi"},
	{"1500-400", "BPJS Kesehatan", "Medis"},
	{"1500-533", "BPOM (Keracunan / Obat Pangan)", "Medis"},
	{"1500-454", "BPJS Ketenagakerjaan", "Sosial"},
	{"021-7992325", "PMI (Palang Merah Indonesia)", "Medis"},
	{"136", "Komnas HAM", "Sosial"},
	{"138", "Dirjen Pajak", "Pemerintahan"},
	{"139", "KAI (Kereta Api Indonesia)", "Transportasi"},
	{"196", "Jasa Marga Tol (alternatif)", "Transportasi"},
}
