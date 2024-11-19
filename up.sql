CREATE TABLE IF NOT EXISTS api_tokens (
    id bigint(20) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name varchar(255) NOT NULL, -- glearning
    token text UNIQUE NOT NULL, -- some_secret_token
    active tinyint(1) NOT NULL DEFAULT '1'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO api_tokens (name, token, active) VALUES ('glearning', '6qPRWChjOmfziYo0dASFKS+vnkZGxHgg', 1);


SELECT
    nilai.id_pd AS id_pd,
    mahasiswa.nik AS nik,
    GROUP_CONCAT(kelaskuliah.id_kls ORDER BY kelaskuliah.id_kls SEPARATOR '|') AS id_kelas,
    GROUP_CONCAT(kelaskuliah.nm_kls ORDER BY kelaskuliah.id_kls SEPARATOR '|') AS nama_kelas,
    GROUP_CONCAT(matakuliah.nm_mk ORDER BY kelaskuliah.id_kls SEPARATOR '|') AS nama_matakuliah,
    GROUP_CONCAT(matakuliah.kode_mk ORDER BY kelaskuliah.id_kls SEPARATOR '|') AS kode_matakuliah,
    GROUP_CONCAT(akt_ajar_dosen.id_ptk ORDER BY akt_ajar_dosen.id_ptk SEPARATOR '|') AS id_dosen_pengajar,
    nilai.smt_ambil AS semester
FROM
    nilai
        JOIN mahasiswa_histori ON mahasiswa_histori.id_pd = nilai.id_pd
        JOIN mahasiswa ON mahasiswa.id = mahasiswa_histori.id_mahasiswa
        JOIN kelaskuliah ON kelaskuliah.id_kls = nilai.id_kls
        JOIN matakuliah_kurikulum ON matakuliah_kurikulum.id_mk_kur = kelaskuliah.id_mk_kur
        JOIN matakuliah ON matakuliah.id_mk = matakuliah_kurikulum.id_mk
        LEFT JOIN akt_ajar_dosen ON akt_ajar_dosen.id_kls = kelaskuliah.id_kls
WHERE
        nilai.smt_ambil = '20241'
GROUP BY
    nilai.id_pd, mahasiswa.nik, nilai.smt_ambil
ORDER BY
    nilai.id_pd;






-- make sure jumlah dosen 1

SELECT
    id_kls,
    COUNT(id_ptk) AS jumlah_dosen
FROM
    akt_ajar_dosen
GROUP BY
    id_kls
HAVING
    COUNT(id_ptk) = 1;


