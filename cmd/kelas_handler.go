package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
	gl "lab.garudacyber.co.id/g-learning-connector"
	"net/http"
	"strings"
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
		NamaKelas       string `json:"nama_kelas"`
		NamaMatakuliah  string `json:"nama_matakuliah"`
		KodeMatakuliah  string `json:"kode_matakuliah"`
		IDDosenPengajar string `json:"id_dosen_pengajar"`
		Jadwal          string `json:"jadwal"`
	}

	ListStudentKelasResponse struct {
		IDPesertaDidik   string             `json:"id_pd"`
		Nik              string             `json:"nik"`
		Semester         string             `json:"semester"`
		KelasPerkuliahan []KelasPerkuliahan `json:"kelas_perkuliahan"`
	}

	ListStudentKelasModel struct {
		IDPesertaDidik     string `json:"id_pd" gorm:"column:id_pd"`
		NIK                string `json:"nik" gorm:"column:nik"`
		IDKelas            string `json:"id_kelas" gorm:"column:id_kelas"`
		NamaKelas          string `json:"nama_kelas" gorm:"column:nama_kelas"`
		NamaMataKuliah     string `json:"nama_matakuliah" gorm:"column:nama_matakuliah"`
		KodeMataKuliah     string `json:"kode_matakuliah" gorm:"column:kode_matakuliah"`
		IDPTKDosenPengajar string `json:"id_dosen_pengajar" gorm:"column:id_dosen_pengajar"`
		Semester           string `json:"semester" gorm:"column:semester"`
		Jadwal             string `json:"jadwal" gorm:"column:jadwal"`
	}
)

var (
	ErrMismatchData = errors.New("Data mismatch: all pipe-separated fields must have the same number of elements")
)

