package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type StasiunKRL struct {
	StasiunId   string
	StasiunNama string
	StasiunKode string
}

type JadwalKRL struct {
	TrainId     string
	KaName      string
	RouteName   string
	DestTime    string
	DestStasiun string
	ColorCode   string
}

type KRLService struct {
	client *http.Client
}

func NewKRLService() *KRLService {
	return &KRLService{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *KRLService) fetch(apiURL string) ([]byte, error) {
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch KRL API: %w", err)
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

const comulineBaseURL = "https://www.api.comuline.com"

// stasiunList adalah daftar statis stasiun KRL Commuter Line Jabodetabek & Yogyakarta.
// ID disesuaikan dengan sistem Comuline untuk kompatibilitas jadwal.
var stasiunList = []StasiunKRL{
	// Lintas Jakarta Kota - Bogor/Nambo
	{StasiunId: "JAKK", StasiunNama: "Jakarta Kota", StasiunKode: "JAKK"},
	{StasiunId: "JAKJ", StasiunNama: "Jayakarta", StasiunKode: "JAKJ"},
	{StasiunId: "MGB", StasiunNama: "Mangga Besar", StasiunKode: "MGB"},
	{StasiunId: "SWB", StasiunNama: "Sawah Besar", StasiunKode: "SWB"},
	{StasiunId: "JUA", StasiunNama: "Juanda", StasiunKode: "JUA"},
	{StasiunId: "GMB", StasiunNama: "Gambir", StasiunKode: "GMB"},
	{StasiunId: "GDD", StasiunNama: "Gondangdia", StasiunKode: "GDD"},
	{StasiunId: "CKI", StasiunNama: "Cikini", StasiunKode: "CKI"},
	{StasiunId: "MRI", StasiunNama: "Manggarai", StasiunKode: "MRI"},
	{StasiunId: "TBT", StasiunNama: "Tebet", StasiunKode: "TBT"},
	{StasiunId: "CWG", StasiunNama: "Cawang", StasiunKode: "CWG"},
	{StasiunId: "DRK", StasiunNama: "Duren Kalibata", StasiunKode: "DRK"},
	{StasiunId: "PSB", StasiunNama: "Pasar Minggu Baru", StasiunKode: "PSB"},
	{StasiunId: "PSM", StasiunNama: "Pasar Minggu", StasiunKode: "PSM"},
	{StasiunId: "TJB", StasiunNama: "Tanjung Barat", StasiunKode: "TJB"},
	{StasiunId: "LNA", StasiunNama: "Lenteng Agung", StasiunKode: "LNA"},
	{StasiunId: "UPN", StasiunNama: "Universitas Pancasila", StasiunKode: "UPN"},
	{StasiunId: "UI", StasiunNama: "Universitas Indonesia", StasiunKode: "UI"},
	{StasiunId: "POC", StasiunNama: "Pondok Cina", StasiunKode: "POC"},
	{StasiunId: "DPB", StasiunNama: "Depok Baru", StasiunKode: "DPB"},
	{StasiunId: "DP", StasiunNama: "Depok", StasiunKode: "DP"},
	{StasiunId: "CTY", StasiunNama: "Citayam", StasiunKode: "CTY"},
	{StasiunId: "BJG", StasiunNama: "Bojong Gede", StasiunKode: "BJG"},
	{StasiunId: "CLB", StasiunNama: "Cilebut", StasiunKode: "CLB"},
	{StasiunId: "BO", StasiunNama: "Bogor", StasiunKode: "BO"},
	// Nambo
	{StasiunId: "NMB", StasiunNama: "Nambo", StasiunKode: "NMB"},
	{StasiunId: "CRB", StasiunNama: "Cibinong", StasiunKode: "CRB"},
	// Lintas Jatinegara - Cikarang
	{StasiunId: "JNG", StasiunNama: "Jatinegara", StasiunKode: "JNG"},
	{StasiunId: "KLD", StasiunNama: "Klender", StasiunKode: "KLD"},
	{StasiunId: "BUA", StasiunNama: "Buaran", StasiunKode: "BUA"},
	{StasiunId: "KLDB", StasiunNama: "Klender Baru", StasiunKode: "KLDB"},
	{StasiunId: "CKG", StasiunNama: "Cakung", StasiunKode: "CKG"},
	{StasiunId: "KRJ", StasiunNama: "Kranji", StasiunKode: "KRJ"},
	{StasiunId: "BKS", StasiunNama: "Bekasi", StasiunKode: "BKS"},
	{StasiunId: "BKST", StasiunNama: "Bekasi Timur", StasiunKode: "BKST"},
	{StasiunId: "TBN", StasiunNama: "Tambun", StasiunKode: "TBN"},
	{StasiunId: "CBT", StasiunNama: "Cibitung", StasiunKode: "CBT"},
	{StasiunId: "CKR", StasiunNama: "Cikarang", StasiunKode: "CKR"},
	// Lintas Tanah Abang - Rangkasbitung
	{StasiunId: "THB", StasiunNama: "Tanah Abang", StasiunKode: "THB"},
	{StasiunId: "PLM", StasiunNama: "Palmerah", StasiunKode: "PLM"},
	{StasiunId: "KBY", StasiunNama: "Kebayoran", StasiunKode: "KBY"},
	{StasiunId: "PDR", StasiunNama: "Pondok Ranji", StasiunKode: "PDR"},
	{StasiunId: "JRG", StasiunNama: "Jurangmangu", StasiunKode: "JRG"},
	{StasiunId: "SDM", StasiunNama: "Sudimara", StasiunKode: "SDM"},
	{StasiunId: "RBT", StasiunNama: "Rawa Buntu", StasiunKode: "RBT"},
	{StasiunId: "SRP", StasiunNama: "Serpong", StasiunKode: "SRP"},
	{StasiunId: "CSK", StasiunNama: "Cisauk", StasiunKode: "CSK"},
	{StasiunId: "CCY", StasiunNama: "Cicayur", StasiunKode: "CCY"},
	{StasiunId: "PPJ", StasiunNama: "Parung Panjang", StasiunKode: "PPJ"},
	{StasiunId: "LGK", StasiunNama: "Legok", StasiunKode: "LGK"},
	{StasiunId: "TNJ", StasiunNama: "Tenjo", StasiunKode: "TNJ"},
	{StasiunId: "TGR", StasiunNama: "Tigaraksa", StasiunKode: "TGR"},
	{StasiunId: "MJA", StasiunNama: "Maja", StasiunKode: "MJA"},
	{StasiunId: "CKY", StasiunNama: "Cikoya", StasiunKode: "CKY"},
	{StasiunId: "CTR", StasiunNama: "Citeras", StasiunKode: "CTR"},
	{StasiunId: "RKB", StasiunNama: "Rangkasbitung", StasiunKode: "RKB"},
	// Lintas Duri - Tangerang
	{StasiunId: "DRI", StasiunNama: "Duri", StasiunKode: "DRI"},
	{StasiunId: "GGL", StasiunNama: "Grogol", StasiunKode: "GGL"},
	{StasiunId: "PSG", StasiunNama: "Pesing", StasiunKode: "PSG"},
	{StasiunId: "TMK", StasiunNama: "Taman Kota", StasiunKode: "TMK"},
	{StasiunId: "BJI", StasiunNama: "Bojong Indah", StasiunKode: "BJI"},
	{StasiunId: "RWB", StasiunNama: "Rawa Buaya", StasiunKode: "RWB"},
	{StasiunId: "KLD2", StasiunNama: "Kalideres", StasiunKode: "KLD2"},
	{StasiunId: "PRS", StasiunNama: "Poris", StasiunKode: "PRS"},
	{StasiunId: "BTC", StasiunNama: "Batu Ceper", StasiunKode: "BTC"},
	{StasiunId: "TTG", StasiunNama: "Tanah Tinggi", StasiunKode: "TTG"},
	{StasiunId: "TNG", StasiunNama: "Tangerang", StasiunKode: "TNG"},
	// Lintas Jakarta Kota - Tanjung Priok
	{StasiunId: "KPB", StasiunNama: "Kampung Bandan", StasiunKode: "KPB"},
	{StasiunId: "ANC", StasiunNama: "Ancol", StasiunKode: "ANC"},
	{StasiunId: "TPK", StasiunNama: "Tanjung Priok", StasiunKode: "TPK"},
	// Transit
	{StasiunId: "SDR", StasiunNama: "Sudirman", StasiunKode: "SDR"},
	// Lintas Yogyakarta
	{StasiunId: "YK", StasiunNama: "Yogyakarta", StasiunKode: "YK"},
	{StasiunId: "LPN", StasiunNama: "Lempuyangan", StasiunKode: "LPN"},
	{StasiunId: "MRI2", StasiunNama: "Maguwo", StasiunKode: "MRI2"},
	{StasiunId: "BWK", StasiunNama: "Brambanan", StasiunKode: "BWK"},
	{StasiunId: "KT", StasiunNama: "Klaten", StasiunKode: "KT"},
}

func (s *KRLService) GetStasiun() ([]StasiunKRL, error) {
	return stasiunList, nil
}

var wib = time.FixedZone("WIB", 7*3600)

func parseDestTime(departsAt string) string {
	// Try RFC3339Nano first (handles fractional seconds like "2024-03-10T09:55:07.213Z")
	if t, err := time.Parse(time.RFC3339Nano, departsAt); err == nil {
		return t.In(wib).Format("15:04")
	}
	if t, err := time.Parse(time.RFC3339, departsAt); err == nil {
		return t.In(wib).Format("15:04")
	}
	// Fallback: extract HH:mm from "...THH:mm:ss..."
	if idx := strings.Index(departsAt, "T"); idx >= 0 && len(departsAt) >= idx+6 {
		return departsAt[idx+1 : idx+6]
	}
	return departsAt
}

func (s *KRLService) GetJadwal(stasiunId, timeFrom, timeTo string) ([]JadwalKRL, error) {
	apiURL := fmt.Sprintf("%s/v1/schedule/%s", comulineBaseURL, stasiunId)

	body, err := s.fetch(apiURL)
	if err != nil {
		return nil, err
	}

	var root struct {
		Metadata struct {
			Success bool `json:"success"`
		} `json:"metadata"`
		Data []struct {
			TrainID              string `json:"train_id"`
			Line                 string `json:"line"`
			Route                string `json:"route"`
			DepartsAt            string `json:"departs_at"`
			StationDestinationID string `json:"station_destination_id"`
			Metadata             struct {
				Origin struct {
					Color *string `json:"color"`
				} `json:"origin"`
			} `json:"metadata"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &root); err != nil {
		return nil, fmt.Errorf("parse jadwal: %w", err)
	}

	result := make([]JadwalKRL, 0, len(root.Data))
	for _, d := range root.Data {
		destTime := parseDestTime(d.DepartsAt)

		if timeFrom != "" && destTime < timeFrom {
			continue
		}
		if timeTo != "" && destTime > timeTo {
			continue
		}

		colorCode := ""
		if d.Metadata.Origin.Color != nil {
			colorCode = *d.Metadata.Origin.Color
		}

		result = append(result, JadwalKRL{
			TrainId:     d.TrainID,
			KaName:      d.Line,
			RouteName:   d.Route,
			DestTime:    destTime,
			DestStasiun: d.StationDestinationID,
			ColorCode:   colorCode,
		})
	}
	return result, nil
}
