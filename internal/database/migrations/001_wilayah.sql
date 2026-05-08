CREATE TABLE IF NOT EXISTS provinsi (
    kode CHAR(2) PRIMARY KEY,
    nama TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS kota (
    kode CHAR(4) PRIMARY KEY,
    provinsi_kode CHAR(2) NOT NULL REFERENCES provinsi(kode),
    nama TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS kecamatan (
    kode CHAR(7) PRIMARY KEY,
    kota_kode CHAR(4) NOT NULL REFERENCES kota(kode),
    nama TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS kelurahan (
    kode CHAR(10) PRIMARY KEY,
    kecamatan_kode CHAR(7) NOT NULL REFERENCES kecamatan(kode),
    nama TEXT NOT NULL,
    kode_pos CHAR(5)
);

CREATE INDEX IF NOT EXISTS idx_kota_provinsi ON kota(provinsi_kode);
CREATE INDEX IF NOT EXISTS idx_kecamatan_kota ON kecamatan(kota_kode);
CREATE INDEX IF NOT EXISTS idx_kelurahan_kecamatan ON kelurahan(kecamatan_kode);
CREATE INDEX IF NOT EXISTS idx_kelurahan_kode_pos ON kelurahan(kode_pos);
CREATE INDEX IF NOT EXISTS idx_provinsi_nama ON provinsi(nama);
CREATE INDEX IF NOT EXISTS idx_kota_nama ON kota(nama);
CREATE INDEX IF NOT EXISTS idx_kecamatan_nama ON kecamatan(nama);
CREATE INDEX IF NOT EXISTS idx_kelurahan_nama ON kelurahan(nama);
