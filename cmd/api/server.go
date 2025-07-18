package main

import (
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
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

func (a *ApplicationServer) SetupHealthCheckRoutes() {
	a.router.Use(healthcheck.New(healthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		LivenessEndpoint: "/api/live",
		ReadinessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		ReadinessEndpoint: "/api/ready",
	}))
}

func (a *ApplicationServer) SetupRoutes() {
	a.router.Get("/api/misca/semesters", a.WithApiKey(), a.ListSemestersMisca)
	a.router.Get("/api/misca/semesters/active", a.WithApiKey(), a.GetActiveSemesterMisca)

	a.router.Get("/api/misca/students", a.WithApiKey(), a.ListStudentsMisca)
	a.router.Get("/api/misca/students/total", a.WithApiKey(), a.GetTotalStudentsMisca)

	a.router.Get("/api/misca/lecturers", a.WithApiKey(), a.ListLecturerMisca)
	a.router.Get("/api/misca/lecturers/total", a.WithApiKey(), a.GetTotalLecturerMisca)

	a.router.Get("/api/misca/classes", a.WithApiKey(), a.ListKelasMisca)
	a.router.Get("/api/misca/classes/total", a.WithApiKey(), a.TotalKelasMisca)

	a.router.Get("/api/misca/student_classes", a.WithApiKey(), a.ListSimpleStudentKelasMisca)
	a.router.Get("/api/misca/student_classes/total", a.WithApiKey(), a.TotalListSimpleStudentKelasMisca)

	a.router.Get("/api/misca/student_classes_details", a.WithApiKey(), a.ListStudentKelasDetailsMisca)
	a.router.Get("/api/misca/student_classes_details/total", a.WithApiKey(), a.GetTotalKelasDetailsMisca)

	a.router.Get("/api/misca/rooms", a.WithApiKey(), a.ListRoomsMisca)
	a.router.Get("/api/misca/rooms/total", a.WithApiKey(), a.GetTotalRoomsMisca)

	a.router.Get("/api/misca/sms", a.WithApiKey(), a.ListSMSMisca)
	a.router.Get("/api/misca/sms/total", a.WithApiKey(), a.GetTotalSMSMisca)
}

func (a *ApplicationServer) Run() {
	hostPort := net.JoinHostPort("127.0.0.1", a.config.AppPort)
	a.logger.With(slog.String("ON", hostPort)).Info("SERVER STARTED")
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
