package main

import (
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	fiberRecovery "github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
	gl "lab.garudacyber.co.id/g-learning-connector"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type ApplicationServer struct {
	config *gl.Config
	logger *slog.Logger
	db     *gorm.DB
	router *fiber.App
}

func NewApplicationServer(db *gorm.DB, logger *slog.Logger, config *gl.Config, router *fiber.App) *ApplicationServer {
	app := ApplicationServer{
		config: config,
		logger: logger,
		db:     db,
		router: router,
	}

	return &app
}

func (a *ApplicationServer) SetupCommonMiddlewares() {
	a.router.Use(cors.New())
	a.router.Use(fiberRecovery.New())
}

func (a *ApplicationServer) SetupHealthCheckRoutes() {
	a.router.Use(healthcheck.New(healthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		LivenessEndpoint: "/live",
		ReadinessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		ReadinessEndpoint: "/ready",
	}))
}

func (a *ApplicationServer) SetupRoutes() {
	a.router.Get("/api/misca/semesters", a.WithApiKey(), a.ListSemesters)
	a.router.Get("/api/misca/semesters/active", a.WithApiKey(), a.GetActiveSemester)

	a.router.Get("/api/misca/students", a.WithApiKey(), a.ListStudents)
	a.router.Get("/api/misca/students/total", a.WithApiKey(), a.GetTotalStudents)

	a.router.Get("/api/misca/lecturers", a.WithApiKey(), a.ListLecturer)
	a.router.Get("/api/misca/lecturers/total", a.WithApiKey(), a.GetTotalLecturer)

	a.router.Get("/api/misca/classes", a.WithApiKey(), a.ListKelas)
	a.router.Get("/api/misca/classes/total", a.WithApiKey(), a.TotalKelas)

	a.router.Get("/api/misca/student_classes", a.WithApiKey(), a.ListSimpleStudentKelas)
	a.router.Get("/api/misca/student_classes/total", a.WithApiKey(), a.TotalListSimpleStudentKelas)

	a.router.Get("/api/misca/student_classes_details", a.WithApiKey(), a.ListStudentKelasDetails)
	a.router.Get("/api/misca/student_classes_details/total", a.WithApiKey(), a.GetTotalKelasDetails)

	a.router.Get("/api/misca/rooms", a.WithApiKey(), a.ListRooms)
	a.router.Get("/api/misca/rooms/total", a.WithApiKey(), a.GetTotalRooms)

	a.router.Get("/api/misca/tes", a.Test)
}

func (a *ApplicationServer) Run() {
	host := "0.0.0.0"
	port := fmt.Sprintf("%s", a.config.AppPort)
	hostPort := fmt.Sprintf("%s:%s", host, port)

	a.logger.With(slog.String("host", host), slog.String("port", port)).Info("Server started")

	err := a.router.Listen(hostPort)
	gl.PanicIfNeeded(err)
}