func convertListKelasModels(models []ListStudentKelasModel) ([]ListStudentKelasResponse, error) {
	responses := make([]ListStudentKelasResponse, 0)

	for _, model := range models {
		// Split setiap kolom berdasarkan '|'
		idKelasList := strings.Split(model.IDKelas, "|")
		namaKelasList := strings.Split(model.NamaKelas, "|")
		namaMatakuliahList := strings.Split(model.NamaMataKuliah, "|")
		kodeMatakuliahList := strings.Split(model.KodeMataKuliah, "|")
		idDosenPengajarList := strings.Split(model.IDPTKDosenPengajar, "|")
		jadwalList := strings.Split(model.Jadwal, "|")

		// Pastikan semua array memiliki panjang yang sama
		maxLen := len(idKelasList)
		if len(namaKelasList) != maxLen || len(namaMatakuliahList) != maxLen || len(kodeMatakuliahList) != maxLen || len(idDosenPengajarList) != maxLen || len(jadwalList) != maxLen {
			return nil, ErrMismatchData
		}

		// Bangun kelas_perkuliahan
		kelasPerkuliahan := make([]KelasPerkuliahan, maxLen)
		for i := 0; i < maxLen; i++ {
			kelasPerkuliahan[i] = KelasPerkuliahan{
				IDKelas:         idKelasList[i],
				NamaKelas:       namaKelasList[i],
				NamaMatakuliah:  namaMatakuliahList[i],
				KodeMatakuliah:  kodeMatakuliahList[i],
				IDDosenPengajar: idDosenPengajarList[i],
				Jadwal:          jadwalList[i],
			}
		}

		// Tambahkan ke response
		responses = append(responses, ListStudentKelasResponse{
			IDPesertaDidik:   model.IDPesertaDidik,
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

func (a *ApplicationServer) ListStudentKelas(c *fiber.Ctx) error {
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
			mahasiswa.nik AS nik,
			GROUP_CONCAT(kelaskuliah.id_kls ORDER BY kelaskuliah.id_kls SEPARATOR '|') AS id_kelas,
			GROUP_CONCAT(kelaskuliah.nm_kls ORDER BY kelaskuliah.id_kls SEPARATOR '|') AS nama_kelas,
			GROUP_CONCAT(matakuliah.nm_mk ORDER BY kelaskuliah.id_kls SEPARATOR '|') AS nama_matakuliah,
			GROUP_CONCAT(matakuliah.kode_mk ORDER BY kelaskuliah.id_kls SEPARATOR '|') AS kode_matakuliah,
			GROUP_CONCAT(akt_ajar_dosen.id_ptk ORDER BY akt_ajar_dosen.id_ptk SEPARATOR '|') AS id_dosen_pengajar,			
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
		Joins("LEFT JOIN akt_ajar_dosen ON akt_ajar_dosen.id_kls = kelaskuliah.id_kls").
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

func (a *ApplicationServer) GetTotalKelas(c *fiber.Ctx) error {
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
		Joins("LEFT JOIN akt_ajar_dosen ON akt_ajar_dosen.id_kls = kelaskuliah.id_kls").
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

type ListKelasResponse struct {
	IDKelas            string `json:"id_kelas" gorm:"column:id_kelas"`
	NamaKelas          string `json:"nama_kelas" gorm:"column:nama_kelas"`
	NamaMataKuliah     string `json:"nama_matakuliah" gorm:"column:nama_matakuliah"`
	KodeMataKuliah     string `json:"kode_matakuliah" gorm:"column:kode_matakuliah"`
	IDPTKDosenPengajar string `json:"id_dosen_pengajar" gorm:"column:id_dosen_pengajar"`
	Semester           string `json:"semester" gorm:"column:semester"`
	Jadwal             string `json:"jadwal" gorm:"column:jadwal"`
}

func (a *ApplicationServer) TotalKelas(c *fiber.Ctx) error {
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

func (a *ApplicationServer) ListKelas(c *fiber.Ctx) error {
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
	listKelas := make([]ListKelasResponse, 0)

	offset := req.Filter.GetOffset()
	limit := req.Filter.GetLimit()

	q := a.db.
		Table("kelaskuliah").
		Select(`
			kelaskuliah.id_kls AS id_kelas,
			kelaskuliah.nm_kls AS nama_kelas,
			matakuliah.nm_mk AS nama_matakuliah,
			matakuliah.kode_mk AS kode_matakuliah,
			akt_ajar_dosen.id_ptk AS id_dosen_pengajar,
			kelaskuliah.id_smt AS semester,
			GROUP_CONCAT(
				CONCAT(
				CASE jadwal.hari
					WHEN '1' THEN 'Senin'
					WHEN '2' THEN 'Selasa'
					WHEN '3' THEN 'Rabu'
					WHEN '4' THEN 'Kamis'
					WHEN '5' THEN 'Jumat'
					WHEN '6' THEN 'Sabtu'
					WHEN '7' THEN 'Minggu'
					ELSE 'Unknown'
				END,
				'-', jadwal.jam_mulai, 
				'-', jadwal.jam_selesai)
				ORDER BY jadwal.hari, jadwal.jam_mulai ASC
				SEPARATOR '|'
			) AS jadwal
		`).
		Joins("JOIN matakuliah_kurikulum ON matakuliah_kurikulum.id_mk_kur = kelaskuliah.id_mk_kur").
		Joins("JOIN matakuliah ON matakuliah.id_mk = matakuliah_kurikulum.id_mk").
		Joins("LEFT JOIN akt_ajar_dosen ON akt_ajar_dosen.id_kls = kelaskuliah.id_kls").
		Joins("LEFT JOIN jadwal ON jadwal.id_kls = kelaskuliah.id_kls").
		Where("kelaskuliah.id_smt = ?", req.Semester).
		Group("kelaskuliah.id_kls")

	// Menambahkan pencarian berdasarkan keyword
	if req.Filter.HasKeyword() {
		q = q.Where("kelaskuliah.nm_kls LIKE ?", "%"+req.Filter.Keyword+"%")
	}

	// Menambahkan sorting
	if req.Filter.HasSort() {
		q = q.Order(clause.OrderByColumn{
			Column: clause.Column{Name: req.Filter.SortBy},
			Desc:   req.Filter.IsDesc(),
		})
	} else {
		q = q.Order("kelaskuliah.id_kls ASC")
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
