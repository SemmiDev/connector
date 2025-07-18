package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
	gl "lab.garudacyber.co.id/g-learning-connector"
)

type (
	GetTotalKelasResponse struct {
		Total int64 `json:"total"`
	}

	ListStudentKelasRequest struct {
		gl.Filter
		Semester string `json:"semester" form:"semester" query:"semester"`
	}

	KelasPerkuliahan struct {
		IDKelas         string `json:"id_kelas"`
		IDSMS           string `json:"id_sms"`
		NamaKelas       string `json:"nama_kelas"`
		NamaMatakuliah  string `json:"nama_matakuliah"`
		KodeMatakuliah  string `json:"kode_matakuliah"`
		IDDosenPengajar string `json:"id_dosen_pengajar"`
		Jadwal          string `json:"jadwal"`
	}

	ListStudentKelasResponse struct {
		IDPesertaDidik   string             `json:"id_pd"`
		IDMahasiswa      string             `json:"id_mahasiswa"`
		Nik              string             `json:"nik"`
		Semester         string             `json:"semester"`
		KelasPerkuliahan []KelasPerkuliahan `json:"kelas_perkuliahan"`
	}

	ListStudentKelasModel struct {
		IDPesertaDidik     string `json:"id_pd" gorm:"column:id_pd"`
		IDMahasiswa        string `json:"id_mahasiswa" gorm:"column:id_mahasiswa"`
		NIK                string `json:"nik" gorm:"column:nik"`
		IDKelas            string `json:"id_kelas" gorm:"column:id_kelas"`
		IDSMS              string `json:"id_sms" gorm:"column:id_sms"`
		NamaKelas          string `json:"nama_kelas" gorm:"column:nama_kelas"`
		NamaMataKuliah     string `json:"nama_matakuliah" gorm:"column:nama_matakuliah"`
		KodeMataKuliah     string `json:"kode_matakuliah" gorm:"column:kode_matakuliah"`
		IDPTKDosenPengajar string `json:"id_dosen_pengajar" gorm:"column:id_dosen_pengajar"`
		Semester           string `json:"semester" gorm:"column:semester"`
		Jadwal             string `json:"jadwal" gorm:"column:jadwal"`
	}

	Ruangan struct {
		IDRuangan      string    `json:"id_ruangan" gorm:"type:varchar(255)"`
		IDSMSRaw       string    `json:"id_sms" gorm:"column:id_sms"`
		NamaRuangan    string    `json:"nama_ruangan" gorm:"type:varchar(255)"`
		IDJenisRuangan string    `json:"id_jenis_ruangan" gorm:"type:varchar(255)"`
		KodeRuangan    string    `json:"kode_ruangan" gorm:"type:varchar(191)"`
		Keterangan     string    `json:"keterangan" gorm:"type:text"`
		Kapasitas      int       `json:"kapasitas" gorm:"type:int"`
		IDInstansi     string    `json:"id_instansi" gorm:"type:char(26)"`
		CreatedAt      time.Time `json:"created_at" gorm:"type:timestamp;default:null"`
		UpdatedAt      time.Time `json:"updated_at" gorm:"type:timestamp;default:null"`
	}

	RuanganResponse struct {
		IDRuangan      string    `json:"id_ruangan" gorm:"type:varchar(255)"`
		IDSMS          []string  `json:"id_sms"`
		NamaRuangan    string    `json:"nama_ruangan" gorm:"type:varchar(255)"`
		IDJenisRuangan string    `json:"id_jenis_ruangan" gorm:"type:varchar(255)"`
		KodeRuangan    string    `json:"kode_ruangan" gorm:"type:varchar(191)"`
		Keterangan     string    `json:"keterangan" gorm:"type:text"`
		Kapasitas      int       `json:"kapasitas" gorm:"type:int"`
		IDInstansi     string    `json:"id_instansi" gorm:"type:char(26)"`
		CreatedAt      time.Time `json:"created_at" gorm:"type:timestamp;default:null"`
		UpdatedAt      time.Time `json:"updated_at" gorm:"type:timestamp;default:null"`
	}

	SMS struct {
		IDSms               string     `json:"id_sms" gorm:"column:id_sms;primaryKey"`
		NmLemb              string     `json:"nm_lemb" gorm:"column:nm_lemb"`
		NmLembInggris       *string    `json:"nm_lemb_inggris" gorm:"column:nm_lemb_inggris"`
		IDJenjangDidik      *int64     `json:"id_jenj_didik" gorm:"column:id_jenj_didik"`
		IDJenisSms          *int64     `json:"id_jns_sms" gorm:"column:id_jns_sms"`
		IDIndukSms          *string    `json:"id_induk_sms" gorm:"column:id_induk_sms"`
		KodeSms             *string    `json:"kode_sms" gorm:"column:kode_sms"`
		UUID                *string    `json:"uuid" gorm:"column:uuid"`
		BukaKrs             *bool      `json:"buka_krs" gorm:"column:buka_krs"`
		BukaNilai           *bool      `json:"buka_nilai" gorm:"column:buka_nilai"`
		BukaKhs             *bool      `json:"buka_khs" gorm:"column:buka_khs"`
		BukaKuesioner       *bool      `json:"buka_kuesioner" gorm:"column:buka_kuesioner"`
		BukaTranskrip       *bool      `json:"buka_transkrip" gorm:"column:buka_transkrip"`
		BukaKartuUjian      *bool      `json:"buka_kartu_ujian" gorm:"column:buka_kartu_ujian"`
		MulaiIsiKrs         *time.Time `json:"mulai_isi_krs" gorm:"column:mulai_isi_krs"`
		AkhirIsiKrs         *time.Time `json:"akhir_isi_krs" gorm:"column:akhir_isi_krs"`
		MulaiIsiNilai       *time.Time `json:"mulai_isi_nilai" gorm:"column:mulai_isi_nilai"`
		AkhirIsiNilai       *time.Time `json:"akhir_isi_nilai" gorm:"column:akhir_isi_nilai"`
		BebanStudi          *string    `json:"beban_studi" gorm:"column:beban_studi"`
		PenjadwalanBlok     *string    `json:"penjadwalan_blok" gorm:"column:penjadwalan_blok"`
		Gelar               *string    `json:"gelar" gorm:"column:gelar"`
		GelarSingkatan      *string    `json:"gelar_singkatan" gorm:"column:gelar_singkatan"`
		JenjKualifikasiKKNI *string    `json:"jenj_kualifikasi_kkni" gorm:"column:jenj_kualifikasi_kkni"`
		PersyaratanMasuk    *string    `json:"persyaratan_penerimaan" gorm:"column:persyaratan_penerimaan"`
		LamaStudi           *string    `json:"lama_studi" gorm:"column:lama_studi"`
		JenjLanjutan        *string    `json:"jenj_lanjutan" gorm:"column:jenj_lanjutan"`
		StatusProfesi       *string    `json:"status_profesi" gorm:"column:status_profesi"`
		KodeNIM             *string    `json:"kode_nim" gorm:"column:kode_nim"`
		NamaJenjangDidik    *string    `json:"nama_jenjang_didik" gorm:"column:nama_jenjang_didik"`
		CreatedAt           *time.Time `json:"created_at" gorm:"column:created_at"`
		UpdatedAt           *time.Time `json:"updated_at" gorm:"column:updated_at"`
	}
)