type BobotNilai struct {
	ID              int        `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	IDSMS           string     `json:"id_sms" gorm:"column:id_sms"`
	NilaiHuruf      string     `json:"nilai_huruf" gorm:"column:nilai_huruf"`
	Min             string     `json:"min" gorm:"column:min"`
	Maks            string     `json:"maks" gorm:"column:maks"`
	NilaiIndeks     string     `json:"nilai_indeks" gorm:"column:nilai_indeks"`
	Angkatan        string     `json:"angkatan" gorm:"column:angkatan"`
	TglMulaiEfektif *time.Time `json:"tgl_mulai_efektif" gorm:"column:tgl_mulai_efektif"`
	TglAkhirEfektif *time.Time `json:"tgl_akhir_efektif" gorm:"column:tgl_akhir_efektif"`
	CreatedAt       *time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt       *time.Time `json:"updated_at" gorm:"column:updated_at"`
}

type KelasKuliahBobotNilai struct {
	ID            int        `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	IDKelas       int        `json:"id_kls" gorm:"column:id_kls"`
	BobotAbsensi  int        `json:"bobot_absensi" gorm:"column:bobot_absensi;default:10"`
	BobotTugas    int        `json:"bobot_tugas" gorm:"column:bobot_tugas;default:20"`
	BobotUTS      int        `json:"bobot_uts" gorm:"column:bobot_uts;default:30"`
	BobotUAS      int        `json:"bobot_uas" gorm:"column:bobot_uas;default:40"`
	BobotJSON     string     `json:"bobot_json" gorm:"column:bobot_json;type:longtext"`
	JmlhPertemuan int        `json:"jmlh_pertemuan" gorm:"column:jmlh_pertemuan;default:16"`
	IDUnsur       *int       `json:"id_unsur" gorm:"column:id_unsur"`
	CreatedAt     *time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt     *time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (k KelasKuliahBobotNilai) Exists() bool {
	return k.ID != 0
}

type UnsurNilai struct {
	ID            int        `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	Nama          string     `json:"nama" gorm:"column:nama"`
	IDSMT         string     `json:"id_smt" gorm:"column:id_smt"`
	IDSMS         string     `json:"id_sms" gorm:"column:id_sms;type:longtext"`
	TipeKuliah    string     `json:"tipe_kuliah" gorm:"column:tipe_kuliah;type:enum('teori','praktikum','teori_praktikum');default:'teori_praktikum'"`
	TipePenilaian string     `json:"tipe_penilaian" gorm:"column:tipe_penilaian;type:enum('detail','angka');default:'detail'"`
	Unsur         string     `json:"unsur" gorm:"column:unsur;type:longtext"`
	CreatedAt     *time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt     *time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (u UnsurNilai) Exists() bool {
	return u.ID != 0
}

func (BobotNilai) TableName() string {
	return "bobot_nilai"
}

func (KelasKuliahBobotNilai) TableName() string {
	return "kelaskuliah_bobot_nilai"
}

func (UnsurNilai) TableName() string {
	return "unsur_nilai"
}

func GetUnsurNilai(db *gorm.DB, idSMS, idSMT, tipeKuliah, tipePenilaian string) (*UnsurNilai, error) {
	var unsur UnsurNilai
	query := db.Where("JSON_CONTAINS(id_sms, ?, '$')", fmt.Sprintf(`"%s"`, idSMS)).
		Where("id_smt = ?", idSMT).
		Where("tipe_penilaian = ?", tipePenilaian)

	if tipePenilaian == "detail" {
		query = query.Where("tipe_kuliah = ? OR tipe_kuliah = ?", tipeKuliah, "teori_praktikum")
	}

	result := query.First(&unsur)
	return &unsur, result.Error
}

func (a *ApplicationServer) Test(c *fiber.Ctx) error {

	idKls := "3448"
	idSmt := 20242

	// ambil informasi kelas kuliah
	type KelasKuliah struct {
		IDKls         int    `json:"id_kls" gorm:"column:id_kls"`
		BukaNilai     int    `json:"buka_nilai" gorm:"column:buka_nilai"`
		BobotAbsensi  int    `json:"bobot_absensi" gorm:"column:bobot_absensi"`
		BobotTugas    int    `json:"bobot_tugas" gorm:"column:bobot_tugas"`
		BobotUTS      int    `json:"bobot_uts" gorm:"column:bobot_uts"`
		BobotUAS      int    `json:"bobot_uas" gorm:"column:bobot_uas"`
		IDMkKur       int    `json:"id_mk_kur" gorm:"column:id_mk_kur"`
		IDSMS         int    `json:"id_sms" gorm:"column:id_sms"`
		IDSMT         int    `json:"id_smt" gorm:"column:id_smt"`
		TipeKuliah    string `json:"tipe_kuliah" gorm:"column:tipe_kuliah"`
		TipePenilaian string `json:"tipe_penilaian" gorm:"column:tipe_penilaian"`
	}

	var kelasKuliah KelasKuliah
	err := a.db.Table("kelaskuliah AS kk").
		Select(`
            kk.id_kls AS id_kls,
            s.buka_nilai AS buka_nilai,
			kbn.bobot_absensi AS bobot_absensi,
            kbn.bobot_tugas AS bobot_tugas,
            kbn.bobot_uts AS bobot_uts,
            kbn.bobot_uas AS bobot_uas,
			kk.id_mk_kur AS id_mk_kur,
            mk.tipe_kuliah AS tipe_kuliah,
            mk.tipe_penilaian AS tipe_penilaian,
            kk.id_sms AS id_sms,
            kk.id_smt AS id_smt
        `).
		Joins("JOIN sms s ON kk.id_sms = s.id_sms").
		Joins("JOIN kelaskuliah_bobot_nilai kbn ON kbn.id_kls = kk.id_kls").
		Joins("JOIN matakuliah_kurikulum mk ON mk.id_mk_kur = kk.id_mk_kur").
		Where("kk.id_kls = ? AND kk.id_smt = ?", idKls, idSmt).
		First(&kelasKuliah).Error

	if err != nil {
		return HandleError(c, err)
	}

	// cek status buka nilai
	if kelasKuliah.BukaNilai != 1 {
		return HandleError(c, errors.New("periode nilai belum dibuka"))
	}

	// ambil informasi bobot nilai kelas kuliah
	var bobotKelasKuliah KelasKuliahBobotNilai
	err = a.db.Table("kelaskuliah_bobot_nilai").Where("id_kls = ?", idKls).First(&bobotKelasKuliah).Error
	if err != nil {
		return HandleError(c, err)
	}

	// jika bobot json kosong, ambil unsur nilai
	if bobotKelasKuliah.BobotJSON == "" {
		unsurNilai, err := GetUnsurNilai(a.db, strconv.Itoa(kelasKuliah.IDSMS), strconv.Itoa(kelasKuliah.IDSMT), kelasKuliah.TipeKuliah, kelasKuliah.TipePenilaian)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return HandleError(c, err)
		}
		if unsurNilai.Exists() {
			bobotKelasKuliah.BobotJSON = unsurNilai.Unsur
			bobotKelasKuliah.IDUnsur = &unsurNilai.ID
		}
	}

	if !bobotKelasKuliah.Exists() {
		bobotKelasKuliah.IDKelas = kelasKuliah.IDKls
		a.db.Save(&bobotKelasKuliah)
	}

	return c.JSON(fiber.Map{
		"kelas_kuliah":             kelasKuliah,
		"kelas_kuliah_bobot_nilai": bobotKelasKuliah,
	})
}

func (a *ApplicationServer) ListRooms(c *fiber.Ctx) error {
	var rooms []Ruangan
	if err := a.db.Table("ruangan").Find(&rooms).Error; err != nil {
		return HandleError(c, err)
	}

	var response []RuanganResponse
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

func (a *ApplicationServer) GetTotalRooms(c *fiber.Ctx) error {
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
