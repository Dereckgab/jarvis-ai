package http

import (
	"context"
	"strconv"

	"jarvis/config"
	"jarvis/internal/application/service"
	"jarvis/internal/interfaces/http/dto"
	appErrors "jarvis/pkg/errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// SystemInfoHandler handles HTTP requests related to system information.
type SystemInfoHandler struct {
	systemInfoService service.SystemInfoService
	cfg               *config.Config
}

// NewSystemInfoHandler creates a new SystemInfoHandler.
func NewSystemInfoHandler(sis service.SystemInfoService, cfg *config.Config) *SystemInfoHandler {
	return &SystemInfoHandler{
		systemInfoService: sis,
		cfg:               cfg,
	}
}

// CollectSystemInfo triggers the collection and saving of current system information.
// @Summary Collect system information
// @Description Collects current system hardware and performance information for the authenticated user
// @Tags SystemInfo
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} dto.SystemInfoResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /system-info/collect [post]
func (h *SystemInfoHandler) CollectSystemInfo(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), h.cfg.App.ReadTimeout)
	defer cancel()

	userID := c.Locals("userID").(uuid.UUID)

	info, err := h.systemInfoService.CollectAndSaveSystemInfo(ctx, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: appErrors.ErrInternalServerError.Error(), Details: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(dto.SystemInfoResponse{
		ID:                    info.ID,
		UserID:                info.UserID,
		Timestamp:             info.Timestamp,
		CPUName:               info.CPUName,
		CPUCores:              info.CPUCores,
		CPUThreads:            info.CPUThreads,
		CPUFrequency:          info.CPUFrequency,
		CPUPercent:            info.CPUPercent,
		TotalMemoryGB:         info.TotalMemoryGB,
		UsedMemoryGB:          info.UsedMemoryGB,
		MemoryPercent:         info.MemoryPercent,
		TotalDiskGB:           info.TotalDiskGB,
		UsedDiskGB:            info.UsedDiskGB,
		DiskPercent:           info.DiskPercent,
		GPUName:               info.GPUName,
		GPUTemperatureC:       info.GPUTemperatureC,
		GPUUtilizationPercent: info.GPUUtilizationPercent,
		BytesSentSec:          info.BytesSentSec,
		BytesRecvSec:          info.BytesRecvSec,
		OSPlatform:            info.OSPlatform,
		OSFamily:              info.OSFamily,
		OSVersion:             info.OSVersion,
		CreatedAt:             info.CreatedAt,
		UpdatedAt:             info.UpdatedAt,
	})
}

// GetLatestSystemInfo retrieves the latest system information for the authenticated user.
// @Summary Get latest system information
// @Description Retrieves the most recent system hardware and performance information for the authenticated user
// @Tags SystemInfo
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} dto.SystemInfoResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /system-info/latest [get]
func (h *SystemInfoHandler) GetLatestSystemInfo(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), h.cfg.App.ReadTimeout)
	defer cancel()

	userID := c.Locals("userID").(uuid.UUID)

	info, err := h.systemInfoService.GetLatestSystemInfo(ctx, userID)
	if err != nil {
		// No data yet — collect on the fly so the dashboard works on first load
		info, err = h.systemInfoService.CollectAndSaveSystemInfo(ctx, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: appErrors.ErrInternalServerError.Error(), Details: err.Error()})
		}
	}

	return c.Status(fiber.StatusOK).JSON(dto.SystemInfoResponse{
		ID:                    info.ID,
		UserID:                info.UserID,
		Timestamp:             info.Timestamp,
		CPUName:               info.CPUName,
		CPUCores:              info.CPUCores,
		CPUThreads:            info.CPUThreads,
		CPUFrequency:          info.CPUFrequency,
		CPUPercent:            info.CPUPercent,
		TotalMemoryGB:         info.TotalMemoryGB,
		UsedMemoryGB:          info.UsedMemoryGB,
		MemoryPercent:         info.MemoryPercent,
		TotalDiskGB:           info.TotalDiskGB,
		UsedDiskGB:            info.UsedDiskGB,
		DiskPercent:           info.DiskPercent,
		GPUName:               info.GPUName,
		GPUTemperatureC:       info.GPUTemperatureC,
		GPUUtilizationPercent: info.GPUUtilizationPercent,
		BytesSentSec:          info.BytesSentSec,
		BytesRecvSec:          info.BytesRecvSec,
		OSPlatform:            info.OSPlatform,
		OSFamily:              info.OSFamily,
		OSVersion:             info.OSVersion,
		CreatedAt:             info.CreatedAt,
		UpdatedAt:             info.UpdatedAt,
	})
}

// GetSystemInfoHistory retrieves a paginated history of system information for the authenticated user.
// @Summary Get system information history
// @Description Retrieves a paginated list of historical system hardware and performance information for the authenticated user
// @Tags SystemInfo
// @Security ApiKeyAuth
// @Produce json
// @Param limit query int false "Limit the number of results" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {array} dto.SystemInfoResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /system-info/history [get]
func (h *SystemInfoHandler) GetSystemInfoHistory(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), h.cfg.App.ReadTimeout)
	defer cancel()

	userID := c.Locals("userID").(uuid.UUID)

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrBadRequest.Error(), Details: "invalid limit parameter"})
	}

	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil || offset < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrBadRequest.Error(), Details: "invalid offset parameter"})
	}

	infos, err := h.systemInfoService.GetSystemInfoHistory(ctx, userID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: appErrors.ErrInternalServerError.Error(), Details: err.Error()})
	}

	var responses []dto.SystemInfoResponse
	for _, info := range infos {
		responses = append(responses, dto.SystemInfoResponse{
			ID:                    info.ID,
			UserID:                info.UserID,
			Timestamp:             info.Timestamp,
			CPUName:               info.CPUName,
			CPUCores:              info.CPUCores,
			CPUThreads:            info.CPUThreads,
			CPUFrequency:          info.CPUFrequency,
			CPUPercent:            info.CPUPercent,
			TotalMemoryGB:         info.TotalMemoryGB,
			UsedMemoryGB:          info.UsedMemoryGB,
			MemoryPercent:         info.MemoryPercent,
			TotalDiskGB:           info.TotalDiskGB,
			UsedDiskGB:            info.UsedDiskGB,
			DiskPercent:           info.DiskPercent,
			GPUName:               info.GPUName,
			GPUTemperatureC:       info.GPUTemperatureC,
			GPUUtilizationPercent: info.GPUUtilizationPercent,
			BytesSentSec:          info.BytesSentSec,
			BytesRecvSec:          info.BytesRecvSec,
			OSPlatform:            info.OSPlatform,
			OSFamily:              info.OSFamily,
			OSVersion:             info.OSVersion,
			CreatedAt:             info.CreatedAt,
			UpdatedAt:             info.UpdatedAt,
		})
	}

	return c.Status(fiber.StatusOK).JSON(responses)
}