var (
	ErrMismatchData = errors.New("Data mismatch: all pipe-separated fields must have the same number of elements")
)

func convertListKelasModels(models []ListStudentKelasModel) ([]ListStudentKelasResponse, error) {
	responses := make([]ListStudentKelasResponse, 0)

	for _, model := range models {
		// Get the length of idKelasList as single source of truth
		idKelasList := strings.Split(model.IDKelas, "|")
		maxLen := len(idKelasList)

		// Safely split other fields and ensure they have at least maxLen elements
		safeGetElement := func(list []string, idx int) string {
			if idx < len(list) {
				return list[idx]
			}
			return "" // Return empty string for missing elements
		}

		// Split all fields
		smsKelasList := strings.Split(model.IDSMS, "|")
		namaKelasList := strings.Split(model.NamaKelas, "|")
		namaMatakuliahList := strings.Split(model.NamaMataKuliah, "|")
		kodeMatakuliahList := strings.Split(model.KodeMataKuliah, "|")
		idDosenPengajarList := strings.Split(model.IDPTKDosenPengajar, "|")

		// Handle jadwal field which might be completely empty
		var jadwalList []string
		if model.Jadwal != "" {
			jadwalList = strings.Split(model.Jadwal, "|")
		} else {
			jadwalList = make([]string, maxLen) // Create empty list with correct length
		}

		// Build kelasPerkuliahan with safe access
		kelasPerkuliahan := make([]KelasPerkuliahan, maxLen)
		for i := 0; i < maxLen; i++ {
			kelasPerkuliahan[i] = KelasPerkuliahan{
				IDKelas:         idKelasList[i],
				IDSMS:           safeGetElement(smsKelasList, i),
				NamaKelas:       safeGetElement(namaKelasList, i),
				NamaMatakuliah:  safeGetElement(namaMatakuliahList, i),
				KodeMatakuliah:  safeGetElement(kodeMatakuliahList, i),
				IDDosenPengajar: safeGetElement(idDosenPengajarList, i),
				Jadwal:          safeGetElement(jadwalList, i),
			}
		}

		// Add to response
		responses = append(responses, ListStudentKelasResponse{
			IDPesertaDidik:   model.IDPesertaDidik,
			IDMahasiswa:      model.IDMahasiswa,
			Nik:              model.NIK,
			Semester:         model.Semester,
			KelasPerkuliahan: kelasPerkuliahan,
		})
	}

	return responses, nil
}

func NewListKelasRequest() *ListStudentKelasRequest {
	return &ListStudentKelasRequest{
		Filter: gl.NewFilterPagination(),
	}
}

