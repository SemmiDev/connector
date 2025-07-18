package main

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type (
	ListSemestersResponse struct {
		ID     string `gorm:"column:id_smt" json:"id"`
		Name   string `gorm:"column:nm_smt" json:"name"`
		Active uint8  `gorm:"column:active" json:"active"`
	}

	GetActiveSemester struct {
		ID   string `gorm:"column:id_smt" json:"id"`
		Name string `gorm:"column:nm_smt" json:"name"`
	}
)

func (a *ApplicationServer) ListSemestersMisca(c *fiber.Ctx) error {
	if IsSmartInstansi(c) {
		return a.ListSemestersSmart(c)
	}

	semesters := make([]ListSemestersResponse, 0)

	err := a.db.
		Table("semester").
		Select(`
			semester.id_smt AS id_smt,
			semester.nm_smt AS nm_smt,
			CASE WHEN setting.param = 'periode_berlaku' THEN 1 ELSE 0 END AS active
		`).
		Joins("LEFT JOIN setting ON semester.id_smt = setting.value AND setting.param = 'periode_berlaku'").
		Find(&semesters).Error

	if err != nil {
		return HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(ApiResponse[ListDataApiResponseWrapper[ListSemestersResponse]]{
		Code:    fiber.StatusOK,
		Status:  http.StatusText(fiber.StatusOK),
		Success: true,
		Message: "Sukses mendapatkan semua data semester",
		Data: ListDataApiResponseWrapper[ListSemestersResponse]{
			List: semesters,
		},
	})
}

func (a *ApplicationServer) GetActiveSemesterMisca(c *fiber.Ctx) error {
	if IsSmartInstansi(c) {
		return a.GetActiveSemesterSmart(c)
	}

	var semester GetActiveSemester

	err := a.db.
		Table("semester").
		Select(`
			semester.id_smt AS id_smt,
			semester.nm_smt AS nm_smt
		`).
		Joins("LEFT JOIN setting ON semester.id_smt = setting.value AND setting.param = 'periode_berlaku'").
		Where("setting.param = 'periode_berlaku'").
		Scan(&semester).
		Error

	if err != nil {
		return HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(ApiResponse[GetActiveSemester]{
		Code:    fiber.StatusOK,
		Status:  http.StatusText(fiber.StatusOK),
		Success: true,
		Message: "Sukses mendapatkan data semester yang aktif",
		Data:    semester,
	})
}

func (a *ApplicationServer) ListSemestersSmart(c *fiber.Ctx) error {
	semesters := make([]ListSemestersResponse, 0)

	err := a.db.
		Table("semester").
		Select(`id_smt, nm_smt, a_periode_aktif AS active`).
		Find(&semesters).Error

	if err != nil {
		return HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(ApiResponse[ListDataApiResponseWrapper[ListSemestersResponse]]{
		Code:    fiber.StatusOK,
		Status:  http.StatusText(fiber.StatusOK),
		Success: true,
		Message: "Sukses mendapatkan semua data semester",
		Data: ListDataApiResponseWrapper[ListSemestersResponse]{
			List: semesters,
		},
	})
}

func (a *ApplicationServer) GetActiveSemesterSmart(c *fiber.Ctx) error {
	var semester GetActiveSemester

	err := a.db.Table("semester").Select(`id_smt, nm_smt`).Where("a_periode_aktif = 1").Scan(&semester).Error
	if err != nil {
		return HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(ApiResponse[GetActiveSemester]{
		Code:    fiber.StatusOK,
		Status:  http.StatusText(fiber.StatusOK),
		Success: true,
		Message: "Sukses mendapatkan data semester yang aktif",
		Data:    semester,
	})
}
