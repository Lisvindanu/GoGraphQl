# IndonesiaQL — Developer Guide

Panduan untuk developer yang melanjutkan pengembangan project ini.

---

## Overview

Dua repo yang perlu diketahui:

| Repo | Stack | Lokasi Lokal | Live |
|------|-------|-------------|------|
| **API (GoGraphQL)** | Go 1.25 + gqlgen | `/Users/lisvindanuu/Sites/RisetGraphql` | `https://gographql.project-n.site/query` |
| **Web (GraphqlWEB)** | Laravel 13 + Inertia + React 19 | `/Users/lisvindanuu/Sites/GraphqlWEB` | *(frontend)* |

GitHub API repo: https://github.com/Lisvindanu/GoGraphQl

---

## 1. API — GoGraphQL

### Setup lokal

```bash
cp .env.example .env
# isi DB_URL dengan PostgreSQL connection string
go run ./cmd/server/
# server jalan di :8080, playground di http://localhost:8080/
```

### Struktur

```
cmd/server/main.go              ← entry point, HTTP mux, service wiring
graph/
  schema/*.graphqls             ← definisi GraphQL schema (schema-first)
  *.resolvers.go                ← implementasi resolver
  resolver.go                   ← struct Resolver + NewResolver()
internal/
  service/*_svc.go              ← business logic
  repository/*_repo.go          ← database queries (PostgreSQL)
  staticdata/*.go               ← data statis (Go slice, no DB)
  staticdata/photos/pahlawan/   ← 165 foto JPG, di-embed ke binary
  cache/inmemory.go             ← in-memory TTL cache
  middleware/                   ← CORS, rate limit, logging, recovery
  database/migrations/          ← SQL migrations
```

### Cara tambah fitur baru

1. Buat file schema di `graph/schema/nama_fitur.graphqls`
2. Jalankan `go run github.com/99designs/gqlgen@v0.17.90 generate`
3. Buat service di `internal/service/nama_svc.go`
4. Implementasi resolver yang di-generate di `graph/nama_fitur.resolvers.go`
5. Daftarkan service di `graph/resolver.go` (struct field + NewResolver param)
6. Wire di `cmd/server/main.go`

### Data sources

| Fitur | Sumber |
|-------|--------|
| Wilayah, NIK, Kode Pos | PostgreSQL (data BPS) |
| Cuaca | BMKG XML real-time |
| Kurs | Bank Indonesia SOAP XML |
| Hari Libur | PostgreSQL (SKB 3 Menteri) |
| Kalender Jawa, Hijriyah, Terbilang | Pure calculation |
| Kode Bank, Plat Nomor, BBM, BPJS, UMR, Bandara | Static slice di `internal/staticdata/` |
| Gempa | BMKG JSON real-time |
| IHSG | Fetch live |
| Gunung Berapi, Pahlawan | Static slice di `internal/staticdata/` |

### Deploy ke VPS

```bash
# Dari lokal: push ke GitHub
git push origin main

# Di VPS (IP: 167.253.158.192)
cd /var/www/gographql-src
git pull
export PATH=$PATH:/usr/local/go/bin
go build -o /var/www/gographql/gographql ./cmd/server/
pm2 restart gographql   # PM2 id: 46, name: gographql
```

Binary jalan di port **8010**, di-proxy Nginx ke domain.

---

## 2. Web — GraphqlWEB

### Setup lokal

```bash
composer install
npm install
cp .env.example .env
php artisan key:generate
npm run dev      # Vite dev server
php artisan serve
```

### Struktur

```
app/
  Http/Controllers/        ← satu controller per halaman
  Services/IndonesiaQLService.php  ← SEMUA query GraphQL ada di sini
resources/js/
  pages/                   ← React pages (Inertia)
  components/              ← shared components (Layout, NavBar, dll)
  types/indonesiaql.ts     ← TypeScript interfaces
routes/web.php             ← route definitions
```

### Pattern (PENTING)

**GraphQL TIDAK boleh di-fetch dari client-side React.** Alurnya:

```
Browser → Laravel Controller → IndonesiaQLService::method() → GraphQL API
                                                               ↓
Browser ← React page (Inertia props) ←── data PHP array ─────┘
```

### Cara tambah halaman baru

1. Tambah method di `app/Services/IndonesiaQLService.php` yang query GraphQL
2. Buat Controller di `app/Http/Controllers/NamaController.php`
3. Buat React page di `resources/js/pages/Nama/Index.tsx`
4. Tambah interface di `resources/js/types/indonesiaql.ts`
5. Tambah route di `routes/web.php`
6. Tambah link di `resources/js/components/NavBar.tsx`

### Halaman yang sudah ada

| URL | Deskripsi |
|-----|-----------|
| `/` | Landing page / welcome |
| `/playground` | GraphQL playground |
| `/wilayah` | Wilayah administratif (provinsi/kota) |
| `/cuaca` | Prakiraan cuaca BMKG |
| `/kurs` | Kurs Bank Indonesia |
| `/hari-libur` | Hari libur nasional |
| `/nik` | Validasi NIK |
| `/kalender-jawa` | Kalender Jawa |
| `/kalender-hijriyah` | Kalender Hijriyah |
| `/terbilang` | Angka ke teks |
| `/kode-bank` | Kode bank Indonesia |
| `/plat-nomor` | Plat nomor kendaraan |
| `/waktu-sholat` | Waktu sholat |
| `/gempa` | Info gempa BMKG |
| `/kode-pos` | Kode pos |
| `/harga-bbm` | Harga BBM Pertamina |
| `/saham` | IHSG |
| `/bpjs` | Iuran BPJS |
| `/rekening` | Validasi nomor rekening |
| `/inflasi` | Data inflasi |
| `/umr` | UMR provinsi |
| `/emas` | Harga emas Antam |
| `/bandara` | Bandara Indonesia |
| `/gunung-berapi` | 150 gunung berapi |
| `/pahlawan` | 191 pahlawan nasional (+ foto) |

### Deploy ke VPS

```bash
cd /var/www/LaravelGQL
git pull
composer install --no-dev
npm run build
php artisan config:cache
php artisan route:cache
php artisan view:cache
```

---

## GraphQL Endpoint

- **URL:** `https://gographql.project-n.site/query`
- **Method:** `POST`
- **Content-Type:** `application/json`

Contoh request:
```bash
curl -X POST https://gographql.project-n.site/query \
  -H "Content-Type: application/json" \
  -d '{"query":"{ pahlawanList { nama tahunDiangkat foto } }"}'
```

---

## Notes

- Foto pahlawan di-embed langsung ke dalam binary Go (tidak perlu file server terpisah)
- Cache: cuaca 30 menit, kurs 1 jam, gempa 5 menit, wilayah permanent
- Rate limit: dikonfigurasi via env `RATE_LIMIT_RPM`
- Untuk tambah data statis (kode bank baru, dll), edit file di `internal/staticdata/` lalu rebuild + deploy
