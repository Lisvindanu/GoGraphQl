# IndonesiaQL

> Open GraphQL API for Indonesian public data and utilities.

IndonesiaQL menyediakan satu endpoint GraphQL untuk mengakses berbagai data publik Indonesia — mulai dari data wilayah administratif, prakiraan cuaca, kurs mata uang, hingga utilitas seperti validasi NIK dan konversi kalender Jawa.

## Features

### 🗺️ Wilayah Administratif
Data lengkap wilayah Indonesia dari BPS — provinsi, kota/kabupaten, kecamatan, hingga kelurahan beserta kode pos. Mendukung pencarian teks bebas.

```graphql
query {
  provinsiList {
    kode
    nama
    kota {
      kode
      nama
    }
  }
}

query {
  searchWilayah(query: "kebayoran", limit: 5) {
    kode
    nama
    tipe
    kota
    provinsi
  }
}
```

### 🌤️ Cuaca BMKG
Prakiraan cuaca resmi dari BMKG untuk seluruh kota di Indonesia — suhu, kelembapan, kondisi cuaca, dan arah/kecepatan angin.

```graphql
query {
  cuaca(provinsiKode: "31", kota: "jakarta pusat") {
    kota
    prakiraan {
      waktu
      suhu
      kelembapan
      cuaca
      kecepatanAngin
      arahAngin
    }
  }
}
```

### 💱 Kurs Bank Indonesia
Nilai tukar mata uang resmi dari Bank Indonesia, diperbarui setiap hari kerja.

```graphql
query {
  kurs(mataUang: "USD") {
    mataUang
    kursBeli
    kursJual
    kursTengah
    tanggal
  }
}

# Semua mata uang sekaligus
query {
  kurs {
    mataUang
    kursTengah
    tanggal
  }
}
```

### 📅 Hari Libur Nasional
Daftar hari libur nasional dan cuti bersama berdasarkan SKB 3 Menteri.

```graphql
query {
  hariLibur(tahun: 2025, bulan: 12) {
    tanggal
    nama
    jenis
  }
}

query {
  hariLiburHariIni {
    nama
    jenis
  }
}
```

### 🪪 Validasi NIK
Validasi format NIK dan ekstraksi informasi: provinsi, kota, tanggal lahir, dan jenis kelamin.

```graphql
query {
  validasiNIK(nik: "3201012501900001") {
    valid
    provinsi
    kota
    tanggalLahir
    jenisKelamin
    errors
  }
}
```

### 📆 Kalender Jawa
Konversi tanggal Masehi ke Kalender Jawa — hari pasaran, wuku, windu, dan tahun Jawa.

```graphql
query {
  kalenderJawa(tanggal: "2025-08-17") {
    hari
    pasaran
    wuku
    tahunJawa
    namaWindu
    tahunDalamWindu
  }
}
```

### 🔢 Terbilang
Konversi angka ke teks Bahasa Indonesia, mendukung desimal dan negatif hingga 999 triliun.

```graphql
query {
  terbilang(angka: 1500000) {
    angka
    terbilang
  }
}
# → "satu juta lima ratus ribu"
```

## Tech Stack

- **Runtime:** Go 1.25
- **GraphQL:** [gqlgen](https://github.com/99designs/gqlgen) v0.17 (schema-first)
- **Database:** PostgreSQL 16 (wilayah & hari libur)
- **Cache:** In-memory TTL cache (cuaca 30 menit, kurs 1 jam)
- **Data sources:** BPS, BMKG, Bank Indonesia

## Endpoint

| | |
|---|---|
| GraphQL | `POST /query` |
| Playground | `GET /` *(dev only)* |
| Health | `GET /health` |

## License

[MIT](./LICENSE)
