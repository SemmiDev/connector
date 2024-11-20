package main

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
	gl "lab.garudacyber.co.id/g-learning-connector"
	"net/http"
)

type (
	GetTotalStudentsResponse struct {
		Total int64 `json:"total"`
	}

	ListStudentsRequest struct {
		gl.Filter
	}

	ListStudentsResponse struct {
		ID        string `json:"id" gorm:"column:id"`
		Name      string `json:"name" gorm:"column:nama_mahasiswa"`
		Gender    string `json:"gender" gorm:"column:jenis_kelamin"`
		NIK       string `json:"nik" gorm:"column:nik"`
		Email     string `json:"email" gorm:"column:email"`
		Handphone string `json:"handphone" gorm:"column:handphone"`
		Telephone string `json:"telephone" gorm:"column:telepon"`
	}
)

func NewListStudentsRequest() *ListStudentsRequest {
	return &ListStudentsRequest{
		Filter: gl.NewFilterPagination(),
	}
}

func (a *ApplicationServer) ListStudents(c *fiber.Ctx) error {
	req := NewListStudentsRequest()
	if err := c.QueryParser(req); err != nil {
		return HandleError(c, err)
	}

	listStudents := make([]ListStudentsResponse, 0)

	offset := req.Filter.GetOffset()
	limit := req.Filter.GetLimit()

	q := a.db.
		Select(`
		mahasiswa_histori.id_pd AS id, 
		mahasiswa.nama_mahasiswa AS nama_mahasiswa, 
		mahasiswa.jenis_kelamin AS jenis_kelamin, 
		mahasiswa.nik AS nik, 
		mahasiswa.email AS email, 
		mahasiswa.handphone AS handphone, 
		mahasiswa.telepon AS telepon`).
		Table("mahasiswa_histori").
		Joins("INNER JOIN mahasiswa ON mahasiswa_histori.id_mahasiswa = mahasiswa.id").
		Where("mahasiswa.nik IS NOT NULL AND mahasiswa.nik != '' AND LENGTH(mahasiswa.nik) = 16 AND mahasiswa_histori.deleted_at IS NULL")

	if req.Filter.HasKeyword() {
		q = q.Where("mahasiswa.nama_mahasiswa LIKE ? OR mahasiswa.nik LIKE ?", "%"+req.Filter.Keyword+"%", "%"+req.Filter.Keyword+"%")
	}

	if req.Filter.HasSort() {
		sortBy := req.Filter.SortBy
		if sortBy == "id" {
			sortBy = "mahasiswa_histori.id_pd"
		} else {
			sortBy = "mahasiswa." + sortBy
		}

		q = q.Order(
			clause.OrderByColumn{
				Column: clause.Column{Name: sortBy},
				Desc:   req.Filter.IsDesc(),
			},
		)
	} else {
		q = q.Order("mahasiswa_histori.created_at ASC")
	}

	// Menghitung jumlah total data tanpa offset dan limit
	var totalData int64
	if err := q.Count(&totalData).Error; err != nil {
		return HandleError(c, err)
	}

	// Menambahkan limit dan offset setelah menghitung total data
	q = q.Offset(int(offset)).Limit(int(limit))

	if err := q.Scan(&listStudents).Error; err != nil {
		return HandleError(c, err)
	}

	pageInfo, err := gl.NewPageInfo(req.Filter.CurrentPage, limit, offset, totalData)
	if err != nil {
		return HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(ApiResponse[ListDataApiResponseWrapper[ListStudentsResponse]]{
		Code:    fiber.StatusOK,
		Status:  http.StatusText(fiber.StatusOK),
		Success: true,
		Message: "Sukses mendapatkan data mahasiswa",
		Data: ListDataApiResponseWrapper[ListStudentsResponse]{
			List:     listStudents,
			PageInfo: pageInfo,
		},
	})
}

func (a *ApplicationServer) GetTotalStudents(c *fiber.Ctx) error {
	var total int64

	err := a.db.
		Table("mahasiswa").
		Where("nik IS NOT NULL AND nik != '' AND LENGTH(nik) = 16 AND deleted_at IS NULL").
		Count(&total).
		Error

	if err != nil {
		return HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(ApiResponse[GetTotalStudentsResponse]{
		Code:    fiber.StatusOK,
		Status:  http.StatusText(fiber.StatusOK),
		Success: true,
		Message: "Sukses mendapatkan total mahasiswa",
		Data: GetTotalStudentsResponse{
			Total: total,
		},
	})
}
