package main

import (
	"net/http"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

func (a *ApplicationServer) ListRoomsMisca(c *fiber.Ctx) error {
	if IsSmartInstansi(c) {
		return a.ListRoomsSmart(c)
	}

	rooms := make([]Ruangan, 0)
	if err := a.db.Table("ruangan").Find(&rooms).Error; err != nil {
		return HandleError(c, err)
	}

	response := make([]RuanganResponse, 0)
	for _, r := range rooms {
		var idsms []string
		json.Unmarshal([]byte(r.IDSMSRaw), &idsms) // parsing string JSON ke slice

		response = append(response, RuanganResponse{
			IDRuangan:      r.IDRuangan,
			IDSMS:          idsms,
			NamaRuangan:    r.NamaRuangan,
			IDJenisRuangan: r.IDJenisRuangan,
			KodeRuangan:    r.KodeRuangan,
			Keterangan:     r.Keterangan,
			Kapasitas:      r.Kapasitas,
			CreatedAt:      r.CreatedAt,
			UpdatedAt:      r.UpdatedAt,
		})
	}

	return c.Status(fiber.StatusOK).JSON(ApiResponse[ListDataApiResponseWrapper[RuanganResponse]]{
		Code:    fiber.StatusOK,
		Status:  http.StatusText(fiber.StatusOK),
		Success: true,
		Message: "Sukses mendapatkan data ruangan",
		Data: ListDataApiResponseWrapper[RuanganResponse]{
			List: response,
		},
	})
}

func (a *ApplicationServer) GetTotalRoomsMisca(c *fiber.Ctx) error {
	if IsSmartInstansi(c) {
		return a.GetTotalRoomsSmart(c)
	}

	var total int64
	err := a.db.Table("ruangan").Count(&total).Error
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

func (a *ApplicationServer) ListRoomsSmart(c *fiber.Ctx) error {
	rooms := make([]Ruangan, 0)
	if err := a.db.
		Select("id_ruangan AS id_ruangan, id_sms AS id_sms, kode_ruangan AS kode_ruangan, kode_ruangan AS nama_ruangan, ket AS keterangan").
		Table("ruangan").Find(&rooms).Error; err != nil {
		return HandleError(c, err)
	}

	response := make([]RuanganResponse, 0)
	for _, r := range rooms {
		// cek jika r.IDSMSRaw ini tidak array -> ["86205","86206","87203","88201"]
		// maka jangan di unmarshal, masukkan langsung value r.IDSMSRaw ke slice idsms
		var idsms []string
		err := json.Unmarshal([]byte(r.IDSMSRaw), &idsms)
		if err != nil {
			// Jika gagal unmarshal DAN string-nya tidak kosong,
			// anggap sebagai ID tunggal dan masukkan ke slice.
			if r.IDSMSRaw != "" {
				idsms = []string{r.IDSMSRaw}
			}
			// Jika string kosong, idsms akan tetap menjadi slice kosong, yang sudah benar.
		}

		response = append(response, RuanganResponse{
			IDRuangan:      r.IDRuangan,
			IDSMS:          idsms,
			NamaRuangan:    r.NamaRuangan,
			IDJenisRuangan: r.IDJenisRuangan,
			KodeRuangan:    r.KodeRuangan,
			Keterangan:     r.Keterangan,
			Kapasitas:      r.Kapasitas,
			CreatedAt:      r.CreatedAt,
			UpdatedAt:      r.UpdatedAt,
		})
	}

	return c.Status(fiber.StatusOK).JSON(ApiResponse[ListDataApiResponseWrapper[RuanganResponse]]{
		Code:    fiber.StatusOK,
		Status:  http.StatusText(fiber.StatusOK),
		Success: true,
		Message: "Sukses mendapatkan data ruangan",
		Data: ListDataApiResponseWrapper[RuanganResponse]{
			List: response,
		},
	})
}

func (a *ApplicationServer) GetTotalRoomsSmart(c *fiber.Ctx) error {
	var total int64
	err := a.db.Table("ruangan").Count(&total).Error
	if err != nil {
		return HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(ApiResponse[GetTotalStudentsResponse]{
		Code:    fiber.StatusOK,
		Status:  http.StatusText(fiber.StatusOK),
		Success: true,
		Message: "Sukses mendapatkan total ruangan",
		Data: GetTotalStudentsResponse{
			Total: total,
		},
	})
}
