package main

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
	gl "lab.garudacyber.co.id/g-learning-connector"
	"net/http"
)

type (
	GetTotalKelasResponse struct {
		Total int64 `json:"total"`
	}

	ListKelasRequest struct {
		gl.Filter
	}

	ListKelasResponse struct {
		ID   string `json:"id" gorm:"column:id_ptk"`
		Name string `json:"name" gorm:"column:nama_dosen"`
	}
)

func NewListKelasRequest() *ListKelasRequest {
	return &ListKelasRequest{
		Filter: gl.NewFilterPagination(),
	}
}

func (a *ApplicationServer) ListKelas(c *fiber.Ctx) error {
	req := NewListKelasRequest()
	if err := c.QueryParser(req); err != nil {
		return HandleError(c, err)
	}

	listLecturer := make([]ListLecturerResponse, 0)

	offset := req.Filter.GetOffset()
	limit := req.Filter.GetLimit()

	q := a.db.
		Select(`id_ptk, nama_dosen, jenis_kelamin, nik, email, handphone, telepon`).
		Table("dosen")

	if req.Filter.HasKeyword() {
		q = q.Where("nama_dosen LIKE ? OR nik LIKE ?", "%"+req.Filter.Keyword+"%", "%"+req.Filter.Keyword+"%")
	}

	if req.Filter.HasSort() {
		q = q.Order(
			clause.OrderByColumn{
				Column: clause.Column{Name: req.Filter.SortBy},
				Desc:   req.Filter.IsDesc(),
			},
		)
	} else {
		q = q.Order("created_at ASC")
	}

	// Menghitung jumlah total data tanpa offset dan limit
	var totalData int64
	if err := q.Count(&totalData).Error; err != nil {
		return HandleError(c, err)
	}

	// Menambahkan limit dan offset setelah menghitung total data
	q = q.Offset(int(offset)).Limit(int(limit))

	if err := q.Scan(&listLecturer).Error; err != nil {
		return HandleError(c, err)
	}

	pageInfo, err := gl.NewPageInfo(req.Filter.CurrentPage, limit, offset, totalData)
	if err != nil {
		return HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(ApiResponse[ListDataApiResponseWrapper[ListLecturerResponse]]{
		Code:    fiber.StatusOK,
		Status:  http.StatusText(fiber.StatusOK),
		Success: true,
		Message: "Sukses mendapatkan data dosen",
		Data: ListDataApiResponseWrapper[ListLecturerResponse]{
			List:     listLecturer,
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
