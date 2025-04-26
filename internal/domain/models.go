// internal/domain/models.go
package domain

import (
	"time"

	"github.com/google/uuid"
)

type Base struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// Attendance tracks employee check-in and check-out
type Attendance struct {
	Base
	OrganizationID uuid.UUID  `json:"organization_id" gorm:"type:uuid;not null"`
	EmployeeID     uuid.UUID  `json:"employee_id" gorm:"type:uuid;not null"`
	CheckIn        *time.Time `json:"check_in"`
	CheckOut       *time.Time `json:"check_out"`
	Date           time.Time  `json:"date" gorm:"not null;type:date"`
	Status         string     `json:"status" gorm:"default:'present'"`
	WorkMode       string     `json:"work_mode" gorm:"default:'office'"`
	Location       string     `json:"location"`
	DeviceInfo     string     `json:"device_info"`
}

// Timesheet records work hours on projects/tasks
type Timesheet struct {
	Base
	OrganizationID uuid.UUID  `json:"organization_id" gorm:"type:uuid;not null"`
	EmployeeID     uuid.UUID  `json:"employee_id" gorm:"type:uuid;not null"`
	ProjectID      *uuid.UUID `json:"project_id,omitempty" gorm:"type:uuid"`
	TaskID         *uuid.UUID `json:"task_id,omitempty" gorm:"type:uuid"`
	Description    string     `json:"description"`
	Date           time.Time  `json:"date" gorm:"not null;type:date"`
	Hours          float64    `json:"hours" gorm:"not null;type:decimal(5,2)"`
	Status         string     `json:"status" gorm:"default:'pending'"`
	Notes          string     `json:"notes"`
	ApprovedBy     *uuid.UUID `json:"approved_by,omitempty" gorm:"type:uuid"`
	ApprovedAt     *time.Time `json:"approved_at"`
}

// QRCode for employee check-in/check-out
type QRCode struct {
	Base
	OrganizationID uuid.UUID  `json:"organization_id" gorm:"type:uuid;not null"`
	EmployeeID     uuid.UUID  `json:"employee_id" gorm:"type:uuid;not null"`
	Code           string     `json:"code" gorm:"unique;not null"`
	ExpiryDate     *time.Time `json:"expiry_date"`
	IsActive       bool       `json:"is_active" gorm:"default:true"`
	LastUsed       *time.Time `json:"last_used"`
}

// Request/Response types
type CheckInRequest struct {
	QRCode     string    `json:"qr_code" binding:"required"`
	Location   string    `json:"location"`
	DeviceInfo string    `json:"device_info"`
	WorkMode   string    `json:"work_mode" binding:"required"`
	Timestamp  time.Time `json:"timestamp"`
}

type CheckOutRequest struct {
	QRCode     string    `json:"qr_code" binding:"required"`
	Location   string    `json:"location"`
	DeviceInfo string    `json:"device_info"`
	Timestamp  time.Time `json:"timestamp"`
}

type CreateTimesheetRequest struct {
	ProjectID   *uuid.UUID `json:"project_id"`
	TaskID      *uuid.UUID `json:"task_id"`
	Description string     `json:"description" binding:"required"`
	Date        string     `json:"date" binding:"required"`
	Hours       float64    `json:"hours" binding:"required,min=0.1,max=24"`
	Notes       string     `json:"notes"`
}

type TimesheetResponse struct {
	ID          uuid.UUID  `json:"id"`
	EmployeeID  uuid.UUID  `json:"employee_id"`
	ProjectID   *uuid.UUID `json:"project_id,omitempty"`
	TaskID      *uuid.UUID `json:"task_id,omitempty"`
	ProjectName string     `json:"project_name,omitempty"`
	TaskName    string     `json:"task_name,omitempty"`
	Description string     `json:"description"`
	Date        time.Time  `json:"date"`
	Hours       float64    `json:"hours"`
	Status      string     `json:"status"`
	Notes       string     `json:"notes"`
	ApprovedBy  *uuid.UUID `json:"approved_by,omitempty"`
	ApprovedAt  *time.Time `json:"approved_at"`
}

type GenerateQRRequest struct {
	EmployeeID uuid.UUID `json:"employee_id" binding:"required"`
	ExpiryDays int       `json:"expiry_days"`
}

// Constants
const (
	AttendanceStatusPresent = "present"
	AttendanceStatusAbsent  = "absent"
	AttendanceStatusLate    = "late"
	AttendanceStatusHalfDay = "half-day"

	WorkModeOffice = "office"
	WorkModeRemote = "remote"
	WorkModeHybrid = "hybrid"

	TimesheetStatusPending  = "pending"
	TimesheetStatusApproved = "approved"
	TimesheetStatusRejected = "rejected"
)
