package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"time"

	"github.com/Axontik/comin-time-service/internal/domain"
	"github.com/Axontik/comin-time-service/internal/repository"
	"github.com/google/uuid"
)

type TimeService interface {
	// QR Code methods
	GenerateQRCode(orgID uuid.UUID, req *domain.GenerateQRRequest) (*domain.QRCode, error)
	ValidateQRCode(code string) (*domain.QRCode, error)
	GetEmployeeQRCodes(orgID, employeeID uuid.UUID) ([]domain.QRCode, error)

	// Attendance methods
	CheckIn(req *domain.CheckInRequest) (*domain.Attendance, error)
	CheckOut(req *domain.CheckOutRequest) (*domain.Attendance, error)
	GetAttendanceByDate(employeeID uuid.UUID, date time.Time) (*domain.Attendance, error)
	GetAttendanceSummary(employeeID uuid.UUID, month, year int) (map[string]int, error)
	ListAttendances(orgID uuid.UUID) ([]domain.Attendance, error)

	// Timesheet methods
	CreateTimesheet(orgID, employeeID uuid.UUID, req *domain.CreateTimesheetRequest) (*domain.Timesheet, error)
	GetTimesheet(id uuid.UUID) (*domain.Timesheet, error)
	UpdateTimesheet(id uuid.UUID, req *domain.CreateTimesheetRequest) (*domain.Timesheet, error)
	DeleteTimesheet(id uuid.UUID) error
	ListTimesheets(orgID, employeeID uuid.UUID, startDate, endDate time.Time) ([]domain.Timesheet, error)
	ApproveTimesheet(id, approverID uuid.UUID) error
	RejectTimesheet(id, approverID uuid.UUID, reason string) error
}

type timeService struct {
	timeRepo repository.TimeRepository
}

func NewTimeService(timeRepo repository.TimeRepository) TimeService {
	return &timeService{
		timeRepo: timeRepo,
	}
}

// Generate QR Code for employee
func (s *timeService) GenerateQRCode(orgID uuid.UUID, req *domain.GenerateQRRequest) (*domain.QRCode, error) {
	// Generate unique code
	codeBytes := make([]byte, 32)
	_, err := rand.Read(codeBytes)
	if err != nil {
		return nil, err
	}
	code := base64.URLEncoding.EncodeToString(codeBytes)

	// Set expiry date if specified
	var expiryDate *time.Time
	if req.ExpiryDays > 0 {
		exp := time.Now().AddDate(0, 0, req.ExpiryDays)
		expiryDate = &exp
	}

	qrCode := &domain.QRCode{
		OrganizationID: orgID,
		EmployeeID:     req.EmployeeID,
		Code:           code,
		ExpiryDate:     expiryDate,
		IsActive:       true,
	}

	if err := s.timeRepo.CreateQRCode(qrCode); err != nil {
		return nil, err
	}

	return qrCode, nil
}

// Validate QR Code
func (s *timeService) ValidateQRCode(code string) (*domain.QRCode, error) {
	qrCode, err := s.timeRepo.GetQRCode(code)
	if err != nil {
		return nil, errors.New("invalid QR code")
	}

	if !qrCode.IsActive {
		return nil, errors.New("QR code is inactive")
	}

	if qrCode.ExpiryDate != nil && time.Now().After(*qrCode.ExpiryDate) {
		return nil, errors.New("QR code has expired")
	}

	return qrCode, nil
}

// Check-in employee
func (s *timeService) CheckIn(req *domain.CheckInRequest) (*domain.Attendance, error) {
	// Validate QR code
	qrCode, err := s.ValidateQRCode(req.QRCode)
	if err != nil {
		return nil, err
	}

	// Use provided timestamp or current time
	checkInTime := req.Timestamp
	if checkInTime.IsZero() {
		checkInTime = time.Now()
	}

	// Check if employee already checked in today
	today := time.Now().Truncate(24 * time.Hour)
	attendance, err := s.timeRepo.GetAttendanceByDate(qrCode.EmployeeID, today)
	if err == nil && attendance.CheckIn != nil {
		return nil, errors.New("already checked in today")
	}

	// Create new attendance record
	attendance = &domain.Attendance{
		OrganizationID: qrCode.OrganizationID,
		EmployeeID:     qrCode.EmployeeID,
		CheckIn:        &checkInTime,
		Date:           today,
		Status:         domain.AttendanceStatusPresent,
		WorkMode:       req.WorkMode,
		Location:       req.Location,
		DeviceInfo:     req.DeviceInfo,
	}

	if err := s.timeRepo.CreateAttendance(attendance); err != nil {
		return nil, err
	}

	// Update QR code last used
	qrCode.LastUsed = &checkInTime
	if err := s.timeRepo.UpdateQRCode(qrCode); err != nil {
		log.Printf("Failed to update LastUsed: %v", err)
	}

	return attendance, nil
}