type ListSimpleStudentKelas struct {
	IDPd        string `json:"id_pd"`
	IDMahasiswa string `json:"id_mahasiswa"`
	NIK         string `json:"nik"`
	IDKelas     string `json:"id_kelas"`
	Semester    string `json:"semester"`
}

func (a *ApplicationServer) ListSimpleStudentKelasMisca(c *fiber.Ctx) error {
	if IsSmartInstansi(c) {
		return a.ListSimpleStudentKelasSmart(c)
	}

	req := NewListKelasRequest()
	if err := c.QueryParser(req); err != nil {
		return HandleError(c, err)
	}

	var activeSemester string

	// Ambil semester aktif
	if err := a.db.
		Table("setting").
		Where("param = ?", "periode_berlaku").
		Select("value").
		Scan(&activeSemester).
		Error; err != nil {
		return HandleError(c, err)
	}

	// Set default semester jika kosong
	if req.Semester == "" {
		req.Semester = activeSemester
	}

	// Model untuk menampung hasil query
	listKelas := make([]ListSimpleStudentKelas, 0)

	offset := req.Filter.GetOffset()
	limit := req.Filter.GetLimit()

	// Query hanya mengambil kolom yang diperlukan
	q := a.db.Table("nilai").
		Select(`
			nilai.id_pd AS id_pd,
			mahasiswa.id AS id_mahasiswa,
			mahasiswa.nik AS nik,
			GROUP_CONCAT(kelaskuliah.id_kls ORDER BY kelaskuliah.id_kls SEPARATOR '|') AS id_kelas,
			nilai.smt_ambil AS semester
		`).
		Joins("JOIN mahasiswa_histori ON mahasiswa_histori.id_pd = nilai.id_pd").
		Joins("JOIN mahasiswa ON mahasiswa.id = mahasiswa_histori.id_mahasiswa").
		Joins("JOIN kelaskuliah ON kelaskuliah.id_kls = nilai.id_kls").
		Where("nilai.smt_ambil = ?", req.Semester).
		Group("nilai.id_pd, mahasiswa.nik, nilai.smt_ambil")

	// Filter berdasarkan keyword
	if req.Filter.HasKeyword() {
		q = q.Where("mahasiswa.nik LIKE ?", "%"+req.Filter.Keyword+"%")
	}

	// Sorting default atau sesuai request
	if req.Filter.HasSort() {
		q = q.Order(clause.OrderByColumn{
			Column: clause.Column{Name: req.Filter.SortBy},
			Desc:   req.Filter.IsDesc(),
		})
	} else {
		q = q.Order("nilai.id_pd ASC")
	}

	// Menghitung jumlah total data
	var totalData int64
	if err := q.Count(&totalData).Error; err != nil {
		return HandleError(c, err)
	}

	// Menambahkan limit dan offset
	q = q.Offset(int(offset)).Limit(int(limit))

	// Eksekusi query
	if err := q.Scan(&listKelas).Error; err != nil {
		return HandleError(c, err)
	}

	// Membuat informasi paginasi
	pageInfo, err := gl.NewPageInfo(req.Filter.CurrentPage, limit, offset, totalData)
	if err != nil {
		return HandleError(c, err)
	}

	// Mengembalikan hasil sebagai JSON
	return c.Status(fiber.StatusOK).JSON(ApiResponse[ListDataApiResponseWrapper[ListSimpleStudentKelas]]{
		Code:    fiber.StatusOK,
		Status:  http.StatusText(fiber.StatusOK),
		Success: true,
		Message: "Sukses mendapatkan data kelas sederhana",
		Data: ListDataApiResponseWrapper[ListSimpleStudentKelas]{
			List:     listKelas,
			PageInfo: pageInfo,
		},
	})
}

