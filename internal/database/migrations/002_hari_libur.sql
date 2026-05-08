CREATE TABLE IF NOT EXISTS hari_libur (
    id SERIAL PRIMARY KEY,
    tanggal DATE NOT NULL,
    nama TEXT NOT NULL,
    tahun INT NOT NULL,
    jenis TEXT NOT NULL CHECK (jenis IN ('nasional', 'cuti_bersama'))
);

CREATE INDEX IF NOT EXISTS idx_hari_libur_tahun ON hari_libur(tahun);
CREATE INDEX IF NOT EXISTS idx_hari_libur_tanggal ON hari_libur(tanggal);
CREATE UNIQUE INDEX IF NOT EXISTS idx_hari_libur_tanggal_jenis ON hari_libur(tanggal, jenis);
