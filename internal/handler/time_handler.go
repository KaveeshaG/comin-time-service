package handler

import (
	"net/http"

	"github.com/Axontik/comin-time-service/internal/domain"
	"github.com/Axontik/comin-time-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TimeHandler struct {
	timeService service.TimeService
}

func NewTimeHandler(timeService service.TimeService) *TimeHandler {
	return &TimeHandler{
		timeService: timeService,
	}
}

// @Summary Generate QR code for employee
// @Tags qr-codes
// @Accept json
// @Produce json
// @Param organization_id path string true "Organization ID"
// @Param request body domain.GenerateQRRequest true "QR code details"
// @Success 201 {object} domain.QRCode
// @Router /organizations/{organization_id}/qr-codes [post]
func (h *TimeHandler) GenerateQRCode(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("organization_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
		return
	}

	var req domain.GenerateQRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	qrCode, err := h.timeService.GenerateQRCode(orgID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, qrCode)
}

// @Summary Check in employee
// @Tags attendance
// @Accept json
// @Produce json
// @Param request body domain.CheckInRequest true "Check-in details"
// @Success 200 {object} domain.Attendance
// @Router /attendance/check-in [post]
func (h *TimeHandler) CheckIn(c *gin.Context) {
	var req domain.CheckInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	attendance, err := h.timeService.CheckIn(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, attendance)
}

// @Summary Check out employee
// @Tags attendance
// @Accept json
// @Produce json
// @Param request body domain.CheckOutRequest true "Check-out details"
// @Success 200 {object} domain.Attendance
// @Router /attendance/check-out [post]
func (h *TimeHandler) CheckOut(c *gin.Context) {
	var req domain.CheckOutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	attendance, err := h.timeService.CheckOut(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, attendance)
}

// @Summary Create timesheet entry
// @Tags timesheets
// @Accept json
// @Produce json
// @Param organization_id path string true "Organization ID"
// @Param employee_id path string true "Employee ID"
// @Param request body domain.CreateTimesheetRequest true "Timesheet details"
// @Success 201 {object} domain.Timesheet
// @Router /organizations/{organization_id}/employees/{employee_id}/timesheets [post]
func (h *TimeHandler) CreateTimesheet(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("organization_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
		return
	}

	employeeID, err := uuid.Parse(c.Param("employee_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee id"})
		return
	}

	var req domain.CreateTimesheetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	timesheet, err := h.timeService.CreateTimesheet(orgID, employeeID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, timesheet)
}

// @Summary Get timesheet by ID
// @Tags timesheets
// @Accept json
// @Produce json
// @Param id path string true "Timesheet ID"
// @Success 200 {object} domain.Timesheet
// @Router /timesheets/{id} [get]
func (h *TimeHandler) GetTimesheet(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid timesheet id"})
		return
	}

	timesheet, err := h.timeService.GetTimesheet(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, timesheet)
}

// @Summary Update timesheet
// @Tags timesheets
// @Accept json
// @Produce json
// @Param id path string true "Timesheet ID"
// @Param request body domain.CreateTimesheetRequest true "Timesheet details"
// @Success 200 {object} domain.Timesheet
// @Router /timesheets/{id} [put]
func (h *TimeHandler) UpdateTimesheet(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid timesheet id"})
		return
	}

	var req domain.CreateTimesheetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	timesheet, err := h.timeService.UpdateTimesheet(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, timesheet)
}

// @Summary Delete timesheet
// @Tags timesheets
// @Accept json
// @Produce json
// @Param id path string true "Timesheet ID"
// @Success 204
// @Router /timesheets/{id} [delete]
func (h *TimeHandler) DeleteTimesheet(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid timesheet id"})
		return
	}

	err = h.timeService.DeleteTimesheet(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// @Summary List attendances
// @Tags attendance
// @Accept json
// @Produce json
// @Param organization_id path string true "Organization ID"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {array} domain.Attendance
// @Router /organizations/{organization_id}/attendance [get]
func (h *TimeHandler) ListAttendances(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("organization_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
		return
	}

	// startDateStr := c.Query("start_date")
	// startDate, err := time.Parse("2006-01-02", startDateStr)
	// if err != nil {
	// 	startDate = time.Time{}
	// }

	// endDateStr := c.Query("end_date")
	// endDate, err := time.Parse("2006-01-02", endDateStr)
	// if err != nil {
	// 	endDate = time.Time{}
	// }

	attendances, err := h.timeService.ListAttendances(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, attendances)
}

// @Summary Get employee QR codes
// @Tags qr-codes
// @Accept json
// @Produce json
// @Param organization_id path string true "Organization ID"
// @Param employee_id path string true "Employee ID"
// @Success 200 {array} domain.QRCode
// @Router /organizations/{organization_id}/qr-codes/{employee_id} [get]
func (h *TimeHandler) GetEmployeeQRCodes(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("organization_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
		return
	}

	employeeID, err := uuid.Parse(c.Param("employee_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee id"})
		return
	}

	qrCodes, err := h.timeService.GetEmployeeQRCodes(orgID, employeeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, qrCodes)
}
