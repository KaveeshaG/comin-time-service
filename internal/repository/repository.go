package repository

import (
	"time"

	"github.com/Axontik/comin-time-service/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TimeRepository interface {
	// Attendance methods
	CreateAttendance(attendance *domain.Attendance) error
	GetAttendanceByDate(employeeID uuid.UUID, date time.Time) (*domain.Attendance, error)
	UpdateAttendance(attendance *domain.Attendance) error
	ListAttendances(orgID, employeeID uuid.UUID, startDate, endDate time.Time) ([]domain.Attendance, error)

	// QR Code methods
	CreateQRCode(qrCode *domain.QRCode) error
	GetQRCode(code string) (*domain.QRCode, error)
	UpdateQRCode(qrCode *domain.QRCode) error
	ListQRCodes(employeeID uuid.UUID) ([]domain.QRCode, error)
	GetEmployeeQRCodes(orgID, employeeID uuid.UUID) ([]domain.QRCode, error)

	// Timesheet methods
	CreateTimesheet(timesheet *domain.Timesheet) error
	GetTimesheet(id uuid.UUID) (*domain.Timesheet, error)
	UpdateTimesheet(timesheet *domain.Timesheet) error
	DeleteTimesheet(id uuid.UUID) error
	ListTimesheets(orgID, employeeID uuid.UUID, startDate, endDate time.Time) ([]domain.Timesheet, error)
	GetTimesheetSummary(orgID, employeeID uuid.UUID, startDate, endDate time.Time) (float64, error)
}

type timeRepository struct {
	db *gorm.DB
}

func NewTimeRepository(db *gorm.DB) TimeRepository {
	return &timeRepository{db: db}
}

func (r *timeRepository) CreateAttendance(attendance *domain.Attendance) error {
	result := r.db.Create(attendance)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *timeRepository) GetAttendanceByDate(employeeID uuid.UUID, date time.Time) (*domain.Attendance, error) {
	attendance := &domain.Attendance{}
	err := r.db.Where("employee_id = ? AND date = ?", employeeID, date).First(attendance).Error
	if err != nil {
		return nil, err
	}
	return attendance, nil
}

func (r *timeRepository) UpdateAttendance(attendance *domain.Attendance) error {
	return r.db.Model(&domain.Attendance{}).Updates(attendance).Error
}

func (r *timeRepository) ListAttendances(orgID, employeeID uuid.UUID, startDate, endDate time.Time) ([]domain.Attendance, error) {
	attendances := []domain.Attendance{}
	err := r.db.Where("organization_id = ? AND employee_id = ? AND date BETWEEN ? AND ?", orgID, employeeID, startDate, endDate).Find(&attendances).Error
	if err != nil {
		return nil, err
	}
	return attendances, nil
}

func (r *timeRepository) CreateQRCode(qrCode *domain.QRCode) error {
	return r.db.Create(qrCode).Error
}

func (r *timeRepository) GetQRCode(code string) (*domain.QRCode, error) {
	qrCode := &domain.QRCode{}
	err := r.db.Where("code = ?", code).First(qrCode).Error
	if err != nil {
		return nil, err
	}
	return qrCode, nil
}

func (r *timeRepository) UpdateQRCode(qrCode *domain.QRCode) error {
	return r.db.Model(&domain.QRCode{}).Updates(qrCode).Error
}

func (r *timeRepository) ListQRCodes(employeeID uuid.UUID) ([]domain.QRCode, error) {
	qrCodes := []domain.QRCode{}
	err := r.db.Where("employee_id = ?", employeeID).Find(&qrCodes).Error
	if err != nil {
		return nil, err
	}
	return qrCodes, nil
}

func (r *timeRepository) CreateTimesheet(timesheet *domain.Timesheet) error {
	return r.db.Create(timesheet).Error
}

func (r *timeRepository) GetTimesheet(id uuid.UUID) (*domain.Timesheet, error) {
	timesheet := &domain.Timesheet{}
	err := r.db.Where("id = ?", id).First(timesheet).Error
	if err != nil {
		return nil, err
	}
	return timesheet, nil
}

func (r *timeRepository) UpdateTimesheet(timesheet *domain.Timesheet) error {
	return r.db.Model(&domain.Timesheet{}).Updates(timesheet).Error
}

func (r *timeRepository) DeleteTimesheet(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&domain.Timesheet{}).Error
}

func (r *timeRepository) ListTimesheets(orgID, employeeID uuid.UUID, startDate, endDate time.Time) ([]domain.Timesheet, error) {
	timesheets := []domain.Timesheet{}
	err := r.db.Where("organization_id = ? AND employee_id = ? AND date BETWEEN ? AND ?", orgID, employeeID, startDate, endDate).Find(&timesheets).Error
	if err != nil {
		return nil, err
	}
	return timesheets, nil
}

func (r *timeRepository) GetTimesheetSummary(orgID, employeeID uuid.UUID, startDate, endDate time.Time) (float64, error) {
	var totalHours float64
	err := r.db.Model(&domain.Timesheet{}).Where("organization_id = ? AND employee_id = ? AND date BETWEEN ? AND ?", orgID, employeeID, startDate, endDate).Select("SUM(hours)").Find(&totalHours).Error
	if err != nil {
		return 0, err
	}
	return totalHours, nil
}

func (r *timeRepository) GetEmployeeQRCodes(orgID, employeeID uuid.UUID) ([]domain.QRCode, error) {
	qrCodes := []domain.QRCode{}
	err := r.db.Where("organization_id = ? AND employee_id = ?", orgID, employeeID).Find(&qrCodes).Error
	if err != nil {
		return nil, err
	}
	return qrCodes, nil
}
