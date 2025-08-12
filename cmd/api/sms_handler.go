package main

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (a *ApplicationServer) ListSMSMisca(c *fiber.Ctx) error {
	if IsSmartInstansi(c) {
		return a.ListSMSSmart(c)
	}

	sms := make([]SMS, 0)
	if err := a.db.Table("sms").
		Select("sms.*,jenjang_pendidikan.nama_jenjang_didik AS nama_jenjang_didik").
		Joins("LEFT JOIN jenjang_pendidikan ON sms.id_jenj_didik = jenjang_pendidikan.id_jenjang_didik").
		Scan(&sms).Error; err != nil {
		return HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(ApiResponse[ListDataApiResponseWrapper[SMS]]{
		Code:    fiber.StatusOK,
		Status:  http.StatusText(fiber.StatusOK),
		Success: true,
		Message: "Sukses mendapatkan data ruangan",
		Data: ListDataApiResponseWrapper[SMS]{
			List: sms,
		},
	})
}

func (a *ApplicationServer) GetTotalSMSMisca(c *fiber.Ctx) error {
	if IsSmartInstansi(c) {
		return a.GetTotalSMSSmart(c)
	}

	var total int64
	err := a.db.Table("sms").Count(&total).Error
	if err != nil {
		return HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(ApiResponse[GetTotalStudentsResponse]{
		Code:    fiber.StatusOK,
		Status:  http.StatusText(fiber.StatusOK),
		Success: true,
		Message: "Sukses mendapatkan total sms",
		Data: GetTotalStudentsResponse{
			Total: total,
		},
	})
}

func (a *ApplicationServer) ListSMSSmart(c *fiber.Ctx) error {
	var sms []SMS
	if err := a.db.Table("sms").
		Select(`
				sms.id_sms AS id_sms,
				sms.nm_lemb AS nm_lemb,
				sms.nm_lemb_english AS nm_lemb_inggris,
				sms.kode_prodi AS kode_sms,
				sms.id_jns_sms,
				jenjang_pendidikan.nm_jenj_didik AS nama_jenjang_didik`).
		Joins("LEFT JOIN jenjang_pendidikan ON sms.id_jenj_didik = jenjang_pendidikan.id_jenj_didik").
		Scan(&sms).Error; err != nil {
		return HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(ApiResponse[ListDataApiResponseWrapper[SMS]]{
		Code:    fiber.StatusOK,
		Status:  http.StatusText(fiber.StatusOK),
		Success: true,
		Message: "Sukses mendapatkan data ruangan",
		Data: ListDataApiResponseWrapper[SMS]{
			List: sms,
		},
	})
}

func (a *ApplicationServer) GetTotalSMSSmart(c *fiber.Ctx) error {
	var total int64
	err := a.db.Table("sms").Count(&total).Error
	if err != nil {
		return HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(ApiResponse[GetTotalStudentsResponse]{
		Code:    fiber.StatusOK,
		Status:  http.StatusText(fiber.StatusOK),
		Success: true,
		Message: "Sukses mendapatkan total sms",
		Data: GetTotalStudentsResponse{
			Total: total,
		},
	})
}
