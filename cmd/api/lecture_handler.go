package main

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
	gl "lab.garudacyber.co.id/g-learning-connector"
)

type (
	GetTotalLecturerResponse struct {
		Total int64 `json:"total"`
	}

	ListLecturerRequest struct {
		gl.Filter
	}

	ListLecturerResponse struct {
		ID        string `json:"id" gorm:"column:id_ptk"`
		Name      string `json:"name" gorm:"column:nama_dosen"`
		Gender    string `json:"gender" gorm:"column:jenis_kelamin"`
		NIK       string `json:"nik" gorm:"column:nik"`
		Email     string `json:"email" gorm:"column:email"`
		Handphone string `json:"handphone" gorm:"column:handphone"`
		Telephone string `json:"telephone" gorm:"column:telepon"`
	}
)

func NewListLecturerRequest() *ListLecturerRequest {
	return &ListLecturerRequest{
		Filter: gl.NewFilterPagination(),
	}
}

func (a *ApplicationServer) ListLecturer(c *fiber.Ctx) error {
	req := NewListLecturerRequest()
	if err := c.QueryParser(req); err != nil {
		return HandleError(c, err)
	}

	listLecturer := make([]ListLecturerResponse, 0)

	offset := req.Filter.GetOffset()
	limit := req.Filter.GetLimit()

	q := a.db.Select(`id_ptk, nama_dosen, jenis_kelamin, nik, email, handphone, telepon`).
		Table("dosen").
		Where("nik IS NOT NULL AND nik != '' AND LENGTH(nik) = 16")

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

func (a *ApplicationServer) GetTotalLecturer(c *fiber.Ctx) error {
	var total int64

	err := a.db.Table("dosen").Count(&total).
		Where("nik IS NOT NULL AND nik != '' AND LENGTH(nik) = 16").
		Error

	if err != nil {
		return HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(ApiResponse[GetTotalLecturerResponse]{
		Code:    fiber.StatusOK,
		Status:  http.StatusText(fiber.StatusOK),
		Success: true,
		Message: "Sukses mendapatkan total dosen",
		Data: GetTotalLecturerResponse{
			Total: total,
		},
	})
}