func (a *ApplicationServer) TotalListSimpleStudentKelasMisca(c *fiber.Ctx) error {
	if IsSmartInstansi(c) {
		return a.TotalListSimpleStudentKelasSmart(c)
	}

	var activeSemester string

	// Ambil semester aktif
	if err := a.db.
		Table("setting").
		Where("param = ?", "periode_berlaku").
		Select("value").
		Scan(&activeSemester).
		Error; err != nil {
		return HandleError(c, err)
	}

	semester := c.Query("semester", activeSemester)

	var total int64

	err := a.db.Table("nilai").
		Joins("JOIN mahasiswa_histori ON mahasiswa_histori.id_pd = nilai.id_pd").
		Joins("JOIN mahasiswa ON mahasiswa.id = mahasiswa_histori.id_mahasiswa").
		Joins("JOIN kelaskuliah ON kelaskuliah.id_kls = nilai.id_kls").
		Where("nilai.smt_ambil = ?", semester).
		Group("nilai.id_pd, mahasiswa.nik, nilai.smt_ambil").
		Count(&total).Error

	if err != nil {
		return HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(ApiResponse[GetTotalKelasResponse]{
		Code:    fiber.StatusOK,
		Status:  http.StatusText(fiber.StatusOK),
		Success: true,
		Message: "Sukses mendapatkan total kelas sederhana",
		Data: GetTotalKelasResponse{
			Total: total,
		},
	})
}

func (a *ApplicationServer) ListStudentKelasDetailsMisca(c *fiber.Ctx) error {
	req := NewListKelasRequest()
	if err := c.QueryParser(req); err != nil {
		return HandleError(c, err)
	}

	var activeSemester string

	if err := a.db.
		Table("setting").
		Where("param = ?", "periode_berlaku").
		Select("value").
		Scan(&activeSemester).
		Error; err != nil {
		return HandleError(c, err)
	}

	// set default semester
	if req.Semester == "" {
		req.Semester = activeSemester
	}

	// Model untuk menampung hasil query
	listKelas := make([]ListStudentKelasModel, 0)

	offset := req.Filter.GetOffset()
	limit := req.Filter.GetLimit()

	// Membuat query
	q := a.db.Table("nilai").
		Select(`
			nilai.id_pd AS id_pd,
			mahasiswa.id AS id_mahasiswa,
			mahasiswa.nik AS nik,
			GROUP_CONCAT(kelaskuliah.id_kls ORDER BY kelaskuliah.id_kls SEPARATOR '|') AS id_kelas,
			GROUP_CONCAT(kelaskuliah.id_sms ORDER BY kelaskuliah.id_sms SEPARATOR '|') AS id_sms,
			GROUP_CONCAT(kelaskuliah.nm_kls ORDER BY kelaskuliah.id_kls SEPARATOR '|') AS nama_kelas,
			GROUP_CONCAT(matakuliah.nm_mk ORDER BY kelaskuliah.id_kls SEPARATOR '|') AS nama_matakuliah,
			GROUP_CONCAT(matakuliah.kode_mk ORDER BY kelaskuliah.id_kls SEPARATOR '|') AS kode_matakuliah,
			GROUP_CONCAT(akt_mengajar_dosen.id_ptk ORDER BY akt_mengajar_dosen.id_ptk SEPARATOR '|') AS id_dosen_pengajar,
			GROUP_CONCAT(
				CONCAT(
					CASE jadwal.hari
						WHEN '0' THEN 'Minggu'
						WHEN '1' THEN 'Senin'
						WHEN '2' THEN 'Selasa'
						WHEN '3' THEN 'Rabu'
						WHEN '4' THEN 'Kamis'
						WHEN '5' THEN 'Jumat'
						WHEN '6' THEN 'Sabtu'
						ELSE 'Unknown'
					END, '-',
					jadwal.jam_mulai, '-',
					jadwal.jam_selesai
				)
				ORDER BY jadwal.id_jadwal SEPARATOR '|'
			) AS jadwal,
			nilai.smt_ambil AS semester
		`).
		Joins("JOIN mahasiswa_histori ON mahasiswa_histori.id_pd = nilai.id_pd").
		Joins("JOIN mahasiswa ON mahasiswa.id = mahasiswa_histori.id_mahasiswa").
		Joins("JOIN kelaskuliah ON kelaskuliah.id_kls = nilai.id_kls").
		Joins("JOIN matakuliah_kurikulum ON matakuliah_kurikulum.id_mk_kur = kelaskuliah.id_mk_kur").
		Joins("JOIN matakuliah ON matakuliah.id_mk = matakuliah_kurikulum.id_mk").
		Joins("LEFT JOIN akt_mengajar_dosen ON akt_mengajar_dosen.id_kls = kelaskuliah.id_kls").
		Joins("LEFT JOIN jadwal ON jadwal.id_kls = kelaskuliah.id_kls").
		Where("nilai.smt_ambil = ?", req.Semester).
		Group("nilai.id_pd, mahasiswa.nik, nilai.smt_ambil")

	// Menambahkan pencarian berdasarkan keyword
	if req.Filter.HasKeyword() {
		q = q.Where("mahasiswa.nik LIKE ?", "%"+req.Filter.Keyword+"%")
	}

	// Menambahkan sorting
	if req.Filter.HasSort() {
		q = q.Order(clause.OrderByColumn{
			Column: clause.Column{Name: req.Filter.SortBy},
			Desc:   req.Filter.IsDesc(),
		})
	} else {
		q = q.Order("nilai.id_pd ASC")
	}

	// Menghitung jumlah total data
	var totalData int64
	if err := q.Count(&totalData).Error; err != nil {
		return HandleError(c, err)
	}

	// Menambahkan limit dan offset
	q = q.Offset(int(offset)).Limit(int(limit))

	// Eksekusi query
	if err := q.Scan(&listKelas).Error; err != nil {
		return HandleError(c, err)
	}

	// Membuat informasi paginasi
	pageInfo, err := gl.NewPageInfo(req.Filter.CurrentPage, limit, offset, totalData)
	if err != nil {
		return HandleError(c, err)
	}

	listKelasResponse, err := convertListKelasModels(listKelas)
	if err != nil {
		return HandleError(c, err)
	}

	// Mengembalikan hasil sebagai JSON
	return c.Status(fiber.StatusOK).JSON(ApiResponse[ListDataApiResponseWrapper[ListStudentKelasResponse]]{
		Code:    fiber.StatusOK,
		Status:  http.StatusText(fiber.StatusOK),
		Success: true,
		Message: "Sukses mendapatkan data kelas",
		Data: ListDataApiResponseWrapper[ListStudentKelasResponse]{
			List:     listKelasResponse,
			PageInfo: pageInfo,
		},
	})
}

func (a *ApplicationServer) GetTotalKelasDetailsMisca(c *fiber.Ctx) error {
	var activeSemester string

	if err := a.db.
		Table("setting").
		Where("param = ?", "periode_berlaku").
		Select("value").
		Scan(&activeSemester).
		Error; err != nil {
		return HandleError(c, err)
	}

	semester := c.Query("semester", activeSemester)

	var total int64

	// Buat query GORM
	err := a.db.Table("nilai").
		Select("nilai.id_pd").
		Joins("JOIN mahasiswa_histori ON mahasiswa_histori.id_pd = nilai.id_pd").
		Joins("JOIN mahasiswa ON mahasiswa.id = mahasiswa_histori.id_mahasiswa").
		Joins("JOIN kelaskuliah ON kelaskuliah.id_kls = nilai.id_kls").
		Joins("JOIN matakuliah_kurikulum ON matakuliah_kurikulum.id_mk_kur = kelaskuliah.id_mk_kur").
		Joins("JOIN matakuliah ON matakuliah.id_mk = matakuliah_kurikulum.id_mk").
		Joins("LEFT JOIN akt_mengajar_dosen ON akt_mengajar_dosen.id_kls = kelaskuliah.id_kls").
		Where("nilai.smt_ambil = ?", semester).
		Group("nilai.id_pd, mahasiswa.nik, nilai.smt_ambil").
		Count(&total).Error // Hitung total

	if err != nil {
		return HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(ApiResponse[GetTotalKelasResponse]{
		Code:    fiber.StatusOK,
		Status:  http.StatusText(fiber.StatusOK),
		Success: true,
		Message: "Sukses mendapatkan total kelas",
		Data: GetTotalKelasResponse{
			Total: total,
		},
	})
}

// ListKelasResponse defines the structure for class list response
type ListKelasResponse struct {
	IDKelas            string   `json:"id_kelas" gorm:"column:id_kelas"`
	IDSMS              string   `json:"id_sms" gorm:"column:id_sms"`
	NamaKelas          string   `json:"nama_kelas" gorm:"column:nama_kelas"`
	NamaMataKuliah     string   `json:"nama_matakuliah" gorm:"column:nama_matakuliah"`
	KodeMataKuliah     string   `json:"kode_matakuliah" gorm:"column:kode_matakuliah"`
	IDDosenPengajar    []string `json:"id_dosen_pengajar"`
	IDDosenPengajarStr string   `gorm:"column:id_dosen_pengajar"`
	Semester           string   `json:"semester" gorm:"column:semester"`
	Jadwal             string   `json:"jadwal" gorm:"column:jadwal"`
	NamaRuangan        string   `json:"nama_ruangan" gorm:"column:nama_ruangan"`
	TotalPertemuan     string   `json:"total_pertemuan" gorm:"column:total_pertemuan"`
}

func (a *ApplicationServer) TotalKelasMisca(c *fiber.Ctx) error {
	if IsSmartInstansi(c) {
		return a.TotalKelasSmart(c)
	}

	var activeSemester string

	if err := a.db.
		Table("setting").
		Where("param = ?", "periode_berlaku").
		Select("value").
		Scan(&activeSemester).
		Error; err != nil {
		return HandleError(c, err)
	}

	semester := c.Query("semester", activeSemester)

	var total int64

	// Buat query GORM
	err := a.db.Table("kelaskuliah").Where("kelaskuliah.id_smt = ?", semester).Count(&total).Error
	if err != nil {
		return HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(ApiResponse[GetTotalKelasResponse]{
		Code:    fiber.StatusOK,
		Status:  http.StatusText(fiber.StatusOK),
		Success: true,
		Message: "Sukses mendapatkan total kelas",
		Data: GetTotalKelasResponse{
			Total: total,
		},
	})
}

// ListKelas handles the listing of classes with multiple lecturers
func (a *ApplicationServer) ListKelasMisca(c *fiber.Ctx) error {
	if IsSmartInstansi(c) {
		return a.ListKelasSmart(c)
	}

	req := NewListKelasRequest()
	if err := c.QueryParser(req); err != nil {
		return HandleError(c, err)
	}

	var activeSemester string

	if err := a.db.
		Table("setting").
		Where("param = ?", "periode_berlaku").
		Select("value").
		Scan(&activeSemester).
		Error; err != nil {
		return HandleError(c, err)
	}

	// Set default semester
	if req.Semester == "" {
		req.Semester = activeSemester
	}

	// Model to hold query results
	listKelas := make([]ListKelasResponse, 0)

	offset := req.Filter.GetOffset()
	limit := req.Filter.GetLimit()

	q := a.db.
		Table("kelaskuliah").
		Select(`
			kelaskuliah.id_kls AS id_kelas,
			kelaskuliah.id_sms AS id_sms,
			kelaskuliah.nm_kls AS nama_kelas,
			matakuliah.nm_mk AS nama_matakuliah,
			matakuliah.kode_mk AS kode_matakuliah,
			GROUP_CONCAT(
				DISTINCT CAST(akt_mengajar_dosen.id_ptk AS CHAR)
				ORDER BY akt_mengajar_dosen.id_ptk ASC
				SEPARATOR '|'
			) AS id_dosen_pengajar,
			kelaskuliah.id_smt AS semester,
			GROUP_CONCAT(
				CONCAT(
					CASE jadwal.hari
						WHEN '0' THEN 'Minggu'
						WHEN '1' THEN 'Senin'
						WHEN '2' THEN 'Selasa'
						WHEN '3' THEN 'Rabu'
						WHEN '4' THEN 'Kamis'
						WHEN '5' THEN 'Jumat'
						WHEN '6' THEN 'Sabtu'
						ELSE 'Unknown'
					END,
					'-', jadwal.jam_mulai,
					'-', jadwal.jam_selesai
				)
				ORDER BY jadwal.hari, jadwal.jam_mulai ASC
				SEPARATOR '|'
			) AS jadwal,
			GROUP_CONCAT(
				ruangan.nama_ruangan
				ORDER BY ruangan.nama_ruangan ASC
				SEPARATOR '|'
			) AS nama_ruangan,
			GROUP_CONCAT(
				akt_mengajar_dosen.temu_rencana
				ORDER BY akt_mengajar_dosen.temu_rencana ASC
				SEPARATOR '|'
			) AS total_pertemuan
		`).
		Joins("JOIN matakuliah_kurikulum ON matakuliah_kurikulum.id_mk_kur = kelaskuliah.id_mk_kur").
		Joins("JOIN matakuliah ON matakuliah.id_mk = matakuliah_kurikulum.id_mk").
		Joins("LEFT JOIN akt_mengajar_dosen ON akt_mengajar_dosen.id_kls = kelaskuliah.id_kls").
		Joins("LEFT JOIN jadwal ON jadwal.id_kls = kelaskuliah.id_kls").
		Joins("LEFT JOIN ruangan ON ruangan.id_ruangan = jadwal.id_ruangan").
		Where("kelaskuliah.id_smt = ?", req.Semester).
		Group("kelaskuliah.id_kls")

	// Add keyword search
	if req.Filter.HasKeyword() {
		q = q.Where("kelaskuliah.nm_kls LIKE ?", "%"+req.Filter.Keyword+"%")
	}

	// Add sorting
	if req.Filter.HasSort() {
		q = q.Order(clause.OrderByColumn{
			Column: clause.Column{Name: req.Filter.SortBy},
			Desc:   req.Filter.IsDesc(),
		})
	} else {
		q = q.Order("kelaskuliah.id_kls ASC")
	}

	// Count total data
	var totalData int64
	if err := q.Count(&totalData).Error; err != nil {
		return HandleError(c, err)
	}

	// Add limit and offset
	q = q.Offset(int(offset)).Limit(int(limit))

	// Execute query
	if err := q.Scan(&listKelas).Error; err != nil {
		return HandleError(c, err)
	}

	// Post-process id_dosen_pengajar to convert pipe-separated string to slice
	for i := range listKelas {
		if listKelas[i].IDDosenPengajarStr != "" {
			listKelas[i].IDDosenPengajar = strings.Split(listKelas[i].IDDosenPengajarStr, "|")
		} else {
			listKelas[i].IDDosenPengajar = []string{}
		}
	}

	// Create pagination info
	pageInfo, err := gl.NewPageInfo(req.Filter.CurrentPage, limit, offset, totalData)
	if err != nil {
		return HandleError(c, err)
	}

	// Return result as JSON
	return c.Status(fiber.StatusOK).JSON(ApiResponse[ListDataApiResponseWrapper[ListKelasResponse]]{
		Code:    fiber.StatusOK,
		Status:  http.StatusText(fiber.StatusOK),
		Success: true,
		Message: "Sukses mendapatkan data kelas",
		Data: ListDataApiResponseWrapper[ListKelasResponse]{
			List:     listKelas,
			PageInfo: pageInfo,
		},
	})
}

func (a *ApplicationServer) TotalKelasSmart(c *fiber.Ctx) error {
	var IDSemesterAktif string
	err := a.db.Table("semester").Select(`id_smt`).Where("a_periode_aktif = 1").Scan(&IDSemesterAktif).Error
	if err != nil {
		return HandleError(c, err)
	}

	semester := c.Query("semester", IDSemesterAktif)
	var total int64

	err = a.db.Table("kelas_kuliah").Where("kelas_kuliah.id_smt = ?", semester).Count(&total).Error
	if err != nil {
		return HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(ApiResponse[GetTotalKelasResponse]{
		Code:    fiber.StatusOK,
		Status:  http.StatusText(fiber.StatusOK),
		Success: true,
		Message: "Sukses mendapatkan total kelas",
		Data: GetTotalKelasResponse{
			Total: total,
		},
	})
}

// ListKelas handles the listing of classes with multiple lecturers
func (a *ApplicationServer) ListKelasSmart(c *fiber.Ctx) error {
	req := NewListKelasRequest()
	if err := c.QueryParser(req); err != nil {
		return HandleError(c, err)
	}

	var activeSemester string
	err := a.db.Table("semester").Select(`id_smt`).Where("a_periode_aktif = 1").Scan(&activeSemester).Error
	if err != nil {
		return HandleError(c, err)
	}

	// Set default semester
	if req.Semester == "" {
		req.Semester = activeSemester
	}

	// Model to hold query results
	listKelas := make([]ListKelasResponse, 0)

	offset := req.Filter.GetOffset()
	limit := req.Filter.GetLimit()

	// Query utama
	q := a.db.
		Table("kelas_kuliah").
		Select(`
			kelas_kuliah.id_kls AS id_kelas,
			kelas_kuliah.id_sms AS id_sms,
			CONCAT_WS(
			' ',
			program.nm_program,
			kelas_kuliah.nm_kls,
			CONCAT('(Pilihan ', kelas_kuliah.pilihan_kelas, ')')
			) AS nama_kelas,
			matkul.nm_mk AS nama_matakuliah,
			matkul.kode_mk AS kode_matakuliah,
			GROUP_CONCAT(
				DISTINCT CAST(akt_ajar_dosen.id_reg_ptk AS CHAR)
				ORDER BY akt_ajar_dosen.id_reg_ptk ASC
				SEPARATOR '|'
			) AS id_dosen_pengajar,
			kelas_kuliah.id_smt AS semester,
			GROUP_CONCAT(
				CONCAT(
					CASE jadwal.hari
						WHEN '0' THEN 'Minggu'
						WHEN '1' THEN 'Senin'
						WHEN '2' THEN 'Selasa'
						WHEN '3' THEN 'Rabu'
						WHEN '4' THEN 'Kamis'
						WHEN '5' THEN 'Jumat'
						WHEN '6' THEN 'Sabtu'
						ELSE 'Unknown'
					END,
					'-', jadwal.jam_mulai,
					'-', jadwal.jam_selesai
				)
				ORDER BY jadwal.hari, jadwal.jam_mulai ASC
				SEPARATOR '|'
			) AS jadwal,
			GROUP_CONCAT(
				ruangan.kode_ruangan
				ORDER BY ruangan.kode_ruangan ASC
				SEPARATOR '|'
			) AS nama_ruangan,
			GROUP_CONCAT(
				akt_ajar_dosen.jml_tm_renc
				ORDER BY akt_ajar_dosen.jml_tm_renc ASC
				SEPARATOR '|'
			) AS total_pertemuan
		`).
		Joins("JOIN matkul ON matkul.id_mk = kelas_kuliah.id_mk").
		Joins("JOIN program ON program.id_program = kelas_kuliah.id_program").
		Joins("LEFT JOIN akt_ajar_dosen ON akt_ajar_dosen.id_kls = kelas_kuliah.id_kls").
		Joins("LEFT JOIN jadwal ON jadwal.id_kls = kelas_kuliah.id_kls").
		Joins("LEFT JOIN ruangan ON ruangan.id_ruangan = jadwal.id_ruangan").
		Where("kelas_kuliah.id_smt = ?", req.Semester).
		Debug().
		Group("kelas_kuliah.id_kls")

	// Add keyword search
	if req.Filter.HasKeyword() {
		q = q.Where("kelas_kuliah.nm_kls LIKE ?", "%"+req.Filter.Keyword+"%")
	}

	// Add sorting
	if req.Filter.HasSort() {
		q = q.Order(clause.OrderByColumn{
			Column: clause.Column{Name: req.Filter.SortBy},
			Desc:   req.Filter.IsDesc(),
		})
	} else {
		q = q.Order("kelas_kuliah.id_kls ASC")
	}

	// Count total data
	var totalData int64
	if err := q.Count(&totalData).Error; err != nil {
		return HandleError(c, err)
	}

	// Add limit and offset
	q = q.Offset(int(offset)).Limit(int(limit))

	// Execute query
	if err := q.Scan(&listKelas).Error; err != nil {
		return HandleError(c, err)
	}

	// Post-process id_dosen_pengajar to convert pipe-separated string to slice
	for i := range listKelas {
		if listKelas[i].IDDosenPengajarStr != "" {
			listKelas[i].IDDosenPengajar = strings.Split(listKelas[i].IDDosenPengajarStr, "|")
		} else {
			listKelas[i].IDDosenPengajar = []string{}
		}
	}

	// Create pagination info
	pageInfo, err := gl.NewPageInfo(req.Filter.CurrentPage, limit, offset, totalData)
	if err != nil {
		return HandleError(c, err)
	}

	// Return result as JSON
	return c.Status(fiber.StatusOK).JSON(ApiResponse[ListDataApiResponseWrapper[ListKelasResponse]]{
		Code:    fiber.StatusOK,
		Status:  http.StatusText(fiber.StatusOK),
		Success: true,
		Message: "Sukses mendapatkan data kelas",
		Data: ListDataApiResponseWrapper[ListKelasResponse]{
			List:     listKelas,
			PageInfo: pageInfo,
		},
	})
}

// smart
func (a *ApplicationServer) ListSimpleStudentKelasSmart(c *fiber.Ctx) error {
	req := NewListKelasRequest()
	if err := c.QueryParser(req); err != nil {
		return HandleError(c, err)
	}

	var activeSemester string
	err := a.db.Table("semester").Select(`id_smt`).Where("a_periode_aktif = 1").Scan(&activeSemester).Error
	if err != nil {
		return HandleError(c, err)
	}

	// Set default semester
	if req.Semester == "" {
		req.Semester = activeSemester
	}

	// Model untuk menampung hasil query
	listKelas := make([]ListSimpleStudentKelas, 0)

	offset := req.Filter.GetOffset()
	limit := req.Filter.GetLimit()

	// Query hanya mengambil kolom yang diperlukan
	q := a.db.Table("nilai").
		Select(`
			nilai.id_reg_pd AS id_pd,
			nilai.id_reg_pd AS id_mahasiswa,
			mahasiswa.nik AS nik,
			GROUP_CONCAT(kelas_kuliah.id_kls ORDER BY kelas_kuliah.id_kls SEPARATOR '|') AS id_kelas,
			kelas_kuliah.id_smt AS semester
		`).
		Joins("JOIN mahasiswa ON mahasiswa.id_pd = nilai.id_reg_pd").
		Joins("JOIN kelas_kuliah ON kelas_kuliah.id_kls = nilai.id_kls").
		Where("kelas_kuliah.id_smt = ?", req.Semester).
		Group("nilai.id_reg_pd, mahasiswa.nik, kelas_kuliah.id_smt")

	// Filter berdasarkan keyword
	if req.Filter.HasKeyword() {
		q = q.Where("mahasiswa.nik LIKE ?", "%"+req.Filter.Keyword+"%")
	}

	// Sorting default atau sesuai request
	if req.Filter.HasSort() {
		q = q.Order(clause.OrderByColumn{
			Column: clause.Column{Name: req.Filter.SortBy},
			Desc:   req.Filter.IsDesc(),
		})
	} else {
		q = q.Order("nilai.id_reg_pd ASC")
	}

	// Menghitung jumlah total data
	var totalData int64
	if err := q.Count(&totalData).Error; err != nil {
		return HandleError(c, err)
	}

	// Menambahkan limit dan offset
	q = q.Offset(int(offset)).Limit(int(limit))

	// Eksekusi query
	if err := q.Scan(&listKelas).Error; err != nil {
		return HandleError(c, err)
	}

	// Membuat informasi paginasi
	pageInfo, err := gl.NewPageInfo(req.Filter.CurrentPage, limit, offset, totalData)
	if err != nil {
		return HandleError(c, err)
	}

	// Mengembalikan hasil sebagai JSON
	return c.Status(fiber.StatusOK).JSON(ApiResponse[ListDataApiResponseWrapper[ListSimpleStudentKelas]]{
		Code:    fiber.StatusOK,
		Status:  http.StatusText(fiber.StatusOK),
		Success: true,
		Message: "Sukses mendapatkan data kelas sederhana",
		Data: ListDataApiResponseWrapper[ListSimpleStudentKelas]{
			List:     listKelas,
			PageInfo: pageInfo,
		},
	})
}

func (a *ApplicationServer) TotalListSimpleStudentKelasSmart(c *fiber.Ctx) error {
	var activeSemester string
	err := a.db.Table("semester").Select(`id_smt`).Where("a_periode_aktif = 1").Scan(&activeSemester).Error
	if err != nil {
		return HandleError(c, err)
	}

	semester := c.Query("semester", activeSemester)

	var total int64

	err = a.db.Table("nilai").
		Select(`
			nilai.id_reg_pd AS id_pd,
			nilai.id_reg_pd AS id_mahasiswa,
			mahasiswa.nik AS nik,
			GROUP_CONCAT(kelas_kuliah.id_kls ORDER BY kelas_kuliah.id_kls SEPARATOR '|') AS id_kelas,
			kelas_kuliah.id_smt AS semester
		`).
		Joins("JOIN mahasiswa ON mahasiswa.id_pd = nilai.id_reg_pd").
		Joins("JOIN kelas_kuliah ON kelas_kuliah.id_kls = nilai.id_kls").
		Where("kelas_kuliah.id_smt = ?", semester).
		Group("nilai.id_reg_pd, mahasiswa.nik, kelas_kuliah.id_smt").
		Count(&total).Error

	if err != nil {
		return HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(ApiResponse[GetTotalKelasResponse]{
		Code:    fiber.StatusOK,
		Status:  http.StatusText(fiber.StatusOK),
		Success: true,
		Message: "Sukses mendapatkan total kelas sederhana",
		Data: GetTotalKelasResponse{
			Total: total,
		},
	})
}
