package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	fiberRecovery "github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
	gl "lab.garudacyber.co.id/g-learning-connector"
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

	a.router.Get("/api/misca/sms", a.WithApiKey(), a.ListSMS)
	a.router.Get("/api/misca/sms/total", a.WithApiKey(), a.GetTotalSMS)

	// SMART
	a.router.Get("/api/smart/semesters", a.WithApiKey(), a.ListSemestersSmart)
	a.router.Get("/api/smart/semesters/active", a.WithApiKey(), a.GetActiveSemesterSmart)

	a.router.Get("/api/smart/lecturers", a.WithApiKey(), a.ListLecturerSmart)
	a.router.Get("/api/smart/lecturers/total", a.WithApiKey(), a.GetTotalLecturerSmart)

	a.router.Get("/api/smart/students", a.WithApiKey(), a.ListStudentsSmart)
	a.router.Get("/api/smart/students/total", a.WithApiKey(), a.GetTotalStudentsSmart)
}

func (a *ApplicationServer) Run() {
	host := "0.0.0.0"
	port := a.config.AppPort
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

func (a *ApplicationServer) ListRooms(c *fiber.Ctx) error {
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

func (a *ApplicationServer) ListSMS(c *fiber.Ctx) error {
	var sms []SMS
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

func (a *ApplicationServer) GetTotalSMS(c *fiber.Ctx) error {
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
