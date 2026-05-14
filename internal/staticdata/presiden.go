package staticdata

type PresidenInfo struct {
	Urutan        int
	Nama          string
	WakilPresiden string
	MulaiJabatan  int
	AkhirJabatan  int // 0 = masih menjabat
}

var DaftarPresiden = []PresidenInfo{
	{1, "Ir. Soekarno", "Drs. Mohammad Hatta", 1945, 1967},
	{2, "H.M. Soeharto", "Sri Sultan Hamengkubuwono IX, H. Adam Malik, Umar Wirahadikusumah, H. Sudharmono, H. Try Sutrisno, B.J. Habibie", 1967, 1998},
	{3, "Prof. Dr. B.J. Habibie", "", 1998, 1999},
	{4, "K.H. Abdurrahman Wahid", "Megawati Soekarnoputri", 1999, 2001},
	{5, "Megawati Soekarnoputri", "Hamzah Haz", 2001, 2004},
	{6, "H. Susilo Bambang Yudhoyono", "M. Jusuf Kalla, Boediono", 2004, 2014},
	{7, "Ir. H. Joko Widodo", "M. Jusuf Kalla, K.H. Ma'ruf Amin", 2014, 2024},
	{8, "H. Prabowo Subianto", "Gibran Rakabuming Raka", 2024, 0},
}