// Check-out employee
func (s *timeService) CheckOut(req *domain.CheckOutRequest) (*domain.Attendance, error) {
	// Validate QR code
	qrCode, err := s.ValidateQRCode(req.QRCode)
	if err != nil {
		return nil, err
	}

	// Use provided timestamp or current time
	checkOutTime := req.Timestamp
	if checkOutTime.IsZero() {
		checkOutTime = time.Now()
	}

	// Find today's attendance record
	today := time.Now().Truncate(24 * time.Hour)
	attendance, err := s.timeRepo.GetAttendanceByDate(qrCode.EmployeeID, today)
	if err != nil {
		return nil, errors.New("no check-in record found for today")
	}

	if attendance.CheckOut != nil {
		return nil, errors.New("already checked out today")
	}

	// Update check-out time
	attendance.CheckOut = &checkOutTime
	attendance.Location = req.Location
	attendance.DeviceInfo = req.DeviceInfo

	if err := s.timeRepo.UpdateAttendance(attendance); err != nil {
		return nil, err
	}

	// Update QR code last used
	qrCode.LastUsed = &checkOutTime
	if err := s.timeRepo.UpdateQRCode(qrCode); err != nil {
		return nil, err
	}

	return attendance, nil
}

func (s *timeService) GetAttendanceByDate(employeeID uuid.UUID, date time.Time) (*domain.Attendance, error) {
	return s.timeRepo.GetAttendanceByDate(employeeID, date)
}

func (s *timeService) GetAttendanceSummary(employeeID uuid.UUID, month, year int) (map[string]int, error) {
	// TO DO: implement attendance summary logic
	return nil, nil
}

func (s *timeService) ListAttendances(orgID uuid.UUID) ([]domain.Attendance, error) {
	return s.timeRepo.ListAttendances(orgID)
}

func (s *timeService) CreateTimesheet(orgID, employeeID uuid.UUID, req *domain.CreateTimesheetRequest) (*domain.Timesheet, error) {
	today := time.Now().Truncate(24 * time.Hour)
	timesheet := &domain.Timesheet{
		OrganizationID: orgID,
		EmployeeID:     employeeID,
		ProjectID:      req.ProjectID,
		TaskID:         req.TaskID,
		Description:    req.Description,
		Date:           today,
		Hours:          req.Hours,
	}
	if err := s.timeRepo.CreateTimesheet(timesheet); err != nil {
		return nil, err
	}
	return timesheet, nil
}

func (s *timeService) GetTimesheet(id uuid.UUID) (*domain.Timesheet, error) {
	return s.timeRepo.GetTimesheet(id)
}

func (s *timeService) UpdateTimesheet(id uuid.UUID, req *domain.CreateTimesheetRequest) (*domain.Timesheet, error) {
	today := time.Now().Truncate(24 * time.Hour)
	timesheet, err := s.timeRepo.GetTimesheet(id)
	if err != nil {
		return nil, err
	}
	timesheet.ProjectID = req.ProjectID
	timesheet.TaskID = req.TaskID
	timesheet.Description = req.Description
	timesheet.Date = today
	timesheet.Hours = req.Hours
	if err := s.timeRepo.UpdateTimesheet(timesheet); err != nil {
		return nil, err
	}
	return timesheet, nil
}

func (s *timeService) DeleteTimesheet(id uuid.UUID) error {
	return s.timeRepo.DeleteTimesheet(id)
}

func (s *timeService) ListTimesheets(orgID, employeeID uuid.UUID, startDate, endDate time.Time) ([]domain.Timesheet, error) {
	return s.timeRepo.ListTimesheets(orgID, employeeID, startDate, endDate)
}

func (s *timeService) ApproveTimesheet(id, approverID uuid.UUID) error {
	today := time.Now().Truncate(24 * time.Hour)
	timesheet, err := s.timeRepo.GetTimesheet(id)
	if err != nil {
		return err
	}
	timesheet.Status = domain.TimesheetStatusApproved
	timesheet.ApprovedBy = &approverID
	timesheet.ApprovedAt = &today
	return s.timeRepo.UpdateTimesheet(timesheet)
}

func (s *timeService) RejectTimesheet(id, approverID uuid.UUID, reason string) error {
	today := time.Now().Truncate(24 * time.Hour)
	timesheet, err := s.timeRepo.GetTimesheet(id)
	if err != nil {
		return err
	}
	timesheet.Status = domain.TimesheetStatusRejected
	timesheet.ApprovedBy = &approverID
	timesheet.ApprovedAt = &today
	timesheet.Notes = reason
	return s.timeRepo.UpdateTimesheet(timesheet)
}

func (s *timeService) GetEmployeeQRCodes(orgID, employeeID uuid.UUID) ([]domain.QRCode, error) {
	return s.timeRepo.GetEmployeeQRCodes(orgID, employeeID)
}
